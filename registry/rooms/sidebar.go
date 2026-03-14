package rooms

import "github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"

// Sidebar:
// +--------+-------------+
// |        |             |
// |sidebar |    main     |
// |        |             |
// +--------+-------------+
// Sidebar renders a fixed-width sidebar and a flexible main area.
func Sidebar(width, height, sideWidth int, sidebar, main Sizable) string {
	return engine.RenderHorizontal(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: sideWidth}, {Kind: engine.Fill}},
		[]Sizable{sidebar, main},
	)
}

// HolyGrail:
// +----------------------+
// | header               |
// +--------+-------------+
// |        |             |
// |sidebar |    main     |
// |        |             |
// +--------+-------------+
// | footer               |
// +----------------------+
// HolyGrail renders header, sidebar+main body, and footer.
func HolyGrail(width, height, sideWidth int, header, sidebar, main, footer Sizable) string {
	bodyH := engine.Max(1, height-2)
	body := Sidebar(width, bodyH, sideWidth, sidebar, main)
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: 1}, {Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{header, Static(body), footer},
	)
}
