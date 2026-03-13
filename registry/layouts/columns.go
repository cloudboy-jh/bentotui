package layouts

import "github.com/cloudboy-jh/bentotui/registry/layouts/internal/engine"

// TripleCol:
// +-----+---------+------+
// |     |         |      |
// | nav |  list   |detail|
// |     |         |      |
// +-----+---------+------+
// TripleCol renders nav, list, and detail columns.
func TripleCol(width, height, navW, listW int, nav, list, detail Sizable) string {
	return engine.RenderHorizontal(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: navW}, {Kind: engine.Fixed, N: listW}, {Kind: engine.Fill}},
		[]Sizable{nav, list, detail},
	)
}
