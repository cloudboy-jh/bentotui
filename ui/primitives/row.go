package primitives

import (
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core/surface"
)

func PaintRow(width int, bg, fg, content string) string {
	if width <= 0 {
		return ""
	}
	line := surface.FitWidth(content, width)
	style := lipgloss.NewStyle().Width(width)
	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}
	if fg != "" {
		style = style.Foreground(lipgloss.Color(fg))
	}
	return style.Render(line)
}

func PaintStyledRow(style lipgloss.Style, width int, content string) string {
	if width <= 0 {
		return ""
	}
	line := surface.FitWidth(content, width)
	return style.Width(width).Render(line)
}
