package layout

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
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
	horizontal bool
	items      []Item
	width      int
	height     int
}

func Horizontal(items ...Item) *Split {
	return &Split{horizontal: true, items: items}
}

func Vertical(items ...Item) *Split {
	return &Split{horizontal: false, items: items}
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
	parts := make([]string, 0, len(s.items))
	for _, item := range s.items {
		parts = append(parts, core.ViewString(item.child.View()))
	}
	if s.horizontal {
		return tea.NewView(joinHorizontal(parts...))
	}
	return tea.NewView(strings.Join(parts, "\n"))
}

func (s *Split) SetSize(width, height int) {
	s.width = width
	s.height = height
	allocs := s.allocations()
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
