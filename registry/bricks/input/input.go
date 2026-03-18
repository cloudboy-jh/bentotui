// Brick: Input
// +-----------------------------------+
// | accent | typed text / placeholder |
// +-----------------------------------+
// Single-line text input.
// Copy this file into your project: bento add input
package input

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Model is a themed single-line text input.
type Model struct {
	input  textinput.Model
	width  int
	height int
	theme  theme.Theme // nil = use theme.CurrentTheme()
}

// New creates a focused, empty input.
func New() *Model {
	ti := textinput.New()
	m := &Model{input: ti}
	return m
}

func (m *Model) SetValue(v string)       { m.input.SetValue(v) }
func (m *Model) Value() string           { return m.input.Value() }
func (m *Model) Focus() tea.Cmd          { return m.input.Focus() }
func (m *Model) Blur()                   { m.input.Blur() }
func (m *Model) IsFocused() bool         { return m.input.Focused() }
func (m *Model) SetPlaceholder(s string) { m.input.Placeholder = s }
func (m *Model) Init() tea.Cmd           { return nil }

// SetTheme updates the theme. Call on ThemeChangedMsg.
func (m *Model) SetTheme(t theme.Theme) {
	m.theme = t
	m.syncStyles()
}

func (m *Model) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

func (m *Model) syncStyles() {
	t := m.activeTheme()
	s := textinput.DefaultStyles(m.input.Focused())
	textStyle := lipgloss.NewStyle().
		Foreground(t.InputFG()).
		Background(t.InputBG())
	placeholderStyle := lipgloss.NewStyle().
		Foreground(t.InputPlaceholder()).
		Background(t.InputBG())
	s.Focused.Text = textStyle
	s.Focused.Placeholder = placeholderStyle
	s.Blurred.Text = textStyle
	s.Blurred.Placeholder = placeholderStyle
	s.Cursor.Color = t.InputCursor()
	m.input.SetStyles(s)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.SetSize(ws.Width, ws.Height)
		return m, nil
	}
	updated, cmd := m.input.Update(msg)
	m.input = updated
	return m, cmd
}

func (m *Model) View() tea.View {
	// Sync styles on every render so theme changes are reflected immediately
	// even if SetTheme wasn't explicitly called.
	m.syncStyles()
	return tea.NewView(m.input.View())
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.SetWidth(width)
}

func (m *Model) GetSize() (int, int) { return m.width, m.height }
