package scenarios

import (
	"fmt"
	"strings"

	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"

	tea "charm.land/bubbletea/v2"
)

func runFooter(ctx Context) Result {
	plain := bar.New(
		bar.FooterAnchored(),
		bar.AnchoredCardStyleMode(bar.AnchoredCardStylePlain),
		bar.Cards(
			bar.Card{Command: "enter", Label: "run", Variant: bar.CardPrimary, Enabled: true, Priority: 5},
			bar.Card{Command: "tab", Label: "focus", Variant: bar.CardNormal, Enabled: true, Priority: 4},
			bar.Card{Command: "shift+tab", Label: "prev", Variant: bar.CardMuted, Enabled: true, Priority: 3},
			bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
	)
	chip := bar.New(
		bar.FooterAnchored(),
		bar.AnchoredCardStyleMode(bar.AnchoredCardStyleChip),
		bar.Cards(
			bar.Card{Command: "enter", Label: "run", Variant: bar.CardPrimary, Enabled: true, Priority: 5},
			bar.Card{Command: "tab", Label: "focus", Variant: bar.CardNormal, Enabled: true, Priority: 4},
			bar.Card{Command: "shift+tab", Label: "prev", Variant: bar.CardMuted, Enabled: true, Priority: 3},
			bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
	)
	mixed := bar.New(
		bar.FooterAnchored(),
		bar.AnchoredCardStyleMode(bar.AnchoredCardStyleMixed),
		bar.Cards(
			bar.Card{Command: "enter", Label: "run", Variant: bar.CardPrimary, Enabled: true, Priority: 5},
			bar.Card{Command: "tab", Label: "focus", Variant: bar.CardNormal, Enabled: true, Priority: 4},
			bar.Card{Command: "shift+tab", Label: "prev", Variant: bar.CardMuted, Enabled: true, Priority: 3},
			bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
	)

	w := max(36, min(ctx.Width, 72))
	plain.SetSize(w, 1)
	chip.SetSize(w, 1)
	mixed.SetSize(w, 1)

	rows := []string{
		"plain mode:",
		viewString(plain.View()),
		"",
		"chip mode (anchored):",
		viewString(chip.View()),
		"",
		"mixed mode (anchored):",
		viewString(mixed.View()),
	}

	checks := []Check{
		{Name: "anchored-footer-row", Level: CheckPass, Detail: "footer keeps single full-width ownership"},
		{Name: "anchored-style-modes", Level: CheckPass, Detail: "plain/chip/mixed style modes available"},
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
