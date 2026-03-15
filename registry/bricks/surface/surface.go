// Brick: Surface:
// +-----------------------------------+
// | Fill(canvas)                      |
// | Draw(x,y,content)                 |
// | DrawCenter(overlay)               |
// +-----------------------------------+
// Root compositor for full-frame paint.
// Package surface provides a deterministic full-terminal paint surface backed
// by the Ultraviolet cell buffer.
//
// Surface is the root canvas for every bento and full-screen layout.
// The contract:
//
//  1. Create one surface per frame sized to the full terminal (width x height).
//  2. Call Fill(bg) once — paints every cell with the canvas background.
//  3. Call Draw(x, y, content) for each component — overlays only the content's
//     own styled cells; surrounding filled cells are never touched.
//  4. Call Render() once and pass the result to tea.NewView.
//
// This keeps the Ultraviolet cell buffer as the single source of truth for
// what every terminal cell contains on each frame — no ANSI whitespace-reset
// bleed, no partial clears, no string concatenation outside this surface.
//
// Copy this file into your project: bento add surface
//
// Dependencies:
//   - github.com/charmbracelet/ultraviolet
package surface

import (
	"image/color"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

// Surface is the full-terminal cell buffer. Build one per frame in View(),
// fill it, draw components onto it, then call Render().
type Surface struct {
	buf    uv.ScreenBuffer
	width  int
	height int
}

// New creates a Surface sized to the full terminal (width x height).
// Call this at the top of every View() with the current terminal dimensions.
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

// Fill paints every cell of the surface with a space styled with bg.
// Always call this first — it is the root canvas layer that every
// component draws on top of.
func (s *Surface) Fill(bg color.Color) {
	cell := &uv.Cell{
		Content: " ",
		Width:   1,
		Style:   uv.Style{Bg: bg},
	}
	s.buf.Fill(cell)
}

// overlayScreen wraps the real ScreenBuffer but intercepts SetCell so that
// nil cells (which StyledString.Draw uses to pre-clear its area) are silently
// dropped. Only real glyph cells are written through. This makes Draw an
// overlay operation — the filled background beneath is never wiped.
type overlayScreen struct {
	uv.ScreenBuffer
}

func (o overlayScreen) SetCell(x, y int, c *uv.Cell) {
	if c == nil {
		// Drop the pre-clear — preserve the filled background beneath.
		return
	}
	// If the incoming cell has no background set, inherit the background
	// from the already-filled cell beneath it. This prevents dark seams
	// around inline text whose Bg is nil after lipgloss renders it.
	if c.Style.Bg == nil {
		if beneath := o.ScreenBuffer.CellAt(x, y); beneath != nil {
			clone := *c
			clone.Style.Bg = beneath.Style.Bg
			c = &clone
		}
	}
	o.ScreenBuffer.SetCell(x, y, c)
}

// Draw places a pre-rendered ANSI string at (x, y) as an overlay.
// Only the content's own styled cells are written; all surrounding cells
// filled by Fill() are left completely untouched.
func (s *Surface) Draw(x, y int, content string) {
	if content == "" || x >= s.width || y >= s.height {
		return
	}
	ss := uv.NewStyledString(content)
	b := ss.Bounds()
	w := min(b.Dx(), s.width-x)
	h := min(b.Dy(), s.height-y)
	if w <= 0 || h <= 0 {
		return
	}
	// Use the overlay screen so pre-clear SetCell(nil) calls are dropped.
	ov := overlayScreen{s.buf}
	ss.Draw(ov, uv.Rect(x, y, w, h))
}

// DrawCenter places a pre-rendered ANSI string centered on the surface.
// Use this for dialogs and overlays — they draw on top of the filled bg.
func (s *Surface) DrawCenter(content string) {
	if content == "" {
		return
	}
	ss := uv.NewStyledString(content)
	b := ss.Bounds()
	x := max(0, (s.width-b.Dx())/2)
	y := max(0, (s.height-b.Dy())/2)
	s.Draw(x, y, content)
}

// Render serializes the cell buffer to an ANSI string for tea.NewView.
// uv.Buffer.Render() emits \r\n; we normalize to \n so Bubble Tea does
// not emit a raw carriage return that resets the cursor column mid-frame.
func (s *Surface) Render() string {
	return strings.ReplaceAll(s.buf.Render(), "\r\n", "\n")
}

// Width returns the surface width in cells.
func (s *Surface) Width() int { return s.width }

// Height returns the surface height in cells.
func (s *Surface) Height() int { return s.height }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
