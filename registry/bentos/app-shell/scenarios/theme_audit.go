package scenarios

import (
	"fmt"
	"strings"

	"github.com/cloudboy-jh/bentotui/theme"
)

func runThemeAudit(ctx Context) Result {
	all := theme.AvailableThemes()
	stableSet := map[string]struct{}{}
	for _, n := range theme.AvailableStableThemes() {
		stableSet[n] = struct{}{}
	}

	limit := min(len(all), 10)
	rows := []string{"theme registry audit (stable-first):", ""}
	for i := 0; i < limit; i++ {
		name := all[i]
		meta, _ := theme.ThemeMetadata(name)
		tier := string(meta.Tier)
		if tier == "" {
			tier = "experimental"
		}
		rows = append(rows, fmt.Sprintf("%2d. %-20s tier=%-12s score=%5.1f", i+1, name, tier, meta.Score))
	}

	stableCount := len(stableSet)
	checks := []Check{{Name: "theme-stable-set", Level: CheckPass, Detail: fmt.Sprintf("stable themes available: %d", stableCount)}}
	if stableCount < 4 {
		checks[0].Level = CheckWarn
		checks[0].Detail = fmt.Sprintf("stable theme count is low (%d)", stableCount)
	}
	if len(all) == 0 {
		checks = append(checks, Check{Name: "theme-registry", Level: CheckFail, Detail: "no themes registered"})
	} else {
		checks = append(checks, Check{Name: "theme-registry", Level: CheckPass, Detail: fmt.Sprintf("themes registered: %d", len(all))})
	}

	metrics := map[string]string{
		"themes": fmt.Sprintf("%d", len(all)),
		"stable": fmt.Sprintf("%d", stableCount),
	}
	return Result{Canvas: fitBlock(strings.Join(rows, "\n"), ctx.Width, ctx.Height), Checks: checks, Metrics: metrics}
}
