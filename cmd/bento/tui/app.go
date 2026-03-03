package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/components/bar"
	"github.com/cloudboy-jh/bentotui/registry/components/list"
	"github.com/cloudboy-jh/bentotui/theme"
)

const (
	maxFrameWidth = 92
	minFrameWidth = 40
	logRows       = 4
	bentoFrameBG  = "#2f2826"
	bentoFrameFG  = "#f4ddca"
	bentoBorder   = "#a77b8f"
	bentoMuted    = "#d5b8a2"
	bentoAccent   = "#ff92b6"
)

// App is the root tea.Model for the bento TUI.
type App struct {
	state  *State
	header *bar.Model
	log    *list.Model
}

// NewApp creates a new TUI app.
func NewApp() *App {
	_, _ = theme.PreviewTheme("bento-rose")

	return &App{
		state: NewState(),
		header: bar.New(
			bar.Left("🍱 bento"),
			bar.Right("v0.2.0"),
		),
		log: list.New(300),
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.state.Width = msg.Width
		a.state.Height = msg.Height
		frameW := clamp(msg.Width-2, minFrameWidth, maxFrameWidth)
		a.header.SetSize(frameW, 1)
		a.log.SetSize(frameW-2, logRows)
		return a, nil

	case tea.KeyMsg:
		return a.handleKey(msg)

	case doctorTickMsg:
		// Continue doctor animation
		if a.state.DoctorRunning && a.state.DoctorIndex < len(a.state.DoctorResults) {
			a.revealNextDoctorCheck()
			return a, a.doctorTickCmd()
		}
		return a, nil
	}

	return a, nil
}

func (a *App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return a, tea.Quit
	case "q":
		if a.state.CurrentView == ViewMenu {
			return a, tea.Quit
		}
		return a.goBack()
	case "esc":
		return a.goBack()
	}

	// Route to current view
	switch a.state.CurrentView {
	case ViewMenu:
		return a.handleMenuKey(msg)
	case ViewInitForm:
		return a.handleInitFormKey(msg)
	case ViewComponentList:
		return a.handleComponentListKey(msg)
	case ViewDoctor:
		return a.handleDoctorKey(msg)
	}

	return a, nil
}

func (a *App) goBack() (tea.Model, tea.Cmd) {
	switch a.state.CurrentView {
	case ViewInitForm, ViewComponentList, ViewDoctor:
		a.state.CurrentView = ViewMenu
		a.state.ClearLog()
	}
	return a, nil
}

func (a *App) View() tea.View {
	if a.state.Width == 0 {
		return tea.NewView("loading bento...")
	}

	frameW := clamp(a.state.Width-2, minFrameWidth, maxFrameWidth)
	a.header.SetSize(frameW, 1)
	a.log.SetSize(frameW-2, logRows)
	contentHeight := 12

	// Build content based on current view
	var content string
	switch a.state.CurrentView {
	case ViewMenu:
		content = a.renderMenu(contentHeight)
	case ViewInitForm:
		content = a.renderInitForm(contentHeight)
	case ViewComponentList:
		content = a.renderComponentList(contentHeight)
	case ViewDoctor:
		content = a.renderDoctor(contentHeight)
	}

	a.log.Clear()
	for _, line := range a.state.LogLines {
		a.log.Append(line)
	}

	logView := viewString(a.log.View())
	logBlock := lipgloss.NewStyle().
		Background(lipgloss.Color("#3a312f")).
		Foreground(lipgloss.Color(bentoFrameFG)).
		Width(frameW).
		Render(logView)

	// Combine all sections
	body := strings.Join([]string{
		viewString(a.header.View()),
		"",
		content,
		lipgloss.NewStyle().Foreground(lipgloss.Color(bentoBorder)).Render(strings.Repeat("─", frameW)),
		logBlock,
	}, "\n")

	frame := lipgloss.NewStyle().
		Width(frameW).
		Background(lipgloss.Color(bentoFrameBG)).
		Foreground(lipgloss.Color(bentoFrameFG)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(bentoBorder)).
		Padding(0, 1).
		Render(body)

	return tea.NewView(frame)
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
