package theme

import "testing"

func TestAvailableThemesStableOrder(t *testing.T) {
	got := AvailableThemes()

	if len(got) != len(builtinThemes) {
		t.Fatalf("expected %d themes, got %d", len(builtinThemes), len(got))
	}

	if len(got) == 0 || got[0] != DefaultName {
		t.Fatalf("expected first theme to be %q, got %q", DefaultName, got[0])
	}

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
		th := Preset(name)
		if err := validateTheme(th); err != nil {
			t.Fatalf("theme %q invalid: %v", name, err)
		}
	}
}

func TestRegisterThemeRejectsMissingRequiredToken(t *testing.T) {
	th := Preset(DefaultName)
	th.Input.Border = ""
	if err := RegisterTheme("bad-theme", th); err == nil {
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

func TestCurrentThemeConcurrentAccess(t *testing.T) {
	// Smoke test: concurrent reads and one write must not race.
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			_ = CurrentTheme()
			done <- struct{}{}
		}()
	}
	go func() {
		_, _ = PreviewTheme(DefaultName)
		done <- struct{}{}
	}()
	for i := 0; i < 11; i++ {
		<-done
	}
}
