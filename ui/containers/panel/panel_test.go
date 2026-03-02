package panel

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

type staticContent struct {
	text          string
	width, height int
}

func (s *staticContent) Init() tea.Cmd { return nil }

func (s *staticContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return s, nil
}

func (s *staticContent) View() tea.View { return tea.NewView(s.text) }

func (s *staticContent) SetSize(width, height int) {
	s.width = width
	s.height = height
}

func (s *staticContent) GetSize() (int, int) { return s.width, s.height }

func TestPanelRendersExactAssignedBounds(t *testing.T) {
	body := &staticContent{text: "row one\nrow two"}
	p := New(Title("Main"), Content(body))
	p.SetSize(40, 8)

	view := core.ViewString(p.View())
	lines := strings.Split(view, "\n")
	if len(lines) != 8 {
		t.Fatalf("expected 8 rows, got %d", len(lines))
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 40 {
			t.Fatalf("expected row %d width 40, got %d", i, w)
		}
	}
}

// TestPanelFullHeightFillWhenContentIsSparse verifies the OpenCode container
// pattern: a panel whose content is shorter than its allocated height must
// still render the full h rows, guaranteeing region-anchored background
// paint and preventing shell-canvas bleed-through.
func TestPanelFullHeightFillWhenContentIsSparse(t *testing.T) {
	// Content has only 1 line; panel is allocated 20 rows.
	body := &staticContent{text: "only one line"}
	p := New(Title("Sparse"), Content(body))
	p.SetSize(50, 20)

	view := core.ViewString(p.View())
	lines := strings.Split(view, "\n")
	if len(lines) != 20 {
		t.Fatalf("expected 20 rows for full-height fill, got %d", len(lines))
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 50 {
			t.Fatalf("row %d: expected width 50, got %d", i, w)
		}
	}
}

// TestPanelNoTitleFullHeightFill checks the same guarantee for a panel
// without a title bar.
func TestPanelNoTitleFullHeightFill(t *testing.T) {
	body := &staticContent{text: "line a"}
	p := New(Content(body))
	p.SetSize(30, 10)

	view := core.ViewString(p.View())
	lines := strings.Split(view, "\n")
	if len(lines) != 10 {
		t.Fatalf("expected 10 rows, got %d", len(lines))
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 30 {
			t.Fatalf("row %d: expected width 30, got %d", i, w)
		}
	}
}
