package primitives

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func RenderFrame(style lipgloss.Style, width, height int, rows []string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	return style.Width(width).Height(height).Render(strings.Join(rows, "\n"))
}
