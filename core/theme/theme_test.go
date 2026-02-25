package theme

import "testing"

func TestAvailableThemesStableOrder(t *testing.T) {
	got := AvailableThemes()
	want := []string{CatppuccinMochaName, DraculaName, OsakaJadeName}
	if len(got) != len(want) {
		t.Fatalf("expected %d themes, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected theme[%d]=%q, got %q", i, want[i], got[i])
		}
	}
}

func TestPresetFallbackUsesDefault(t *testing.T) {
	got := Preset("not-a-theme")
	want := Preset(DefaultName)
	if got != want {
		t.Fatalf("expected unknown preset to fall back to %q", DefaultName)
	}
}

func TestSetThemeRejectsUnknownName(t *testing.T) {
	if _, err := SetTheme("does-not-exist"); err == nil {
		t.Fatal("expected error for unknown theme name")
	}
}

func TestPresetsHaveSolidSurfaceTokens(t *testing.T) {
	for _, name := range AvailableThemes() {
		theme := Preset(name)
		if theme.Background == "" || theme.PanelBG == "" || theme.ElementBG == "" || theme.ModalBG == "" {
			t.Fatalf("theme %q missing layered background tokens", name)
		}
		if theme.Surface == "" || theme.SurfaceMuted == "" {
			t.Fatalf("theme %q missing background/surface tokens", name)
		}
		if theme.Border == "" || theme.BorderSubtle == "" || theme.BorderFocused == "" {
			t.Fatalf("theme %q missing border tokens", name)
		}
		if theme.SelectionBG == "" || theme.SelectionText == "" || theme.InputBG == "" || theme.InputBorder == "" {
			t.Fatalf("theme %q missing interaction tokens", name)
		}
		if theme.StatusBG == "" || theme.DialogBG == "" || theme.Scrim == "" {
			t.Fatalf("theme %q missing status/dialog/scrim tokens", name)
		}
	}
}
