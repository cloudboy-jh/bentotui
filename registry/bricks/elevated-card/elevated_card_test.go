package elevatedcard

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type mockContent struct{ text string }

func (m *mockContent) Init() tea.Cmd                           { return nil }
func (m *mockContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *mockContent) View() tea.View                          { return tea.NewView(m.text) }
func (m *mockContent) SetSize(width, height int)               {}

func TestElevatedCardWidthExact(t *testing.T) {
	c := New(Title("Section"), Content(&mockContent{text: "line one\nline two"}))
	c.SetSize(44, 8)
	out := viewString(c.View())
	lines := strings.Split(out, "\n")
	if len(lines) != 8 {
		t.Fatalf("expected 8 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if got := lipgloss.Width(line); got != 44 {
			t.Fatalf("line %d width mismatch: got %d want 44", i, got)
		}
	}
}

func TestElevatedCardMetaAndFooterRender(t *testing.T) {
	c := New(
		Title("Section"),
		Meta("meta line"),
		Footer("footer line"),
		Content(&mockContent{text: "line one\nline two"}),
	)
	c.SetSize(48, 10)
	out := viewString(c.View())
	plain := ansi.Strip(out)
	if !strings.Contains(plain, "meta line") {
		t.Fatalf("expected meta line in output")
	}
	if !strings.Contains(plain, "footer line") {
		t.Fatalf("expected footer line in output")
	}
}
