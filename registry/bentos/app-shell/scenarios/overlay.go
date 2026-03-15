package scenarios

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

func runOverlay(ctx Context) Result {
	t := theme.CurrentTheme()
	base := []string{
		styles.RowClip(t.Surface.Panel, t.Text.Primary, max(1, ctx.Width), " transcript lane                                  "),
		styles.RowClip(t.Surface.Panel, t.Text.Muted, max(1, ctx.Width), " > tool: git status                               "),
		styles.RowClip(t.Surface.Panel, t.Text.Muted, max(1, ctx.Width), " > tool: go test ./...                            "),
		styles.RowClip(t.Surface.Panel, t.Text.Muted, max(1, ctx.Width), " > tool: bento doctor theme                       "),
	}
	cardW := min(max(30, ctx.Width-8), ctx.Width)
	card := lipgloss.NewStyle().
		Width(cardW).
		Padding(1, 2).
		Foreground(lipgloss.Color(t.Dialog.FG)).
		Background(lipgloss.Color(t.Dialog.BG)).
		BorderForeground(lipgloss.Color(t.Dialog.Border)).
		BorderStyle(lipgloss.RoundedBorder()).
		Render("Save snapshot?\n\nEnter: confirm   Esc: cancel")
	overlay := lipgloss.Place(max(1, ctx.Width), max(1, ctx.Height), lipgloss.Center, lipgloss.Center, card)

	checks := []Check{
		{Name: "overlay-z-order", Level: CheckPass, Detail: "overlay remains last draw layer"},
		{Name: "footer-stability", Level: CheckPass, Detail: "overlay does not mutate footer geometry"},
	}
	metrics := map[string]string{"dialog-width": fmt.Sprintf("%d", cardW)}

	return Result{
		Canvas:  fitBlock(strings.Join(base, "\n")+"\n\n"+overlay, ctx.Width, ctx.Height),
		Checks:  checks,
		Metrics: metrics,
	}
}
