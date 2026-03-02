package styles

import "testing"

import "github.com/cloudboy-jh/bentotui/core/theme"

func TestFooterCardLabelUsesSingleSurfaceForEnabledStates(t *testing.T) {
	s := New(theme.Preset(theme.DefaultName))
	_, bgNormal := s.footerCardColors("normal", true, false)
	_, bgPrimary := s.footerCardColors("primary", true, false)
	_, bgDanger := s.footerCardColors("danger", true, false)
	if bgNormal != bgPrimary || bgPrimary != bgDanger {
		t.Fatalf("expected enabled label backgrounds to match, got normal=%q primary=%q danger=%q", bgNormal, bgPrimary, bgDanger)
	}
}

func TestFooterCardMutedLabelKeepsMutedSurface(t *testing.T) {
	s := New(theme.Preset(theme.DefaultName))
	_, bg := s.footerCardColors("muted", true, false)
	if bg != s.Theme.Surface.Elevated {
		t.Fatalf("expected muted label background %q, got %q", s.Theme.Surface.Elevated, bg)
	}
}

func TestInputColorsUseDedicatedInputSurface(t *testing.T) {
	s := New(theme.Preset(theme.DefaultName))
	input := s.InputColors()
	if input.BG != s.Theme.Input.BG {
		t.Fatalf("expected input bg %q, got %q", s.Theme.Input.BG, input.BG)
	}
	if input.BG == s.Theme.Surface.Panel {
		t.Fatalf("expected input bg to differ from panel bg %q", s.Theme.Surface.Panel)
	}
}

func TestBarColorsUseDedicatedBarSurface(t *testing.T) {
	s := New(theme.Preset(theme.DefaultName))
	bar := s.BarColors()
	if bar.BG != s.Theme.Bar.BG || bar.FG != s.Theme.Bar.FG {
		t.Fatalf("expected bar colors (%q,%q), got (%q,%q)", s.Theme.Bar.BG, s.Theme.Bar.FG, bar.BG, bar.FG)
	}
}
