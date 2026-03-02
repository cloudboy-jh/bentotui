package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func runInit(args []string) {
	fmt.Println("bento init")
	fmt.Println()

	appName := ""
	if len(args) > 0 && args[0] != "--help" && args[0] != "-h" {
		appName = strings.TrimSpace(args[0])
	}
	if appName == "" {
		appName = prompt("  App name", "my-bento-app")
	}
	if appName == "" {
		appName = "my-bento-app"
	}

	modulePath := prompt("  Module path", fmt.Sprintf("example.com/%s", appName))
	if modulePath == "" {
		modulePath = fmt.Sprintf("example.com/%s", appName)
	}

	fmt.Println()
	fmt.Printf("  Creating %s/\n", appName)

	if _, err := os.Stat(appName); err == nil {
		fatal("directory %q already exists", appName)
	}

	check(os.MkdirAll(appName, 0755), "create app directory")

	data := struct {
		AppName string
		Module  string
	}{AppName: appName, Module: modulePath}

	writeTemplate(filepath.Join(appName, "go.mod"), goModTemplate, data)
	writeTemplate(filepath.Join(appName, "main.go"), mainGoTemplate, data)

	fmt.Println()
	fmt.Println("  Running go mod tidy...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = appName
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err := tidyCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: go mod tidy failed: %v\n", err)
		fmt.Println("  You may need to run it manually after editing go.mod.")
	} else {
		fmt.Printf("  %s/go.sum\n", appName)
	}

	fmt.Println()
	fmt.Println("  Done. Next steps:")
	fmt.Println()
	fmt.Printf("    cd %s\n", appName)
	fmt.Println("    bento add panel bar  # copy components you want")
	fmt.Println("    go run .")
	fmt.Println()
}

func writeTemplate(path, tmplStr string, data any) {
	t, err := template.New("").Parse(tmplStr)
	check(err, "parse template for "+path)
	f, err := os.Create(path)
	check(err, "create "+path)
	defer f.Close()
	check(t.Execute(f, data), "write "+path)
	fmt.Printf("  %s\n", path)
}

func prompt(question, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s (%s): ", question, defaultVal)
	} else {
		fmt.Printf("%s: ", question)
	}
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

// ── templates ─────────────────────────────────────────────────────────────────

const goModTemplate = `module {{.Module}}

go 1.23

require (
	charm.land/bubbletea/v2 v2.0.0-rc.2
	charm.land/lipgloss/v2 v2.0.0-beta.3.0.20251106192539-4b304240aab7
	github.com/cloudboy-jh/bentotui v0.2.0
)
`

// mainGoTemplate generates a minimal but complete BentoTUI app.
// It uses only the stable module deps (layout, theme) so it compiles
// without requiring any `bento add` step. Users add registry components
// as needed.
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
