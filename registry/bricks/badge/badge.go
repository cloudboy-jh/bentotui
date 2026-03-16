// Brick: Badge:
// +-----------------------------+
// | [ variant label text ]      |
// +-----------------------------+
// Inline status/emphasis pill.
// Package badge provides a small inline themed label.
// Copy this file into your project: bento add badge
package badge

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Variant controls badge color mapping.
type Variant string

const (
	VariantNeutral Variant = "neutral"
	VariantInfo    Variant = "info"
	VariantSuccess Variant = "success"
	VariantWarning Variant = "warning"
	VariantDanger  Variant = "danger"
	VariantAccent  Variant = "accent"
)

type Model struct {
	text    string
	variant Variant
	bold    bool
	width   int
	height  int
}

func New(text string) *Model {
	return &Model{text: text, variant: VariantAccent, bold: true}
}

func (m *Model) SetText(text string)                 { m.text = text }
func (m *Model) SetVariant(v Variant)                { m.variant = v }
func (m *Model) SetBold(b bool)                      { m.bold = b }
func (m *Model) SetSize(width, height int)           { m.width, m.height = width, height }
func (m *Model) GetSize() (int, int)                 { return m.width, m.height }
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	fg, bg := badgeColors(t, m.variant)
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg)).
		Padding(0, 1).
		Bold(m.bold)
	return tea.NewView(style.Render(m.text))
}

func badgeColors(t theme.Theme, v Variant) (fg, bg string) {
	switch v {
	case VariantNeutral:
		return pick(t.Text.Primary, t.Text.Inverse), pick(t.Surface.Panel, t.Surface.Canvas)
	case VariantInfo:
		return pick(t.Text.Inverse, t.Text.Primary), pick(t.State.Info, t.Text.Accent)
	case VariantSuccess:
		return pick(t.Text.Inverse, t.Text.Primary), pick(t.State.Success, t.Text.Accent)
	case VariantWarning:
		return pick(t.Text.Inverse, t.Text.Primary), pick(t.State.Warning, t.Text.Accent)
	case VariantDanger:
		return pick(t.Text.Inverse, t.Text.Primary), pick(t.State.Danger, t.Text.Accent)
	default:
		return pick(t.Selection.FG, t.Text.Inverse), pick(t.Selection.BG, t.Text.Accent)
	}
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
