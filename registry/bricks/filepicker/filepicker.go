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
	if strings.TrimSpace(startDir) != "" {
		fp.CurrentDirectory = startDir
	}
	fp.AutoHeight = false
	fp.SetHeight(10)
	return &Model{picker: fp, focused: true}
}

func (m *Model) Init() tea.Cmd { return m.picker.Init() }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}
	updated, cmd := m.picker.Update(msg)
	m.picker = updated
	if ok, path := m.picker.DidSelectFile(msg); ok {
		m.selectedPath = path
		m.status = "selected " + filepath.Base(path)
	} else if ok, path := m.picker.DidSelectDisabledFile(msg); ok {
		m.status = "blocked " + filepath.Base(path)
	}
	return m, cmd
}

func (m *Model) View() tea.View {
	m.syncStyles()
	return tea.NewView(m.picker.View())
}

func (m *Model) SetSize(width, height int) {
	if width > 0 {
		m.width = width
	}
	if height > 0 {
		m.height = height
	}
	if m.height > 0 {
		m.picker.SetHeight(m.height)
	}
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
	m.picker.CurrentDirectory = path
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
	s.Selected = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Selection.BG, t.Text.Accent))).Bold(true)
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
