package theme

// Semantic surface ladder:
// +-------------------------------+
// | surface.canvas                |
// +-------------------------------+
// | surface.panel                 |
// +-------------------------------+
// | surface.interactive           |
// +-------------------------------+
//
// Elevated card slab model:
// +-------------------------------+
// | card.frameBG                  |
// | card.headerBG                 |
// | card.bodyBG                   |
// | card.footerBG                 |
// | card.shadowBG                 |
// +-------------------------------+
// `card.focusEdgeBG` is the accent stripe used for focused cards.
// Footer anchored colors remain separate semantic tokens so command rows do not
// need to reuse selection-highlight intensity.

import (
	"fmt"

	tint "github.com/lrstanley/bubbletint/v2"
)

type SurfaceTokens struct {
	Canvas      string
	Panel       string
	Overlay     string
	Interactive string
}

type TextTokens struct {
	Primary string
	Muted   string
	Inverse string
	Accent  string
}

type BorderTokens struct {
	Normal string
	Subtle string
	Focus  string
}

type StateTokens struct {
	Info    string
	Success string
	Warning string
	Danger  string
}

type SelectionTokens struct {
	BG string
	FG string
}

type InputTokens struct {
	BG          string
	FG          string
	Placeholder string
	Cursor      string
	Border      string
}

type BarTokens struct {
	BG string
	FG string
}

type FooterTokens struct {
	AnchoredBG    string
	AnchoredFG    string
	AnchoredMuted string
}

type DialogTokens struct {
	BG     string
	FG     string
	Border string
	Scrim  string
}

type CardTokens struct {
	HeaderBG    string
	BodyBG      string
	FooterBG    string
	FrameBG     string
	FrameFG     string
	ShadowBG    string
	FocusEdgeBG string
}

type Theme struct {
	Name      string
	Surface   SurfaceTokens
	Text      TextTokens
	Border    BorderTokens
	State     StateTokens
	Selection SelectionTokens
	Input     InputTokens
	Bar       BarTokens
	Footer    FooterTokens
	Dialog    DialogTokens
	Card      CardTokens
}

type ThemeTier string

const (
	ThemeTierStable       ThemeTier = "stable"
	ThemeTierExperimental ThemeTier = "experimental"
)

type ThemeMeta struct {
	Tier  ThemeTier
	Score float64
}

const (
	minSurfacePanelCanvasDelta      = 0.045
	minSurfaceInteractivePanelDelta = 0.030
	minCardHeaderBodyDelta          = 0.025
	minCardFrameBodyDelta           = 0.020
	minCardShadowCanvasDelta        = 0.015
	minCardFocusEdgeFrameDelta      = 0.040
)

const (
	DefaultName         = "catppuccin-mocha"
	CatppuccinMochaName = "catppuccin-mocha"
	DraculaName         = "dracula"
)

func AvailableThemes() []string {
	return availableThemeNames()
}

func AvailableStableThemes() []string {
	return availableStableThemeNames()
}

func ThemeMetadata(name string) (ThemeMeta, bool) {
	return themeMetadata(name)
}

func Preset(name string) Theme {
	t, ok := presetTheme(name)
	if !ok {
		t, _ = presetTheme(DefaultName)
	}
	return t
}

func presetTheme(name string) (Theme, bool) {
	t, ok := builtinThemes[name]
	if !ok {
		return Theme{}, false
	}
	if t.Name == "" {
		t.Name = name
	}
	return t, true
}

func validateTheme(t Theme) error {
	// ── required tokens ───────────────────────────────────────────────────────
	required := []struct {
		label string
		value string
	}{
		{"surface.canvas", t.Surface.Canvas},
		{"surface.panel", t.Surface.Panel},
		{"surface.overlay", t.Surface.Overlay},
		{"surface.interactive", t.Surface.Interactive},
		{"text.primary", t.Text.Primary},
		{"text.muted", t.Text.Muted},
		{"text.inverse", t.Text.Inverse},
		{"text.accent", t.Text.Accent},
		{"border.normal", t.Border.Normal},
		{"border.subtle", t.Border.Subtle},
		{"border.focus", t.Border.Focus},
		{"state.info", t.State.Info},
		{"state.success", t.State.Success},
		{"state.warning", t.State.Warning},
		{"state.danger", t.State.Danger},
		{"selection.bg", t.Selection.BG},
		{"selection.fg", t.Selection.FG},
		{"input.bg", t.Input.BG},
		{"input.fg", t.Input.FG},
		{"input.placeholder", t.Input.Placeholder},
		{"input.cursor", t.Input.Cursor},
		{"input.border", t.Input.Border},
		{"bar.bg", t.Bar.BG},
		{"bar.fg", t.Bar.FG},
		{"dialog.bg", t.Dialog.BG},
		{"dialog.fg", t.Dialog.FG},
		{"dialog.border", t.Dialog.Border},
		{"dialog.scrim", t.Dialog.Scrim},
		{"card.headerBG", t.Card.HeaderBG},
		{"card.bodyBG", t.Card.BodyBG},
		{"card.footerBG", t.Card.FooterBG},
		{"card.frameBG", t.Card.FrameBG},
		{"card.frameFG", t.Card.FrameFG},
		{"card.shadowBG", t.Card.ShadowBG},
		{"card.focusEdgeBG", t.Card.FocusEdgeBG},
	}
	for _, c := range required {
		if c.value == "" {
			return fmt.Errorf("theme token %q is required", c.label)
		}
	}

	footerProvided := 0
	for _, v := range []string{t.Footer.AnchoredBG, t.Footer.AnchoredFG, t.Footer.AnchoredMuted} {
		if v != "" {
			footerProvided++
		}
	}
	if footerProvided > 0 && footerProvided < 3 {
		return fmt.Errorf("footer anchored tokens must be set together: footer.anchoredBG, footer.anchoredFG, footer.anchoredMuted")
	}

	// ── layer separation checks ───────────────────────────────────────────────
	// Key layer pairs must be visually distinct. Thresholds are calibrated to
	// the minimum detectable contrast in dark terminal themes:
	//   - input.bg vs canvas: 0.03 (raised surface, subtle but visible)
	//   - selection.bg vs canvas: 0.05 (must pop clearly)
	//   - selection.bg vs input.bg: 0.05 (selected row must stand out from field)
	//   - dialog.bg vs canvas: 0.03 (dialog body is a raised surface)
	layerPairs := []struct {
		labelA, labelB string
		a, b           string
		minDelta       float64
	}{
		{"surface.panel", "surface.canvas", t.Surface.Panel, t.Surface.Canvas, minSurfacePanelCanvasDelta},
		{"surface.interactive", "surface.panel", t.Surface.Interactive, t.Surface.Panel, minSurfaceInteractivePanelDelta},
		{"input.bg", "surface.canvas", t.Input.BG, t.Surface.Canvas, 0.03},
		{"selection.bg", "surface.canvas", t.Selection.BG, t.Surface.Canvas, 0.05},
		{"selection.bg", "input.bg", t.Selection.BG, t.Input.BG, 0.05},
		{"dialog.bg", "surface.canvas", t.Dialog.BG, t.Surface.Canvas, 0.03},
		{"bar.bg", "surface.canvas", t.Bar.BG, t.Surface.Canvas, 0.02},
		{"card.headerBG", "card.bodyBG", t.Card.HeaderBG, t.Card.BodyBG, minCardHeaderBodyDelta},
		{"card.frameBG", "card.bodyBG", t.Card.FrameBG, t.Card.BodyBG, minCardFrameBodyDelta},
		{"card.shadowBG", "surface.canvas", t.Card.ShadowBG, t.Surface.Canvas, minCardShadowCanvasDelta},
		{"card.focusEdgeBG", "card.frameBG", t.Card.FocusEdgeBG, t.Card.FrameBG, minCardFocusEdgeFrameDelta},
	}
	for _, p := range layerPairs {
		delta := lumDelta(p.a, p.b)
		if delta < p.minDelta {
			return fmt.Errorf(
				"theme tokens %q and %q are too similar (luminance delta %.3f < %.3f) — increase contrast",
				p.labelA, p.labelB, delta, p.minDelta,
			)
		}
	}

	return nil
}

// builtinThemes drives the theme registry. All entries are derived from
// well-known iTerm2-Color-Schemes palettes via the bubbletint adapter so
// that every theme has professionally designed contrast ratios.
var builtinThemes = map[string]Theme{
	// Catppuccin family
	"catppuccin-mocha":     fromTint(tint.TintCatppuccinMocha, "catppuccin-mocha"),
	"catppuccin-macchiato": fromTint(tint.TintCatppuccinMacchiato, "catppuccin-macchiato"),
	"catppuccin-frappe":    fromTint(tint.TintCatppuccinFrappe, "catppuccin-frappe"),

	// Dracula family — TintDraculaPlus has better contrast than base Dracula
	"dracula": fromTint(tint.TintDraculaPlus, "dracula"),

	// Tokyo Night family
	"tokyo-night":       fromTint(tint.TintTokyoNight, "tokyo-night"),
	"tokyo-night-storm": fromTint(tint.TintTokyoNightStorm, "tokyo-night-storm"),

	// Nordic / cool
	"nord": fromTint(tint.TintNord, "nord"),

	// Warm / retro
	"bento-rose":   bentoRoseTheme(),
	"gruvbox-dark": fromTint(tint.TintGruvboxDark, "gruvbox-dark"),
	"monokai-pro":  fromTint(tint.TintMonokaiPro, "monokai-pro"),

	// Earthy / artistic
	"kanagawa":   fromTint(tint.TintKanagawa, "kanagawa"),
	"rose-pine":  fromTint(tint.TintRosePine, "rose-pine"),
	"ayu-mirage": fromTint(tint.TintAyuMirage, "ayu-mirage"),

	// Editor-inspired
	"one-dark":       fromTint(tint.TintOneDark, "one-dark"),
	"material-ocean": fromTint(tint.TintMaterialOcean, "material-ocean"),
	"github-dark":    fromTint(tint.TintGitHubDark, "github-dark"),
}
