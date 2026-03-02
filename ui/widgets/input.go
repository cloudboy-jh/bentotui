package widgets

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

// Input is a text input field with theming.
type Input struct {
	box
	input textinput.Model
	theme theme.Theme
}

// NewInput creates an input field.
func NewInput() *Input {
	ti := textinput.New()
	t := theme.CurrentTheme()

	i := &Input{
		box:   *newBox(),
		input: ti,
		theme: t,
	}
	i.updateStyle()
	return i
}

func (i *Input) updateStyle() {
	i.input.SetStyles(styles.New(i.theme).InputStyles())
}

// SetValue sets the input text.
func (i *Input) SetValue(v string) {
	i.input.SetValue(v)
}

// Value returns the current text.
func (i *Input) Value() string {
	return i.input.Value()
}

// Focus focuses the input.
func (i *Input) Focus() tea.Cmd {
	return i.input.Focus()
}

// Blur removes focus.
func (i *Input) Blur() {
	i.input.Blur()
}

// SetTheme updates the theme.
func (i *Input) SetTheme(t theme.Theme) {
	i.theme = t
	i.updateStyle()
}

func (i *Input) Init() tea.Cmd {
	return nil
}

func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := i.input.Update(msg)
	i.input = updated
	return i, cmd
}

func (i *Input) View() tea.View {
	return tea.NewView(i.input.View())
}

func (i *Input) SetSize(width, height int) {
	i.width = width
	i.height = height
	i.input.SetWidth(width)
}

func (i *Input) GetSize() (int, int) {
	return i.width, i.height
}

// HeightConstraint returns FixedHeight(1) - input needs exactly 1 line.
func (i *Input) HeightConstraint(width int) HeightConstraint {
	return FixedHeight(1)
}

var _ core.Component = (*Input)(nil)
var _ core.Sizeable = (*Input)(nil)
var _ HeightConstraintAware = (*Input)(nil)
