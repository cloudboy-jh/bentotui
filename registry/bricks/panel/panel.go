// Brick: Panel:
// +-----------------------------------+
// | title bar                          |
// | separator                          |
// | content rows                       |
// +-----------------------------------+
// Focusable surface container.
// Package panel provides a themed, titled, focusable content container.
// Copy this file into your project: bento add panel
//
// Dependencies (real Go module imports, not copied):
//   - charm.land/bubbletea/v2
//   - charm.land/lipgloss/v2
//   - github.com/cloudboy-jh/bentotui/theme
//   - github.com/cloudboy-jh/bentotui/theme/styles
package panel

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/theme/styles"
)

// Model is a themed, titled, focusable content container.
//
// Layout (top to bottom):
//
//	title bar row  — Interactive bg, badge on left  (1 row, only when Title set)
//	separator row  — Border.Normal fg, full-width ─── (1 row, only when Title set)
//	content rows   — Panel bg, focus-accent left stripe when focused
type Model struct {
	title    string
	content  tea.Model
	elevated bool
	focused  bool
	width    int
	height   int
}

// Option configures a Model at construction time.
type Option func(*Model)

func Title(title string) Option  { return func(m *Model) { m.title = title } }
func Content(c tea.Model) Option { return func(m *Model) { m.content = c } }
func Elevated() Option           { return func(m *Model) { m.elevated = true } }
func Scrollable(_ bool) Option   { return func(_ *Model) {} } // reserved

// New constructs a panel with the given options.
func New(opts ...Option) *Model {
	m := &Model{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// ── tea.Model ─────────────────────────────────────────────────────────────────

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
	updated, cmd := m.content.Update(msg)
	m.content = updated
	return m, cmd
}

func (m *Model) View() tea.View {
	// Always read theme at render time — never cache it.
	t := theme.CurrentTheme()
	sys := styles.New(t)

	w := m.width
	h := m.height
	if w <= 0 {
		w = 30
	}
	if h <= 0 {
		h = 6
	}

	var panelBG string
	switch {
	case m.elevated:
		panelBG = pick(t.Surface.Elevated, t.Surface.Panel)
	case m.focused:
		panelBG = pick(t.Surface.Interactive, t.Surface.Panel)
	default:
		panelBG = pick(t.Surface.Panel, t.Surface.Elevated)
	}
	panelFG := t.Text.Primary

	rows := make([]string, 0, h)

	// ── title bar + separator ─────────────────────────────────────────────────
	titleRows := 0
	if m.title != "" && h > 0 {
		titleBar := renderTitleRow(m.title, w, panelBG, panelFG, m.focused, t)
		rows = append(rows, titleBar)
		titleRows++

		if h > 1 {
			sep := styles.RowClip(panelBG, pick(t.Border.Subtle, t.Text.Muted), w, strings.Repeat("─", w))
			rows = append(rows, sep)
			titleRows++
		}
	}

	// ── content rows ──────────────────────────────────────────────────────────
	body := ""
	if m.content != nil {
		body = viewString(m.content.View())
	}
	contentLines := strings.Split(body, "\n")

	for len(rows) < h {
		idx := len(rows) - titleRows
		line := ""
		if idx >= 0 && idx < len(contentLines) {
			line = contentLines[idx]
		}
		rows = append(rows, contentRow(line, w, panelBG, panelFG, m.focused, sys))
	}

	return tea.NewView(strings.Join(rows, "\n"))
}

// ── sizing ────────────────────────────────────────────────────────────────────

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	titleRows := 0
	if m.title != "" {
		titleRows = 2 // title bar + separator
	}
	contentW := max(0, width-2) // 1-cell left + right margin
	contentH := max(0, height-titleRows)
	if s, ok := m.content.(interface{ SetSize(int, int) }); ok {
		s.SetSize(contentW, contentH)
	}
}

func (m *Model) GetSize() (int, int) { return m.width, m.height }

// ── focus ─────────────────────────────────────────────────────────────────────

func (m *Model) Focus()          { m.focused = true }
func (m *Model) Blur()           { m.focused = false }
func (m *Model) IsFocused() bool { return m.focused }

// ── internal helpers ──────────────────────────────────────────────────────────

// contentRow renders one content row with guaranteed full-width background.
// Content is stripped of ANSI before styling so the single Width().Render()
// call owns every cell's color, preventing background bleed-through.
func contentRow(line string, w int, bg, fg string, focused bool, sys styles.System) string {
	if w <= 0 {
		return ""
	}
	plain := styles.ClipANSI(ansi.Strip(line), max(0, w-2))
	if focused && w > 1 {
		accent := sys.FocusAccent().Render(" ")
		rest := styles.RowClip(bg, fg, w-1, " "+plain)
		return accent + rest
	}
	return styles.RowClip(bg, fg, w, " "+plain)
}

func renderTitleRow(title string, width int, bg, fg string, focused bool, t theme.Theme) string {
	if width <= 0 {
		return ""
	}
	badgeFG := pick(t.Text.Muted, fg)
	if focused {
		badgeFG = pick(t.Text.Accent, fg)
	}
	badge := " " + ansi.Strip(title) + " "
	badgeW := min(width, len(badge))
	left := styles.RowClip(bg, badgeFG, badgeW, badge)
	if badgeW == width {
		return left
	}
	right := styles.RowClip(bg, fg, width-badgeW, "")
	return left + right
}

// viewString extracts a plain string from a tea.View.
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
