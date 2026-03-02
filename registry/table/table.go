// Package table provides a simple header+rows table widget.
// Copy this file into your project: bento add table
//
// The widget calls theme.CurrentTheme() in View() — never stores theme state.
// Dependencies:
//   - charm.land/bubbletea/v2
//   - charm.land/lipgloss/v2
//   - github.com/cloudboy-jh/bentotui/theme
package table

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Model displays tabular data with a styled header row.
type Model struct {
	width   int
	height  int
	headers []string
	rows    [][]string
}

// New creates a table with the given column headers.
func New(headers ...string) *Model {
	return &Model{headers: headers}
}

// AddRow appends a data row. Extra cells are ignored; missing cells are blank.
func (t *Model) AddRow(cells ...string) {
	t.rows = append(t.rows, cells)
}

// Clear removes all data rows (keeps headers).
func (t *Model) Clear() { t.rows = nil }

func (t *Model) Init() tea.Cmd                           { return nil }
func (t *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }

func (t *Model) View() tea.View {
	if t.width <= 0 || len(t.headers) == 0 {
		return tea.NewView("")
	}

	th := theme.CurrentTheme()
	colWidth := t.width / len(t.headers)
	if colWidth < 1 {
		colWidth = 1
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(th.Text.Primary))
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(th.Border.Normal))
	cellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(th.Text.Primary))

	lines := make([]string, 0, 2+len(t.rows))

	headerCells := make([]string, len(t.headers))
	for i, h := range t.headers {
		headerCells[i] = headerStyle.Render(clip(h, colWidth))
	}
	lines = append(lines, strings.Join(headerCells, " "))
	lines = append(lines, borderStyle.Render(strings.Repeat("─", t.width)))

	for _, row := range t.rows {
		cells := make([]string, len(t.headers))
		for i := range t.headers {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			cells[i] = cellStyle.Render(clip(cell, colWidth))
		}
		lines = append(lines, strings.Join(cells, " "))
	}

	return tea.NewView(strings.Join(lines, "\n"))
}

func (t *Model) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Model) GetSize() (int, int) { return t.width, t.height }

func clip(s string, width int) string {
	r := []rune(s)
	if len(r) > width {
		return string(r[:width])
	}
	return s
}
