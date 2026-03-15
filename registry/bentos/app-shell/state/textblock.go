package state

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

type textBlock struct {
	text   string
	width  int
	height int
}

func (t *textBlock) Init() tea.Cmd                           { return nil }
func (t *textBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t *textBlock) SetSize(width, height int)               { t.width, t.height = width, height }
func (t *textBlock) SetText(v string)                        { t.text = v }

func (t *textBlock) View() tea.View {
	if t.height <= 0 {
		return tea.NewView("")
	}
	lines := strings.Split(t.text, "\n")
	if len(lines) > t.height {
		lines = lines[:t.height]
	}
	return tea.NewView(strings.Join(lines, "\n"))
}
