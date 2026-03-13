package layouts

import "github.com/cloudboy-jh/bentotui/registry/layouts/internal/engine"

// DrawerRight:
// +--------------+-------+
// |              |       |
// |    main      |drawer |
// |              |       |
// +--------------+-------+
// DrawerRight renders main content and a fixed-width right drawer.
func DrawerRight(width, height, drawerW int, main, drawer Sizable) string {
	return engine.RenderHorizontal(width, height,
		[]engine.Spec{{Kind: engine.Fill}, {Kind: engine.Fixed, N: drawerW}},
		[]Sizable{main, drawer},
	)
}

// DrawerChrome:
// +----------------------+
// | header               |
// +--------------+-------+
// |              |       |
// |    main      |drawer |
// |              |       |
// +--------------+-------+
// | footer               |
// +----------------------+
// DrawerChrome renders header, body with right drawer, and footer.
func DrawerChrome(width, height, drawerW int, header, main, drawer, footer Sizable) string {
	bodyH := engine.Max(1, height-2)
	body := DrawerRight(width, bodyH, drawerW, main, drawer)

	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fixed, N: 1}, {Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{header, Static(body), footer},
	)
}
