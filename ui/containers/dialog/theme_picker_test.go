package dialog

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
)

func TestThemePickerImplementsSizeable(t *testing.T) {
	var _ core.Sizeable = NewThemePicker()
}

func TestThemePickerRowsMatchAssignedWidth(t *testing.T) {
	p := NewThemePicker()
	p.SetSize(36, 12)

	view := core.ViewString(p.View())
	lines := strings.Split(view, "\n")
	if len(lines) == 0 {
		t.Fatal("expected non-empty theme picker view")
	}
	for i, line := range lines {
		if w := lipgloss.Width(line); w != 36 {
			t.Fatalf("expected row %d width 36, got %d", i, w)
		}
	}
}

func TestThemePickerSelectionPreviewAndEscRevert(t *testing.T) {
	_, _ = theme.SetTheme(theme.CatppuccinMochaName)
	p := NewThemePicker()

	updated, cmd := p.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	next := updated.(*ThemePicker)
	msg := runThemeCmd(cmd)
	if _, ok := msg.(theme.ThemeChangedMsg); !ok {
		t.Fatalf("expected ThemeChangedMsg on preview, got %T", msg)
	}
	if theme.CurrentThemeName() == theme.CatppuccinMochaName {
		t.Fatal("expected selection move to preview a different theme")
	}

	_, cmd = next.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	_ = runThemeCmd(cmd)
	if theme.CurrentThemeName() != theme.CatppuccinMochaName {
		t.Fatal("expected esc to revert to base theme")
	}
}

func TestThemePickerEnterCommitsTheme(t *testing.T) {
	_, _ = theme.SetTheme(theme.CatppuccinMochaName)
	p := NewThemePicker()
	updated, _ := p.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	next := updated.(*ThemePicker)
	_, cmd := next.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	_ = runThemeCmd(cmd)
	if theme.CurrentThemeName() == theme.CatppuccinMochaName {
		t.Fatal("expected enter to commit selected theme")
	}
}

func runThemeCmd(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}
