package rooms

import "github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"

// Rail:
// +--------+-------------+
// |        |             |
// |  rail  |    main     |
// |        |             |
// +--------+-------------+
// Rail renders a fixed-width rail and a flexible main area.
func Rail(width, height, railWidth int, rail, main Sizable) string {
	return engine.RenderHorizontal(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: railWidth}, {Kind: engine.Fill}},
		[]Sizable{rail, main},
	)
}

// HolyGrail:
// +----------------------+
// | header               |
// +--------+-------------+
// |        |             |
// |  rail  |    main     |
// |        |             |
// +--------+-------------+
// | footer               |
// +----------------------+
// HolyGrail renders header, rail+main body, and footer.
func HolyGrail(width, height, railWidth int, header, rail, main, footer Sizable) string {
	bodyH := engine.Max(1, height-2)
	body := Rail(width, bodyH, railWidth, rail, main)
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: 1}, {Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{header, Static(body), footer},
	)
}
