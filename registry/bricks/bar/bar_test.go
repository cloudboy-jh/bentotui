package bar

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
)

func TestCompactCardsAndPriorityOverflow(t *testing.T) {
	b := New(
		FooterAnchored(),
		Left("app"),
		Cards(
			Card{Command: "r", Label: "refresh", Enabled: true, Priority: 4},
			Card{Command: "s", Label: "save", Enabled: true, Priority: 3},
			Card{Command: "q", Label: "quit", Enabled: true, Priority: 2},
			Card{Command: "x", Label: "debug", Enabled: true, Priority: 1},
		),
		CompactCards(),
	)
	b.SetSize(24, 1)

	out := ansi.Strip(viewString(b.View()))
	if !strings.Contains(out, "q") {
		t.Fatalf("expected high priority card retained, got %q", out)
	}
	if strings.Contains(out, "debug") {
		t.Fatalf("expected lowest priority label dropped first, got %q", out)
	}
	if lipgloss.Width(out) != 24 {
		t.Fatalf("expected output width 24, got %d", lipgloss.Width(out))
	}
}

func TestAnchoredIgnoredForNonFooterRole(t *testing.T) {
	b := New(RoleTopBar(), Left("top"), Right("meta"))
	b.SetAnchored(true)
	b.SetSize(20, 1)
	out := ansi.Strip(viewString(b.View()))
	if !strings.Contains(out, "top") || !strings.Contains(out, "meta") {
		t.Fatalf("unexpected top bar render: %q", out)
	}
}

func TestStatusPillRenderedAsSingleUnit(t *testing.T) {
	b := New(RoleTopBar(), StatusPill("LIVE"), Left("app"))
	b.SetSize(24, 1)
	out := ansi.Strip(viewString(b.View()))
	if !strings.Contains(out, "LIVE") {
		t.Fatalf("expected status pill text, got %q", out)
	}
	if strings.Contains(out, "mode") {
		t.Fatalf("expected no split status label, got %q", out)
	}
}

func TestAnchoredCardStyleModes(t *testing.T) {
	tm := theme.CurrentTheme()
	card := Card{Command: "k", Label: "save", Variant: CardPrimary, Enabled: true}

	plain := renderCard(tm, card, true, true, true)
	chip := renderCard(tm, card, true, true, true)
	mixed := renderCard(tm, card, true, true, true)

	if ansi.Strip(plain) != "k save" {
		t.Fatalf("expected plain style text, got %q", ansi.Strip(plain))
	}
	if ansi.Strip(chip) != "k save" {
		t.Fatalf("expected chip style text, got %q", ansi.Strip(chip))
	}
	if ansi.Strip(mixed) != "k save" {
		t.Fatalf("expected mixed style text, got %q", ansi.Strip(mixed))
	}
	if plain == "k save" {
		t.Fatalf("expected anchored card to be styled, got raw text")
	}

	muted := Card{Command: "q", Label: "quit", Variant: CardMuted, Enabled: true}
	mixedMuted := renderCard(tm, muted, true, true, true)
	plainMuted := renderCard(tm, muted, true, true, true)
	if ansi.Strip(mixedMuted) != "q quit" || ansi.Strip(plainMuted) != "q quit" {
		t.Fatalf("expected anchored muted card text pair, got mixed=%q plain=%q", ansi.Strip(mixedMuted), ansi.Strip(plainMuted))
	}
}

func viewString(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprint(v.Content)
}
