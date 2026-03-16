package styles

// Row paint contract:
// +-------------------------------+
// | ClipANSI(content, width)      |
// +-------------------------------+
// | Row(bg, fg, width, clipped)   |
// +-------------------------------+
// Every rendered row must own its full width with explicit background cells.
// Use Row() or RowClip() for any surface-facing row output.

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
)

type System struct {
	Theme theme.Theme
}

type SurfaceColors struct {
	BG string
	FG string
}

type CardSlabColors struct {
	HeaderBG string
	BodyBG   string
	FooterBG string
	FrameBG  string
	FrameFG  string
	ShadowBG string
	FocusBG  string
}

type Spacing struct {
	XXS int
	XS  int
	SM  int
	MD  int
	LG  int
}

var LayoutSpacing = Spacing{
	XXS: 0,
	XS:  1,
	SM:  2,
	MD:  3,
	LG:  4,
}

func New(t theme.Theme) System { return System{Theme: t} }

func (s System) BarColors() SurfaceColors {
	return SurfaceColors{BG: s.Theme.Bar.BG, FG: s.Theme.Bar.FG}
}

func (s System) StatusRowColors(role string, anchored bool) SurfaceColors {
	switch role {
	case "subheader":
		return SurfaceColors{
			BG: pick(s.Theme.Surface.Panel, s.Theme.Bar.BG),
			FG: pick(s.Theme.Text.Muted, s.Theme.Bar.FG),
		}
	case "footer":
		if anchored {
			if s.Theme.Footer.AnchoredBG != "" && s.Theme.Footer.AnchoredFG != "" {
				return SurfaceColors{
					BG: pick(s.Theme.Footer.AnchoredBG, s.Theme.Selection.BG),
					FG: pick(s.Theme.Footer.AnchoredFG, s.Theme.Selection.FG),
				}
			}
			return SurfaceColors{
				BG: pick(s.Theme.Selection.BG, s.Theme.Border.Focus),
				FG: pick(s.Theme.Selection.FG, s.Theme.Text.Inverse),
			}
		}
		return SurfaceColors{
			BG: pick(s.Theme.Surface.Panel, s.Theme.Bar.BG),
			FG: pick(s.Theme.Text.Primary, s.Theme.Bar.FG),
		}
	default:
		return s.BarColors()
	}
}

func (s System) InputColors() SurfaceColors {
	return SurfaceColors{BG: s.Theme.Input.BG, FG: s.Theme.Input.FG}
}

func (s System) PanelFrame(focused bool) lipgloss.Style {
	bg := pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas)
	if focused {
		bg = pick(s.Theme.Surface.Interactive, bg)
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

func (s System) PanelTitleBar(focused bool) lipgloss.Style {
	bg := pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas)
	fg := s.Theme.Text.Muted
	if focused {
		fg = pick(s.Theme.Text.Primary, s.Theme.Text.Muted)
	}
	return lipgloss.NewStyle().Background(lipgloss.Color(bg)).Foreground(lipgloss.Color(fg))
}

func (s System) PanelTitleBadge(focused bool) lipgloss.Style {
	bg := pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas)
	fg := pick(s.Theme.Text.Muted, s.Theme.Text.Primary)
	if focused {
		fg = pick(s.Theme.Text.Accent, s.Theme.Text.Primary)
	}
	return lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(fg))
}

func (s System) StatusBar() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Theme.Bar.FG)).
		Background(lipgloss.Color(s.Theme.Bar.BG))
}

func (s System) StatusPillMuted() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.Color(pick(s.Theme.Text.Primary, s.Theme.Bar.FG))).
		Background(lipgloss.Color(pick(s.Theme.Surface.Panel, s.Theme.Bar.BG)))
}

func (s System) FooterCardCommand(variant string, enabled bool) lipgloss.Style {
	fg, bg := s.footerCardColors(variant, enabled, true)
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg))
}

// FooterCardCommandAnchored renders command emphasis with no background so
// anchored footer rows remain a single solid strip.
func (s System) FooterCardCommandAnchored(variant string, enabled bool) lipgloss.Style {
	fg, _ := s.footerCardColors(variant, enabled, true)
	st := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(fg))
	if !enabled {
		return st.Faint(true)
	}
	return st
}

func (s System) FooterCardLabel(variant string, enabled bool) lipgloss.Style {
	fg, bg := s.footerCardColors(variant, enabled, false)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg))
}

// FooterCardLabelAnchored renders labels with no background so anchored footer
// rows remain visually continuous.
func (s System) FooterCardLabelAnchored(variant string, enabled bool) lipgloss.Style {
	fg, _ := s.footerCardColors(variant, enabled, false)
	st := lipgloss.NewStyle().Foreground(lipgloss.Color(fg))
	if !enabled {
		return st.Faint(true)
	}
	return st
}

// DialogFrame returns the base style for the dialog outer box.
// No padding is set here — padding rows are rendered explicitly by
// renderDialogFrame so every cell carries Dialog.BG via a Width() call.
func (s System) DialogFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Theme.Dialog.BG)).
		Foreground(lipgloss.Color(s.Theme.Dialog.FG))
}

func (s System) DialogHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(s.Theme.Text.Primary)).
		Background(lipgloss.Color(s.Theme.Dialog.BG))
}

func (s System) DialogEscHint() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Theme.Text.Muted)).
		Background(lipgloss.Color(s.Theme.Dialog.BG))
}

// DialogSearchRow is the style for the search input row inside a dialog.
// Uses Input.BG so it clearly stands out from the dialog body (Dialog.BG).
func (s System) DialogSearchRow() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Theme.Input.BG)).
		Foreground(lipgloss.Color(s.Theme.Input.FG))
}

// DialogListRow is the style for an unselected list row inside a dialog.
func (s System) DialogListRow() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Theme.Dialog.BG)).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

// DialogListRowSelected is the style for a selected list row.
func (s System) DialogListRowSelected() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Theme.Selection.BG)).
		Foreground(lipgloss.Color(s.Theme.Selection.FG))
}

func (s System) ElevatedCardColors(focused bool) CardSlabColors {
	focus := pick(s.Theme.Card.FocusEdgeBG, s.Theme.Border.Focus)
	if !focused {
		focus = pick(s.Theme.Card.FrameBG, s.Theme.Border.Normal)
	}
	return CardSlabColors{
		HeaderBG: pick(s.Theme.Card.HeaderBG, s.Theme.Surface.Interactive),
		BodyBG:   pick(s.Theme.Card.BodyBG, s.Theme.Surface.Panel),
		FooterBG: pick(s.Theme.Card.FooterBG, s.Theme.Surface.Panel),
		FrameBG:  pick(s.Theme.Card.FrameBG, s.Theme.Surface.Panel),
		FrameFG:  pick(s.Theme.Card.FrameFG, s.Theme.Text.Primary),
		ShadowBG: pick(s.Theme.Card.ShadowBG, s.Theme.Surface.Canvas),
		FocusBG:  focus,
	}
}

func (s System) ListItem(selected bool) lipgloss.Style {
	if selected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color(pick(s.Theme.Selection.BG, s.Theme.Border.Focus))).
			Foreground(lipgloss.Color(pick(s.Theme.Selection.FG, s.Theme.Text.Inverse)))
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas))).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

func (s System) CurrentMarker() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Text.Accent)).Bold(true)
}

func (s System) PaletteGroupHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(pick(s.Theme.Text.Accent, s.Theme.State.Info)))
}

func (s System) PaletteItem(selected bool) lipgloss.Style {
	if selected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color(pick(s.Theme.Selection.BG, s.Theme.Border.Focus))).
			Foreground(lipgloss.Color(pick(s.Theme.Selection.FG, s.Theme.Text.Inverse)))
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.Dialog.BG, s.Theme.Surface.Panel))).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

func (s System) PaletteKeybind() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Text.Muted))
}

func (s System) ActionButton(active bool) lipgloss.Style {
	if active {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(pick(s.Theme.Selection.FG, s.Theme.Text.Inverse))).
			Background(lipgloss.Color(pick(s.Theme.Selection.BG, s.Theme.Border.Focus))).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Theme.Text.Primary)).
		Background(lipgloss.Color(pick(s.Theme.Surface.Interactive, s.Theme.Surface.Panel))).
		Padding(0, 1).
		MarginRight(1)
}

// PanelBorder returns a rounded-border style colored by focus/elevation state.
func (s System) PanelBorder(focused, elevated bool) lipgloss.Style {
	var color string
	switch {
	case focused:
		color = pick(s.Theme.Border.Focus, s.Theme.Text.Accent)
	case elevated:
		color = pick(s.Theme.Border.Subtle, s.Theme.Border.Normal)
	default:
		color = pick(s.Theme.Border.Normal, s.Theme.Border.Subtle)
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color))
}

// Divider returns a style for full-width separator lines.
func (s System) Divider() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Border.Normal))
}

// SubtleDivider returns a style for low-contrast separator lines.
func (s System) SubtleDivider() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Border.Subtle))
}

// FocusAccent returns a style for the focused panel left-edge stripe.
func (s System) FocusAccent() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Theme.Border.Focus)).
		Foreground(lipgloss.Color(pick(s.Theme.Text.Inverse, s.Theme.Surface.Canvas)))
}

// ElevatedFrame returns a style for secondary/nested panels.
func (s System) ElevatedFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.Card.BodyBG, s.Theme.Surface.Panel))).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

// DashboardPanel returns a low-noise container style.
func (s System) DashboardPanel() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas))).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

// DashboardElevated returns a raised container style.
func (s System) DashboardElevated() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.Card.BodyBG, s.Theme.Surface.Panel))).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

// BorderSubtle returns a subtle border style for low-contrast separation.
func (s System) BorderSubtle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(pick(s.Theme.Border.Subtle, s.Theme.Border.Normal)))
}

// InputStyles returns a fully-themed textinput.Styles struct.
// Every sub-style carries Background(Input.BG) so that UV cell-level compositing
// sees an explicit Bg on every character — no cell falls back to canvas color.
func (s System) InputStyles() textinput.Styles {
	st := textinput.DefaultStyles(true)
	ibg := lipgloss.Color(s.Theme.Input.BG)

	st.Focused.Prompt = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Border.Focus)).Bold(true)
	st.Focused.Text = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Input.FG))
	st.Focused.Placeholder = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Input.Placeholder))
	st.Focused.Suggestion = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Text.Accent))

	st.Blurred.Prompt = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Text.Muted))
	st.Blurred.Text = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Input.FG))
	st.Blurred.Placeholder = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Input.Placeholder))
	st.Blurred.Suggestion = lipgloss.NewStyle().
		Background(ibg).
		Foreground(lipgloss.Color(s.Theme.Text.Muted))

	st.Cursor.Color = lipgloss.Color(s.Theme.Input.Cursor)
	st.Cursor.Blink = true
	return st
}

// Row returns a fully-painted row string of exactly width cells.
//
// This is the canonical way to render any row in a component or bento.
// It guarantees every cell has an explicit Bg set so the Ultraviolet surface
// overlay does not fall back to the canvas color for padding/whitespace cells.
//
// Rule: never use lipgloss.PlaceHorizontal or bare Render(content) for rows
// that sit on a surface — always go through Row() or an equivalent
// .Background().Width(w).Render() chain.
func Row(bg, fg string, width int, content string) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(fg)).
		Width(width).
		Render(content)
}

// ClipANSI truncates styled text to width cells safely.
func ClipANSI(content string, width int) string {
	if width <= 0 {
		return ""
	}
	return ansi.Truncate(content, width, "")
}

// RowClip clips ANSI content first, then paints an exact-width row.
func RowClip(bg, fg string, width int, content string) string {
	if width <= 0 {
		return ""
	}
	return Row(bg, fg, width, ClipANSI(content, width))
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func (s System) footerCardColors(variant string, enabled bool, commandPart bool) (fg string, bg string) {
	if s.Theme.Footer.AnchoredMuted != "" && variant == "muted" {
		if !enabled {
			return pick(s.Theme.Footer.AnchoredMuted, s.Theme.Text.Muted), pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas)
		}
		return pick(s.Theme.Footer.AnchoredMuted, s.Theme.Text.Muted), pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas)
	}
	if !enabled {
		return pick(s.Theme.Text.Muted, s.Theme.Text.Primary), pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas)
	}
	if !commandPart && variant != "muted" {
		return pick(s.Theme.Text.Primary, s.Theme.Text.Primary), pick(s.Theme.Surface.Interactive, s.Theme.Surface.Panel)
	}

	switch variant {
	case "primary":
		return pick(s.Theme.Selection.FG, s.Theme.Text.Inverse), pick(s.Theme.Selection.BG, s.Theme.Text.Accent)
	case "danger":
		return pick(s.Theme.Text.Inverse, s.Theme.Text.Primary), pick(s.Theme.State.Danger, s.Theme.Text.Accent)
	case "muted":
		return pick(s.Theme.Text.Muted, s.Theme.Text.Primary), pick(s.Theme.Surface.Panel, s.Theme.Surface.Canvas)
	default:
		return pick(s.Theme.Selection.FG, s.Theme.Text.Inverse), pick(s.Theme.Border.Focus, s.Theme.Text.Accent)
	}
}
