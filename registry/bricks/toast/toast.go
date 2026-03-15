// Brick: Toast:
// +-----------------------------------+
// | info message                       |
// | success message                    |
// +-----------------------------------+
// Stacked transient notifications.
// Package toast provides a lightweight stacked notification model.
// Copy this file into your project: bento add toast
package toast

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

type Variant string

const (
	Info    Variant = "info"
	Success Variant = "success"
	Warning Variant = "warning"
	Danger  Variant = "danger"
)

type Item struct {
	ID      int
	Text    string
	Variant Variant
}

type Model struct {
	items  []Item
	nextID int
	width  int
	height int
	max    int
}

func New(maxVisible int) *Model {
	if maxVisible <= 0 {
		maxVisible = 3
	}
	return &Model{max: maxVisible, nextID: 1}
}

func (m *Model) Init() tea.Cmd                       { return nil }
func (m *Model) Update(tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) GetSize() (int, int) {
	return m.width, min(len(m.items), m.max)
}

func (m *Model) Push(text string, variant Variant) int {
	id := m.nextID
	m.nextID++
	m.items = append(m.items, Item{ID: id, Text: text, Variant: variant})
	if len(m.items) > m.max {
		m.items = m.items[len(m.items)-m.max:]
	}
	return id
}

func (m *Model) Dismiss(id int) {
	filtered := m.items[:0]
	for _, item := range m.items {
		if item.ID != id {
			filtered = append(filtered, item)
		}
	}
	m.items = filtered
}

func (m *Model) Clear() { m.items = nil }

func (m *Model) Items() []Item {
	return append([]Item(nil), m.items...)
}

func (m *Model) View() tea.View {
	if len(m.items) == 0 {
		return tea.NewView("")
	}
	t := theme.CurrentTheme()
	width := m.width
	if width <= 0 {
		width = 48
	}
	count := min(len(m.items), m.max)
	if m.height > 0 {
		count = min(count, m.height)
	}
	start := len(m.items) - count
	if start < 0 {
		start = 0
	}

	lines := make([]string, 0, count)
	for _, item := range m.items[start:] {
		bg, fg := toastColors(t, item.Variant)
		content := fmt.Sprintf("! %s", item.Text)
		lines = append(lines, styles.Row(bg, fg, width, content))
	}
	return tea.NewView(strings.Join(lines, "\n"))
}

func toastColors(t theme.Theme, v Variant) (bg, fg string) {
	fg = pick(t.Text.Inverse, t.Text.Primary)
	switch v {
	case Success:
		bg = pick(t.State.Success, t.Selection.BG)
	case Warning:
		bg = pick(t.State.Warning, t.Selection.BG)
	case Danger:
		bg = pick(t.State.Danger, t.Selection.BG)
	default:
		bg = pick(t.State.Info, t.Selection.BG)
	}
	return bg, fg
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
