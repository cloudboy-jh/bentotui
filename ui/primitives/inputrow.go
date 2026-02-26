package primitives

import "strings"

func RenderInputRow(width int, bg, fg, content string) string {
	return RenderRow(width, bg, fg, content)
}

func RenderInputRowInset(width int, bg, fg, content string, inset int) string {
	if inset < 0 {
		inset = 0
	}
	return RenderRow(width, bg, fg, strings.Repeat(" ", inset)+content)
}
