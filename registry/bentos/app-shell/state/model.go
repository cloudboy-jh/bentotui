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

type FocusOwner string

const (
	FocusLeft   FocusOwner = "left"
	FocusCenter FocusOwner = "center"
	FocusRight  FocusOwner = "right"
	FocusFooter FocusOwner = "footer"
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
	showKeymap  bool
	stressStep  int
	status      string

	focusOwner FocusOwner

	navPanel    *panel.Model
	canvasPanel *panel.Model
	diagPanel   *panel.Model
	navText     *textBlock
	canvasText  *textBlock
	diagText    *textBlock
	footer      *bar.Model

	themeOrder []string
	themeIdx   int

	checks  []scenarios.Check
	metrics map[string]string
}

func NewModel() *Model {
	navTxt := &textBlock{}
	canvasTxt := &textBlock{}
	diagTxt := &textBlock{}

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
		diagText:    diagTxt,
		themeOrder:  names,
		themeIdx:    idx,
		status:      "ready",
		scenarioIdx: 0,
		presetIdx:   1,
		focusOwner:  FocusCenter,
		showKeymap:  false,
		checks:      nil,
		metrics:     map[string]string{},
	}

	m.navPanel = panel.New(panel.Title("Scenarios"), panel.Content(navTxt), panel.Elevated())
	m.canvasPanel = panel.New(panel.Title("Validation Canvas"), panel.Content(canvasTxt))
	m.diagPanel = panel.New(panel.Title("Diagnostics"), panel.Content(diagTxt), panel.Elevated())
	m.footer = bar.New(bar.FooterAnchored(), bar.Left("validation bento"), bar.Cards(ui.FooterCards()...), bar.CompactCards())
	m.syncFocusStyles()
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
		case "j", "down":
			m.selectScenario(m.scenarioIdx + 1)
		case "k", "up":
			m.selectScenario(m.scenarioIdx - 1)
		case "l", "right":
			m.shiftPreset(1)
		case "h", "left":
			m.shiftPreset(-1)
		case "t":
			m.shiftTheme(1)
		case "d":
			m.paintDebug = !m.paintDebug
			m.status = ternary(m.paintDebug, "paint debug enabled", "paint debug disabled")
		case "s":
			m.snapshot = !m.snapshot
			m.status = ternary(m.snapshot, "snapshot mode enabled", "snapshot mode disabled")
		case "m":
			m.showKeymap = !m.showKeymap
			m.status = ternary(m.showKeymap, "keymap visible", "keymap hidden")
		case "]", "tab":
			m.shiftFocus(1)
		case "[", "shift+tab":
			m.shiftFocus(-1)
		case "r":
			m.stressStep++
			m.status = fmt.Sprintf("stress step -> %d", m.stressStep)
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
	m.syncDiagnostics(bodyH)
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
	if m.width < 120 {
		return rooms.Sidebar(m.width, bodyH, navW, m.navPanel, m.canvasPanel)
	}

	innerW := max(1, m.width-navW)
	diagW := clamp(innerW/3, 30, 44)
	main := rooms.DrawerRight(innerW, bodyH, diagW, m.canvasPanel, m.diagPanel, rooms.WithGutter(1), rooms.WithDivider("subtle"))
	return rooms.Sidebar(m.width, bodyH, navW, m.navPanel, rooms.Static(main))
}

func (m *Model) syncScenarioText(bodyH int) {
	if len(m.scenarios) == 0 {
		m.canvasText.SetText("no scenarios")
		return
	}
	preset := viewportPresets[m.presetIdx]
	s := m.scenarios[m.scenarioIdx]
	ctx := scenarios.Context{
		Width:      min(max(32, m.width-44), preset.Width),
		Height:     min(max(10, bodyH-4), preset.Height),
		Viewport:   preset,
		PaintDebug: m.paintDebug,
		Snapshot:   m.snapshot,
		FocusOwner: string(m.focusOwner),
		StressStep: m.stressStep,
	}

	r := s.Run(ctx)
	m.checks = append([]scenarios.Check{}, r.Checks...)
	m.checks = append(m.checks, validateCanvasFrame(r.Canvas, ctx.Width, ctx.Height)...)
	m.metrics = r.Metrics

	header := ui.CanvasHeader(m.scenarioIdx+1, len(m.scenarios), s.Title, s.Description, m.width < 120)
	header = append(header, r.Canvas)

	if m.width < 120 {
		header = append(header, "", m.compactDiagnosticsLine())
	}
	m.canvasText.SetText(strings.Join(header, "\n"))
}

func (m *Model) syncDiagnostics(bodyH int) {
	if len(m.scenarios) == 0 {
		m.diagText.SetText("no diagnostics")
		return
	}
	p := viewportPresets[m.presetIdx]
	s := m.scenarios[m.scenarioIdx]
	d := ui.DiagnosticsInput{
		TerminalW:     m.width,
		TerminalH:     m.height,
		BodyH:         bodyH,
		Viewport:      p,
		ThemeName:     theme.CurrentThemeName(),
		ScenarioID:    s.ID,
		FocusOwner:    string(m.focusOwner),
		PaintDebug:    m.paintDebug,
		Snapshot:      m.snapshot,
		ShowKeymap:    m.showKeymap,
		Status:        m.status,
		ContrastScore: contrastScore(theme.CurrentTheme()),
		Checks:        m.checks,
		Metrics:       m.metrics,
	}
	m.diagText.SetText(ui.DiagnosticsText(d))
}

func (m *Model) syncFooterLine() {
	p := viewportPresets[m.presetIdx]
	pass, warn, fail := summarizeChecks(m.checks)
	m.footer.SetLeft(fmt.Sprintf("%s | %s | p:%d w:%d f:%d", m.scenarios[m.scenarioIdx].ID, p.Name, pass, warn, fail))
	m.footer.SetRight(fmt.Sprintf("theme:%s snapshot:%t focus:%s", theme.CurrentThemeName(), m.snapshot, m.focusOwner))
}

func (m *Model) compactDiagnosticsLine() string {
	p := viewportPresets[m.presetIdx]
	pass, warn, fail := summarizeChecks(m.checks)
	return fmt.Sprintf("diag: %s | %s(%dx%d) | %s | p:%d w:%d f:%d | focus:%s | %s",
		m.scenarios[m.scenarioIdx].ID,
		p.Name,
		p.Width,
		p.Height,
		theme.CurrentThemeName(),
		pass,
		warn,
		fail,
		m.focusOwner,
		m.status,
	)
}

func (m *Model) syncNavText() {
	m.navText.SetText(ui.SelectorText(m.scenarios, m.scenarioIdx, string(m.focusOwner)))
}

func (m *Model) syncFocusStyles() {
	m.navPanel.Blur()
	m.canvasPanel.Blur()
	m.diagPanel.Blur()
	if m.focusOwner == FocusLeft {
		m.navPanel.Focus()
	}
	if m.focusOwner == FocusCenter {
		m.canvasPanel.Focus()
	}
	if m.focusOwner == FocusRight {
		m.diagPanel.Focus()
	}
}

func (m *Model) shiftFocus(step int) {
	order := []FocusOwner{FocusLeft, FocusCenter, FocusRight, FocusFooter}
	idx := 0
	for i, v := range order {
		if v == m.focusOwner {
			idx = i
			break
		}
	}
	idx = (idx + step + len(order)) % len(order)
	m.focusOwner = order[idx]
	m.syncFocusStyles()
	m.status = "focus -> " + string(m.focusOwner)
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

func contrastScore(t theme.Theme) string {
	delta := func(a, b string) float64 {
		ar, ag, ab := parseHex(a)
		br, bg, bb := parseHex(b)
		la := luminance(ar, ag, ab)
		lb := luminance(br, bg, bb)
		if la > lb {
			return la - lb
		}
		return lb - la
	}
	vals := []float64{
		delta(t.Surface.Panel, t.Surface.Canvas),
		delta(t.Surface.Elevated, t.Surface.Panel),
		delta(t.Surface.Interactive, t.Surface.Panel),
		delta(t.Selection.BG, t.Surface.Canvas),
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	avg := sum / float64(len(vals))
	grade := "pass"
	if avg < 0.04 {
		grade = "warn"
	}
	if avg < 0.02 {
		grade = "fail"
	}
	return fmt.Sprintf("%s %.3f", grade, avg)
}

func parseHex(s string) (float64, float64, float64) {
	if len(s) != 7 || s[0] != '#' {
		return 0, 0, 0
	}
	r, _ := strconv.ParseInt(s[1:3], 16, 64)
	g, _ := strconv.ParseInt(s[3:5], 16, 64)
	b, _ := strconv.ParseInt(s[5:7], 16, 64)
	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0
}

func luminance(r, g, b float64) float64 {
	return 0.299*r + 0.587*g + 0.114*b
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
		checks = append(checks, scenarios.Check{
			Name:   "canvas-line-count",
			Level:  scenarios.CheckWarn,
			Detail: fmt.Sprintf("expected %d lines, got %d", height, len(lines)),
		})
	}
	for i, line := range lines {
		if lipgloss.Width(line) != width {
			checks = append(checks, scenarios.Check{
				Name:   "canvas-row-width",
				Level:  scenarios.CheckFail,
				Detail: fmt.Sprintf("line %d width mismatch: got %d want %d", i, lipgloss.Width(line), width),
			})
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
