// Brick: Tabs
// +-----------------------------------+
// | [Tab A] [Tab B] [Tab C]          |
// +-----------------------------------+
// Keyboard-navigable tab row.
// Copy this file into your project: bento add tabs
package tabs

import (
	"strings"

	bubbleskey "charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/paginator"
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

type Tab struct {
	ID    string
	Label string
}

type KeyMap struct {
	Prev bubbleskey.Binding
	Next bubbleskey.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Prev: bubbleskey.NewBinding(bubbleskey.WithKeys("left", "h"), bubbleskey.WithHelp("left", "prev tab")),
		Next: bubbleskey.NewBinding(bubbleskey.WithKeys("right", "l"), bubbleskey.WithHelp("right", "next tab")),
	}
}

type Model struct {
	tabs    []Tab
	active  int
	focused bool
	width   int
	keys    KeyMap
	pager   paginator.Model
	theme   theme.Theme // nil = use theme.CurrentTheme()
}

func New(tabs ...Tab) *Model {
	p := paginator.New(paginator.WithTotalPages(max(1, len(tabs))), paginator.WithPerPage(1))
	p.Type = paginator.Dots
	return &Model{tabs: append([]Tab(nil), tabs...), keys: DefaultKeyMap(), pager: p}
}

func (m *Model) SetTabs(tabs []Tab) {
	m.tabs = append([]Tab(nil), tabs...)
	if len(m.tabs) == 0 {
		m.active = 0
		m.pager.TotalPages = 1
		m.pager.Page = 0
		return
	}
	if m.active >= len(m.tabs) {
		m.active = len(m.tabs) - 1
	}
	m.pager.TotalPages = len(m.tabs)
	m.pager.Page = m.active
}

func (m *Model) SetActive(i int) {
	if i >= 0 && i < len(m.tabs) {
		m.active = i
		m.pager.Page = i
	}
}

func (m *Model) Active() int          { return m.active }
func (m *Model) Focus()               { m.focused = true }
func (m *Model) Blur()                { m.focused = false }
func (m *Model) IsFocused() bool      { return m.focused }
func (m *Model) Init() tea.Cmd        { return nil }
func (m *Model) SetSize(width, _ int) { m.width = width }
func (m *Model) GetSize() (int, int)  { return m.width, 1 }

// SetTheme updates the theme. Call on ThemeChangedMsg.
func (m *Model) SetTheme(t theme.Theme) { m.theme = t }

func (m *Model) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.focused || len(m.tabs) == 0 {
		return m, nil
	}
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch {
	case bubbleskey.Matches(k, m.keys.Prev):
		if m.active > 0 {
			m.active--
			m.pager.Page = m.active
		}
	case bubbleskey.Matches(k, m.keys.Next):
		if m.active < len(m.tabs)-1 {
			m.active++
			m.pager.Page = m.active
		}
	}
	return m, nil
}

func (m *Model) View() tea.View {
	t := m.activeTheme()
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
		} else {
			parts = append(parts, " "+label+" ")
		}
	}
	line := strings.Join(parts, " ")
	if m.width > 0 {
		bg := t.BackgroundPanel()
		fg := t.Text()
		if m.focused {
			bg = t.BackgroundInteractive()
		}
		if strings.TrimSpace(line) == "" {
			line = m.pager.View()
		}
		return tea.NewView(styles.Row(bg, fg, m.width, line))
	}
	return tea.NewView(line)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
