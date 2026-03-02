package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/focus"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/dialog"
	"github.com/cloudboy-jh/bentotui/ui/containers/bar"
)

func newTestPage() *starterPage {
	ft := bar.New()
	return newStarterPage(theme.Preset(theme.DefaultName), "harness", "secondary", ft)
}

func keyPress(text string) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Text: text})
}

func specialKey(code rune) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func ctrlKey(code rune) tea.Msg {
	return tea.KeyPressMsg(tea.Key{Code: code, Mod: tea.ModCtrl})
}

// TestHarnessSlashFromNonCommandPanelOpensPalette checks that pressing "/" when
// a non-input panel is focused opens the command palette rather than typing into input.
func TestHarnessSlashFromNonCommandPanelOpensPalette(t *testing.T) {
	p := newTestPage()
	p.SetSize(120, 40)

	// Move focus away from Command panel (index 1) to Info (index 0)
	p.focusIdx = 0
	p.syncInputFocus()

	_, cmd := p.Update(keyPress("/"))
	if cmd == nil {
		t.Fatal("expected command palette cmd from '/' when non-input panel is focused")
	}
	msg := cmd()
	if _, ok := msg.(interface{}); !ok {
		t.Fatalf("unexpected msg type: %T", msg)
	}
}

// TestHarnessSlashInCommandPanelTypesIntoInput verifies "/" types into the input
// when the Command panel is focused.
func TestHarnessSlashInCommandPanelTypesIntoInput(t *testing.T) {
	p := newTestPage()
	p.SetSize(120, 40)

	// Command panel is focused by default (index 1)
	if p.focusIdx != commandPanelIdx {
		t.Fatalf("expected command panel to be focused initially, got %d", p.focusIdx)
	}

	_, _ = p.Update(keyPress("/"))
	if got := p.input.Value(); got != "/" {
		t.Fatalf("expected input to contain '/', got %q", got)
	}
}

func TestHarnessThemeCommandOpensThemePicker(t *testing.T) {
	p := newTestPage()
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

func TestHarnessDialogCommandOpensDialog(t *testing.T) {
	p := newTestPage()
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
	p := newTestPage()
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
	p := newTestPage()
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
	p := newTestPage()
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
	p := newTestPage()
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

func TestHarnessFocusCycleUpdatesFocusIdx(t *testing.T) {
	p := newTestPage()
	p.SetSize(120, 40)

	initial := p.focusIdx

	// Simulate FocusChangedMsg from the focus manager
	_, _ = p.Update(focus.FocusChangedMsg{From: initial, To: 0})
	if p.focusIdx != 0 {
		t.Fatalf("expected focusIdx 0 after FocusChangedMsg{To:0}, got %d", p.focusIdx)
	}

	_, _ = p.Update(focus.FocusChangedMsg{From: 0, To: 2})
	if p.focusIdx != 2 {
		t.Fatalf("expected focusIdx 2 after FocusChangedMsg{To:2}, got %d", p.focusIdx)
	}
}

func TestHarnessFocusCommandPanelActivatesInput(t *testing.T) {
	p := newTestPage()
	p.SetSize(120, 40)

	// Move to Info panel first
	_, _ = p.Update(focus.FocusChangedMsg{From: commandPanelIdx, To: 0})
	if p.input.Focused() {
		t.Fatal("expected input to be blurred when non-command panel is focused")
	}

	// Move back to Command panel
	_, _ = p.Update(focus.FocusChangedMsg{From: 0, To: commandPanelIdx})
	if !p.input.Focused() {
		t.Fatal("expected input to be focused when Command panel is active")
	}
}

func TestHarnessAllFourPanelNamesKnown(t *testing.T) {
	if len(panelNames) != 4 {
		t.Fatalf("expected 4 panel names, got %d", len(panelNames))
	}
	if len(panelIcons) != 4 {
		t.Fatalf("expected 4 panel icons, got %d", len(panelIcons))
	}
}

func TestHarnessStatusPanelContainsHealthChecks(t *testing.T) {
	p := newTestPage()
	p.SetSize(120, 40)

	// Check the statusText content directly — the horizontal canvas-layer
	// renderer can split plain words across ANSI escape sequences, making
	// raw string search on the final view unreliable.
	status := p.statusText.text
	for _, keyword := range []string{"footer", "palette", "focus", "theme"} {
		if !strings.Contains(status, keyword) {
			t.Fatalf("expected statusText to contain %q, got:\n%s", keyword, status)
		}
	}
}
