package widgets

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
)

// Table displays data in rows and columns.
type Table struct {
	width   int
	height  int
	headers []string
	rows    [][]string
	theme   theme.Theme
}

// NewTable creates a table with headers.
func NewTable(headers ...string) *Table {
	return &Table{
		headers: headers,
		theme:   theme.CurrentTheme(),
	}
}

// AddRow adds a data row.
func (t *Table) AddRow(cells ...string) {
	t.rows = append(t.rows, cells)
}

// Clear removes all rows.
func (t *Table) Clear() {
	t.rows = nil
}

func (t *Table) Init() tea.Cmd { return nil }

func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return t, nil
}

func (t *Table) View() tea.View {
	if t.width <= 0 || t.height <= 0 {
		return tea.NewView("")
	}

	var lines []string
	colWidth := t.width / len(t.headers)
	if colWidth < 1 {
		colWidth = 1
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(t.theme.Text.Primary))
	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.theme.Border.Normal))
	cellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.theme.Text.Primary))

	var headerCells []string
	for _, h := range t.headers {
		if len(h) > colWidth {
			h = h[:colWidth]
		}
		headerCells = append(headerCells, headerStyle.Render(h))
	}
	lines = append(lines, strings.Join(headerCells, " "))
	lines = append(lines, borderStyle.Render(strings.Repeat("─", t.width)))

	for _, row := range t.rows {
		var cells []string
		for i, cell := range row {
			if i < len(t.headers) {
				if len(cell) > colWidth {
					cell = cell[:colWidth]
				}
				cells = append(cells, cellStyle.Render(cell))
			}
		}
		lines = append(lines, strings.Join(cells, " "))
	}

	return tea.NewView(strings.Join(lines, "\n"))
}

func (t *Table) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Table) GetSize() (int, int) {
	return t.width, t.height
}

// SetTheme updates the theme.
func (t *Table) SetTheme(theme theme.Theme) {
	t.theme = theme
}

// HeightConstraint returns MinHeight based on header + rows.
func (t *Table) HeightConstraint(width int) HeightConstraint {
	minRows := 2 + len(t.rows)
	if minRows < 3 {
		minRows = 3
	}
	return MinHeight(minRows)
}

var _ core.Component = (*Table)(nil)
var _ core.Sizeable = (*Table)(nil)
var _ HeightConstraintAware = (*Table)(nil)
