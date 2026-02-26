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

func TestFooterActionCollapseUsesKeyOnlyThenDropsFromEnd(t *testing.T) {
	m := New(
		Actions(
			Action{Key: "tab", Label: "focus", Variant: ActionPrimary, Enabled: true},
			Action{Key: "enter", Label: "submit", Variant: ActionNormal, Enabled: true},
		),
	)
	m.SetSize(7, 1)
	view := core.ViewString(m.View())
	if strings.Contains(view, "focus") || strings.Contains(view, "submit") {
		t.Fatal("expected labels to drop before truncation from end")
	}
	if !strings.Contains(view, "tab") {
		t.Fatal("expected first key to remain visible")
	}
	if strings.Contains(view, "enter") {
		t.Fatal("expected trailing key to be dropped from end when width is tight")
	}
}

func TestFooterHarnessThreeChipLayout(t *testing.T) {
	m := New(
		LeftAction(Action{Key: "/dialog", Label: "custom", Variant: ActionMuted, Enabled: true}),
		Actions(Action{Key: "/theme", Label: "picker", Variant: ActionPrimary, Enabled: true}),
		RightAction(Action{Key: "/page", Label: "swap", Variant: ActionNormal, Enabled: true}),
	)
	m.SetSize(120, 1)
	view := core.ViewString(m.View())
	if !strings.Contains(view, "/dialog") {
		t.Fatal("expected left chip to render")
	}
	if !strings.Contains(view, "/theme") {
		t.Fatal("expected middle chip to render")
	}
	if !strings.Contains(view, "/page") {
		t.Fatalf("expected right chip to render, got %q", view)
	}
	if strings.Index(view, "/dialog") > strings.Index(view, "/theme") || strings.Index(view, "/theme") > strings.Index(view, "/page") {
		t.Fatal("expected left -> middle -> right segment order")
	}
}
