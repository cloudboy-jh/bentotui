package footer

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

func TestFooterRendersSingleLine(t *testing.T) {
	m := New(
		Left("left"),
		Right("right"),
		Cards(Card{Command: "/focus", Label: "focus", Variant: CardPrimary, Enabled: true}),
	)
	m.SetSize(80, 4)
	view := core.ViewString(m.View())
	if strings.Contains(view, "\n") {
		t.Fatal("expected single-line footer output")
	}
	if lipgloss.Width(view) != 80 {
		t.Fatalf("expected footer width 80, got %d", lipgloss.Width(view))
	}
	_, h := m.GetSize()
	if h != 1 {
		t.Fatalf("expected footer height to stay 1, got %d", h)
	}
}

func TestFooterKeepsRightSegmentUnderTightWidth(t *testing.T) {
	m := New(
		Left("very long left context that should truncate first"),
		Right("RIGHT"),
		Cards(
			Card{Command: "/focus", Label: "focus", Variant: CardPrimary, Enabled: true},
			Card{Command: "/run", Label: "run", Variant: CardNormal, Enabled: true},
		),
	)
	m.SetSize(20, 1)
	view := core.ViewString(m.View())
	if !strings.Contains(view, "RIGHT") {
		t.Fatal("expected right segment to be preserved")
	}
	if lipgloss.Width(view) != 20 {
		t.Fatalf("expected width 20, got %d", lipgloss.Width(view))
	}
}

func TestFooterCardCollapseUsesCommandOnlyThenDropsFromEnd(t *testing.T) {
	m := New(
		Cards(
			Card{Command: "/f", Label: "focus", Variant: CardPrimary, Enabled: true},
			Card{Command: "/s", Label: "submit", Variant: CardNormal, Enabled: true},
		),
	)
	m.SetSize(3, 1)
	view := core.ViewString(m.View())
	if strings.Contains(view, "focus") || strings.Contains(view, "submit") {
		t.Fatal("expected labels to drop before truncation from end")
	}
	if !strings.Contains(view, "/f") {
		t.Fatal("expected first card command to remain visible")
	}
	if strings.Contains(view, "/s") {
		t.Fatal("expected trailing card command to be dropped from end when width is tight")
	}
}

func TestFooterHarnessThreeCardLayout(t *testing.T) {
	m := New(
		LeftCard(Card{Command: "/pr", Label: "pull requests", Variant: CardMuted, Enabled: true}),
		Cards(Card{Command: "/issue", Label: "issues", Variant: CardPrimary, Enabled: true}),
		RightCard(Card{Command: "/branch", Label: "branches", Variant: CardNormal, Enabled: true}),
	)
	m.SetSize(120, 1)
	view := core.ViewString(m.View())
	if !strings.Contains(view, "/pr") {
		t.Fatal("expected left card to render")
	}
	if !strings.Contains(view, "/issue") {
		t.Fatal("expected middle card to render")
	}
	if !strings.Contains(view, "/branch") {
		t.Fatalf("expected right card to render, got %q", view)
	}
	if strings.Index(view, "/pr") > strings.Index(view, "/issue") || strings.Index(view, "/issue") > strings.Index(view, "/branch") {
		t.Fatal("expected left -> middle -> right segment order")
	}
}
