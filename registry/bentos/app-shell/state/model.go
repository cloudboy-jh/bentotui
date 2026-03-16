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

	navPanel    *panel.Model
	canvasPanel *panel.Model
	navText     *textBlock
	canvasText  *textBlock
	footer      *bar.Model

	themeOrder []string
	themeIdx   int

	checks []scenarios.Check
}

func NewModel() *Model {
	navTxt := &textBlock{}
	canvasTxt := &textBlock{}

	names := theme.AvailableThemes()
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
		canvasText:  canvasTxt,
		themeOrder:  names,
		themeIdx:    idx,
		status:      "ready",
		scenarioIdx: 0,
		presetIdx:   1,
	}

	m.navPanel = panel.New(panel.Title("Scenarios"), panel.Content(navTxt), panel.Elevated())
	m.canvasPanel = panel.New(panel.Title("Validation Canvas"), panel.Content(canvasTxt))
	m.navPanel.Focus()
	m.footer = bar.New(bar.FooterAnchored(), bar.Left("validation bento"), bar.Cards(ui.FooterCards()...), bar.CompactCards())
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

	bodyH := max(1, m.height-1)
	m.syncNavText()
	m.syncScenarioText(bodyH)
	m.syncFooterLine()

	body := m.layoutBody(bodyH)
	screen := rooms.Focus(m.width, m.height, rooms.Static(body), m.footer)

	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)
	surf.Draw(0, 0, screen)

	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func (m *Model) layoutBody(bodyH int) string {
	if m.width < 84 {
		return rooms.VSplit(m.width, bodyH, m.navPanel, m.canvasPanel)
	}
	navW := clamp(m.width/5, 24, 34)
	return rooms.Sidebar(m.width, bodyH, navW, m.navPanel, m.canvasPanel)
}

func (m *Model) syncScenarioText(bodyH int) {
	if len(m.scenarios) == 0 {
		m.canvasText.SetText("no scenarios")
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
	m.checks = append(m.checks, validateCanvasFrame(r.Canvas, ctx.Width, ctx.Height)...)

	header := ui.CanvasHeader(m.scenarioIdx+1, len(m.scenarios), s.Title, s.Description, m.width < 84)
	header = append(header, r.Canvas, "", m.inlineSummary())
	m.canvasText.SetText(strings.Join(header, "\n"))
}

func (m *Model) inlineSummary() string {
	p := viewportPresets[m.presetIdx]
	pass, warn, fail := summarizeChecks(m.checks)
	return fmt.Sprintf("summary: viewport=%s(%dx%d) theme=%s checks=p:%d w:%d f:%d status=%s",
		p.Name, p.Width, p.Height, theme.CurrentThemeName(), pass, warn, fail, m.status)
}

func (m *Model) syncFooterLine() {
	p := viewportPresets[m.presetIdx]
	pass, warn, fail := summarizeChecks(m.checks)
	m.footer.SetLeft(fmt.Sprintf("%s | %s | p:%d w:%d f:%d", m.scenarios[m.scenarioIdx].ID, p.Name, pass, warn, fail))
	m.footer.SetRight(fmt.Sprintf("theme:%s snapshot:%t", theme.CurrentThemeName(), m.snapshot))
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

func validateCanvasFrame(canvas string, width, height int) []scenarios.Check {
	checks := []scenarios.Check{}
	if width <= 0 || height <= 0 {
		return append(checks, scenarios.Check{Name: "canvas-size", Level: scenarios.CheckFail, Detail: "invalid scenario canvas size"})
	}
	lines := strings.Split(canvas, "\n")
	if len(lines) != height {
		checks = append(checks, scenarios.Check{Name: "canvas-line-count", Level: scenarios.CheckWarn, Detail: fmt.Sprintf("expected %d lines, got %d", height, len(lines))})
	}
	for i, line := range lines {
		if lipgloss.Width(line) != width {
			checks = append(checks, scenarios.Check{Name: "canvas-row-width", Level: scenarios.CheckFail, Detail: fmt.Sprintf("line %d width mismatch: got %d want %d", i, lipgloss.Width(line), width)})
			break
		}
	}
	if len(checks) == 0 {
		checks = append(checks, scenarios.Check{Name: "canvas-row-width", Level: scenarios.CheckPass, Detail: "all scenario rows are width-exact"})
	}
	return checks
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
