// Brick: Text:
// +-----------------------------+
// | static label text           |
// +-----------------------------+
// Simple themed text label.
// Package text provides a static text display widget.
// Copy this file into your project: bento add text
//
// The widget calls theme.CurrentTheme() in View() — never stores theme state.
// Dependencies:
//   - charm.land/bubbletea/v2
//   - charm.land/lipgloss/v2
//   - github.com/cloudboy-jh/bentotui/theme
package text

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Model displays a static string in the theme's primary text color.
// For colored or styled text, build a lipgloss.Style directly in your app.
type Model struct {
	text   string
	width  int
	height int
}

// New creates a text widget with initial content.
func New(text string) *Model { return &Model{text: text} }

// SetText updates the displayed text.
func (t *Model) SetText(text string) { t.text = text }

func (t *Model) Init() tea.Cmd                           { return nil }
func (t *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }

func (t *Model) View() tea.View {
	th := theme.CurrentTheme()
	style := lipgloss.NewStyle().Foreground(th.Text())
	return tea.NewView(style.Render(t.text))
}

func (t *Model) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Model) GetSize() (int, int) { return t.width, t.height }
