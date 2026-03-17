package table

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestCompactBorderlessColumnAlign(t *testing.T) {
	tb := New("NAME", "STATUS", "LATENCY")
	tb.SetCompact(true)
	tb.SetBorderless(true)
	tb.SetColumnWidth(0, 8)
	tb.SetColumnAlign(2, AlignRight)
	tb.AddRow("api", "healthy", "36ms")
	tb.SetSize(30, 4)

	out := ansi.Strip(viewString(tb.View()))
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "NAME") || !strings.Contains(lines[0], "LATENCY") {
		t.Fatalf("expected multi-column header row, got %q", lines[0])
	}
	if !strings.HasSuffix(lines[1], "36ms") {
		t.Fatalf("expected right aligned latency, got %q", lines[1])
	}
	if lipgloss.Width(lines[1]) > 30 {
		t.Fatalf("expected row width <= 30, got %d", lipgloss.Width(lines[1]))
	}
}

func TestColumnPriorityShrinkKeepsNumericAlignment(t *testing.T) {
	tb := New("SERVICE", "OWNER", "P95", "ERR%", "DEPLOY")
	tb.SetCompact(true)
	tb.SetBorderless(true)
	tb.SetColumnAlign(2, AlignRight)
	tb.SetColumnAlign(3, AlignRight)
	tb.SetColumnMinWidth(0, 8)
	tb.SetColumnMinWidth(1, 6)
	tb.SetColumnMinWidth(4, 6)
	tb.SetColumnPriority(4, 5)
	tb.SetColumnPriority(1, 4)
	tb.AddRow("checkout-api", "kai", "112ms", "1.7", "27m ago")
	tb.SetSize(34, 4)

	out := ansi.Strip(viewString(tb.View()))
	lines := strings.Split(out, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
	if lipgloss.Width(lines[1]) > 34 {
		t.Fatalf("expected row width <= 34, got %d", lipgloss.Width(lines[1]))
	}
	if !strings.Contains(lines[1], "112") || !strings.Contains(lines[1], "1.7") {
		t.Fatalf("expected numeric columns preserved under shrink, got %q", lines[1])
	}
}

func TestFocusBlurControlsKeyboardMovement(t *testing.T) {
	tb := New("NAME")
	tb.AddRow("one")
	tb.AddRow("two")
	tb.SetSize(20, 4)

	tb.Blur()
	_, _ = tb.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	out := ansi.Strip(viewString(tb.View()))
	if !strings.Contains(out, "one") {
		t.Fatalf("expected first row selected while blurred")
	}

	tb.Focus()
	_, _ = tb.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	out = ansi.Strip(viewString(tb.View()))
	if !strings.Contains(out, "two") {
		t.Fatalf("expected keyboard movement after focus")
	}
}

func TestWindowSizeMsgAppliesSize(t *testing.T) {
	tb := New("A", "B")
	_, _ = tb.Update(tea.WindowSizeMsg{Width: 50, Height: 8})
	w, h := tb.GetSize()
	if w != 50 || h != 8 {
		t.Fatalf("expected size 50x8, got %dx%d", w, h)
	}
}

func viewString(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprint(v.Content)
}
