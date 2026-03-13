package engine

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

// Sizable is implemented by components that can render at a given size.
type Sizable interface {
	SetSize(width, height int)
	View() tea.View
}

type staticCell struct {
	content string
}

// Static wraps a plain string as a layout cell.
func Static(s string) Sizable {
	return &staticCell{content: s}
}

func (s *staticCell) SetSize(width, height int) {}

func (s *staticCell) View() tea.View {
	return tea.NewView(s.content)
}

type renderFuncCell struct {
	fn     func(width, height int) string
	width  int
	height int
}

// RenderFunc wraps a function as a layout cell.
func RenderFunc(fn func(width, height int) string) Sizable {
	if fn == nil {
		fn = func(_, _ int) string { return "" }
	}
	return &renderFuncCell{fn: fn}
}

func (r *renderFuncCell) SetSize(width, height int) {
	r.width = width
	r.height = height
}

func (r *renderFuncCell) View() tea.View {
	return tea.NewView(r.fn(r.width, r.height))
}

func ViewString(v tea.View) string {
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
