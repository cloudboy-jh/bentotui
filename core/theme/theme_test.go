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

func TestPresetsDefineAllRequiredTokens(t *testing.T) {
	for _, name := range AvailableThemes() {
		theme := Preset(name)
		if err := validateTheme(theme); err != nil {
			t.Fatalf("theme %q invalid: %v", name, err)
		}
	}
}

func TestRegisterThemeRejectsMissingRequiredToken(t *testing.T) {
	theme := Preset(DefaultName)
	theme.Input.Border = ""
	if err := RegisterTheme("bad-theme", theme); err == nil {
		t.Fatal("expected register to fail when required token is missing")
	}
}

func TestBuiltinsPreserveStitchedSurfaceSeparation(t *testing.T) {
	for _, name := range AvailableThemes() {
		th := Preset(name)
		if th.Input.BG == th.Surface.Panel {
			t.Fatalf("theme %q input bg must differ from panel bg", name)
		}
		if th.Selection.BG == th.Surface.Panel {
			t.Fatalf("theme %q selection bg must differ from panel bg", name)
		}
		if th.Dialog.BG == th.Surface.Panel {
			t.Fatalf("theme %q dialog bg must differ from panel bg", name)
		}
	}
}
