package scenarios

import (
	"fmt"
	"strings"
)

func runStress(ctx Context) Result {
	rows := []string{
		fmt.Sprintf("viewport preset: %s (%dx%d)", ctx.Viewport.Name, ctx.Viewport.Width, ctx.Viewport.Height),
		fmt.Sprintf("canvas target: %dx%d", ctx.Width, ctx.Height),
		fmt.Sprintf("stress step: %d", ctx.StressStep),
		fmt.Sprintf("focus owner: %s", ctx.FocusOwner),
		"",
		ruler(ctx.Width),
		rulerIndex(ctx.Width),
		ruler(ctx.Width),
		"",
		"Expected: no clipping tears, no stale color strips, stable footer lane.",
	}

	level := CheckPass
	if ctx.StressStep%7 == 0 && ctx.StressStep > 0 {
		level = CheckWarn
	}
	checks := []Check{{Name: "stress-seam-watch", Level: level, Detail: "manual stress watchpoint"}}
	metrics := map[string]string{"stress-step": fmt.Sprintf("%d", ctx.StressStep)}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
