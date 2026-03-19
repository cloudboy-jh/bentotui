package filterbar

import (
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"
	"github.com/cloudboy-jh/bentotui/registry/bricks/input"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Model is a copy-and-own composition recipe for a filter row and footer keybind strip.
type Model struct {
	Input  *input.Model
	Footer *bar.Model
}

// New creates a filter-bar recipe with default keybind cards.
func New(appName string, t theme.Theme) *Model {
	if t == nil {
		t = theme.CurrentTheme()
	}

	inp := input.New()
	inp.SetPlaceholder("Filter results...")
	inp.SetTheme(t)

	footer := bar.New(
		bar.FooterAnchored(),
		bar.Left("~ "+appName),
		bar.Cards(
			bar.Card{Command: "enter", Label: "apply", Variant: bar.CardPrimary, Enabled: true, Priority: 3},
			bar.Card{Command: "esc", Label: "clear", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
		bar.WithTheme(t),
	)

	return &Model{Input: inp, Footer: footer}
}

// SetSize keeps recipe internals sized together.
func (m *Model) SetSize(width int) {
	m.Input.SetSize(max(1, width-4), 1)
	m.Footer.SetSize(width, 1)
}

// Focus returns the command for focusing the filter input.
func (m *Model) Focus() tea.Cmd { return m.Input.Focus() }

// SetTheme updates all recipe internals.
func (m *Model) SetTheme(t theme.Theme) {
	m.Input.SetTheme(t)
	m.Footer.SetTheme(t)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
