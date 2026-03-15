package ui

import "fmt"

func CanvasHeader(active, total int, title, description string, compact bool) []string {
	rows := []string{
		fmt.Sprintf("[%d/%d] %s", active, total, title),
		description,
		"",
	}
	if compact {
		rows = append(rows, "compact mode: diagnostics collapsed into canvas", "")
	}
	return rows
}
