package panel

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Option func(*Model)

type Model struct {
	title      string
	content    core.Component
	scrollable bool
	focused    bool
	theme      theme.Theme
	width      int
	height     int
}

func New(opts ...Option) *Model {
	m := &Model{theme: theme.Preset("amber")}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func Title(title string) Option {
	return func(m *Model) { m.title = title }
}

func Border(_ any) Option { return func(m *Model) {} }

func Content(content core.Component) Option {
	return func(m *Model) { m.content = content }
}

func Scrollable(v bool) Option {
	return func(m *Model) { m.scrollable = v }
}

func Theme(t theme.Theme) Option {
	return func(m *Model) { m.theme = t }
}

func (m *Model) Init() tea.Cmd {
	if m.content == nil {
		return nil
	}
	return m.content.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.content == nil {
		return m, nil
	}
	updated, cmd := m.content.Update(msg)
	if next, ok := updated.(core.Component); ok {
		m.content = next
	}
	return m, cmd
}

func (m *Model) View() tea.View {
	body := ""
	if m.content != nil {
		body = core.ViewString(m.content.View())
	}
	outerWidth := max(20, m.width)
	if m.width <= 0 {
		outerWidth = max(30, maxLineWidth(body)+4)
	}
	outerHeight := max(5, m.height)
	if m.height <= 0 {
		outerHeight = max(8, len(strings.Split(body, "\n"))+4)
	}

	innerWidth := max(0, outerWidth-2)
	innerHeight := max(0, outerHeight-2)
	rows := make([]string, 0, innerHeight)
	contentRows := strings.Split(body, "\n")
	contentStart := 0
	if m.title != "" && innerHeight > 0 {
		rows = append(rows, renderTitleRow(m.title, innerWidth, m.theme))
		contentStart = 1
	}
	for len(rows) < innerHeight {
		idx := len(rows) - contentStart
		if idx >= 0 && idx < len(contentRows) {
			rows = append(rows, fitWidth(contentRows[idx], innerWidth))
			continue
		}
		rows = append(rows, fitWidth("", innerWidth))
	}

	borderColor := m.theme.Border
	if m.focused {
		borderColor = m.theme.BorderFocused
	}
	boxStyle := lipgloss.NewStyle().
		Width(outerWidth).
		Height(outerHeight).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Background(lipgloss.Color(m.theme.Surface)).
		Foreground(lipgloss.Color(m.theme.Text))

	return tea.NewView(boxStyle.Render(strings.Join(rows, "\n")))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if s, ok := m.content.(core.Sizeable); ok {
		s.SetSize(max(0, width-2), max(0, height-2))
	}
}

func (m *Model) GetSize() (width, height int) { return m.width, m.height }

func (m *Model) Focus() { m.focused = true }

func (m *Model) Blur() { m.focused = false }

func (m *Model) IsFocused() bool { return m.focused }

func renderTitleRow(title string, width int, t theme.Theme) string {
	if width <= 0 {
		return ""
	}
	badge := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.TitleText)).
		Background(lipgloss.Color(t.TitleBG)).
		Bold(true).
		Padding(0, 1).
		Render(title)
	base := lipgloss.NewStyle().
		Width(width).
		Foreground(lipgloss.Color(t.Muted)).
		Background(lipgloss.Color(t.SurfaceMuted)).
		Render(lipgloss.PlaceHorizontal(width, lipgloss.Left, badge))
	return base
}

func fitWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	s = lipgloss.NewStyle().MaxWidth(width).Render(s)
	return lipgloss.PlaceHorizontal(width, lipgloss.Left, s)
}

func maxLineWidth(s string) int {
	lines := strings.Split(s, "\n")
	m := 0
	for _, line := range lines {
		w := lipgloss.Width(line)
		if w > m {
			m = w
		}
	}
	return m
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
