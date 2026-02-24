package main

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/app"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/dialog"
	"github.com/cloudboy-jh/bentotui/focus"
	"github.com/cloudboy-jh/bentotui/layout"
	"github.com/cloudboy-jh/bentotui/panel"
	"github.com/cloudboy-jh/bentotui/router"
	"github.com/cloudboy-jh/bentotui/statusbar"
	"github.com/cloudboy-jh/bentotui/theme"
)

const (
	compactWidthBreakpoint  = 116
	compactHeightBreakpoint = 34
)

func main() {
	t := theme.Preset("amber")
	status := statusbar.New(
		statusbar.Left("BentoTUI internal control room"),
		statusbar.Right("1 Home  2 Inspect  Tab Focus  O Dialog  Q Quit"),
	)

	m := bentotui.New(
		bentotui.WithTheme(t),
		app.WithStatus(status),
		bentotui.WithPages(
			bentotui.Page("home", func() core.Page { return newHomePage(t) }),
			bentotui.Page("inspect", func() core.Page { return newInspectPage(t) }),
		),
		bentotui.WithStatusBar(true),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("bentotui test-tui failed: %v\n", err)
	}
}

type textBlock struct {
	text string
}

func newTextBlock(text string) *textBlock { return &textBlock{text: text} }
func (b *textBlock) SetText(text string)  { b.text = text }
func (b *textBlock) Init() tea.Cmd        { return nil }
func (b *textBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_ = msg
	return b, nil
}
func (b *textBlock) View() tea.View { return tea.NewView(b.text) }

type homePage struct {
	theme theme.Theme

	root                                             *layout.Split
	headerPanel, controlsPanel, statePanel, logPanel *panel.Model
	notesPanel                                       *panel.Model
	headerText, controlsText, stateText, logText     *textBlock
	notesText                                        *textBlock
	focus                                            *focus.Manager

	compact       bool
	width, height int
	counter       int
	logs          []string
}

func newHomePage(t theme.Theme) *homePage {
	p := &homePage{
		theme:        t,
		headerText:   newTextBlock(""),
		controlsText: newTextBlock(""),
		stateText:    newTextBlock(""),
		logText:      newTextBlock(""),
		notesText:    newTextBlock(""),
	}
	p.headerPanel = panel.New(panel.Theme(t), panel.Title("Session"), panel.Content(p.headerText))
	p.controlsPanel = panel.New(panel.Theme(t), panel.Title("Controls"), panel.Content(p.controlsText))
	p.statePanel = panel.New(panel.Theme(t), panel.Title("State"), panel.Content(p.stateText))
	p.logPanel = panel.New(panel.Theme(t), panel.Title("Activity"), panel.Content(p.logText))
	p.notesPanel = panel.New(panel.Theme(t), panel.Title("Notes"), panel.Content(p.notesText))
	p.rebuildLayout()
	p.refresh()
	return p
}

func (p *homePage) Init() tea.Cmd { return nil }

func (p *homePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "q", "ctrl+c":
			return p, tea.Quit
		case "1":
			return p, func() tea.Msg { return router.Navigate("home") }
		case "2":
			return p, func() tea.Msg { return router.Navigate("inspect") }
		case "+", "=":
			p.counter++
			p.log("counter incremented")
		case "-", "_":
			p.counter--
			p.log("counter decremented")
		case "r":
			p.counter = 0
			p.log("counter reset")
		case "o":
			return p, func() tea.Msg {
				return dialog.Open(dialog.Custom{
					DialogTitle: "Control Room Overlay",
					Content:     newTextBlock("Overlay layer is active on Home.\nPress Enter or Esc to close."),
					Width:       56,
					Height:      8,
				})
			}
		}
	}

	if p.focus != nil {
		_, _ = p.focus.Update(msg)
		if keyMsg, ok := msg.(tea.KeyMsg); ok && key.Matches(keyMsg, p.focus.Bindings()...) {
			p.log("focus moved")
		}
	}

	p.refresh()
	_, cmd := p.root.Update(msg)
	return p, cmd
}

func (p *homePage) View() tea.View { return p.root.View() }

func (p *homePage) SetSize(width, height int) {
	p.width, p.height = width, height
	p.updateCompact()
	p.root.SetSize(width, height)
}

func (p *homePage) GetSize() (int, int) { return p.width, p.height }
func (p *homePage) Title() string       { return "Home" }

func (p *homePage) updateCompact() {
	next := p.width < compactWidthBreakpoint || p.height < compactHeightBreakpoint
	if next == p.compact {
		return
	}
	p.compact = next
	p.rebuildLayout()
}

func (p *homePage) rebuildLayout() {
	if p.compact {
		main := layout.Vertical(
			layout.Fixed(8, p.controlsPanel),
			layout.Fixed(7, p.statePanel),
			layout.Flex(1, p.logPanel),
		)
		p.root = layout.Vertical(
			layout.Fixed(4, p.headerPanel),
			layout.Flex(1, main),
		)
		p.focus = focus.New(focus.Ring(p.statePanel, p.logPanel))
	} else {
		left := layout.Vertical(
			layout.Fixed(10, p.controlsPanel),
			layout.Flex(1, p.logPanel),
		)
		right := layout.Vertical(
			layout.Fixed(9, p.statePanel),
			layout.Flex(1, p.notesPanel),
		)
		main := layout.Horizontal(layout.Fixed(40, left), layout.Flex(1, right))
		p.root = layout.Vertical(
			layout.Fixed(4, p.headerPanel),
			layout.Flex(1, main),
		)
		p.focus = focus.New(focus.Ring(p.statePanel, p.logPanel, p.notesPanel))
	}
	if p.width > 0 && p.height > 0 {
		p.root.SetSize(p.width, p.height)
	}
}

func (p *homePage) refresh() {
	mode := "wide"
	if p.compact {
		mode = "compact"
	}
	p.headerText.SetText(strings.Join([]string{
		fmt.Sprintf("Route: home  |  Mode: %s", mode),
		"Use Tab / Shift+Tab to cycle focus panels.",
	}, "\n"))

	p.controlsText.SetText(strings.Join([]string{
		"1: home",
		"2: inspect",
		"tab: next pane",
		"shift+tab: prev pane",
		"+: increment",
		"-: decrement",
		"r: reset",
		"o: dialog",
		"q: quit",
	}, "\n"))

	p.stateText.SetText(strings.Join([]string{
		fmt.Sprintf("Counter: %d", p.counter),
		fmt.Sprintf("Events: %d", len(p.logs)),
		fmt.Sprintf("Focus: %s", p.focusName()),
		fmt.Sprintf("Time: %s", time.Now().Format("15:04:05")),
	}, "\n"))

	p.notesText.SetText(strings.Join([]string{
		"Solid-surface validation:",
		"- focus borders",
		"- router + cache",
		"- dialog overlay",
		"- compact behavior",
	}, "\n"))

	if len(p.logs) == 0 {
		p.logText.SetText("No events yet.")
	} else {
		p.logText.SetText(strings.Join(p.logs, "\n"))
	}
}

func (p *homePage) log(s string) {
	entry := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), s)
	p.logs = append([]string{entry}, p.logs...)
	if len(p.logs) > 18 {
		p.logs = p.logs[:18]
	}
}

func (p *homePage) focusName() string {
	switch p.focus.Focused() {
	case p.statePanel:
		return "state"
	case p.logPanel:
		return "activity"
	case p.notesPanel:
		return "notes"
	default:
		return "unknown"
	}
}

type inspectPage struct {
	theme theme.Theme

	root                                           *layout.Split
	headerPanel, summaryPanel, checkPanel          *panel.Model
	resultPanel                                    *panel.Model
	headerText, summaryText, checkText, resultText *textBlock
	focus                                          *focus.Manager

	compact       bool
	width, height int
	runs          int
	lastResult    string
}

func newInspectPage(t theme.Theme) *inspectPage {
	p := &inspectPage{
		theme:       t,
		headerText:  newTextBlock(""),
		summaryText: newTextBlock(""),
		checkText:   newTextBlock(""),
		resultText:  newTextBlock(""),
		lastResult:  "ready",
	}
	p.headerPanel = panel.New(panel.Theme(t), panel.Title("Session"), panel.Content(p.headerText))
	p.summaryPanel = panel.New(panel.Theme(t), panel.Title("Inspection"), panel.Content(p.summaryText))
	p.checkPanel = panel.New(panel.Theme(t), panel.Title("Checks"), panel.Content(p.checkText))
	p.resultPanel = panel.New(panel.Theme(t), panel.Title("Latest Result"), panel.Content(p.resultText))
	p.rebuildLayout()
	p.refresh()
	return p
}

func (p *inspectPage) Init() tea.Cmd { return nil }

func (p *inspectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.KeyMsg:
		switch v.String() {
		case "q", "ctrl+c":
			return p, tea.Quit
		case "1":
			return p, func() tea.Msg { return router.Navigate("home") }
		case "2":
			return p, func() tea.Msg { return router.Navigate("inspect") }
		case "o":
			return p, func() tea.Msg {
				return dialog.Open(dialog.Custom{
					DialogTitle: "Inspect Overlay",
					Content:     newTextBlock("Overlay works on Inspect too."),
					Width:       48,
					Height:      8,
				})
			}
		case "c":
			p.runs++
			p.lastResult = fmt.Sprintf("check run #%d at %s", p.runs, time.Now().Format("15:04:05"))
		}
	}

	if p.focus != nil {
		_, _ = p.focus.Update(msg)
	}

	p.refresh()
	_, cmd := p.root.Update(msg)
	return p, cmd
}

func (p *inspectPage) View() tea.View { return p.root.View() }

func (p *inspectPage) SetSize(width, height int) {
	p.width, p.height = width, height
	p.updateCompact()
	p.root.SetSize(width, height)
}

func (p *inspectPage) GetSize() (int, int) { return p.width, p.height }
func (p *inspectPage) Title() string       { return "Inspect" }

func (p *inspectPage) updateCompact() {
	next := p.width < compactWidthBreakpoint || p.height < compactHeightBreakpoint
	if next == p.compact {
		return
	}
	p.compact = next
	p.rebuildLayout()
}

func (p *inspectPage) rebuildLayout() {
	if p.compact {
		main := layout.Vertical(
			layout.Fixed(11, p.summaryPanel),
			layout.Fixed(8, p.checkPanel),
			layout.Flex(1, p.resultPanel),
		)
		p.root = layout.Vertical(
			layout.Fixed(4, p.headerPanel),
			layout.Flex(1, main),
		)
	} else {
		right := layout.Vertical(layout.Fixed(10, p.checkPanel), layout.Flex(1, p.resultPanel))
		main := layout.Horizontal(layout.Flex(2, p.summaryPanel), layout.Flex(1, right))
		p.root = layout.Vertical(
			layout.Fixed(4, p.headerPanel),
			layout.Flex(1, main),
		)
	}
	p.focus = focus.New(focus.Ring(p.summaryPanel, p.checkPanel, p.resultPanel))
	if p.width > 0 && p.height > 0 {
		p.root.SetSize(p.width, p.height)
	}
}

func (p *inspectPage) refresh() {
	mode := "wide"
	if p.compact {
		mode = "compact"
	}
	p.headerText.SetText(strings.Join([]string{
		fmt.Sprintf("Route: inspect  |  Mode: %s", mode),
		"Inspect route validates split + focus + overlay behavior.",
	}, "\n"))

	p.summaryText.SetText(strings.Join([]string{
		fmt.Sprintf("Manual checks: %d", p.runs),
		fmt.Sprintf("Focus: %s", p.focusName()),
		"",
		"Validation list:",
		"- router caching",
		"- split sizing",
		"- focused borders",
		"- compact layout",
		"- overlay composition",
	}, "\n"))

	p.checkText.SetText(strings.Join([]string{
		"1: home",
		"2: inspect",
		"tab: focus cycle",
		"c: run check",
		"o: dialog",
		"q: quit",
	}, "\n"))

	p.resultText.SetText(strings.Join([]string{"Last result:", p.lastResult}, "\n"))
}

func (p *inspectPage) focusName() string {
	switch p.focus.Focused() {
	case p.summaryPanel:
		return "inspection"
	case p.checkPanel:
		return "checks"
	case p.resultPanel:
		return "result"
	default:
		return "unknown"
	}
}
