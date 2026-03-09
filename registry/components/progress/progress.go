// Package progress provides a themed horizontal progress bar.
// Copy this file into your project: bento add progress
package progress

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Model struct {
	value       float64
	width       int
	label       string
	showPercent bool
}

func New(width int) *Model {
	if width < 10 {
		width = 24
	}
	return &Model{width: width, showPercent: true}
}

func (m *Model) SetValue(v float64)                  { m.value = clamp01(v) }
func (m *Model) Value() float64                      { return m.value }
func (m *Model) SetLabel(v string)                   { m.label = v }
func (m *Model) SetShowPercent(v bool)               { m.showPercent = v }
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *Model) SetSize(width, _ int) {
	if width > 0 {
		m.width = width
	}
}
func (m *Model) GetSize() (int, int) { return m.width, 1 }

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	barWidth := m.width
	if barWidth < 10 {
		barWidth = 24
	}
	filled := int(float64(barWidth) * m.value)
	if filled > barWidth {
		filled = barWidth
	}
	if filled < 0 {
		filled = 0
	}
	empty := barWidth - filled

	fillStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Selection.BG, t.Text.Accent)))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Border.Subtle, t.Text.Muted)))
	bar := fillStyle.Render(repeat("#", filled)) + emptyStyle.Render(repeat("-", empty))

	out := bar
	if m.showPercent {
		out += fmt.Sprintf(" %3.0f%%", m.value*100)
	}
	if m.label != "" {
		out = m.label + " " + out
	}
	return tea.NewView(out)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
