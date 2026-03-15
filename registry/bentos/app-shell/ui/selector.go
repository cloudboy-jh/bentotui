package ui

import (
	"fmt"
	"strings"

	"github.com/cloudboy-jh/bentotui/registry/bentos/app-shell/scenarios"
)

func SelectorText(items []scenarios.Definition, active int, focusOwner string) string {
	rows := make([]string, 0, len(items)+3)
	rows = append(rows, "Validation scenarios", "")
	for i, s := range items {
		prefix := "  "
		if i == active {
			prefix = "> "
		}
		rows = append(rows, fmt.Sprintf("%s%d. %s", prefix, i+1, s.Title))
	}
	rows = append(rows, "", "focus owner: "+focusOwner)
	return strings.Join(rows, "\n")
}
