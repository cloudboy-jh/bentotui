package rooms

import "github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"

// Pancake:
// +----------------------+
// | header               |
// +----------------------+
// |                      |
// |       content        |
// |                      |
// +----------------------+
// | footer               |
// +----------------------+
// Pancake renders header, content, and footer.
func Pancake(width, height int, header, content, footer Sizable) string {
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: 1}, {Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{header, content, footer},
	)
}

// TopbarPancake:
// +----------------------+
// | topbar               |
// +----------------------+
// | header               |
// +----------------------+
// |                      |
// |       content        |
// |                      |
// +----------------------+
// | footer               |
// +----------------------+
// TopbarPancake renders topbar, header, content, and footer.
func TopbarPancake(width, height int, topbar, header, content, footer Sizable) string {
	return Frame(width, height, topbar, header, content, footer)
}
