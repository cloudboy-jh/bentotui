// Brick: Package Manager:
// +-----------------------------------+
// | spinner Installing pkg   [====]   |
// +-----------------------------------+
// Sequential package install flow inspired by Bubble Tea's package-manager example.
// Copy this file into your project: bento add package-manager
package packagemanager

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

type installedPkgMsg string
type failedPkgMsg struct {
	pkg string
	err error
}

type installFunc func(string) tea.Cmd

type Model struct {
	packages   []string
	index      int
	width      int
	height     int
	spinner    spinner.Model
	progress   progress.Model
	done       bool
	err        error
	quitOnDone bool
	installer  installFunc
}

func New(packages []string) *Model {
	sp := spinner.New()
	pg := progress.New(
		progress.WithDefaultBlend(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	m := &Model{
		packages:   append([]string(nil), packages...),
		spinner:    sp,
		progress:   pg,
		width:      1,
		height:     1,
		quitOnDone: true,
		installer:  downloadAndInstall,
	}
	m.syncStyles()
	return m
}

func (m *Model) Init() tea.Cmd {
	if len(m.packages) == 0 {
		m.done = true
		return nil
	}
	return tea.Batch(m.installer(m.packages[m.index]), m.spinner.Tick)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case installedPkgMsg:
		if len(m.packages) == 0 {
			m.done = true
			return m, nil
		}
		if m.index >= len(m.packages)-1 {
			m.done = true
			if m.quitOnDone {
				return m, tea.Quit
			}
			return m, nil
		}
		m.index++
		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.packages)))
		return m, tea.Batch(progressCmd, m.installer(m.packages[m.index]))
	case failedPkgMsg:
		m.err = msg.err
		m.done = true
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) View() tea.View {
	m.syncStyles()
	n := len(m.packages)
	if m.err != nil {
		return tea.NewView(lipgloss.NewStyle().Render("Failed: " + m.err.Error()))
	}
	if n == 0 {
		return tea.NewView("No packages.")
	}
	if m.done {
		return tea.NewView(fmt.Sprintf("Done! Installed %d packages.", n))
	}

	w := lipgloss.Width(fmt.Sprintf("%d", n))
	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)
	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))
	t := theme.CurrentTheme()
	pkgName := lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Accent, t.Text.Primary))).Render(m.packages[m.index])
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Installing " + pkgName)
	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	line := spin + info + gap + prog + pkgCount
	return tea.NewView(lipgloss.NewStyle().Width(max(1, m.width)).Render(line))
}

func (m *Model) SetSize(width, height int) {
	m.width = max(1, width)
	m.height = max(1, height)
	barWidth := max(10, m.width/3)
	m.progress.SetWidth(barWidth)
}

func (m *Model) GetSize() (int, int) { return m.width, m.height }

func (m *Model) Done() bool { return m.done }

func (m *Model) Error() error { return m.err }

func (m *Model) SetInstaller(fn func(string) tea.Cmd) {
	if fn == nil {
		return
	}
	m.installer = fn
}

func (m *Model) SetQuitOnDone(v bool) { m.quitOnDone = v }

func (m *Model) syncStyles() {
	t := theme.CurrentTheme()
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Accent, t.Text.Primary)))
	m.progress.Full = '█'
	m.progress.Empty = '░'
	m.progress.FullColor = lipgloss.Color(pick(t.Selection.BG, t.Border.Focus))
	m.progress.EmptyColor = lipgloss.Color(pick(t.Border.Subtle, t.Border.Normal))
}

func downloadAndInstall(pkg string) tea.Cmd {
	d := time.Millisecond * time.Duration(rand.Intn(500)) //nolint:gosec
	return tea.Tick(d, func(time.Time) tea.Msg {
		return installedPkgMsg(pkg)
	})
}

func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
