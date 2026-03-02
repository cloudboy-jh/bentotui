package primitives

import (
	"testing"

	"charm.land/lipgloss/v2"
)

func TestRenderRowKeepsAssignedWidthWithStyledContent(t *testing.T) {
	content := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("hello")
	row := RenderRow(24, "1", "7", content)
	if got := lipgloss.Width(row); got != 24 {
		t.Fatalf("expected width 24, got %d", got)
	}
}

func TestRenderStyledRowKeepsAssignedWidthWithStyledContent(t *testing.T) {
	style := lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("7"))
	content := lipgloss.NewStyle().Bold(true).Render("title")
	row := RenderStyledRow(style, 18, content)
	if got := lipgloss.Width(row); got != 18 {
		t.Fatalf("expected width 18, got %d", got)
	}
}
