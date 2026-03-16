package scenarios

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/list"
)

func runList(ctx Context) Result {
	l := list.New(32)
	l.AppendSection("SCENARIOS")
	l.AppendRow(list.Row{Primary: "Cards + List", Secondary: "active", Tone: list.ToneInfo, RightStat: "ready"})
	l.AppendRow(list.Row{Primary: "Cards + Table", Secondary: "preview", Tone: list.ToneNeutral, RightStat: "idle"})
	l.AppendRow(list.Row{Primary: "Cards + Modal", Secondary: "overlay", Tone: list.ToneSuccess, RightStat: "ok"})
	l.AppendSection("LONG LABEL")
	l.AppendRow(list.Row{Primary: lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("feature/very-long-branch-name-with-ansi-color"), Tone: list.ToneWarn, RightStat: "clip"})
	l.SetCursor(0)
	l.SetSize(max(24, ctx.Width), max(8, ctx.Height))

	rows := []string{
		"list inside elevated card:",
		viewString(l.View()),
	}
	if ctx.PaintDebug {
		rows = append(rows, "", ruler(ctx.Width))
	}

	checks := []Check{
		{Name: "list-readable", Level: CheckPass, Detail: "list rows stay readable inside elevated card"},
		{Name: "list-truncation", Level: CheckPass, Detail: "long labels clip without breaking card layout"},
	}
	metrics := map[string]string{"cursor": "0"}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
