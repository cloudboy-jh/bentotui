package layouts

import "github.com/cloudboy-jh/bentotui/registry/layouts/internal/engine"

// Frame renders the canonical Bento screen grammar:
// +----------------------+
// | top                  |
// +----------------------+
// | subheader            |
// +----------------------+
// |                      |
// |       body           |
// |                      |
// +----------------------+
// | subfooter            |
// +----------------------+
func Frame(width, height int, top, subheader, body, subfooter Sizable) string {
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: 1}, {Kind: engine.Fixed, N: 1}, {Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{top, subheader, body, subfooter},
	)
}

// FrameMainDrawer renders Frame grammar with a right-side body drawer.
func FrameMainDrawer(width, height, drawerW int, top, subheader, main, drawer, subfooter Sizable) string {
	bodyH := engine.Max(1, height-3)
	body := DrawerRight(width, bodyH, drawerW, main, drawer)
	return Frame(width, height, top, subheader, Static(body), subfooter)
}

// FrameTriple renders Frame grammar with a nav/list/detail body split.
func FrameTriple(width, height, navW, listW int, top, subheader, nav, list, detail, subfooter Sizable) string {
	bodyH := engine.Max(1, height-3)
	body := TripleCol(width, bodyH, navW, listW, nav, list, detail)
	return Frame(width, height, top, subheader, Static(body), subfooter)
}
