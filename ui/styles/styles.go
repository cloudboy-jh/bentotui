package styles

import (
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core/theme"
)

type System struct {
	Theme theme.Theme
}

type SurfaceColors struct {
	BG string
	FG string
}

func New(t theme.Theme) System { return System{Theme: t} }

func (s System) BarColors() SurfaceColors {
	return SurfaceColors{BG: s.Theme.Bar.BG, FG: s.Theme.Bar.FG}
}

func (s System) InputColors() SurfaceColors {
	return SurfaceColors{BG: s.Theme.Input.BG, FG: s.Theme.Input.FG}
}

func (s System) PanelFrame(focused bool) lipgloss.Style {
	bg := pick(s.Theme.Surface.Panel, s.Theme.Surface.Elevated)
	if focused {
		bg = pick(s.Theme.Surface.Interactive, bg)
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

func (s System) PanelTitleBar(focused bool) lipgloss.Style {
	bg := pick(s.Theme.Surface.Interactive, s.Theme.Surface.Elevated)
	fg := s.Theme.Text.Muted
	if focused {
		bg = pick(s.Theme.Selection.BG, s.Theme.Border.Focus)
		fg = pick(s.Theme.Selection.FG, s.Theme.Text.Inverse)
	}
	return lipgloss.NewStyle().Background(lipgloss.Color(bg)).Foreground(lipgloss.Color(fg))
}

func (s System) PanelTitleBadge(focused bool) lipgloss.Style {
	bg := pick(s.Theme.Text.Accent, s.Theme.State.Info)
	fg := pick(s.Theme.Text.Inverse, s.Theme.Surface.Canvas)
	if focused {
		bg = pick(s.Theme.Border.Focus, s.Theme.Text.Accent)
		fg = pick(s.Theme.Selection.FG, s.Theme.Text.Inverse)
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

func (s System) FooterCardCommand(variant string, enabled bool) lipgloss.Style {
	fg, bg := s.footerCardColors(variant, enabled, true)
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg))
}

func (s System) FooterCardLabel(variant string, enabled bool) lipgloss.Style {
	fg, bg := s.footerCardColors(variant, enabled, false)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg))
}

func (s System) DialogFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Theme.Dialog.BG)).
		Foreground(lipgloss.Color(s.Theme.Dialog.FG)).
		Padding(1, 2)
}

func (s System) DialogHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(s.Theme.Text.Primary)).
		Background(lipgloss.Color(pick(s.Theme.Surface.Interactive, s.Theme.Surface.Elevated)))
}

func (s System) DialogEscHint() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Text.Muted))
}

func (s System) ListItem(selected bool) lipgloss.Style {
	if selected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color(pick(s.Theme.Selection.BG, s.Theme.Border.Focus))).
			Foreground(lipgloss.Color(pick(s.Theme.Selection.FG, s.Theme.Text.Inverse)))
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.Surface.Panel, s.Theme.Surface.Elevated))).
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
		Background(lipgloss.Color(pick(s.Theme.Dialog.BG, s.Theme.Surface.Elevated))).
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
		Background(lipgloss.Color(pick(s.Theme.Surface.Interactive, s.Theme.Surface.Elevated))).
		Padding(0, 1).
		MarginRight(1)
}

// ── missing token methods ─────────────────────────────────────────────────────

// Divider returns a style for full-width separator lines (─── rows).
// Uses Border.Normal as foreground so separators are visible but not loud.
func (s System) Divider() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Border.Normal))
}

// SubtleDivider returns a style for low-contrast separator lines.
// Uses Border.Subtle — for gutters and de-emphasized section breaks.
func (s System) SubtleDivider() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Border.Subtle))
}

// FocusAccent returns a style for the focused panel left-edge stripe.
// One-cell wide column with Border.Focus background.
func (s System) FocusAccent() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(s.Theme.Border.Focus)).
		Foreground(lipgloss.Color(pick(s.Theme.Text.Inverse, s.Theme.Surface.Canvas)))
}

// ElevatedFrame returns a style for secondary/nested panels.
// Uses Surface.Elevated instead of Surface.Panel — creates visual depth.
func (s System) ElevatedFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.Surface.Elevated, s.Theme.Surface.Panel))).
		Foreground(lipgloss.Color(s.Theme.Text.Primary))
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func (s System) footerCardColors(variant string, enabled bool, commandPart bool) (fg string, bg string) {
	if !enabled {
		return pick(s.Theme.Text.Muted, s.Theme.Text.Primary), pick(s.Theme.Surface.Elevated, s.Theme.Surface.Panel)
	}
	if !commandPart && variant != "muted" {
		return pick(s.Theme.Text.Primary, s.Theme.Text.Primary), pick(s.Theme.Surface.Interactive, s.Theme.Surface.Elevated)
	}

	switch variant {
	case "primary":
		return pick(s.Theme.Selection.FG, s.Theme.Text.Inverse), pick(s.Theme.Selection.BG, s.Theme.Text.Accent)
	case "danger":
		return pick(s.Theme.Text.Inverse, s.Theme.Text.Primary), pick(s.Theme.State.Danger, s.Theme.Text.Accent)
	case "muted":
		return pick(s.Theme.Text.Muted, s.Theme.Text.Primary), pick(s.Theme.Surface.Elevated, s.Theme.Surface.Panel)
	default:
		return pick(s.Theme.Selection.FG, s.Theme.Text.Inverse), pick(s.Theme.Border.Focus, s.Theme.Text.Accent)
	}
}
