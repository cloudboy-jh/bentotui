package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Theme is the interface every theme implements.
// Components call methods on Theme — they never touch hex strings directly.
// Pass a Theme to any brick via WithTheme(t) or SetTheme(t).
// Use theme.CurrentTheme() when you want the app-level active theme.
type Theme interface {
	// Surface
	Background() color.Color            // root canvas / terminal bg
	BackgroundPanel() color.Color       // raised component surface
	BackgroundOverlay() color.Color     // modal / dialog body
	BackgroundInteractive() color.Color // hover / focus tinted

	// Card chrome
	CardChrome() color.Color    // card header band
	CardBody() color.Color      // card content slab
	CardFrameFG() color.Color   // card title / frame text
	CardFocusEdge() color.Color // accent stripe on focused card

	// Text
	Text() color.Color        // primary body text
	TextMuted() color.Color   // secondary / dim text
	TextInverse() color.Color // text on accent / selection bg
	TextAccent() color.Color  // highlight / accent text

	// Border
	BorderNormal() color.Color
	BorderSubtle() color.Color
	BorderFocus() color.Color

	// State
	Success() color.Color
	Warning() color.Color
	Error() color.Color
	Info() color.Color

	// Selection
	SelectionBG() color.Color
	SelectionFG() color.Color

	// Input
	InputBG() color.Color
	InputFG() color.Color
	InputPlaceholder() color.Color
	InputCursor() color.Color
	InputBorder() color.Color

	// Bar
	BarBG() color.Color
	BarFG() color.Color

	// Footer (anchored)
	FooterBG() color.Color
	FooterFG() color.Color
	FooterMuted() color.Color

	// Dialog
	DialogBG() color.Color
	DialogFG() color.Color
	DialogBorder() color.Color
	DialogScrim() color.Color

	// Name
	Name() string
}

// BaseTheme provides a default implementation of Theme.
// Embed in concrete theme structs and fill the exported color fields.
type BaseTheme struct {
	ThemeName string

	BackgroundColor            color.Color
	BackgroundPanelColor       color.Color
	BackgroundOverlayColor     color.Color
	BackgroundInteractiveColor color.Color

	CardChromeColor    color.Color
	CardBodyColor      color.Color
	CardFrameFGColor   color.Color
	CardFocusEdgeColor color.Color

	TextColor        color.Color
	TextMutedColor   color.Color
	TextInverseColor color.Color
	TextAccentColor  color.Color

	BorderNormalColor color.Color
	BorderSubtleColor color.Color
	BorderFocusColor  color.Color

	SuccessColor color.Color
	WarningColor color.Color
	ErrorColor   color.Color
	InfoColor    color.Color

	SelectionBGColor color.Color
	SelectionFGColor color.Color

	InputBGColor          color.Color
	InputFGColor          color.Color
	InputPlaceholderColor color.Color
	InputCursorColor      color.Color
	InputBorderColor      color.Color

	BarBGColor color.Color
	BarFGColor color.Color

	FooterBGColor    color.Color
	FooterFGColor    color.Color
	FooterMutedColor color.Color

	DialogBGColor     color.Color
	DialogFGColor     color.Color
	DialogBorderColor color.Color
	DialogScrimColor  color.Color
}

func (t *BaseTheme) Name() string { return t.ThemeName }

func (t *BaseTheme) Background() color.Color            { return t.BackgroundColor }
func (t *BaseTheme) BackgroundPanel() color.Color       { return t.BackgroundPanelColor }
func (t *BaseTheme) BackgroundOverlay() color.Color     { return t.BackgroundOverlayColor }
func (t *BaseTheme) BackgroundInteractive() color.Color { return t.BackgroundInteractiveColor }

func (t *BaseTheme) CardChrome() color.Color    { return t.CardChromeColor }
func (t *BaseTheme) CardBody() color.Color      { return t.CardBodyColor }
func (t *BaseTheme) CardFrameFG() color.Color   { return t.CardFrameFGColor }
func (t *BaseTheme) CardFocusEdge() color.Color { return t.CardFocusEdgeColor }

func (t *BaseTheme) Text() color.Color        { return t.TextColor }
func (t *BaseTheme) TextMuted() color.Color   { return t.TextMutedColor }
func (t *BaseTheme) TextInverse() color.Color { return t.TextInverseColor }
func (t *BaseTheme) TextAccent() color.Color  { return t.TextAccentColor }

func (t *BaseTheme) BorderNormal() color.Color { return t.BorderNormalColor }
func (t *BaseTheme) BorderSubtle() color.Color { return t.BorderSubtleColor }
func (t *BaseTheme) BorderFocus() color.Color  { return t.BorderFocusColor }

func (t *BaseTheme) Success() color.Color { return t.SuccessColor }
func (t *BaseTheme) Warning() color.Color { return t.WarningColor }
func (t *BaseTheme) Error() color.Color   { return t.ErrorColor }
func (t *BaseTheme) Info() color.Color    { return t.InfoColor }

func (t *BaseTheme) SelectionBG() color.Color { return t.SelectionBGColor }
func (t *BaseTheme) SelectionFG() color.Color { return t.SelectionFGColor }

func (t *BaseTheme) InputBG() color.Color          { return t.InputBGColor }
func (t *BaseTheme) InputFG() color.Color          { return t.InputFGColor }
func (t *BaseTheme) InputPlaceholder() color.Color { return t.InputPlaceholderColor }
func (t *BaseTheme) InputCursor() color.Color      { return t.InputCursorColor }
func (t *BaseTheme) InputBorder() color.Color      { return t.InputBorderColor }

func (t *BaseTheme) BarBG() color.Color { return t.BarBGColor }
func (t *BaseTheme) BarFG() color.Color { return t.BarFGColor }

func (t *BaseTheme) FooterBG() color.Color    { return t.FooterBGColor }
func (t *BaseTheme) FooterFG() color.Color    { return t.FooterFGColor }
func (t *BaseTheme) FooterMuted() color.Color { return t.FooterMutedColor }

func (t *BaseTheme) DialogBG() color.Color     { return t.DialogBGColor }
func (t *BaseTheme) DialogFG() color.Color     { return t.DialogFGColor }
func (t *BaseTheme) DialogBorder() color.Color { return t.DialogBorderColor }
func (t *BaseTheme) DialogScrim() color.Color  { return t.DialogScrimColor }

// h converts a hex string to a color.Color via lipgloss.Color.
func h(hex string) color.Color {
	return lipgloss.Color(hex)
}

// Preset returns a named built-in theme. Falls back to CatppuccinMocha if
// the name is not found.
func Preset(name string) Theme {
	if t, ok := presets[name]; ok {
		return t
	}
	return presets[DefaultName]
}

// Names returns all built-in preset names, default first.
func Names() []string {
	out := make([]string, 0, len(presets))
	out = append(out, DefaultName)
	for name := range presets {
		if name != DefaultName {
			out = append(out, name)
		}
	}
	return out
}

const DefaultName = "catppuccin-mocha"
