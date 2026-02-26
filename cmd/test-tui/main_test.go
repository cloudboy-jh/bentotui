package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/dialog"
)

func keyPress(text string) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Text: text})
}

func specialKey(code rune) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func ctrlKey(code rune) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Code: code, Mod: tea.ModCtrl})
}

func TestHarnessSlashTypesIntoInput(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName), "harness", "secondary")
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
	p := newHarnessPage(theme.Preset(theme.DefaultName), "harness", "secondary")
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

func TestHarnessDialogCommandsOpenDialogs(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName), "harness", "secondary")
	p.SetSize(120, 40)
	p.input.SetValue("/dialog")

	_, cmd := p.Update(specialKey(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("expected custom dialog command for /dialog")
	}
	open := cmd()
	if _, ok := open.(dialog.OpenMsg); !ok {
		t.Fatalf("expected dialog.OpenMsg for /dialog, got %T", open)
	}

}

func TestHarnessPageCommandNavigatesToSecondary(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName), "harness", "secondary")
	p.SetSize(120, 40)
	p.input.SetValue("/page")

	_, cmd := p.Update(specialKey(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("expected navigate command for /page")
	}
	msg := cmd()
	nav, ok := msg.(core.NavigateMsg)
	if !ok {
		t.Fatalf("expected core.NavigateMsg, got %T", msg)
	}
	if nav.Page != "secondary" {
		t.Fatalf("expected page secondary, got %q", nav.Page)
	}
}

func TestHarnessInputAcceptsDAndQCharacters(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName), "harness", "secondary")
	p.SetSize(120, 40)

	_, cmd := p.Update(keyPress("d"))
	if cmd != nil {
		if _, ok := cmd().(tea.QuitMsg); ok {
			t.Fatalf("did not expect quit from 'd' while input focused")
		}
	}
	_, cmd = p.Update(keyPress("q"))
	if cmd != nil {
		if _, ok := cmd().(tea.QuitMsg); ok {
			t.Fatalf("did not expect quit from 'q' while input focused")
		}
	}
	if got := p.input.Value(); got != "dq" {
		t.Fatalf("expected input value dq, got %q", got)
	}
}

func TestHarnessEnterSubmitsInput(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName), "harness", "secondary")
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

func TestHarnessViewRowsMatchViewportWidth(t *testing.T) {
	p := newHarnessPage(theme.Preset(theme.DefaultName), "harness", "secondary")
	p.SetSize(120, 40)

	view := core.ViewString(p.View())
	lines := strings.Split(view, "\n")
	if len(lines) == 0 {
		t.Fatal("expected non-empty harness view")
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 120 {
			t.Fatalf("expected row %d width 120, got %d", i, w)
		}
	}
}
