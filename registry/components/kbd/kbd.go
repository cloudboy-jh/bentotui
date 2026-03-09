// Package kbd provides a keyboard shortcut display pair.
// Copy this file into your project: bento add kbd
package kbd

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Model struct {
	command string
	label   string
	active  bool
	variant string
}

func New(command, label string) *Model {
	return &Model{command: command, label: label, active: true, variant: "normal"}
}

func (m *Model) SetCommand(v string)                 { m.command = v }
func (m *Model) SetLabel(v string)                   { m.label = v }
func (m *Model) SetActive(v bool)                    { m.active = v }
func (m *Model) SetVariant(v string)                 { m.variant = strings.TrimSpace(v) }
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	sys := styles.New(t)
	cmd := sys.FooterCardCommand(m.variant, m.active).Render(m.command)
	if strings.TrimSpace(m.label) == "" {
		return tea.NewView(cmd)
	}
	lbl := sys.FooterCardLabel(m.variant, m.active).Render(m.label)
	return tea.NewView(cmd + " " + lbl)
}
