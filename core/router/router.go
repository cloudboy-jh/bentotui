package router

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

type Factory func() core.Page

type Route struct {
	Name    string
	Factory Factory
}

func Page(name string, factory Factory) Route {
	return Route{Name: name, Factory: factory}
}

type Model struct {
	routes  map[string]Factory
	cache   map[string]core.Page
	current string
	width   int
	height  int
}

func New(routes ...Route) *Model {
	m := &Model{
		routes: make(map[string]Factory, len(routes)),
		cache:  make(map[string]core.Page, len(routes)),
	}
	for i, r := range routes {
		m.routes[r.Name] = r.Factory
		if i == 0 {
			m.current = r.Name
		}
	}
	_ = m.ensureCurrent()
	return m
}

func Navigate(page string) tea.Msg {
	return core.Navigate(page)
}

func (m *Model) Init() tea.Cmd {
	if page := m.Current(); page != nil {
		return page.Init()
	}
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case core.NavigateMsg:
		if _, ok := m.routes[v.Page]; ok {
			m.current = v.Page
			_ = m.ensureCurrent()
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(v.Width, v.Height)
		return m, nil
	}

	page := m.Current()
	if page == nil {
		return m, nil
	}

	updated, cmd := page.Update(msg)
	if next, ok := updated.(core.Page); ok {
		m.cache[m.current] = next
	}

	return m, cmd
}

func (m *Model) View() tea.View {
	if page := m.Current(); page != nil {
		return page.View()
	}
	return tea.NewView("")
}

func (m *Model) Current() core.Page {
	_ = m.ensureCurrent()
	return m.cache[m.current]
}

func (m *Model) CurrentName() string {
	return m.current
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if page := m.Current(); page != nil {
		page.SetSize(width, height)
	}
}

func (m *Model) GetSize() (width, height int) {
	return m.width, m.height
}

func (m *Model) ensureCurrent() error {
	if m.current == "" {
		return nil
	}
	if _, ok := m.cache[m.current]; ok {
		return nil
	}
	factory, ok := m.routes[m.current]
	if !ok {
		return fmt.Errorf("route %q not registered", m.current)
	}
	page := factory()
	page.SetSize(m.width, m.height)
	m.cache[m.current] = page
	return nil
}
