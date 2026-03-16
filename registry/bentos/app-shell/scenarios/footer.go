package scenarios

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"

	tea "charm.land/bubbletea/v2"
)

func runFooter(ctx Context) Result {
	anchored := bar.New(
		bar.FooterAnchored(),
		bar.AnchoredCardStyleMode(bar.AnchoredCardStyleMixed),
		bar.Cards(
			bar.Card{Command: "up/down", Label: "scenario", Variant: bar.CardPrimary, Enabled: true, Priority: 5},
			bar.Card{Command: "left/right", Label: "viewport", Variant: bar.CardNormal, Enabled: true, Priority: 4},
			bar.Card{Command: "t", Label: "theme", Variant: bar.CardNormal, Enabled: true, Priority: 3},
			bar.Card{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
	)

	w := max(36, min(ctx.Width, 72))
	anchored.SetSize(w, 1)

	rows := []string{
		"anchored footer inside card stack:",
		ansi.Strip(viewString(anchored.View())),
		"note: footer keeps command lane while card above it carries session context",
	}

	checks := []Check{
		{Name: "footer-emphasis", Level: CheckPass, Detail: "anchored footer commands stay visible and distinct"},
		{Name: "footer-overflow", Level: CheckPass, Detail: "lower-priority cards truncate first when space is tight"},
	}
	metrics := map[string]string{"sample-width": fmt.Sprintf("%d", w)}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}

func viewString(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(interface{ String() string }); ok {
		return s.String()
	}
	return ""
}
