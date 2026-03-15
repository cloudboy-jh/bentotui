package rooms

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/rooms/internal/engine"
)

type mockCell struct {
	content string
	width   int
	height  int
}

func (m *mockCell) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *mockCell) View() tea.View {
	return tea.NewView(m.content)
}

func TestAllocateFixedFillAndRemainder(t *testing.T) {
	specs := []engine.Spec{
		{Kind: engine.Fixed, N: 1},
		{Kind: engine.Fill},
		{Kind: engine.Fill},
	}

	got := engine.Allocate(specs, 10)
	if len(got) != 3 {
		t.Fatalf("expected 3 cells, got %d", len(got))
	}

	if got[0] != 1 || got[1] != 4 || got[2] != 5 {
		t.Fatalf("unexpected allocation: %#v", got)
	}
}

func TestAllocateRatio(t *testing.T) {
	specs := []engine.Spec{
		{Kind: engine.Ratio, N: 1},
		{Kind: engine.Ratio, N: 2},
	}

	got := engine.Allocate(specs, 9)
	if got[0] != 3 || got[1] != 6 {
		t.Fatalf("unexpected ratio allocation: %#v", got)
	}
}

func TestConstrainExactDimensions(t *testing.T) {
	got := engine.Constrain("abcdef\nxy", 4, 3)
	assertExact(t, got, 4, 3)

	lines := strings.Split(got, "\n")
	if lines[0] != "abcd" {
		t.Fatalf("expected first line to be truncated to abcd, got %q", lines[0])
	}
}

func TestAllLayoutsExactDimensions(t *testing.T) {
	w, h := 37, 11
	b := Static("body")
	f := Static("footer")
	hd := Static("header")
	tb := Static("topbar")
	sb := Static("sidebar")
	mn := Static("main")
	l := Static("left")
	r := Static("right")
	n := Static("nav")
	ls := Static("list")
	d := Static("detail")
	tl := Static("tl")
	tr := Static("tr")
	bl := Static("bl")
	br := Static("br")
	pr := Static("primary")
	st := Static("strip")
	bg := Static("background")
	md := Static("modal")

	outputs := []string{
		Frame(w, h, tb, hd, b, f),
		FrameMainDrawer(w, h, 8, tb, hd, mn, sb, f),
		FrameTriple(w, h, 6, 10, tb, hd, n, ls, d, f),
		Focus(w, h, b, f),
		Pancake(w, h, hd, b, f),
		TopbarPancake(w, h, tb, hd, b, f),
		Sidebar(w, h, 9, sb, mn),
		HolyGrail(w, h, 9, hd, sb, mn, f),
		HSplit(w, h, l, r),
		VSplit(w, h, hd, b),
		HSplitFooter(w, h, l, r, f),
		TripleCol(w, h, 6, 10, n, ls, d),
		Dashboard2x2(w, h, tl, tr, bl, br),
		Dashboard2x2Footer(w, h, tl, tr, bl, br, f),
		DrawerRight(w, h, 8, mn, sb),
		DrawerChrome(w, h, 8, hd, mn, sb, f),
		Modal(w, h, 15, 5, bg, md),
		BigTopStrip(w, h, 2, pr, st),
	}

	for i, out := range outputs {
		assertExact(t, out, w, h)
		if out == "" {
			t.Fatalf("layout %d returned empty output", i)
		}
	}
}

func TestModalCentersOverlay(t *testing.T) {
	out := Modal(7, 5, 3, 1, Static(strings.Repeat(".", 7)), Static("XXX"))
	assertExact(t, out, 7, 5)

	lines := strings.Split(out, "\n")
	if !strings.Contains(lines[2], "XXX") {
		t.Fatalf("expected modal content in centered row, got %q", lines[2])
	}
}

func TestHSplitWithGutterAndDivider(t *testing.T) {
	out := HSplit(21, 4, Static("L"), Static("R"), WithGutter(1), WithDivider("normal"))
	assertExact(t, out, 21, 4)
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		if !strings.Contains(line, "|") {
			t.Fatalf("expected divider in line %d: %q", i, line)
		}
	}
}

func TestDrawerRightWithSubtleDivider(t *testing.T) {
	out := DrawerRight(24, 4, 8, Static("main"), Static("drawer"), WithGutter(1), WithDivider("subtle"))
	assertExact(t, out, 24, 4)
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		if !strings.Contains(line, ".") {
			t.Fatalf("expected subtle divider in line %d: %q", i, line)
		}
	}
}

func assertExact(t *testing.T, out string, width, height int) {
	t.Helper()

	lines := strings.Split(out, "\n")
	if len(lines) != height {
		t.Fatalf("expected %d lines, got %d", height, len(lines))
	}

	for i, line := range lines {
		if lipgloss.Width(line) != width {
			t.Fatalf("line %d width = %d, expected %d", i, lipgloss.Width(line), width)
		}
	}
}
