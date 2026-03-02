package widgets

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
)

// Text displays static text content.
type Text struct {
	width  int
	height int
	text   string
	theme  theme.Theme
}

// NewText creates a text widget with initial content.
func NewText(text string) *Text {
	return &Text{text: text, theme: theme.CurrentTheme()}
}

// SetText updates the displayed text.
func (t *Text) SetText(text string) {
	t.text = text
}

func (t *Text) Init() tea.Cmd { return nil }

func (t *Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *Text) View() tea.View {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.theme.Text.Primary))
	return tea.NewView(style.Render(t.text))
}

func (t *Text) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Text) GetSize() (int, int) {
	return t.width, t.height
}

// SetTheme updates the theme.
func (t *Text) SetTheme(theme theme.Theme) {
	t.theme = theme
}

// HeightConstraint returns FixedHeight based on number of lines in text.
func (t *Text) HeightConstraint(width int) HeightConstraint {
	lines := strings.Count(t.text, "\n") + 1
	if lines < 1 {
		lines = 1
	}
	return FixedHeight(lines)
}

var _ core.Component = (*Text)(nil)
var _ core.Sizeable = (*Text)(nil)
var _ HeightConstraintAware = (*Text)(nil)
