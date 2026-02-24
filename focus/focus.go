package focus

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

type Manager struct {
	ring []core.Focusable
	idx  int
	next key.Binding
	prev key.Binding
}

type Option func(*Manager)

func Ring(components ...core.Focusable) Option {
	return func(m *Manager) {
		m.ring = components
		m.idx = 0
		m.applyFocus()
	}
}

func Keys(next, prev key.Binding) Option {
	return func(m *Manager) {
		m.next = next
		m.prev = prev
	}
}

func New(opts ...Option) *Manager {
	m := &Manager{
		next: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next panel")),
		prev: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev panel")),
	}
	for _, opt := range opts {
		opt(m)
	}
	m.applyFocus()
	return m
}

func (m *Manager) Init() tea.Cmd { return nil }

func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, m.next):
			m.Next()
		case key.Matches(keyMsg, m.prev):
			m.Prev()
		}
	}
	return m, nil
}

func (m *Manager) View() tea.View { return tea.NewView("") }

func (m *Manager) Next() {
	if len(m.ring) == 0 {
		return
	}
	m.idx = (m.idx + 1) % len(m.ring)
	m.applyFocus()
}

func (m *Manager) Prev() {
	if len(m.ring) == 0 {
		return
	}
	m.idx--
	if m.idx < 0 {
		m.idx = len(m.ring) - 1
	}
	m.applyFocus()
}

func (m *Manager) Focused() core.Focusable {
	if len(m.ring) == 0 {
		return nil
	}
	return m.ring[m.idx]
}

func (m *Manager) Bindings() []key.Binding {
	return []key.Binding{m.next, m.prev}
}

func (m *Manager) applyFocus() {
	for i, c := range m.ring {
		if i == m.idx {
			c.Focus()
			continue
		}
		c.Blur()
	}
}
