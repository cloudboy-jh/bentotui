// Brick: FilePicker:
// +-----------------------------------+
// | > src/                            |
// |   go.mod                          |
// |   main.go                         |
// +-----------------------------------+
// File and directory picker backed by bubbles/filepicker.
// Copy this file into your project: bento add filepicker
package filepicker

import (
	"path/filepath"
	"strings"

	bubblesfilepicker "charm.land/bubbles/v2/filepicker"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Model struct {
	picker       bubblesfilepicker.Model
	width        int
	height       int
	focused      bool
	selectedPath string
	status       string
}

func New(startDir string) *Model {
	fp := bubblesfilepicker.New()
	fp.CurrentDirectory = cleanPath(startDir)
	fp.AutoHeight = false
	fp.SetHeight(10)
	return &Model{picker: fp, focused: true, width: 1, height: 10}
}

func (m *Model) Init() tea.Cmd { return m.picker.Init() }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil
	}
	if !m.focused {
		if _, ok := msg.(tea.KeyMsg); ok {
			return m, nil
		}
	}
	updated, cmd := m.picker.Update(msg)
	m.picker = updated
	if ok, path := m.picker.DidSelectFile(msg); ok {
		m.selectedPath = cleanPath(path)
		m.status = "selected " + filepath.Base(path)
	} else if ok, path := m.picker.DidSelectDisabledFile(msg); ok {
		m.selectedPath = ""
		m.status = "blocked " + filepath.Base(path)
	}
	return m, cmd
}

func (m *Model) View() tea.View {
	m.syncStyles()
	out := m.picker.View()
	if m.width > 0 {
		out = clipRows(out, m.width)
	}
	return tea.NewView(out)
}

func (m *Model) SetSize(width, height int) {
	m.width = max(1, width)
	m.height = max(1, height)
	m.picker.SetHeight(m.height)
}

func (m *Model) GetSize() (int, int) {
	return m.width, max(1, m.picker.Height())
}

func (m *Model) Focus()          { m.focused = true }
func (m *Model) Blur()           { m.focused = false }
func (m *Model) IsFocused() bool { return m.focused }

func (m *Model) SetDirectory(path string) {
	if strings.TrimSpace(path) == "" {
		return
	}
	m.picker.CurrentDirectory = cleanPath(path)
}

func (m *Model) SetAllowedTypes(exts ...string) {
	m.picker.AllowedTypes = append([]string(nil), exts...)
}

func (m *Model) SetShowHidden(v bool)       { m.picker.ShowHidden = v }
func (m *Model) SetAllowDirectories(v bool) { m.picker.DirAllowed = v }
func (m *Model) SetAllowFiles(v bool)       { m.picker.FileAllowed = v }

func (m *Model) CurrentDirectory() string { return m.picker.CurrentDirectory }
func (m *Model) HighlightedPath() string  { return m.picker.HighlightedPath() }
func (m *Model) SelectedPath() string     { return m.selectedPath }
func (m *Model) Status() string           { return m.status }

func (m *Model) syncStyles() {
	t := theme.CurrentTheme()
	s := m.picker.Styles
	s.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Accent, t.Border.Focus))).Bold(true)
	s.DisabledCursor = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Muted, t.Border.Subtle)))
	s.Directory = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Accent, t.Text.Primary)))
	s.File = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Primary, t.Text.Primary)))
	s.DisabledFile = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Muted, t.Text.Muted)))
	s.Permission = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Muted, t.Border.Normal)))
	s.FileSize = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Muted, t.Text.Muted))).Align(lipgloss.Right)
	s.Symlink = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Accent, t.Text.Primary)))
	s.Selected = lipgloss.NewStyle().
		Background(lipgloss.Color(pick(t.Selection.BG, t.Border.Focus))).
		Foreground(lipgloss.Color(pick(t.Selection.FG, t.Text.Inverse))).
		Bold(true)
	s.DisabledSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Border.Subtle, t.Text.Muted)))
	s.EmptyDirectory = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Muted, t.Text.Muted))).SetString("No files")
	m.picker.Styles = s
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func cleanPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "."
	}
	return filepath.Clean(trimmed)
}

func clipRows(content string, width int) string {
	if width <= 0 {
		return ""
	}
	lines := strings.Split(content, "\n")
	clipped := make([]string, 0, len(lines))
	for _, line := range lines {
		clipped = append(clipped, ansi.Truncate(line, width, ""))
	}
	return strings.Join(clipped, "\n")
}
