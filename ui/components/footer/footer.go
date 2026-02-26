package footer

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/surface"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
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
	leftA   *Action
	rightA  *Action
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
func LeftAction(a Action) Option {
	return func(m *Model) { m.leftA = copyAction(a) }
}
func RightAction(a Action) Option {
	return func(m *Model) { m.rightA = copyAction(a) }
}
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
	left := m.renderLeftSegment()
	actions := m.renderActionBlock(-1)
	rightRaw := m.renderRightSegment()
	if m.width == 0 {
		line := strings.TrimSpace(strings.Join(nonEmpty(left, actions, rightRaw), "  "))
		return tea.NewView(styles.New(m.theme).StatusBar().Render(line))
	}
	right := rightRaw
	rightWidth := lipgloss.Width(right)
	if rightWidth > m.width {
		right = clipWidth(rightRaw, max(0, m.width))
		rightWidth = m.width
	}
	if rightWidth >= m.width {
		return m.renderLine(right)
	}

	leftArea := max(0, m.width-rightWidth)
	if rightWidth > 0 && leftArea > 0 {
		leftArea--
	}
	leftBlock := ""
	actionBlock := ""
	if left != "" && leftArea > 0 {
		leftBlock = clipWidth(left, leftArea)
		leftArea -= lipgloss.Width(leftBlock)
	}
	if leftArea > 0 {
		if leftBlock != "" {
			leftArea--
		}
		actionBlock = m.renderActionBlock(leftArea)
	}

	leftSide := strings.TrimSpace(strings.Join(nonEmpty(leftBlock, actionBlock), " "))
	line := strings.TrimSpace(strings.Join(nonEmpty(leftSide, right), " "))
	return m.renderLine(line)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = 1
}

func (m *Model) GetSize() (width, height int) {
	return m.width, m.height
}

func (m *Model) SetTheme(t theme.Theme) {
	m.theme = t
}

func (m *Model) SetActions(actions []Action) {
	m.actions = append([]Action(nil), actions...)
}

func (m *Model) SetLeftAction(a Action) {
	m.leftA = copyAction(a)
}

func (m *Model) SetRightAction(a Action) {
	m.rightA = copyAction(a)
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

func (m *Model) renderLeftSegment() string {
	parts := make([]string, 0, 2)
	if m.leftA != nil {
		parts = append(parts, m.renderChip(*m.leftA, false))
	}
	help := m.helpText()
	text := strings.TrimSpace(strings.Join(nonEmpty(m.left, help), "  "))
	if text != "" {
		parts = append(parts, text)
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func (m *Model) renderRightSegment() string {
	parts := make([]string, 0, 2)
	if m.right != "" {
		parts = append(parts, m.right)
	}
	if m.rightA != nil {
		parts = append(parts, m.renderChip(*m.rightA, false))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func (m *Model) renderLine(text string) tea.View {
	style := styles.New(m.theme).StatusBar()
	if m.width == 0 {
		return tea.NewView(style.Render(text))
	}
	line := primitives.PaintRow(m.width, m.theme.StatusBG, m.theme.StatusText, text)
	return tea.NewView(line)
}

func (m *Model) renderActions(width int, keyOnly bool) string {
	rendered, _ := m.renderActionsForWidth(width, keyOnly)
	return rendered
}

func (m *Model) renderActionsForWidth(width int, keyOnly bool) (string, bool) {
	if len(m.actions) == 0 {
		return "", true
	}
	rendered := make([]string, 0, len(m.actions))
	used := 0
	allFit := true
	validCount := 0
	for _, a := range m.actions {
		if strings.TrimSpace(a.Key) == "" {
			continue
		}
		validCount++
		chip := m.renderChip(a, keyOnly)

		if width >= 0 {
			w := lipgloss.Width(chip)
			sep := 0
			if len(rendered) > 0 {
				sep = 1
			}
			if used+sep+w > width {
				allFit = false
				break
			}
			used += sep + w
		}
		rendered = append(rendered, chip)
	}
	if len(rendered) < validCount {
		allFit = false
	}
	return strings.Join(rendered, " "), allFit
}

func (m *Model) renderActionBlock(width int) string {
	full, allFit := m.renderActionsForWidth(width, false)
	if allFit {
		return full
	}
	keyOnly, _ := m.renderActionsForWidth(width, true)
	return keyOnly
}

func (m *Model) renderChip(a Action, keyOnly bool) string {
	return primitives.ActionChip(m.theme, string(a.Variant), a.Enabled, a.Key, a.Label, keyOnly)
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

func clipWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().MaxWidth(width).Render(s)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func copyAction(a Action) *Action {
	b := a
	return &b
}
