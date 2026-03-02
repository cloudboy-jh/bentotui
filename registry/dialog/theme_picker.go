package dialog

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/cloudboy-jh/bentotui/styles"
	"github.com/cloudboy-jh/bentotui/theme"
)

// ThemePicker is a searchable theme selection dialog.
// It previews themes live as the user navigates and reverts on ESC.
type ThemePicker struct {
	allThemes []string
	filtered  []string
	selected  int
	width     int
	height    int
	themeName string // current preview/selection
	baseTheme string // theme at dialog-open time — reverted on ESC
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

	p := &ThemePicker{
		allThemes: names,
		filtered:  append([]string(nil), names...),
		themeName: cur,
		baseTheme: cur,
		search:    in,
	}
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
		t, err := theme.PreviewTheme(p.baseTheme)
		if err != nil {
			return p, func() tea.Msg { return Close() }
		}
		p.themeName = p.baseTheme
		p.alignSelectionToCurrent()
		p.syncStyles()
		return p, tea.Batch(
			func() tea.Msg { return theme.ThemeChangedMsg{Name: p.baseTheme, Theme: t} },
			func() tea.Msg { return Close() },
		)
	case "up", "k":
		if p.selected > 0 {
			p.selected--
			return p, p.previewSelectedCmd()
		}
		return p, nil
	case "down", "j":
		if p.selected < len(p.filtered)-1 {
			p.selected++
			return p, p.previewSelectedCmd()
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
		p.baseTheme = name
		p.syncStyles()
		return p, tea.Batch(
			func() tea.Msg { return theme.ThemeChangedMsg{Name: name, Theme: t} },
			func() tea.Msg { return Close() },
		)
	}

	updated, cmd := p.search.Update(keyMsg)
	p.search = updated
	before := p.selected
	p.refilter()
	if p.selected != before {
		return p, tea.Batch(cmd, p.previewSelectedCmd())
	}
	return p, cmd
}

func (p *ThemePicker) View() tea.View {
	t := theme.CurrentTheme()
	contentWidth := maxv(24, p.width)
	rows := make([]string, 0, 10)

	dialogBG := t.Dialog.BG
	inputColors := styles.New(t).InputColors()

	// Search input row
	inputContent := fitWidth(p.search.View(), maxv(1, contentWidth-2))
	rows = append(rows, renderRow(contentWidth, inputColors.BG, inputColors.FG, " "+inputContent))
	// Blank separator
	rows = append(rows, renderRow(contentWidth, dialogBG, "", ""))

	if len(p.filtered) == 0 {
		rows = append(rows, renderRow(contentWidth, dialogBG, t.Text.Muted, "No matching themes"))
	} else {
		maxRows := maxv(1, p.height-4)
		start := 0
		if p.selected >= maxRows {
			start = p.selected - maxRows + 1
		}
		end := minv(len(p.filtered), start+maxRows)
		for i := start; i < end; i++ {
			name := p.filtered[i]
			selected := i == p.selected
			isCurrent := name == p.themeName

			rowBG := dialogBG
			rowFG := t.Text.Primary
			if selected {
				rowBG = pick(t.Border.Focus, t.Selection.BG)
				rowFG = pick(t.Text.Inverse, t.Selection.FG)
			}

			var row string
			if isCurrent {
				bulletStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(pick(t.Text.Accent, t.Border.Focus))).
					Background(lipgloss.Color(rowBG))
				nameStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(rowFG)).
					Background(lipgloss.Color(rowBG))
				row = bulletStyle.Render("● ") + nameStyle.Width(contentWidth-2).Render(name)
			} else {
				marker := "  "
				if selected {
					marker = "> "
				}
				row = lipgloss.NewStyle().
					Foreground(lipgloss.Color(rowFG)).
					Background(lipgloss.Color(rowBG)).
					Width(contentWidth).
					Render(marker + name)
			}
			rows = append(rows, row)
		}
	}

	return tea.NewView(strings.Join(rows, "\n"))
}

func (p *ThemePicker) SetSize(width, height int) {
	p.width = maxv(1, width)
	p.height = maxv(1, height)
	p.syncStyles()
}

func (p *ThemePicker) GetSize() (int, int) { return p.width, p.height }
func (p *ThemePicker) Title() string       { return "Themes" }

func (p *ThemePicker) syncStyles() {
	t := theme.CurrentTheme()
	s := textinput.DefaultStyles(true)
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Primary))
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Input.Placeholder))
	s.Blurred = s.Focused
	s.Cursor.Color = lipgloss.Color(t.Input.Cursor)
	p.search.SetStyles(s)
	p.search.SetWidth(maxv(10, p.width-2))
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
	for i, name := range p.filtered {
		if name == p.themeName {
			p.selected = i
			return
		}
	}
	p.selected = 0
}

func (p *ThemePicker) previewSelectedCmd() tea.Cmd {
	if len(p.filtered) == 0 || p.selected < 0 || p.selected >= len(p.filtered) {
		return nil
	}
	name := p.filtered[p.selected]
	if name == p.themeName {
		return nil
	}
	t, err := theme.PreviewTheme(name)
	if err != nil {
		return nil
	}
	p.themeName = name
	p.syncStyles()
	return func() tea.Msg { return theme.ThemeChangedMsg{Name: name, Theme: t} }
}

// renderRow renders a full-width row with guaranteed background fill.
func renderRow(width int, bg, fg, content string) string {
	if width <= 0 {
		return ""
	}
	style := lipgloss.NewStyle()
	if bg != "" {
		style = style.Background(lipgloss.Color(bg))
	}
	if fg != "" {
		style = style.Foreground(lipgloss.Color(fg))
	}
	return style.Width(width).Render(content)
}

func maxv(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minv(a, b int) int {
	if a < b {
		return a
	}
	return b
}
