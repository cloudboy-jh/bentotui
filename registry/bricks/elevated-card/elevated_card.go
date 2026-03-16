// Brick: ElevatedCard:
// +-----------------------------------+
// | title                              |
// | divider                            |
// | content rows                       |
// +-----------------------------------+
// Raised surface container for sectioned dashboard/app regions.
// Copy this file into your project: bento add elevated-card
package elevatedcard

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

type Model struct {
	title   string
	meta    string
	footer  string
	content tea.Model
	focused bool
	inset   int
	width   int
	height  int
}

type Option func(*Model)

func Title(v string) Option      { return func(m *Model) { m.title = v } }
func Meta(v string) Option       { return func(m *Model) { m.meta = v } }
func Footer(v string) Option     { return func(m *Model) { m.footer = v } }
func Content(v tea.Model) Option { return func(m *Model) { m.content = v } }
func Inset(n int) Option         { return func(m *Model) { m.inset = clamp(n, 0, 4) } }

func New(opts ...Option) *Model {
	m := &Model{inset: 1}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Model) Init() tea.Cmd {
	if m.content == nil {
		return nil
	}
	return m.content.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.content == nil {
		return m, nil
	}
	u, cmd := m.content.Update(msg)
	m.content = u
	return m, cmd
}

func (m *Model) View() tea.View {
	t := theme.CurrentTheme()
	w, h := m.width, m.height
	if w <= 0 {
		w = 40
	}
	if h <= 0 {
		h = 8
	}

	outerBG := pick(t.Surface.Panel, t.Surface.Canvas)
	cardColors := styles.New(t).ElevatedCardColors(m.focused)
	borderBG := cardColors.FrameBG
	headerBG := cardColors.HeaderBG
	bodyBG := cardColors.BodyBG
	footerBG := cardColors.FooterBG
	shadowBG := cardColors.ShadowBG
	titleFG := cardColors.FrameFG
	bodyFG := t.Text.Primary
	metaFG := pick(t.Text.Muted, t.Text.Primary)
	footerFG := pick(t.Text.Muted, t.Text.Primary)
	leftEdgeBG := cardColors.FocusBG

	inset := clamp(m.inset, 0, min(max(0, (w-6)/2), max(0, (h-6)/2)))
	shadow := w >= 10 && h >= 7
	cardW := w - (inset * 2)
	cardH := h - (inset * 2)
	if shadow {
		cardW--
		cardH--
	}
	if cardW < 6 || cardH < 5 {
		shadow = false
		cardW = w - (inset * 2)
		cardH = h - (inset * 2)
	}
	if cardW < 4 || cardH < 4 {
		inset = 0
		shadow = false
		cardW = w
		cardH = h
	}

	body := ""
	if m.content != nil {
		body = viewString(m.content.View())
	}
	bodyLines := strings.Split(body, "\n")
	for i := range bodyLines {
		bodyLines[i] = ansi.Strip(bodyLines[i])
	}

	metaRow := strings.TrimSpace(m.meta) != ""
	footerRows := strings.TrimSpace(m.footer) != ""
	reserved := 4 // top + header + divider + bottom
	if metaRow {
		reserved++
	}
	if footerRows {
		reserved += 2
	}
	if reserved > cardH {
		if footerRows {
			footerRows = false
			reserved -= 2
		}
	}
	if reserved > cardH {
		metaRow = false
		reserved = 4
	}
	contentRows := max(1, cardH-reserved)

	cardRows := make([]string, 0, cardH)
	cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, borderBG, borderBG, titleFG, ""))
	cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, headerBG, borderBG, titleFG, ansi.Strip(m.title)))
	if metaRow {
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, bodyBG, borderBG, metaFG, " "+ansi.Strip(m.meta)))
	}
	cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, borderBG, borderBG, metaFG, ""))

	for i := 0; i < contentRows; i++ {
		line := ""
		if i < len(bodyLines) {
			line = bodyLines[i]
		}
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, bodyBG, borderBG, bodyFG, " "+line))
	}

	if footerRows {
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, borderBG, borderBG, metaFG, ""))
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, footerBG, borderBG, footerFG, " "+ansi.Strip(m.footer)))
	}
	for len(cardRows) < cardH-1 {
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, bodyBG, borderBG, bodyFG, ""))
	}
	cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, borderBG, borderBG, titleFG, ""))

	rows := make([]string, 0, h)
	for y := 0; y < h; y++ {
		line := ""
		if inset > 0 {
			line += styles.RowClip(outerBG, bodyFG, inset, "")
		}

		inCardY := y >= inset && y < inset+cardH
		shadowRow := shadow && y == inset+cardH
		shadowCol := shadow && y >= inset+1 && y <= inset+cardH

		used := inset
		if inCardY {
			line += cardRows[y-inset]
			used += cardW
		} else if shadowRow {
			line += styles.RowClip(shadowBG, bodyFG, cardW, "")
			used += cardW
		}
		if shadowCol {
			line += styles.RowClip(shadowBG, bodyFG, 1, "")
			used++
		}
		if used < w {
			line += styles.RowClip(outerBG, bodyFG, w-used, "")
		}
		rows = append(rows, line)
	}

	return tea.NewView(strings.Join(rows, "\n"))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if s, ok := m.content.(interface{ SetSize(int, int) }); ok {
		inset := clamp(m.inset, 0, min(max(0, (width-6)/2), max(0, (height-6)/2)))
		shadow := width >= 10 && height >= 7
		cardW := width - (inset * 2)
		cardH := height - (inset * 2)
		if shadow {
			cardW--
			cardH--
		}
		if cardW < 6 || cardH < 5 {
			shadow = false
			cardW = width - (inset * 2)
			cardH = height - (inset * 2)
		}

		reserved := 4
		if strings.TrimSpace(m.meta) != "" {
			reserved++
		}
		if strings.TrimSpace(m.footer) != "" {
			reserved += 2
		}
		contentH := max(1, cardH-reserved)
		s.SetSize(max(0, cardW-2), max(0, contentH))
	}
}

func (m *Model) SetTitle(v string)  { m.title = v }
func (m *Model) SetMeta(v string)   { m.meta = v }
func (m *Model) SetFooter(v string) { m.footer = v }

func (m *Model) Focus()          { m.focused = true }
func (m *Model) Blur()           { m.focused = false }
func (m *Model) IsFocused() bool { return m.focused }

func slabRow(width int, leftBG, centerBG, rightBG, fg, content string) string {
	if width <= 0 {
		return ""
	}
	if width == 1 {
		return styles.RowClip(centerBG, fg, 1, content)
	}
	left := styles.RowClip(leftBG, fg, 1, "")
	if width == 2 {
		right := styles.RowClip(rightBG, fg, 1, "")
		return left + right
	}
	mid := styles.RowClip(centerBG, fg, width-2, content)
	right := styles.RowClip(rightBG, fg, 1, "")
	return left + mid + right
}

func viewString(v tea.View) string {
	if v.Content == nil {
		return ""
	}
	if r, ok := v.Content.(interface{ Render() string }); ok {
		return r.Render()
	}
	if s, ok := v.Content.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(v.Content)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
