package dialogcmp

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/dialog"
	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/theme"
)

type ThemePicker struct {
	allThemes []string
	filtered  []string
	selected  int
	width     int
	height    int
	themeName string
	search    textinput.Model
}

func NewThemePicker() *ThemePicker {
	names := theme.AvailableThemes()
	cur := theme.CurrentThemeName()
	in := textinput.New()
	in.Placeholder = "Search"
	in.Prompt = ""
	in.ShowSuggestions = false
	in.Focus()

	p := &ThemePicker{allThemes: names, filtered: append([]string(nil), names...), themeName: cur, search: in}
	p.alignSelectionToCurrent()
	p.syncStyles()
	return p
}

func (p *ThemePicker) Init() tea.Cmd { return p.search.Focus() }

func (p *ThemePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if changed, ok := msg.(theme.ThemeChangedMsg); ok {
		p.themeName = changed.Name
		p.alignSelectionToCurrent()
		return p, nil
	}

	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return p, nil
	}

	switch keyMsg.String() {
	case "esc":
		return p, func() tea.Msg { return dialog.Close() }
	case "up", "k":
		if p.selected > 0 {
			p.selected--
		}
		return p, nil
	case "down", "j":
		if p.selected < len(p.filtered)-1 {
			p.selected++
		}
		return p, nil
	case "enter":
		if len(p.filtered) == 0 {
			return p, nil
		}
		name := p.filtered[p.selected]
		t, err := theme.SetTheme(name)
		if err != nil {
			return p, nil
		}
		p.themeName = name
		return p, tea.Batch(
			func() tea.Msg { return theme.ThemeChangedMsg{Name: name, Theme: t} },
			func() tea.Msg { return dialog.Close() },
		)
	}

	updated, cmd := p.search.Update(keyMsg)
	p.search = updated
	p.refilter()
	return p, cmd
}

func (p *ThemePicker) View() tea.View {
	t := theme.CurrentTheme()
	sys := styles.New(t)
	rows := make([]string, 0, 10)
	rows = append(rows, "Search")
	rows = append(rows, inputContainer(p.search.View(), t))
	rows = append(rows, "")

	if len(p.filtered) == 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted)).Render("No matching themes"))
	} else {
		maxRows := max(5, p.height-10)
		start := 0
		if p.selected >= maxRows {
			start = p.selected - maxRows + 1
		}
		end := min(len(p.filtered), start+maxRows)
		for i := start; i < end; i++ {
			name := p.filtered[i]
			selected := i == p.selected
			marker := " "
			if name == p.themeName {
				marker = sys.CurrentMarker().Render("●")
			}
			line := fmt.Sprintf("%s %s", marker, name)
			rows = append(rows, sys.ListItem(selected).Width(max(28, p.width-8)).Render(line))
		}
	}

	rows = append(rows, "", lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted)).Render("enter apply  esc close"))
	return tea.NewView(strings.Join(rows, "\n"))
}

func (p *ThemePicker) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.syncStyles()
}

func (p *ThemePicker) Title() string { return "Themes" }

func (p *ThemePicker) syncStyles() {
	t := theme.CurrentTheme()
	s := textinput.DefaultStyles(true)
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
	s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text))
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))
	s.Blurred = s.Focused
	s.Cursor.Color = lipgloss.Color(t.Accent)
	p.search.SetStyles(s)
	p.search.SetWidth(max(16, p.width-12))
}

func (p *ThemePicker) refilter() {
	query := strings.ToLower(strings.TrimSpace(p.search.Value()))
	if query == "" {
		p.filtered = append([]string(nil), p.allThemes...)
		p.alignSelectionToCurrent()
		return
	}

	prevSelected := ""
	if len(p.filtered) > 0 && p.selected < len(p.filtered) {
		prevSelected = p.filtered[p.selected]
	}

	next := make([]string, 0, len(p.allThemes))
	for _, name := range p.allThemes {
		if strings.Contains(strings.ToLower(name), query) {
			next = append(next, name)
		}
	}
	p.filtered = next
	if len(p.filtered) == 0 {
		p.selected = 0
		return
	}
	for i, name := range p.filtered {
		if name == prevSelected {
			p.selected = i
			return
		}
	}
	p.selected = 0
}

func (p *ThemePicker) alignSelectionToCurrent() {
	if len(p.filtered) == 0 {
		p.selected = 0
		return
	}
	for i, name := range p.filtered {
		if name == p.themeName {
			p.selected = i
			return
		}
	}
	p.selected = 0
}

func inputContainer(view string, t theme.Theme) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(pick(t.InputBG, t.ElementBG, t.SurfaceMuted))).
		Foreground(lipgloss.Color(t.Text)).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(pick(t.InputBorder, t.BorderFocused, t.Border))).
		Padding(0, 1).
		Render(view)
}

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
