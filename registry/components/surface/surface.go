// Package surface provides a deterministic full-screen paint surface backed
// by the Ultraviolet cell buffer. It is the registry-level answer to ANSI
// whitespace-reset bleed: instead of composing styled strings with
// PlaceHorizontal, callers draw lipgloss-rendered blocks onto a pre-filled
// cell buffer so every cell is explicitly painted before Bubble Tea flushes
// the frame.
//
// Usage pattern (in your root model's View):
//
//	s := surface.New(m.width, m.height)
//	s.Fill(t.Surface.Canvas)          // paint every cell with bg color
//	s.Draw(0, 0, bodyContent)         // place body string at top-left
//	s.DrawCenter(m.width, m.height, dialogStr) // center a dialog overlay
//
//	v := tea.NewView(s.Render())
//	v.AltScreen = true
//	v.BackgroundColor = lipgloss.Color(t.Surface.Canvas)
//	return v
//
// Copy this file into your project: bento add surface
//
// Dependencies:
//   - charm.land/bubbletea/v2
//   - github.com/charmbracelet/ultraviolet
package surface

import (
	"image/color"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

// Surface is a full-screen cell buffer. Build one per frame in View(),
// fill it with a background color, draw your rendered component strings
// onto it, then call Render() to get the final frame string.
type Surface struct {
	buf    uv.ScreenBuffer
	width  int
	height int
}

// New creates a Surface sized to width × height cells.
func New(width, height int) *Surface {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	return &Surface{
		buf:    uv.NewScreenBuffer(width, height),
		width:  width,
		height: height,
	}
}

// Fill paints every cell in the surface with a space character styled with
// the given background color. Pass a lipgloss.Color string value cast to
// color.Color, or use image/color directly.
// Typically called once at the start of View() before any Draw calls.
func (s *Surface) Fill(bg color.Color) {
	cell := &uv.Cell{
		Content: " ",
		Width:   1,
		Style:   uv.Style{Bg: bg},
	}
	s.buf.Fill(cell)
}

// Draw places a pre-rendered ANSI string at position (x, y) on the surface.
// The string is parsed cell-by-cell so styled blocks from lipgloss land
// exactly where you intend them, with no surrounding whitespace resets.
func (s *Surface) Draw(x, y int, content string) {
	if content == "" {
		return
	}
	area := uv.Rect(x, y, s.width, s.height)
	uv.NewStyledString(content).Draw(s.buf, area)
}

// DrawCenter places a pre-rendered ANSI string centered within the surface.
// Use this for dialogs and overlays.
func (s *Surface) DrawCenter(content string) {
	if content == "" {
		return
	}
	lines := strings.Split(content, "\n")
	contentH := len(lines)
	contentW := 0
	for _, l := range lines {
		// Use lipgloss width-aware measurement via the styled string.
		ss := uv.NewStyledString(l)
		if w := ss.UnicodeWidth(); w > contentW {
			contentW = w
		}
	}
	x := max(0, (s.width-contentW)/2)
	y := max(0, (s.height-contentH)/2)
	s.Draw(x, y, content)
}

// Render returns the final ANSI frame string to pass to tea.NewView.
// Bubble Tea normalises \r\n, so we emit raw buffer output.
func (s *Surface) Render() string {
	return s.buf.Render()
}

// Width returns the surface width in cells.
func (s *Surface) Width() int { return s.width }

// Height returns the surface height in cells.
func (s *Surface) Height() int { return s.height }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
