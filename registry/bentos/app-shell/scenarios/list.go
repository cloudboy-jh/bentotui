package scenarios

import (
	"strings"

	"github.com/cloudboy-jh/bentotui/registry/bricks/list"
)

func runList(ctx Context) Result {
	l := list.New(32)
	l.SetDensity(list.DensityCompact)
	l.AppendSection("WORKSPACES")
	l.AppendRow(list.Row{Primary: "Checkout API", Secondary: "pager active", Tone: list.ToneDanger, RightStat: "2m"})
	l.AppendRow(list.Row{Primary: "Billing Events", Secondary: "healthy", Tone: list.ToneSuccess, RightStat: "9m"})
	l.AppendRow(list.Row{Primary: "Growth Experiments", Secondary: "in review", Tone: list.ToneInfo, RightStat: "18m"})
	l.AppendSection("BACKLOG")
	l.AppendRow(list.Row{Primary: "Customer imports from legacy warehouse with long labels", Secondary: "triage", Tone: list.ToneWarn, RightStat: "27m"})
	l.SetCursor(1)
	l.SetSize(max(24, ctx.Width), max(8, ctx.Height))

	rows := []string{
		"workspace list inside elevated card:",
		viewString(l.View()),
	}
	if ctx.PaintDebug {
		rows = append(rows, "", ruler(ctx.Width))
	}

	checks := []Check{
		{Name: "list-scannable", Level: CheckPass, Detail: "status lane and right-time lane stay easy to scan"},
		{Name: "list-selected-row", Level: CheckPass, Detail: "selected row remains obvious in dense mode"},
		{Name: "list-overflow", Level: CheckPass, Detail: "long labels truncate while preserving right stats"},
	}
	metrics := map[string]string{"cursor": "1", "density": "compact"}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
