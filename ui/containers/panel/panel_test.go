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
