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
	fmt.Println("🍱 bento")
	fmt.Println()

	// Determine app name — from args or interactive prompt
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

	// Module path — interactive with placeholder default
	modulePath := prompt("  Module path", fmt.Sprintf("example.com/%s", appName))
	if modulePath == "" {
		modulePath = fmt.Sprintf("example.com/%s", appName)
	}

	fmt.Println()
	fmt.Printf("  Creating %s/\n", appName)

	// Refuse to clobber an existing directory
	if _, err := os.Stat(appName); err == nil {
		fatal("directory %q already exists", appName)
	}

	check(os.MkdirAll(appName, 0755), "create app directory")

	data := struct {
		AppName string
		Module  string
	}{
		AppName: appName,
		Module:  modulePath,
	}

	writeTemplate(filepath.Join(appName, "go.mod"), goModTemplate, data)
	writeTemplate(filepath.Join(appName, "main.go"), mainGoTemplate, data)

	// Run go mod tidy inside the new directory
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
		fmt.Printf("    🍱 %s/go.sum\n", appName)
	}

	fmt.Println()
	fmt.Println("  Done. Your BentoTUI app is ready.")
	fmt.Println()
	fmt.Printf("    cd %s\n", appName)
	fmt.Println("    go run .")
	fmt.Println()
}

// writeTemplate executes a text/template string and writes the result to path,
// printing a 🍱 confirmation line on success.
func writeTemplate(path, tmplStr string, data any) {
	t, err := template.New("").Parse(tmplStr)
	check(err, "parse template for "+path)

	f, err := os.Create(path)
	check(err, "create "+path)
	defer f.Close()

	check(t.Execute(f, data), "write "+path)
	fmt.Printf("    🍱 %s\n", path)
}

// prompt prints a question, shows a default in parens, reads a line.
// Returns the default if the user hits enter with no input.
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

// ── embedded template strings ────────────────────────────────────────────────

const goModTemplate = `module {{.Module}}

go 1.25

require (
	github.com/cloudboy-jh/bentotui latest
)
`

const mainGoTemplate = `package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/bar"
	"github.com/cloudboy-jh/bentotui/ui/containers/panel"
)

func main() {
	t := theme.CurrentTheme()

	hd := bar.New(
		bar.LeftCard(bar.Card{Command: "{{.AppName}}", Variant: bar.CardMuted, Enabled: true}),
		bar.RightCard(bar.Card{Command: "theme", Label: theme.CurrentThemeName(), Variant: bar.CardPrimary, Enabled: true}),
	)

	ft := bar.New(
		bar.Cards(
			bar.Card{Command: "/theme", Label: "switch theme", Variant: bar.CardPrimary, Enabled: true},
		),
		bar.RightCard(bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted, Enabled: true}),
	)

	m := bentotui.New(
		bentotui.WithTheme(t),
		bentotui.WithHeader(hd),
		bentotui.WithFooter(ft),
		bentotui.WithPages(
			bentotui.Page("home", func() core.Page { return newHomePage(theme.CurrentTheme()) }),
		),
		bentotui.WithHeaderBar(true),
		bentotui.WithFooterBar(true),
		bentotui.WithFullScreen(true),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("{{.AppName}} failed: %v\n", err)
	}
}

type homePage struct {
	panel  *panel.Model
	theme  theme.Theme
	width  int
	height int
}

func newHomePage(t theme.Theme) *homePage {
	p := &homePage{theme: t}
	p.panel = panel.New(panel.Theme(t), panel.Title("Home"))
	return p
}

func (p *homePage) Init() tea.Cmd { return nil }

func (p *homePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case theme.ThemeChangedMsg:
		p.theme = v.Theme
		p.panel.SetTheme(v.Theme)
	case tea.KeyMsg:
		if v.String() == "ctrl+c" {
			return p, tea.Quit
		}
	}
	_, cmd := p.panel.Update(msg)
	return p, cmd
}

func (p *homePage) View() tea.View { return p.panel.View() }

func (p *homePage) SetSize(w, h int) {
	p.width = w
	p.height = h
	p.panel.SetSize(w, h)
}

func (p *homePage) GetSize() (int, int) { return p.width, p.height }
func (p *homePage) Title() string       { return "Home" }
`
