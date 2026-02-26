package primitives

import "strings"

func PaintInputRow(width int, bg, fg, content string) string {
	return PaintRow(width, bg, fg, content)
}

func PaintInputRowInset(width int, bg, fg, content string, inset int) string {
	if inset < 0 {
		inset = 0
	}
	return PaintRow(width, bg, fg, strings.Repeat(" ", inset)+content)
}
