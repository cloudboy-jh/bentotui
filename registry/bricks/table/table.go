// Brick: Table
// +-----------------------------------+
// | col A | col B | col C            |
// +-----------------------------------+
// | row data                          |
// +-----------------------------------+
// Header + row data renderer.
// Copy this file into your project: bento add table
package table

import (
	"sort"

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

type VisualStyle string

const (
	VisualClean VisualStyle = "clean"
	VisualGrid  VisualStyle = "grid"
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
	visual     VisualStyle
	focused    bool
	inner      bubblestable.Model
	theme      theme.Theme // nil = use theme.CurrentTheme()
}

// New creates a table with the given column headers.
func New(headers ...string) *Model {
	cols := make([]Column, len(headers))
	for i, h := range headers {
		cols[i] = Column{Header: h, Align: AlignLeft, Min: 4, Pri: 1}
	}
	inner := bubblestable.New(
		bubblestable.WithColumns([]bubblestable.Column{{Title: "", Width: 1}}),
		bubblestable.WithRows(nil),
		bubblestable.WithWidth(1),
		bubblestable.WithHeight(2),
		bubblestable.WithFocused(true),
	)
	t := &Model{columns: cols, inner: inner, focused: true, width: 1, height: 2, visual: VisualClean}
	t.syncData()
	return t
}

func (t *Model) AddRow(cells ...string) { t.rows = append(t.rows, cells); t.syncData() }
func (t *Model) Clear()                 { t.rows = nil; t.syncData() }

func (t *Model) SetBorderless(v bool) { t.borderless = v; t.syncData() }
func (t *Model) SetCompact(v bool)    { t.compact = v; t.syncData() }
func (t *Model) SetVisualStyle(v VisualStyle) {
	if v == VisualGrid {
		t.visual = VisualGrid
	} else {
		t.visual = VisualClean
	}
}

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
	t.syncData()
}

func (t *Model) SetColumnWidth(index, width int) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	t.columns[index].Width = width
	t.syncData()
}

func (t *Model) SetColumnAlign(index int, align Align) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	t.columns[index].Align = align
	t.syncData()
}

func (t *Model) SetColumnMinWidth(index, width int) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	if width < 1 {
		width = 1
	}
	t.columns[index].Min = width
	t.syncData()
}

func (t *Model) SetColumnPriority(index, priority int) {
	if index < 0 || index >= len(t.columns) {
		return
	}
	if priority < 1 {
		priority = 1
	}
	t.columns[index].Pri = priority
	t.syncData()
}

// SetTheme updates the theme. Call on ThemeChangedMsg.
func (t *Model) SetTheme(th theme.Theme) { t.theme = th }

func (t *Model) activeTheme() theme.Theme {
	if t.theme != nil {
		return t.theme
	}
	return theme.CurrentTheme()
}

func (t *Model) Init() tea.Cmd { return nil }

func (t *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.SetSize(msg.Width, msg.Height)
		return t, nil
	}
	if !t.focused {
		if _, ok := msg.(tea.KeyMsg); ok {
			return t, nil
		}
	}
	updated, cmd := t.inner.Update(msg)
	t.inner = updated
	return t, cmd
}

func (t *Model) View() tea.View {
	if t.width <= 0 || t.height <= 0 || len(t.columns) == 0 {
		return tea.NewView("")
	}
	t.applyTheme()
	return tea.NewView(t.inner.View())
}

func (t *Model) SetSize(width, height int) {
	t.width = max(1, width)
	t.height = max(2, height)
	t.syncData()
}

func (t *Model) GetSize() (int, int) { return t.width, t.height }

func (t *Model) Focus() {
	t.focused = true
	t.inner.Focus()
}

func (t *Model) Blur() {
	t.focused = false
	t.inner.Blur()
}

func (t *Model) IsFocused() bool { return t.focused }

// applyTheme sets lipgloss styles on the inner bubbles table at render time.
//
// KEY RULE — Cell must NOT carry Background:
// bubbles/table renderRow() calls Cell.Render() on each cell individually,
// then joins cells, then calls Selected.Render() on the joined string.
// If Cell.Background is set, those escape codes are already embedded before
// Selected re-paints — causing color bleed on padding chars. Leave Background
// unpainted on Cell; Selected owns the whole-row background.
func (t *Model) applyTheme() {
	th := t.activeTheme()
	s := bubblestable.DefaultStyles()

	headerBG := th.BackgroundPanel()
	selectedBG := th.SelectionBG()
	selectedFG := th.SelectionFG()
	textFG := th.Text()

	pad := 1
	if t.compact || t.borderless {
		pad = 0
	}

	s.Header = lipgloss.NewStyle().
		Bold(true).
		Padding(0, pad).
		Foreground(textFG).
		Background(headerBG)

	// NO Background on Cell — see doc above.
	s.Cell = lipgloss.NewStyle().
		Padding(0, pad).
		Foreground(textFG)

	if t.focused {
		s.Selected = lipgloss.NewStyle().
			Bold(true).
			Padding(0, pad).
			Foreground(selectedFG).
			Background(selectedBG)
	} else {
		s.Selected = lipgloss.NewStyle().
			Padding(0, pad).
			Foreground(th.TextAccent())
	}

	if t.visual == VisualGrid && !t.borderless {
		gridBorder := lipgloss.NormalBorder()
		headerBorderFG := th.BorderFocus()
		cellBorderFG := th.BorderSubtle()
		s.Header = s.Header.
			BorderStyle(gridBorder).
			BorderBottom(true).
			BorderRight(true).
			BorderForeground(headerBorderFG)
		s.Cell = s.Cell.
			BorderStyle(gridBorder).
			BorderBottom(true).
			BorderRight(true).
			BorderForeground(cellBorderFG)
		s.Selected = s.Selected.
			BorderStyle(gridBorder).
			BorderBottom(true).
			BorderRight(true).
			BorderForeground(headerBorderFG)
	}

	t.inner.SetStyles(s)
}

func (t *Model) syncData() {
	colWidths := t.computeColumnWidths(max(1, t.width))
	cols := make([]bubblestable.Column, 0, len(t.columns))
	for i, col := range t.columns {
		cols = append(cols, bubblestable.Column{Title: col.Header, Width: colWidths[i]})
	}

	rows := make([]bubblestable.Row, 0, len(t.rows))
	for _, row := range t.rows {
		cells := make([]string, len(t.columns))
		for i := range t.columns {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			cells[i] = alignText(cell, colWidths[i], t.columns[i].Align)
		}
		rows = append(rows, bubblestable.Row(cells))
	}

	t.inner.SetColumns(cols)
	t.inner.SetRows(rows)
	t.inner.SetWidth(max(1, t.width))
	t.inner.SetHeight(max(2, t.height))
	if t.focused {
		t.inner.Focus()
	} else {
		t.inner.Blur()
	}
}

func (t *Model) computeColumnWidths(totalWidth int) []int {
	count := len(t.columns)
	if count == 0 {
		return nil
	}
	widths := make([]int, count)
	mins := make([]int, count)
	available := totalWidth
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
		s = ansi.Truncate(s, width, "…")
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
