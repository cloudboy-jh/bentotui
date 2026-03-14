package rooms

import "github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"

// Focus:
// +----------------------+
// |                      |
// |       content        |
// |                      |
// +----------------------+
// | footer               |
// +----------------------+
// Focus renders full content plus a one-row footer.
func Focus(width, height int, content, footer Sizable) string {
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{content, footer},
	)
}
