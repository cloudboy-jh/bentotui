package ui

import (
	"fmt"
	"strings"
)

func SelectorText(sections []string, active int, themeName string, progress float64, compact bool) string {
	rows := make([]string, 0, len(sections)+10)
	rows = append(rows, "Sections")
	for i, s := range sections {
		prefix := "  "
		if i == active {
			prefix = "> "
		}
		rows = append(rows, fmt.Sprintf("%s%d. %s", prefix, i+1, s))
	}
	rows = append(rows,
		"",
		"Rail",
		fmt.Sprintf("  Theme     %s", themeName),
		fmt.Sprintf("  Progress  %3.0f%%", progress*100),
		fmt.Sprintf("  Table     %s", ternary(compact, "compact", "comfortable")),
	)
	return strings.Join(rows, "\n")
}

func ternary[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}
