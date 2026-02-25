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
		Actions(Action{Key: "tab", Label: "focus", Variant: ActionPrimary, Enabled: true}),
	)
	m.SetSize(80, 1)
	view := core.ViewString(m.View())
	if strings.Contains(view, "\n") {
		t.Fatal("expected single-line footer output")
	}
	if lipgloss.Width(view) != 80 {
		t.Fatalf("expected footer width 80, got %d", lipgloss.Width(view))
	}
}

func TestFooterKeepsRightSegmentUnderTightWidth(t *testing.T) {
	m := New(
		Left("very long left context that should truncate first"),
		Right("RIGHT"),
		Actions(
			Action{Key: "tab", Label: "focus", Variant: ActionPrimary, Enabled: true},
			Action{Key: "enter", Label: "run", Variant: ActionNormal, Enabled: true},
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

func TestFooterActionsFallbackToKeyOnly(t *testing.T) {
	m := New(
		Actions(
			Action{Key: "tab", Label: "focus", Variant: ActionPrimary, Enabled: true},
			Action{Key: "enter", Label: "run", Variant: ActionNormal, Enabled: true},
		),
	)
	m.SetSize(12, 1)
	view := core.ViewString(m.View())
	if strings.Contains(view, "focus") || strings.Contains(view, "run") {
		t.Fatal("expected label text to be dropped under tight width")
	}
	if lipgloss.Width(view) != 12 {
		t.Fatalf("expected width 12, got %d", lipgloss.Width(view))
	}
}
