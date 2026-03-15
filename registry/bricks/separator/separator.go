// Brick: Separator:
// +-----------------------------+
// |-----------------------------|
// +-----------------------------+
// Horizontal or vertical divider.
// Package separator provides horizontal and vertical themed rules.
// Copy this file into your project: bento add separator
package separator

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Orientation string

const (
	Horizontal Orientation = "horizontal"
	Vertical   Orientation = "vertical"
)

type Model struct {
	orientation Orientation
	length      int
	width       int
	height      int
}

func New(orientation Orientation, length int) *Model {
	if length <= 0 {
		length = 1
	}
	if orientation == "" {
		orientation = Horizontal
	}
	return &Model{orientation: orientation, length: length}
}

func (m *Model) SetOrientation(v Orientation) { m.orientation = v }
func (m *Model) SetLength(v int) {
	if v > 0 {
		m.length = v
	}
}
func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.orientation == Vertical {
		if height > 0 {
			m.length = height
		}
		return
	}
	if width > 0 {
		m.length = width
	}
}
func (m *Model) GetSize() (int, int) {
	if m.orientation == Vertical {
		return 1, m.length
	}
	return m.length, 1
}

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Border.Subtle, t.Border.Normal)))
	if m.orientation == Vertical {
		parts := make([]string, 0, m.length)
		for i := 0; i < m.length; i++ {
			parts = append(parts, "|")
		}
		return tea.NewView(style.Render(strings.Join(parts, "\n")))
	}
	return tea.NewView(style.Render(strings.Repeat("-", m.length)))
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
