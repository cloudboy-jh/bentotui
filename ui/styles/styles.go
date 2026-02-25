package styles

import (
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core/theme"
)

type System struct {
	Theme theme.Theme
}

func New(t theme.Theme) System { return System{Theme: t} }

func (s System) PanelFrame(focused bool) lipgloss.Style {
	bg := pick(s.Theme.PanelBG, s.Theme.Surface)
	if focused {
		bg = pick(s.Theme.ElementBG, bg)
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(s.Theme.Text))
}

func (s System) PanelTitleBar(focused bool) lipgloss.Style {
	bg := pick(s.Theme.ElementBG, s.Theme.SurfaceMuted)
	fg := s.Theme.Muted
	if focused {
		bg = pick(s.Theme.SelectionBG, s.Theme.BorderFocused)
		fg = pick(s.Theme.SelectionText, s.Theme.Background)
	}
	return lipgloss.NewStyle().Background(lipgloss.Color(bg)).Foreground(lipgloss.Color(fg))
}

func (s System) PanelTitleChip(focused bool) lipgloss.Style {
	bg := pick(s.Theme.TitleBG, s.Theme.Accent)
	fg := pick(s.Theme.TitleText, s.Theme.Background)
	if focused {
		bg = pick(s.Theme.BorderFocused, s.Theme.Accent)
		fg = pick(s.Theme.SelectionText, s.Theme.Background)
	}
	return lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(fg))
}

func (s System) StatusBar() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Theme.StatusText)).
		Background(lipgloss.Color(s.Theme.StatusBG))
}

func (s System) FooterActionKey(variant string, enabled bool) lipgloss.Style {
	fg, bg := s.footerActionColors(variant, enabled, true)
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg))
}

func (s System) FooterActionLabel(variant string, enabled bool) lipgloss.Style {
	fg, bg := s.footerActionColors(variant, enabled, false)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg))
}

func (s System) DialogFrame() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.ModalBG, s.Theme.DialogBG))).
		Foreground(lipgloss.Color(s.Theme.DialogText)).
		Padding(1, 2)
}

func (s System) DialogHeader() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(s.Theme.Text)).
		Background(lipgloss.Color(pick(s.Theme.ElementBG, s.Theme.SurfaceMuted)))
}

func (s System) DialogEscHint() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Muted))
}

func (s System) ListItem(selected bool) lipgloss.Style {
	if selected {
		return lipgloss.NewStyle().
			Background(lipgloss.Color(pick(s.Theme.SelectionBG, s.Theme.BorderFocused))).
			Foreground(lipgloss.Color(pick(s.Theme.SelectionText, s.Theme.Background)))
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(s.Theme.PanelBG, s.Theme.Surface))).
		Foreground(lipgloss.Color(s.Theme.Text))
}

func (s System) CurrentMarker() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(s.Theme.Accent)).Bold(true)
}

func (s System) ActionButton(active bool) lipgloss.Style {
	if active {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(pick(s.Theme.SelectionText, s.Theme.Background))).
			Background(lipgloss.Color(pick(s.Theme.SelectionBG, s.Theme.BorderFocused))).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Theme.Text)).
		Background(lipgloss.Color(pick(s.Theme.ElementBG, s.Theme.SurfaceMuted))).
		Padding(0, 1).
		MarginRight(1)
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func (s System) footerActionColors(variant string, enabled bool, keyPart bool) (fg string, bg string) {
	if !enabled {
		return pick(s.Theme.Muted, s.Theme.Text), pick(s.Theme.SurfaceMuted, s.Theme.PanelBG)
	}

	switch variant {
	case "primary":
		if keyPart {
			return pick(s.Theme.SelectionText, s.Theme.Background), pick(s.Theme.SelectionBG, s.Theme.Accent)
		}
		return pick(s.Theme.Text, s.Theme.Text), pick(s.Theme.ElementBG, s.Theme.SurfaceMuted)
	case "danger":
		if keyPart {
			return pick(s.Theme.Background, s.Theme.Text), pick(s.Theme.Error, s.Theme.Accent)
		}
		return pick(s.Theme.Text, s.Theme.Text), pick(s.Theme.ElementBG, s.Theme.SurfaceMuted)
	case "muted":
		if keyPart {
			return pick(s.Theme.Muted, s.Theme.Text), pick(s.Theme.SurfaceMuted, s.Theme.PanelBG)
		}
		return pick(s.Theme.Muted, s.Theme.Text), pick(s.Theme.SurfaceMuted, s.Theme.PanelBG)
	default:
		if keyPart {
			return pick(s.Theme.SelectionText, s.Theme.Background), pick(s.Theme.BorderFocused, s.Theme.Accent)
		}
		return pick(s.Theme.Text, s.Theme.Text), pick(s.Theme.ElementBG, s.Theme.SurfaceMuted)
	}
}
