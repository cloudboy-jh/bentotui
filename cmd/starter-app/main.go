// starter-app is the BentoTUI home screen.
// Wordmark · accented input block · kbd hints · tip · status bar.
// Run with: go run ./cmd/starter-app
package main

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/bar"
	"github.com/cloudboy-jh/bentotui/registry/bricks/dialog"
	"github.com/cloudboy-jh/bentotui/registry/bricks/input"
	"github.com/cloudboy-jh/bentotui/registry/bricks/surface"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
	"github.com/cloudboy-jh/bentotui/theme"
)

const version = "v0.5.0"

// wordmark is large ASCII art rendered centered in the upper body.
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

// ── model ─────────────────────────────────────────────────────────────────────

type model struct {
	theme     theme.Theme
	inputBox  *input.Model
	footerBar *bar.Model
	dialogs   *dialog.Manager
	width     int
	height    int
	inputW    int // cached input block width — set on WindowSizeMsg
}

func newModel() *model {
	t := theme.CurrentTheme()
	inp := input.New()
	inp.SetPlaceholder(`Ask anything… /theme  /dialog`)
	inp.SetTheme(t)
	foot := bar.New(
		bar.FooterAnchored(),
		bar.Left("~ bentotui:main"),
		bar.Cards(
			bar.Card{Command: "enter", Label: "submit", Variant: bar.CardPrimary, Enabled: true, Priority: 3},
			bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted, Enabled: true, Priority: 2},
		),
		bar.CompactCards(),
		bar.WithTheme(t),
	)
	return &model{
		theme:     t,
		inputBox:  inp,
		footerBar: foot,
		dialogs:   dialog.New(),
	}
}

func (m *model) Init() tea.Cmd {
	return m.inputBox.Focus()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Dialog manager gets first shot — it owns esc/enter while open.
	if m.dialogs.IsOpen() {
		updated, cmd := m.dialogs.Update(msg)
		m.dialogs = updated.(*dialog.Manager)
		if tc, ok := msg.(theme.ThemeChangedMsg); ok {
			m.onThemeChange(tc)
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	// OpenMsg must reach the manager even when no dialog is currently open.
	case dialog.OpenMsg:
		updated, cmd := m.dialogs.Update(msg)
		m.dialogs = updated.(*dialog.Manager)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
		updated, cmd := m.inputBox.Update(msg)
		m.inputBox = updated.(*input.Model)
		return m, cmd
	}
	return m, nil
}

func (m *model) View() tea.View {
	t := m.theme
	canvasColor := t.Background()

	// Always enter alt screen from frame 1 — before WindowSizeMsg arrives.
	if m.width == 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvasColor
		return v
	}

	dim := lipgloss.NewStyle().Foreground(t.TextMuted())
	bright := lipgloss.NewStyle().Foreground(t.Text())

	wm := lipgloss.NewStyle().
		Foreground(t.TextAccent()).
		Bold(true).
		Render(wordmark)

	inputBlockW := m.inputW
	if inputBlockW == 0 {
		inputBlockW = clamp(m.width*6/10, 50, 90)
	}
	contentW := max(1, inputBlockW-5)
	inputStr := viewString(m.inputBox.View())

	mkRow := func(fg color.Color, content string) string {
		return lipgloss.NewStyle().
			Background(t.InputBG()).
			Foreground(fg).
			PaddingLeft(2).PaddingRight(2).
			Width(contentW).
			Render(content)
	}
	blankRow := lipgloss.NewStyle().Background(t.InputBG()).Width(contentW + 4).Render(" ")
	inner := lipgloss.JoinVertical(lipgloss.Left,
		blankRow,
		mkRow(t.InputFG(), inputStr),
		mkRow(t.TextMuted(), "add   card   list   input   table   dialog"),
		blankRow,
	)
	block := lipgloss.NewStyle().
		Background(t.InputBG()).
		Border(lipgloss.Border{Left: "┃"}, false, false, false, true).
		BorderForeground(t.BorderFocus()).
		Width(inputBlockW - 1).
		Render(inner)

	kbdStr := dim.Render("tab ") + bright.Render("bricks") +
		dim.Render("   ⌘K ") + bright.Render("commands")

	dot := lipgloss.NewStyle().Foreground(t.Info()).Render("● Tip")
	tipStr := dot + dim.Render("  Build pages with rooms and compose with bricks")

	body := rooms.RenderFunc(func(width, height int) string {
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

	screen := rooms.Focus(m.width, m.height, body, m.footerBar)
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

// onThemeChange keeps the status bar in sync with the active theme.
func (m *model) onThemeChange(msg theme.ThemeChangedMsg) {
	if msg.Theme == nil {
		return
	}
	m.theme = msg.Theme
	m.inputBox.SetTheme(m.theme)
	m.footerBar.SetTheme(m.theme)
}

// ── commands ──────────────────────────────────────────────────────────────────

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

// ── helpers ───────────────────────────────────────────────────────────────────

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
