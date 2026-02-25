package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/dialog"
)

func keyPress(text string) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Text: text})
}

func specialKey(code rune) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func TestHarnessSlashTypesIntoInput(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName))
	p.SetSize(120, 40)

	_, cmd := p.Update(keyPress("/"))
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(theme.OpenThemePickerMsg); ok {
			t.Fatal("did not expect theme picker command from raw '/' key")
		}
	}
	if got := p.input.Value(); got != "/" {
		t.Fatalf("expected input to contain '/', got %q", got)
	}
}

func TestHarnessThemeCommandOpensThemePicker(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName))
	p.SetSize(120, 40)
	p.input.SetValue("/theme")

	_, cmd := p.Update(specialKey(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("expected theme picker command from /theme")
	}
	msg := cmd()
	if _, ok := msg.(theme.OpenThemePickerMsg); !ok {
		t.Fatalf("expected theme.OpenThemePickerMsg, got %T", msg)
	}
}

func TestHarnessHotkeysOpenDialogs(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName))
	p.SetSize(120, 40)
	_, _ = p.Update(specialKey(tea.KeyTab))
	if got := p.focusName(); got != "actions" {
		t.Fatalf("expected focus to move to actions, got %q", got)
	}

	_, cmd := p.Update(keyPress("d"))
	if cmd == nil {
		t.Fatal("expected custom dialog command for 'd' on actions focus")
	}
	open := cmd()
	if _, ok := open.(dialog.OpenMsg); !ok {
		t.Fatalf("expected dialog.OpenMsg for 'd', got %T", open)
	}

	_, cmd = p.Update(keyPress("x"))
	if cmd == nil {
		t.Fatal("expected confirm dialog command for 'x' on actions focus")
	}
	open = cmd()
	if _, ok := open.(dialog.OpenMsg); !ok {
		t.Fatalf("expected dialog.OpenMsg for 'x', got %T", open)
	}
}

func TestHarnessInputAcceptsDAndQCharacters(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName))
	p.SetSize(120, 40)

	_, cmd := p.Update(keyPress("d"))
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(dialog.OpenMsg); ok {
			t.Fatalf("did not expect dialog open from 'd' while input focused")
		}
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Fatalf("did not expect quit from 'd' while input focused")
		}
	}
	_, cmd = p.Update(keyPress("q"))
	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(dialog.OpenMsg); ok {
			t.Fatalf("did not expect dialog open from 'q' while input focused")
		}
		if _, ok := msg.(tea.QuitMsg); ok {
			t.Fatalf("did not expect quit from 'q' while input focused")
		}
	}
	if got := p.input.Value(); got != "dq" {
		t.Fatalf("expected input value dq, got %q", got)
	}
}

func TestHarnessEnterSubmitsInput(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName))
	p.SetSize(120, 40)
	p.input.SetValue("hello bento")

	_, cmd := p.Update(specialKey(tea.KeyEnter))
	if cmd != nil {
		t.Fatalf("expected nil command for plain text submit, got non-nil")
	}
	if len(p.events) == 0 {
		t.Fatal("expected event log entry after submit")
	}
	if !strings.Contains(p.events[0], "submitted: hello bento") {
		t.Fatalf("unexpected submit log entry: %q", p.events[0])
	}
	if p.input.Value() != "" {
		t.Fatalf("expected input to clear after submit, got %q", p.input.Value())
	}
}

func TestHarnessActionCycleAndRun(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName))
	p.SetSize(120, 40)

	_, _ = p.Update(specialKey(tea.KeyTab))
	if got := p.focusName(); got != "actions" {
		t.Fatalf("expected focus to move to actions, got %q", got)
	}

	_, _ = p.Update(specialKey(tea.KeyRight))
	if p.actionIdx != 1 {
		t.Fatalf("expected actionIdx 1 after right key, got %d", p.actionIdx)
	}

	_, cmd := p.Update(specialKey(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("expected dialog command from selected action")
	}
	open := cmd()
	if _, ok := open.(dialog.OpenMsg); !ok {
		t.Fatalf("expected dialog.OpenMsg from action run, got %T", open)
	}
}
