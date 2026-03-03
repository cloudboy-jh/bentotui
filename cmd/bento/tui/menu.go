package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var menuItems = []string{
	"Initialize New Project",
	"Add Components",
	"Run Doctor",
	"Quit",
}

func (a *App) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if a.state.MenuSelection > 0 {
			a.state.MenuSelection--
		}
	case "down", "j":
		if a.state.MenuSelection < len(menuItems)-1 {
			a.state.MenuSelection++
		}
	case "enter":
		switch a.state.MenuSelection {
		case 0:
			a.state.CurrentView = ViewInitForm
			a.state.AppName = ""
			a.state.Module = ""
			a.state.FormFocus = 0
			a.state.ClearLog()
		case 1:
			a.state.CurrentView = ViewComponentList
			a.state.ComponentCursor = 0
			a.state.ClearLog()
		case 2:
			a.state.CurrentView = ViewDoctor
			a.startDoctor()
			if len(a.state.DoctorResults) > 0 {
				a.revealNextDoctorCheck()
				return a, a.doctorTickCmd()
			}
		case 3:
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a *App) renderMenu(height int) string {
	_ = height
	lines := make([]string, 0, len(menuItems)+2)
	lines = append(lines, "Select an action:")
	lines = append(lines, "")

	// Menu items
	for i, item := range menuItems {
		prefix := "  "
		if i == a.state.MenuSelection {
			prefix = "> "
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoFrameFG))
		if i == a.state.MenuSelection {
			style = style.Bold(true).Foreground(lipgloss.Color(bentoAccent))
		}

		line := style.Render(prefix + item)
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color(bentoMuted)).Render("↑/↓ move  enter select  q quit"))
	return strings.Join(lines, "\n")
}
