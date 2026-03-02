// starter-app is the BentoTUI home screen.
// Wordmark · accented input block · kbd hints · tip · status bar.
// Run with: go run ./cmd/starter-app
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
	// OpenMsg must reach the manager even when no dialog is currently open.
	case dialog.OpenMsg:
		updated, cmd := m.dialogs.Update(msg)
		m.dialogs = updated.(*dialog.Manager)
		return m, cmd

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
	t := theme.CurrentTheme()
	canvasColor := lipgloss.Color(t.Surface.Canvas)

	// Always enter alt screen from frame 1 — before WindowSizeMsg arrives.
	if m.width == 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvasColor
		return v
	}

	bodyH := max(0, m.height-1) // minus status bar row

	// ── Surface: every cell is explicitly painted — no ANSI whitespace resets ──
	// Fill paints each cell with the canvas background color first, then we
	// draw lipgloss-rendered blocks on top via Ultraviolet's StyledString draw.
	surf := surface.New(m.width, bodyH)
	surf.Fill(canvasColor)

	dim := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	bright := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Primary))

	// ── wordmark ──────────────────────────────────────────────────────────────
	wm := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Text.Accent)).
		Bold(true).
		Render(wordmark)
	wmW := lipgloss.Width(wm)
	wmH := strings.Count(wm, "\n") + 1

	// ── input block ───────────────────────────────────────────────────────────
	badges := dim.Render("add   panel   list   input   table   dialog")
	inputStr := viewString(m.inputBox.View())
	inner := lipgloss.JoinVertical(lipgloss.Left, inputStr, badges)
	block := lipgloss.NewStyle().
		Border(lipgloss.Border{Left: "┃"}, false, false, false, true).
		BorderForeground(lipgloss.Color(t.Border.Focus)).
		PaddingTop(1).PaddingBottom(1).
		PaddingLeft(2).PaddingRight(2).
		Render(inner)
	blockW := lipgloss.Width(block)
	blockH := lipgloss.Height(block)

	// ── kbd hints ─────────────────────────────────────────────────────────────
	kbdStr := dim.Render("tab ") + bright.Render("components") +
		dim.Render("   ⌘K ") + bright.Render("commands")
	kbdW := lipgloss.Width(kbdStr)

	// ── tip ───────────────────────────────────────────────────────────────────
	dot := lipgloss.NewStyle().Foreground(lipgloss.Color(t.State.Info)).Render("● Tip")
	tipStr := dot + dim.Render("  Run bento init to scaffold a new TUI app")
	tipW := lipgloss.Width(tipStr)

	// ── vertical centering ────────────────────────────────────────────────────
	// Layout: wordmark(6) + gap(2) + block(4) + gap(1) + kbd(1) + gap(1) + tip(1) = 16
	const contentH = 16
	topPad := max(0, (bodyH-contentH)/2)

	// Draw each element at its centered X, calculated Y position.
	// surface.Draw(x, y, content) — no whitespace padding needed.
	y := topPad

	// wordmark — centered horizontally
	surf.Draw(max(0, (m.width-wmW)/2), y, wm)
	y += wmH + 2

	// input block — centered horizontally
	surf.Draw(max(0, (m.width-blockW)/2), y, block)
	y += blockH + 1

	// kbd hints — right-aligned (2 cell margin from edge)
	surf.Draw(max(0, m.width-kbdW-2), y, kbdStr)
	y += 2

	// tip — centered
	surf.Draw(max(0, (m.width-tipW)/2), y, tipStr)

	// ── dialog overlay ────────────────────────────────────────────────────────
	if m.dialogs.IsOpen() {
		dlgStr := viewString(m.dialogs.View())
		surf.DrawCenter(dlgStr)
	}

	// ── status bar ────────────────────────────────────────────────────────────
	statusStr := viewString(m.statusBar.View())

	// Render the surface, append a newline + status bar row.
	// surface.Render() uses \r\n between lines (raw buffer output);
	// Bubble Tea normalises this correctly.
	screen := surf.Render() + "\r\n" + statusStr

	v := tea.NewView(screen)
	v.AltScreen = true
	v.BackgroundColor = canvasColor
	return v
}

// onThemeChange keeps the status bar in sync with the active theme.
func (m *model) onThemeChange(msg theme.ThemeChangedMsg) {
	m.statusBar.SetRight(fmt.Sprintf("theme: %s  %s", msg.Name, version))
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
