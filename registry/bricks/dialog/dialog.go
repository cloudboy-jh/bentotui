// Brick: Dialog Manager
// +-------------------------------+
// |      +---------------+        |
// |      | dialog frame  |        |
// |      +---------------+        |
// +-------------------------------+
// Centers and routes modal dialogs.
// Copy this directory into your project: bento add dialog
package dialog

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/theme"
)

// Dialog is implemented by any modal that the Manager can host.
type Dialog interface {
	tea.Model
	SetSize(width, height int)
	Title() string
}

// OpenMsg signals the Manager to open a dialog.
type OpenMsg struct{ Dialog Dialog }

// CloseMsg signals the Manager to close the active dialog.
type CloseMsg struct{}

func Open(dialog Dialog) tea.Msg { return OpenMsg{Dialog: dialog} }
func Close() tea.Msg             { return CloseMsg{} }

// Manager hosts zero or one active Dialog.
type Manager struct {
	active Dialog
	width  int
	height int
	theme  theme.Theme // nil = use theme.CurrentTheme()
}

func New() *Manager { return &Manager{} }

// SetTheme updates the theme. Call on ThemeChangedMsg.
func (m *Manager) SetTheme(t theme.Theme) { m.theme = t }

func (m *Manager) activeTheme() theme.Theme {
	if m.theme != nil {
		return m.theme
	}
	return theme.CurrentTheme()
}

func (m *Manager) Init() tea.Cmd { return nil }

func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case OpenMsg:
		m.active = v.Dialog
		if m.active != nil {
			m.active = resizeDialog(m.active, m.width, m.height)
		}
		return m, nil
	case CloseMsg:
		m.active = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(v.Width, v.Height)
		return m, nil
	}

	if m.active != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				switch m.active.(type) {
				case Confirm, *Confirm:
					m.active = nil
					return m, nil
				}
			case "enter":
				switch active := m.active.(type) {
				case Confirm:
					m.active = nil
					if active.OnConfirm != nil {
						return m, func() tea.Msg { return active.OnConfirm() }
					}
					return m, nil
				case *Confirm:
					m.active = nil
					if active != nil && active.OnConfirm != nil {
						return m, func() tea.Msg { return active.OnConfirm() }
					}
					return m, nil
				}
			}
		}
	}

	if m.active == nil {
		return m, nil
	}
	updated, cmd := m.active.Update(msg)
	if next, ok := updated.(Dialog); ok {
		m.active = next
	}
	return m, cmd
}

func (m *Manager) View() tea.View {
	if m.active == nil {
		return tea.NewView("")
	}
	return tea.NewView(viewString(m.active.View()))
}

func (m *Manager) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.active != nil {
		m.active = resizeDialog(m.active, width, height)
	}
}

func (m *Manager) IsOpen() bool { return m.active != nil }

// ── Confirm ───────────────────────────────────────────────────────────────────

type Confirm struct {
	DialogTitle string
	Message     string
	OnConfirm   func() tea.Msg
	width       int
	height      int
}

func (c Confirm) Init() tea.Cmd                           { return nil }
func (c Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return c, nil }

func (c Confirm) View() tea.View {
	text := c.Message
	if strings.TrimSpace(text) == "" {
		text = "Confirm?"
	}
	content := strings.Join([]string{text, "", "Enter confirm  Esc cancel"}, "\n")
	view := renderDialogFrame(c.DialogTitle, content, max(48, c.width/2), 0, theme.CurrentTheme())
	return tea.NewView(view)
}

func (c Confirm) SetSize(width, height int) { c.width = width; c.height = height }
func (c Confirm) Title() string             { return c.DialogTitle }

// ── Custom ────────────────────────────────────────────────────────────────────

type Custom struct {
	DialogTitle string
	Content     tea.Model
	Width       int
	Height      int
}

func (c Custom) Init() tea.Cmd {
	if c.Content == nil {
		return nil
	}
	return c.Content.Init()
}

func (c Custom) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if c.Content == nil {
		return c, nil
	}
	updated, cmd := c.Content.Update(msg)
	c.Content = updated
	return c, cmd
}

func (c Custom) View() tea.View {
	body := ""
	if c.Content != nil {
		body = viewString(c.Content.View())
	}
	view := renderDialogFrame(c.DialogTitle, body, max(50, c.Width), max(12, c.Height), theme.CurrentTheme())
	return tea.NewView(view)
}

func (c Custom) SetSize(width, height int) { applyCustomSize(&c, width, height) }
func (c Custom) Title() string             { return c.DialogTitle }

func resizeDialog(d Dialog, width, height int) Dialog {
	switch v := d.(type) {
	case Confirm:
		v.width = width
		v.height = height
		return v
	case *Confirm:
		if v != nil {
			v.width = width
			v.height = height
		}
		return v
	case Custom:
		applyCustomSize(&v, width, height)
		return v
	case *Custom:
		if v != nil {
			applyCustomSize(v, width, height)
		}
		return v
	default:
		d.SetSize(width, height)
		return d
	}
}

func applyCustomSize(c *Custom, width, height int) {
	if c == nil {
		return
	}
	if c.Width == 0 {
		c.Width = width / 2
	}
	if c.Height == 0 {
		c.Height = height / 2
	}
	maxWidth := max(20, width-4)
	maxHeight := max(8, height-4)
	minWidth := min(36, maxWidth)
	minHeight := min(12, maxHeight)
	c.Width = clamp(c.Width, minWidth, maxWidth)
	c.Height = clamp(c.Height, minHeight, maxHeight)
	if s, ok := c.Content.(interface{ SetSize(int, int) }); ok {
		s.SetSize(max(1, c.Width-4), max(1, c.Height-4))
	}
}

// ── rendering ─────────────────────────────────────────────────────────────────

func renderDialogFrame(title, content string, width, height int, t theme.Theme) string {
	bg := t.DialogBG()
	fg := t.DialogFG()
	mutedFG := t.TextMuted()
	accentFG := t.TextAccent()

	if strings.TrimSpace(title) == "" {
		title = "Dialog"
	}
	if width <= 0 {
		width = 48
	}
	if height <= 0 {
		height = 14
	}

	innerWidth := max(1, width-4)

	base := lipgloss.NewStyle().Background(bg)
	mkRow := func(rowFG color.Color, rowContent string) string {
		return base.
			Foreground(rowFG).
			PaddingLeft(2).PaddingRight(2).
			Width(innerWidth).
			Render(rowContent)
	}
	blankRow := base.Width(width).Render("")

	// Header: title left, "esc" right
	rightWidth := 3
	leftWidth := max(1, innerWidth-rightWidth)
	titleCell := base.Foreground(accentFG).Bold(true).Width(leftWidth).Render(title)
	escCell := base.Foreground(mutedFG).Width(rightWidth).Render("esc")
	header := mkRow(fg, titleCell+escCell)

	body := strings.TrimRight(content, "\n")
	if strings.TrimSpace(body) == "" {
		body = " "
	}
	rawBodyLines := clipLines(body, innerWidth, bg, fg)

	bodyMax := max(1, height-4)
	if len(rawBodyLines) > bodyMax {
		rawBodyLines = rawBodyLines[:bodyMax]
	}
	bodyRows := make([]string, len(rawBodyLines))
	for i, line := range rawBodyLines {
		bodyRows[i] = mkRow(fg, line)
	}

	allRows := make([]string, 0, height)
	allRows = append(allRows, blankRow)
	allRows = append(allRows, header)
	allRows = append(allRows, blankRow)
	allRows = append(allRows, bodyRows...)
	allRows = append(allRows, blankRow)
	for len(allRows) < height {
		allRows = append(allRows, blankRow)
	}

	return base.Width(width).Render(strings.Join(allRows, "\n"))
}

func clipLines(content string, width int, bg, fg color.Color) []string {
	lines := strings.Split(content, "\n")
	clipped := make([]string, 0, len(lines))
	for _, line := range lines {
		clipped = append(clipped, lipgloss.NewStyle().
			Background(bg).
			Foreground(fg).
			Width(width).
			MaxWidth(width).
			Render(line))
	}
	return clipped
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

func clamp(v, minV, maxV int) int {
	if maxV < minV {
		return minV
	}
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}
