package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/cmd/bento/logic"
)

func (a *App) handleInitFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if a.state.FormFocus > 0 {
			a.state.FormFocus--
		}
	case "down", "j", "tab":
		if a.state.FormFocus < 2 {
			a.state.FormFocus++
		}
	case "backspace":
		switch a.state.FormFocus {
		case 0:
			if len(a.state.AppName) > 0 {
				a.state.AppName = a.state.AppName[:len(a.state.AppName)-1]
			}
		case 1:
			if len(a.state.Module) > 0 {
				a.state.Module = a.state.Module[:len(a.state.Module)-1]
			}
		}
	case "enter":
		if a.state.FormFocus == 2 {
			// Submit form
			return a.submitInitForm()
		}
		// Move to next field
		if a.state.FormFocus < 2 {
			a.state.FormFocus++
		}
	default:
		// Add character to current field
		if len(msg.String()) == 1 {
			switch a.state.FormFocus {
			case 0:
				a.state.AppName += msg.String()
			case 1:
				a.state.Module += msg.String()
			}
		}
	}
	return a, nil
}

func (a *App) renderInitForm(height int) string {
	_ = height

	// Form fields
	fields := []struct {
		label   string
		value   string
		hint    string
		focused bool
	}{
		{"App Name", a.state.AppName, "my-bento-app", a.state.FormFocus == 0},
		{"Module", a.state.Module, "example.com/my-bento-app", a.state.FormFocus == 1},
		{"[ Submit ]", "", "", a.state.FormFocus == 2},
	}

	lines := make([]string, 0, 14)
	lines = append(lines, "Init project:")
	lines = append(lines, "")

	// Fields
	for i, f := range fields {
		if f.label == "[ Submit ]" {
			// Submit button
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoFrameFG))
			if f.focused {
				style = style.Bold(true).Foreground(lipgloss.Color(bentoAccent))
			}
			lines = append(lines, "")
			lines = append(lines, style.Render("  "+f.label))
		} else {
			// Input field
			labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoMuted))
			if f.focused {
				labelStyle = labelStyle.Bold(true).Foreground(lipgloss.Color(bentoFrameFG))
			}

			// Show value or hint
			displayValue := f.value
			if displayValue == "" && !f.focused {
				displayValue = f.hint
			}

			// Input box styling
			inputStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("#4a3a35")).
				Foreground(lipgloss.Color(bentoFrameFG)).
				Width(40)

			if f.value == "" && !f.focused {
				inputStyle = inputStyle.Foreground(lipgloss.Color(bentoMuted))
			}

			lines = append(lines, labelStyle.Render("  "+f.label+":"))
			lines = append(lines, "  "+inputStyle.Render(displayValue))
			lines = append(lines, "")
		}

		_ = i // suppress unused warning
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color(bentoMuted)).Render("tab/↑/↓ move  enter submit  esc back"))
	return strings.Join(lines, "\n")
}

func (a *App) submitInitForm() (tea.Model, tea.Cmd) {
	// Set defaults if empty
	appName := a.state.AppName
	if appName == "" {
		appName = "my-bento-app"
	}
	module := a.state.Module
	if module == "" {
		module = "example.com/" + appName
	}

	a.state.AddLog("Creating project: " + appName)

	cfg := logic.ProjectConfig{
		AppName: appName,
		Module:  module,
	}

	created, err := logic.ScaffoldProject(cfg)
	if err != nil {
		a.state.AddLog("Error: " + err.Error())
		return a, nil
	}

	for _, f := range created {
		a.state.AddLog("Created: " + f)
	}

	a.state.AddLog("Done.")

	return a, nil
}
