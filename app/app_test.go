package app

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/dialog"
	"github.com/cloudboy-jh/bentotui/router"
	"github.com/cloudboy-jh/bentotui/theme"
)

type testPage struct {
	title  string
	width  int
	height int
}

func (p *testPage) Init() tea.Cmd { return nil }

func (p *testPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "2":
			return p, func() tea.Msg { return router.Navigate("two") }
		case "q":
			return p, tea.Quit
		}
	}
	return p, nil
}

func (p *testPage) View() tea.View { return tea.NewView("") }

func (p *testPage) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *testPage) GetSize() (int, int) { return p.width, p.height }

func (p *testPage) Title() string { return p.title }

func makePage(title string) func() core.Page {
	return func() core.Page { return &testPage{title: title} }
}

func keyPress(text string) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Text: text})
}

func TestKeyNavigatesRoutes(t *testing.T) {
	m := New(
		WithPages(
			router.Page("one", makePage("one")),
			router.Page("two", makePage("two")),
		),
	)

	_, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if got := m.Router().CurrentName(); got != "one" {
		t.Fatalf("expected initial route one, got %q", got)
	}

	_, cmd := m.Update(keyPress("2"))
	if cmd == nil {
		t.Fatal("expected navigation command from key '2'")
	}
	navMsg := cmd()
	_, _ = m.Update(navMsg)

	if got := m.Router().CurrentName(); got != "two" {
		t.Fatalf("expected route two after key '2', got %q", got)
	}
}

func TestKeyQuitCmdPropagates(t *testing.T) {
	m := New(
		WithPages(
			router.Page("one", makePage("one")),
		),
	)

	_, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	_, cmd := m.Update(keyPress("q"))
	if cmd == nil {
		t.Fatal("expected quit command from key 'q'")
	}
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", msg)
	}
}

func TestOpenThemePickerMessageOpensThemeDialog(t *testing.T) {
	m := New(
		WithPages(
			router.Page("one", makePage("one")),
		),
	)

	_, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	_, cmd := m.Update(theme.OpenThemePicker())
	if cmd == nil {
		t.Fatal("expected dialog open command from OpenThemePicker message")
	}
	open := cmd()
	if _, ok := open.(dialog.OpenMsg); !ok {
		t.Fatalf("expected dialog.OpenMsg, got %T", open)
	}

	_, _ = m.Update(open)
	if !m.Dialogs().IsOpen() {
		t.Fatal("expected theme dialog to be open")
	}
}
