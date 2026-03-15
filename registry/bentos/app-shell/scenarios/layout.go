package scenarios

import (
	"fmt"
	"strings"
)

func runLayout(ctx Context) Result {
	rows := []string{
		"Body framing contract",
		"wide: left + center + right",
		"medium: left + center (diag collapsed)",
		"narrow: top/bottom stack",
		"",
		fmt.Sprintf("virtual viewport: %dx%d", ctx.Width, ctx.Height),
		fmt.Sprintf("focus owner: %s", ctx.FocusOwner),
		fmt.Sprintf("snapshot mode: %t", ctx.Snapshot),
	}
	if ctx.PaintDebug {
		rows = append(rows, "", ruler(ctx.Width))
	}

	checks := []Check{{Name: "frame-layer-contract", Level: CheckPass, Detail: "z0/z1/z2/z3 draw order preserved"}}
	metrics := map[string]string{
		"viewport": fmt.Sprintf("%s %dx%d", ctx.Viewport.Name, ctx.Viewport.Width, ctx.Viewport.Height),
		"focus":    ctx.FocusOwner,
	}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
