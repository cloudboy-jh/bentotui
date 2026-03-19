package emptystatepane

import (
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/card"
	"github.com/cloudboy-jh/bentotui/registry/bricks/text"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Model is a copy-and-own recipe for a centered empty-state pane.
type Model struct {
	Card *card.Model
	Body *text.Model
}

// New creates a titled empty-state card with message content.
func New(title, message string, t theme.Theme) *Model {
	if t == nil {
		t = theme.CurrentTheme()
	}
	body := text.New(message)
	panel := card.New(
		card.Title(title),
		card.Flat(),
		card.Content(body),
		card.WithTheme(t),
	)
	return &Model{Card: panel, Body: body}
}

// SetSize forwards size to the recipe container.
func (m *Model) SetSize(width, height int) {
	m.Card.SetSize(width, height)
}

// View satisfies tea.Model-style usage in page composition.
func (m *Model) View() tea.View { return m.Card.View() }

// SetTheme updates the card theme for live switching.
func (m *Model) SetTheme(t theme.Theme) {
	m.Card.SetTheme(t)
}
