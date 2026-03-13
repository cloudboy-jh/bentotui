package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/components/bar"
	"github.com/cloudboy-jh/bentotui/registry/components/dialog"
	"github.com/cloudboy-jh/bentotui/registry/components/input"
	"github.com/cloudboy-jh/bentotui/registry/components/surface"
	"github.com/cloudboy-jh/bentotui/registry/layouts"
	"github.com/cloudboy-jh/bentotui/theme"
)

const version = "v0.3.0"

const wordmark = "" +
	"██████╗ ███████╗███╗   ██╗████████╗ ██████╗ \n" +
	"██╔══██╗██╔════╝████╗  ██║╚══██╔══╝██╔═══██╗\n" +
	"██████╔╝█████╗  ██╔██╗ ██║   ██║   ██║   ██║\n" +
	"██╔══██╗██╔══╝  ██║╚██╗██║   ██║   ██║   ██║\n" +
	"██████╔╝███████╗██║ ╚████║   ██║   ╚██████╔╝\n" +
	"╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝    ╚═════╝ "

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

type model struct {
	inputBox  *input.Model
	topBar    *bar.Model
	metaBar   *bar.Model
	footerBar *bar.Model
	dialogs   *dialog.Manager
	width     int
	height    int
	inputW    int
}

func newModel() *model {
	inp := input.New()
	inp.SetPlaceholder(`Ask anything… /theme  /dialog`)
	top := bar.New(
		bar.RoleTopBar(),
		bar.StatusPill("LIVE"),
		bar.Left("bentotui home-screen"),
		bar.Right("context: examples"),
	)
	meta := bar.New(
		bar.RoleSubBar(),
		bar.Left("starter grammar: frame"),
		bar.Right(fmt.Sprintf("theme: %s", theme.CurrentThemeName())),
	)
	foot := bar.New(
		bar.FooterAnchored(),
		bar.Left("~ registry/bentos/home-screen"),
		bar.Cards(
			bar.Card{Command: "enter", Label: "submit", Variant: bar.CardPrimary, Enabled: true, Priority: 3},
			bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
	)
	return &model{inputBox: inp, topBar: top, metaBar: meta, footerBar: foot, dialogs: dialog.New()}
}

func (m *model) Init() tea.Cmd {
	return m.inputBox.Focus()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dialogs.IsOpen() {
		u, cmd := m.dialogs.Update(msg)
		m.dialogs = u.(*dialog.Manager)
		if tc, ok := msg.(theme.ThemeChangedMsg); ok {
			m.onThemeChange(tc)
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case dialog.OpenMsg:
		u, cmd := m.dialogs.Update(msg)
		m.dialogs = u.(*dialog.Manager)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.topBar.SetSize(msg.Width, 1)
		m.metaBar.SetSize(msg.Width, 1)
		m.footerBar.SetSize(msg.Width, 1)
		m.dialogs.SetSize(msg.Width, msg.Height)
		m.inputW = clamp(m.width*6/10, 50, 90)
		m.inputBox.SetSize(m.inputW-5, 1)
		return m, nil

	case theme.ThemeChangedMsg:
		m.onThemeChange(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			return m, nil
		case "enter":
			val := strings.TrimSpace(m.inputBox.Value())
			if val == "" {
				return m, nil
			}
			m.inputBox.SetValue("")
			switch val {
			case "/theme":
				return m, openThemePicker()
			case "/dialog":
				return m, openSampleDialog()
			}
			return m, nil
		}

		u, cmd := m.inputBox.Update(msg)
		m.inputBox = u.(*input.Model)
		return m, cmd
	}

	return m, nil
}

func (m *model) View() tea.View {
	t := theme.CurrentTheme()
	canvasColor := lipgloss.Color(t.Surface.Canvas)

	if m.width == 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvasColor
		return v
	}

	dim := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	bright := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Primary))

	wm := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text.Accent)).
		Bold(true).
		Render(wordmark)
	inputBlockW := m.inputW
	if inputBlockW == 0 {
		inputBlockW = clamp(m.width*6/10, 50, 90)
	}
	contentW := max(1, inputBlockW-5)
	inputStr := viewString(m.inputBox.View())

	mkRow := func(fg, content string) string {
		return lipgloss.NewStyle().
			Background(lipgloss.Color(t.Input.BG)).
			Foreground(lipgloss.Color(fg)).
			PaddingLeft(2).PaddingRight(2).
			Width(contentW).
			Render(content)
	}

	blankRow := lipgloss.NewStyle().Background(lipgloss.Color(t.Input.BG)).Width(contentW + 4).Render(" ")
	inner := lipgloss.JoinVertical(lipgloss.Left,
		blankRow,
		mkRow(t.Input.FG, inputStr),
		mkRow(t.Text.Muted, "add   panel   list   input   table   dialog"),
		blankRow,
	)

	block := lipgloss.NewStyle().
		Background(lipgloss.Color(t.Input.BG)).
		Border(lipgloss.Border{Left: "┃"}, false, false, false, true).
		BorderForeground(lipgloss.Color(t.Border.Focus)).
		Width(inputBlockW - 1).
		Render(inner)

	kbdStr := dim.Render("tab ") + bright.Render("components") +
		dim.Render("   ⌘K ") + bright.Render("commands")

	dot := lipgloss.NewStyle().Foreground(lipgloss.Color(t.State.Info)).Render("● Tip")
	tipStr := dot + dim.Render("  Run bento init to scaffold a new TUI app")
	body := layouts.RenderFunc(func(width, height int) string {
		center := func(s string) string {
			return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(s)
		}
		right := func(s string) string {
			line := lipgloss.NewStyle().Width(max(1, width-2)).Align(lipgloss.Right).Render(s)
			if width > 1 {
				return " " + line
			}
			return line
		}
		stack := strings.Join([]string{
			center(wm),
			"",
			center(block),
			"",
			right(kbdStr),
			"",
			center(tipStr),
		}, "\n")
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, stack)
	})

	screen := layouts.Frame(m.width, m.height, m.topBar, m.metaBar, body, m.footerBar)
	surf := surface.New(m.width, m.height)
	surf.Fill(canvasColor)
	surf.Draw(0, 0, screen)
	if m.dialogs.IsOpen() {
		surf.DrawCenter(viewString(m.dialogs.View()))
	}

	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvasColor
	return v
}

func (m *model) onThemeChange(msg theme.ThemeChangedMsg) {
	m.metaBar.SetRight(fmt.Sprintf("theme: %s", msg.Name))
}

func openThemePicker() tea.Cmd {
	return func() tea.Msg {
		h := len(theme.AvailableThemes()) + 8
		return dialog.Open(dialog.Custom{
			DialogTitle: "Themes",
			Content:     dialog.NewThemePicker(),
			Width:       44,
			Height:      h,
		})
	}
}

func openSampleDialog() tea.Cmd {
	return func() tea.Msg {
		return dialog.Open(dialog.Confirm{
			DialogTitle: "Hello from BentoTUI",
			Message:     "This is a Confirm dialog.\nPress Enter to confirm or Esc to cancel.",
		})
	}
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
