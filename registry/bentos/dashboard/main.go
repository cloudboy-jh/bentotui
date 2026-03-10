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
		metricATxt: mATxt,
		metricBTxt: mBTxt,
		metricCTxt: mCTxt,
		tableTxt:   tTxt,
		table:      t,
		badgeA:     b1,
		badgeB:     b2,
		badgeC:     b3,
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
		m.topBar.SetSize(m.width, 1)
		m.botBar.SetSize(m.width, 1)
		m.layoutPanels()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			seedTable(m.table, true)
			m.metricATxt.SetText("1.91M total")
			m.metricBTxt.SetText("0.31%")
			m.metricCTxt.SetText("Last deploy: 24m ago")
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

	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)

	m.tableTxt.SetText(viewString(m.table.View()))

	surf.Draw(0, 0, viewString(m.topBar.View()))

	bodyH := max(0, m.height-2)
	cardH := 6
	gap := 1

	if m.width >= 96 {
		cardW, _ := m.metricA.GetSize()
		y := 1
		x1 := gap
		x2 := x1 + cardW + gap
		x3 := x2 + cardW + gap

		surf.Draw(x1, y, viewString(m.metricA.View()))
		surf.Draw(x2, y, viewString(m.metricB.View()))
		surf.Draw(x3, y, viewString(m.metricC.View()))

		surf.Draw(x1+2, y+3, viewString(m.badgeA.View()))
		surf.Draw(x2+2, y+3, viewString(m.badgeB.View()))
		surf.Draw(x3+2, y+3, viewString(m.badgeC.View()))

		tableY := y + cardH + 1
		surf.Draw(1, tableY, viewString(m.tableP.View()))

		legend := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted)).Render("Tip: press r to refresh sample data")
		surf.Draw(max(1, m.width-lipgloss.Width(legend)-2), max(1, tableY+max(6, bodyH-cardH-1)-1), legend)
	} else {
		y := 1
		surf.Draw(1, y, viewString(m.metricA.View()))
		surf.Draw(3, y+3, viewString(m.badgeA.View()))
		y += cardH
		surf.Draw(1, y, viewString(m.metricB.View()))
		surf.Draw(3, y+3, viewString(m.badgeB.View()))
		y += cardH
		surf.Draw(1, y, viewString(m.metricC.View()))
		surf.Draw(3, y+3, viewString(m.badgeC.View()))
		y += cardH

		if y < m.height-2 {
			surf.Draw(1, y, viewString(m.tableP.View()))
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
	cardH := 6

	if m.width >= 96 {
		cardW := max(18, (m.width-4)/3)
		m.metricA.SetSize(cardW, cardH)
		m.metricB.SetSize(cardW, cardH)
		m.metricC.SetSize(cardW, cardH)
		tableH := max(6, bodyH-cardH-1)
		m.table.SetSize(max(10, m.width-4), max(4, tableH-3))
		m.tableP.SetSize(max(20, m.width-2), tableH)
		return
	}

	cardW := max(20, m.width-2)
	m.metricA.SetSize(cardW, cardH)
	m.metricB.SetSize(cardW, cardH)
	m.metricC.SetSize(cardW, cardH)

	tableH := max(6, bodyH-(cardH*3))
	m.table.SetSize(max(10, m.width-4), max(4, tableH-3))
	m.tableP.SetSize(max(20, m.width-2), tableH)
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
