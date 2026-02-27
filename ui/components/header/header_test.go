package header

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

func TestHeaderRendersSingleLine(t *testing.T) {
	m := New(
		Left("left"),
		Right("right"),
		Cards(Card{Command: "/focus", Label: "focus", Variant: CardPrimary, Enabled: true}),
	)
	m.SetSize(80, 4)
	view := core.ViewString(m.View())
	if strings.Contains(view, "\n") {
		t.Fatal("expected single-line header output")
	}
	if lipgloss.Width(view) != 80 {
		t.Fatalf("expected header width 80, got %d", lipgloss.Width(view))
	}
	_, h := m.GetSize()
	if h != 1 {
		t.Fatalf("expected header height to stay 1, got %d", h)
	}
}

func TestHeaderKeepsRightSegmentUnderTightWidth(t *testing.T) {
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
