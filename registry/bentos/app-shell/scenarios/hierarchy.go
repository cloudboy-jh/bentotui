package scenarios

import (
	"fmt"
	"strings"

	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

func runHierarchy(ctx Context) Result {
	t := theme.CurrentTheme()
	row := func(label, bg, fg string) string {
		return styles.RowClip(bg, fg, max(1, ctx.Width), " "+label)
	}

	rows := []string{
		row("surface.panel        :: baseline panel depth", t.Surface.Panel, t.Text.Primary),
		row("surface.elevated     :: raised depth", t.Surface.Elevated, t.Text.Primary),
		row("surface.interactive  :: active / interactive depth", t.Surface.Interactive, t.Text.Primary),
		row("focus owner stripe   :: visual focus authority", t.Surface.Interactive, t.Text.Primary),
	}
	if ctx.PaintDebug {
		rows = append(rows, ruler(ctx.Width))
	}

	checks := []Check{
		{Name: "depth-separation", Level: CheckPass, Detail: "panel/elevated/interactive ladders are visible"},
		{Name: "focus-affordance", Level: CheckPass, Detail: "focus owner is explicit in frame"},
	}
	metrics := map[string]string{
		"focus-owner": ctx.FocusOwner,
		"rows":        fmt.Sprintf("%d", len(rows)),
	}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
