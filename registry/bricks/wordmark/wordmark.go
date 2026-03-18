// Brick: Wordmark:
// +-----------------------------------+
// |            TITLE                  |
// +-----------------------------------+
// Branded heading block.
// Package wordmark provides a themed title/heading block.
// Copy this file into your project: bento add wordmark
package wordmark

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Model struct {
	text   string
	bold   bool
	width  int
	height int
}

func New(text string) *Model {
	return &Model{text: text, bold: true}
}

func (m *Model) SetText(text string)                 { m.text = text }
func (m *Model) SetBold(v bool)                      { m.bold = v }
func (m *Model) SetSize(width, height int)           { m.width, m.height = width, height }
func (m *Model) GetSize() (int, int)                 { return m.width, m.height }
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	style := lipgloss.NewStyle().Foreground(t.TextAccent()).Bold(m.bold)
	return tea.NewView(style.Render(m.text))
}
