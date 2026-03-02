package widgets

import (
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

// VerticalStack arranges multiple widgets vertically.
type VerticalStack struct {
	children []core.Component
	width    int
	height   int
}

// NewVerticalStack creates a new vertical stack.
func NewVerticalStack(children ...core.Component) *VerticalStack {
	return &VerticalStack{children: children}
}

func (v *VerticalStack) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, child := range v.children {
		cmd := child.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (v *VerticalStack) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for i, child := range v.children {
		updated, cmd := child.Update(msg)
		if next, ok := updated.(core.Component); ok {
			v.children[i] = next
		}
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return v, tea.Batch(cmds...)
}

func (v *VerticalStack) View() tea.View {
	views := make([]string, len(v.children))
	for i, child := range v.children {
		views[i] = core.ViewString(child.View())
	}
	result := ""
	for i, view := range views {
		if i > 0 {
			result += "\n"
		}
		result += view
	}
	return tea.NewView(result)
}

func (v *VerticalStack) SetSize(width, height int) {
	v.width = width
	v.height = height

	allocations := v.allocateHeights(height, width)
	for i, child := range v.children {
		if sizeable, ok := child.(core.Sizeable); ok {
			sizeable.SetSize(width, allocations[i])
		}
	}
}

func (v *VerticalStack) GetSize() (int, int) {
	return v.width, v.height
}

func (v *VerticalStack) allocateHeights(totalH, width int) []int {
	n := len(v.children)
	if n == 0 {
		return nil
	}

	allocations := make([]int, n)
	constraints := make([]HeightConstraint, n)
	for i, child := range v.children {
		if aware, ok := child.(HeightConstraintAware); ok {
			constraints[i] = aware.HeightConstraint(width)
		} else {
			constraints[i] = FlexHeight(1)
		}
	}

	remaining := totalH
	flexTotal := 0
	for i, c := range constraints {
		switch c.Type {
		case HeightFixed:
			allocations[i] = c.Value
			remaining -= c.Value
		case HeightMin:
			allocations[i] = c.Value
			remaining -= c.Value
		case HeightMax, HeightFlex:
			flexTotal += c.Value
		}
	}

	if remaining <= 0 || flexTotal == 0 {
		return allocations
	}

	for i, c := range constraints {
		if c.Type == HeightFlex {
			allocations[i] = remaining * c.Value / flexTotal
		} else if c.Type == HeightMax && allocations[i] == 0 {
			share := remaining / flexTotal
			if share > c.Value {
				share = c.Value
			}
			allocations[i] = share
		}
	}

	return allocations
}

var _ core.Component = (*VerticalStack)(nil)
var _ core.Sizeable = (*VerticalStack)(nil)
