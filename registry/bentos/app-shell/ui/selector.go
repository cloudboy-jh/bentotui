package ui

import (
	"fmt"
	"strings"

	"github.com/cloudboy-jh/bentotui/registry/bentos/app-shell/scenarios"
)

func SelectorText(items []scenarios.Definition, active int) string {
	rows := make([]string, 0, len(items)+2)
	rows = append(rows, "Validation scenarios", "")
	for i, s := range items {
		prefix := "  "
		if i == active {
			prefix = "> "
		}
		rows = append(rows, fmt.Sprintf("%s%d. %s", prefix, i+1, s.Title))
	}
	return strings.Join(rows, "\n")
}
