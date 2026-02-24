package layout

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type stubComponent struct {
	width  int
	height int
}

func (s *stubComponent) Init() tea.Cmd                           { return nil }
func (s *stubComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s *stubComponent) View() tea.View                          { return tea.NewView("") }
func (s *stubComponent) SetSize(width, height int) {
	s.width = width
	s.height = height
}
func (s *stubComponent) GetSize() (width, height int) { return s.width, s.height }

func TestHorizontalAllocations(t *testing.T) {
	a := &stubComponent{}
	b := &stubComponent{}
	c := &stubComponent{}

	s := Horizontal(
		Fixed(10, a),
		Flex(1, b),
		Flex(2, c),
	)

	s.SetSize(70, 12)

	if a.width != 10 {
		t.Fatalf("expected fixed width 10, got %d", a.width)
	}
	if b.width != 20 {
		t.Fatalf("expected flex(1) width 20, got %d", b.width)
	}
	if c.width != 40 {
		t.Fatalf("expected flex(2) width 40, got %d", c.width)
	}
}

func TestVerticalAllocations(t *testing.T) {
	a := &stubComponent{}
	b := &stubComponent{}

	s := Vertical(
		Fixed(3, a),
		Flex(1, b),
	)

	s.SetSize(44, 20)

	if a.height != 3 {
		t.Fatalf("expected fixed height 3, got %d", a.height)
	}
	if b.height != 17 {
		t.Fatalf("expected flex height 17, got %d", b.height)
	}
}
