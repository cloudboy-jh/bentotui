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

	rail := pick(t.Border.Subtle, t.Border.Normal)
	if m.focused {
		rail = pick(t.Border.Focus, rail)
	}
	outerBG := pick(t.Surface.Panel, t.Surface.Canvas)
	innerBG := pick(t.Surface.Elevated, t.Surface.Panel)
	titleFG := pick(t.Text.Primary, t.Text.Muted)
	bodyFG := t.Text.Primary
	metaFG := pick(t.Text.Muted, t.Text.Primary)
	footerFG := pick(t.Text.Muted, t.Text.Primary)
	inset := clamp(m.inset, 0, min(max(0, (w-2)/2), max(0, (h-2)/2)))
	innerW := w - (inset * 2)
	innerH := h - (inset * 2)
	if innerW < 2 || innerH < 3 {
		inset = 0
		innerW = w
		innerH = h
	}

	innerRows := make([]string, 0, innerH)
	innerRows = append(innerRows, row(innerW, rail, innerBG, titleFG, " "+m.title))
	innerRows = append(innerRows, row(innerW, rail, innerBG, pick(t.Border.Subtle, t.Text.Muted), strings.Repeat("─", max(0, innerW-2))))
	if strings.TrimSpace(m.meta) != "" {
		innerRows = append(innerRows, row(innerW, rail, innerBG, metaFG, " "+ansi.Strip(m.meta)))
	}

	body := ""
	if m.content != nil {
		body = viewString(m.content.View())
	}
	lines := strings.Split(body, "\n")
	baseRows := 2
	if strings.TrimSpace(m.meta) != "" {
		baseRows++
	}
	footerRows := 0
	if strings.TrimSpace(m.footer) != "" {
		footerRows = 2 // divider + footer
	}
	contentLimit := max(0, innerH-baseRows-footerRows)
	for i := 0; i < contentLimit; i++ {
		idx := i
		line := ""
		if idx >= 0 && idx < len(lines) {
			line = ansi.Strip(lines[idx])
		}
		innerRows = append(innerRows, row(innerW, rail, innerBG, bodyFG, " "+line))
	}
	if strings.TrimSpace(m.footer) != "" {
		innerRows = append(innerRows, row(innerW, rail, innerBG, pick(t.Border.Subtle, t.Text.Muted), strings.Repeat("─", max(0, innerW-2))))
		innerRows = append(innerRows, row(innerW, rail, innerBG, footerFG, " "+ansi.Strip(m.footer)))
	}
	for len(innerRows) < innerH {
		innerRows = append(innerRows, row(innerW, rail, innerBG, bodyFG, ""))
	}

	rows := make([]string, 0, h)
	for y := 0; y < h; y++ {
		if y < inset || y >= inset+innerH {
			rows = append(rows, styles.RowClip(outerBG, bodyFG, w, ""))
			continue
		}
		inner := innerRows[y-inset]
		if inset == 0 {
			rows = append(rows, inner)
			continue
		}
		left := styles.RowClip(outerBG, bodyFG, inset, "")
		right := styles.RowClip(outerBG, bodyFG, w-inset-innerW, "")
		rows = append(rows, left+inner+right)
	}

	return tea.NewView(strings.Join(rows, "\n"))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if s, ok := m.content.(interface{ SetSize(int, int) }); ok {
		inset := clamp(m.inset, 0, min(max(0, (width-2)/2), max(0, (height-2)/2)))
		innerW := width - (inset * 2)
		innerH := height - (inset * 2)
		s.SetSize(max(0, innerW-2), max(0, innerH-2))
	}
}

func (m *Model) SetTitle(v string)  { m.title = v }
func (m *Model) SetMeta(v string)   { m.meta = v }
func (m *Model) SetFooter(v string) { m.footer = v }

func (m *Model) Focus()          { m.focused = true }
func (m *Model) Blur()           { m.focused = false }
func (m *Model) IsFocused() bool { return m.focused }

func row(width int, rail, bg, fg, content string) string {
	if width <= 0 {
		return ""
	}
	if width == 1 {
		return styles.RowClip(bg, fg, 1, "┃")
	}
	left := styles.RowClip(bg, rail, 1, "┃")
	right := styles.RowClip(bg, fg, width-1, content)
	return left + right
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
