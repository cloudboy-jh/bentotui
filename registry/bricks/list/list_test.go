package list

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func TestDefaultFormatterSectionAndStat(t *testing.T) {
	l := New(20)
	l.SetSize(24, 5)
	l.AppendSection("services")
	l.AppendRow(Row{Label: "api", Status: "ok", Stat: "36ms"})
	l.AppendRow(Row{Label: "workers", Status: "warn", Stat: "112ms"})

	out := ansi.Strip(viewString(l.View()))
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "  SERVICES") {
		t.Fatalf("expected section header row, got %q", lines[0])
	}
	if !strings.HasSuffix(lines[1], "36ms") {
		t.Fatalf("expected right stat in row, got %q", lines[1])
	}
	if lipgloss.Width(lines[1]) != 24 {
		t.Fatalf("expected row width 24, got %d", lipgloss.Width(lines[1]))
	}
}

func TestCustomFormatterUsed(t *testing.T) {
	l := New(10)
	l.SetFormatter(func(row Row, selected bool, width int) string {
		if row.Kind == RowSection {
			return "SECTION"
		}
		return "ROW"
	})
	l.AppendSection("x")
	l.Append("item")

	out := viewString(l.View())
	if !strings.Contains(out, "SECTION") || !strings.Contains(out, "ROW") {
		t.Fatalf("expected custom formatter output, got %q", out)
	}
}

func TestStructuredRowFieldsFallbackAndAlignment(t *testing.T) {
	l := New(20)
	l.SetSize(36, 4)
	l.AppendRow(Row{Primary: "api", Secondary: "health", Tone: ToneSuccess, RightStat: "36ms", SelectedStyle: SelectedSubtle})
	l.SetCursor(0)

	out := ansi.Strip(viewString(l.View()))
	lines := strings.Split(out, "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "[success] api - health") {
		t.Fatalf("expected tone/primary/secondary text, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[0], "* ") {
		t.Fatalf("expected subtle selected marker, got %q", lines[0])
	}
	if !strings.HasSuffix(lines[0], "36ms") {
		t.Fatalf("expected right stat alignment, got %q", lines[0])
	}
	if lipgloss.Width(lines[0]) != 36 {
		t.Fatalf("expected row width 36, got %d", lipgloss.Width(lines[0]))
	}
}

func TestCompactDensityDropsSecondaryButKeepsStat(t *testing.T) {
	l := New(20)
	l.SetSize(28, 3)
	l.SetDensity(DensityCompact)
	l.AppendRow(Row{Primary: "customer-sync", Secondary: "slow query", Tone: ToneWarn, RightStat: "27m"})
	out := ansi.Strip(viewString(l.View()))
	if strings.Contains(out, "slow query") {
		t.Fatalf("expected compact density to drop secondary text, got %q", out)
	}
	if !strings.Contains(out, "[warn]") || !strings.HasSuffix(strings.TrimSpace(out), "27m") {
		t.Fatalf("expected tone and right stat in compact row, got %q", out)
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
