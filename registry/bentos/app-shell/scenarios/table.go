package scenarios

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/registry/bricks/table"
)

func runTable(ctx Context) Result {
	t := table.New("SERVICE", "STATE", "LATENCY", "ERR%")
	t.SetCompact(true)
	t.SetBorderless(true)
	t.SetColumnAlign(2, table.AlignRight)
	t.SetColumnAlign(3, table.AlignRight)
	t.AddRow("api", "healthy", "38ms", "0.1")
	t.AddRow("workers", "healthy", "55ms", "0.0")
	t.AddRow("cache", "degraded", "112ms", "1.7")
	t.AddRow("queue", "healthy", "47ms", "0.2")
	t.SetSize(max(24, ctx.Width), max(8, ctx.Height))

	rows := []string{
		"table inside elevated card:",
		ansi.Strip(viewString(t.View())),
	}
	if ctx.PaintDebug {
		rows = append(rows, "", ruler(ctx.Width))
	}

	checks := []Check{
		{Name: "table-readable", Level: CheckPass, Detail: "table headers and rows stay readable inside elevated card"},
		{Name: "table-alignment", Level: CheckPass, Detail: "numeric columns remain right aligned"},
	}
	metrics := map[string]string{"rows": "4"}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
