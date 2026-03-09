// Package select provides a themed single-choice picker.
// Copy this file into your project: bento add select
package selectx

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Item struct {
	Label string
	Value string
}

type Model struct {
	items       []Item
	cursor      int
	selected    int
	open        bool
	focused     bool
	placeholder string
	width       int
	height      int
}

func New(items ...Item) *Model {
	return &Model{items: append([]Item(nil), items...), selected: -1, placeholder: "Select..."}
}

func (m *Model) SetItems(items []Item) {
	m.items = append([]Item(nil), items...)
	if len(m.items) == 0 {
		m.cursor = 0
		m.selected = -1
		return
	}
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.selected >= len(m.items) {
		m.selected = -1
	}
}

func (m *Model) SetPlaceholder(v string) { m.placeholder = v }
func (m *Model) Focus()                  { m.focused = true }
func (m *Model) Blur()                   { m.focused = false; m.open = false }
func (m *Model) IsFocused() bool         { return m.focused }
func (m *Model) Open() {
	if len(m.items) > 0 {
		m.open = true
	}
}
func (m *Model) Close() { m.open = false }
func (m *Model) ToggleOpen() {
	if m.open {
		m.Close()
	} else {
		m.Open()
	}
}
func (m *Model) SetSize(width, height int) {
	if width > 0 {
		m.width = width
	}
	m.height = height
}
func (m *Model) GetSize() (int, int) {
	if m.open {
		return m.width, m.visibleCount() + 1
	}
	return m.width, 1
}
func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Selected() (Item, bool) {
	if m.selected < 0 || m.selected >= len(m.items) {
		return Item{}, false
	}
	return m.items[m.selected], true
}

func (m *Model) Value() string {
	item, ok := m.Selected()
	if !ok {
		return ""
	}
	if strings.TrimSpace(item.Value) != "" {
		return item.Value
	}
	return item.Label
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "enter", " ":
		if m.open {
			m.selected = m.cursor
			m.open = false
		} else {
			m.Open()
		}
	case "esc":
		m.open = false
	case "down", "j":
		if !m.open {
			m.Open()
			break
		}
		if m.cursor < len(m.items)-1 {
			m.cursor++
		}
	case "up", "k":
		if !m.open {
			m.Open()
			break
		}
		if m.cursor > 0 {
			m.cursor--
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	w := m.width
	if w <= 0 {
		w = 28
	}

	head := m.placeholder
	if item, ok := m.Selected(); ok {
		head = item.Label
	}
	if len(m.items) == 0 {
		head = "No options"
	}
	caret := " v"
	if m.open {
		caret = " ^"
	}
	body := []string{styles.Row(pick(t.Input.BG, t.Surface.Elevated), pick(t.Input.FG, t.Text.Primary), w, head+caret)}
	if !m.open || len(m.items) == 0 {
		return tea.NewView(strings.Join(body, "\n"))
	}

	count := m.visibleCount()
	for i := 0; i < count; i++ {
		idx := i
		line := "  " + m.items[idx].Label
		if idx == m.selected {
			line = "* " + m.items[idx].Label
		}
		if idx == m.cursor {
			body = append(body, styles.Row(pick(t.Selection.BG, t.Border.Focus), pick(t.Selection.FG, t.Text.Inverse), w, line))
			continue
		}
		body = append(body, styles.Row(pick(t.Surface.Panel, t.Surface.Elevated), pick(t.Text.Primary, t.Text.Primary), w, line))
	}
	return tea.NewView(strings.Join(body, "\n"))
}

func (m *Model) visibleCount() int {
	if m.height <= 1 {
		return len(m.items)
	}
	maxRows := m.height - 1
	if maxRows > len(m.items) {
		maxRows = len(m.items)
	}
	if maxRows < 0 {
		maxRows = 0
	}
	return maxRows
}

func (m *Model) String() string {
	item, ok := m.Selected()
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s", item.Label)
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
