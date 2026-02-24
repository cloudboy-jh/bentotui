package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/dialog"
	"github.com/cloudboy-jh/bentotui/theme"
)

func keyPress(text string) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Text: text})
}

func TestHomePageKeyBindings(t *testing.T) {
	p := newHomePage(theme.Preset("amber"))
	p.SetSize(120, 40)

	_, cmd := p.Update(keyPress("2"))
	if cmd == nil {
		t.Fatal("expected navigation cmd for key '2'")
	}
	nav := cmd()
	if msg, ok := nav.(core.NavigateMsg); !ok || msg.Page != "inspect" {
		t.Fatalf("expected navigate to inspect, got %T %#v", nav, nav)
	}

	_, cmd = p.Update(keyPress("o"))
	if cmd == nil {
		t.Fatal("expected dialog open cmd for key 'o'")
	}
	open := cmd()
	if _, ok := open.(dialog.OpenMsg); !ok {
		t.Fatalf("expected dialog.OpenMsg, got %T", open)
	}

	_, cmd = p.Update(keyPress("q"))
	if cmd == nil {
		t.Fatal("expected quit cmd for key 'q'")
	}
	quit := cmd()
	if _, ok := quit.(tea.QuitMsg); !ok {
		t.Fatalf("expected tea.QuitMsg, got %T", quit)
	}
}

func TestInspectPageKeyBindings(t *testing.T) {
	p := newInspectPage(theme.Preset("amber"))
	p.SetSize(120, 40)

	_, cmd := p.Update(keyPress("1"))
	if cmd == nil {
		t.Fatal("expected navigation cmd for key '1'")
	}
	nav := cmd()
	if msg, ok := nav.(core.NavigateMsg); !ok || msg.Page != "home" {
		t.Fatalf("expected navigate to home, got %T %#v", nav, nav)
	}

	_, cmd = p.Update(keyPress("o"))
	if cmd == nil {
		t.Fatal("expected dialog open cmd for key 'o'")
	}
	open := cmd()
	if _, ok := open.(dialog.OpenMsg); !ok {
		t.Fatalf("expected dialog.OpenMsg, got %T", open)
	}
}
