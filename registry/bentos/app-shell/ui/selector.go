package ui

import (
	"fmt"
	"strings"

	"github.com/cloudboy-jh/bentotui/registry/bentos/app-shell/scenarios"
)

func SelectorText(items []scenarios.Definition, active int) string {
	rows := make([]string, 0, len(items)+12)
	rows = append(rows, "Component scenarios", "")
	for i, s := range items {
		prefix := "  "
		if i == active {
			prefix = "> "
		}
		rows = append(rows, fmt.Sprintf("%s%d. %s", prefix, i+1, s.Title))
	}
	rows = append(rows,
		"",
		"Areas",
		"  1. Rail/Nav",
		"  2. Main/Canvas",
		"  3. Checks",
		"  4. Context",
		"  5. Session",
		"  6. Commands",
	)
	return strings.Join(rows, "\n")
}
