package layouts

import "github.com/cloudboy-jh/bentotui/registry/layouts/internal/engine"

// Dashboard2x2:
// +----------+-----------+
// |    tl    |    tr     |
// +----------+-----------+
// |    bl    |    br     |
// +----------+-----------+
// Dashboard2x2 renders four equal quadrants.
func Dashboard2x2(width, height int, tl, tr, bl, br Sizable) string {
	rowH := engine.Max(1, height/2)
	top := HSplit(width, rowH, tl, tr)
	bottom := HSplit(width, engine.Max(1, height-rowH), bl, br)

	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Ratio, N: 1}, {Kind: engine.Ratio, N: 1}},
		[]Sizable{Static(top), Static(bottom)},
	)
}

// Dashboard2x2Footer:
// +----------+-----------+
// |    tl    |    tr     |
// +----------+-----------+
// |    bl    |    br     |
// +----------+-----------+
// | footer               |
// +----------------------+
// Dashboard2x2Footer renders four equal quadrants plus a footer.
func Dashboard2x2Footer(width, height int, tl, tr, bl, br, footer Sizable) string {
	bodyH := engine.Max(1, height-1)
	body := Dashboard2x2(width, bodyH, tl, tr, bl, br)

	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{Static(body), footer},
	)
}
