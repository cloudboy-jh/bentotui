// Brick: List:
// +-----------------------------------+
// | row 1                              |
// | row 2                              |
// | ...                                |
// +-----------------------------------+
// Scrollable plain-text row list.
// Package list provides a list widget backed by bubbles/list.
// Copy this file into your project: bento add list
package list

import (
	"io"
	"strings"

	bubbleslist "charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// RowKind marks whether a row is normal content or a section header.
type RowKind int

const (
	RowItem RowKind = iota
	RowSection
)

type Tone string

const (
	ToneNeutral Tone = "neutral"
	ToneInfo    Tone = "info"
	ToneSuccess Tone = "success"
	ToneWarn    Tone = "warn"
	ToneDanger  Tone = "danger"
)

type SelectedStyle string

const (
	SelectedDefault SelectedStyle = "default"
	SelectedSubtle  SelectedStyle = "subtle"
)

// Row is a structured list row.
type Row struct {
	Kind          RowKind
	Primary       string
	Secondary     string
	RightStat     string
	Tone          Tone
	SelectedStyle SelectedStyle
	Label         string
	Status        string
	Stat          string
	Section       string
}

// RowFormatter renders one row at the given width.
type RowFormatter func(row Row, selected bool, width int) string

// Model is a scrollable list that shows the last N items that fit in height.
// It produces plain text — no ANSI styling — so the parent panel can paint
// the background without bleed-through.
type Model struct {
	width     int
	height    int
	rows      []Row
	max       int
	cursor    int
	density   Density
	formatter RowFormatter
	inner     bubbleslist.Model
	delegate  *rowDelegate
}

type Density string

const (
	DensityComfortable Density = "comfortable"
	DensityCompact     Density = "compact"
)

type listItem struct{ row Row }

func (i listItem) FilterValue() string {
	if strings.TrimSpace(i.row.Primary) != "" {
		return i.row.Primary
	}
	if strings.TrimSpace(i.row.Label) != "" {
		return i.row.Label
	}
	if strings.TrimSpace(i.row.Section) != "" {
		return i.row.Section
	}
	return ""
}

type rowDelegate struct{ owner *Model }

func (d rowDelegate) Height() int  { return 1 }
func (d rowDelegate) Spacing() int { return 0 }
func (d rowDelegate) Update(msg tea.Msg, m *bubbleslist.Model) tea.Cmd {
	return nil
}

func (d rowDelegate) Render(w io.Writer, m bubbleslist.Model, index int, item bubbleslist.Item) {
	li, ok := item.(listItem)
	if !ok {
		return
	}
	selected := index == m.Index()
	line := d.owner.renderRow(li.row, selected, m.Width())
	_, _ = io.WriteString(w, line)
}

// New creates a list with an optional cap on stored items.
// maxItems <= 0 defaults to 200.
func New(maxItems int) *Model {
	if maxItems <= 0 {
		maxItems = 200
	}
	l := &Model{max: maxItems, density: DensityComfortable}
	d := &rowDelegate{owner: l}
	inner := bubbleslist.New([]bubbleslist.Item{}, *d, 1, 1)
	inner.SetShowTitle(false)
	inner.SetShowFilter(false)
	inner.SetShowStatusBar(false)
	inner.SetShowPagination(false)
	inner.SetShowHelp(false)
	inner.SetFilteringEnabled(false)
	inner.DisableQuitKeybindings()
	l.inner = inner
	l.delegate = d
	return l
}

// Append adds an item to the bottom of the list.
func (l *Model) Append(item string) {
	l.AppendRow(Row{Kind: RowItem, Label: item})
}

// AppendRow adds a structured row to the bottom of the list.
func (l *Model) AppendRow(row Row) {
	if row.Kind == RowSection && strings.TrimSpace(row.Section) == "" {
		row.Section = row.Label
	}
	l.rows = append(l.rows, row)
	if len(l.rows) > l.max {
		l.rows = l.rows[1:]
	}
	l.syncInner()
}

// AppendSection adds a section/header row to the bottom.
func (l *Model) AppendSection(title string) {
	l.AppendRow(Row{Kind: RowSection, Section: title})
}

// Prepend adds an item to the top of the list.
func (l *Model) Prepend(item string) {
	l.PrependRow(Row{Kind: RowItem, Label: item})
}

// PrependRow adds a structured row to the top of the list.
func (l *Model) PrependRow(row Row) {
	if row.Kind == RowSection && strings.TrimSpace(row.Section) == "" {
		row.Section = row.Label
	}
	l.rows = append([]Row{row}, l.rows...)
	if len(l.rows) > l.max {
		l.rows = l.rows[:l.max]
	}
	l.syncInner()
}

// PrependSection adds a section/header row to the top.
func (l *Model) PrependSection(title string) {
	l.PrependRow(Row{Kind: RowSection, Section: title})
}

// Clear removes all items.
func (l *Model) Clear() {
	l.rows = nil
	l.cursor = 0
	l.syncInner()
}

// Items returns a copy of the current item list.
func (l *Model) Items() []string {
	out := make([]string, 0, len(l.rows))
	for _, row := range l.rows {
		if row.Kind == RowSection {
			continue
		}
		out = append(out, row.Label)
	}
	return out
}

// SetFormatter sets a custom row formatter.
func (l *Model) SetFormatter(f RowFormatter) {
	l.formatter = f
	l.syncInner()
}

// SetDensity controls row verbosity in default formatter.
func (l *Model) SetDensity(v Density) {
	switch v {
	case DensityCompact:
		l.density = DensityCompact
	default:
		l.density = DensityComfortable
	}
	l.syncInner()
}

// SetCursor sets the selected item index (item rows only).
func (l *Model) SetCursor(i int) {
	if i < 0 {
		i = 0
	}
	l.cursor = i
	l.syncInner()
}

func (l *Model) Init() tea.Cmd { return nil }

func (l *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updated, cmd := l.inner.Update(msg)
	l.inner = updated
	return l, cmd
}

func (l *Model) View() tea.View {
	if len(l.rows) == 0 {
		return tea.NewView("")
	}
	l.syncInner()
	out := l.inner.View()
	parts := strings.Split(out, "\n")
	clean := make([]string, 0, len(parts))
	for _, line := range parts {
		if strings.TrimSpace(line) == "" {
			continue
		}
		clean = append(clean, line)
	}
	return tea.NewView(strings.Join(clean, "\n"))
}

func (l *Model) SetSize(width, height int) {
	l.width = width
	l.height = height
	l.inner.SetSize(max(1, width), max(1, height))
	l.syncInner()
}

func (l *Model) GetSize() (int, int) { return l.width, l.height }

func (l *Model) renderRow(row Row, selected bool, width int) string {
	if l.formatter != nil {
		return l.formatter(row, selected, width)
	}
	return defaultFormatter(row, selected, width, l.density)
}

func defaultFormatter(row Row, selected bool, width int, density Density) string {
	if row.Kind == RowSection {
		title := strings.TrimSpace(row.Section)
		if title == "" {
			title = strings.TrimSpace(row.Label)
		}
		if width > 0 {
			prefix := strings.ToUpper(title)
			if prefix == "" {
				prefix = "section"
			}
			line := "  " + prefix + " "
			if lipgloss.Width(line) < width {
				line += strings.Repeat("-", width-lipgloss.Width(line))
			}
			return line
		}
		if title == "" {
			return "  --"
		}
		return "  " + strings.ToUpper(title)
	}

	prefix := "  "
	if selected {
		if row.SelectedStyle == SelectedSubtle {
			prefix = "* "
		} else {
			prefix = "> "
		}
	}

	primary := strings.TrimSpace(row.Primary)
	if primary == "" {
		primary = strings.TrimSpace(row.Label)
	}
	secondary := strings.TrimSpace(row.Secondary)
	status := strings.TrimSpace(row.Status)
	if status == "" && row.Tone != "" && row.Tone != ToneNeutral {
		status = string(row.Tone)
	}
	if density == DensityCompact {
		secondary = ""
	}
	leftParts := make([]string, 0, 2)
	if status != "" {
		leftParts = append(leftParts, "["+status+"]")
	}
	if secondary != "" {
		primary = primary + " - " + secondary
	}
	leftParts = append(leftParts, primary)
	left := prefix + strings.TrimSpace(strings.Join(leftParts, " "))
	stat := strings.TrimSpace(row.RightStat)
	if stat == "" {
		stat = strings.TrimSpace(row.Stat)
	}
	if stat == "" || width <= 0 {
		if width > 0 {
			return ansi.Truncate(left, width, "")
		}
		return left
	}
	return fitLeftRight(left, stat, width)
}

func fitLeftRight(left, right string, width int) string {
	if width <= 0 {
		return left
	}
	rw := lipgloss.Width(right)
	if rw >= width {
		return ansi.Truncate(right, width, "")
	}
	gap := 1
	maxLeft := width - rw - gap
	if maxLeft < 1 {
		maxLeft = 1
	}
	l := ansi.Truncate(left, maxLeft, "")
	space := width - lipgloss.Width(l) - rw
	if space < 1 {
		space = 1
	}
	return l + strings.Repeat(" ", space) + right
}

func (l *Model) syncInner() {
	items := make([]bubbleslist.Item, 0, len(l.rows))
	selectedRowIdx := 0
	itemIdx := 0
	for i, row := range l.rows {
		items = append(items, listItem{row: row})
		if row.Kind != RowItem {
			continue
		}
		if itemIdx == l.cursor {
			selectedRowIdx = i
		}
		itemIdx++
	}
	if cmd := l.inner.SetItems(items); cmd != nil {
		_ = cmd
	}
	l.inner.Select(selectedRowIdx)
	w := l.width
	if w <= 0 {
		w = 1
		for _, row := range l.rows {
			line := l.renderRow(row, false, 0)
			if lw := lipgloss.Width(line); lw > w {
				w = lw
			}
		}
	}
	h := l.height
	if h <= 0 {
		h = len(l.rows) + 2
	}
	l.inner.SetSize(max(1, w), max(1, h))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
