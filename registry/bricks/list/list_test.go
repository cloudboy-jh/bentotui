package list

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

func nonEmptyLines(s string) []string {
	parts := strings.Split(s, "\n")
	lines := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) == "" {
			continue
		}
		lines = append(lines, p)
	}
	return lines
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

func TestSectionAndStatRendered(t *testing.T) {
	l := New(20)
	l.SetSize(24, 5)
	l.AppendSection("services")
	l.AppendRow(Row{Label: "api", Status: "ok", Stat: "36ms"})
	l.AppendRow(Row{Label: "workers", Status: "warn", Stat: "112ms"})

	out := ansi.Strip(viewString(l.View()))
	lines := nonEmptyLines(out)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d:\n%s", len(lines), out)
	}
	if !strings.Contains(lines[0], "SERVICES") {
		t.Fatalf("expected section header row, got %q", lines[0])
	}
	if !strings.Contains(lines[1], "36ms") {
		t.Fatalf("expected right stat in row, got %q", lines[1])
	}
}

func TestCustomFormatterUsed(t *testing.T) {
	l := New(10)
	l.SetSize(16, 4)
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

func TestStructuredRowFieldsRendered(t *testing.T) {
	l := New(20)
	l.SetSize(36, 4)
	l.AppendRow(Row{Primary: "api", Secondary: "health", Tone: ToneSuccess, RightStat: "36ms", SelectedStyle: SelectedSubtle})
	l.SetCursor(0)

	out := ansi.Strip(viewString(l.View()))
	lines := nonEmptyLines(out)
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d:\n%s", len(lines), out)
	}
	if !strings.Contains(lines[0], "api") {
		t.Fatalf("expected primary text in row, got %q", lines[0])
	}
	if !strings.Contains(lines[0], "36ms") {
		t.Fatalf("expected right stat in row, got %q", lines[0])
	}
}

func TestCompactDensityDropsSecondary(t *testing.T) {
	l := New(20)
	l.SetSize(28, 3)
	l.SetDensity(DensityCompact)
	l.AppendRow(Row{Primary: "customer-sync", Secondary: "slow query", Tone: ToneWarn, RightStat: "27m"})
	out := ansi.Strip(viewString(l.View()))
	if strings.Contains(out, "slow query") {
		t.Fatalf("expected compact density to drop secondary text, got %q", out)
	}
	if !strings.Contains(out, "[warn]") {
		t.Fatalf("expected tone in compact row, got %q", out)
	}
}

func TestKeyboardNavigationPersistsAcrossView(t *testing.T) {
	l := New(20)
	l.SetSize(24, 4)
	l.Append("one")
	l.Append("two")
	l.Append("three")

	_, _ = l.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	out := ansi.Strip(viewString(l.View()))
	lines := nonEmptyLines(out)
	found := false
	for _, line := range lines {
		if strings.Contains(line, "two") && strings.Contains(line, "> ") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected 'two' selected after down key, got:\n%s", out)
	}
}

func TestBlurBlocksNavigation(t *testing.T) {
	l := New(20)
	l.SetSize(24, 4)
	l.Append("one")
	l.Append("two")

	l.Blur()
	_, _ = l.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	out := ansi.Strip(viewString(l.View()))
	if strings.Contains(out, "> two") {
		t.Fatalf("expected cursor to stay on 'one' while blurred, got %q", out)
	}
}

func TestWindowSizeMsgAppliesSize(t *testing.T) {
	l := New(20)
	_, _ = l.Update(tea.WindowSizeMsg{Width: 40, Height: 7})
	w, h := l.GetSize()
	if w != 40 || h != 7 {
		t.Fatalf("expected size 40x7, got %dx%d", w, h)
	}
}
