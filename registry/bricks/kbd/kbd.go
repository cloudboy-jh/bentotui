// Brick: KBD
// +-----------------------------+
// | cmd  label                  |
// +-----------------------------+
// Compact keybind hint pair.
// Copy this file into your project: bento add kbd
package kbd

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Model struct {
	command string
	label   string
	active  bool
	variant string
	width   int
	height  int
	theme   theme.Theme // nil = use theme.CurrentTheme()
}

func New(command, label string) *Model {
	return &Model{command: command, label: label, active: true, variant: "normal"}
}

func (m *Model) SetCommand(v string)                 { m.command = v }
func (m *Model) SetLabel(v string)                   { m.label = v }
func (m *Model) SetActive(v bool)                    { m.active = v }
func (m *Model) SetVariant(v string)                 { m.variant = strings.TrimSpace(v) }
func (m *Model) SetSize(w, h int)                    { m.width, m.height = w, h }
func (m *Model) GetSize() (int, int)                 { return m.width, m.height }
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *Model) SetTheme(t theme.Theme)              { m.theme = t }

func (m *Model) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

func (m *Model) View() tea.View {
	t := m.activeTheme()

	var cmdFG, cmdBG, lblFG color.Color
	switch {
	case !m.active:
		cmdFG = t.FooterMuted()
		cmdBG = t.FooterBG()
		lblFG = t.FooterMuted()
	case m.variant == "danger":
		cmdFG = t.TextInverse()
		cmdBG = t.Error()
		lblFG = t.FooterMuted()
	case m.variant == "primary":
		cmdFG = t.SelectionFG()
		cmdBG = t.SelectionBG()
		lblFG = t.FooterFG()
	default:
		cmdFG = t.FooterFG()
		cmdBG = t.FooterBG()
		lblFG = t.FooterMuted()
	}

	cmdStyle := lipgloss.NewStyle().Bold(true).Foreground(cmdFG).Background(cmdBG)
	cmd := cmdStyle.Render(m.command)
	if strings.TrimSpace(m.label) == "" {
		return tea.NewView(cmd)
	}
	lblStyle := lipgloss.NewStyle().Foreground(lblFG)
	lbl := lblStyle.Render(m.label)
	return tea.NewView(cmd + " " + lbl)
}
