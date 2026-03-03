package theme

import (
	"math"

	tint "github.com/lrstanley/bubbletint/v2"
)

// fromTint maps a bubbletint palette to BentoTUI semantic tokens.
//
// Layer hierarchy (darkest → lightest for dark themes):
//
//	Canvas      — terminal root background (filled by surface.Fill)
//	Panel       — default component body (input block, panels)
//	Elevated    — secondary surfaces (sidebars, nested panels)
//	Overlay     — modal/dialog body
//	Interactive — hover/focus tinted surfaces
//
// The adapter guarantees visual separation between adjacent layers by
// checking relative luminance delta. If a mapped pair is too close it
// falls back to a brighter/darker variant automatically.
func fromTint(t *tint.Tint, name string) Theme {
	if t == nil {
		return Theme{Name: name}
	}

	// ── palette slots ─────────────────────────────────────────────────────────
	// Named for clarity — we reference these below when building token pairs.
	bg := hex(t.Bg, "#1e1e2e")            // true app background
	bbk := hex(t.BrightBlack, "#313244")  // first elevated surface
	blk := hex(t.Black, "#181825")        // deepest dark (darker than bg)
	fg := hex(t.Fg, "#cdd6f4")            // primary text
	wht := hex(t.White, "#6c7086")        // muted text
	bwht := hex(t.BrightWhite, "#bac2de") // secondary text
	bacc := hex(t.BrightBlue, "#89b4fa")  // accent / focus
	bcya := hex(t.BrightCyan, "#89dceb")  // secondary accent
	bgrn := hex(t.BrightGreen, "#a6e3a1")
	bred := hex(t.BrightRed, "#f38ba8")
	byel := hex(t.BrightYellow, "#f9e2af")
	bpur := hex(t.BrightPurple, "#cba4f7")
	cur := orFallback(t.Cursor, t.BrightBlue, "#89b4fa")

	// ── layer assignment ──────────────────────────────────────────────────────
	// Canvas is always t.Bg (the true terminal background).
	// Panel must contrast against Canvas — use BrightBlack (slightly raised).
	// Elevated is darker than Canvas for depth — use Black.
	// Overlay (dialog) uses BrightBlack as base, guaranteed distinct from panel.
	// Interactive is a lighter tint of the accent for hover surfaces.
	canvas := bg
	panel := ensureDistinct(bbk, canvas, blk, bacc)    // raised above canvas
	elevated := ensureDistinct(blk, canvas, bbk, bacc) // below canvas for depth
	overlay := ensureDistinct(bbk, canvas, blk, bacc)  // dialog body
	interactive := bbk                                 // focus-tinted surface

	// Input BG must contrast against canvas but remain a dark surface color —
	// prefer BrightBlack (raised surface), fall back to Black if BrightBlack
	// is too close to canvas. Never fall back to the accent (bacc) — that
	// would collide with selectionBG.
	inputBG := ensureDistinctMin(bbk, canvas, 0.03, blk, elevated)

	// Selection must always be the brightest available accent — clearly
	// distinguishable from both canvas and inputBG.
	// Prefer BrightBlue → BrightCyan → BrightPurple → BrightYellow.
	// Never use surface slots (BrightBlack/Black) which can collide with inputBG.
	selectionBG := pickDistinctFrom([]string{bacc, bcya, bpur, byel}, []string{canvas, inputBG}, 0.05)
	selectionFG := blk // dark text on bright selection

	// ── assemble ──────────────────────────────────────────────────────────────
	return Theme{
		Name: name,
		Surface: SurfaceTokens{
			Canvas:      canvas,
			Panel:       panel,
			Elevated:    elevated,
			Overlay:     overlay,
			Interactive: interactive,
		},
		Text: TextTokens{
			Primary: fg,
			Muted:   pick(wht, bwht),
			Inverse: pick(blk, canvas),
			Accent:  pick(bacc, bcya),
		},
		Border: BorderTokens{
			Normal: pick(bbk, blk),
			Subtle: pick(blk, bbk),
			Focus:  pick(bacc, bcya),
		},
		State: StateTokens{
			Info:    pick(bacc, bcya),
			Success: bgrn,
			Warning: byel,
			Danger:  bred,
		},
		Selection: SelectionTokens{
			BG: selectionBG,
			FG: selectionFG,
		},
		Input: InputTokens{
			BG:          inputBG,
			FG:          fg,
			Placeholder: pick(wht, bwht),
			Cursor:      cur,
			Border:      pick(bacc, bpur),
		},
		Bar: BarTokens{
			BG: pick(blk, canvas),
			FG: pick(fg, bwht),
		},
		Dialog: DialogTokens{
			BG:     overlay,
			FG:     fg,
			Border: pick(bacc, bcya),
			Scrim:  pick(blk, canvas),
		},
	}
}

// ensureDistinct returns candidate if luminance delta vs base >= 0.06,
// else tries fallback1, then fallback2.
func ensureDistinct(candidate, base, fallback1, fallback2 string) string {
	return ensureDistinctMin(candidate, base, 0.06, fallback1, fallback2)
}

// ensureDistinctMin is like ensureDistinct but with a configurable min delta.
func ensureDistinctMin(candidate, base string, minDelta float64, fallback1, fallback2 string) string {
	if lumDelta(candidate, base) >= minDelta {
		return candidate
	}
	if lumDelta(fallback1, base) >= minDelta {
		return fallback1
	}
	return fallback2
}

// pickDistinctFrom returns the first candidate that exceeds minDelta vs ALL
// colors in the exclusions slice. Falls back to the last candidate.
func pickDistinctFrom(candidates, exclusions []string, minDelta float64) string {
	for _, c := range candidates {
		ok := true
		for _, ex := range exclusions {
			if lumDelta(c, ex) < minDelta {
				ok = false
				break
			}
		}
		if ok {
			return c
		}
	}
	return candidates[len(candidates)-1]
}

// lumDelta returns the absolute difference in relative luminance between
// two hex color strings. Range [0, 1].
func lumDelta(a, b string) float64 {
	return math.Abs(luminance(a) - luminance(b))
}

// luminance computes the WCAG relative luminance of a hex color string.
// Input: "#rrggbb" or "#rgb". Returns [0, 1].
func luminance(hex string) float64 {
	c := tint.FromHex(hex)
	if c == nil {
		return 0
	}
	r := linearize(float64(c.R) / 255.0)
	g := linearize(float64(c.G) / 255.0)
	b := linearize(float64(c.B) / 255.0)
	return 0.2126*r + 0.7152*g + 0.0722*b
}

func linearize(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func pick(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func hex(c *tint.Color, fallback string) string {
	if c == nil {
		return fallback
	}
	return c.Hex()
}

func orFallback(a, b *tint.Color, fallback string) string {
	if a != nil {
		return a.Hex()
	}
	if b != nil {
		return b.Hex()
	}
	return fallback
}
