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
	"github.com/cloudboy-jh/bentotui/ui/containers/dialog"
	"github.com/cloudboy-jh/bentotui/ui/containers/footer"
	"github.com/cloudboy-jh/bentotui/ui/containers/header"
	"github.com/cloudboy-jh/bentotui/ui/containers/panel"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

func main() {
	t := theme.CurrentTheme()
	hd := header.New(
		header.LeftCard(header.Card{Command: "starter", Label: "bentotui", Variant: header.CardMuted, Enabled: true}),
		header.Cards(
			header.Card{Command: "page", Label: "harness", Variant: header.CardNormal, Enabled: true},
		),
		header.RightCard(header.Card{Command: "theme", Label: theme.CurrentThemeName(), Variant: header.CardPrimary, Enabled: true}),
	)
	ft := footer.New(
		footer.LeftCard(footer.Card{Command: "/dialog", Label: "open dialog", Variant: footer.CardNormal, Enabled: true}),
		footer.Cards(
			footer.Card{Command: "/theme", Label: "switch theme", Variant: footer.CardPrimary, Enabled: true},
		),
		footer.RightCard(footer.Card{Command: "/page", Label: "next page", Variant: footer.CardMuted, Enabled: true}),
	)

	m := bentotui.New(
		bentotui.WithTheme(t),
		app.WithHeader(hd),
		app.WithFooter(ft),
		bentotui.WithPages(
			bentotui.Page("harness", func() core.Page { return newStarterPage(theme.CurrentTheme(), "harness", "secondary") }),
			bentotui.Page("secondary", func() core.Page { return newStarterPage(theme.CurrentTheme(), "secondary", "harness") }),
		),
		bentotui.WithHeaderBar(true),
		bentotui.WithFooterBar(true),
		bentotui.WithFullScreen(true),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("bentotui starter-app failed: %v\n", err)
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

type starterPage struct {
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

func newStarterPage(t theme.Theme, pageName, nextPage string) *starterPage {
	in := textinput.New()
	in.Prompt = "> "
	in.Placeholder = "Type text, /dialog, /theme, or /page"
	in.ShowSuggestions = true
	in.SetSuggestions([]string{"/dialog", "/theme", "/page"})
	in.SetStyles(inputStyles(t))

	p := &starterPage{
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

func (p *starterPage) Init() tea.Cmd {
	return p.input.Focus()
}

func (p *starterPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (p *starterPage) View() tea.View { return p.root.View() }

func (p *starterPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.root.SetSize(width, height)
	p.updateInputWidth()
	p.refresh()
}

func (p *starterPage) GetSize() (int, int) { return p.width, p.height }
func (p *starterPage) Title() string       { return "Starter" }

func (p *starterPage) rebuildLayout() {
	p.root = layout.Vertical(
		layout.Fixed(6, p.headerPanel),
		layout.Flex(1, p.mainPanel),
	)
	if p.width > 0 && p.height > 0 {
		p.root.SetSize(p.width, p.height)
		p.updateInputWidth()
	}
}

func (p *starterPage) applyTheme(t theme.Theme) {
	p.theme = t
	p.headerPanel.SetTheme(t)
	p.mainPanel.SetTheme(t)
	p.input.SetStyles(inputStyles(t))
}

func (p *starterPage) submitInput() tea.Cmd {
	text := strings.TrimSpace(p.input.Value())
	if text == "" {
		return nil
	}
	p.input.SetValue("")

	switch text {
	case "/theme", "/issue":
		p.log("command accepted: /theme")
		return openThemePickerCmd()
	case "/dialog", "/pr":
		p.log("command accepted: /dialog")
		return openCustomDialogCmd()
	case "/page", "/branch":
		p.log("command accepted: /page -> " + p.nextPage)
		return navigateToCmd(p.nextPage)
	case "/":
		p.log("command pending: type /dialog, /theme, or /page")
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

func (p *starterPage) updateInputWidth() {
	p.input.SetWidth(max(20, p.mainContentWidth()-2))
}

func (p *starterPage) refresh() {
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
		"Commands: /dialog  /theme  /page",
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
	input := styles.New(t).InputColors()
	return primitives.RenderInputRow(width, input.BG, input.FG, view)
}

func sectionDivider(width int, t theme.Theme) string {
	if width <= 0 {
		return ""
	}
	line := strings.Repeat("-", width)
	return lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Border.Subtle, t.Text.Muted))).Render(line)
}

func (p *starterPage) mainContentWidth() int {
	mainWidth, _ := p.mainPanel.GetSize()
	if mainWidth <= 2 {
		return 0
	}
	return mainWidth - 2
}

func (p *starterPage) log(s string) {
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
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Border.Focus)).Bold(true)
	s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Input.FG))
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Input.Placeholder))
	s.Focused.Suggestion = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Accent))

	s.Blurred.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	s.Blurred.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Input.FG))
	s.Blurred.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Input.Placeholder))
	s.Blurred.Suggestion = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))

	s.Cursor.Color = lipgloss.Color(t.Input.Cursor)
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
