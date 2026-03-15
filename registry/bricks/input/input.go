// Brick: Input:
// +-----------------------------------+
// | accent | typed text / placeholder |
// +-----------------------------------+
// Single-line text input.
// Package input provides a themed text input widget wrapping bubbles/textinput.
// Copy this file into your project: bento add input
//
// Styles are updated from theme.CurrentTheme() on every View() call so that
// live theme switching works without any SetTheme() propagation.
// Dependencies:
//   - charm.land/bubbletea/v2
//   - charm.land/bubbles/v2
//   - charm.land/lipgloss/v2
//   - github.com/cloudboy-jh/bentotui/theme
//   - github.com/cloudboy-jh/bentotui/theme/styles
package input

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

// Model is a themed single-line text input.
type Model struct {
	input  textinput.Model
	width  int
	height int
}

// New creates a focused, empty input.
func New() *Model {
	ti := textinput.New()
	m := &Model{input: ti}
	return m
}

// SetValue sets the input text.
func (m *Model) SetValue(v string) { m.input.SetValue(v) }

// Value returns the current text.
func (m *Model) Value() string { return m.input.Value() }

// Focus focuses the input.
func (m *Model) Focus() tea.Cmd { return m.input.Focus() }

// Blur removes focus.
func (m *Model) Blur() { m.input.Blur() }

// SetPlaceholder sets the placeholder text.
func (m *Model) SetPlaceholder(s string) { m.input.Placeholder = s }

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.input.Update(msg)
	m.input = updated
	return m, cmd
}

func (m *Model) View() tea.View {
	// Update styles every render so theme switching is instant.
	m.input.SetStyles(styles.New(theme.CurrentTheme()).InputStyles())
	return tea.NewView(m.input.View())
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.SetWidth(width)
}

func (m *Model) GetSize() (int, int) { return m.width, m.height }
