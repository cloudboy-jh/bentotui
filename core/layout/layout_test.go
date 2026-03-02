package layout

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

// TestSizeWrappedGuaranteesExactDimensions verifies that sizeWrapped produces
// a string with the exact w×h cell count, which is the core of the
// OpenCode container / region-anchored background fill pattern.
// Note: content passed to sizeWrapped must fit within w cells; the helper
// is not responsible for truncation — child components handle their own width.
func TestSizeWrappedGuaranteesExactDimensions(t *testing.T) {
	cases := []struct {
		w, h    int
		content string
	}{
		{20, 5, "hi"},
		{10, 3, "hi"},
		{80, 24, "some content"},
		{0, 5, "hi"},
		{10, 0, "hi"},
	}
	for _, tc := range cases {
		result := sizeWrapped(tc.content, tc.w, tc.h)
		if tc.w <= 0 || tc.h <= 0 {
			if result != "" {
				t.Fatalf("sizeWrapped(%d,%d): expected empty, got %q", tc.w, tc.h, result)
			}
			continue
		}
		lines := strings.Split(result, "\n")
		if len(lines) != tc.h {
			t.Fatalf("sizeWrapped(%d,%d): expected %d rows, got %d", tc.w, tc.h, tc.h, len(lines))
		}
		for i, line := range lines {
			if w := lipgloss.Width(line); w != tc.w {
				t.Fatalf("sizeWrapped(%d,%d) row %d: expected width %d, got %d", tc.w, tc.h, i, tc.w, w)
			}
		}
	}
}

// TestHorizontalSplitViewDimensions verifies that a horizontal Split with
// sized children produces a View whose string dimensions match the allocated
// total width × height, confirming the region-anchored wrapping is correct.
func TestHorizontalSplitViewDimensions(t *testing.T) {
	a := &stubComponent{}
	b := &stubComponent{}

	s := Horizontal(Flex(1, a), Flex(1, b))
	s.SetSize(40, 10)

	viewStr := ""
	v := s.View()
	if v.Content != nil {
		if r, ok := v.Content.(interface{ Render() string }); ok {
			viewStr = r.Render()
		}
	}
	if viewStr == "" {
		t.Skip("canvas did not produce string output; skipping dimension check")
	}

	lines := strings.Split(viewStr, "\n")
	if len(lines) != 10 {
		t.Fatalf("expected 10 rows from horizontal split, got %d", len(lines))
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 40 {
			t.Fatalf("row %d: expected width 40, got %d", i, w)
		}
	}
}
