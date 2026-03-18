// Brick: Progress:
// +-----------------------------------+
// | label [####------] 42%            |
// +-----------------------------------+
// Horizontal progress indicator.
// Package progress provides a themed horizontal progress bar backed by
// bubbles/progress.
// Copy this file into your project: bento add progress
package progress

import (
	"fmt"

	bubblesprogress "charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Model struct {
	value       float64
	width       int
	label       string
	showPercent bool
	bar         bubblesprogress.Model
}

func New(width int) *Model {
	if width < 10 {
		width = 24
	}
	b := bubblesprogress.New(bubblesprogress.WithWidth(width), bubblesprogress.WithFillCharacters('█', '─'))
	return &Model{width: width, showPercent: true, bar: b}
}

func (m *Model) SetValue(v float64) { m.value = clamp01(v) }
func (m *Model) Value() float64     { return m.value }
func (m *Model) SetLabel(v string)  { m.label = v }
func (m *Model) SetShowPercent(v bool) {
	m.showPercent = v
	if v {
		m.bar.ShowPercentage = true
		return
	}
	m.bar.ShowPercentage = false
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := m.bar.Update(msg)
	m.bar = updated
	return m, cmd
}

func (m *Model) SetSize(width, _ int) {
	if width <= 0 {
		return
	}
	m.width = width
	m.bar.SetWidth(width)
}

func (m *Model) GetSize() (int, int) { return m.width, 1 }

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	m.bar.FullColor = t.SelectionBG()
	m.bar.EmptyColor = t.BorderSubtle()
	m.bar.PercentageStyle = lipgloss.NewStyle().Foreground(t.TextMuted())
	m.bar.ShowPercentage = m.showPercent
	if m.width > 0 {
		m.bar.SetWidth(m.width)
	}
	line := m.bar.ViewAs(m.value)
	if m.label != "" {
		line = fmt.Sprintf("%s %s", m.label, line)
	}
	return tea.NewView(line)
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

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
