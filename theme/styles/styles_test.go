package styles

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
)

func TestClipANSIWidth(t *testing.T) {
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("ABCDEFGHIJKLMN")
	clipped := ClipANSI(styled, 10)
	if lipgloss.Width(clipped) != 10 {
		t.Fatalf("expected clipped width 10, got %d", lipgloss.Width(clipped))
	}
	if ansi.Strip(clipped) != "ABCDEFGHIJ" {
		t.Fatalf("expected ansi-safe clip to ABCDEFGHIJ, got %q", ansi.Strip(clipped))
	}
}

func TestRowClipWidth(t *testing.T) {
	th := theme.Preset(theme.DefaultName)
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("ABCDEFGHIJKLMN")
	row := RowClip(th.Background(), th.Text(), 8, styled)
	if lipgloss.Width(row) != 8 {
		t.Fatalf("expected row width 8, got %d", lipgloss.Width(row))
	}
}

func TestRowWidth(t *testing.T) {
	th := theme.Preset(theme.DefaultName)
	row := Row(th.BackgroundPanel(), th.Text(), 20, "hello")
	if lipgloss.Width(row) != 20 {
		t.Fatalf("expected row width 20, got %d", lipgloss.Width(row))
	}
}
