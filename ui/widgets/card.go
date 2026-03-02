package widgets

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
)

// Card displays a key/value badge, typically for keybindings.
type Card struct {
	box
	command string
	label   string
	theme   theme.Theme
}

// NewCard creates a card with command and label.
func NewCard(command, label string) *Card {
	c := &Card{
		box:     *newBox(),
		command: command,
		label:   label,
		theme:   theme.CurrentTheme(),
	}
	c.updateStyle()
	return c
}

func (c *Card) updateStyle() {
	c.bg = lipgloss.Color(c.theme.Surface.Elevated)
	c.fg = lipgloss.Color(c.theme.Text.Primary)
}

func (c *Card) Init() tea.Cmd { return nil }

func (c *Card) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

// View renders the card as "command label".
func (c *Card) View() tea.View {
	sys := lipgloss.NewStyle().
		Background(lipgloss.Color(c.theme.Surface.Elevated))
	cmdStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(c.theme.Text.Accent))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(c.theme.Text.Muted))

	content := cmdStyle.Render(c.command) + " " + labelStyle.Render(c.label)

	return tea.NewView(sys.Render(content))
}

func (c *Card) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *Card) GetSize() (int, int) {
	return c.width, c.height
}

// SetTheme updates the theme.
func (c *Card) SetTheme(t theme.Theme) {
	c.theme = t
}

// HeightConstraint returns FixedHeight(1) - card renders on one line.
func (c *Card) HeightConstraint(width int) HeightConstraint {
	return FixedHeight(1)
}

var _ core.Component = (*Card)(nil)
var _ core.Sizeable = (*Card)(nil)
var _ HeightConstraintAware = (*Card)(nil)
