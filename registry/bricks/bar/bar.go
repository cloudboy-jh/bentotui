// Brick: Bar:
// +--------------------------------------------------+
// | left/status            cards            right     |
// +--------------------------------------------------+
// Single-row command or metadata strip.
// Package bar provides a single-row themed status/navigation bar.
// Copy this file into your project: bento add bar
//
// Dependencies (real Go module imports, not copied):
//   - charm.land/bubbletea/v2
//   - charm.land/lipgloss/v2
//   - github.com/cloudboy-jh/bentotui/theme
//   - github.com/cloudboy-jh/bentotui/theme/styles
package bar

import (
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

// CardVariant controls the visual weight of a card's command badge.
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

// Card is a single keybinding hint rendered inside a bar.
type Card struct {
	Command  string
	Label    string
	Variant  CardVariant
	Enabled  bool
	Priority int
}

// Option configures a Model at construction time.
type Option func(*Model)

// Model is the bar component. Stateless with respect to theme — always reads
// theme.CurrentTheme() in View().
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
}

// New constructs a bar with the given options.
func New(opts ...Option) *Model {
	m := &Model{role: RoleTop, footerMode: FooterModeNormal, anchoredCard: AnchoredCardStylePlain}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// ── options ───────────────────────────────────────────────────────────────────

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

func CompactCards() Option { return func(m *Model) { m.compactCards = true } }
func RoleTopBar() Option   { return func(m *Model) { m.role = RoleTop } }
func RoleSubBar() Option   { return func(m *Model) { m.role = RoleSubheader } }
func RoleFooterBar() Option {
	return func(m *Model) { m.role = RoleFooter }
}
func FooterAnchored() Option {
	return func(m *Model) {
		m.role = RoleFooter
		m.footerMode = FooterModeAnchored
	}
}
func AnchoredCardStyleMode(style AnchoredCardStyle) Option {
	return func(m *Model) { m.anchoredCard = style }
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
		colors := styles.New(t).StatusRowColors(string(m.role), m.role == RoleFooter && m.footerMode == FooterModeAnchored)
		// No width — single solid background, no seam.
		return tea.NewView(lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.FG)).
			Background(lipgloss.Color(colors.BG)).
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

// ── sizing ────────────────────────────────────────────────────────────────────

func (m *Model) SetSize(width, _ int) { m.width = width }
func (m *Model) GetSize() (int, int)  { return m.width, 1 }

// ── mutations ─────────────────────────────────────────────────────────────────

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
func (m *Model) SetAnchoredCardStyle(style AnchoredCardStyle) {
	m.anchoredCard = style
}

// ── rendering ─────────────────────────────────────────────────────────────────

func (m *Model) renderLine(t theme.Theme, text string) tea.View {
	bar := styles.New(t).StatusRowColors(string(m.role), m.role == RoleFooter && m.footerMode == FooterModeAnchored)
	return tea.NewView(renderRow(m.width, bar.BG, bar.FG, text))
}

// renderRow renders a full-width bar row with a single lipgloss call.
// Background fills every cell — no canvas seams, no ANSI bleed.
// The entire width is owned by one style.Width() call so the gap between
// left and right segments is the same background color as the rest of the bar.
func renderRow(width int, bg, fg, content string) string {
	if width <= 0 {
		return ""
	}
	return styles.Row(bg, fg, width, styles.ClipANSI(content, width))
}

func (m *Model) renderLeftSegment(t theme.Theme) string {
	parts := make([]string, 0, 2)
	if m.statusPill != "" {
		parts = append(parts, styles.New(t).StatusPillMuted().Render(m.statusPill))
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

// renderCard renders a command/label card pair for the bar.
func renderCard(t theme.Theme, c Card, showLabel bool, compact bool, anchored bool) string {
	if anchored {
		sys := styles.New(t)
		commandPart := sys.FooterCardCommandAnchored(string(c.Variant), c.Enabled).Render(c.Command)
		label := strings.TrimSpace(c.Label)
		if !showLabel || label == "" {
			return commandPart
		}
		labelPart := sys.FooterCardLabelAnchored(string(c.Variant), c.Enabled).Render(label)
		if compact {
			if strings.EqualFold(label, c.Command) || strings.HasPrefix(strings.ToLower(label), strings.ToLower(c.Command)+" ") {
				return labelPart
			}
			return commandPart + " " + labelPart
		}
		return commandPart + " " + labelPart
	}

	sys := styles.New(t)
	commandPart := sys.FooterCardCommand(string(c.Variant), c.Enabled).Render(c.Command)
	label := strings.TrimSpace(c.Label)
	if !showLabel || label == "" {
		return commandPart
	}
	labelPart := sys.FooterCardLabel(string(c.Variant), c.Enabled).Render(label)
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
	return styles.ClipANSI(s, width)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// composeAlignedLine builds the raw content string for a bar row.
// It uses plain spaces to position left and right segments — the actual
// background color is applied by renderRow() wrapping this output with
// a single lipgloss Width() call, so every cell (including the gap) is
// the same solid background. Never call this outside of renderLine.
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
		// renderRow's Width() will right-justify by filling the full width,
		// but we still need the right content anchored at the right edge.
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
