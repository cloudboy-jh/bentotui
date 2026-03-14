// Package checkbox provides a themed boolean toggle input.
// Copy this file into your project: bento add checkbox
package checkbox

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Model struct {
	label   string
	checked bool
	focused bool
	width   int
}

func New(label string) *Model {
	return &Model{label: label}
}

func (m *Model) Toggle()              { m.checked = !m.checked }
func (m *Model) Checked() bool        { return m.checked }
func (m *Model) SetChecked(v bool)    { m.checked = v }
func (m *Model) SetLabel(v string)    { m.label = v }
func (m *Model) Focus()               { m.focused = true }
func (m *Model) Blur()                { m.focused = false }
func (m *Model) IsFocused() bool      { return m.focused }
func (m *Model) Init() tea.Cmd        { return nil }
func (m *Model) SetSize(width, _ int) { m.width = width }
func (m *Model) GetSize() (int, int)  { return m.width, 1 }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case " ", "enter":
		m.Toggle()
	}
	return m, nil
}

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	mark := "[ ]"
	if m.checked {
		mark = "[x]"
	}
	content := strings.TrimSpace(mark + " " + m.label)
	if m.width > 0 {
		bg := pick(t.Input.BG, t.Surface.Elevated)
		fg := pick(t.Input.FG, t.Text.Primary)
		if m.focused {
			bg = pick(t.Selection.BG, t.Border.Focus)
			fg = pick(t.Selection.FG, t.Text.Inverse)
		}
		return tea.NewView(styles.Row(bg, fg, m.width, content))
	}
	if m.focused {
		return tea.NewView(styles.New(t).ActionButton(true).Render(content))
	}
	return tea.NewView(content)
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
