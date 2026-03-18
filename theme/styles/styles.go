package styles

// styles provides pure rendering utilities for bento bricks.
// There is no System struct, no theme dependency in this package.
// Bricks use the Theme interface directly — styles just provides
// the row-painting helpers that enforce the UV cell contract.

import (
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"image/color"
)

// Row returns a fully-painted row string of exactly width cells.
//
// This is the canonical way to render any row in a component or bento.
// Every cell has an explicit Bg so the Ultraviolet surface overlay does not
// fall back to the canvas color for padding/whitespace cells.
//
// Rule: never use lipgloss.PlaceHorizontal or bare Render(content) for rows
// that sit on a surface — always go through Row() or an equivalent
// .Background().Width(w).Render() chain.
func Row(bg, fg color.Color, width int, content string) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Background(bg).
		Foreground(fg).
		Width(width).
		Render(content)
}

// ClipANSI truncates styled text to width cells safely.
func ClipANSI(content string, width int) string {
	if width <= 0 {
		return ""
	}
	return ansi.Truncate(content, width, "")
}

// RowClip clips ANSI content first, then paints an exact-width row.
func RowClip(bg, fg color.Color, width int, content string) string {
	if width <= 0 {
		return ""
	}
	return Row(bg, fg, width, ClipANSI(content, width))
}

// InputStyles returns a fully-themed textinput.Styles struct.
// Kept here as a utility so bricks don't need to import bubbles/textinput directly.
