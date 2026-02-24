package surface

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func Fill(width, height int, bg string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	block := blankBlock(width, height)
	if bg == "" {
		return block
	}
	return lipgloss.NewStyle().Background(lipgloss.Color(bg)).Render(block)
}

func Region(content string, width, height int, bg, fg string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	rows := make([]string, 0, height)
	lines := strings.Split(content, "\n")
	for i := 0; i < height; i++ {
		if i >= len(lines) {
			rows = append(rows, strings.Repeat(" ", width))
			continue
		}
		rows = append(rows, FitWidth(lines[i], width))
	}
	block := strings.Join(rows, "\n")

	style := lipgloss.NewStyle()
	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}
	if fg != "" {
		style = style.Foreground(lipgloss.Color(fg))
	}
	if bg == "" && fg == "" {
		return block
	}
	return style.Render(block)
}

func FitWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	s = lipgloss.NewStyle().MaxWidth(width).Render(s)
	return lipgloss.PlaceHorizontal(width, lipgloss.Left, s)
}

func blankBlock(width, height int) string {
	line := strings.Repeat(" ", width)
	rows := make([]string, height)
	for i := range rows {
		rows[i] = line
	}
	return strings.Join(rows, "\n")
}
