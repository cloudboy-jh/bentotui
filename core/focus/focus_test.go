package focus

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

type focusSpy struct {
	focused bool
}

func (s *focusSpy) Init() tea.Cmd                           { return nil }
func (s *focusSpy) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s *focusSpy) View() tea.View                          { return tea.NewView("") }
func (s *focusSpy) Focus()                                  { s.focused = true }
func (s *focusSpy) Blur()                                   { s.focused = false }
func (s *focusSpy) IsFocused() bool                         { return s.focused }

func TestSetRingIgnoresNilEntries(t *testing.T) {
	a := &focusSpy{}
	b := &focusSpy{}
	m := New()
	_ = m.SetRing(nil, a, nil, b)
	if m.Focused() != a {
		t.Fatal("expected first non-nil entry to be focused")
	}
}

func TestSetIndexWrapAndClamp(t *testing.T) {
	a := &focusSpy{}
	b := &focusSpy{}
	m := New(Ring(a, b))
	msg := runCmd(m.SetIndex(-1))
	changed, ok := msg.(FocusChangedMsg)
	if !ok || changed.To != 1 {
		t.Fatalf("expected wrapped index to 1, got %#v", msg)
	}

	m.SetWrap(false)
	_ = m.SetIndex(50)
	if m.Focused() != b {
		t.Fatal("expected clamped index to remain last entry")
	}
}

func TestFocusByDisabledReturnsNoChange(t *testing.T) {
	a := &focusSpy{}
	b := &focusSpy{}
	m := New(Ring(a, b), Enabled(false))
	if cmd := m.FocusBy(1); cmd != nil {
		t.Fatal("expected nil cmd when manager is disabled")
	}
	if m.Focused() != a {
		t.Fatal("expected focus to remain unchanged when disabled")
	}
}

func TestNextPrevEmitFocusChanged(t *testing.T) {
	a := &focusSpy{}
	b := &focusSpy{}
	m := New(Ring(a, b))

	msg := runCmd(m.Next())
	changed, ok := msg.(FocusChangedMsg)
	if !ok {
		t.Fatalf("expected FocusChangedMsg, got %T", msg)
	}
	if changed.From != 0 || changed.To != 1 {
		t.Fatalf("expected from 0 to 1, got %+v", changed)
	}

	msg = runCmd(m.Prev())
	changed, ok = msg.(FocusChangedMsg)
	if !ok {
		t.Fatalf("expected FocusChangedMsg, got %T", msg)
	}
	if changed.From != 1 || changed.To != 0 {
		t.Fatalf("expected from 1 to 0, got %+v", changed)
	}
}

func runCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}
