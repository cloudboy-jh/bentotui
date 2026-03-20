package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/registry/bricks/surface"
	"github.com/cloudboy-jh/bentotui/registry/recipes/vimstatus"
	"github.com/cloudboy-jh/bentotui/registry/rooms"
	"github.com/cloudboy-jh/bentotui/theme"
)

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

type model struct {
	theme     theme.Theme
	width     int
	height    int
	body      rooms.Sizable
	vimFooter *vimstatus.Model
	cfg       vimstatus.Config
	modes     []string
	modeIdx   int
}

func newModel() *model {
	t := theme.CurrentTheme()
	cfg := vimstatus.Config{
		Mode:      "NORMAL",
		Branch:    "main",
		Context:   "registry/bentos/vimstatus-demo/main.go",
		Position:  "49:30",
		Scroll:    "27%",
		ShowClock: true,
	}
	vf := vimstatus.New(t)
	vf.SetConfig(cfg)

	body := rooms.RenderFunc(func(width, height int) string {
		return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, "VimStatus Demo")
	})

	return &model{
		theme:     t,
		body:      body,
		vimFooter: vf,
		cfg:       cfg,
		modes:     []string{"NORMAL", "INSERT", "VISUAL", "COMMAND"},
	}
}

func (m *model) Init() tea.Cmd {
	return m.vimFooter.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.body.SetSize(msg.Width, max(1, msg.Height-1))
		m.vimFooter.SetSize(msg.Width, 1)

	case theme.ThemeChangedMsg:
		if msg.Theme != nil {
			m.theme = msg.Theme
			m.vimFooter.SetTheme(msg.Theme)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "i":
			m.modeIdx = (m.modeIdx + 1) % len(m.modes)
			m.cfg.Mode = m.modes[m.modeIdx]
			m.vimFooter.SetConfig(m.cfg)
		}
	}

	u, cmd := m.vimFooter.Update(msg)
	if next, ok := u.(*vimstatus.Model); ok {
		m.vimFooter = next
	}
	return m, cmd
}

func (m *model) View() tea.View {
	canvas := m.theme.Background()
	if m.width <= 0 || m.height <= 0 {
		v := tea.NewView("")
		v.AltScreen = true
		v.BackgroundColor = canvas
		return v
	}

	screen := rooms.Focus(m.width, m.height, m.body, m.vimFooter)
	surf := surface.New(m.width, m.height)
	surf.Fill(canvas)
	surf.Draw(0, 0, screen)

	v := tea.NewView(surf.Render())
	v.AltScreen = true
	v.BackgroundColor = canvas
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
