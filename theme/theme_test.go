package theme_test

import (
	"testing"

	"github.com/cloudboy-jh/bentotui/theme"
)

func TestPresetReturnsTheme(t *testing.T) {
	for _, name := range theme.Names() {
		th := theme.Preset(name)
		if th == nil {
			t.Errorf("Preset(%q) returned nil", name)
			continue
		}
		if th.Name() != name {
			t.Errorf("Preset(%q).Name() = %q, want %q", name, th.Name(), name)
		}
	}
}

func TestPresetFallback(t *testing.T) {
	th := theme.Preset("nonexistent-theme-xyz")
	if th == nil {
		t.Fatal("Preset fallback returned nil")
	}
	if th.Name() != theme.DefaultName {
		t.Errorf("Preset fallback Name() = %q, want %q", th.Name(), theme.DefaultName)
	}
}

func TestNamesContainsDefault(t *testing.T) {
	names := theme.Names()
	if len(names) == 0 {
		t.Fatal("Names() returned empty slice")
	}
	if names[0] != theme.DefaultName {
		t.Errorf("Names()[0] = %q, want default %q", names[0], theme.DefaultName)
	}
}

func TestAllPresetsHaveNonNilColors(t *testing.T) {
	for _, name := range theme.Names() {
		th := theme.Preset(name)
		checks := []struct {
			method string
			val    interface {
				RGBA() (uint32, uint32, uint32, uint32)
			}
		}{
			{"Background", th.Background()},
			{"BackgroundPanel", th.BackgroundPanel()},
			{"Text", th.Text()},
			{"TextMuted", th.TextMuted()},
			{"BorderFocus", th.BorderFocus()},
			{"SelectionBG", th.SelectionBG()},
			{"InputBG", th.InputBG()},
			{"DialogBG", th.DialogBG()},
		}
		for _, c := range checks {
			if c.val == nil {
				t.Errorf("Preset(%q).%s() returned nil", name, c.method)
			}
		}
	}
}

func TestManagerCurrentTheme(t *testing.T) {
	th := theme.CurrentTheme()
	if th == nil {
		t.Fatal("CurrentTheme() returned nil")
	}
}

func TestManagerSetTheme(t *testing.T) {
	original := theme.CurrentThemeName()

	_, err := theme.SetTheme("dracula")
	if err != nil {
		t.Fatalf("SetTheme(dracula) error: %v", err)
	}
	if theme.CurrentThemeName() != "dracula" {
		t.Errorf("after SetTheme(dracula), CurrentThemeName() = %q", theme.CurrentThemeName())
	}

	_, _ = theme.SetTheme(original)
}

func TestManagerSetThemeUnknown(t *testing.T) {
	_, err := theme.SetTheme("this-theme-does-not-exist")
	if err == nil {
		t.Error("SetTheme with unknown name should return an error")
	}
}

func TestAvailableThemes(t *testing.T) {
	names := theme.AvailableThemes()
	if len(names) == 0 {
		t.Fatal("AvailableThemes() returned empty slice")
	}
	if names[0] != theme.DefaultName {
		t.Errorf("AvailableThemes()[0] = %q, want %q", names[0], theme.DefaultName)
	}
}

func TestRegisterCustomTheme(t *testing.T) {
	custom := theme.Preset(theme.DefaultName)
	err := theme.RegisterTheme("test-custom", custom)
	if err != nil {
		t.Fatalf("RegisterTheme error: %v", err)
	}
	got := theme.Preset("test-custom")
	if got == nil {
		t.Fatal("Preset(test-custom) returned nil after registration")
	}
}
