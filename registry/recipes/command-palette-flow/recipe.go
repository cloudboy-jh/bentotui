package commandpaletteflow

import (
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/dialog"
)

// Open returns a command that opens the command palette dialog.
// Copy and own this recipe, then wire your app-specific command list.
func Open(commands []dialog.Command) tea.Cmd {
	return func() tea.Msg {
		palette := dialog.NewCommandPalette(commands)
		return dialog.Open(dialog.Custom{
			DialogTitle: "Commands",
			Content:     palette,
			Width:       56,
			Height:      18,
		})
	}
}
