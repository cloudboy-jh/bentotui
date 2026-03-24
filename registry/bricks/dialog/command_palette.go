// Brick: Command Palette
// +-----------------------------------+
// | search row                        |
// +-----------------------------------+
// | grouped command list              |
// +-----------------------------------+
// Searchable grouped action picker.
package dialog

import (
	"image/color"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/cloudboy-jh/bentotui/theme"
)

// Command is a single entry in the command palette.
type Command struct {
	Label   string // e.g. "Switch theme"
	Group   string // e.g. "System" — empty = "Commands"
	Keybind string // e.g. "ctrl+t" — optional, shown right-aligned
	Action  func() tea.Msg
}

// CommandPalette is a searchable, grouped command picker dialog.
type CommandPalette struct {
	commands []Command
	filtered []Command
	groups   []string
	selected int
	search   textinput.Model
	width    int
	height   int
}

// NewCommandPalette creates a CommandPalette pre-loaded with commands.
func NewCommandPalette(commands []Command) *CommandPalette {
	in := textinput.New()
	in.Placeholder = "Search"
	in.Prompt = ""
	in.ShowSuggestions = false
	in.Focus()

	p := &CommandPalette{
		commands: commands,
		selected: 0,
		search:   in,
	}
	p.refilter()
	p.syncStyles()
	return p
}

func (p *CommandPalette) Init() tea.Cmd { return p.search.Focus() }

func (p *CommandPalette) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		cmd := p.filtered[p.selected]
		if cmd.Action == nil {
			return p, func() tea.Msg { return Close() }
		}
		action := cmd.Action
		return p, tea.Sequence(
			func() tea.Msg { return Close() },
			func() tea.Msg { return action() },
		)
	}

	updated, cmd := p.search.Update(keyMsg)
	p.search = updated
	p.refilter()
	return p, cmd
}

func (p *CommandPalette) View() tea.View {
	t := theme.CurrentTheme()
	contentWidth := maxv(24, p.width)

	dbg := t.DialogBG()
	dfg := t.DialogFG()
	muted := t.TextMuted()

	baseRow := func(fg color.Color, content string) string {
		return lipgloss.NewStyle().
			Background(dbg).
			Foreground(fg).
			Width(contentWidth).
			Render(content)
	}

	rows := make([]string, 0, 16)

	// Search input row
	inputContent := paletteClip(p.search.View(), maxv(1, contentWidth-2))
	rows = append(rows, baseRow(dfg, " "+inputContent))
	rows = append(rows, baseRow(muted, ""))

	if len(p.filtered) == 0 {
		rows = append(rows, baseRow(muted, "  No matching commands"))
	} else {
		maxVisible := maxv(1, p.height-4)

		type entry struct {
			isGroup bool
			label   string
			keybind string
			idx     int
		}
		entries := make([]entry, 0, len(p.filtered)+len(p.groups))
		for _, g := range p.groups {
			entries = append(entries, entry{isGroup: true, label: g, idx: -1})
			for i, c := range p.filtered {
				if c.Group == g {
					entries = append(entries, entry{label: c.Label, keybind: c.Keybind, idx: i})
				}
			}
		}

		selectedEntry := 0
		for i, e := range entries {
			if !e.isGroup && e.idx == p.selected {
				selectedEntry = i
				break
			}
		}

		start := 0
		if selectedEntry >= maxVisible {
			start = selectedEntry - maxVisible + 1
		}
		end := minv(len(entries), start+maxVisible)

		for i := start; i < end; i++ {
			e := entries[i]
			if e.isGroup {
				groupStyle := lipgloss.NewStyle().
					Background(dbg).
					Foreground(muted).
					Bold(true).
					Width(contentWidth)
				rows = append(rows, groupStyle.Render(" "+e.label))
				continue
			}
			selected := e.idx == p.selected
			rows = append(rows, renderPaletteCommandRow(t, e.label, e.keybind, contentWidth, selected))
		}
	}

	// Navigation hint
	rows = append(rows,
		baseRow(muted, ""),
		baseRow(muted, "  ↑↓ navigate  enter run  esc close"),
	)

	return tea.NewView(strings.Join(rows, "\n"))
}

func (p *CommandPalette) SetSize(width, height int) {
	p.width = maxv(1, width)
	p.height = maxv(1, height)
	p.syncStyles()
}

func (p *CommandPalette) GetSize() (int, int) { return p.width, p.height }
func (p *CommandPalette) Title() string       { return "Commands" }

func (p *CommandPalette) syncStyles() {
	t := theme.CurrentTheme()
	s := textinput.DefaultStyles(true)
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(t.TextMuted())
	s.Focused.Text = lipgloss.NewStyle().Foreground(t.DialogFG())
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(t.InputPlaceholder())
	s.Blurred = s.Focused
	s.Cursor.Color = t.InputCursor()
	p.search.SetStyles(s)
	p.search.SetWidth(maxv(10, p.width-2))
}

func (p *CommandPalette) refilter() {
	query := strings.ToLower(strings.TrimSpace(p.search.Value()))

	next := make([]Command, 0, len(p.commands))
	for _, c := range p.commands {
		if query == "" ||
			strings.Contains(strings.ToLower(c.Label), query) ||
			strings.Contains(strings.ToLower(c.Group), query) {
			next = append(next, c)
		}
	}
	p.filtered = next

	seen := make(map[string]bool)
	groups := make([]string, 0, 4)
	for _, c := range p.filtered {
		g := c.Group
		if g == "" {
			g = "Commands"
		}
		if !seen[g] {
			seen[g] = true
			groups = append(groups, g)
		}
	}
	p.groups = groups

	if p.selected >= len(p.filtered) {
		p.selected = maxv(0, len(p.filtered)-1)
	}
}

func renderPaletteCommandRow(t theme.Theme, label, keybind string, width int, selected bool) string {
	var itemBG, itemFG, keybindFG color.Color
	if selected {
		itemBG = t.SelectionBG()
		itemFG = t.SelectionFG()
		keybindFG = t.SelectionFG()
	} else {
		itemBG = t.DialogBG()
		itemFG = t.DialogFG()
		keybindFG = t.TextMuted()
	}

	if keybind == "" {
		return lipgloss.NewStyle().Background(itemBG).Foreground(itemFG).Width(width).Render(" " + label)
	}

	keybindW := lipgloss.Width(keybind)
	labelMaxW := maxv(1, width-keybindW-3) // 1 lead + 2 sep
	labelClipped := label
	if lipgloss.Width(label) > labelMaxW {
		labelClipped = paletteClip(label, labelMaxW)
	}

	actualLabelW := lipgloss.Width(labelClipped)
	gap := maxv(1, width-1-actualLabelW-2-keybindW)
	line := " " + labelClipped + strings.Repeat(" ", gap) + "  " +
		lipgloss.NewStyle().Foreground(keybindFG).Render(keybind)
	return lipgloss.NewStyle().Background(itemBG).Foreground(itemFG).Width(width).Render(line)
}

func paletteClip(s string, width int) string {
	if width <= 0 {
		return ""
	}
	return lipgloss.NewStyle().MaxWidth(width).Width(width).Render(s)
}
