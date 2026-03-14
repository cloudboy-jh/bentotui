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
		cfg.Module = fmt.Sprintf("example.com/%s", filepath.Base(cfg.AppName))
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
	github.com/cloudboy-jh/bentotui v0.3.4
)
`

const mainGoTemplate = `package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"
	"github.com/cloudboy-jh/bentotui/registry/bricks/dialog"
	"github.com/cloudboy-jh/bentotui/registry/bricks/input"
	"github.com/cloudboy-jh/bentotui/registry/bricks/surface"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
	"github.com/cloudboy-jh/bentotui/theme"
)

const version = "v0.3.4"
const wordmark = "" +
	"██████╗ ███████╗███╗   ██╗████████╗ ██████╗ \n" +
	"██╔══██╗██╔════╝████╗  ██║╚══██╔══╝██╔═══██╗\n" +
	"██████╔╝█████╗  ██╔██╗ ██║   ██║   ██║   ██║\n" +
	"██╔══██╗██╔══╝  ██║╚██╗██║   ██║   ██║   ██║\n" +
	"██████╔╝███████╗██║ ╚████║   ██║   ╚██████╔╝\n" +
	"╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝ "

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

type model struct {
	inputBox  *input.Model
	footerBar *bar.Model
	dialogs   *dialog.Manager
	width     int
	height    int
	inputW    int
}

func newModel() *model {
	inp := input.New()
	inp.SetPlaceholder("Ask anything... /theme /dialog")
	foot := bar.New(
		bar.FooterAnchored(),
		bar.Left("~ {{.AppName}}"),
		bar.Cards(
			bar.Card{Command: "enter", Label: "submit", Variant: bar.CardPrimary, Enabled: true, Priority: 3},
			bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
	)
	return &model{inputBox: inp, footerBar: foot, dialogs: dialog.New()}
}

func (m *model) Init() tea.Cmd { return m.inputBox.Focus() }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dialogs.IsOpen() {
		u, cmd := m.dialogs.Update(msg)
		m.dialogs = u.(*dialog.Manager)
		if tc, ok := msg.(theme.ThemeChangedMsg); ok {
			m.onThemeChange(tc)
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case dialog.OpenMsg:
		u, cmd := m.dialogs.Update(msg)
		m.dialogs = u.(*dialog.Manager)
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footerBar.SetSize(msg.Width, 1)
		m.dialogs.SetSize(msg.Width, msg.Height)
		m.inputW = clamp(m.width*6/10, 50, 90)
		m.inputBox.SetSize(m.inputW-5, 1)
		return m, nil
	case theme.ThemeChangedMsg:
		m.onThemeChange(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			val := strings.TrimSpace(m.inputBox.Value())
			if val == "" {
				return m, nil
			}
			m.inputBox.SetValue("")
			switch val {
			case "/theme":
				return m, openThemePicker()
			case "/dialog":
				return m, openSampleDialog()
			}
			return m, nil
		}
		u, cmd := m.inputBox.Update(msg)
		m.inputBox = u.(*input.Model)
		return m, cmd
	}
	return m, nil
}

func (m *model) View() tea.View {
	t := theme.CurrentTheme()
	canvasColor := lipgloss.Color(t.Surface.Canvas)
	if m.width == 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvasColor
		return v
	}

	dim := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	bright := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Primary))

	wm := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Accent)).Bold(true).Render(wordmark)

	inputBlockW := m.inputW
	if inputBlockW == 0 {
		inputBlockW = clamp(m.width*6/10, 50, 90)
	}
	contentW := max(1, inputBlockW-5)
	inputStr := viewString(m.inputBox.View())
	mkRow := func(fg, content string) string {
		return lipgloss.NewStyle().
			Background(lipgloss.Color(t.Input.BG)).
			Foreground(lipgloss.Color(fg)).
			PaddingLeft(2).PaddingRight(2).
			Width(contentW).
			Render(content)
	}
	blankRow := lipgloss.NewStyle().Background(lipgloss.Color(t.Input.BG)).Width(contentW + 4).Render(" ")
	inner := lipgloss.JoinVertical(lipgloss.Left,
		blankRow,
		mkRow(t.Input.FG, inputStr),
		mkRow(t.Text.Muted, "type /theme or /dialog"),
		blankRow,
	)
	block := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Input.BG)).
		Border(lipgloss.Border{Left: "|"}, false, false, false, true).
		BorderForeground(lipgloss.Color(t.Border.Focus)).
		Width(inputBlockW - 1).
		Render(inner)

	kbdStr := dim.Render("enter ") + bright.Render("submit") + dim.Render("  ctrl+c ") + bright.Render("quit")
	tipDot := lipgloss.NewStyle().Foreground(lipgloss.Color(t.State.Info)).Render("* Tip")
	tipStr := tipDot + dim.Render("  This file is yours. Edit anything.")

	body := rooms.RenderFunc(func(width, height int) string {
		center := func(s string) string {
			return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(s)
		}
		right := func(s string) string {
			line := lipgloss.NewStyle().Width(max(1, width-2)).Align(lipgloss.Right).Render(s)
			if width > 1 {
				return " " + line
			}
			return line
		}
		stack := strings.Join([]string{
			center(wm),
			"",
			center(block),
			"",
			right(kbdStr),
			"",
			center(tipStr),
		}, "\n")
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, stack)
	})

	screen := rooms.Focus(m.width, m.height, body, m.footerBar)
	surf := surface.New(m.width, m.height)
	surf.Fill(canvasColor)
	surf.Draw(0, 0, screen)
	if m.dialogs.IsOpen() {
		surf.DrawCenter(viewString(m.dialogs.View()))
	}

	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvasColor
	return v
}

func (m *model) onThemeChange(msg theme.ThemeChangedMsg) {
	_ = msg
}

func openThemePicker() tea.Cmd {
	return func() tea.Msg {
		h := len(theme.AvailableThemes()) + 8
		return dialog.Open(dialog.Custom{DialogTitle: "Themes", Content: dialog.NewThemePicker(), Width: 44, Height: h})
	}
}

func openSampleDialog() tea.Cmd {
	return func() tea.Msg {
		return dialog.Open(dialog.Confirm{DialogTitle: "Hello", Message: "Press Enter to confirm or Esc to cancel."})
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
`
