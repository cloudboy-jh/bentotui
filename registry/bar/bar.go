// Package bar provides a single-row themed status/navigation bar.
// Copy this file into your project: bento add bar
//
// Dependencies (real Go module imports, not copied):
//   - charm.land/bubbletea/v2
//   - charm.land/lipgloss/v2
//   - github.com/cloudboy-jh/bentotui/theme
//   - github.com/cloudboy-jh/bentotui/styles
package bar

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/theme"
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

// Model is the bar component. Stateless with respect to theme — always reads
// theme.CurrentTheme() in View().
type Model struct {
	left      string
	right     string
	leftCard  *Card
	rightCard *Card
	cards     []Card
	width     int
}

// New constructs a bar with the given options.
func New(opts ...Option) *Model {
	m := &Model{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// ── options ───────────────────────────────────────────────────────────────────

func Left(v string) Option  { return func(m *Model) { m.left = v } }
func Right(v string) Option { return func(m *Model) { m.right = v } }

func LeftCard(c Card) Option  { return func(m *Model) { cp := c; m.leftCard = &cp } }
func RightCard(c Card) Option { return func(m *Model) { cp := c; m.rightCard = &cp } }

func Cards(cards ...Card) Option {
	return func(m *Model) { m.cards = append([]Card(nil), cards...) }
}

// ── tea.Model ─────────────────────────────────────────────────────────────────

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.SetSize(ws.Width, 1)
	}
	return m, nil
}

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()

	left := m.renderLeftSegment(t)
	cards := m.renderCardBlock(t, -1)
	rightRaw := m.renderRightSegment(t)

	if m.width == 0 {
		line := strings.TrimSpace(strings.Join(nonEmpty(left, cards, rightRaw), "  "))
		return tea.NewView(styles.New(t).StatusBar().Render(line))
	}

	right := rightRaw
	rightWidth := lipgloss.Width(right)
	if rightWidth > m.width {
		right = clipWidth(rightRaw, max(0, m.width))
		rightWidth = m.width
	}
	if rightWidth >= m.width {
		return m.renderLine(t, right)
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
		cardBlock = m.renderCardBlock(t, leftArea)
	}

	leftSide := strings.TrimSpace(strings.Join(nonEmpty(leftBlock, cardBlock), " "))
	line := strings.TrimSpace(strings.Join(nonEmpty(leftSide, right), " "))
	return m.renderLine(t, line)
}

// ── sizing ────────────────────────────────────────────────────────────────────

func (m *Model) SetSize(width, _ int) { m.width = width }
func (m *Model) GetSize() (int, int)  { return m.width, 1 }

// ── mutations ─────────────────────────────────────────────────────────────────

func (m *Model) SetLeft(v string)      { m.left = v }
func (m *Model) SetRight(v string)     { m.right = v }
func (m *Model) SetCards(cards []Card) { m.cards = append([]Card(nil), cards...) }
func (m *Model) SetLeftCard(c Card)    { cp := c; m.leftCard = &cp }
func (m *Model) SetRightCard(c Card)   { cp := c; m.rightCard = &cp }

// ── rendering ─────────────────────────────────────────────────────────────────

func (m *Model) renderLine(t theme.Theme, text string) tea.View {
	bar := styles.New(t).BarColors()
	return tea.NewView(renderRow(m.width, bar.BG, bar.FG, text))
}

// renderRow renders a full-width bar row with a single lipgloss call.
// Background fills every cell — no canvas layers, no ANSI bleed.
func renderRow(width int, bg, fg, content string) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(fg)).
		Width(width).
		Render(content)
}

func (m *Model) renderLeftSegment(t theme.Theme) string {
	parts := make([]string, 0, 2)
	if m.leftCard != nil {
		parts = append(parts, renderCard(t, *m.leftCard, false))
	}
	if m.left != "" {
		parts = append(parts, m.left)
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func (m *Model) renderRightSegment(t theme.Theme) string {
	parts := make([]string, 0, 2)
	if m.right != "" {
		parts = append(parts, m.right)
	}
	if m.rightCard != nil {
		parts = append(parts, renderCard(t, *m.rightCard, false))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func (m *Model) renderCardBlock(t theme.Theme, width int) string {
	full, allFit := m.renderCardsForWidth(t, width, false)
	if allFit {
		return full
	}
	commandOnly, _ := m.renderCardsForWidth(t, width, true)
	return commandOnly
}

func (m *Model) renderCardsForWidth(t theme.Theme, width int, commandOnly bool) (string, bool) {
	if len(m.cards) == 0 {
		return "", true
	}
	rendered := make([]string, 0, len(m.cards))
	used := 0
	allFit := true
	validCount := 0
	for _, c := range m.cards {
		if strings.TrimSpace(c.Command) == "" {
			continue
		}
		validCount++
		card := renderCard(t, c, commandOnly)
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

// renderCard renders a command/label card pair for the bar.
func renderCard(t theme.Theme, c Card, commandOnly bool) string {
	sys := styles.New(t)
	commandPart := sys.FooterCardCommand(string(c.Variant), c.Enabled).Render(c.Command)
	if commandOnly || strings.TrimSpace(c.Label) == "" {
		return commandPart
	}
	return commandPart + " " + sys.FooterCardLabel(string(c.Variant), c.Enabled).Render(c.Label)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
