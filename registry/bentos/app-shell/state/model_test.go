package state

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestModelViewExactDimensions(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	out := viewFromTea(m.View())
	lines := strings.Split(out, "\n")
	if len(lines) != 30 {
		t.Fatalf("expected 30 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got != 100 {
			t.Fatalf("line %d width mismatch: got %d want 100", i, got)
		}
	}
}

func TestModelViewContainsWorkspaceFooterAndCards(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	out := viewFromTea(m.View())
	checks := []string{"Services", "Queue", "Progress", "overview | queue:1", "theme:", "ctrl+k"}
	for _, token := range checks {
		if !strings.Contains(out, token) {
			t.Fatalf("expected token %q in render output", token)
		}
	}
}

func TestModelSmallViewportKeepsFooterBar(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 70, Height: 12})
	out := viewFromTea(m.View())
	lines := strings.Split(out, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[len(lines)-1]) == "" {
		t.Fatalf("expected footer bar right status in small viewport")
	}
}

func viewFromTea(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(v.Content)
}
