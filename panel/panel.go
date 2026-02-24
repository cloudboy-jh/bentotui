package panel

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

type Option func(*Model)

type Model struct {
	title      string
	content    core.Component
	scrollable bool
	width      int
	height     int
}

func New(opts ...Option) *Model {
	m := &Model{}
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
	lines := strings.Split(body, "\n")
	innerWidth := 0
	for _, line := range lines {
		if len(line) > innerWidth {
			innerWidth = len(line)
		}
	}
	if m.width > 2 && m.width-2 > innerWidth {
		innerWidth = m.width - 2
	}
	border := "+" + strings.Repeat("-", innerWidth) + "+"
	out := []string{border}
	if m.title != "" {
		title := m.title
		if len(title) > innerWidth {
			title = title[:innerWidth]
		}
		out = append(out, "|"+padRight(title, innerWidth)+"|")
		out = append(out, border)
	}
	for _, line := range lines {
		out = append(out, "|"+padRight(line, innerWidth)+"|")
	}
	out = append(out, border)
	return tea.NewView(strings.Join(out, "\n"))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if s, ok := m.content.(core.Sizeable); ok {
		s.SetSize(width-2, height-2)
	}
}

func (m *Model) GetSize() (width, height int) { return m.width, m.height }

func padRight(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}
