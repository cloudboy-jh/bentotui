package layouts

import "github.com/cloudboy-jh/bentotui/registry/layouts/internal/engine"

// BigTopStrip:
// +----------------------+
// |                      |
// |      primary         |
// |                      |
// +----------------------+
// | strip                |
// +----------------------+
// BigTopStrip renders a large primary area and a fixed-height bottom strip.
func BigTopStrip(width, height, stripH int, primary, strip Sizable) string {
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fill}, {Kind: engine.Fixed, N: stripH}},
		[]Sizable{primary, strip},
	)
}
