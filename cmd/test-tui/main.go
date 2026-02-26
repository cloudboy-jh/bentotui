package main

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/app"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/layout"
	"github.com/cloudboy-jh/bentotui/core/surface"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/dialog"
	"github.com/cloudboy-jh/bentotui/ui/components/footer"
	"github.com/cloudboy-jh/bentotui/ui/components/panel"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
)

func main() {
	t := theme.CurrentTheme()
	ft := footer.New(
		footer.LeftCard(footer.Card{Command: "/pr", Label: " pull requests", Variant: footer.CardNormal, Enabled: true}),
		footer.Cards(
			footer.Card{Command: "/issue", Label: " issues", Variant: footer.CardPrimary, Enabled: true},
		),
		footer.RightCard(footer.Card{Command: "/branch", Label: " branches", Variant: footer.CardMuted, Enabled: true}),
	)

	m := bentotui.New(
		bentotui.WithTheme(t),
		app.WithFooter(ft),
		bentotui.WithPages(
			bentotui.Page("harness", func() core.Page { return newHarnessPage(t, "harness", "secondary") }),
			bentotui.Page("secondary", func() core.Page { return newHarnessPage(t, "secondary", "harness") }),
		),
		bentotui.WithFooterBar(true),
		bentotui.WithFullScreen(true),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("bentotui test-tui failed: %v\n", err)
	}
}

type textBlock struct {
	text   string
	width  int
	height int
}

func newTextBlock(text string) *textBlock { return &textBlock{text: text} }
func (b *textBlock) SetText(text string)  { b.text = text }
func (b *textBlock) Init() tea.Cmd        { return nil }
func (b *textBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return b, nil
}
func (b *textBlock) View() tea.View {
	if b.width <= 0 || b.height <= 0 {
		return tea.NewView(b.text)
	}
	return tea.NewView(surface.Region(b.text, b.width, b.height, "", ""))
}
func (b *textBlock) SetSize(width, height int) {
	b.width = width
	b.height = height
}
func (b *textBlock) GetSize() (int, int) { return b.width, b.height }

type harnessPage struct {
	root *layout.Split

	headerPanel *panel.Model
	mainPanel   *panel.Model

	headerText *textBlock
	mainText   *textBlock
	input      textinput.Model

	theme     theme.Theme
	themeName string
	pageName  string
	nextPage  string
	width     int
	height    int
	startedAt time.Time
	events    []string
}

func newHarnessPage(t theme.Theme, pageName, nextPage string) *harnessPage {
	in := textinput.New()
	in.Prompt = "> "
	in.Placeholder = "Type text, /pr, /issue, or /branch"
	in.ShowSuggestions = true
	in.SetSuggestions([]string{"/pr", "/issue", "/branch", "/dialog", "/theme", "/page"})
	in.SetStyles(inputStyles(t))

	p := &harnessPage{
		headerText: newTextBlock(""),
		mainText:   newTextBlock(""),
		input:      in,
		theme:      t,
		themeName:  theme.CurrentThemeName(),
		pageName:   pageName,
		nextPage:   nextPage,
		startedAt:  time.Now(),
	}
	_ = p.input.Focus()

	p.headerPanel = panel.New(panel.Theme(t), panel.Title("Header"), panel.Content(p.headerText))
	p.mainPanel = panel.New(panel.Theme(t), panel.Title("Main"), panel.Content(p.mainText))

	p.rebuildLayout()
	p.refresh()
	return p
}

func (p *harnessPage) Init() tea.Cmd {
	return p.input.Focus()
}

func (p *harnessPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case theme.ThemeChangedMsg:
		p.themeName = v.Name
		p.applyTheme(v.Theme)
		p.log("theme switched to " + v.Name)
	case tea.KeyMsg:
		if v.String() == "ctrl+c" {
			return p, tea.Quit
		}

		if v.String() == "enter" {
			cmd := p.submitInput()
			p.refresh()
			return p, cmd
		}

		updated, cmd := p.input.Update(v)
		p.input = updated
		p.refresh()
		return p, cmd
	}

	_, layoutCmd := p.root.Update(msg)
	p.refresh()
	return p, layoutCmd
}

func (p *harnessPage) View() tea.View { return p.root.View() }

func (p *harnessPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.root.SetSize(width, height)
	p.updateInputWidth()
	p.refresh()
}

func (p *harnessPage) GetSize() (int, int) { return p.width, p.height }
func (p *harnessPage) Title() string       { return "Harness" }

func (p *harnessPage) rebuildLayout() {
	p.root = layout.Vertical(
		layout.Fixed(6, p.headerPanel),
		layout.Flex(1, p.mainPanel),
	)
	if p.width > 0 && p.height > 0 {
		p.root.SetSize(p.width, p.height)
		p.updateInputWidth()
	}
}

func (p *harnessPage) applyTheme(t theme.Theme) {
	p.theme = t
	p.headerPanel.SetTheme(t)
	p.mainPanel.SetTheme(t)
	p.input.SetStyles(inputStyles(t))
}

func (p *harnessPage) submitInput() tea.Cmd {
	text := strings.TrimSpace(p.input.Value())
	if text == "" {
		return nil
	}
	p.input.SetValue("")

	switch text {
	case "/theme", "/issue":
		p.log("command accepted: /issue")
		return openThemePickerCmd()
	case "/dialog", "/pr":
		p.log("command accepted: /pr")
		return openCustomDialogCmd()
	case "/page", "/branch":
		p.log("command accepted: /branch -> " + p.nextPage)
		return navigateToCmd(p.nextPage)
	case "/":
		p.log("command pending: type /pr, /issue, or /branch")
		return nil
	default:
		if strings.HasPrefix(text, "/") {
			p.log("unknown command: " + text)
			return nil
		}
		p.log("submitted: " + text)
		return nil
	}
}

func (p *harnessPage) updateInputWidth() {
	panelWidth, _ := p.mainPanel.GetSize()
	p.input.SetWidth(max(20, panelWidth-8))
}

func (p *harnessPage) refresh() {
	headerLines := []string{
		"Command-driven validation harness",
		fmt.Sprintf("Page: %s   Next: %s", p.pageName, p.nextPage),
		fmt.Sprintf("Theme: %s", p.themeName),
		fmt.Sprintf("Viewport: %dx%d   Uptime: %s", p.width, p.height, time.Since(p.startedAt).Round(time.Second)),
	}
	p.headerText.SetText(strings.Join(headerLines, "\n"))

	mainLines := []string{
		"Prompt",
		inputSurface(p.input.View(), p.mainContentWidth(), p.theme),
		sectionDivider(p.mainContentWidth(), p.theme),
		"",
		"Cards:   /pr pull requests  /issue issues  /branch branches",
		"Slash:   /pr               /issue         /branch",
		"",
		"Recent Events:",
	}
	if len(p.events) == 0 {
		mainLines = append(mainLines, "- none yet")
	} else {
		mainLines = append(mainLines, p.events...)
	}
	p.mainText.SetText(strings.Join(mainLines, "\n"))
}

func inputSurface(view string, width int, t theme.Theme) string {
	if width <= 0 {
		return ""
	}
	bg := pick(t.InputBG, t.ElementBG, t.SurfaceMuted)
	return primitives.RenderInputRow(width, bg, t.Text, view)
}

func sectionDivider(width int, t theme.Theme) string {
	if width <= 0 {
		return ""
	}
	line := strings.Repeat("-", width)
	return lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.BorderSubtle, t.Muted))).Render(line)
}

func (p *harnessPage) mainContentWidth() int {
	mainWidth, _ := p.mainPanel.GetSize()
	if mainWidth <= 2 {
		return 0
	}
	return mainWidth - 2
}

func (p *harnessPage) log(s string) {
	entry := fmt.Sprintf("- [%s] %s", time.Now().Format("15:04:05"), s)
	p.events = append([]string{entry}, p.events...)
	if len(p.events) > 10 {
		p.events = p.events[:10]
	}
}

func openThemePickerCmd() tea.Cmd {
	return func() tea.Msg { return theme.OpenThemePicker() }
}

func openCustomDialogCmd() tea.Cmd {
	return func() tea.Msg {
		return dialog.Open(dialog.Custom{
			DialogTitle: "Custom Dialog",
			Content:     newTextBlock("Custom dialog rendering is healthy.\nPress Enter or Esc to close."),
			Width:       62,
			Height:      9,
		})
	}
}

func navigateToCmd(page string) tea.Cmd {
	return func() tea.Msg { return core.Navigate(page) }
}

func inputStyles(t theme.Theme) textinput.Styles {
	s := textinput.DefaultStyles(true)
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.BorderFocused)).Bold(true)
	s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text))
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
	s.Focused.Suggestion = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent))

	s.Blurred.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
	s.Blurred.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text))
	s.Blurred.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
	s.Blurred.Suggestion = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))

	s.Cursor.Color = lipgloss.Color(t.BorderFocused)
	s.Cursor.Blink = true
	return s
}

func pick(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
