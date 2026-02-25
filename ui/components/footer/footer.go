package footer

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/surface"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

type Option func(*Model)

type ActionVariant string

const (
	ActionNormal  ActionVariant = "normal"
	ActionPrimary ActionVariant = "primary"
	ActionMuted   ActionVariant = "muted"
	ActionDanger  ActionVariant = "danger"
)

type Action struct {
	Key     string
	Label   string
	Variant ActionVariant
	Enabled bool
}

type Model struct {
	left    string
	right   string
	help    core.Bindable
	actions []Action
	theme   theme.Theme
	width   int
	height  int
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
func Actions(actions ...Action) Option {
	return func(m *Model) {
		m.actions = append([]Action(nil), actions...)
	}
}
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
	actions := m.renderActions(-1, false)
	if m.width == 0 {
		line := strings.TrimSpace(strings.Join(nonEmpty(left, actions, m.right), "  "))
		return tea.NewView(styles.New(m.theme).StatusBar().Render(line))
	}
	right := fitWidth(m.right, max(0, m.width))
	rightWidth := lipgloss.Width(right)
	if rightWidth >= m.width {
		return tea.NewView(m.renderLine(right))
	}

	leftArea := max(0, m.width-rightWidth-1)
	leftBlock := ""
	actionBlock := ""
	if left != "" {
		leftBlock = fitWidth(left, leftArea)
		leftArea -= lipgloss.Width(leftBlock)
	}
	if leftArea > 0 {
		if leftBlock != "" {
			leftArea--
		}
		actionBlock = m.renderActions(leftArea, false)
		if actionBlock == "" {
			actionBlock = m.renderActions(leftArea, true)
		}
	}

	leftSide := strings.TrimSpace(strings.Join(nonEmpty(leftBlock, actionBlock), " "))
	line := strings.TrimSpace(strings.Join(nonEmpty(leftSide, right), " "))
	return tea.NewView(m.renderLine(line))
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

func (m *Model) renderActions(width int, keyOnly bool) string {
	if len(m.actions) == 0 {
		return ""
	}
	sys := styles.New(m.theme)
	rendered := make([]string, 0, len(m.actions))
	used := 0
	for _, a := range m.actions {
		if strings.TrimSpace(a.Key) == "" {
			continue
		}
		enabled := a.Enabled
		keyPart := sys.FooterActionKey(string(a.Variant), enabled).Render(a.Key)
		chip := keyPart
		if !keyOnly && strings.TrimSpace(a.Label) != "" {
			chip = keyPart + " " + sys.FooterActionLabel(string(a.Variant), enabled).Render(a.Label)
		}

		if width >= 0 {
			w := lipgloss.Width(chip)
			sep := 0
			if len(rendered) > 0 {
				sep = 1
			}
			if used+sep+w > width {
				break
			}
			used += sep + w
		}
		rendered = append(rendered, chip)
	}
	return strings.Join(rendered, " ")
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

func fitWidth(s string, width int) string {
	return surface.FitWidth(s, width)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
