package focus

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
)

type Manager struct {
	ring    []core.Focusable
	idx     int
	enabled bool
	wrap    bool
	next    key.Binding
	prev    key.Binding
}

type Option func(*Manager)

func Ring(components ...core.Focusable) Option {
	return func(m *Manager) {
		m.SetRing(components...)
	}
}

func Keys(next, prev key.Binding) Option {
	return func(m *Manager) {
		m.next = next
		m.prev = prev
	}
}

func Enabled(v bool) Option {
	return func(m *Manager) {
		m.enabled = v
	}
}

func Wrap(v bool) Option {
	return func(m *Manager) {
		m.wrap = v
	}
}

func New(opts ...Option) *Manager {
	m := &Manager{
		enabled: true,
		wrap:    true,
		next:    key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next panel")),
		prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev panel")),
	}
	for _, opt := range opts {
		opt(m)
	}
	_ = m.applyFocus(-1)
	return m
}

func (m *Manager) Init() tea.Cmd { return nil }

func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.enabled {
		return m, nil
	}
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(keyMsg, m.next):
			return m, m.Next()
		case key.Matches(keyMsg, m.prev):
			return m, m.Prev()
		}
	}
	return m, nil
}

func (m *Manager) View() tea.View { return tea.NewView("") }

func (m *Manager) Next() tea.Cmd {
	return m.FocusBy(1)
}

func (m *Manager) Prev() tea.Cmd {
	return m.FocusBy(-1)
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

func (m *Manager) SetEnabled(v bool) {
	m.enabled = v
}

func (m *Manager) SetWrap(v bool) {
	m.wrap = v
}

func (m *Manager) SetRing(components ...core.Focusable) tea.Cmd {
	m.ring = sanitizeRing(components)
	if len(m.ring) == 0 {
		m.idx = 0
		return nil
	}
	if m.idx >= len(m.ring) {
		m.idx = len(m.ring) - 1
	}
	if m.idx < 0 {
		m.idx = 0
	}
	return m.applyFocus(-1)
}

func (m *Manager) SetIndex(idx int) tea.Cmd {
	if len(m.ring) == 0 {
		m.idx = 0
		return nil
	}
	next := idx
	if m.wrap {
		next = wrapIndex(idx, len(m.ring))
	} else {
		next = clamp(idx, 0, len(m.ring)-1)
	}
	if next == m.idx {
		return nil
	}
	from := m.idx
	m.idx = next
	return m.applyFocus(from)
}

func (m *Manager) FocusBy(delta int) tea.Cmd {
	if !m.enabled || len(m.ring) == 0 || delta == 0 {
		return nil
	}
	return m.SetIndex(m.idx + delta)
}

func (m *Manager) applyFocus(from int) tea.Cmd {
	for i, c := range m.ring {
		if c == nil {
			continue
		}
		if i == m.idx {
			c.Focus()
			continue
		}
		c.Blur()
	}
	if from == m.idx {
		return nil
	}
	to := m.idx
	return func() tea.Msg { return FocusChangedMsg{From: from, To: to} }
}

func sanitizeRing(components []core.Focusable) []core.Focusable {
	out := make([]core.Focusable, 0, len(components))
	for _, c := range components {
		if c != nil {
			out = append(out, c)
		}
	}
	return out
}

func clamp(v, minV, maxV int) int {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func wrapIndex(idx, length int) int {
	if length <= 0 {
		return 0
	}
	v := idx % length
	if v < 0 {
		v += length
	}
	return v
}
