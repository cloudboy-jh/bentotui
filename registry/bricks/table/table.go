// Brick: Table:
// +-----------------------------------+
// | col A | col B | col C            |
// +-----------------------------------+
// | row data                          |
// +-----------------------------------+
// Header + row data renderer.
// Package table provides a table widget backed by bubbles/table.
// Copy this file into your project: bento add table
package table

import (
	"sort"
	"strings"

	bubblestable "charm.land/bubbles/v2/table"
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
	Min    int
	Pri    int
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
	inner      bubblestable.Model
}

// New creates a table with the given column headers.
func New(headers ...string) *Model {
	cols := make([]Column, len(headers))
	for i, h := range headers {
		cols[i] = Column{Header: h, Align: AlignLeft, Min: 4, Pri: 1}
	}
	inner := bubblestable.New(bubblestable.WithColumns([]bubblestable.Column{{Title: "", Width: 1}}), bubblestable.WithRows(nil), bubblestable.WithWidth(1), bubblestable.WithHeight(1))
	inner.Blur()
	return &Model{columns: cols, inner: inner}
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
	if col.Min <= 0 {
		col.Min = 4
	}
	if col.Pri <= 0 {
		col.Pri = 1
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

func (t *Model) SetColumnMinWidth(index, width int) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	if width < 1 {
		width = 1
	}
	t.columns[index].Min = width
}

func (t *Model) SetColumnPriority(index, priority int) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	if priority < 1 {
		priority = 1
	}
	t.columns[index].Pri = priority
}

func (t *Model) Init() tea.Cmd { return nil }

func (t *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := t.inner.Update(msg)
	t.inner = updated
	return t, cmd
}

func (t *Model) View() tea.View {
	if t.width <= 0 || len(t.columns) == 0 {
		return tea.NewView("")
	}
	t.syncInner()
	return tea.NewView(t.inner.View())
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
		return strings.Repeat("-", t.width)
	}
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("-", w)
	}
	if t.compact {
		return strings.Join(parts, " ")
	}
	return strings.Join(parts, strings.Repeat("-", len(sep)))
}

func (t *Model) computeColumnWidths(totalWidth, sepWidth int) []int {
	count := len(t.columns)
	if count == 0 {
		return nil
	}
	widths := make([]int, count)
	mins := make([]int, count)
	seps := (count - 1) * sepWidth
	available := totalWidth - seps
	if available < count {
		available = count
	}

	base := 0
	flexCount := 0
	for i, col := range t.columns {
		minW := col.Min
		if minW < 1 {
			minW = 1
		}
		mins[i] = minW
		if col.Width > 0 {
			w := col.Width
			if w < minW {
				w = minW
			}
			widths[i] = w
			base += w
			continue
		}
		headerMin := lipgloss.Width(col.Header) + 1
		if headerMin < minW {
			headerMin = minW
		}
		widths[i] = headerMin
		base += headerMin
		flexCount++
	}

	if base > available {
		idxs := make([]int, len(t.columns))
		for i := range t.columns {
			idxs[i] = i
		}
		sort.Slice(idxs, func(i, j int) bool {
			a := t.columns[idxs[i]]
			b := t.columns[idxs[j]]
			if a.Pri != b.Pri {
				return a.Pri > b.Pri
			}
			return idxs[i] > idxs[j]
		})
		over := base - available
		for over > 0 {
			changed := false
			for _, idx := range idxs {
				if widths[idx] > mins[idx] {
					widths[idx]--
					over--
					changed = true
					if over == 0 {
						break
					}
				}
			}
			if !changed {
				break
			}
		}
	}

	used := 0
	for _, w := range widths {
		used += w
	}

	if flexCount > 0 && used < available {
		extra := available - used
		share := extra / flexCount
		rem := extra % flexCount
		for i, col := range t.columns {
			if col.Width > 0 {
				continue
			}
			widths[i] += share
			if rem > 0 {
				widths[i]++
				rem--
			}
		}
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

func (t *Model) syncInner() {
	th := theme.CurrentTheme()
	styles := bubblestable.DefaultStyles()
	styles.Header = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(th.Text.Primary))
	styles.Cell = lipgloss.NewStyle().Foreground(lipgloss.Color(th.Text.Primary))
	styles.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color(th.Text.Primary))
	t.inner.SetStyles(styles)

	sep := t.separator()
	colWidths := t.computeColumnWidths(t.width, len(sep))

	headerCells := make([]string, len(t.columns))
	for i, col := range t.columns {
		headerCells[i] = alignText(col.Header, colWidths[i], col.Align)
	}
	headerLine := strings.Join(headerCells, sep)
	headerLine = lipgloss.NewStyle().Width(max(1, t.width)).MaxWidth(max(1, t.width)).Render(headerLine)

	rows := make([]bubblestable.Row, 0, len(t.rows)+1)
	if !t.compact {
		rows = append(rows, bubblestable.Row{t.divider(colWidths, sep)})
	}
	for _, row := range t.rows {
		cells := make([]string, len(t.columns))
		for i := range t.columns {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			cells[i] = alignText(cell, colWidths[i], t.columns[i].Align)
		}
		line := strings.Join(cells, sep)
		line = lipgloss.NewStyle().Width(max(1, t.width)).MaxWidth(max(1, t.width)).Render(line)
		rows = append(rows, bubblestable.Row{line})
	}

	t.inner.SetColumns([]bubblestable.Column{{Title: headerLine, Width: max(1, t.width)}})
	t.inner.SetRows(rows)
	t.inner.SetWidth(max(1, t.width))
	t.inner.SetHeight(max(2, t.height))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
