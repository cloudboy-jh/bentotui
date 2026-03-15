package panel

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type mockContent struct {
	text string
}

func (m *mockContent) Init() tea.Cmd                           { return nil }
func (m *mockContent) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m *mockContent) View() tea.View                          { return tea.NewView(m.text) }
func (m *mockContent) SetSize(width, height int)               {}

func TestPanelRowsRemainWidthExactWithANSIContent(t *testing.T) {
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("colored-segment")
	m := New(Title("Diagnostics"), Content(&mockContent{text: "prefix " + styled + " suffix"}))
	m.SetSize(48, 8)

	out := viewString(m.View())
	lines := strings.Split(out, "\n")
	if len(lines) != 8 {
		t.Fatalf("expected 8 lines, got %d", len(lines))
	}

	for i, line := range lines {
		if got := lipgloss.Width(line); got != 48 {
			t.Fatalf("line %d width mismatch: got %d want 48", i, got)
		}
	}
}

func TestPanelStripsInputANSISequencesFromContent(t *testing.T) {
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("ansi-input")
	m := New(Content(&mockContent{text: "x " + styled + " y"}))
	m.SetSize(36, 4)

	out := viewString(m.View())
	if strings.Contains(out, "38;5;196") {
		t.Fatalf("expected input ansi sequence to be stripped from content rows")
	}
}
