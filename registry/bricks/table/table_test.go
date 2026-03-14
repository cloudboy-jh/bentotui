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
	if strings.Contains(lines[0], "|") {
		t.Fatalf("expected borderless header, got %q", lines[0])
	}
	if !strings.HasSuffix(lines[1], "36ms") {
		t.Fatalf("expected right aligned latency, got %q", lines[1])
	}
	if lipgloss.Width(lines[1]) != 30 {
		t.Fatalf("expected width 30, got %d", lipgloss.Width(lines[1]))
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
