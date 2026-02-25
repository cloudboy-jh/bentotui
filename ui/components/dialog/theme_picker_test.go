package dialog

import (
	"testing"

	"github.com/cloudboy-jh/bentotui/core"
)

func TestThemePickerImplementsSizeable(t *testing.T) {
	var _ core.Sizeable = NewThemePicker()
}
