// Brick: Select:
// +-----------------------------------+
// | current value v                    |
// | option 1                           |
// | option 2                           |
// +-----------------------------------+
// Single-choice picker.
// Package select provides a themed single-choice picker backed by bubbles/list.
// Copy this file into your project: bento add select
package selectx

import (
	"fmt"
	"io"
	"strings"

	bubbleskey "charm.land/bubbles/v2/key"
	bubbleslist "charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

type Item struct {
	Label string
	Value string
}

type KeyMap struct {
	Toggle bubbleskey.Binding
	Close  bubbleskey.Binding
	Up     bubbleskey.Binding
	Down   bubbleskey.Binding
	Choose bubbleskey.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Toggle: bubbleskey.NewBinding(bubbleskey.WithKeys("enter", " "), bubbleskey.WithHelp("enter", "open")),
		Close:  bubbleskey.NewBinding(bubbleskey.WithKeys("esc"), bubbleskey.WithHelp("esc", "close")),
		Up:     bubbleskey.NewBinding(bubbleskey.WithKeys("up", "k"), bubbleskey.WithHelp("up", "up")),
		Down:   bubbleskey.NewBinding(bubbleskey.WithKeys("down", "j"), bubbleskey.WithHelp("down", "down")),
		Choose: bubbleskey.NewBinding(bubbleskey.WithKeys("enter"), bubbleskey.WithHelp("enter", "choose")),
	}
}

type selectItem struct {
	label string
	value string
}

func (i selectItem) FilterValue() string { return i.label }

// selectDelegate renders items with theme-aware colors baked in at the delegate
// level — not via post-hoc string scanning. The owner pointer gives access to
// current theme tokens at render time.
type selectDelegate struct{ owner *Model }

func (d selectDelegate) Height() int  { return 1 }
func (d selectDelegate) Spacing() int { return 0 }
func (d selectDelegate) Update(msg tea.Msg, m *bubbleslist.Model) tea.Cmd {
	return nil
}

func (d selectDelegate) Render(w io.Writer, m bubbleslist.Model, index int, item bubbleslist.Item) {
	opt, ok := item.(selectItem)
	if !ok {
		return
	}
	t := theme.CurrentTheme()
	width := d.owner.width
	if width <= 0 {
		width = 28
	}
	isSelected := index == m.Index()
	prefix := "  "
	if isSelected {
		prefix = "> "
	}
	content := prefix + opt.label
	if isSelected {
		line := lipgloss.NewStyle().
			Background(t.SelectionBG()).
			Foreground(t.SelectionFG()).
			Bold(true).
			Width(width).
			Render(content)
		_, _ = io.WriteString(w, line)
	} else {
		line := lipgloss.NewStyle().
			Background(t.BackgroundPanel()).
			Foreground(t.Text()).
			Width(width).
			Render(content)
		_, _ = io.WriteString(w, line)
	}
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
	inner       bubbleslist.Model
	keys        KeyMap
}

func New(items ...Item) *Model {
	m := &Model{
		items:       append([]Item(nil), items...),
		selected:    -1,
		placeholder: "Select...",
		keys:        DefaultKeyMap(),
	}
	inner := bubbleslist.New([]bubbleslist.Item{}, selectDelegate{owner: m}, 24, 4)
	inner.SetShowTitle(false)
	inner.SetShowFilter(false)
	inner.SetShowStatusBar(false)
	inner.SetShowPagination(false)
	inner.SetShowHelp(false)
	inner.SetFilteringEnabled(false)
	inner.DisableQuitKeybindings()
	m.inner = inner
	m.syncItems()
	return m
}

func (m *Model) SetItems(items []Item) {
	m.items = append([]Item(nil), items...)
	if len(m.items) == 0 {
		m.cursor = 0
		m.selected = -1
		m.syncItems()
		return
	}
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.selected >= len(m.items) {
		m.selected = -1
	}
	m.syncItems()
}

func (m *Model) SetPlaceholder(v string) { m.placeholder = v }
func (m *Model) Focus()                  { m.focused = true }
func (m *Model) Blur()                   { m.focused = false; m.open = false }
func (m *Model) IsFocused() bool         { return m.focused }
func (m *Model) Open() {
	if len(m.items) > 0 {
		m.open = true
		m.syncItems()
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
	m.syncItems()
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

	if !m.open {
		if bubbleskey.Matches(k, m.keys.Toggle, m.keys.Choose) {
			m.Open()
		}
		return m, nil
	}

	if bubbleskey.Matches(k, m.keys.Close) {
		m.Close()
		return m, nil
	}

	updated, cmd := m.inner.Update(msg)
	m.inner = updated
	m.cursor = m.inner.Index()
	if bubbleskey.Matches(k, m.keys.Choose) {
		m.selected = m.cursor
		m.open = false
	}
	return m, cmd
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
	// Header row: uses input surface tokens.
	headerRow := styles.Row(t.InputBG(), t.InputFG(), w, head+caret)

	if !m.open || len(m.items) == 0 {
		return tea.NewView(headerRow)
	}

	// Delegate renders each row with correct bg/fg already applied —
	// no post-hoc string scanning needed. Just consume the output directly.
	menu := m.inner.View()
	rows := []string{headerRow}
	for _, line := range strings.Split(menu, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		rows = append(rows, line)
	}
	return tea.NewView(strings.Join(rows, "\n"))
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

func (m *Model) syncItems() {
	items := make([]bubbleslist.Item, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, selectItem{label: item.Label, value: item.Value})
	}
	if cmd := m.inner.SetItems(items); cmd != nil {
		_ = cmd
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	if m.cursor >= len(items) && len(items) > 0 {
		m.cursor = len(items) - 1
	}
	m.inner.Select(m.cursor)
	w := max(10, m.width)
	h := max(1, m.visibleCount())
	m.inner.SetSize(w, h)
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
