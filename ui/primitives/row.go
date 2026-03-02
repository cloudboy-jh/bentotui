package primitives

import (
	"charm.land/lipgloss/v2"
)

func RenderRow(width int, bg, fg, content string) string {
	if width <= 0 {
		return ""
	}
	style := lipgloss.NewStyle()
	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}
	if fg != "" {
		style = style.Foreground(lipgloss.Color(fg))
	}
	return renderStyledContent(style, width, content)
}

func RenderStyledRow(style lipgloss.Style, width int, content string) string {
	if width <= 0 {
		return ""
	}
	return renderStyledContent(style, width, content)
}

func clipRowContent(content string, width int) string {
	if content == "" || width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().MaxWidth(width).Render(content)
}

func renderStyledContent(style lipgloss.Style, width int, content string) string {
	base := style.Render(lipgloss.NewStyle().Width(width).Render(""))
	line := clipRowContent(content, width)
	if line == "" {
		return base
	}
	return lipgloss.NewCanvas(
		lipgloss.NewLayer(base).X(0).Y(0).Z(0),
		lipgloss.NewLayer(line).X(0).Y(0).Z(1),
	).Render()
}
