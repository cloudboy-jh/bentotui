package statusbar

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/surface"
	"github.com/cloudboy-jh/bentotui/theme"
)

type Option func(*Model)

type Model struct {
	left   string
	right  string
	help   core.Bindable
	theme  theme.Theme
	width  int
	height int
}

func New(opts ...Option) *Model {
	m := &Model{theme: theme.Preset(theme.DefaultName)}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func Left(v string) Option  { return func(m *Model) { m.left = v } }
func Right(v string) Option { return func(m *Model) { m.right = v } }
func HelpFrom(b core.Bindable) Option {
	return func(m *Model) { m.help = b }
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(v.Width, 1)
	}
	return m, nil
}

func (m *Model) View() tea.View {
	help := m.helpText()
	left := strings.TrimSpace(strings.Join([]string{m.left, help}, "  "))
	if m.width == 0 {
		return tea.NewView(styles.New(m.theme).StatusBar().Render(left))
	}
	right := fitWidth(m.right, max(0, m.width))
	rightWidth := lipgloss.Width(right)
	if rightWidth >= m.width {
		return tea.NewView(m.renderLine(right))
	}
	leftWidth := max(0, m.width-rightWidth-1)
	leftBlock := fitWidth(left, leftWidth)
	return tea.NewView(m.renderLine(leftBlock + " " + right))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) GetSize() (width, height int) {
	return m.width, m.height
}

func (m *Model) SetTheme(t theme.Theme) {
	m.theme = t
}

func (m *Model) helpText() string {
	if m.help == nil {
		return ""
	}
	bindings := m.help.Bindings()
	parts := make([]string, 0, len(bindings))
	for _, b := range bindings {
		if !b.Enabled() {
			continue
		}
		h := b.Help()
		if h.Key == "" || h.Desc == "" {
			continue
		}
		parts = append(parts, key.NewBinding(key.WithKeys(h.Key), key.WithHelp(h.Key, h.Desc)).Help().Key+": "+h.Desc)
	}
	return strings.Join(parts, " • ")
}

func (m *Model) renderLine(text string) string {
	style := styles.New(m.theme).StatusBar()
	if m.width == 0 {
		return style.Render(text)
	}
	style = style.Width(max(0, m.width))
	text = fitWidth(text, m.width)
	return style.Render(text)
}

func fitWidth(s string, width int) string {
	return surface.FitWidth(s, width)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
