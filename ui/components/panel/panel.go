package panel

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/surface"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/ui/styles"
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
	m := &Model{theme: theme.Preset(theme.DefaultName)}
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
	sys := styles.New(m.theme)
	body := ""
	if m.content != nil {
		body = core.ViewString(m.content.View())
	}
	outerWidth := m.width
	if outerWidth <= 0 {
		outerWidth = max(30, maxLineWidth(body)+4)
	}
	outerHeight := m.height
	if outerHeight <= 0 {
		outerHeight = max(8, len(strings.Split(body, "\n"))+4)
	}

	innerWidth := max(0, outerWidth-2)
	innerHeight := max(0, outerHeight-2)
	rows := make([]string, 0, innerHeight)
	contentRows := strings.Split(body, "\n")
	contentStart := 0
	if m.title != "" && innerHeight > 0 {
		rows = append(rows, renderTitleRow(m.title, innerWidth, m.theme, m.focused))
		contentStart = 1
	}
	for len(rows) < innerHeight {
		idx := len(rows) - contentStart
		if idx >= 0 && idx < len(contentRows) {
			rows = append(rows, padContentRow(contentRows[idx], innerWidth))
			continue
		}
		rows = append(rows, padContentRow("", innerWidth))
	}

	boxStyle := sys.PanelFrame(m.focused).
		Width(outerWidth).
		Height(outerHeight)

	return tea.NewView(boxStyle.Render(strings.Join(rows, "\n")))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	contentWidth := max(0, width-4)
	contentHeight := max(0, height-2)
	if m.title != "" {
		contentHeight = max(0, contentHeight-1)
	}
	if s, ok := m.content.(core.Sizeable); ok {
		s.SetSize(contentWidth, contentHeight)
	}
}

func (m *Model) GetSize() (width, height int) { return m.width, m.height }

func (m *Model) SetTheme(t theme.Theme) { m.theme = t }

func (m *Model) Focus() { m.focused = true }

func (m *Model) Blur() { m.focused = false }

func (m *Model) IsFocused() bool { return m.focused }

func renderTitleRow(title string, width int, t theme.Theme, focused bool) string {
	if width <= 0 {
		return ""
	}
	sys := styles.New(t)
	badge := sys.PanelTitleChip(focused).Render(title)
	return sys.PanelTitleBar(focused).
		Width(width).
		Render(lipgloss.PlaceHorizontal(width, lipgloss.Left, badge))
}

func fitWidth(s string, width int) string {
	return surface.FitWidth(s, width)
}

func padContentRow(s string, innerWidth int) string {
	if innerWidth <= 0 {
		return ""
	}
	if innerWidth == 1 {
		return " "
	}
	if innerWidth == 2 {
		return "  "
	}
	return " " + fitWidth(s, innerWidth-2) + " "
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
