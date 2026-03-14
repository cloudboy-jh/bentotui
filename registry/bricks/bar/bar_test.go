package bar

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
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
