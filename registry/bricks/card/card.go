// Brick: Card
// +-----------------------------------+
// | title                [meta]        |
// +-----------------------------------+
// | content rows                       |
// +-----------------------------------+
// | footer                             |
// +-----------------------------------+
//
// The one content-container brick in BentoTUI.
//
// ElevationRaised (default) — chrome header band + body slab.
//
//	Visually lifted off the canvas. Use for dashboard widgets,
//	file previews, code panels.
//
// ElevationFlat — plain titled container, flush with parent surface.
//
//	Title row + separator + content. Use for sidebars, panes,
//	split-view regions.
//
// Copy this file into your project: bento add card
package card

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

// Elevation controls the visual depth of the card.
type Elevation string

const (
	ElevationRaised Elevation = "raised" // chrome band + body slab (default)
	ElevationFlat   Elevation = "flat"   // plain titled container, no chrome band
)

type Model struct {
	title     string
	meta      string
	footer    string
	content   tea.Model
	elevation Elevation
	focused   bool
	inset     int
	width     int
	height    int
	theme     theme.Theme // nil = use theme.CurrentTheme()
}

type Option func(*Model)

func Title(v string) Option      { return func(m *Model) { m.title = v } }
func Meta(v string) Option       { return func(m *Model) { m.meta = v } }
func Footer(v string) Option     { return func(m *Model) { m.footer = v } }
func Content(v tea.Model) Option { return func(m *Model) { m.content = v } }
func Inset(n int) Option         { return func(m *Model) { m.inset = clamp(n, 0, 4) } }

// Flat sets ElevationFlat — plain titled container with separator.
func Flat() Option { return func(m *Model) { m.elevation = ElevationFlat } }

// Raised sets ElevationRaised — chrome band + body slab. This is the default.
func Raised() Option { return func(m *Model) { m.elevation = ElevationRaised } }

// WithTheme sets the theme for this card instance.
// If not set, falls back to theme.CurrentTheme().
func WithTheme(t theme.Theme) Option {
	return func(m *Model) { m.theme = t }
}

func New(opts ...Option) *Model {
	m := &Model{inset: 0, elevation: ElevationRaised}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// SetTheme updates the theme. Call from your app's Update() on ThemeChangedMsg.
func (m *Model) SetTheme(t theme.Theme) { m.theme = t }

func (m *Model) Init() tea.Cmd {
	if m.content == nil {
		return nil
	}
	return m.content.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.SetSize(ws.Width, ws.Height)
		return m, nil
	}
	if !m.focused {
		if _, ok := msg.(tea.KeyMsg); ok {
			return m, nil
		}
	}
	if m.content == nil {
		return m, nil
	}
	u, cmd := m.content.Update(msg)
	m.content = u
	return m, cmd
}

func (m *Model) View() tea.View {
	if m.elevation == ElevationFlat {
		return m.viewFlat()
	}
	return m.viewRaised()
}

// viewRaised renders the chrome header band + body slab layout.
func (m *Model) viewRaised() tea.View {
	t := m.activeTheme()
	w, h := m.width, m.height
	if w <= 0 {
		w = 40
	}
	if h <= 0 {
		h = 8
	}

	outerBG := t.BackgroundPanel()
	chromeBG := t.CardChrome()
	bodyBG := t.CardBody()
	titleFG := t.CardFrameFG()
	bodyFG := t.Text()
	metaFG := t.TextMuted()
	footerFG := t.TextMuted()

	var leftEdgeBG color.Color
	if m.focused {
		leftEdgeBG = t.CardFocusEdge()
	} else {
		leftEdgeBG = t.CardChrome()
	}

	inset := clamp(m.inset, 0, min(max(0, (w-6)/2), max(0, (h-6)/2)))
	cardW := w - (inset * 2)
	cardH := h - (inset * 2)
	if cardW < 4 || cardH < 4 {
		inset = 0
		cardW = w
		cardH = h
	}

	body := ""
	if m.content != nil {
		body = viewStr(m.content.View())
	}
	bodyLines := strings.Split(body, "\n")
	for i := range bodyLines {
		bodyLines[i] = ansi.Strip(bodyLines[i])
	}

	hasMeta := strings.TrimSpace(m.meta) != ""
	hasFooter := strings.TrimSpace(m.footer) != ""
	reserved := 2
	if hasMeta {
		reserved++
	}
	if hasFooter {
		reserved++
	}
	if reserved > cardH && hasFooter {
		hasFooter = false
		reserved--
	}
	if reserved > cardH {
		hasMeta = false
		reserved = 2
	}
	contentRows := max(1, cardH-reserved)

	cardRows := make([]string, 0, cardH)
	cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, chromeBG, chromeBG, titleFG, " "+ansi.Strip(m.title)))
	if hasMeta {
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, chromeBG, chromeBG, metaFG, " "+ansi.Strip(m.meta)))
	}
	for i := 0; i < contentRows; i++ {
		line := ""
		if i < len(bodyLines) {
			line = bodyLines[i]
		}
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, bodyBG, chromeBG, bodyFG, " "+line))
	}
	if hasFooter {
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, chromeBG, chromeBG, footerFG, " "+ansi.Strip(m.footer)))
	}
	for len(cardRows) < cardH-1 {
		cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, bodyBG, chromeBG, bodyFG, ""))
	}
	cardRows = append(cardRows, slabRow(cardW, leftEdgeBG, chromeBG, chromeBG, titleFG, ""))

	rows := make([]string, 0, h)
	for y := 0; y < h; y++ {
		line := ""
		if inset > 0 {
			line += styles.RowClip(outerBG, bodyFG, inset, "")
		}
		inCardY := y >= inset && y < inset+cardH
		used := inset
		if inCardY {
			line += cardRows[y-inset]
			used += cardW
		}
		if used < w {
			line += styles.RowClip(outerBG, bodyFG, w-used, "")
		}
		rows = append(rows, line)
	}

	return tea.NewView(strings.Join(rows, "\n"))
}

// viewFlat renders a plain titled container with separator (former panel).
func (m *Model) viewFlat() tea.View {
	t := m.activeTheme()
	w, h := m.width, m.height
	if w <= 0 {
		w = 30
	}
	if h <= 0 {
		h = 6
	}

	var bg color.Color
	if m.focused {
		bg = t.BackgroundInteractive()
	} else {
		bg = t.BackgroundPanel()
	}
	fg := t.Text()

	rows := make([]string, 0, h)
	titleRows := 0

	if m.title != "" && h > 0 {
		var titleFG color.Color
		if m.focused {
			titleFG = t.TextAccent()
		} else {
			titleFG = t.TextMuted()
		}
		title := " " + ansi.Strip(m.title) + " "
		titleW := min(w, lipglossWidth(title))
		left := styles.RowClip(bg, titleFG, titleW, title)
		if titleW < w {
			left += styles.RowClip(bg, fg, w-titleW, "")
		}
		rows = append(rows, left)
		titleRows++

		if h > 1 {
			sep := styles.RowClip(bg, t.BorderSubtle(), w, strings.Repeat("─", w))
			rows = append(rows, sep)
			titleRows++
		}
	}

	body := ""
	if m.content != nil {
		body = viewStr(m.content.View())
	}
	contentLines := strings.Split(body, "\n")

	for len(rows) < h {
		idx := len(rows) - titleRows
		line := ""
		if idx >= 0 && idx < len(contentLines) {
			line = ansi.Strip(contentLines[idx])
		}
		if m.focused && w > 1 {
			accent := styles.RowClip(t.BorderFocus(), t.TextInverse(), 1, "")
			rest := styles.RowClip(bg, fg, w-1, " "+line)
			rows = append(rows, accent+rest)
		} else {
			rows = append(rows, styles.RowClip(bg, fg, w, " "+line))
		}
	}

	return tea.NewView(strings.Join(rows, "\n"))
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if s, ok := m.content.(interface{ SetSize(int, int) }); ok {
		titleRows := 0
		if m.title != "" {
			if m.elevation == ElevationFlat {
				titleRows = 2
			} else {
				titleRows = 1
				if strings.TrimSpace(m.meta) != "" {
					titleRows++
				}
			}
		}
		var contentH int
		if m.elevation == ElevationFlat {
			contentH = max(0, height-titleRows)
		} else {
			contentH = max(0, height-titleRows-1)
		}
		s.SetSize(max(0, width-2), max(0, contentH))
	}
}

func (m *Model) GetSize() (int, int) { return m.width, m.height }

func (m *Model) SetTitle(v string)  { m.title = v }
func (m *Model) SetMeta(v string)   { m.meta = v }
func (m *Model) SetFooter(v string) { m.footer = v }

func (m *Model) Focus()          { m.focused = true }
func (m *Model) Blur()           { m.focused = false }
func (m *Model) IsFocused() bool { return m.focused }

func (m *Model) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

// slabRow renders a row with a 1-cell left accent edge, center slab, 1-cell right edge.
func slabRow(width int, leftBG, centerBG, rightBG, fg color.Color, content string) string {
	if width <= 0 {
		return ""
	}
	if width == 1 {
		return styles.RowClip(centerBG, fg, 1, content)
	}
	left := styles.RowClip(leftBG, fg, 1, "")
	if width == 2 {
		return left + styles.RowClip(rightBG, fg, 1, "")
	}
	mid := styles.RowClip(centerBG, fg, width-2, content)
	right := styles.RowClip(rightBG, fg, 1, "")
	return left + mid + right
}

func viewStr(v tea.View) string {
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

// lipglossWidth returns the display width of a string (ANSI-safe).
func lipglossWidth(s string) int {
	return len([]rune(ansi.Strip(s)))
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
