// Brick: Tabs:
// +-----------------------------------+
// | [Tab A] [Tab B] [Tab C]          |
// +-----------------------------------+
// Keyboard-navigable tab row.
// Package tabs provides a keyboard-navigable tab row.
// Copy this file into your project: bento add tabs
package tabs

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Tab struct {
	ID    string
	Label string
}

type Model struct {
	tabs    []Tab
	active  int
	focused bool
	width   int
}

func New(tabs ...Tab) *Model {
	return &Model{tabs: append([]Tab(nil), tabs...)}
}

func (m *Model) SetTabs(tabs []Tab) {
	m.tabs = append([]Tab(nil), tabs...)
	if len(m.tabs) == 0 {
		m.active = 0
		return
	}
	if m.active >= len(m.tabs) {
		m.active = len(m.tabs) - 1
	}
}

func (m *Model) SetActive(i int) {
	if i >= 0 && i < len(m.tabs) {
		m.active = i
	}
}

func (m *Model) Active() int          { return m.active }
func (m *Model) Focus()               { m.focused = true }
func (m *Model) Blur()                { m.focused = false }
func (m *Model) IsFocused() bool      { return m.focused }
func (m *Model) Init() tea.Cmd        { return nil }
func (m *Model) SetSize(width, _ int) { m.width = width }
func (m *Model) GetSize() (int, int)  { return m.width, 1 }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.focused || len(m.tabs) == 0 {
		return m, nil
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch k.String() {
	case "left", "h":
		if m.active > 0 {
			m.active--
		}
	case "right", "l":
		if m.active < len(m.tabs)-1 {
			m.active++
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	if len(m.tabs) == 0 {
		return tea.NewView("")
	}
	parts := make([]string, 0, len(m.tabs))
	for i, tab := range m.tabs {
		label := tab.Label
		if strings.TrimSpace(label) == "" {
			label = tab.ID
		}
		if i == m.active {
			parts = append(parts, "["+label+"]")
			continue
		}
		parts = append(parts, " "+label+" ")
	}
	line := strings.Join(parts, " ")
	if m.width > 0 {
		bg := pick(t.Surface.Panel, t.Surface.Elevated)
		fg := pick(t.Text.Primary, t.Text.Primary)
		if m.focused {
			bg = pick(t.Surface.Interactive, t.Surface.Panel)
		}
		return tea.NewView(styles.Row(bg, fg, m.width, line))
	}
	return tea.NewView(line)
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
