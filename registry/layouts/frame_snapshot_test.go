package layouts

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestFrameSnapshotNarrow(t *testing.T) {
	out := Frame(12, 6, Static("TOP"), Static("META"), Static("BODY"), Static("FOOT"))
	got := snapshotView(out)
	want := strings.Join([]string{
		"TOP.........",
		"META........",
		"BODY........",
		"............",
		"............",
		"FOOT........",
	}, "\n")
	if got != want {
		t.Fatalf("narrow snapshot mismatch\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}

func TestFrameSnapshotNormal(t *testing.T) {
	out := Frame(20, 8, Static("TOP"), Static("META"), Static("BODY"), Static("FOOT"))
	got := snapshotView(out)
	want := strings.Join([]string{
		"TOP.................",
		"META................",
		"BODY................",
		"....................",
		"....................",
		"....................",
		"....................",
		"FOOT................",
	}, "\n")
	if got != want {
		t.Fatalf("normal snapshot mismatch\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}

func TestFrameSnapshotWide(t *testing.T) {
	out := Frame(32, 10, Static("TOP"), Static("META"), Static("BODY"), Static("FOOT"))
	got := snapshotView(out)
	want := strings.Join([]string{
		"TOP.............................",
		"META............................",
		"BODY............................",
		"................................",
		"................................",
		"................................",
		"................................",
		"................................",
		"................................",
		"FOOT............................",
	}, "\n")
	if got != want {
		t.Fatalf("wide snapshot mismatch\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}

func TestFrameANSISafeTruncation(t *testing.T) {
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("ABCDEFGHIJKLMN")
	out := Frame(10, 5, Static("TOP"), Static("META"), Static(styled), Static("FOOT"))
	lines := strings.Split(out, "\n")
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if lipgloss.Width(line) != 10 {
			t.Fatalf("line %d width=%d, expected 10", i, lipgloss.Width(line))
		}
	}
	if ansi.Strip(lines[2]) != "ABCDEFGHIJ" {
		t.Fatalf("expected ANSI-safe truncation to ABCDEFGHIJ, got %q", ansi.Strip(lines[2]))
	}
}

func snapshotView(s string) string {
	return strings.ReplaceAll(s, " ", ".")
}
