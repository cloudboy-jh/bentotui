// starter-app is the bentotui component showcase.
// It demonstrates every registry component with live theme switching.
// Run with: go run ./cmd/starter-app
package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/layout"
	"github.com/cloudboy-jh/bentotui/registry/bar"
	"github.com/cloudboy-jh/bentotui/registry/dialog"
	"github.com/cloudboy-jh/bentotui/registry/input"
	"github.com/cloudboy-jh/bentotui/registry/list"
	"github.com/cloudboy-jh/bentotui/registry/panel"
	"github.com/cloudboy-jh/bentotui/registry/table"
	"github.com/cloudboy-jh/bentotui/registry/text"
	"github.com/cloudboy-jh/bentotui/theme"
)

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

// ── root model ────────────────────────────────────────────────────────────────

type model struct {
	// Layout
	root   *layout.Split
	header *bar.Model
	footer *bar.Model
	body   *layout.Split

	// Panels
	infoPanel   *panel.Model
	eventsPanel *panel.Model
	inputPanel  *panel.Model
	tablePanel  *panel.Model

	// Widgets
	eventsList *list.Model
	inputBox   *input.Model
	dataTable  *table.Model

	// Dialogs
	dialogs *dialog.Manager

	// Commands for the palette
	commands []dialog.Command

	width  int
	height int
}

const (
	headerH = 1
	footerH = 1
)

func newModel() *model {
	// Events log
	events := list.New(50)
	events.Append("bentotui starter-app")
	events.Append("ctrl+t or /theme → theme picker")
	events.Append("ctrl+p or /command → command palette")
	events.Append("ctrl+d or /dialog → sample dialog")

	// Input
	inp := input.New()
	inp.SetPlaceholder("Type here… or /theme, /command, /dialog")

	// Table
	tbl := table.New("Component", "Status", "Notes")
	tbl.AddRow("panel", "✓", "themed container")
	tbl.AddRow("bar", "✓", "header/footer row")
	tbl.AddRow("dialog", "✓", "modal overlay")
	tbl.AddRow("list", "✓", "plain-text log")
	tbl.AddRow("table", "✓", "header + rows")
	tbl.AddRow("text", "✓", "static label")
	tbl.AddRow("input", "✓", "text field")

	// Panels
	infoP := panel.New(
		panel.Title("bentotui"),
		panel.Content(text.New(infoText())),
	)
	eventsP := panel.New(
		panel.Title("Events"),
		panel.Content(events),
	)
	inputP := panel.New(
		panel.Title("Input"),
		panel.Content(inp),
		panel.Elevated(),
	)
	tableP := panel.New(
		panel.Title("Components"),
		panel.Content(tbl),
	)

	// Layout: left column (info + input), right column (table + events)
	left := layout.Vertical(
		layout.Flex(1, infoP),
		layout.Fixed(4, inputP),
	)
	right := layout.Vertical(
		layout.Flex(1, tableP),
		layout.Fixed(8, eventsP),
	)
	body := layout.Horizontal(
		layout.Flex(1, left),
		layout.Flex(2, right),
	)

	// Header + footer bars
	hdr := bar.New(
		bar.Left("bentotui showcase"),
		bar.Right(fmt.Sprintf("theme: %s", theme.CurrentThemeName())),
	)
	ftr := bar.New(
		bar.Cards(
			bar.Card{Command: "ctrl+t  /theme", Label: "theme", Enabled: true},
			bar.Card{Command: "ctrl+p  /command", Label: "palette", Enabled: true},
			bar.Card{Command: "ctrl+d  /dialog", Label: "dialog", Enabled: true},
			bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardDanger, Enabled: true},
		),
	)

	// Root layout: header / body / footer
	root := layout.Vertical(
		layout.Fixed(headerH, hdr),
		layout.Flex(1, body),
		layout.Fixed(footerH, ftr),
	)

	// Dialog manager + command palette commands
	commands := []dialog.Command{
		{Label: "Switch theme", Group: "App", Keybind: "ctrl+t", Action: func() tea.Msg {
			return dialog.Open(dialog.Custom{
				DialogTitle: "Themes",
				Content:     dialog.NewThemePicker(),
			})
		}},
		{Label: "Sample confirm", Group: "App", Action: func() tea.Msg {
			return dialog.Open(dialog.Confirm{
				DialogTitle: "Confirm",
				Message:     "This is a confirm dialog.\nPress Enter to confirm.",
			})
		}},
		{Label: "Quit", Group: "App", Keybind: "ctrl+c", Action: func() tea.Msg {
			return tea.Quit()
		}},
	}

	return &model{
		root:        root,
		header:      hdr,
		footer:      ftr,
		body:        body,
		infoPanel:   infoP,
		eventsPanel: eventsP,
		inputPanel:  inputP,
		tablePanel:  tableP,
		eventsList:  events,
		inputBox:    inp,
		dataTable:   tbl,
		dialogs:     dialog.New(),
		commands:    commands,
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.root.Init(),
		m.inputBox.Focus(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Dialog manager gets first shot at every message.
	if m.dialogs.IsOpen() {
		updated, cmd := m.dialogs.Update(msg)
		m.dialogs = updated.(*dialog.Manager)
		// Handle dialog-generated theme changes.
		if tc, ok := msg.(theme.ThemeChangedMsg); ok {
			m.onThemeChange(tc)
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.root.SetSize(m.width, m.height)
		// Propagate body size to dialog manager.
		bodyH := max(0, m.height-headerH-footerH)
		m.dialogs.SetSize(m.width, bodyH)
		return m, nil

	case theme.ThemeChangedMsg:
		m.onThemeChange(msg)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+t":
			m.eventsList.Append("opening theme picker")
			return m, func() tea.Msg {
				return dialog.Open(dialog.Custom{
					DialogTitle: "Themes",
					Content:     dialog.NewThemePicker(),
				})
			}
		case "ctrl+p":
			m.eventsList.Append("opening command palette")
			return m, func() tea.Msg {
				return dialog.Open(dialog.Custom{
					DialogTitle: "Commands",
					Content:     dialog.NewCommandPalette(m.commands),
				})
			}
		case "ctrl+d":
			m.eventsList.Append("opening dialog")
			return m, func() tea.Msg {
				return dialog.Open(dialog.Confirm{
					DialogTitle: "Hello",
					Message:     "This is a Confirm dialog.\nPress Enter or Esc.",
				})
			}
		case "enter":
			val := strings.TrimSpace(m.inputBox.Value())
			if val != "" {
				m.inputBox.SetValue("")
				if cmd := m.checkInputCommand(val); cmd != nil {
					return m, cmd
				}
				m.eventsList.Append("> " + val)
			}
			return m, nil
		}
		// Pass other keys to the input box.
		updated, cmd := m.inputBox.Update(msg)
		m.inputBox = updated.(*input.Model)
		return m, cmd
	}

	// Propagate to layout for window size and other system messages.
	updated, cmd := m.root.Update(msg)
	m.root = updated.(*layout.Split)
	return m, cmd
}

func (m *model) View() tea.View {
	t := theme.CurrentTheme()

	// Render body (panels).
	bodyStr := viewString(m.body.View())

	// If a dialog is open, overlay it centered on the body.
	if m.dialogs.IsOpen() {
		dialogStr := viewString(m.dialogs.View())
		bodyH := max(0, m.height-headerH-footerH)
		centered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center,
			lipgloss.PlaceVertical(bodyH, lipgloss.Center, dialogStr))
		bodyStr = centered
	}

	headerStr := viewString(m.header.View())
	footerStr := viewString(m.footer.View())

	// Paint the full screen: header / body / footer.
	screen := lipgloss.JoinVertical(lipgloss.Top, headerStr, bodyStr, footerStr)

	v := tea.NewView(screen)
	v.AltScreen = true
	v.BackgroundColor = lipgloss.Color(t.Surface.Canvas)
	return v
}

// onThemeChange updates gutter colors and header text when the theme changes.
func (m *model) onThemeChange(msg theme.ThemeChangedMsg) {
	m.body.SetGutterColor(msg.Theme.Border.Subtle)
	m.header.SetRight(fmt.Sprintf("theme: %s", msg.Name))
	m.eventsList.Append("theme → " + msg.Name)
}

// checkInputCommand routes slash commands typed into the input field.
// Returns a non-nil Cmd if the input was a slash command, nil otherwise.
func (m *model) checkInputCommand(val string) tea.Cmd {
	switch val {
	case "/theme":
		m.eventsList.Append("opening theme picker")
		return func() tea.Msg {
			return dialog.Open(dialog.Custom{
				DialogTitle: "Themes",
				Content:     dialog.NewThemePicker(),
			})
		}
	case "/command":
		m.eventsList.Append("opening command palette")
		return func() tea.Msg {
			return dialog.Open(dialog.Custom{
				DialogTitle: "Commands",
				Content:     dialog.NewCommandPalette(m.commands),
			})
		}
	case "/dialog":
		m.eventsList.Append("opening dialog")
		return func() tea.Msg {
			return dialog.Open(dialog.Confirm{
				DialogTitle: "Hello",
				Message:     "This is a Confirm dialog.\nPress Enter or Esc.",
			})
		}
	}
	return nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func infoText() string {
	return strings.Join([]string{
		"Copy-and-own registry",
		"for Go TUIs.",
		"",
		"Copy components into",
		"your project and own",
		"them completely.",
		"",
		"bento add panel",
		"bento add bar",
		"bento add dialog",
		"bento add list",
		"bento add table",
		"bento add input",
		"bento add text",
	}, "\n")
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
