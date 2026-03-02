package layout

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
)

type Kind int

const (
	fixed Kind = iota
	flex
)

type Item struct {
	kind   Kind
	size   int
	weight int
	child  core.Component
}

func Fixed(size int, child core.Component) Item {
	if size < 0 {
		size = 0
	}
	return Item{kind: fixed, size: size, child: child}
}

func Flex(weight int, child core.Component) Item {
	if weight < 1 {
		weight = 1
	}
	return Item{kind: flex, weight: weight, child: child}
}

type Split struct {
	horizontal  bool
	items       []Item
	width       int
	height      int
	gutterColor string // if non-empty, insert 1-cell gutters between children
}

func Horizontal(items ...Item) *Split {
	return &Split{horizontal: true, items: items}
}

func Vertical(items ...Item) *Split {
	return &Split{horizontal: false, items: items}
}

// WithGutterColor inserts a 1-cell separator between children painted with
// the given hex color string. The gutter cells are subtracted from the total
// before flex allocations are computed, so children still fill exactly their
// allocated region.
func (s *Split) WithGutterColor(color string) *Split {
	s.gutterColor = color
	return s
}

func (s *Split) Init() tea.Cmd { return nil }

func (s *Split) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, len(s.items))
	for i := range s.items {
		updated, cmd := s.items[i].child.Update(msg)
		if next, ok := updated.(core.Component); ok {
			s.items[i].child = next
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	if len(cmds) == 0 {
		return s, nil
	}
	return s, tea.Batch(cmds...)
}

func (s *Split) View() tea.View {
	if len(s.items) == 0 {
		return tea.NewView("")
	}

	if s.width <= 0 || s.height <= 0 {
		parts := make([]string, 0, len(s.items))
		for _, item := range s.items {
			parts = append(parts, core.ViewString(item.child.View()))
		}
		if s.horizontal {
			return tea.NewView(joinHorizontal(parts...))
		}
		return tea.NewView(lipgloss.JoinVertical(lipgloss.Top, parts...))
	}

	allocs := s.allocations()
	gutterSize := s.gutterCells()
	layers := make([]*lipgloss.Layer, 0, len(s.items)*2)
	if s.horizontal {
		x := 0
		for i, item := range s.items {
			w := max(0, allocs[i])
			h := s.height
			// Wrap each child in an exact-size style so that the layer occupies
			// its full allocated w×h region in the canvas. Without this, children
			// that render fewer rows than their allocation leave gaps that the
			// canvas background bleeds through (OpenCode container pattern).
			content := sizeWrapped(core.ViewString(item.child.View()), w, h)
			layer := lipgloss.NewLayer(tea.NewView(content).Content).
				X(x).
				Y(0).
				Z(i * 2)
			layers = append(layers, layer)
			x += w
			// Insert gutter between children (not after the last one).
			if gutterSize > 0 && i < len(s.items)-1 {
				gutter := primitives.Fill(gutterSize, h, s.gutterColor)
				gl := lipgloss.NewLayer(tea.NewView(gutter).Content).
					X(x).
					Y(0).
					Z(i*2 + 1)
				layers = append(layers, gl)
				x += gutterSize
			}
		}
		return tea.NewView(lipgloss.NewCanvas(layers...))
	}

	y := 0
	for i, item := range s.items {
		w := s.width
		h := max(0, allocs[i])
		// Same region-anchored wrapping for vertical splits.
		content := sizeWrapped(core.ViewString(item.child.View()), w, h)
		layer := lipgloss.NewLayer(tea.NewView(content).Content).
			X(0).
			Y(y).
			Z(i * 2)
		layers = append(layers, layer)
		y += h
		// Insert gutter between children (not after the last one).
		if gutterSize > 0 && i < len(s.items)-1 {
			gutter := primitives.Fill(w, gutterSize, s.gutterColor)
			gl := lipgloss.NewLayer(tea.NewView(gutter).Content).
				X(0).
				Y(y).
				Z(i*2 + 1)
			layers = append(layers, gl)
			y += gutterSize
		}
	}

	return tea.NewView(lipgloss.NewCanvas(layers...))
}

// sizeWrapped forces a string into an exact width×height block so that the
// canvas layer for this region covers every cell in its allocation. This is
// the OpenCode container pattern applied to BentoTUI's canvas compositor.
// The child component is responsible for painting its own background color;
// this wrapper only guarantees the layer dimensions are exact.
func sizeWrapped(s string, w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	return lipgloss.NewStyle().Width(w).Height(h).Render(s)
}

func (s *Split) SetSize(width, height int) {
	s.width = width
	s.height = height
	allocs := s.allocations() // already gutter-subtracted
	for i := range s.items {
		sz, ok := s.items[i].child.(core.Sizeable)
		if !ok {
			continue
		}
		if s.horizontal {
			sz.SetSize(allocs[i], height)
		} else {
			sz.SetSize(width, allocs[i])
		}
	}
}

func (s *Split) GetSize() (width, height int) {
	return s.width, s.height
}

// gutterCells returns the size of each gutter (1 if gutterColor set, else 0).
func (s *Split) gutterCells() int {
	if s.gutterColor == "" || len(s.items) < 2 {
		return 0
	}
	return 1
}

func (s *Split) allocations() []int {
	count := len(s.items)
	if count == 0 {
		return nil
	}

	total := s.width
	if !s.horizontal {
		total = s.height
	}
	if total < 0 {
		total = 0
	}

	// Subtract gutter cells from the total before distributing to children.
	gutterTotal := s.gutterCells() * max(0, count-1)
	total = max(0, total-gutterTotal)

	out := make([]int, count)
	remaining := total
	weightSum := 0
	for i, item := range s.items {
		if item.kind == fixed {
			w := min(item.size, remaining)
			out[i] = w
			remaining -= w
			continue
		}
		weightSum += item.weight
	}

	if remaining == 0 || weightSum == 0 {
		return out
	}

	assigned := 0
	lastFlex := -1
	for i, item := range s.items {
		if item.kind != flex {
			continue
		}
		lastFlex = i
		w := remaining * item.weight / weightSum
		out[i] = w
		assigned += w
	}

	if lastFlex >= 0 {
		out[lastFlex] += remaining - assigned
	}

	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DebugLayout returns the current split allocations for tests/debugging.
func (s *Split) DebugLayout() string {
	allocs := s.allocations()
	parts := make([]string, len(allocs))
	for i, n := range allocs {
		parts[i] = strings.Repeat("#", max(n, 0))
	}
	return strings.Join(parts, "|")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func joinHorizontal(parts ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
