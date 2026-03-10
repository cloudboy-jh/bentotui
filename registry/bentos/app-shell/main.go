package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/components/bar"
	"github.com/cloudboy-jh/bentotui/registry/components/panel"
	"github.com/cloudboy-jh/bentotui/registry/components/surface"
	"github.com/cloudboy-jh/bentotui/registry/components/tabs"
	"github.com/cloudboy-jh/bentotui/theme"
)

type focusScope int

const (
	focusNav focusScope = iota
	focusTabs
)

type textBlock struct {
	text   string
	width  int
	height int
}

func (t *textBlock) Init() tea.Cmd                           { return nil }
func (t *textBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return t, nil }
func (t *textBlock) SetSize(width, height int)               { t.width, t.height = width, height }
func (t *textBlock) SetText(v string)                        { t.text = v }
func (t *textBlock) View() tea.View {
	if t.height <= 0 {
		return tea.NewView("")
	}
	lines := strings.Split(t.text, "\n")
	if len(lines) > t.height {
		lines = lines[:t.height]
	}
	return tea.NewView(strings.Join(lines, "\n"))
}

type model struct {
	width   int
	height  int
	scope   focusScope
	nav     []string
	cursor  int
	tabs    *tabs.Model
	topBar  *bar.Model
	botBar  *bar.Model
	left    *panel.Model
	center  *panel.Model
	right   *panel.Model
	leftTxt *textBlock
	mainTxt *textBlock
	metaTxt *textBlock
}

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func newModel() *model {
	tabModel := tabs.New(
		tabs.Tab{ID: "overview", Label: "Overview"},
		tabs.Tab{ID: "jobs", Label: "Jobs"},
		tabs.Tab{ID: "alerts", Label: "Alerts"},
	)
	tabModel.Blur()

	leftTxt := &textBlock{}
	mainTxt := &textBlock{}
	metaTxt := &textBlock{text: "Context\n- workspace: demo\n- branch: main\n- health: green"}

	m := &model{
		scope:   focusNav,
		nav:     []string{"Home", "Deploys", "Logs", "Settings", "About"},
		cursor:  0,
		tabs:    tabModel,
		topBar:  bar.New(bar.Left("bento app-shell"), bar.Right("workspace: demo")),
		botBar:  bar.New(),
		leftTxt: leftTxt,
		mainTxt: mainTxt,
		metaTxt: metaTxt,
	}

	m.left = panel.New(panel.Title("Navigation"), panel.Content(leftTxt))
	m.center = panel.New(panel.Title("Body"), panel.Content(mainTxt))
	m.right = panel.New(panel.Title("Inspector"), panel.Content(metaTxt), panel.Elevated())
	m.syncFocus()
	m.syncFooter()
	return m
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.topBar.SetSize(m.width, 1)
		m.botBar.SetSize(m.width, 1)
		m.layoutPanels()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.scope == focusNav {
				m.scope = focusTabs
			} else {
				m.scope = focusNav
			}
			m.syncFocus()
			m.syncFooter()
			return m, nil
		}

		if m.scope == focusNav {
			switch msg.String() {
			case "j", "down":
				if m.cursor < len(m.nav)-1 {
					m.cursor++
				}
				return m, nil
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}
				return m, nil
			}
		}

		u, cmd := m.tabs.Update(msg)
		m.tabs = u.(*tabs.Model)
		return m, cmd
	}

	return m, nil
}

func (m *model) View() tea.View {
	t := theme.CurrentTheme()
	canvas := lipgloss.Color(t.Surface.Canvas)

	if m.width == 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvas
		return v
	}

	bodyH := max(0, m.height-2)
	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)

	m.leftTxt.SetText(m.navText())
	m.mainTxt.SetText(m.bodyText())

	surf.Draw(0, 0, viewString(m.topBar.View()))

	if m.width < 78 {
		navH := clamp(bodyH/3, 6, 10)
		contentH := max(4, bodyH-navH)
		m.left.SetSize(max(20, m.width), navH)
		m.center.SetSize(max(20, m.width), contentH)
		surf.Draw(0, 1, viewString(m.left.View()))
		surf.Draw(0, 1+navH, viewString(m.center.View()))
	} else {
		surf.Draw(0, 1, viewString(m.left.View()))
		leftW, _ := m.left.GetSize()
		surf.Draw(leftW+1, 1, viewString(m.center.View()))
		if m.width >= 108 {
			centerW, _ := m.center.GetSize()
			surf.Draw(leftW+centerW+2, 1, viewString(m.right.View()))
		}
	}

	surf.Draw(0, m.height-1, viewString(m.botBar.View()))

	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func (m *model) layoutPanels() {
	bodyH := max(0, m.height-2)
	if m.width < 78 {
		m.left.SetSize(max(20, m.width), clamp(bodyH/3, 6, 10))
		m.center.SetSize(max(20, m.width), max(4, bodyH-clamp(bodyH/3, 6, 10)))
		m.right.SetSize(0, 0)
		return
	}

	leftW := clamp(m.width/4, 24, 32)
	if m.width >= 108 {
		rightW := clamp(m.width/5, 22, 28)
		centerW := max(26, m.width-leftW-rightW-2)
		m.left.SetSize(leftW, bodyH)
		m.center.SetSize(centerW, bodyH)
		m.right.SetSize(rightW, bodyH)
		return
	}

	centerW := max(26, m.width-leftW-1)
	m.left.SetSize(leftW, bodyH)
	m.center.SetSize(centerW, bodyH)
	m.right.SetSize(0, 0)
}

func (m *model) syncFocus() {
	if m.scope == focusNav {
		m.tabs.Blur()
		m.left.Focus()
		m.center.Blur()
		return
	}
	m.tabs.Focus()
	m.left.Blur()
	m.center.Focus()
}

func (m *model) syncFooter() {
	if m.scope == focusNav {
		m.botBar.SetLeft("scope: nav")
		m.botBar.SetCards([]bar.Card{
			{Command: "j/k", Label: "move", Variant: bar.CardPrimary, Enabled: true},
			{Command: "tab", Label: "focus tabs", Variant: bar.CardNormal, Enabled: true},
			{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true},
		})
		return
	}

	m.botBar.SetLeft("scope: tabs")
	m.botBar.SetCards([]bar.Card{
		{Command: "h/l", Label: "switch tab", Variant: bar.CardPrimary, Enabled: true},
		{Command: "tab", Label: "focus nav", Variant: bar.CardNormal, Enabled: true},
		{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true},
	})
}

func (m *model) navText() string {
	lines := make([]string, 0, len(m.nav))
	for i, n := range m.nav {
		prefix := "  "
		if i == m.cursor {
			prefix = "> "
		}
		lines = append(lines, prefix+n)
	}
	return strings.Join(lines, "\n")
}

func (m *model) bodyText() string {
	active := []string{"Overview", "Jobs", "Alerts"}[m.tabs.Active()]
	return strings.Join([]string{
		viewString(m.tabs.View()),
		"",
		"Reusable application shell",
		"- selected nav: " + m.nav[m.cursor],
		"- active tab: " + active,
		"",
		"This bento demonstrates a stable",
		"header/sidebar/body/status layout.",
	}, "\n")
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
