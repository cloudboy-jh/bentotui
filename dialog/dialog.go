package dialog

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Dialog interface {
	core.Component
	SetSize(width, height int)
	Title() string
}

type OpenMsg struct {
	Dialog Dialog
}

type CloseMsg struct{}

func Open(dialog Dialog) tea.Msg { return OpenMsg{Dialog: dialog} }
func Close() tea.Msg             { return CloseMsg{} }

type Manager struct {
	active Dialog
	theme  theme.Theme
	width  int
	height int
}

func New() *Manager { return &Manager{theme: theme.Preset("amber")} }

func (m *Manager) Init() tea.Cmd { return nil }

func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case OpenMsg:
		switch d := v.Dialog.(type) {
		case Confirm:
			d.theme = m.theme
			m.active = d
		case *Confirm:
			d.theme = m.theme
			m.active = d
		case Custom:
			d.theme = m.theme
			m.active = d
		case *Custom:
			d.theme = m.theme
			m.active = d
		default:
			m.active = v.Dialog
		}
		if m.active != nil {
			m.active.SetSize(m.width, m.height)
		}
		return m, nil
	case CloseMsg:
		m.active = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(v.Width, v.Height)
		return m, nil
	}

	if m.active != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				m.active = nil
				return m, nil
			case "enter":
				switch active := m.active.(type) {
				case Confirm:
					m.active = nil
					if active.OnConfirm != nil {
						return m, func() tea.Msg { return active.OnConfirm() }
					}
					return m, nil
				case *Confirm:
					m.active = nil
					if active != nil && active.OnConfirm != nil {
						return m, func() tea.Msg { return active.OnConfirm() }
					}
					return m, nil
				default:
					m.active = nil
					return m, nil
				}
			}
		}
	}

	if m.active == nil {
		return m, nil
	}
	updated, cmd := m.active.Update(msg)
	if next, ok := updated.(Dialog); ok {
		m.active = next
	}
	return m, cmd
}

func (m *Manager) View() tea.View {
	if m.active == nil {
		return tea.NewView("")
	}
	view := core.ViewString(m.active.View())
	return tea.NewView(view)
}

func (m *Manager) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.active != nil {
		m.active.SetSize(width, height)
	}
}

func (m *Manager) IsOpen() bool { return m.active != nil }

func (m *Manager) SetTheme(t theme.Theme) {
	m.theme = t
}

type Confirm struct {
	DialogTitle string
	Message     string
	OnConfirm   func() tea.Msg
	theme       theme.Theme
	width       int
	height      int
}

func (c Confirm) Init() tea.Cmd { return nil }

func (c Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return c, nil
}

func (c Confirm) View() tea.View {
	text := c.Message
	if strings.TrimSpace(text) == "" {
		text = "Confirm?"
	}
	t := c.theme
	if t.Accent == "" {
		t = theme.Preset("amber")
	}
	content := strings.Join([]string{
		c.DialogTitle,
		"",
		text,
		"",
		"Press Enter to confirm, Esc to cancel",
	}, "\n")
	view := renderDialogFrame(content, max(42, c.width/2), 0, t)
	return tea.NewView(view)
}

func (c Confirm) SetSize(width, height int) {
	_ = width
	_ = height
}

func (c Confirm) Title() string { return c.DialogTitle }

type Custom struct {
	DialogTitle string
	Content     core.Component
	Width       int
	Height      int
	theme       theme.Theme
}

func (c Custom) Init() tea.Cmd {
	if c.Content == nil {
		return nil
	}
	return c.Content.Init()
}

func (c Custom) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if c.Content == nil {
		return c, nil
	}
	updated, cmd := c.Content.Update(msg)
	if next, ok := updated.(core.Component); ok {
		c.Content = next
	}
	return c, cmd
}

func (c Custom) View() tea.View {
	t := c.theme
	if t.Accent == "" {
		t = theme.Preset("amber")
	}
	body := ""
	if c.Content != nil {
		body = core.ViewString(c.Content.View())
	}
	if c.DialogTitle != "" {
		body = c.DialogTitle + "\n\n" + body
	}
	view := renderDialogFrame(body, max(42, c.Width), max(10, c.Height), t)
	return tea.NewView(view)
}

func (c Custom) SetSize(width, height int) {
	if c.Width == 0 {
		c.Width = width / 2
	}
	if c.Height == 0 {
		c.Height = height / 2
	}
	if s, ok := c.Content.(core.Sizeable); ok {
		s.SetSize(c.Width-4, c.Height-4)
	}
}

func (c Custom) Title() string { return c.DialogTitle }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func renderDialogFrame(content string, width, height int, t theme.Theme) string {
	style := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.DialogBorder)).
		Background(lipgloss.Color(t.DialogBG)).
		Foreground(lipgloss.Color(t.DialogText)).
		Padding(1, 2)
	if height > 0 {
		style = style.Height(height)
	}
	return style.Render(content)
}
