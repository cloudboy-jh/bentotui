// Brick: KBD:
// +-----------------------------+
// | cmd  label                  |
// +-----------------------------+
// Compact keybind hint pair.
// Package kbd provides a keyboard shortcut display pair.
// Copy this file into your project: bento add kbd
package kbd

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

type Model struct {
	command string
	label   string
	active  bool
	variant string
	width   int
	height  int
}

func New(command, label string) *Model {
	return &Model{command: command, label: label, active: true, variant: "normal"}
}

func (m *Model) SetCommand(v string)                 { m.command = v }
func (m *Model) SetLabel(v string)                   { m.label = v }
func (m *Model) SetActive(v bool)                    { m.active = v }
func (m *Model) SetVariant(v string)                 { m.variant = strings.TrimSpace(v) }
func (m *Model) SetSize(width, height int)           { m.width, m.height = width, height }
func (m *Model) GetSize() (int, int)                 { return m.width, m.height }
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	sys := styles.New(t)
	cmd := sys.FooterCardCommandAnchored(m.variant, m.active).Render(m.command)
	if strings.TrimSpace(m.label) == "" {
		return tea.NewView(cmd)
	}
	lbl := sys.FooterCardLabelAnchored(m.variant, m.active).Render(m.label)
	return tea.NewView(cmd + " " + lbl)
}
