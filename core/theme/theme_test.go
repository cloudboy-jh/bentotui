package theme

import "testing"

func TestAvailableThemesStableOrder(t *testing.T) {
	got := AvailableThemes()

	// Must have exactly as many entries as registered builtins.
	if len(got) != len(builtinThemes) {
		t.Fatalf("expected %d themes, got %d", len(builtinThemes), len(got))
	}

	// Default theme must always be first.
	if len(got) == 0 || got[0] != DefaultName {
		t.Fatalf("expected first theme to be %q, got %q", DefaultName, got[0])
	}

	// Remaining entries (index 1 onward) must be sorted ascending.
	for i := 1; i < len(got)-1; i++ {
		if got[i] > got[i+1] {
			t.Fatalf("themes not sorted at index %d/%d: %q > %q", i, i+1, got[i], got[i+1])
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
