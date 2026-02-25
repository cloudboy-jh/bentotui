package dialog

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/core/surface"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/styles"
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
		return p, func() tea.Msg { return Close() }
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
			func() tea.Msg { return Close() },
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
	contentWidth := maxInt(24, p.width)
	rows := make([]string, 0, 10)
	rows = append(rows, surface.FitWidth("Search", contentWidth))
	rows = append(rows, inputContainer(p.search.View(), contentWidth, t))
	rows = append(rows, "")

	if len(p.filtered) == 0 {
		rows = append(rows, lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted)).Render(surface.FitWidth("No matching themes", contentWidth)))
	} else {
		maxRows := maxInt(1, p.height-5)
		start := 0
		if p.selected >= maxRows {
			start = p.selected - maxRows + 1
		}
		end := minInt(len(p.filtered), start+maxRows)
		for i := start; i < end; i++ {
			name := p.filtered[i]
			selected := i == p.selected
			marker := " "
			if name == p.themeName {
				marker = sys.CurrentMarker().Render("●")
			}
			line := fmt.Sprintf("%s %s", marker, name)
			rows = append(rows, sys.ListItem(selected).Width(contentWidth).Render(surface.FitWidth(line, contentWidth)))
		}
	}

	rows = append(rows, "", lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted)).Render(surface.FitWidth("enter apply  esc close", contentWidth)))
	return tea.NewView(strings.Join(rows, "\n"))
}

func (p *ThemePicker) SetSize(width, height int) {
	p.width = maxInt(1, width)
	p.height = maxInt(1, height)
	p.syncStyles()
}

func (p *ThemePicker) GetSize() (int, int) { return p.width, p.height }

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
	p.search.SetWidth(maxInt(10, p.width-2))
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

func inputContainer(view string, width int, t theme.Theme) string {
	return lipgloss.NewStyle().
		Width(width).
		Background(lipgloss.Color(pick(t.InputBG, t.ElementBG, t.SurfaceMuted))).
		Foreground(lipgloss.Color(t.Text)).
		Padding(0, 1).
		Render(surface.FitWidth(view, maxInt(1, width-2)))
}

func pick(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
