package panel

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

// Option configures a Model at construction time.
type Option func(*Model)

// Model is a themed, titled, focusable content container.
// It composes three layers from top to bottom:
//
//	title bar row  — Surface.Interactive bg, badge on left  (1 row, optional)
//	separator row  — Border.Normal fg, full-width ───        (1 row, when title set)
//	content rows   — Surface.Panel bg, left-edge accent when focused
type Model struct {
	title    string
	content  core.Component
	elevated bool // use Surface.Elevated instead of Surface.Panel
	focused  bool
	theme    theme.Theme
	width    int
	height   int
}

func New(opts ...Option) *Model {
	m := &Model{theme: theme.Preset(theme.DefaultName)}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// ── options ───────────────────────────────────────────────────────────────────

func Title(title string) Option       { return func(m *Model) { m.title = title } }
func Content(c core.Component) Option { return func(m *Model) { m.content = c } }
func Elevated() Option                { return func(m *Model) { m.elevated = true } }
func Scrollable(_ bool) Option        { return func(_ *Model) {} } // reserved
func Theme(t theme.Theme) Option      { return func(m *Model) { m.theme = t } }

// Border is kept for API compatibility — border rendering is handled internally.
func Border(_ any) Option { return func(_ *Model) {} }

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
	if next, ok := updated.(core.Component); ok {
		m.content = next
	}
	return m, cmd
}

func (m *Model) View() tea.View {
	sys := styles.New(m.theme)

	w := m.width
	h := m.height
	if w <= 0 {
		w = 30
	}
	if h <= 0 {
		h = 6
	}

	// Choose surface background directly from theme tokens.
	var panelBG string
	if m.elevated {
		panelBG = pick(m.theme.Surface.Elevated, m.theme.Surface.Panel)
	} else if m.focused {
		panelBG = pick(m.theme.Surface.Interactive, m.theme.Surface.Panel)
	} else {
		panelBG = pick(m.theme.Surface.Panel, m.theme.Surface.Elevated)
	}

	rows := make([]string, 0, h)

	// ── title bar (row 0) + separator (row 1) ────────────────────────────────
	titleRows := 0
	if m.title != "" && h > 0 {
		// Title bar: full-width row with badge rendered via PanelTitleBadge.
		badge := sys.PanelTitleBadge(m.focused).Render(m.title)
		titleBar := primitives.RenderStyledRow(sys.PanelTitleBar(m.focused), w, badge)
		rows = append(rows, titleBar)
		titleRows++

		// Separator: ─── line in Border.Normal color.
		if h > 1 {
			sep := sys.Divider().Render(strings.Repeat("─", w))
			rows = append(rows, sep)
			titleRows++
		}
	}

	// ── content rows ─────────────────────────────────────────────────────────
	body := ""
	if m.content != nil {
		body = core.ViewString(m.content.View())
	}
	contentLines := strings.Split(body, "\n")

	for len(rows) < h {
		idx := len(rows) - titleRows
		line := ""
		if idx >= 0 && idx < len(contentLines) {
			line = contentLines[idx]
		}
		rows = append(rows, m.contentRow(line, w))
	}

	// ── assemble with explicit background fill ────────────────────────────────
	frameStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(panelBG)).
		Foreground(lipgloss.Color(m.theme.Text.Primary)).
		Width(w).
		Height(h)

	return tea.NewView(primitives.RenderFrame(frameStyle, w, h, rows))
}

// contentRow renders one content row with a 1-cell left margin (or focus
// accent stripe when focused) and 1-cell right margin.
func (m *Model) contentRow(line string, w int) string {
	if w <= 0 {
		return ""
	}
	if w == 1 {
		return " "
	}
	if m.focused && w > 2 {
		// Focus accent: 1-cell Border.Focus stripe on the left edge.
		sys := styles.New(m.theme)
		accent := sys.FocusAccent().Render(" ")
		content := primitives.FitWidth(line, w-2)
		return accent + content + " "
	}
	// Normal: 1-space left margin, content, 1-space right margin.
	return " " + primitives.FitWidth(line, w-2) + " "
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
	if s, ok := m.content.(core.Sizeable); ok {
		s.SetSize(contentW, contentH)
	}
}

func (m *Model) GetSize() (int, int) { return m.width, m.height }

// ── focus + theme ─────────────────────────────────────────────────────────────

func (m *Model) Focus()          { m.focused = true }
func (m *Model) Blur()           { m.focused = false }
func (m *Model) IsFocused() bool { return m.focused }
func (m *Model) SetTheme(t theme.Theme) {
	m.theme = t
	if s, ok := m.content.(interface{ SetTheme(theme.Theme) }); ok {
		s.SetTheme(t)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// pick returns v if non-empty, otherwise fallback.
func pick(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
