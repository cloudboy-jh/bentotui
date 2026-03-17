package theme

import (
	"fmt"
	"math"

	tint "github.com/lrstanley/bubbletint/v2"
)

// fromTint maps a bubbletint palette to BentoTUI semantic tokens.
//
// Layer hierarchy (darkest → lightest for dark themes):
//
//	Canvas      — terminal root background (filled by surface.Fill)
//	Panel       — default component body (input block, panels)
//	Overlay     — modal/dialog body
//	Interactive — hover/focus tinted surfaces
//	Card.*      — elevated-card slab tones (chrome/body)
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
	// Overlay (dialog) uses BrightBlack as base, guaranteed distinct from panel.
	// Interactive is a lighter tint of the accent for hover surfaces.
	canvas := bg
	panel := pick(bbk, blk)
	overlay := ensureDelta(pick(bbk, blk), canvas, 0.03)
	interactive := pick(bbk, bacc)
	panel, interactive = normalizeSurfaces(canvas, panel, interactive)
	muted := ensureDelta(pick(wht, bwht), fg, minTextPrimaryMutedDelta)

	// Input BG must contrast against canvas but remain a dark surface color —
	// prefer BrightBlack (raised surface), fall back to Black if BrightBlack
	// is too close to canvas. Never fall back to the accent (bacc) — that
	// would collide with selectionBG.
	inputBG := ensureDelta(pick(bbk, blk), canvas, 0.03)

	// Selection must always be the brightest available accent — clearly
	// distinguishable from both canvas and inputBG.
	// Prefer BrightBlue → BrightCyan → BrightPurple → BrightYellow.
	// Never use surface slots (BrightBlack/Black) which can collide with inputBG.
	selectionBG := pickDistinctFrom([]string{bacc, bcya, bpur, byel}, []string{canvas, inputBG}, 0.05)
	selectionFG := blk // dark text on bright selection
	barBG := ensureDistinctMin(panel, canvas, 0.02, interactive, blk)
	footerBG := ensureDelta(pick(selectionBG, barBG), panel, 0.03)
	footerFG := ensureDelta(pick(fg, bwht), footerBG, minFooterFGToBGDelta)
	footerMuted := ensureDelta(muted, footerBG, minFooterMutedToBGDelta)
	if lumDelta(footerFG, footerMuted) < minFooterFGMutedDelta {
		footerMuted = ensureDelta(blendHex(footerMuted, footerBG, 0.35), footerFG, minFooterFGMutedDelta)
		footerMuted = ensureDelta(footerMuted, footerBG, minFooterMutedToBGDelta)
	}

	cardBody := ensureDelta(pick(blk, panel), panel, 0.04)
	cardChrome := ensureDelta(pick(panel, interactive), cardBody, minCardChromeBodyDelta)
	cardFrameFG := ensureDelta(pick(bwht, fg), cardChrome, 0.10)
	cardFocusEdge := ensureDelta(pick(bacc, bcya), cardChrome, minCardFocusEdgeChromeDelta)

	// ── assemble ──────────────────────────────────────────────────────────────
	return Theme{
		Name: name,
		Surface: SurfaceTokens{
			Canvas:      canvas,
			Panel:       panel,
			Overlay:     overlay,
			Interactive: interactive,
		},
		Text: TextTokens{
			Primary: fg,
			Muted:   muted,
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
			Placeholder: muted,
			Cursor:      cur,
			Border:      pick(bacc, bpur),
		},
		Bar: BarTokens{
			BG: pick(barBG, panel),
			FG: pick(fg, bwht),
		},
		Footer: FooterTokens{
			AnchoredBG:    footerBG,
			AnchoredFG:    footerFG,
			AnchoredMuted: footerMuted,
		},
		Dialog: DialogTokens{
			BG:     overlay,
			FG:     fg,
			Border: pick(bacc, bcya),
			Scrim:  pick(blk, canvas),
		},
		Card: CardTokens{
			ChromeBG:    cardChrome,
			BodyBG:      cardBody,
			FrameFG:     cardFrameFG,
			FocusEdgeBG: cardFocusEdge,
		},
	}
}

func normalizeSurfaces(canvas, panel, interactive string) (string, string) {
	for i := 0; i < 4; i++ {
		panel = ensureDelta(panel, canvas, minSurfacePanelCanvasDelta)
		interactive = ensureDelta(interactive, panel, minSurfaceInteractivePanelDelta)
		interactive = ensureDelta(interactive, canvas, 0.05)
	}
	return panel, interactive
}

func ensureDelta(candidate, base string, minDelta float64) string {
	if lumDelta(candidate, base) >= minDelta {
		return candidate
	}
	best := candidate
	bestDiff := 2.0
	for _, target := range []string{"#000000", "#ffffff"} {
		for step := 1; step <= 20; step++ {
			alpha := float64(step) / 20.0
			next := blendHex(candidate, target, alpha)
			delta := lumDelta(next, base)
			if delta >= minDelta {
				if alpha < bestDiff {
					best = next
					bestDiff = alpha
				}
				break
			}
		}
	}
	return best
}

func blendHex(from, to string, alpha float64) string {
	a := tint.FromHex(from)
	b := tint.FromHex(to)
	if a == nil || b == nil {
		return from
	}
	if alpha < 0 {
		alpha = 0
	}
	if alpha > 1 {
		alpha = 1
	}
	r := int(math.Round(float64(a.R) + (float64(b.R)-float64(a.R))*alpha))
	g := int(math.Round(float64(a.G) + (float64(b.G)-float64(a.G))*alpha))
	bv := int(math.Round(float64(a.B) + (float64(b.B)-float64(a.B))*alpha))
	return fmt.Sprintf("#%02x%02x%02x", clamp8(r), clamp8(g), clamp8(bv))
}

func clamp8(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
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
