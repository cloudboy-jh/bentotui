// Brick: Table:
// +-----------------------------------+
// | col A | col B | col C            |
// +-----------------------------------+
// | row data                          |
// +-----------------------------------+
// Header + row data renderer.
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
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Align string

const (
	AlignLeft   Align = "left"
	AlignCenter Align = "center"
	AlignRight  Align = "right"
)

type Column struct {
	Header string
	Width  int
	Align  Align
}

// Model displays tabular data with a styled header row.
type Model struct {
	width      int
	height     int
	columns    []Column
	rows       [][]string
	borderless bool
	compact    bool
}

// New creates a table with the given column headers.
func New(headers ...string) *Model {
	cols := make([]Column, len(headers))
	for i, h := range headers {
		cols[i] = Column{Header: h, Align: AlignLeft}
	}
	return &Model{columns: cols}
}

// AddRow appends a data row. Extra cells are ignored; missing cells are blank.
func (t *Model) AddRow(cells ...string) {
	t.rows = append(t.rows, cells)
}

// Clear removes all data rows (keeps headers).
func (t *Model) Clear() { t.rows = nil }

func (t *Model) SetBorderless(v bool) { t.borderless = v }
func (t *Model) SetCompact(v bool)    { t.compact = v }

func (t *Model) SetColumn(index int, col Column) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	if col.Align == "" {
		col.Align = AlignLeft
	}
	t.columns[index] = col
}

func (t *Model) SetColumnWidth(index, width int) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	t.columns[index].Width = width
}

func (t *Model) SetColumnAlign(index int, align Align) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	t.columns[index].Align = align
}

func (t *Model) Init() tea.Cmd                           { return nil }
func (t *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }

func (t *Model) View() tea.View {
	if t.width <= 0 || len(t.columns) == 0 {
		return tea.NewView("")
	}

	th := theme.CurrentTheme()
	sep := t.separator()
	colWidths := t.computeColumnWidths(t.width, len(sep))

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(th.Text.Primary))
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(th.Border.Normal))
	cellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(th.Text.Primary))

	lines := make([]string, 0, 2+len(t.rows))

	headerCells := make([]string, len(t.columns))
	for i, col := range t.columns {
		headerCells[i] = headerStyle.Render(alignText(col.Header, colWidths[i], col.Align))
	}
	lines = append(lines, strings.Join(headerCells, sep))
	if !t.compact {
		lines = append(lines, borderStyle.Render(t.divider(colWidths, sep)))
	}

	rowsToRender := t.rows
	if t.height > 0 {
		reserved := 1
		if !t.compact {
			reserved++
		}
		maxRows := t.height - reserved
		if maxRows < 0 {
			maxRows = 0
		}
		if len(rowsToRender) > maxRows {
			rowsToRender = rowsToRender[:maxRows]
		}
	}

	for _, row := range rowsToRender {
		cells := make([]string, len(t.columns))
		for i := range t.columns {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			cells[i] = cellStyle.Render(alignText(cell, colWidths[i], t.columns[i].Align))
		}
		lines = append(lines, strings.Join(cells, sep))
	}

	return tea.NewView(strings.Join(lines, "\n"))
}

func (t *Model) SetSize(width, height int) {
	t.width = width
	t.height = height
}

func (t *Model) GetSize() (int, int) { return t.width, t.height }

func (t *Model) separator() string {
	if t.compact {
		return " "
	}
	if t.borderless {
		return "  "
	}
	return " | "
}

func (t *Model) divider(widths []int, sep string) string {
	if t.borderless {
		return strings.Repeat("─", t.width)
	}
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("─", w)
	}
	if t.compact {
		return strings.Join(parts, " ")
	}
	return strings.Join(parts, strings.Repeat("─", len(sep)))
}

func (t *Model) computeColumnWidths(totalWidth, sepWidth int) []int {
	count := len(t.columns)
	if count == 0 {
		return nil
	}
	widths := make([]int, count)
	seps := (count - 1) * sepWidth
	available := totalWidth - seps
	if available < count {
		available = count
	}

	fixed := 0
	flexCount := 0
	for i, col := range t.columns {
		if col.Width > 0 {
			widths[i] = col.Width
			fixed += col.Width
			continue
		}
		flexCount++
	}

	remaining := available - fixed
	if remaining < flexCount {
		remaining = flexCount
	}
	share := 1
	if flexCount > 0 {
		share = remaining / flexCount
		if share < 1 {
			share = 1
		}
	}
	lastFlex := -1
	for i, col := range t.columns {
		if col.Width > 0 {
			continue
		}
		widths[i] = share
		lastFlex = i
	}

	used := 0
	for _, w := range widths {
		used += w
	}
	if lastFlex >= 0 {
		widths[lastFlex] += available - used
	}
	for i := range widths {
		if widths[i] < 1 {
			widths[i] = 1
		}
	}
	return widths
}

func alignText(s string, width int, align Align) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) > width {
		s = ansi.Truncate(s, width, "")
	}
	cell := lipgloss.NewStyle().Width(width)
	switch align {
	case AlignRight:
		return cell.Align(lipgloss.Right).Render(s)
	case AlignCenter:
		return cell.Align(lipgloss.Center).Render(s)
	default:
		return cell.Align(lipgloss.Left).Render(s)
	}
}
