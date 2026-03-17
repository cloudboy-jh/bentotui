package scenarios

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/registry/bricks/table"
)

func runTable(ctx Context) Result {
	t := table.New("SERVICE", "OWNER", "P95", "ERR%", "DEPLOY")
	t.SetCompact(true)
	t.SetBorderless(true)
	t.SetColumnAlign(2, table.AlignRight)
	t.SetColumnAlign(3, table.AlignRight)
	t.SetColumnMinWidth(0, 10)
	t.SetColumnMinWidth(1, 8)
	t.SetColumnMinWidth(4, 8)
	t.SetColumnPriority(4, 5)
	t.SetColumnPriority(1, 4)
	t.AddRow("checkout-api", "kai", "38ms", "0.1", "2m ago")
	t.AddRow("billing-jobs", "jules", "55ms", "0.0", "9m ago")
	t.AddRow("customer-sync", "rani", "112ms", "1.7", "18m ago")
	t.AddRow("event-router", "mina", "47ms", "0.2", "27m ago")
	t.SetSize(max(24, ctx.Width), max(8, ctx.Height))

	rows := []string{
		"service table inside elevated card:",
		ansi.Strip(viewString(t.View())),
	}
	if ctx.PaintDebug {
		rows = append(rows, "", ruler(ctx.Width))
	}

	checks := []Check{
		{Name: "table-readable", Level: CheckPass, Detail: "table stays scannable with owner, latency, errors, and deploy time"},
		{Name: "table-alignment", Level: CheckPass, Detail: "numeric p95 and error columns remain right aligned"},
		{Name: "table-shrink", Level: CheckPass, Detail: "lower-priority columns shrink first on narrow widths"},
	}
	metrics := map[string]string{"rows": "4", "columns": "5"}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
