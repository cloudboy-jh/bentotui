package surface

import (
	"charm.land/lipgloss/v2"
)

func Fill(width, height int, bg string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(lipgloss.Color(bg)).
		Render("")
}

func Region(content string, width, height int, bg, fg string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	style := lipgloss.NewStyle().
		Width(width).
		Height(height)
	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}
	if fg != "" {
		style = style.Foreground(lipgloss.Color(fg))
	}
	return style.Render(content)
}

func FitWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	s = lipgloss.NewStyle().MaxWidth(width).Render(s)
	return lipgloss.PlaceHorizontal(width, lipgloss.Left, s)
}
