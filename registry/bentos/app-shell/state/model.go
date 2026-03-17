package state

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bentos/app-shell/scenarios"
	"github.com/cloudboy-jh/bentotui/registry/bentos/app-shell/ui"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"
	elevatedcard "github.com/cloudboy-jh/bentotui/registry/bricks/elevated-card"
	"github.com/cloudboy-jh/bentotui/registry/bricks/panel"
	"github.com/cloudboy-jh/bentotui/registry/bricks/surface"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
	"github.com/cloudboy-jh/bentotui/theme"
)

var viewportPresets = []scenarios.Viewport{
	{Name: "narrow", Width: 80, Height: 24},
	{Name: "medium", Width: 100, Height: 30},
	{Name: "wide", Width: 140, Height: 42},
}

type Model struct {
	width       int
	height      int
	scenarios   []scenarios.Definition
	scenarioIdx int
	presetIdx   int
	paintDebug  bool
	snapshot    bool
	status      string

	navPanel   *panel.Model
	canvasCard *elevatedcard.Model
	centerDeck *centerDeck
	navText    *textBlock
	footerText *textBlock
	footer     *bar.Model
	footerCard *elevatedcard.Model

	themeOrder []string
	themeIdx   int

	checks []scenarios.Check
}

func NewModel() *Model {
	navTxt := &textBlock{}
	footerTxt := &textBlock{}
	deck := newCenterDeck()

	names := theme.AvailableStableThemes()
	if len(names) == 0 {
		names = theme.AvailableThemes()
	}
	cur := theme.CurrentThemeName()
	idx := 0
	for i, n := range names {
		if n == cur {
			idx = i
			break
		}
	}

	m := &Model{
		scenarios:   scenarios.All(),
		navText:     navTxt,
		centerDeck:  deck,
		footerText:  footerTxt,
		themeOrder:  names,
		themeIdx:    idx,
		status:      "ready",
		scenarioIdx: 0,
		presetIdx:   1,
	}

	m.navPanel = panel.New(panel.Title("Scenarios"), panel.Content(navTxt), panel.Elevated())
	m.canvasCard = elevatedcard.New(elevatedcard.Title("UI Sandbox"), elevatedcard.Content(deck), elevatedcard.Inset(1))
	m.footerCard = elevatedcard.New(elevatedcard.Title("Session"), elevatedcard.Content(footerTxt), elevatedcard.CardVariant(elevatedcard.VariantDense), elevatedcard.Inset(1))
	m.navPanel.Focus()
	m.footer = bar.New(
		bar.FooterAnchored(),
		bar.AnchoredCardStyleMode(bar.AnchoredCardStyleMixed),
		bar.Left("validation bento"),
		bar.Cards(ui.FooterCards()...),
		bar.CompactCards(),
	)
	m.syncNavText()
	return m
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.SetSize(msg.Width, 1)
		if cardH := m.footerCardRows(); cardH > 0 {
			m.footerCard.SetSize(msg.Width, cardH)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.selectScenario(m.scenarioIdx - 1)
		case "down":
			m.selectScenario(m.scenarioIdx + 1)
		case "left":
			m.shiftPreset(-1)
		case "right":
			m.shiftPreset(1)
		case "t":
			m.shiftTheme(1)
		case "d":
			m.paintDebug = !m.paintDebug
			m.status = ternary(m.paintDebug, "paint debug enabled", "paint debug disabled")
		case "s":
			m.snapshot = !m.snapshot
			m.status = ternary(m.snapshot, "snapshot mode enabled", "snapshot mode disabled")
		default:
			if n, err := strconv.Atoi(msg.String()); err == nil {
				m.selectScenario(n - 1)
			}
		}
		return m, nil
	}
	return m, nil
}

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	canvas := lipgloss.Color(t.Surface.Canvas)

	if m.width == 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvas
		return v
	}

	m.syncNavText()
	cardH := m.footerCardRows()
	bodyH := max(1, m.height-1-cardH)
	m.syncScenarioText(bodyH)
	m.syncFooterLine()
	m.syncFooterCard()

	m.footer.SetSize(m.width, 1)
	if cardH > 0 {
		m.footerCard.SetSize(m.width, cardH)
	}
	navW := clamp(m.width/5, 24, 34)
	screen := rooms.RailFooterStack(m.width, m.height, navW, cardH, m.navPanel, m.canvasCard, m.footerCard, m.footer)

	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)
	surf.Draw(0, 0, screen)

	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func (m *Model) syncScenarioText(bodyH int) {
	if len(m.scenarios) == 0 {
		m.centerDeck.SetOutput("Scenario Output", "no scenarios")
		m.centerDeck.SetChecks("checks: none")
		m.centerDeck.SetMetrics("context unavailable")
		return
	}
	preset := viewportPresets[m.presetIdx]
	s := m.scenarios[m.scenarioIdx]
	ctx := scenarios.Context{
		Width:      min(max(32, m.width-40), preset.Width),
		Height:     min(max(10, bodyH-6), preset.Height),
		Viewport:   preset,
		PaintDebug: m.paintDebug,
		Snapshot:   m.snapshot,
		FocusOwner: "center",
		StressStep: m.themeIdx*10 + m.presetIdx,
	}

	r := s.Run(ctx)
	m.checks = append([]scenarios.Check{}, r.Checks...)

	title := fmt.Sprintf("[%d/%d] %s", m.scenarioIdx+1, len(m.scenarios), s.Title)
	m.canvasCard.SetTitle(title)
	m.canvasCard.SetMeta(s.Description)
	m.canvasCard.SetFooter(m.inlineSummary())
	body := compactBodyLines(r.Canvas)
	if m.paintDebug {
		body = append(body, strings.Repeat("-", max(12, min(ctx.Width-4, 72))))
	}
	m.centerDeck.SetOutput("Scenario Output", strings.Join(body, "\n"))

	m.centerDeck.SetChecks(formatChecksForCard(m.checks))
	m.centerDeck.SetMetrics(fmt.Sprintf("scenario=%s\nviewport=%s(%dx%d)\ntheme=%s\nstatus=%s",
		s.ID,
		preset.Name,
		preset.Width,
		preset.Height,
		theme.CurrentThemeName(),
		m.status,
	))
}

func (m *Model) inlineSummary() string {
	p := viewportPresets[m.presetIdx]
	return fmt.Sprintf("summary: %s | viewport=%s(%dx%d) | theme=%s | %s",
		m.scenarios[m.scenarioIdx].ID,
		p.Name,
		p.Width,
		p.Height,
		theme.CurrentThemeName(),
		m.status,
	)
}

func (m *Model) syncFooterLine() {
	p := viewportPresets[m.presetIdx]
	m.footer.SetLeft(fmt.Sprintf("%s | %s", m.scenarios[m.scenarioIdx].ID, p.Name))
	m.footer.SetRight(fmt.Sprintf("theme:%s snapshot:%t", theme.CurrentThemeName(), m.snapshot))
}

func (m *Model) syncFooterCard() {
	p := viewportPresets[m.presetIdx]
	pass, warn, fail := summarizeChecks(m.checks)
	m.footerText.SetText(fmt.Sprintf("scenario=%s  status=%s\nviewport=%s(%dx%d)  theme=%s\nchecks: pass=%d warn=%d fail=%d",
		m.scenarios[m.scenarioIdx].ID,
		m.status,
		p.Name,
		p.Width,
		p.Height,
		theme.CurrentThemeName(),
		pass,
		warn,
		fail,
	))
}

func formatChecksForCard(checks []scenarios.Check) string {
	if len(checks) == 0 {
		return "No checks yet"
	}
	lines := make([]string, 0, len(checks))
	for i := 0; i < len(checks) && i < 6; i++ {
		c := checks[i]
		prefix := "ok"
		switch c.Level {
		case scenarios.CheckWarn:
			prefix = "warn"
		case scenarios.CheckFail:
			prefix = "fail"
		}
		detail := strings.TrimSpace(c.Detail)
		if detail == "" {
			detail = c.Name
		}
		lines = append(lines, fmt.Sprintf("%s %s", prefix, detail))
	}
	return strings.Join(lines, "\n")
}

func (m *Model) syncNavText() {
	m.navText.SetText(ui.SelectorText(m.scenarios, m.scenarioIdx))
}

func (m *Model) selectScenario(idx int) {
	if len(m.scenarios) == 0 {
		return
	}
	if idx < 0 {
		idx = len(m.scenarios) - 1
	}
	if idx >= len(m.scenarios) {
		idx = 0
	}
	m.scenarioIdx = idx
	m.status = "scenario -> " + m.scenarios[idx].ID
}

func (m *Model) shiftPreset(step int) {
	if len(viewportPresets) == 0 {
		return
	}
	m.presetIdx = (m.presetIdx + step + len(viewportPresets)) % len(viewportPresets)
	p := viewportPresets[m.presetIdx]
	m.status = fmt.Sprintf("viewport -> %s (%dx%d)", p.Name, p.Width, p.Height)
}

func (m *Model) shiftTheme(step int) {
	if len(m.themeOrder) == 0 {
		m.status = "theme registry empty"
		return
	}
	m.themeIdx = (m.themeIdx + step + len(m.themeOrder)) % len(m.themeOrder)
	name := m.themeOrder[m.themeIdx]
	if _, err := theme.SetTheme(name); err != nil {
		m.status = "theme error: " + err.Error()
		return
	}
	m.status = "theme -> " + name
}

func summarizeChecks(checks []scenarios.Check) (pass, warn, fail int) {
	for _, c := range checks {
		switch c.Level {
		case scenarios.CheckWarn:
			warn++
		case scenarios.CheckFail:
			fail++
		default:
			pass++
		}
	}
	return pass, warn, fail
}

func compactBodyLines(canvas string) []string {
	lines := strings.Split(canvas, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " ")
		if strings.TrimSpace(trimmed) == "" {
			continue
		}
		out = append(out, trimmed)
	}
	if len(out) == 0 {
		return []string{"(empty scenario output)"}
	}
	return out
}

func ternary[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const footerCardHeight = 4

func (m *Model) footerCardRows() int {
	if m.height <= 2 {
		return 0
	}
	h := footerCardHeight
	if h > m.height-2 {
		h = m.height - 2
	}
	if h < 1 {
		return 0
	}
	return h
}
