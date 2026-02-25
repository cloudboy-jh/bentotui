package main

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/app"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/focus"
	"github.com/cloudboy-jh/bentotui/layout"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/dialog"
	"github.com/cloudboy-jh/bentotui/ui/components/footer"
	"github.com/cloudboy-jh/bentotui/ui/components/panel"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

var actionLabels = []string{"Theme Picker", "Custom Dialog", "Confirm Dialog"}

func main() {
	t := theme.CurrentTheme()
	ft := footer.New(
		footer.Left("BentoTUI theme harness"),
		footer.Right("tab:focus  <-/->:action  enter:submit/run  /:command  d/x/q:actions  ctrl+c:quit"),
	)

	m := bentotui.New(
		bentotui.WithTheme(t),
		app.WithFooter(ft),
		bentotui.WithPages(
			bentotui.Page("harness", func() core.Page { return newHarnessPage(t) }),
		),
		bentotui.WithFooterBar(true),
		bentotui.WithFullScreen(true),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("bentotui test-tui failed: %v\n", err)
	}
}

type textBlock struct{ text string }

func newTextBlock(text string) *textBlock { return &textBlock{text: text} }
func (b *textBlock) SetText(text string)  { b.text = text }
func (b *textBlock) Init() tea.Cmd        { return nil }
func (b *textBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return b, nil
}
func (b *textBlock) View() tea.View { return tea.NewView(b.text) }

type harnessEventMsg struct{ text string }

type harnessPage struct {
	root *layout.Split

	headerPanel  *panel.Model
	mainPanel    *panel.Model
	actionsPanel *panel.Model

	headerText  *textBlock
	mainText    *textBlock
	actionsText *textBlock

	focus *focus.Manager
	input textinput.Model

	theme     theme.Theme
	themeName string
	width     int
	height    int
	startedAt time.Time
	events    []string
	actionIdx int
}

func newHarnessPage(t theme.Theme) *harnessPage {
	in := textinput.New()
	in.Prompt = "> "
	in.Placeholder = "Type text, /theme, /dialog, or /confirm"
	in.ShowSuggestions = true
	in.SetSuggestions([]string{"/theme", "/dialog", "/confirm"})
	in.SetStyles(inputStyles(t))

	p := &harnessPage{
		headerText:  newTextBlock(""),
		mainText:    newTextBlock(""),
		actionsText: newTextBlock(""),
		input:       in,
		theme:       t,
		themeName:   theme.CurrentThemeName(),
		startedAt:   time.Now(),
	}

	p.headerPanel = panel.New(panel.Theme(t), panel.Title("Header"), panel.Content(p.headerText))
	p.mainPanel = panel.New(panel.Theme(t), panel.Title("Main"), panel.Content(p.mainText))
	p.actionsPanel = panel.New(panel.Theme(t), panel.Title("Dialog Test Actions"), panel.Content(p.actionsText))

	p.rebuildLayout()
	p.syncInputFocus()
	p.refresh()
	return p
}

func (p *harnessPage) Init() tea.Cmd {
	return p.syncInputFocus()
}

func (p *harnessPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case theme.ThemeChangedMsg:
		p.themeName = v.Name
		p.applyTheme(v.Theme)
		p.log("theme switched to " + v.Name)
	case harnessEventMsg:
		if strings.TrimSpace(v.text) != "" {
			p.log(v.text)
		}
	case tea.KeyMsg:
		if key.Matches(v, p.focus.Bindings()...) {
			_, _ = p.focus.Update(v)
			cmd := p.syncInputFocus()
			p.log("focus moved to " + p.focusName())
			p.refresh()
			return p, cmd
		}

		switch v.String() {
		case "ctrl+c":
			return p, tea.Quit
		}

		if p.focusOnActions() {
			switch v.String() {
			case "left", "h":
				p.actionIdx = (p.actionIdx + len(actionLabels) - 1) % len(actionLabels)
			case "right", "l":
				p.actionIdx = (p.actionIdx + 1) % len(actionLabels)
			case "enter":
				p.log("ran action: " + strings.ToLower(actionLabels[p.actionIdx]))
				p.refresh()
				return p, p.runSelectedActionCmd()
			case "d":
				p.log("opened custom dialog via actions hotkey")
				p.refresh()
				return p, openCustomDialogCmd()
			case "x":
				p.log("opened confirm dialog via actions hotkey")
				p.refresh()
				return p, openConfirmDialogCmd()
			case "q":
				return p, tea.Quit
			}
			p.refresh()
			return p, nil
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
		layout.Fixed(7, p.actionsPanel),
	)
	p.focus = focus.New(focus.Ring(p.mainPanel, p.actionsPanel))
	if p.width > 0 && p.height > 0 {
		p.root.SetSize(p.width, p.height)
		p.updateInputWidth()
	}
}

func (p *harnessPage) applyTheme(t theme.Theme) {
	p.theme = t
	p.headerPanel.SetTheme(t)
	p.mainPanel.SetTheme(t)
	p.actionsPanel.SetTheme(t)
	p.input.SetStyles(inputStyles(t))
}

func (p *harnessPage) syncInputFocus() tea.Cmd {
	if p.focusOnActions() {
		p.input.Blur()
		return nil
	}
	return p.input.Focus()
}

func (p *harnessPage) focusOnActions() bool {
	if p.focus == nil {
		return false
	}
	return p.focus.Focused() == p.actionsPanel
}

func (p *harnessPage) focusName() string {
	if p.focusOnActions() {
		return "actions"
	}
	return "input"
}

func (p *harnessPage) runSelectedActionCmd() tea.Cmd {
	switch p.actionIdx {
	case 0:
		return openThemePickerCmd()
	case 1:
		return openCustomDialogCmd()
	default:
		return openConfirmDialogCmd()
	}
}

func (p *harnessPage) submitInput() tea.Cmd {
	text := strings.TrimSpace(p.input.Value())
	if text == "" {
		return nil
	}
	p.input.SetValue("")

	switch text {
	case "/theme":
		p.log("command accepted: /theme")
		return openThemePickerCmd()
	case "/dialog":
		p.log("command accepted: /dialog")
		return openCustomDialogCmd()
	case "/confirm":
		p.log("command accepted: /confirm")
		return openConfirmDialogCmd()
	case "/":
		p.log("command pending: type /theme, /dialog, or /confirm")
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
		"Primitive-first validation harness",
		fmt.Sprintf("Theme: %s", p.themeName),
		fmt.Sprintf("Focus: %s   Viewport: %dx%d   Uptime: %s", p.focusName(), p.width, p.height, time.Since(p.startedAt).Round(time.Second)),
	}
	p.headerText.SetText(strings.Join(headerLines, "\n"))

	mainLines := []string{
		"Prompt",
		inputSurface(p.input.View(), p.theme),
		"",
		"Commands: /theme  /dialog  /confirm",
		"",
		"Recent Events:",
	}
	if len(p.events) == 0 {
		mainLines = append(mainLines, "- none yet")
	} else {
		mainLines = append(mainLines, p.events...)
	}
	p.mainText.SetText(strings.Join(mainLines, "\n"))

	actionLines := []string{
		p.renderActionButtons(),
		"",
		"Tab switches focus between input and actions.",
		"Left/Right selects action. Enter runs selected action.",
		"Commands run on Enter. / starts command text in input.",
		"Hotkeys: d/x/q on actions focus  ctrl+c global quit",
	}
	p.actionsText.SetText(strings.Join(actionLines, "\n"))
}

func (p *harnessPage) renderActionButtons() string {
	sys := styles.New(p.theme)

	buttons := make([]string, 0, len(actionLabels))
	for i, label := range actionLabels {
		buttons = append(buttons, sys.ActionButton(i == p.actionIdx).Render(label))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, buttons...)
}

func inputSurface(view string, t theme.Theme) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(t.InputBG, t.ElementBG, t.SurfaceMuted))).
		Foreground(lipgloss.Color(t.Text)).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(pick(t.InputBorder, t.BorderFocused, t.Border))).
		Padding(0, 1).
		Render(view)
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

func openConfirmDialogCmd() tea.Cmd {
	return func() tea.Msg {
		return dialog.Open(dialog.Confirm{
			DialogTitle: "Confirm Action",
			Message:     "Press Enter to emit a confirm event.",
			OnConfirm: func() tea.Msg {
				return harnessEventMsg{text: "confirm accepted"}
			},
		})
	}
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
