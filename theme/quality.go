package theme

var curatedStableThemes = map[string]struct{}{
	"catppuccin-mocha":  {},
	"tokyo-night":       {},
	"tokyo-night-storm": {},
	"one-dark":          {},
	"github-dark":       {},
	"nord":              {},
	"bento-rose":        {},
}

func classifyThemeTier(name string, t Theme) ThemeTier {
	if _, ok := curatedStableThemes[name]; ok {
		return ThemeTierStable
	}
	if themeQualityScore(t) >= 82 {
		return ThemeTierStable
	}
	return ThemeTierExperimental
}

func themeQualityScore(t Theme) float64 {
	type pair struct {
		a, b string
		min  float64
		w    float64
	}
	pairs := []pair{
		{t.Surface.Panel, t.Surface.Canvas, minSurfacePanelCanvasDelta, 1.2},
		{t.Surface.Interactive, t.Surface.Panel, minSurfaceInteractivePanelDelta, 1.0},
		{t.Input.BG, t.Surface.Canvas, 0.03, 1.0},
		{t.Selection.BG, t.Surface.Canvas, 0.05, 1.2},
		{t.Selection.BG, t.Input.BG, 0.05, 1.0},
		{t.Dialog.BG, t.Surface.Canvas, 0.03, 1.0},
		{t.Card.HeaderBG, t.Card.BodyBG, minCardHeaderBodyDelta, 1.1},
		{t.Card.FrameBG, t.Card.BodyBG, minCardFrameBodyDelta, 1.0},
		{t.Card.FocusEdgeBG, t.Card.FrameBG, minCardFocusEdgeFrameDelta, 1.3},
		{t.Card.ShadowBG, t.Surface.Canvas, minCardShadowCanvasDelta, 0.7},
		{t.Card.FrameFG, t.Card.FrameBG, 0.10, 1.0},
		{t.Text.Primary, t.Card.BodyBG, 0.10, 1.0},
	}

	totalW := 0.0
	acc := 0.0
	for _, p := range pairs {
		d := lumDelta(p.a, p.b)
		r := d / (p.min * 2.0)
		if r > 1.0 {
			r = 1.0
		}
		if r < 0 {
			r = 0
		}
		acc += r * p.w
		totalW += p.w
	}
	if totalW == 0 {
		return 0
	}
	return (acc / totalW) * 100.0
}
