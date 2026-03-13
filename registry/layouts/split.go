package layouts

import "github.com/cloudboy-jh/bentotui/registry/layouts/internal/engine"

// HSplit:
// +----------+-----------+
// |          |           |
// |   left   |   right   |
// |          |           |
// +----------+-----------+
// HSplit renders two equal side-by-side panels.
func HSplit(width, height int, left, right Sizable) string {
	return engine.RenderHorizontal(width, height,
		[]engine.Spec{{Kind: engine.Ratio, N: 1}, {Kind: engine.Ratio, N: 1}},
		[]Sizable{left, right},
	)
}

// VSplit:
// +----------------------+
// |         top          |
// |                      |
// +----------------------+
// |       bottom         |
// |                      |
// +----------------------+
// VSplit renders two equal stacked panels.
func VSplit(width, height int, top, bottom Sizable) string {
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Ratio, N: 1}, {Kind: engine.Ratio, N: 1}},
		[]Sizable{top, bottom},
	)
}

// HSplitFooter:
// +----------+-----------+
// |          |           |
// |   left   |   right   |
// |          |           |
// +----------+-----------+
// | footer               |
// +----------------------+
// HSplitFooter renders two equal side-by-side panels plus a footer.
func HSplitFooter(width, height int, left, right, footer Sizable) string {
	bodyH := engine.Max(1, height-1)
	body := HSplit(width, bodyH, left, right)
	return engine.RenderVertical(width, height,
		[]engine.Spec{{Kind: engine.Fill}, {Kind: engine.Fixed, N: 1}},
		[]Sizable{Static(body), footer},
	)
}
