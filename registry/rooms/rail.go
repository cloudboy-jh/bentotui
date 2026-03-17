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

// RailFooterStack:
// +--------+-------------+
// |  rail  |    main     |
// +----------------------+
// |   footer card rows   |
// +----------------------+
// |      footer bar      |
// +----------------------+
// RailFooterStack renders rail+main body with optional footer-card rows and a
// required final footer bar row.
func RailFooterStack(width, height, railWidth, footerCardRows int, rail, main, footerCard, footerBar Sizable) string {
	bodyH := engine.Max(1, height-1-footerCardRows)
	body := Rail(width, bodyH, railWidth, rail, main)
	if footerCardRows <= 0 || footerCard == nil {
		return engine.RenderVertical(width, height,
			[]engine.Spec{{Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
			[]Sizable{Static(body), footerBar},
		)
	}
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fill}, {Kind: engine.Fixed, N: footerCardRows}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{Static(body), footerCard, footerBar},
	)
}
