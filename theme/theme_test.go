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

	prevTier := ThemeTierStable
	for i := 1; i < len(got); i++ {
		meta, ok := ThemeMetadata(got[i])
		if !ok {
			t.Fatalf("missing theme metadata for %q", got[i])
		}
		if prevTier == ThemeTierExperimental && meta.Tier == ThemeTierStable {
			t.Fatalf("stable theme %q appears after experimental themes", got[i])
		}
		if i < len(got)-1 {
			nextMeta, _ := ThemeMetadata(got[i+1])
			if meta.Tier == nextMeta.Tier && got[i] > got[i+1] {
				t.Fatalf("themes not sorted within tier at %d/%d: %q > %q", i, i+1, got[i], got[i+1])
			}
		}
		prevTier = meta.Tier
	}
}

func TestAvailableStableThemesOnlyReturnsStable(t *testing.T) {
	names := AvailableStableThemes()
	if len(names) == 0 {
		t.Fatal("expected at least one stable theme")
	}
	for _, name := range names {
		meta, ok := ThemeMetadata(name)
		if !ok {
			t.Fatalf("missing metadata for stable theme %q", name)
		}
		if meta.Tier != ThemeTierStable {
			t.Fatalf("expected stable tier for %q, got %q", name, meta.Tier)
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

// TestBuiltinsLayerContrast verifies that key token pairs in every built-in
// theme meet the minimum luminance-delta thresholds defined in validateTheme.
// This replaces the old string-equality check — two tokens can share a hex
// value only if the adapter deliberately maps them to the same palette slot,
// which validateTheme will catch via lumDelta.
func TestBuiltinsLayerContrast(t *testing.T) {
	for _, name := range AvailableThemes() {
		th := Preset(name)
		pairs := []struct {
			a, b     string
			la, lb   string
			minDelta float64
		}{
			{th.Surface.Panel, th.Surface.Canvas, "surface.panel", "surface.canvas", minSurfacePanelCanvasDelta},
			{th.Surface.Interactive, th.Surface.Panel, "surface.interactive", "surface.panel", minSurfaceInteractivePanelDelta},
			{th.Input.BG, th.Surface.Canvas, "input.bg", "surface.canvas", 0.03},
			{th.Selection.BG, th.Surface.Canvas, "selection.bg", "surface.canvas", 0.05},
			{th.Selection.BG, th.Input.BG, "selection.bg", "input.bg", 0.05},
			{th.Dialog.BG, th.Surface.Canvas, "dialog.bg", "surface.canvas", 0.03},
			{th.Bar.BG, th.Surface.Canvas, "bar.bg", "surface.canvas", 0.02},
			{th.Card.HeaderBG, th.Card.BodyBG, "card.headerBG", "card.bodyBG", minCardHeaderBodyDelta},
			{th.Card.FrameBG, th.Card.BodyBG, "card.frameBG", "card.bodyBG", minCardFrameBodyDelta},
			{th.Card.ShadowBG, th.Surface.Canvas, "card.shadowBG", "surface.canvas", minCardShadowCanvasDelta},
			{th.Card.FocusEdgeBG, th.Card.FrameBG, "card.focusEdgeBG", "card.frameBG", minCardFocusEdgeFrameDelta},
		}
		for _, p := range pairs {
			delta := lumDelta(p.a, p.b)
			if delta < p.minDelta {
				t.Errorf("theme %q: %s vs %s luminance delta %.3f < %.3f",
					name, p.la, p.lb, delta, p.minDelta)
			}
		}
	}
}

func TestValidateThemeRejectsPartialFooterTokens(t *testing.T) {
	th := Preset(DefaultName)
	th.Footer = FooterTokens{}
	th.Footer.AnchoredBG = th.Selection.BG
	if err := validateTheme(th); err == nil {
		t.Fatal("expected validation to fail when footer anchored tokens are partially defined")
	}
}

func TestValidateThemeAllowsFooterTokensUnset(t *testing.T) {
	th := Preset(DefaultName)
	th.Footer = FooterTokens{}
	if err := validateTheme(th); err != nil {
		t.Fatalf("expected validation to allow missing footer tokens, got: %v", err)
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
