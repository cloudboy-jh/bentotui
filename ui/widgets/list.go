package widgets

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/core/theme"
)

// List displays a scrollable list of items.
type List struct {
	width  int
	height int
	items  []string
	max    int
	theme  theme.Theme
}

// NewList creates a list with optional max items.
func NewList(maxItems int) *List {
	if maxItems <= 0 {
		maxItems = 100
	}
	return &List{max: maxItems, theme: theme.CurrentTheme()}
}

// Append adds an item to the list.
func (l *List) Append(item string) {
	l.items = append(l.items, item)
	if len(l.items) > l.max {
		l.items = l.items[1:] // Remove oldest
	}
}

// Prepend adds an item to the beginning.
func (l *List) Prepend(item string) {
	l.items = append([]string{item}, l.items...)
	if len(l.items) > l.max {
		l.items = l.items[:l.max]
	}
}

// Clear removes all items.
func (l *List) Clear() {
	l.items = nil
}

func (l *List) Init() tea.Cmd { return nil }

func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l *List) View() tea.View {
	if l.width <= 0 || l.height <= 0 {
		return tea.NewView(strings.Join(l.items, "\n"))
	}

	// Show last N items that fit
	lines := make([]string, 0, l.height)
	start := len(l.items) - l.height
	if start < 0 {
		start = 0
	}

	baseStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(l.theme.Text.Primary))

	for i := start; i < len(l.items); i++ {
		line := l.items[i]
		if len(line) > l.width {
			line = line[:l.width]
		}
		lines = append(lines, baseStyle.Render(line))
	}

	return tea.NewView(strings.Join(lines, "\n"))
}

func (l *List) SetSize(width, height int) {
	l.width = width
	l.height = height
}

func (l *List) GetSize() (int, int) {
	return l.width, l.height
}

// SetTheme updates the theme.
func (l *List) SetTheme(t theme.Theme) {
	l.theme = t
}

// HeightConstraint returns FlexHeight(1) - list wants to fill remaining space.
func (l *List) HeightConstraint(width int) HeightConstraint {
	return FlexHeight(1)
}

var _ core.Component = (*List)(nil)
var _ core.Sizeable = (*List)(nil)
var _ HeightConstraintAware = (*List)(nil)
