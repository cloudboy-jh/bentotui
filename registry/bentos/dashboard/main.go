package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/components/badge"
	"github.com/cloudboy-jh/bentotui/registry/components/bar"
	"github.com/cloudboy-jh/bentotui/registry/components/panel"
	"github.com/cloudboy-jh/bentotui/registry/components/surface"
	"github.com/cloudboy-jh/bentotui/registry/components/table"
	"github.com/cloudboy-jh/bentotui/registry/layouts"
	"github.com/cloudboy-jh/bentotui/theme"
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
	width  int
	height int

	topBar *bar.Model
	botBar *bar.Model

	metricA *panel.Model
	metricB *panel.Model
	metricC *panel.Model
	tableP  *panel.Model

	metricATxt *textBlock
	metricBTxt *textBlock
	metricCTxt *textBlock
	tableTxt   *textBlock

	table  *table.Model
	badgeA *badge.Model
	badgeB *badge.Model
	badgeC *badge.Model

	metricAValue string
	metricBValue string
	metricCValue string
}

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func newModel() *model {
	t := table.New("SERVICE", "STATUS", "LATENCY", "ERR%")
	seedTable(t, false)

	b1 := badge.New("+12.4%")
	b1.SetVariant(badge.VariantSuccess)
	b2 := badge.New("-3.1%")
	b2.SetVariant(badge.VariantWarning)
	b3 := badge.New("stable")
	b3.SetVariant(badge.VariantInfo)

	mATxt := &textBlock{text: "1.82M total"}
	mBTxt := &textBlock{text: "0.42%"}
	mCTxt := &textBlock{text: "Last deploy: 23m ago"}
	tTxt := &textBlock{}

	m := &model{
		topBar: bar.New(
			bar.Left("bento dashboard"),
			bar.Right("range: 24h"),
		),
		botBar: bar.New(
			bar.Left("cards + table composition"),
			bar.Cards(
				bar.Card{Command: "r", Label: "refresh", Variant: bar.CardPrimary, Enabled: true},
				bar.Card{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true},
			),
		),
		metricATxt:   mATxt,
		metricBTxt:   mBTxt,
		metricCTxt:   mCTxt,
		tableTxt:     tTxt,
		table:        t,
		badgeA:       b1,
		badgeB:       b2,
		badgeC:       b3,
		metricAValue: "1.82M total",
		metricBValue: "0.42%",
		metricCValue: "Last deploy: 23m ago",
	}

	m.metricA = panel.New(panel.Title("Requests"), panel.Content(mATxt))
	m.metricB = panel.New(panel.Title("Errors"), panel.Content(mBTxt))
	m.metricC = panel.New(panel.Title("Deploy"), panel.Content(mCTxt))
	m.tableP = panel.New(panel.Title("Service Health"), panel.Content(tTxt), panel.Elevated())

	return m
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			seedTable(m.table, true)
			m.metricAValue = "1.91M total"
			m.metricBValue = "0.31%"
			m.metricCValue = "Last deploy: 24m ago"
		}
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

	bodyH := max(1, m.height-2)
	m.metricATxt.SetText(m.metricAValue + "\n" + viewString(m.badgeA.View()))
	m.metricBTxt.SetText(m.metricBValue + "\n" + viewString(m.badgeB.View()))
	m.metricCTxt.SetText(m.metricCValue + "\n" + viewString(m.badgeC.View()))
	m.tableTxt.SetText(viewString(m.table.View()))

	body := ""
	if m.width >= 96 {
		body = layouts.Dashboard2x2(m.width, bodyH, m.metricA, m.metricB, m.metricC, m.tableP)
	} else {
		top := layouts.VSplit(m.width, max(1, bodyH/2), m.metricA, m.metricB)
		bottom := layouts.VSplit(m.width, max(1, bodyH-bodyH/2), m.metricC, m.tableP)
		body = layouts.VSplit(m.width, bodyH, layouts.Static(top), layouts.Static(bottom))
	}

	screen := layouts.Pancake(m.width, m.height, m.topBar, layouts.Static(body), m.botBar)
	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)
	surf.Draw(0, 0, screen)
	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func seedTable(t *table.Model, refreshed bool) {
	t.Clear()
	if refreshed {
		t.AddRow("api", "healthy", "36ms", "0.0")
		t.AddRow("workers", "healthy", "51ms", "0.1")
		t.AddRow("cache", "healthy", "74ms", "0.4")
		t.AddRow("queue", "healthy", "42ms", "0.0")
		return
	}
	t.AddRow("api", "healthy", "38ms", "0.1")
	t.AddRow("workers", "healthy", "55ms", "0.0")
	t.AddRow("cache", "degraded", "112ms", "1.7")
	t.AddRow("queue", "healthy", "47ms", "0.2")
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
