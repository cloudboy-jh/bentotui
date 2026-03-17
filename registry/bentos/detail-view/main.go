package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"
	elevatedcard "github.com/cloudboy-jh/bentotui/registry/bricks/elevated-card"
	"github.com/cloudboy-jh/bentotui/registry/bricks/list"
	"github.com/cloudboy-jh/bentotui/registry/bricks/surface"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
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

type item struct {
	title string
	kind  string
	meta  string
}

type model struct {
	width  int
	height int

	items  []item
	cursor int

	navList     *list.Model
	navCard     *elevatedcard.Model
	detailText  *textBlock
	detailCard  *elevatedcard.Model
	sessionText *textBlock
	sessionCard *elevatedcard.Model
	footer      *bar.Model
}

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func newModel() *model {
	items := []item{
		{title: "Account", kind: "section", meta: "identity, profile, notifications"},
		{title: "Billing", kind: "section", meta: "plans, seats, invoices"},
		{title: "Integrations", kind: "section", meta: "github, slack, webhooks"},
		{title: "Security", kind: "section", meta: "2FA, sessions, tokens"},
	}

	nav := list.New(24)
	for _, it := range items {
		nav.AppendRow(list.Row{Primary: it.title, Secondary: it.kind, Tone: list.ToneInfo, RightStat: "open"})
	}
	nav.SetCursor(0)

	navCard := elevatedcard.New(
		elevatedcard.Title("Sections"),
		elevatedcard.CardVariant(elevatedcard.VariantDense),
		elevatedcard.Content(nav),
	)
	detailText := &textBlock{}
	detailCard := elevatedcard.New(
		elevatedcard.Title("Detail"),
		elevatedcard.CardVariant(elevatedcard.VariantEmphasis),
		elevatedcard.Content(detailText),
	)
	sessionText := &textBlock{}
	sessionCard := elevatedcard.New(
		elevatedcard.Title("Session"),
		elevatedcard.CardVariant(elevatedcard.VariantDense),
		elevatedcard.Content(sessionText),
	)

	footer := bar.New(
		bar.FooterAnchored(),
		bar.Left("detail-view"),
		bar.Cards(
			bar.Card{Command: "up/down", Label: "select", Variant: bar.CardPrimary, Enabled: true, Priority: 4},
			bar.Card{Command: "enter", Label: "open", Variant: bar.CardNormal, Enabled: true, Priority: 3},
			bar.Card{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
	)

	m := &model{
		items:       items,
		navList:     nav,
		navCard:     navCard,
		detailText:  detailText,
		detailCard:  detailCard,
		sessionText: sessionText,
		sessionCard: sessionCard,
		footer:      footer,
	}
	m.syncDetail()
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.navList.SetCursor(m.cursor)
				m.syncDetail()
			}
		case "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.navList.SetCursor(m.cursor)
				m.syncDetail()
			}
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

	cardRows := 4
	bodyH := max(1, m.height-1-cardRows)
	railW := clamp(m.width/4, 24, 36)
	m.navCard.SetSize(railW, bodyH)
	m.detailCard.SetSize(max(1, m.width-railW), bodyH)
	m.sessionCard.SetSize(m.width, cardRows)
	m.footer.SetSize(m.width, 1)

	screen := rooms.RailFooterStack(m.width, m.height, railW, cardRows, m.navCard, m.detailCard, m.sessionCard, m.footer)
	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)
	surf.Draw(0, 0, screen)
	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func (m *model) syncDetail() {
	if len(m.items) == 0 {
		m.detailText.SetText("No sections")
		m.sessionText.SetText("cursor=0 status=empty")
		return
	}
	it := m.items[m.cursor]
	m.detailCard.SetTitle(it.title)
	m.detailCard.SetMeta("detail-view reference bento")
	m.detailText.SetText(fmt.Sprintf("section: %s\nkind: %s\nnotes: %s", it.title, it.kind, it.meta))
	m.sessionText.SetText(fmt.Sprintf("selected=%s\nindex=%d/%d", it.title, m.cursor+1, len(m.items)))
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
