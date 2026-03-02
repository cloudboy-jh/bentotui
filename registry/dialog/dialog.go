// Package dialog provides modal dialog management for bentotui apps.
// Copy this directory into your project: bento add dialog
//
// Dependencies (real Go module imports, not copied):
//   - charm.land/bubbletea/v2
//   - charm.land/lipgloss/v2
//   - github.com/cloudboy-jh/bentotui/theme
//   - github.com/cloudboy-jh/bentotui/styles
package dialog

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/styles"
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

// Manager hosts zero or one active Dialog. Place it in your root model and
// render it as a canvas overlay above the rest of your app.
type Manager struct {
	active Dialog
	width  int
	height int
}

func New() *Manager { return &Manager{} }

func (m *Manager) Init() tea.Cmd { return nil }

func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case OpenMsg:
		m.active = v.Dialog
		if m.active != nil {
			m.active.SetSize(m.width, m.height)
		}
		return m, nil
	case CloseMsg:
		m.active = nil
		return m, nil
	case tea.WindowSizeMsg:
		m.SetSize(v.Width, v.Height)
		return m, nil
	}

	// Auto-handle Confirm dialogs for esc/enter.
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
		m.active.SetSize(width, height)
	}
}

func (m *Manager) IsOpen() bool { return m.active != nil }

// ── Confirm ───────────────────────────────────────────────────────────────────

// Confirm is a simple yes/no dialog. The Manager handles enter/esc automatically.
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

func (c Confirm) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c Confirm) Title() string { return c.DialogTitle }

// ── Custom ────────────────────────────────────────────────────────────────────

// Custom wraps any tea.Model as a dialog with a title frame.
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

func (c Custom) SetSize(width, height int) {
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

func (c Custom) Title() string { return c.DialogTitle }

// ── rendering ─────────────────────────────────────────────────────────────────

func renderDialogFrame(title, content string, width, height int, t theme.Theme) string {
	sys := styles.New(t)
	headerTitle := title
	if strings.TrimSpace(headerTitle) == "" {
		headerTitle = "Dialog"
	}
	if width <= 0 {
		width = 48
	}
	if height <= 0 {
		height = 14
	}
	innerWidth := max(1, width-4)
	innerHeight := max(1, height-2)

	rightWidth := 3
	leftWidth := max(1, innerWidth-rightWidth-1)
	header := sys.DialogHeader().Render(fitWidth(headerTitle, leftWidth)) + " " + sys.DialogEscHint().Render(fitWidth("esc", rightWidth))
	header = fitWidth(header, innerWidth)

	body := strings.TrimRight(content, "\n")
	if strings.TrimSpace(body) == "" {
		body = " "
	}
	bodyLines := clipLines(body, innerWidth)
	bodyMax := max(1, innerHeight-2)
	if len(bodyLines) > bodyMax {
		bodyLines = bodyLines[:bodyMax]
	}
	joined := strings.Join(append([]string{header, ""}, bodyLines...), "\n")
	style := sys.DialogFrame().Width(width).Height(height)
	return style.Render(joined)
}

func clipLines(content string, width int) []string {
	lines := strings.Split(content, "\n")
	clipped := make([]string, 0, len(lines))
	for _, line := range lines {
		clipped = append(clipped, fitWidth(line, width))
	}
	return clipped
}

// fitWidth clips s to MaxWidth then places it left-aligned at exact width.
func fitWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	s = lipgloss.NewStyle().MaxWidth(width).Render(s)
	return lipgloss.PlaceHorizontal(width, lipgloss.Left, s)
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

// ── helpers ───────────────────────────────────────────────────────────────────

func pick(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
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
