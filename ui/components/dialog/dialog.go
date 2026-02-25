package dialog

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/surface"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/styles"
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

func New() *Manager { return &Manager{theme: theme.Preset(theme.DefaultName)} }

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
		t = theme.Preset(theme.DefaultName)
	}
	content := strings.Join([]string{text, "", "Enter confirm  Esc cancel"}, "\n")
	view := renderDialogFrame(c.DialogTitle, content, max(48, c.width/2), 0, t)
	return tea.NewView(view)
}

func (c Confirm) SetSize(width, height int) {
	c.width = width
	c.height = height
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
		t = theme.Preset(theme.DefaultName)
	}
	body := ""
	if c.Content != nil {
		body = core.ViewString(c.Content.View())
	}
	view := renderDialogFrame(c.DialogTitle, body, max(50, c.Width), max(12, c.Height), t)
	return tea.NewView(view)
}

func (c Custom) SetSize(width, height int) {
	if c.Width == 0 {
		c.Width = width / 2
	}
	if c.Height == 0 {
		c.Height = height / 2
	}
	maxWidth := max(20, width-4)
	maxHeight := max(8, height-4)
	minWidth := min(36, maxWidth)
	minHeight := min(12, maxHeight)
	c.Width = clamp(c.Width, minWidth, maxWidth)
	c.Height = clamp(c.Height, minHeight, maxHeight)
	if s, ok := c.Content.(core.Sizeable); ok {
		s.SetSize(max(1, c.Width-4), max(1, c.Height-4))
	}
}

func (c Custom) Title() string { return c.DialogTitle }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func renderDialogFrame(title, content string, width, height int, t theme.Theme) string {
	sys := styles.New(t)
	headerTitle := title
	if strings.TrimSpace(headerTitle) == "" {
		headerTitle = "Dialog"
	}
	if width <= 0 {
		width = 48
	}
	if height <= 0 {
		height = 14
	}
	innerWidth := max(1, width-4)
	innerHeight := max(1, height-2)

	rightWidth := 3
	leftWidth := max(1, innerWidth-rightWidth-1)
	header := sys.DialogHeader().Render(surface.FitWidth(headerTitle, leftWidth)) + " " + sys.DialogEscHint().Render(surface.FitWidth("esc", rightWidth))
	header = surface.FitWidth(header, innerWidth)

	body := strings.TrimRight(content, "\n")
	if strings.TrimSpace(body) == "" {
		body = " "
	}
	bodyLines := clipLines(body, innerWidth)
	bodyMax := max(1, innerHeight-2)
	if len(bodyLines) > bodyMax {
		bodyLines = bodyLines[:bodyMax]
	}
	joined := strings.Join(append([]string{header, ""}, bodyLines...), "\n")
	style := sys.DialogFrame().Width(width).Height(height)
	return style.Render(joined)
}

func clipLines(content string, width int) []string {
	lines := strings.Split(content, "\n")
	clipped := make([]string, 0, len(lines))
	for _, line := range lines {
		clipped = append(clipped, surface.FitWidth(line, width))
	}
	return clipped
}

func clamp(v, minV, maxV int) int {
	if maxV < minV {
		return minV
	}
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
