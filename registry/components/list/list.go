// Package list provides a scrollable log-style list widget.
// Copy this file into your project: bento add list
//
// The widget returns plain text — the containing panel applies all color.
// Dependencies:
//   - charm.land/bubbletea/v2
//   - charm.land/lipgloss/v2
//   - github.com/charmbracelet/x/ansi
package list

import (
	"strings"

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

// Row is a structured list row.
type Row struct {
	Kind    RowKind
	Label   string
	Status  string
	Stat    string
	Section string
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
	formatter RowFormatter
}

// New creates a list with an optional cap on stored items.
// maxItems <= 0 defaults to 200.
func New(maxItems int) *Model {
	if maxItems <= 0 {
		maxItems = 200
	}
	return &Model{max: maxItems}
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
}

// PrependSection adds a section/header row to the top.
func (l *Model) PrependSection(title string) {
	l.PrependRow(Row{Kind: RowSection, Section: title})
}

// Clear removes all items.
func (l *Model) Clear() { l.rows = nil }

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
func (l *Model) SetFormatter(f RowFormatter) { l.formatter = f }

// SetCursor sets the selected item index (item rows only).
func (l *Model) SetCursor(i int) {
	if i < 0 {
		i = 0
	}
	l.cursor = i
}

func (l *Model) Init() tea.Cmd                           { return nil }
func (l *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return l, nil }

func (l *Model) View() tea.View {
	if len(l.rows) == 0 {
		return tea.NewView("")
	}
	if l.width <= 0 || l.height <= 0 {
		lines := make([]string, 0, len(l.rows))
		itemIdx := 0
		for _, row := range l.rows {
			selected := row.Kind == RowItem && itemIdx == l.cursor
			lines = append(lines, l.renderRow(row, selected, 0))
			if row.Kind == RowItem {
				itemIdx++
			}
		}
		return tea.NewView(strings.Join(lines, "\n"))
	}

	start := len(l.rows) - l.height
	if start < 0 {
		start = 0
	}
	lines := make([]string, 0, l.height)

	itemIdx := 0
	for i := 0; i < start; i++ {
		if l.rows[i].Kind == RowItem {
			itemIdx++
		}
	}

	for i := start; i < len(l.rows); i++ {
		row := l.rows[i]
		selected := row.Kind == RowItem && itemIdx == l.cursor
		line := l.renderRow(row, selected, l.width)
		if lipgloss.Width(line) > l.width {
			line = ansi.Truncate(line, l.width, "")
		}
		lines = append(lines, line)
		if row.Kind == RowItem {
			itemIdx++
		}
	}
	return tea.NewView(strings.Join(lines, "\n"))
}

func (l *Model) SetSize(width, height int) {
	l.width = width
	l.height = height
}

func (l *Model) GetSize() (int, int) { return l.width, l.height }

func (l *Model) renderRow(row Row, selected bool, width int) string {
	if l.formatter != nil {
		return l.formatter(row, selected, width)
	}
	return defaultFormatter(row, selected, width)
}

func defaultFormatter(row Row, selected bool, width int) string {
	if row.Kind == RowSection {
		title := strings.TrimSpace(row.Section)
		if title == "" {
			title = strings.TrimSpace(row.Label)
		}
		if width > 0 {
			prefix := title
			if prefix == "" {
				prefix = "section"
			}
			line := "-- " + prefix + " "
			if lipgloss.Width(line) < width {
				line += strings.Repeat("-", width-lipgloss.Width(line))
			}
			return line
		}
		if title == "" {
			return "--"
		}
		return "-- " + title
	}

	prefix := "  "
	if selected {
		prefix = "> "
	}
	leftParts := make([]string, 0, 2)
	if strings.TrimSpace(row.Status) != "" {
		leftParts = append(leftParts, row.Status)
	}
	leftParts = append(leftParts, row.Label)
	left := prefix + strings.TrimSpace(strings.Join(leftParts, " "))
	stat := strings.TrimSpace(row.Stat)
	if stat == "" || width <= 0 {
		return left
	}
	if lipgloss.Width(left)+1+lipgloss.Width(stat) > width {
		return left
	}
	return left + strings.Repeat(" ", width-lipgloss.Width(left)-lipgloss.Width(stat)) + stat
}
