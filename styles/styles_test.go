package styles

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
)

func TestStatusRowColorsAnchoredUsesFooterTokens(t *testing.T) {
	th := theme.Preset(theme.DefaultName)
	th.Footer = theme.FooterTokens{
		AnchoredBG:    "#112233",
		AnchoredFG:    "#ddeeff",
		AnchoredMuted: "#99aabb",
	}

	got := New(th).StatusRowColors("footer", true)
	if got.BG != th.Footer.AnchoredBG || got.FG != th.Footer.AnchoredFG {
		t.Fatalf("expected anchored footer colors from footer tokens, got bg=%s fg=%s", got.BG, got.FG)
	}
}

func TestStatusRowColorsAnchoredFallsBackToSelection(t *testing.T) {
	th := theme.Preset(theme.DefaultName)
	th.Footer = theme.FooterTokens{}

	got := New(th).StatusRowColors("footer", true)
	if got.BG != th.Selection.BG || got.FG != th.Selection.FG {
		t.Fatalf("expected anchored fallback to selection colors, got bg=%s fg=%s", got.BG, got.FG)
	}
}

func TestClipANSIAndRowClip(t *testing.T) {
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("ABCDEFGHIJKLMN")
	clipped := ClipANSI(styled, 10)
	if lipgloss.Width(clipped) != 10 {
		t.Fatalf("expected clipped width 10, got %d", lipgloss.Width(clipped))
	}
	if ansi.Strip(clipped) != "ABCDEFGHIJ" {
		t.Fatalf("expected ansi-safe clip to ABCDEFGHIJ, got %q", ansi.Strip(clipped))
	}

	row := RowClip("#000000", "#ffffff", 8, styled)
	if lipgloss.Width(row) != 8 {
		t.Fatalf("expected row width 8, got %d", lipgloss.Width(row))
	}
}
