package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/layout"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/dialog"
	"github.com/cloudboy-jh/bentotui/ui/containers/panel"
	"github.com/cloudboy-jh/bentotui/ui/widgets"
)

func main() {
	page := NewStarterPage()

	m := bentotui.New(
		bentotui.WithTheme(theme.CurrentTheme()),
		bentotui.WithPages(
			bentotui.Page("starter", func() core.Page { return page }),
		),
		bentotui.WithCommands(
			bentotui.Command{Label: "Open dialog", Keybind: "/dialog", Action: func() tea.Msg { return openCustomDialogCmd()() }},
			bentotui.Command{Label: "Switch theme", Keybind: "/theme", Action: func() tea.Msg { return openThemePickerCmd()() }},
			bentotui.Command{Label: "Command palette", Keybind: "/command", Action: func() tea.Msg { return openCommandPaletteCmd()() }},
		),
		bentotui.WithFooterBar(true),
	)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func openCustomDialogCmd() tea.Cmd {
	return func() tea.Msg {
		return dialog.Open(dialog.Custom{
			DialogTitle: "Custom Dialog",
			Content:     widgets.NewText("This is a custom dialog.\nPress Enter or Esc to close."),
			Width:       62,
			Height:      9,
		})
	}
}

func openThemePickerCmd() tea.Cmd {
	return func() tea.Msg {
		return theme.OpenThemePicker()
	}
}

// StarterPage demonstrates the new layout/widgets architecture.
type StarterPage struct {
	layout       *layout.Split
	input        *widgets.Input
	eventsList   *widgets.List
	infoPanel    *panel.Model
	commandPanel *panel.Model
	statusPanel  *panel.Model
	eventsPanel  *panel.Model
	width        int
	height       int
}

func NewStarterPage() *StarterPage {
	// Create content using new widgets
	infoList := widgets.NewList(50)
	infoList.Append("Page: starter")
	infoList.Append("Theme: " + theme.CurrentThemeName())
	infoList.Append("Status: Ready")

	eventsList := widgets.NewList(20)
	eventsList.Append("Application started")
	// Add command hints
	eventsList.Append("Type /command, /theme, or /dialog")

	input := widgets.NewInput()
	input.SetValue("")

	// Create panels using OLD panel API with NEW widgets as content
	infoPanel := panel.New(
		panel.Title("Info"),
		panel.Content(infoList),
	)

	commandPanel := panel.New(
		panel.Title("Command"),
		panel.Content(input),
	)

	statusPanel := panel.New(
		panel.Title("Status"),
		panel.Content(widgets.NewText("All systems operational")),
	)

	eventsPanel := panel.New(
		panel.Title("Events"),
		panel.Content(eventsList),
	)

	// Build layout using OLD canvas-based API
	// Right column: Status (fixed 8 rows) + Events (flex remaining)
	rightColumn := layout.Vertical(
		layout.Fixed(8, statusPanel),
		layout.Flex(1, eventsPanel),
	)

	// Main layout: 3 columns with gutter
	mainLayout := layout.Horizontal(
		layout.Flex(1, infoPanel),
		layout.Flex(2, commandPanel),
		layout.Flex(1, rightColumn),
	).WithGutterColor(theme.CurrentTheme().Border.Subtle)

	return &StarterPage{
		layout:       mainLayout,
		input:        input,
		eventsList:   eventsList,
		infoPanel:    infoPanel,
		commandPanel: commandPanel,
		statusPanel:  statusPanel,
		eventsPanel:  eventsPanel,
	}
}

func (p *StarterPage) Init() tea.Cmd {
	// Initialize layout and focus the input
	return tea.Batch(
		p.layout.Init(),
		p.input.Focus(),
	)
}

func (p *StarterPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case theme.ThemeChangedMsg:
		// Apply theme to all panels and widgets
		p.infoPanel.SetTheme(msg.Theme)
		p.commandPanel.SetTheme(msg.Theme)
		p.statusPanel.SetTheme(msg.Theme)
		p.eventsPanel.SetTheme(msg.Theme)
		p.input.SetTheme(msg.Theme)
		return p, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return p, tea.Quit
		}
		// Check if this is a printable character or special key we want to handle
		switch msg.String() {
		case "enter":
			// Check input for commands before passing to layout
			if cmd := p.checkInputCommand(); cmd != nil {
				return p, cmd
			}
			return p.layout.Update(msg)
		case "backspace", "delete", "left", "right", "home", "end":
			// Pass these to input directly
			updated, cmd := p.input.Update(msg)
			if inp, ok := updated.(*widgets.Input); ok {
				p.input = inp
			}
			return p, cmd
		default:
			// Pass all other keys (including "/") to input
			updated, cmd := p.input.Update(msg)
			if inp, ok := updated.(*widgets.Input); ok {
				p.input = inp
			}
			return p, cmd
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.layout.SetSize(p.width, p.height)
	}

	return p.layout.Update(msg)
}

func (p *StarterPage) checkInputCommand() tea.Cmd {
	text := p.input.Value()
	p.input.SetValue("") // Clear input after checking
	switch text {
	case "/command":
		return openCommandPaletteCmd()
	case "/theme":
		return openThemePickerCmd()
	case "/dialog":
		return openCustomDialogCmd()
	}
	return nil
}

func openCommandPaletteCmd() tea.Cmd {
	return func() tea.Msg {
		// Create a command palette with registered commands
		commands := []dialog.Command{
			{Label: "Open dialog", Group: "Suggested", Keybind: "/dialog", Action: func() tea.Msg { return openCustomDialogCmd()() }},
			{Label: "Switch theme", Group: "System", Keybind: "/theme", Action: func() tea.Msg { return openThemePickerCmd()() }},
		}
		palette := dialog.NewCommandPalette(commands)
		return dialog.Open(palette)
	}
}

func (p *StarterPage) View() tea.View {
	return p.layout.View()
}

func (p *StarterPage) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.layout.SetSize(width, height)
}

func (p *StarterPage) GetSize() (int, int) {
	return p.width, p.height
}

func (p *StarterPage) Title() string {
	return "BentoTUI Starter"
}

// Ensure StarterPage implements required interfaces
var _ core.Page = (*StarterPage)(nil)
var _ core.Component = (*StarterPage)(nil)
var _ core.Sizeable = (*StarterPage)(nil)
