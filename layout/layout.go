// Package layout provides composable split-pane layout primitives built on
// bubbletea and lipgloss. It has no dependencies on other bentotui packages
// and can be used standalone.
//
// Usage:
//
//	root := layout.Horizontal(
//	    layout.Fixed(30, sidebar),
//	    layout.Flex(1, main),
//	)
//	root.SetSize(termW, termH)
package layout

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Model is the minimum interface that layout children must satisfy.
// It matches tea.Model exactly — any bubbletea model works.
type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}

// Sizeable is implemented by models that accept explicit size allocation.
// Split calls SetSize on children that implement this interface.
type Sizeable interface {
	Model
	SetSize(width, height int)
	GetSize() (width, height int)
}

type kind int

const (
	fixed kind = iota
	flex
)

// Item is a layout child with its sizing rule. Construct with Fixed or Flex.
type Item struct {
	k      kind
	size   int
	weight int
	child  Model
}

// Fixed allocates an exact number of columns (horizontal) or rows (vertical).
func Fixed(size int, child Model) Item {
	if size < 0 {
		size = 0
	}
	return Item{k: fixed, size: size, child: child}
}

// Flex allocates remaining space proportionally by weight.
func Flex(weight int, child Model) Item {
	if weight < 1 {
		weight = 1
	}
	return Item{k: flex, weight: weight, child: child}
}

// Split is a horizontal or vertical layout container.
type Split struct {
	horizontal  bool
	items       []Item
	width       int
	height      int
	gutterColor string
}

// Horizontal creates a left-to-right split.
func Horizontal(items ...Item) *Split {
	return &Split{horizontal: true, items: items}
}

// Vertical creates a top-to-bottom split.
func Vertical(items ...Item) *Split {
	return &Split{horizontal: false, items: items}
}

// WithGutterColor inserts a 1-cell separator between children painted with
// the given hex color string.
func (s *Split) WithGutterColor(color string) *Split {
	s.gutterColor = color
	return s
}

// SetGutterColor updates the gutter color. Returns s for chaining.
func (s *Split) SetGutterColor(color string) *Split {
	s.gutterColor = color
	return s
}

func (s *Split) Init() tea.Cmd { return nil }

func (s *Split) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, len(s.items))
	for i := range s.items {
		updated, cmd := s.items[i].child.Update(msg)
		if next, ok := updated.(Model); ok {
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
			parts = append(parts, viewString(item.child.View()))
		}
		if s.horizontal {
			return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, parts...))
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
			content := sizeWrapped(viewString(item.child.View()), w, h)
			layer := lipgloss.NewLayer(tea.NewView(content).Content).X(x).Y(0).Z(i * 2)
			layers = append(layers, layer)
			x += w
			if gutterSize > 0 && i < len(s.items)-1 {
				gutter := fill(gutterSize, h, s.gutterColor)
				gl := lipgloss.NewLayer(tea.NewView(gutter).Content).X(x).Y(0).Z(i*2 + 1)
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
		content := sizeWrapped(viewString(item.child.View()), w, h)
		layer := lipgloss.NewLayer(tea.NewView(content).Content).X(0).Y(y).Z(i * 2)
		layers = append(layers, layer)
		y += h
		if gutterSize > 0 && i < len(s.items)-1 {
			gutter := fill(w, gutterSize, s.gutterColor)
			gl := lipgloss.NewLayer(tea.NewView(gutter).Content).X(0).Y(y).Z(i*2 + 1)
			layers = append(layers, gl)
			y += gutterSize
		}
	}
	return tea.NewView(lipgloss.NewCanvas(layers...))
}

// SetSize distributes width/height to children according to their allocation.
func (s *Split) SetSize(width, height int) {
	s.width = width
	s.height = height
	allocs := s.allocations()
	for i := range s.items {
		sz, ok := s.items[i].child.(Sizeable)
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

// GetSize returns the last size passed to SetSize.
func (s *Split) GetSize() (width, height int) {
	return s.width, s.height
}

// DebugLayout returns a visual representation of current allocations.
func (s *Split) DebugLayout() string {
	allocs := s.allocations()
	parts := make([]string, len(allocs))
	for i, n := range allocs {
		parts[i] = strings.Repeat("#", max(n, 0))
	}
	return strings.Join(parts, "|")
}

// ── internal helpers ──────────────────────────────────────────────────────────

// viewString extracts a plain string from a tea.View. It checks for a
// Render() string method (lipgloss Canvas), then fmt.Stringer, then
// falls back to fmt.Sprint.
func viewString(v tea.View) string {
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

// fill returns a width×height block of spaces painted with bg.
func fill(width, height int, bg string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	line := strings.Repeat(" ", width)
	rows := make([]string, height)
	for i := range rows {
		rows[i] = line
	}
	block := strings.Join(rows, "\n")
	if bg == "" {
		return block
	}
	return lipgloss.NewStyle().Background(lipgloss.Color(bg)).Render(block)
}

// sizeWrapped forces a string into an exact width×height block so that the
// canvas layer for this region covers every cell in its allocation.
func sizeWrapped(s string, w, h int) string {
	if w <= 0 || h <= 0 {
		return ""
	}
	return lipgloss.NewStyle().Width(w).Height(h).Render(s)
}

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

	gutterTotal := s.gutterCells() * max(0, count-1)
	total = max(0, total-gutterTotal)

	out := make([]int, count)
	remaining := total
	weightSum := 0
	for i, item := range s.items {
		if item.k == fixed {
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
		if item.k != flex {
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
