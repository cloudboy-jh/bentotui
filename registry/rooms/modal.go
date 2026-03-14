package rooms

import "github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"

// Modal:
// +----------------------+
// |  +----------------+  |
// |  |     modal      |  |
// |  +----------------+  |
// +----------------------+
// Modal renders background content with a centered modal overlay.
func Modal(width, height, modalW, modalH int, background, modal Sizable) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	bg := renderCell(background, width, height)
	mw := engine.Min(engine.Max(1, modalW), width)
	mh := engine.Min(engine.Max(1, modalH), height)
	fg := renderCell(modal, mw, mh)

	x := engine.Max(0, (width-mw)/2)
	y := engine.Max(0, (height-mh)/2)

	return engine.Constrain(engine.Overlay(bg, fg, x, y), width, height)
}

func renderCell(cell Sizable, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	cell.SetSize(width, height)
	return engine.Constrain(engine.ViewString(cell.View()), width, height)
}
