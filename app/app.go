package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/dialog"
	"github.com/cloudboy-jh/bentotui/router"
	"github.com/cloudboy-jh/bentotui/statusbar"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Option func(*Model)

type Model struct {
	router        *router.Model
	dialogs       *dialog.Manager
	status        *statusbar.Model
	theme         theme.Theme
	showStatusBar bool
	width         int
	height        int
}

func New(opts ...Option) *Model {
	m := &Model{
		router:        router.New(),
		dialogs:       dialog.New(),
		status:        statusbar.New(),
		theme:         theme.Preset("amber"),
		showStatusBar: true,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithTheme(t theme.Theme) Option {
	return func(m *Model) { m.theme = t }
}

func WithPages(routes ...router.Route) Option {
	return func(m *Model) {
		m.router = router.New(routes...)
	}
}

func WithStatusBar(v bool) Option {
	return func(m *Model) { m.showStatusBar = v }
}

func WithStatus(model *statusbar.Model) Option {
	return func(m *Model) {
		if model != nil {
			m.status = model
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.router.Init(), m.dialogs.Init(), m.status.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = v.Width
		m.height = v.Height
		bodyHeight := m.height
		if m.showStatusBar {
			bodyHeight--
		}
		if bodyHeight < 0 {
			bodyHeight = 0
		}
		m.router.SetSize(m.width, bodyHeight)
		m.dialogs.SetSize(m.width, bodyHeight)
		m.status.SetSize(m.width, 1)
	}

	// Dialog manager processes all messages to open/close overlays.
	_, dialogCmd := m.dialogs.Update(msg)

	var pageCmd tea.Cmd
	if !m.dialogs.IsOpen() {
		_, pageCmd = m.router.Update(msg)
	}

	_, statusCmd := m.status.Update(msg)

	return m, tea.Batch(dialogCmd, pageCmd, statusCmd)
}

func (m *Model) View() tea.View {
	body := core.ViewString(m.router.View())
	if m.showStatusBar {
		body = strings.Join([]string{body, core.ViewString(m.status.View())}, "\n")
	}
	if m.dialogs.IsOpen() {
		return tea.NewView(strings.Join([]string{body, core.ViewString(m.dialogs.View())}, "\n"))
	}
	return tea.NewView(body)
}

func (m *Model) Router() *router.Model {
	return m.router
}

func (m *Model) Dialogs() *dialog.Manager {
	return m.dialogs
}
