package state

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bentos/app-shell/ui"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"
	"github.com/cloudboy-jh/bentotui/registry/bricks/dialog"
	"github.com/cloudboy-jh/bentotui/registry/bricks/surface"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
	"github.com/cloudboy-jh/bentotui/theme"
)

type setThemeMsg struct{ Name string }
type setSectionMsg struct{ Index int }
type toggleCompactMsg struct{}
type pulseProgressMsg struct{}

type Model struct {
	width  int
	height int

	sections   []string
	sectionIdx int
	queueIdx   int
	compact    bool
	progress   float64
	status     string

	centerDeck *centerDeck
	footer     *bar.Model
	dialogs    *dialog.Manager

	themeOrder []string
	themeIdx   int
}

func NewModel() *Model {
	deck := newCenterDeck()

	themes := theme.AvailableThemes()
	cur := theme.CurrentThemeName()
	themeIdx := 0
	for i, n := range themes {
		if n == cur {
			themeIdx = i
			break
		}
	}

	m := &Model{
		sections:   []string{"Overview", "Services", "Queue", "Progress"},
		centerDeck: deck,
		themeOrder: themes,
		themeIdx:   themeIdx,
		sectionIdx: 0,
		queueIdx:   0,
		compact:    true,
		progress:   0.62,
		status:     "ready",
		dialogs:    dialog.New(),
	}

	m.footer = bar.New(
		bar.FooterAnchored(),
		bar.Left("workspace"),
		bar.Cards(ui.FooterCards()...),
		bar.CompactCards(),
	)
	m.syncAll()
	return m
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.footer.SetSize(msg.Width, 1)
		m.dialogs.SetSize(msg.Width, msg.Height)
		return m, nil

	case dialog.OpenMsg, dialog.CloseMsg:
		u, cmd := m.dialogs.Update(msg)
		m.dialogs = u.(*dialog.Manager)
		return m, cmd

	case setThemeMsg:
		m.applyTheme(msg.Name)
		return m, nil

	case setSectionMsg:
		m.setSection(msg.Index)
		return m, nil

	case toggleCompactMsg:
		m.compact = !m.compact
		m.status = ternary(m.compact, "table compact", "table comfortable")
		return m, nil

	case pulseProgressMsg:
		m.progress += 0.07
		if m.progress > 1 {
			m.progress = 0.08
		}
		m.status = fmt.Sprintf("progress %.0f%%", m.progress*100)
		return m, nil

	case tea.KeyMsg:
		if m.dialogs.IsOpen() {
			u, cmd := m.dialogs.Update(msg)
			m.dialogs = u.(*dialog.Manager)
			return m, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.setSection(m.sectionIdx - 1)
		case "down":
			m.setSection(m.sectionIdx + 1)
		case "left":
			m.queueIdx = max(0, m.queueIdx-1)
			m.status = "queue cursor <-"
		case "right":
			m.queueIdx = min(3, m.queueIdx+1)
			m.status = "queue cursor ->"
		case "enter":
			m.progress += 0.05
			if m.progress > 1 {
				m.progress = 0.1
			}
			m.status = fmt.Sprintf("pulse %.0f%%", m.progress*100)
		case "t":
			m.shiftTheme(1)
		case "c":
			m.compact = !m.compact
			m.status = ternary(m.compact, "table compact", "table comfortable")
		case "ctrl+k":
			return m, m.openPalette()
		default:
			if n, err := strconv.Atoi(msg.String()); err == nil {
				m.setSection(n - 1)
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

	m.syncAll()

	screen := rooms.Focus(m.width, m.height, m.centerDeck, m.footer)

	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)
	surf.Draw(0, 0, screen)
	if m.dialogs.IsOpen() {
		surf.DrawCenter(viewString(m.dialogs.View()))
	}

	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func (m *Model) syncAll() {
	if len(m.sections) == 0 {
		m.sections = []string{"Overview"}
		m.sectionIdx = 0
	}
	if m.sectionIdx < 0 {
		m.sectionIdx = 0
	}
	if m.sectionIdx >= len(m.sections) {
		m.sectionIdx = len(m.sections) - 1
	}

	m.centerDeck.SetActiveSection(m.sections[m.sectionIdx])
	m.centerDeck.SetQueueCursor(m.queueIdx)
	m.centerDeck.SetCompact(m.compact)
	m.centerDeck.SetProgress(m.progress, fmt.Sprintf("section: %s", strings.ToLower(m.sections[m.sectionIdx])))
	m.centerDeck.SetWorkspaceMeta(strings.ToLower(m.sections[m.sectionIdx]), m.compact)
	m.syncFooterLine()
}

func (m *Model) syncFooterLine() {
	left := fmt.Sprintf("%s | queue:%d", strings.ToLower(m.sections[m.sectionIdx]), m.queueIdx+1)
	right := fmt.Sprintf("theme:%s compact:%t %s", theme.CurrentThemeName(), m.compact, m.status)
	m.footer.SetLeft(left)
	m.footer.SetRight(right)
}

func (m *Model) setSection(idx int) {
	if len(m.sections) == 0 {
		return
	}
	if idx < 0 {
		idx = len(m.sections) - 1
	}
	if idx >= len(m.sections) {
		idx = 0
	}
	m.sectionIdx = idx
	m.status = "section -> " + strings.ToLower(m.sections[idx])
}

func (m *Model) shiftTheme(step int) {
	if len(m.themeOrder) == 0 {
		m.status = "theme registry empty"
		return
	}
	m.themeIdx = (m.themeIdx + step + len(m.themeOrder)) % len(m.themeOrder)
	m.applyTheme(m.themeOrder[m.themeIdx])
}

func (m *Model) applyTheme(name string) {
	if _, err := theme.SetTheme(name); err != nil {
		m.status = "theme error: " + err.Error()
		return
	}
	for i, n := range m.themeOrder {
		if n == name {
			m.themeIdx = i
			break
		}
	}
	m.status = "theme -> " + name
}

func (m *Model) openPalette() tea.Cmd {
	commands := make([]dialog.Command, 0, len(m.sections)+len(m.themeOrder)+4)
	for i, section := range m.sections {
		idx := i
		commands = append(commands, dialog.Command{
			Label:   "Go to " + section,
			Group:   "Navigate",
			Keybind: strconv.Itoa(i + 1),
			Action:  func() tea.Msg { return setSectionMsg{Index: idx} },
		})
	}
	commands = append(commands,
		dialog.Command{Label: "Toggle compact table", Group: "View", Keybind: "c", Action: func() tea.Msg { return toggleCompactMsg{} }},
		dialog.Command{Label: "Pulse progress", Group: "View", Keybind: "enter", Action: func() tea.Msg { return pulseProgressMsg{} }},
	)
	for _, name := range m.themeOrder {
		themeName := name
		commands = append(commands, dialog.Command{
			Label:   "Switch to " + themeName,
			Group:   "Theme",
			Keybind: "t",
			Action:  func() tea.Msg { return setThemeMsg{Name: themeName} },
		})
	}

	return func() tea.Msg {
		palette := dialog.NewCommandPalette(commands)
		h := clamp(min(24, m.height-4), 12, 24)
		return dialog.Open(dialog.Custom{
			DialogTitle: "Command Palette",
			Content:     palette,
			Width:       56,
			Height:      h,
		})
	}
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
