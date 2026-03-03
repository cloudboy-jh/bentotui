package tui

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/cmd/bento/logic"
)

// doctorTickMsg is sent to animate the doctor checks.
type doctorTickMsg struct{}

func (a *App) startDoctor() {
	a.state.DoctorRunning = true
	a.state.DoctorIndex = 0
	a.state.DoctorResults = make([]DoctorCheck, 0)
	a.state.AddLog("Running doctor checks...")
	report := logic.RunDoctor()
	for _, r := range report.Results {
		a.state.DoctorResults = append(a.state.DoctorResults, DoctorCheck{
			Label: r.Label,
			Pass:  r.Pass,
			Note:  r.Note,
			Shown: false,
		})
	}
}

func (a *App) handleDoctorKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "r":
		// Retry
		a.startDoctor()
		return a, a.doctorTickCmd()
	}

	if !a.state.DoctorRunning && a.state.DoctorIndex >= len(a.state.DoctorResults) {
		// All done, allow navigation
		return a, nil
	}

	return a, nil
}

func (a *App) doctorTickCmd() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return doctorTickMsg{}
	})
}

func (a *App) renderDoctor(height int) string {
	_ = height

	lines := make([]string, 0, len(a.state.DoctorResults)+5)

	// Title
	lines = append(lines, "  Doctor Checks:")
	lines = append(lines, "")

	for _, check := range a.state.DoctorResults {
		if check.Shown {
			var icon, color string
			if check.Pass {
				icon = "✓"
				color = "#9cd67a"
			} else {
				icon = "✗"
				color = "#ff7e9f"
			}

			iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
			labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoFrameFG))

			line := "  " + iconStyle.Render("["+icon+"]") + " " + labelStyle.Render(check.Label)
			if !check.Pass && check.Note != "" {
				noteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoMuted))
				line += " - " + noteStyle.Render(check.Note)
			}

			lines = append(lines, line)
		}
	}

	// Check if we're done
	if a.state.DoctorRunning && a.state.DoctorIndex >= len(a.state.DoctorResults) {
		a.state.DoctorRunning = false

		// Log summary
		allPass := true
		for _, r := range a.state.DoctorResults {
			if !r.Pass {
				allPass = false
				break
			}
		}

		if allPass {
			a.state.AddLog("All checks passed!")
		} else {
			a.state.AddLog("Some checks failed.")
		}
	}

	// Show spinner if running
	if a.state.DoctorRunning {
		lines = append(lines, "")
		lines = append(lines, "  Checking...")
	}

	lines = append(lines, "")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoMuted))
	if a.state.DoctorRunning {
		lines = append(lines, footerStyle.Render("  esc back"))
	} else {
		lines = append(lines, footerStyle.Render("  r retry  esc back"))
	}
	return strings.Join(lines, "\n")
}

func (a *App) revealNextDoctorCheck() {
	if a.state.DoctorIndex >= len(a.state.DoctorResults) {
		a.state.DoctorRunning = false
		allPass := true
		for _, r := range a.state.DoctorResults {
			if !r.Pass {
				allPass = false
				break
			}
		}
		if allPass {
			a.state.AddLog("All checks passed!")
		} else {
			a.state.AddLog("Some checks failed.")
		}
		return
	}

	check := a.state.DoctorResults[a.state.DoctorIndex]
	a.state.DoctorResults[a.state.DoctorIndex].Shown = true
	a.state.DoctorIndex++
	icon := "✓"
	if !check.Pass {
		icon = "✗"
	}
	a.state.AddLog("[" + icon + "] " + check.Label)

	if a.state.DoctorIndex >= len(a.state.DoctorResults) {
		a.state.DoctorRunning = false
		allPass := true
		for _, r := range a.state.DoctorResults {
			if !r.Pass {
				allPass = false
				break
			}
		}
		if allPass {
			a.state.AddLog("All checks passed!")
		} else {
			a.state.AddLog("Some checks failed.")
		}
	}
}
