package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"
	"github.com/cloudboy-jh/bentotui/registry/bricks/card"
	"github.com/cloudboy-jh/bentotui/registry/bricks/filepicker"
	"github.com/cloudboy-jh/bentotui/registry/bricks/list"
	packagemanager "github.com/cloudboy-jh/bentotui/registry/bricks/package-manager"
	"github.com/cloudboy-jh/bentotui/registry/bricks/surface"
	"github.com/cloudboy-jh/bentotui/registry/bricks/table"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
	"github.com/cloudboy-jh/bentotui/theme"
)

type model struct {
	theme  theme.Theme
	width  int
	height int
	active int

	// theme cycling
	themeOrder []string
	themeIdx   int

	footer *bar.Model

	list       *list.Model
	table      *table.Model
	filepicker *filepicker.Model
	pkg        *packagemanager.Model

	listCard *card.Model
	tblCard  *card.Model
	fpCard   *card.Model
	pkgCard  *card.Model
}

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func newModel() *model {
	// Default to catppuccin-mocha for high-contrast visual testing baseline.
	if _, err := theme.SetTheme("catppuccin-mocha"); err != nil {
		// fall back to whatever is registered
		_ = err
	}

	l := list.New(200)
	l.AppendSection("SERVICES")
	l.AppendRow(list.Row{Primary: "api", Secondary: "healthy", Tone: list.ToneSuccess, RightStat: "36ms"})
	l.AppendRow(list.Row{Primary: "cache", Secondary: "degraded", Tone: list.ToneWarn, RightStat: "112ms"})
	l.AppendRow(list.Row{Primary: "queue", Secondary: "healthy", Tone: list.ToneSuccess, RightStat: "42ms"})

	tb := table.New("SERVICE", "STATUS", "LATENCY", "ERR%")
	tb.SetVisualStyle(table.VisualGrid)
	tb.SetColumnAlign(2, table.AlignRight)
	tb.SetColumnAlign(3, table.AlignRight)
	tb.AddRow("api", "healthy", "38ms", "0.1")
	tb.AddRow("workers", "healthy", "55ms", "0.0")
	tb.AddRow("cache", "degraded", "112ms", "1.7")
	tb.AddRow("queue", "healthy", "47ms", "0.2")

	fp := filepicker.New(".")
	fp.SetAllowDirectories(true)
	fp.SetAllowFiles(true)

	pm := packagemanager.New([]string{"bubbletea", "bubbles", "lipgloss", "bentotui"})
	pm.SetQuitOnDone(false)

	themes := theme.AvailableThemes()
	themeIdx := 0
	for i, n := range themes {
		if n == theme.CurrentThemeName() {
			themeIdx = i
			break
		}
	}

	m := &model{
		theme:      theme.CurrentTheme(),
		active:     0,
		themeOrder: themes,
		themeIdx:   themeIdx,
		list:       l,
		table:      tb,
		filepicker: fp,
		pkg:        pm,
		footer: bar.New(
			bar.FooterAnchored(),
			bar.Left("focus: list"),
			bar.Cards(
				bar.Card{Command: "arrows", Label: "focus", Variant: bar.CardPrimary, Enabled: true, Priority: 5},
				bar.Card{Command: "t", Label: "theme", Variant: bar.CardNormal, Enabled: true, Priority: 4},
				bar.Card{Command: "T", Label: "theme←", Variant: bar.CardNormal, Enabled: true, Priority: 3},
				bar.Card{Command: "enter", Label: "select", Variant: bar.CardNormal, Enabled: true, Priority: 2},
				bar.Card{Command: "q", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 1},
			),
			bar.CompactCards(),
		),
	}

	m.listCard = card.New(card.Title("List"), card.Content(m.list))
	m.tblCard = card.New(card.Title("Table"), card.Content(m.table))
	m.fpCard = card.New(card.Title("File Picker"), card.Content(m.filepicker))
	m.pkgCard = card.New(card.Title("Package Manager"), card.Content(m.pkg))
	m.applyTheme()
	m.syncFocus()
	return m
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.filepicker.Init(), m.pkg.Init())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		// panel focus navigation
		case "left":
			if m.active%2 == 1 {
				m.active--
				m.syncFocus()
			}
			return m, nil
		case "right":
			if m.active%2 == 0 {
				m.active++
				m.syncFocus()
			}
			return m, nil
		case "up":
			if m.active >= 2 {
				m.active -= 2
				m.syncFocus()
			}
			return m, nil
		case "down":
			if m.active <= 1 {
				m.active += 2
				m.syncFocus()
			}
			return m, nil

		// theme cycling
		case "t":
			m.shiftTheme(1)
			return m, nil
		case "T":
			m.shiftTheme(-1)
			return m, nil
		}

		return m, m.updateActive(msg)
	}

	// Non-key messages: always route to filepicker and package-manager
	// since they need ticker/init messages regardless of focus.
	var cmds []tea.Cmd
	if cmd := m.updateFilepicker(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	if cmd := m.updatePackageManager(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	t := m.theme
	canvas := t.Background()

	if m.width <= 0 || m.height <= 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvas
		return v
	}

	bodyH := max(1, m.height-1)
	m.layout()
	m.syncFooter()

	body := rooms.Dashboard2x2(m.width, bodyH, m.listCard, m.tblCard, m.fpCard, m.pkgCard)
	screen := rooms.Focus(m.width, m.height, rooms.Static(body), m.footer)

	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)
	surf.Draw(0, 0, screen)
	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func (m *model) syncFooter() {
	m.footer.SetLeft("focus: " + m.activeLabel())
	m.footer.SetRight("theme: " + m.theme.Name())

	if m.filepicker.Status() != "" {
		m.fpCard.SetMeta(m.filepicker.Status())
	}
	if m.pkg.Done() {
		m.pkgCard.SetMeta("done")
	}
}

func (m *model) layout() {
	if m.width <= 0 || m.height <= 0 {
		return
	}
	bodyH := max(1, m.height-1)
	rowH := max(1, bodyH/2)
	bottomH := max(1, bodyH-rowH)
	leftW := max(1, m.width/2)
	rightW := max(1, m.width-leftW)

	m.listCard.SetSize(leftW, rowH)
	m.tblCard.SetSize(rightW, rowH)
	m.fpCard.SetSize(leftW, bottomH)
	m.pkgCard.SetSize(rightW, bottomH)
}

func (m *model) syncFocus() {

	m.list.Blur()
	m.table.Blur()
	m.filepicker.Blur()

	switch m.active {
	case 0:
		m.list.Focus()

	case 1:
		m.table.Focus()

	case 2:
		m.filepicker.Focus()

	case 3:

	}
}

func (m *model) shiftTheme(step int) {
	if len(m.themeOrder) == 0 {
		return
	}
	m.themeIdx = (m.themeIdx + step + len(m.themeOrder)) % len(m.themeOrder)
	t, err := theme.SetTheme(m.themeOrder[m.themeIdx])
	if err != nil {
		// skip invalid theme, advance again
		m.themeIdx = (m.themeIdx + step + len(m.themeOrder)) % len(m.themeOrder)
		t, err = theme.SetTheme(m.themeOrder[m.themeIdx])
		if err != nil {
			return
		}
	}
	m.theme = t
	m.applyTheme()
}

func (m *model) applyTheme() {
	m.footer.SetTheme(m.theme)
	m.listCard.SetTheme(m.theme)
	m.tblCard.SetTheme(m.theme)
	m.fpCard.SetTheme(m.theme)
	m.pkgCard.SetTheme(m.theme)
}

func (m *model) updateActive(msg tea.Msg) tea.Cmd {
	switch m.active {
	case 0:
		u, cmd := m.list.Update(msg)
		if next, ok := u.(*list.Model); ok {
			m.list = next
		}
		return cmd
	case 1:
		u, cmd := m.table.Update(msg)
		if next, ok := u.(*table.Model); ok {
			m.table = next
		}
		return cmd
	case 2:
		return m.updateFilepicker(msg)
	case 3:
		return m.updatePackageManager(msg)
	}
	return nil
}

func (m *model) updateFilepicker(msg tea.Msg) tea.Cmd {
	u, cmd := m.filepicker.Update(msg)
	if next, ok := u.(*filepicker.Model); ok {
		m.filepicker = next
	}
	return cmd
}

func (m *model) updatePackageManager(msg tea.Msg) tea.Cmd {
	u, cmd := m.pkg.Update(msg)
	if next, ok := u.(*packagemanager.Model); ok {
		m.pkg = next
	}
	return cmd
}

func (m *model) activeLabel() string {
	switch m.active {
	case 0:
		return "list"
	case 1:
		return "table"
	case 2:
		return "filepicker"
	case 3:
		return "package-manager"
	default:
		return "unknown"
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
