// starter-app is the BentoTUI home screen.
// Wordmark · accented input block · kbd hints · tip · status bar.
// Run with: go run ./cmd/starter-app
package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bar"
	"github.com/cloudboy-jh/bentotui/registry/dialog"
	"github.com/cloudboy-jh/bentotui/registry/input"
	"github.com/cloudboy-jh/bentotui/theme"
)

const version = "v0.2.0"

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
	inputBox  *input.Model
	statusBar *bar.Model
	dialogs   *dialog.Manager
	width     int
	height    int
}

func newModel() *model {
	inp := input.New()
	inp.SetPlaceholder(`Ask anything… /theme  /dialog`)
	sb := bar.New(
		bar.Left("~  bentotui:main"),
		bar.Right(fmt.Sprintf("theme: %s  %s", theme.CurrentThemeName(), version)),
	)
	return &model{
		inputBox:  inp,
		statusBar: sb,
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusBar.SetSize(msg.Width, 1)
		m.dialogs.SetSize(msg.Width, max(0, msg.Height-1))
		inputW := clamp(m.width*6/10, 50, 90)
		m.inputBox.SetSize(inputW-5, 1)
		return m, nil

	case theme.ThemeChangedMsg:
		m.onThemeChange(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+t":
			return m, openThemePicker()
		case "ctrl+k", "tab":
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
	if m.width == 0 {
		return tea.NewView("")
	}
	t := theme.CurrentTheme()

	// ── wordmark ──────────────────────────────────────────────────────────────
	wm := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text.Accent)).
		Bold(true).
		Render(wordmark)
	wmLines := strings.Split(
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, wm), "\n")

	// ── input block ───────────────────────────────────────────────────────────
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	badges := dim.Render("add   panel   list   input   table   dialog")
	inputStr := viewString(m.inputBox.View())
	inner := lipgloss.JoinVertical(lipgloss.Left, inputStr, badges)
	block := lipgloss.NewStyle().
		Border(lipgloss.Border{Left: "┃"}, false, false, false, true).
		BorderForeground(lipgloss.Color(t.Border.Focus)).
		PaddingTop(1).PaddingBottom(1).
		PaddingLeft(2).PaddingRight(2).
		Render(inner)
	blockLines := strings.Split(
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, block), "\n")

	// ── kbd hints ─────────────────────────────────────────────────────────────
	bright := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Primary))
	kbdStr := dim.Render("tab ") + bright.Render("components") +
		dim.Render("   ⌘K ") + bright.Render("commands")
	kbdLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Right, kbdStr+"  ")

	// ── tip ───────────────────────────────────────────────────────────────────
	dot := lipgloss.NewStyle().Foreground(lipgloss.Color(t.State.Info)).Render("● Tip")
	tip := lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
		dot+dim.Render("  Run bento init to scaffold a new TUI app"))

	// ── assemble body ─────────────────────────────────────────────────────────
	// Content rows: wordmark(6) + 2 blank + inputBlock(4) + 1 blank + kbd(1) + 1 blank + tip(1) = 16
	const contentH = 16
	bodyH := max(0, m.height-1) // minus status bar
	topPad := max(0, (bodyH-contentH)/2)
	botPad := max(0, bodyH-contentH-topPad)

	rows := make([]string, 0, bodyH+2)
	for i := 0; i < topPad; i++ {
		rows = append(rows, "")
	}
	rows = append(rows, wmLines...)
	rows = append(rows, "", "")
	rows = append(rows, blockLines...)
	rows = append(rows, "")
	rows = append(rows, kbdLine)
	rows = append(rows, "")
	rows = append(rows, tip)
	for i := 0; i < botPad; i++ {
		rows = append(rows, "")
	}
	bodyStr := strings.Join(rows, "\n")

	// ── dialog overlay ────────────────────────────────────────────────────────
	if m.dialogs.IsOpen() {
		dlgStr := viewString(m.dialogs.View())
		bodyStr = lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
			lipgloss.PlaceVertical(bodyH, lipgloss.Center, dlgStr))
	}

	// ── status bar ────────────────────────────────────────────────────────────
	statusStr := viewString(m.statusBar.View())
	screen := lipgloss.JoinVertical(lipgloss.Top, bodyStr, statusStr)

	v := tea.NewView(screen)
	v.AltScreen = true
	v.BackgroundColor = lipgloss.Color(t.Surface.Canvas)
	return v
}

// onThemeChange keeps the status bar in sync with the active theme.
func (m *model) onThemeChange(msg theme.ThemeChangedMsg) {
	m.statusBar.SetRight(fmt.Sprintf("theme: %s  %s", msg.Name, version))
}

// ── commands ──────────────────────────────────────────────────────────────────

func openThemePicker() tea.Cmd {
	return func() tea.Msg {
		return dialog.Open(dialog.Custom{
			DialogTitle: "Themes",
			Content:     dialog.NewThemePicker(),
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
