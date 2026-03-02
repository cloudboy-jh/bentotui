// Package bar provides the single themed bar component used for both
// header and footer positions. The shell places two instances: one at the
// top of the viewport and one at the bottom. Both share identical API,
// rendering, and card truncation behavior.
package bar

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/focus"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

// CardVariant controls the visual weight of a card's command badge.
type CardVariant string

const (
	CardNormal  CardVariant = "normal"
	CardPrimary CardVariant = "primary"
	CardMuted   CardVariant = "muted"
	CardDanger  CardVariant = "danger"
)

// Card is a single keybinding hint rendered inside a bar.
type Card struct {
	Command string
	Label   string
	Variant CardVariant
	Enabled bool
}

// Option configures a Model at construction time.
type Option func(*Model)

// Model is the bar component. Instantiate it twice — once for the header
// position and once for the footer position. The shell determines placement
// via layer Z-order; the bar itself is position-agnostic.
type Model struct {
	left      string
	right     string
	leftCard  *Card
	rightCard *Card
	help      core.Bindable
	cards     []Card
	theme     theme.Theme
	width     int
	height    int
	focusIdx  int
}

// New constructs a Model with the given options.
func New(opts ...Option) *Model {
	m := &Model{theme: theme.Preset(theme.DefaultName), focusIdx: -1}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// ── options ───────────────────────────────────────────────────────────────────

func Left(v string) Option  { return func(m *Model) { m.left = v } }
func Right(v string) Option { return func(m *Model) { m.right = v } }

func LeftCard(c Card) Option  { return func(m *Model) { m.leftCard = copyCard(c) } }
func RightCard(c Card) Option { return func(m *Model) { m.rightCard = copyCard(c) } }

func Cards(cards ...Card) Option {
	return func(m *Model) { m.cards = append([]Card(nil), cards...) }
}

func HelpFrom(b core.Bindable) Option {
	return func(m *Model) { m.help = b }
}

// ── tea.Model ─────────────────────────────────────────────────────────────────

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(v.Width, 1)
	case focus.FocusChangedMsg:
		m.focusIdx = v.To
	}
	return m, nil
}

func (m *Model) View() tea.View {
	left := m.renderLeftSegment()
	cards := m.renderCardBlock(-1)
	rightRaw := m.renderRightSegment()

	if m.width == 0 {
		line := strings.TrimSpace(strings.Join(nonEmpty(left, cards, rightRaw), "  "))
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
	cardBlock := ""
	if left != "" && leftArea > 0 {
		leftBlock = clipWidth(left, leftArea)
		leftArea -= lipgloss.Width(leftBlock)
	}
	if leftArea > 0 {
		if leftBlock != "" {
			leftArea--
		}
		cardBlock = m.renderCardBlock(leftArea)
	}

	leftSide := strings.TrimSpace(strings.Join(nonEmpty(leftBlock, cardBlock), " "))
	line := strings.TrimSpace(strings.Join(nonEmpty(leftSide, right), " "))
	return m.renderLine(line)
}

// ── sizing ────────────────────────────────────────────────────────────────────

func (m *Model) SetSize(width, _ int) {
	m.width = width
	m.height = 1
}

func (m *Model) GetSize() (int, int) { return m.width, m.height }

// ── theme ─────────────────────────────────────────────────────────────────────

func (m *Model) SetTheme(t theme.Theme) { m.theme = t }

// ── mutations ─────────────────────────────────────────────────────────────────

func (m *Model) SetCards(cards []Card) {
	m.cards = append([]Card(nil), cards...)
}

func (m *Model) SetLeftCard(c Card) {
	m.leftCard = copyCard(c)
}

func (m *Model) SetRightCard(c Card) {
	m.rightCard = copyCard(c)
}

// ── rendering ─────────────────────────────────────────────────────────────────

func (m *Model) renderLine(text string) tea.View {
	if m.width == 0 {
		return tea.NewView(styles.New(m.theme).StatusBar().Render(text))
	}
	bar := styles.New(m.theme).BarColors()
	return tea.NewView(primitives.RenderRow(m.width, bar.BG, bar.FG, text))
}

func (m *Model) renderLeftSegment() string {
	parts := make([]string, 0, 2)
	if m.leftCard != nil {
		parts = append(parts, m.renderCard(*m.leftCard, false))
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
	if m.rightCard != nil {
		parts = append(parts, m.renderCard(*m.rightCard, false))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func (m *Model) renderCardBlock(width int) string {
	full, allFit := m.renderCardsForWidth(width, false)
	if allFit {
		return full
	}
	commandOnly, _ := m.renderCardsForWidth(width, true)
	return commandOnly
}

func (m *Model) renderCardsForWidth(width int, commandOnly bool) (string, bool) {
	if len(m.cards) == 0 {
		return "", true
	}
	rendered := make([]string, 0, len(m.cards))
	used := 0
	allFit := true
	validCount := 0
	for i, c := range m.cards {
		if strings.TrimSpace(c.Command) == "" {
			continue
		}
		validCount++
		cardModel := c
		if i == m.focusIdx {
			cardModel.Variant = CardPrimary
		}
		card := m.renderCard(cardModel, commandOnly)
		if width >= 0 {
			w := lipgloss.Width(card)
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
		rendered = append(rendered, card)
	}
	if len(rendered) < validCount {
		allFit = false
	}
	return strings.Join(rendered, " "), allFit
}

func (m *Model) renderCard(c Card, commandOnly bool) string {
	return primitives.Card(m.theme, string(c.Variant), c.Enabled, c.Command, c.Label, commandOnly)
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

// ── helpers ───────────────────────────────────────────────────────────────────

func nonEmpty(parts ...string) []string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return out
}

func clipWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().MaxWidth(width).Render(s)
}

func copyCard(c Card) *Card { b := c; return &b }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
