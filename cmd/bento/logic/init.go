// Package logic provides business logic for the bento CLI.
// These functions are UI-agnostic and can be used by both CLI and TUI modes.
package logic

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// ProjectConfig holds the configuration for a new project.
type ProjectConfig struct {
	AppName string
	Module  string
}

// ScaffoldProject creates a new BentoTUI project with the given configuration.
// Returns a slice of created file paths and any error encountered.
func ScaffoldProject(cfg ProjectConfig) ([]string, error) {
	if cfg.AppName == "" {
		cfg.AppName = "my-bento-app"
	}
	if cfg.Module == "" {
		cfg.Module = fmt.Sprintf("example.com/%s", cfg.AppName)
	}

	// Check if directory already exists
	if _, err := os.Stat(cfg.AppName); err == nil {
		return nil, fmt.Errorf("directory %q already exists", cfg.AppName)
	}

	// Create project directory
	if err := os.MkdirAll(cfg.AppName, 0755); err != nil {
		return nil, fmt.Errorf("create app directory: %w", err)
	}

	created := []string{}

	// Write go.mod
	goModPath := filepath.Join(cfg.AppName, "go.mod")
	if err := writeTemplate(goModPath, goModTemplate, cfg); err != nil {
		return created, fmt.Errorf("write go.mod: %w", err)
	}
	created = append(created, goModPath)

	// Write main.go
	mainPath := filepath.Join(cfg.AppName, "main.go")
	if err := writeTemplate(mainPath, mainGoTemplate, cfg); err != nil {
		return created, fmt.Errorf("write main.go: %w", err)
	}
	created = append(created, mainPath)

	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = cfg.AppName
	if err := tidyCmd.Run(); err != nil {
		// Don't fail - warn only
		return created, nil
	}

	// go.sum was created
	created = append(created, filepath.Join(cfg.AppName, "go.sum"))

	return created, nil
}

func writeTemplate(path, tmplStr string, data any) error {
	t, err := template.New("").Parse(tmplStr)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, data)
}

const goModTemplate = `module {{.Module}}

go 1.23

require (
	charm.land/bubbletea/v2 v2.0.0-rc.2
	charm.land/lipgloss/v2 v2.0.0-beta.3.0.20251106192539-4b304240aab7
	github.com/cloudboy-jh/bentotui v0.2.0
)
`

const mainGoTemplate = `package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/layout"
	"github.com/cloudboy-jh/bentotui/theme"
)

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

// model is the root Bubble Tea model. Add your panels and widgets here.
type model struct {
	root *layout.Split
	w, h int
}

func newModel() *model {
	// Placeholder content — replace with real registry components.
	// Run: bento add panel bar
	placeholder := &placeholder{}

	root := layout.Vertical(
		layout.Flex(1, placeholder),
	)
	return &model{root: root}
}

func (m *model) Init() tea.Cmd { return m.root.Init() }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.root.SetSize(m.w, m.h)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	updated, cmd := m.root.Update(msg)
	m.root = updated.(*layout.Split)
	return m, cmd
}

func (m *model) View() tea.View {
	t := theme.CurrentTheme()
	v := m.root.View()
	result := tea.NewView(viewString(v))
	result.AltScreen = true
	result.BackgroundColor = lipgloss.Color(t.Surface.Canvas)
	return result
}

// placeholder renders a welcome message until you add real components.
type placeholder struct{ w, h int }

func (p *placeholder) Init() tea.Cmd                          { return nil }
func (p *placeholder) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return p, nil }
func (p *placeholder) View() tea.View {
	t := theme.CurrentTheme()
	lines := []string{
		"{{.AppName}}",
		"",
		"Add components:  bento add panel bar dialog",
		"Switch theme:    theme.SetTheme(\"dracula\")",
		"Quit:            ctrl+c",
	}
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text.Primary)).
		Background(lipgloss.Color(t.Surface.Panel)).
		Width(p.w).Height(p.h)
	return tea.NewView(style.Render(strings.Join(lines, "\n")))
}
func (p *placeholder) SetSize(w, h int) { p.w = w; p.h = h }
func (p *placeholder) GetSize() (int, int) { return p.w, p.h }

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
`
