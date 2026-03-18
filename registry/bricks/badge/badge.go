// Brick: Badge
// +-----------------------------+
// | [ variant label text ]      |
// +-----------------------------+
// Inline status/emphasis pill.
// Copy this file into your project: bento add badge
package badge

import (
	"image/color"

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
	theme   theme.Theme // nil = use theme.CurrentTheme()
}

func New(text string) *Model {
	return &Model{text: text, variant: VariantAccent, bold: true}
}

func (m *Model) SetText(text string)                 { m.text = text }
func (m *Model) SetVariant(v Variant)                { m.variant = v }
func (m *Model) SetBold(b bool)                      { m.bold = b }
func (m *Model) SetSize(w, h int)                    { m.width, m.height = w, h }
func (m *Model) GetSize() (int, int)                 { return m.width, m.height }
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

// SetTheme updates the theme. Call on ThemeChangedMsg.
func (m *Model) SetTheme(t theme.Theme) { m.theme = t }

func (m *Model) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

func (m *Model) View() tea.View {
	t := m.activeTheme()
	fg, bg := badgeColors(t, m.variant)
	style := lipgloss.NewStyle().
		Foreground(fg).
		Background(bg).
		Padding(0, 1).
		Bold(m.bold)
	return tea.NewView(style.Render(m.text))
}

func badgeColors(t theme.Theme, v Variant) (fg, bg color.Color) {
	switch v {
	case VariantNeutral:
		return t.Text(), t.BackgroundPanel()
	case VariantInfo:
		return t.TextInverse(), t.Info()
	case VariantSuccess:
		return t.TextInverse(), t.Success()
	case VariantWarning:
		return t.TextInverse(), t.Warning()
	case VariantDanger:
		return t.TextInverse(), t.Error()
	default: // VariantAccent
		return t.SelectionFG(), t.SelectionBG()
	}
}
