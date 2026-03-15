// Brick: Theme Picker:
// +-----------------------------------+
// | themes list                        |
// | live preview on move               |
// +-----------------------------------+
// Interactive theme selection dialog.
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
	sys := styles.New(t)
	contentWidth := maxv(24, p.width)
	rows := make([]string, 0, 10)

	// Search row — rendered entirely with Dialog.BG so no dark cells bleed through.
	// We build the display manually from Value()/Position() instead of using
	// p.search.View(), which may produce cells with Bg=nil or a dark default Bg.
	searchRow := p.renderSearchRow(t, sys, contentWidth)
	rows = append(rows, searchRow)
	// Blank separator using dialog bg
	rows = append(rows, sys.DialogListRow().Width(contentWidth).Render(""))

	if len(p.filtered) == 0 {
		rows = append(rows, sys.DialogListRow().Width(contentWidth).Render("  No matching themes"))
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

			var row string
			if selected {
				// Selected row: Selection.BG/FG — always highest contrast
				rowStyle := sys.DialogListRowSelected()
				if isCurrent {
					bullet := rowStyle.Render("● ")
					rest := rowStyle.Width(contentWidth - 2).Render(name)
					row = bullet + rest
				} else {
					row = rowStyle.Width(contentWidth).Render("  " + name)
				}
			} else if isCurrent {
				// Current (not selected): accent bullet, dialog bg
				bullet := lipgloss.NewStyle().
					Foreground(lipgloss.Color(t.Text.Accent)).
					Background(lipgloss.Color(t.Dialog.BG)).
					Render("● ")
				rest := sys.DialogListRow().Width(contentWidth - 2).Render(name)
				row = bullet + rest
			} else {
				row = sys.DialogListRow().Width(contentWidth).Render("  " + name)
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

// renderSearchRow builds the search input row entirely from Dialog.BG-backed
// styles. We never use p.search.View() because textinput may produce cells with
// Bg=nil or a dark default background that bleeds through the UV surface overlay.
func (p *ThemePicker) renderSearchRow(t theme.Theme, sys styles.System, width int) string {
	dbg := t.Dialog.BG
	dfg := t.Text.Primary
	muted := t.Text.Muted
	cursor := t.Input.Cursor

	val := []rune(p.search.Value())
	pos := p.search.Position()
	if pos > len(val) {
		pos = len(val)
	}

	// Available text width = width - 2 (leading spaces).
	textW := maxv(1, width-2)

	base := lipgloss.NewStyle().Background(lipgloss.Color(dbg))
	textStyle := base.Foreground(lipgloss.Color(dfg))
	mutedStyle := base.Foreground(lipgloss.Color(muted))
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(cursor)).
		Foreground(lipgloss.Color(dbg))

	var result string
	if len(val) == 0 {
		// Placeholder: cursor on first char, rest muted.
		ph := []rune(p.search.Placeholder)
		if len(ph) == 0 {
			ph = []rune{' '}
		}
		cursorChar := string(ph[0])
		rest := ""
		if len(ph) > 1 {
			rest = string(ph[1:])
		}
		result = cursorStyle.Render(cursorChar) + mutedStyle.Render(rest)
	} else {
		before := string(val[:pos])
		after := ""
		cursorChar := " " // cursor at end — block on space
		if pos < len(val) {
			cursorChar = string(val[pos])
			after = string(val[pos+1:])
		}
		result = textStyle.Render(before) + cursorStyle.Render(cursorChar) + textStyle.Render(after)
	}

	// Clip + pad to textW, then add the 2-space indent — all on Dialog.BG.
	clipped := lipgloss.NewStyle().
		Background(lipgloss.Color(dbg)).
		MaxWidth(textW).
		Width(textW).
		Render(result)

	return sys.DialogListRow().Width(width).Render("  " + clipped)
}

func (p *ThemePicker) syncStyles() {
	t := theme.CurrentTheme()
	// The textinput model handles key input only — View() is never called on it.
	// Minimal styles: just keep cursor color for any internal blinking state.
	s := textinput.DefaultStyles(true)
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
