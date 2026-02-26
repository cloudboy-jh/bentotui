package dialog

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
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
