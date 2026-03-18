// Brick: Bar:
// +--------------------------------------------------+
// | left/status            cards            right     |
// +--------------------------------------------------+
// Single-row command or metadata strip.
// Copy this file into your project: bento add bar
package bar

import (
	"image/color"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

type CardVariant string

const (
	CardNormal  CardVariant = "normal"
	CardPrimary CardVariant = "primary"
	CardMuted   CardVariant = "muted"
	CardDanger  CardVariant = "danger"
)

type Role string

const (
	RoleTop       Role = "top"
	RoleSubheader Role = "subheader"
	RoleFooter    Role = "footer"
)

type FooterMode string

const (
	FooterModeNormal   FooterMode = "normal"
	FooterModeAnchored FooterMode = "anchored"
)

type AnchoredCardStyle string

const (
	AnchoredCardStylePlain AnchoredCardStyle = "plain"
	AnchoredCardStyleChip  AnchoredCardStyle = "chip"
	AnchoredCardStyleMixed AnchoredCardStyle = "mixed"
)

type Card struct {
	Command  string
	Label    string
	Variant  CardVariant
	Enabled  bool
	Priority int
}

type Option func(*Model)

// Model is the bar component. Reads theme from WithTheme() or falls back
// to theme.CurrentTheme() — never stores global state directly.
type Model struct {
	left         string
	statusPill   string
	right        string
	leftCard     *Card
	rightCard    *Card
	cards        []Card
	width        int
	compactCards bool
	role         Role
	footerMode   FooterMode
	anchoredCard AnchoredCardStyle
	theme        theme.Theme // nil = use theme.CurrentTheme()
}

func New(opts ...Option) *Model {
	m := &Model{role: RoleTop, footerMode: FooterModeNormal, anchoredCard: AnchoredCardStylePlain}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func Left(v string) Option  { return func(m *Model) { m.left = v } }
func Right(v string) Option { return func(m *Model) { m.right = v } }
func StatusPill(v string) Option {
	return func(m *Model) { m.statusPill = strings.TrimSpace(v) }
}
func LeftCard(c Card) Option  { return func(m *Model) { cp := c; m.leftCard = &cp } }
func RightCard(c Card) Option { return func(m *Model) { cp := c; m.rightCard = &cp } }
func Cards(cards ...Card) Option {
	return func(m *Model) { m.cards = append([]Card(nil), cards...) }
}
func CompactCards() Option  { return func(m *Model) { m.compactCards = true } }
func RoleTopBar() Option    { return func(m *Model) { m.role = RoleTop } }
func RoleSubBar() Option    { return func(m *Model) { m.role = RoleSubheader } }
func RoleFooterBar() Option { return func(m *Model) { m.role = RoleFooter } }
func FooterAnchored() Option {
	return func(m *Model) {
		m.role = RoleFooter
		m.footerMode = FooterModeAnchored
	}
}
func AnchoredCardStyleMode(style AnchoredCardStyle) Option {
	return func(m *Model) { m.anchoredCard = style }
}

// WithTheme sets the theme for this bar instance.
func WithTheme(t theme.Theme) Option {
	return func(m *Model) { m.theme = t }
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.SetSize(ws.Width, 1)
	}
	return m, nil
}

func (m *Model) View() tea.View {
	t := m.activeTheme()

	left := m.renderLeftSegment(t)
	cards := m.renderCardBlock(t, -1)
	rightRaw := m.renderRightSegment(t)

	if m.width == 0 {
		line := strings.TrimSpace(strings.Join(nonEmpty(left, cards, rightRaw), "  "))
		bg, fg := m.rowColors(t)
		return tea.NewView(lipgloss.NewStyle().
			Foreground(fg).
			Background(bg).
			Render(line))
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
	line := composeAlignedLine(m.width, leftSide, right)
	return m.renderLine(t, line)
}

func (m *Model) SetSize(width, _ int) { m.width = width }
func (m *Model) GetSize() (int, int)  { return m.width, 1 }

func (m *Model) SetLeft(v string)       { m.left = v }
func (m *Model) SetRight(v string)      { m.right = v }
func (m *Model) SetStatusPill(v string) { m.statusPill = strings.TrimSpace(v) }
func (m *Model) SetCards(cards []Card)  { m.cards = append([]Card(nil), cards...) }
func (m *Model) SetLeftCard(c Card)     { cp := c; m.leftCard = &cp }
func (m *Model) SetRightCard(c Card)    { cp := c; m.rightCard = &cp }
func (m *Model) SetCompactCards(v bool) { m.compactCards = v }
func (m *Model) SetRole(role Role)      { m.role = role }
func (m *Model) SetFooterMode(mode FooterMode) {
	m.footerMode = mode
}
func (m *Model) SetAnchored(v bool) {
	if v {
		m.footerMode = FooterModeAnchored
	} else {
		m.footerMode = FooterModeNormal
	}
}
func (m *Model) SetAnchoredCardStyle(style AnchoredCardStyle) { m.anchoredCard = style }

// SetTheme updates the theme. Call on ThemeChangedMsg.
func (m *Model) SetTheme(t theme.Theme) { m.theme = t }

func (m *Model) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

func (m *Model) rowColors(t theme.Theme) (bg, fg color.Color) {
	anchored := m.role == RoleFooter && m.footerMode == FooterModeAnchored
	switch m.role {
	case RoleSubheader:
		return t.BackgroundPanel(), t.TextMuted()
	case RoleFooter:
		if anchored {
			return t.FooterBG(), t.FooterFG()
		}
		return t.BackgroundPanel(), t.Text()
	default:
		return t.BarBG(), t.BarFG()
	}
}

func (m *Model) renderLine(t theme.Theme, text string) tea.View {
	bg, fg := m.rowColors(t)
	return tea.NewView(renderRow(m.width, bg, fg, text))
}

func renderRow(width int, bg, fg color.Color, content string) string {
	if width <= 0 {
		return ""
	}
	return styles.Row(bg, fg, width, styles.ClipANSI(content, width))
}

func (m *Model) renderLeftSegment(t theme.Theme) string {
	parts := make([]string, 0, 2)
	if m.statusPill != "" {
		parts = append(parts, lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Foreground(t.Text()).
			Background(t.BackgroundPanel()).
			Render(m.statusPill))
	}
	if m.leftCard != nil {
		parts = append(parts, renderCard(t, *m.leftCard, true, m.compactMode(), m.anchoredMode()))
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
		parts = append(parts, renderCard(t, *m.rightCard, true, m.compactMode(), m.anchoredMode()))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

func (m *Model) renderCardBlock(t theme.Theme, width int) string {
	if len(m.cards) == 0 {
		return ""
	}
	entries := make([]cardEntry, 0, len(m.cards))
	for i, c := range m.cards {
		if strings.TrimSpace(c.Command) == "" {
			continue
		}
		entries = append(entries, cardEntry{card: c, index: i, showLabel: true})
	}
	if len(entries) == 0 {
		return ""
	}
	if width < 0 {
		return joinEntries(t, entries, m.compactMode(), m.anchoredMode())
	}
	rendered := joinEntries(t, entries, m.compactMode(), m.anchoredMode())
	if lipgloss.Width(rendered) <= width {
		return rendered
	}
	order := truncateOrder(entries)
	for _, idx := range order {
		entries[idx].showLabel = false
		rendered = joinEntries(t, entries, m.compactMode(), m.anchoredMode())
		if lipgloss.Width(rendered) <= width {
			return rendered
		}
	}
	keep := make([]bool, len(entries))
	for i := range keep {
		keep[i] = true
	}
	for _, idx := range order {
		keep[idx] = false
		rendered = joinKeptEntries(t, entries, keep, m.compactMode(), m.anchoredMode())
		if lipgloss.Width(rendered) <= width {
			return rendered
		}
	}
	return clipWidth(rendered, width)
}

func (m *Model) compactMode() bool {
	return m.compactCards || (m.role == RoleFooter && m.footerMode == FooterModeAnchored)
}
func (m *Model) anchoredMode() bool {
	return m.role == RoleFooter && m.footerMode == FooterModeAnchored
}

func renderCard(t theme.Theme, c Card, showLabel bool, compact bool, anchored bool) string {
	if anchored {
		var commandFG color.Color
		switch c.Variant {
		case CardDanger:
			commandFG = t.Error()
		case CardMuted:
			commandFG = t.FooterMuted()
		default:
			commandFG = t.FooterFG()
		}
		if !c.Enabled {
			commandFG = t.FooterMuted()
		}
		cmdStyle := lipgloss.NewStyle().Bold(true).Foreground(commandFG)
		commandPart := cmdStyle.Render(c.Command)
		label := strings.TrimSpace(c.Label)
		if !showLabel || label == "" {
			return commandPart
		}
		labelStyle := lipgloss.NewStyle().Foreground(t.FooterMuted())
		labelPart := labelStyle.Render(label)
		if compact {
			if strings.EqualFold(label, c.Command) || strings.HasPrefix(strings.ToLower(label), strings.ToLower(c.Command)+" ") {
				return labelPart
			}
			return commandPart + " " + labelPart
		}
		return commandPart + " " + labelPart
	}

	var commandFG, commandBG color.Color
	switch c.Variant {
	case CardPrimary:
		commandFG, commandBG = t.SelectionFG(), t.SelectionBG()
	case CardDanger:
		commandFG, commandBG = t.TextInverse(), t.Error()
	case CardMuted:
		commandFG, commandBG = t.TextMuted(), t.BackgroundPanel()
	default:
		commandFG, commandBG = t.SelectionFG(), t.BorderFocus()
	}
	if !c.Enabled {
		commandFG, commandBG = t.TextMuted(), t.BackgroundPanel()
	}
	cmdStyle := lipgloss.NewStyle().Bold(true).Foreground(commandFG).Background(commandBG)
	commandPart := cmdStyle.Render(c.Command)
	label := strings.TrimSpace(c.Label)
	if !showLabel || label == "" {
		return commandPart
	}
	labelStyle := lipgloss.NewStyle().
		Foreground(t.Text()).
		Background(t.BackgroundInteractive())
	if !c.Enabled {
		labelStyle = lipgloss.NewStyle().Foreground(t.TextMuted()).Background(t.BackgroundPanel())
	}
	labelPart := labelStyle.Render(label)
	if compact {
		if strings.EqualFold(label, c.Command) || strings.HasPrefix(strings.ToLower(label), strings.ToLower(c.Command)+" ") {
			return labelPart
		}
		return commandPart + "·" + labelPart
	}
	return commandPart + " " + labelPart
}

type cardEntry struct {
	card      Card
	index     int
	showLabel bool
}

func joinEntries(t theme.Theme, entries []cardEntry, compact bool, anchored bool) string {
	parts := make([]string, len(entries))
	for i, e := range entries {
		parts[i] = renderCard(t, e.card, e.showLabel, compact, anchored)
	}
	return strings.Join(parts, " ")
}

func joinKeptEntries(t theme.Theme, entries []cardEntry, keep []bool, compact bool, anchored bool) string {
	parts := make([]string, 0, len(entries))
	for i, e := range entries {
		if !keep[i] {
			continue
		}
		parts = append(parts, renderCard(t, e.card, e.showLabel, compact, anchored))
	}
	return strings.Join(parts, " ")
}

func truncateOrder(entries []cardEntry) []int {
	idx := make([]int, len(entries))
	for i := range entries {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		a := entries[idx[i]]
		b := entries[idx[j]]
		if a.card.Priority != b.card.Priority {
			return a.card.Priority < b.card.Priority
		}
		return a.index > b.index
	})
	return idx
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

func clipWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return styles.ClipANSI(s, width)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func composeAlignedLine(width int, left, right string) string {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if width <= 0 {
		return strings.TrimSpace(strings.Join(nonEmpty(left, right), " "))
	}
	if right == "" {
		return left
	}
	if left == "" {
		pad := max(0, width-lipgloss.Width(right))
		return strings.Repeat(" ", pad) + right
	}
	lw := lipgloss.Width(left)
	rw := lipgloss.Width(right)
	pad := width - lw - rw
	if pad < 1 {
		pad = 1
	}
	return left + strings.Repeat(" ", pad) + right
}
