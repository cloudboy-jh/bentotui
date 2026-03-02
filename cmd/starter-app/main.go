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
	"github.com/cloudboy-jh/bentotui/core/focus"
	"github.com/cloudboy-jh/bentotui/core/layout"
	"github.com/cloudboy-jh/bentotui/core/palette"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/bar"
	"github.com/cloudboy-jh/bentotui/ui/containers/dialog"
	"github.com/cloudboy-jh/bentotui/ui/containers/panel"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

// Nerd font glyphs
const (
	iconBento   = "󰁺" // bento box identity
	iconInfo    = "󰋱" // info circle
	iconCommand = ""  // prompt / terminal
	iconStatus  = ""  // pulse / health
	iconEvents  = ""  // list / log
	iconPalette = ""  // grid / commands
	iconFocus   = ""  // tab / cycle
	iconPage    = ""  // arrow right / navigate
	iconQuit    = ""  // power / exit
	iconTheme   = ""  // palette / brush
	iconSep     = ""  // powerline right arrow separator
)

// panelNames maps focus index → panel display name for the footer right card.
var panelNames = []string{"Info", "Command", "Status", "Events"}

// panelIcons maps focus index → nerd glyph for the footer right card.
var panelIcons = []string{iconInfo, iconCommand, iconStatus, iconEvents}

func main() {
	t := theme.CurrentTheme()

	ft := bar.New(
		bar.LeftCard(bar.Card{
			Command: iconBento + " bento",
			Variant: bar.CardMuted,
			Enabled: true,
		}),
		bar.Cards(
			bar.Card{Command: iconPalette + " /", Label: "commands", Variant: bar.CardNormal, Enabled: true},
			bar.Card{Command: iconFocus + " tab", Label: "focus", Variant: bar.CardNormal, Enabled: true},
			bar.Card{Command: iconPage + " /page", Label: "harness", Variant: bar.CardMuted, Enabled: true},
			bar.Card{Command: iconQuit + " ctrl+c", Label: "quit", Variant: bar.CardDanger, Enabled: true},
		),
		bar.RightCard(bar.Card{
			Command: iconInfo + " Info",
			Variant: bar.CardPrimary,
			Enabled: true,
		}),
	)

	m := bentotui.New(
		bentotui.WithTheme(t),
		app.WithFooter(ft),
		bentotui.WithPages(
			bentotui.Page("harness", func() core.Page {
				return newStarterPage(theme.CurrentTheme(), "harness", "secondary", ft)
			}),
			bentotui.Page("secondary", func() core.Page {
				return newStarterPage(theme.CurrentTheme(), "secondary", "harness", ft)
			}),
		),
		bentotui.WithCommands(
			dialog.Command{Label: "Open dialog", Group: "Suggested", Keybind: "/dialog", Action: func() tea.Msg { return openCustomDialogCmd()() }},
			dialog.Command{Label: "Switch theme", Group: "System", Keybind: "/theme", Action: func() tea.Msg { return theme.OpenThemePicker() }},
			dialog.Command{Label: "Next page", Group: "Session", Keybind: "/page", Action: func() tea.Msg { return core.Navigate("secondary") }},
		),
		bentotui.WithHeaderBar(false),
		bentotui.WithFooterBar(true),
		bentotui.WithFullScreen(true),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("bentotui starter-app failed: %v\n", err)
	}
}

// ── textBlock ────────────────────────────────────────────────────────────────

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
	return tea.NewView(primitives.Region(b.text, b.width, b.height, "", ""))
}
func (b *textBlock) SetSize(width, height int) {
	b.width = width
	b.height = height
}
func (b *textBlock) GetSize() (int, int) { return b.width, b.height }

// ── starterPage ──────────────────────────────────────────────────────────────

type starterPage struct {
	root *layout.Split

	// panels
	infoPanel    *panel.Model
	commandPanel *panel.Model
	statusPanel  *panel.Model
	eventsPanel  *panel.Model

	// content blocks
	infoText    *textBlock
	commandText *textBlock
	statusText  *textBlock
	eventsText  *textBlock

	// input
	input textinput.Model

	// focus manager
	fm       *focus.Manager
	focusIdx int // index into panelNames / panelIcons

	// state
	footer    *bar.Model
	theme     theme.Theme
	themeName string
	pageName  string
	nextPage  string
	width     int
	height    int
	startedAt time.Time
	events    []string
}

// commandPanelIdx is the index in the focus ring that owns the text input.
const commandPanelIdx = 1

func newStarterPage(t theme.Theme, pageName, nextPage string, ft *bar.Model) *starterPage {
	in := textinput.New()
	in.Prompt = "> "
	in.Placeholder = "Type /  /theme  /dialog  /page"
	in.ShowSuggestions = true
	in.SetSuggestions([]string{"/dialog", "/theme", "/page", "/"})
	in.SetStyles(inputStyles(t))

	infoText := newTextBlock("")
	commandText := newTextBlock("")
	statusText := newTextBlock("")
	eventsText := newTextBlock("")

	infoPanel := panel.New(panel.Theme(t), panel.Title(iconInfo+" Info"))
	commandPanel := panel.New(panel.Theme(t), panel.Title(iconCommand+" Command"))
	statusPanel := panel.New(panel.Theme(t), panel.Title(iconStatus+" Status"), panel.Elevated())
	eventsPanel := panel.New(panel.Theme(t), panel.Title(iconEvents+" Events"), panel.Elevated())

	// Set content after panels are constructed
	infoPanel = panel.New(panel.Theme(t), panel.Title(iconInfo+" Info"), panel.Content(infoText))
	commandPanel = panel.New(panel.Theme(t), panel.Title(iconCommand+" Command"), panel.Content(commandText))
	statusPanel = panel.New(panel.Theme(t), panel.Title(iconStatus+" Status"), panel.Elevated(), panel.Content(statusText))
	eventsPanel = panel.New(panel.Theme(t), panel.Title(iconEvents+" Events"), panel.Elevated(), panel.Content(eventsText))

	fm := focus.New(
		focus.Ring(infoPanel, commandPanel, statusPanel, eventsPanel),
		focus.Keys(
			key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next panel")),
			key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev panel")),
		),
		focus.Wrap(true),
		focus.Enabled(true),
	)

	p := &starterPage{
		infoPanel:    infoPanel,
		commandPanel: commandPanel,
		statusPanel:  statusPanel,
		eventsPanel:  eventsPanel,
		infoText:     infoText,
		commandText:  commandText,
		statusText:   statusText,
		eventsText:   eventsText,
		input:        in,
		fm:           fm,
		focusIdx:     0,
		footer:       ft,
		theme:        t,
		themeName:    theme.CurrentThemeName(),
		pageName:     pageName,
		nextPage:     nextPage,
		startedAt:    time.Now(),
	}

	// Start with Command panel focused so input is live
	_ = fm.SetIndex(commandPanelIdx)
	p.focusIdx = commandPanelIdx
	_ = p.input.Focus()

	p.rebuildLayout()
	p.refresh()
	return p
}

func (p *starterPage) Init() tea.Cmd {
	return p.input.Focus()
}

func (p *starterPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch v := msg.(type) {
	case theme.ThemeChangedMsg:
		p.themeName = v.Name
		p.applyTheme(v.Theme)
		p.log("theme switched to " + v.Name)
		p.refresh()
		return p, nil

	case focus.FocusChangedMsg:
		p.focusIdx = v.To
		p.syncInputFocus()
		p.updateFooterFocusState()
		p.refresh()
		return p, nil

	case tea.KeyMsg:
		if v.String() == "ctrl+c" {
			return p, tea.Quit
		}

		// Tab / shift+tab handled by focus manager regardless of which panel is focused
		if v.String() == "tab" || v.String() == "shift+tab" {
			_, fmCmd := p.fm.Update(v)
			cmds = append(cmds, fmCmd)
			p.refresh()
			return p, tea.Batch(cmds...)
		}

		// Only the command panel receives text input
		if p.focusIdx == commandPanelIdx {
			if v.String() == "enter" {
				cmd := p.submitInput()
				p.refresh()
				return p, cmd
			}
			updated, cmd := p.input.Update(v)
			p.input = updated
			cmds = append(cmds, cmd)
			p.refresh()
			return p, tea.Batch(cmds...)
		}

		// Global slash shortcut even from non-command panels
		if v.String() == "/" {
			p.log("command palette opened")
			p.refresh()
			return p, openCommandPaletteCmd()
		}
	}

	_, layoutCmd := p.root.Update(msg)
	cmds = append(cmds, layoutCmd)
	p.refresh()
	return p, tea.Batch(cmds...)
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

// ── layout ───────────────────────────────────────────────────────────────────

func (p *starterPage) rebuildLayout() {
	gutterColor := pick(p.theme.Border.Subtle, p.theme.Border.Normal)

	// Right column: fixed status top, flex events bottom, with subtle gutter.
	rightCol := layout.Vertical(
		layout.Fixed(8, p.statusPanel),
		layout.Flex(1, p.eventsPanel),
	).WithGutterColor(gutterColor)

	// Root: 3 columns 1:2:1 with subtle gutters between each column.
	p.root = layout.Horizontal(
		layout.Flex(1, p.infoPanel),
		layout.Flex(2, p.commandPanel),
		layout.Flex(1, rightCol),
	).WithGutterColor(gutterColor)

	if p.width > 0 && p.height > 0 {
		p.root.SetSize(p.width, p.height)
		p.updateInputWidth()
	}
}

// ── theme ────────────────────────────────────────────────────────────────────

func (p *starterPage) applyTheme(t theme.Theme) {
	p.theme = t
	p.infoPanel.SetTheme(t)
	p.commandPanel.SetTheme(t)
	p.statusPanel.SetTheme(t)
	p.eventsPanel.SetTheme(t)
	p.input.SetStyles(inputStyles(t))
}

// ── focus helpers ─────────────────────────────────────────────────────────────

// syncInputFocus focuses or blurs the text input based on which panel is active.
func (p *starterPage) syncInputFocus() {
	if p.focusIdx == commandPanelIdx {
		_ = p.input.Focus()
	} else {
		p.input.Blur()
	}
}

// updateFooterFocusState updates:
//   - right card to show the focused panel name + icon
//   - page toggle card variant (muted when input active, normal otherwise)
func (p *starterPage) updateFooterFocusState() {
	if p.footer == nil {
		return
	}

	// Right card: focused panel name
	name := "Info"
	icon := iconInfo
	if p.focusIdx >= 0 && p.focusIdx < len(panelNames) {
		name = panelNames[p.focusIdx]
		icon = panelIcons[p.focusIdx]
	}
	p.footer.SetRightCard(bar.Card{
		Command: icon + " " + name,
		Variant: bar.CardPrimary,
		Enabled: true,
	})

	// Page toggle card (index 2 in the cards slice): muted when input is focused
	pageVariant := bar.CardNormal
	if p.focusIdx == commandPanelIdx {
		pageVariant = bar.CardMuted
	}
	p.footer.SetCards([]bar.Card{
		{Command: iconPalette + " /", Label: "commands", Variant: bar.CardNormal, Enabled: true},
		{Command: iconFocus + " tab", Label: "focus", Variant: bar.CardNormal, Enabled: true},
		{Command: iconPage + " /page", Label: p.nextPage, Variant: pageVariant, Enabled: true},
		{Command: iconQuit + " ctrl+c", Label: "quit", Variant: bar.CardDanger, Enabled: true},
	})
}

// ── input ────────────────────────────────────────────────────────────────────

func (p *starterPage) updateInputWidth() {
	w, _ := p.commandPanel.GetSize()
	contentW := max(20, w-4)
	p.input.SetWidth(contentW)
}

func (p *starterPage) submitInput() tea.Cmd {
	text := strings.TrimSpace(p.input.Value())
	if text == "" {
		return nil
	}
	p.input.SetValue("")

	switch text {
	case "/theme":
		p.log("command: /theme")
		return openThemePickerCmd()
	case "/dialog":
		p.log("command: /dialog")
		return openCustomDialogCmd()
	case "/page":
		p.log("command: /page → " + p.nextPage)
		return navigateToCmd(p.nextPage)
	case "/":
		p.log("command palette opened")
		return openCommandPaletteCmd()
	default:
		if strings.HasPrefix(text, "/") {
			p.log("unknown command: " + text)
			return nil
		}
		p.log("submitted: " + text)
		return nil
	}
}

// ── refresh ───────────────────────────────────────────────────────────────────

func (p *starterPage) refresh() {
	p.refreshInfo()
	p.refreshCommand()
	p.refreshStatus()
	p.refreshEvents()
}

func (p *starterPage) refreshInfo() {
	uptime := time.Since(p.startedAt).Round(time.Second)
	lines := []string{
		"",
		fmt.Sprintf("  page     %s", p.pageName),
		fmt.Sprintf("  next     %s", p.nextPage),
		"",
		fmt.Sprintf("  theme    %s", p.themeName),
		"",
		fmt.Sprintf("  size     %dx%d", p.width, p.height),
		fmt.Sprintf("  uptime   %s", uptime),
		"",
		fmt.Sprintf("  focus    %s", safeGet(panelNames, p.focusIdx)),
	}
	p.infoText.SetText(strings.Join(lines, "\n"))
}

func (p *starterPage) refreshCommand() {
	w, _ := p.commandPanel.GetSize()
	contentW := max(1, w-4)

	inputRow := inputSurface(p.input.View(), contentW, p.theme)
	divider := sectionDivider(contentW, p.theme)

	lines := []string{
		"",
		inputRow,
		divider,
		"",
		"  " + iconPalette + "  /          commands",
		"  " + iconTheme + "  /theme      switch theme",
		"  " + iconCommand + "  /dialog     open dialog",
		"  " + iconPage + "  /page       next page",
		"  " + iconFocus + "  tab         cycle focus",
	}
	p.commandText.SetText(strings.Join(lines, "\n"))
}

func (p *starterPage) refreshStatus() {
	lines := []string{
		"",
		"  header     hidden",
		"  footer     " + checkmark(true),
		"  palette    " + checkmark(true),
		"  focus      " + checkmark(true),
		"  theme      " + checkmark(true),
	}
	p.statusText.SetText(strings.Join(lines, "\n"))
}

func (p *starterPage) refreshEvents() {
	lines := []string{""}
	if len(p.events) == 0 {
		lines = append(lines, "  no events yet")
	} else {
		for _, e := range p.events {
			lines = append(lines, "  "+e)
		}
	}
	p.eventsText.SetText(strings.Join(lines, "\n"))
}

func (p *starterPage) log(s string) {
	entry := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), s)
	p.events = append([]string{entry}, p.events...)
	if len(p.events) > 20 {
		p.events = p.events[:20]
	}
}

// ── style helpers ─────────────────────────────────────────────────────────────

func inputSurface(view string, width int, t theme.Theme) string {
	if width <= 0 {
		return ""
	}
	input := styles.New(t).InputColors()
	return primitives.RenderRow(width, input.BG, input.FG, view)
}

func sectionDivider(width int, t theme.Theme) string {
	if width <= 0 {
		return ""
	}
	line := strings.Repeat("─", width)
	return lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Border.Subtle, t.Text.Muted))).Render(line)
}

func checkmark(ok bool) string {
	if ok {
		return "✓"
	}
	return "✗"
}

func safeGet(s []string, i int) string {
	if i >= 0 && i < len(s) {
		return s[i]
	}
	return ""
}

// ── cmd helpers ───────────────────────────────────────────────────────────────

func openThemePickerCmd() tea.Cmd {
	return func() tea.Msg { return theme.OpenThemePicker() }
}

func openCommandPaletteCmd() tea.Cmd {
	return func() tea.Msg { return palette.OpenCommandPalette() }
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

// ── style builders ────────────────────────────────────────────────────────────

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

// ── misc ──────────────────────────────────────────────────────────────────────

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
