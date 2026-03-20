package state

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/dialog"
	"github.com/cloudboy-jh/bentotui/theme"
)

func TestModelViewExactDimensions(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	out := viewFromTea(m.View())
	lines := strings.Split(out, "\n")
	if len(lines) != 30 {
		t.Fatalf("expected 30 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got != 100 {
			t.Fatalf("line %d width mismatch: got %d want 100", i, got)
		}
	}
}

func TestModelViewContainsWorkspaceFooterAndCards(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	out := viewFromTea(m.View())
	checks := []string{"Services", "Queue", "Progress", "overview | queue:1", "theme:", "ctrl+k"}
	for _, token := range checks {
		if !strings.Contains(out, token) {
			t.Fatalf("expected token %q in render output", token)
		}
	}
}

func TestModelSmallViewportKeepsFooterBar(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 70, Height: 12})
	out := viewFromTea(m.View())
	lines := strings.Split(out, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[len(lines)-1]) == "" {
		t.Fatalf("expected footer bar right status in small viewport")
	}
}

func TestDialogLifecycleOpenClose(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	if m.dialogs.IsOpen() {
		t.Fatalf("expected dialogs closed initially")
	}

	_, _ = m.Update(dialog.OpenMsg{Dialog: dialog.Confirm{DialogTitle: "Test", Message: "hello"}})
	if !m.dialogs.IsOpen() {
		t.Fatalf("expected dialog manager to open on OpenMsg")
	}

	_, _ = m.Update(dialog.CloseMsg{})
	if m.dialogs.IsOpen() {
		t.Fatalf("expected dialog manager to close on CloseMsg")
	}
}

func TestPaletteThemePickerFlowOpensDialog(t *testing.T) {
	m := NewModel()
	_, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	openCmd := m.openPalette()
	openMsg := openCmd()
	_, _ = m.Update(openMsg)
	if !m.dialogs.IsOpen() {
		t.Fatalf("expected command palette dialog to be open")
	}

	_, _ = m.Update(dialog.CloseMsg{})
	_, cmd := m.Update(openThemePickerMsg{})
	if cmd == nil {
		t.Fatalf("expected theme picker action to emit open command")
	}

	msg := cmd()
	_, _ = m.Update(msg)
	if !m.dialogs.IsOpen() {
		t.Fatalf("expected theme picker dialog to open via palette flow")
	}

	out := viewFromTea(m.dialogs.View())
	if !strings.Contains(out, "Themes") {
		t.Fatalf("expected active dialog to be theme picker")
	}
}

func TestThemeChangedUpdatesModelTheme(t *testing.T) {
	m := NewModel()
	original := m.theme.Name()
	defer func() {
		_, _ = theme.SetTheme(original)
	}()
	next := ""
	for _, name := range theme.AvailableThemes() {
		if name != original {
			next = name
			break
		}
	}
	if next == "" {
		t.Fatalf("expected at least two themes to test theme change")
	}

	themeValue, err := theme.SetTheme(next)
	if err != nil {
		t.Fatalf("set theme: %v", err)
	}
	_, _ = m.Update(theme.ThemeChangedMsg{Name: next, Theme: themeValue})

	if m.theme == nil || m.theme.Name() != next {
		t.Fatalf("expected model theme to update to %q", next)
	}
}

func viewFromTea(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(v.Content)
}
