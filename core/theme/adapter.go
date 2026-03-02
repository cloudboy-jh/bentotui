package theme

import tint "github.com/lrstanley/bubbletint/v2"

// fromTint converts a bubbletint.Tint (the 16-color ANSI palette from
// iTerm2-Color-Schemes) into a BentoTUI Theme by mapping each semantic
// token slot to the appropriate ANSI color role.
//
// Mapping rationale:
//   - Surface.Canvas   = Black        (darkest slot — outermost shell bg)
//   - Surface.Panel    = Bg           (terminal background — default panel bg)
//   - Surface.Elevated = BrightBlack  (one step lighter — secondary panels)
//   - Surface.Overlay  = Black        (modal backdrop — max darkness)
//   - Surface.Interactive = BrightBlack (same as Elevated; focus stripe signals active)
//   - Text.Primary     = Fg           (terminal foreground — always readable on Bg)
//   - Text.Muted       = White        (dim white — softer than BrightWhite)
//   - Text.Inverse     = Black        (dark text on accent/selection backgrounds)
//   - Text.Accent      = BrightBlue   (hero color in the vast majority of dark tints)
//   - Border.Normal    = BrightBlack  (visible but not loud against Panel)
//   - Border.Subtle    = Black        (barely visible — same depth as canvas)
//   - Border.Focus     = BrightBlue   (matches accent)
//   - State.*          = Bright ANSI semantics (info=blue, success=green, warn=yellow, danger=red)
//   - Selection.*      = SelectionBg tint field (falls back to BrightBlue) + Black fg
//   - Input.*          = BrightBlack bg, Fg text, BrightCyan border
//   - Bar.*            = Black bg (recedes), Fg text
//   - Dialog.*         = BrightBlack bg (lifted above Panel), BrightBlue border
func fromTint(t *tint.Tint, name string) Theme {
	if t == nil {
		// Defensive: return the zero value; validateTheme will catch it.
		return Theme{Name: name}
	}
	return Theme{
		Name: name,
		Surface: SurfaceTokens{
			Canvas:      hex(t.Black, "#000000"),
			Panel:       hex(t.Bg, "#1a1a2e"),
			Elevated:    hex(t.BrightBlack, "#2a2a3e"),
			Overlay:     hex(t.Black, "#000000"),
			Interactive: hex(t.BrightBlack, "#2a2a3e"),
		},
		Text: TextTokens{
			Primary: hex(t.Fg, "#e0e0e0"),
			Muted:   hex(t.White, "#a0a0b0"),
			Inverse: hex(t.Black, "#000000"),
			Accent:  hex(t.BrightBlue, "#89b4fa"),
		},
		Border: BorderTokens{
			Normal: hex(t.BrightBlack, "#444444"),
			Subtle: hex(t.Black, "#222222"),
			Focus:  hex(t.BrightBlue, "#89b4fa"),
		},
		State: StateTokens{
			Info:    hex(t.BrightBlue, "#89b4fa"),
			Success: hex(t.BrightGreen, "#a6e3a1"),
			Warning: hex(t.BrightYellow, "#f9e2af"),
			Danger:  hex(t.BrightRed, "#f38ba8"),
		},
		Selection: SelectionTokens{
			BG: orFallback(t.SelectionBg, t.BrightBlue, "#89b4fa"),
			FG: hex(t.Black, "#000000"),
		},
		Input: InputTokens{
			BG:          hex(t.BrightBlack, "#2a2a3e"),
			FG:          hex(t.Fg, "#e0e0e0"),
			Placeholder: hex(t.White, "#a0a0b0"),
			Cursor:      orFallback(t.Cursor, t.BrightBlue, "#89b4fa"),
			Border:      hex(t.BrightCyan, "#74c7ec"),
		},
		Bar: BarTokens{
			BG: hex(t.Black, "#000000"),
			FG: hex(t.Fg, "#e0e0e0"),
		},
		Dialog: DialogTokens{
			BG:     hex(t.BrightBlack, "#2a2a3e"),
			FG:     hex(t.Fg, "#e0e0e0"),
			Border: hex(t.BrightBlue, "#89b4fa"),
			Scrim:  hex(t.Black, "#000000"),
		},
	}
}

// hex returns the hex string of a *tint.Color, or fallback if nil.
func hex(c *tint.Color, fallback string) string {
	if c == nil {
		return fallback
	}
	return c.Hex()
}

// orFallback returns the hex of a if non-nil, else the hex of b if non-nil,
// else the hardcoded fallback string.
func orFallback(a, b *tint.Color, fallback string) string {
	if a != nil {
		return a.Hex()
	}
	if b != nil {
		return b.Hex()
	}
	return fallback
}
