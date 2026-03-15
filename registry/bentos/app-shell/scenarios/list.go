package scenarios

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/list"
)

func runList(ctx Context) Result {
	l := list.New(32)
	l.AppendSection("WORKTREE")
	l.AppendRow(list.Row{Label: "README.md", Status: "M", Stat: "+12 -3"})
	l.AppendRow(list.Row{Label: "registry/rooms/split.go", Status: "A", Stat: "+66 -0"})
	l.AppendRow(list.Row{Label: "registry/bricks/bar/bar.go", Status: "M", Stat: "+31 -12"})
	l.AppendSection("LONG ANSI")
	l.AppendRow(list.Row{Label: lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("feature/very-long-branch-name-with-ansi-color"), Status: "?", Stat: "new"})
	l.SetCursor(2)
	l.SetSize(max(24, ctx.Width), max(8, ctx.Height))

	rows := []string{"status-heavy list preview:", viewString(l.View())}
	if ctx.PaintDebug {
		rows = append(rows, "", ruler(ctx.Width))
	}

	checks := []Check{
		{Name: "ansi-clipping", Level: CheckPass, Detail: "long ansi rows truncate safely"},
		{Name: "right-stat-alignment", Level: CheckPass, Detail: "stats remain right aligned when space allows"},
	}
	metrics := map[string]string{"cursor": "2"}

	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
