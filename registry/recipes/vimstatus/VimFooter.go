package vimstatus

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/registry/bricks/badge"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

type Config struct {
	Mode      string
	Branch    string
	Context   string
	Position  string
	Scroll    string
	ShowClock bool
}

type Model struct {
	theme     theme.Theme
	width     int
	height    int
	cfg       Config
	clock     string
	modeBadge *badge.Model
}

type clockTickMsg struct {
	when time.Time
}

func New(t theme.Theme) *Model {
	if t == nil {
		t = theme.CurrentTheme()
	}
	b := badge.New("")
	b.SetTheme(t)
	b.SetBold(true)
	m := &Model{theme: t, height: 1, modeBadge: b}
	m.syncModeBadge()
	return m
}

func (m *Model) SetConfig(cfg Config) {
	m.cfg = cfg
	if m.cfg.ShowClock {
		m.clock = time.Now().Format("15:04")
	}
	m.syncModeBadge()
}

func (m *Model) SetTheme(t theme.Theme) {
	m.theme = t
	if m.modeBadge != nil {
		m.modeBadge.SetTheme(m.activeTheme())
	}
}

func (m *Model) SetSize(width, _ int) {
	m.width = width
	m.height = 1
}

func (m *Model) View() tea.View {
	t := m.activeTheme()
	left := m.renderLeft(t)
	right := m.renderRight(t)
	line := composeAlignedLine(m.width, left, right)
	return tea.NewView(styles.Row(t.FooterBG(), t.FooterFG(), m.width, styles.ClipANSI(line, m.width)))
}

func (m *Model) Init() tea.Cmd {
	if !m.cfg.ShowClock {
		return nil
	}
	if m.clock == "" {
		m.clock = time.Now().Format("15:04")
	}
	return tickClock()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case clockTickMsg:
		if !m.cfg.ShowClock {
			return m, nil
		}
		m.clock = msg.when.Format("15:04")
		return m, tickClock()
	}
	return m, nil
}

func (m *Model) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

func (m *Model) renderLeft(t theme.Theme) string {
	parts := make([]string, 0, 3)
	if pill := m.renderModeBadge(); pill != "" {
		parts = append(parts, pill)
	}
	if v := strings.TrimSpace(m.cfg.Branch); v != "" {
		parts = append(parts, lipgloss.NewStyle().Foreground(t.FooterMuted()).Render(v))
	}
	if v := strings.TrimSpace(m.cfg.Context); v != "" {
		parts = append(parts, lipgloss.NewStyle().Foreground(t.FooterFG()).Render(v))
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderRight(t theme.Theme) string {
	parts := make([]string, 0, 3)
	muted := lipgloss.NewStyle().Foreground(t.FooterMuted())
	if v := strings.TrimSpace(m.cfg.Position); v != "" {
		parts = append(parts, muted.Render(v))
	}
	if v := strings.TrimSpace(m.cfg.Scroll); v != "" {
		parts = append(parts, muted.Render(v))
	}
	if m.cfg.ShowClock {
		clockText := m.clock
		if clockText == "" {
			clockText = time.Now().Format("15:04")
		}
		parts = append(parts, lipgloss.NewStyle().Bold(true).Foreground(t.FooterFG()).Render(clockText))
	}
	return strings.Join(parts, " ")
}

func (m *Model) renderModeBadge() string {
	mode := strings.ToUpper(strings.TrimSpace(m.cfg.Mode))
	if mode == "" {
		return ""
	}
	if m.modeBadge == nil {
		return ""
	}
	return viewString(m.modeBadge.View())
}

func (m *Model) syncModeBadge() {
	if m.modeBadge == nil {
		m.modeBadge = badge.New("")
	}
	m.modeBadge.SetTheme(m.activeTheme())

	mode := strings.ToUpper(strings.TrimSpace(m.cfg.Mode))
	m.modeBadge.SetText(mode)

	variant := badge.VariantAccent
	switch mode {
	case "INSERT":
		variant = badge.VariantSuccess
	case "VISUAL":
		variant = badge.VariantWarning
	case "COMMAND":
		variant = badge.VariantInfo
	}
	m.modeBadge.SetVariant(variant)
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

func tickClock() tea.Cmd {
	return tea.Tick(time.Second, func(now time.Time) tea.Msg {
		return clockTickMsg{when: now}
	})
}

func composeAlignedLine(width int, left, right string) string {
	if width <= 0 {
		line := strings.TrimSpace(strings.Join(nonEmpty(left, right), " "))
		return line
	}
	if right == "" {
		return left
	}
	if left == "" {
		pad := max(0, width-visibleWidth(right))
		return strings.Repeat(" ", pad) + right
	}
	lw := visibleWidth(left)
	rw := visibleWidth(right)
	pad := width - lw - rw
	if pad < 1 {
		pad = 1
	}
	return left + strings.Repeat(" ", pad) + right
}

func visibleWidth(s string) int {
	return lipgloss.Width(ansi.Strip(s))
}

func nonEmpty(parts ...string) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return out
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
