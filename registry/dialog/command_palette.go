package dialog

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/cloudboy-jh/bentotui/styles"
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
		return p, tea.Batch(
			func() tea.Msg { return action() },
			func() tea.Msg { return Close() },
		)
	}

	updated, cmd := p.search.Update(keyMsg)
	p.search = updated
	p.refilter()
	return p, cmd
}

func (p *CommandPalette) View() tea.View {
	t := theme.CurrentTheme()
	sys := styles.New(t)
	contentWidth := maxv(24, p.width)

	rows := make([]string, 0, 16)

	// Search input
	inputColors := styles.New(t).InputColors()
	inputContent := fitWidth(p.search.View(), maxv(1, contentWidth-2))
	rows = append(rows, renderRow(contentWidth, inputColors.BG, inputColors.FG, " "+inputContent))
	rows = append(rows, renderRow(contentWidth, "", "", ""))

	if len(p.filtered) == 0 {
		rows = append(rows, renderRow(contentWidth, "", t.Text.Muted, "No matching commands"))
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
				rows = append(rows, renderStyledRow(sys.PaletteGroupHeader(), contentWidth, " "+e.label))
				continue
			}
			selected := e.idx == p.selected
			rows = append(rows, renderCommandRow(e.label, e.keybind, contentWidth, sys, selected))
		}
	}

	// Navigation hint
	rows = append(rows,
		renderRow(contentWidth, "", "", ""),
		renderRow(contentWidth, "", t.Text.Muted, "↑↓ navigate  enter run  esc close"),
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
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Primary))
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Input.Placeholder))
	s.Blurred = s.Focused
	s.Cursor.Color = lipgloss.Color(t.Input.Cursor)
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

// renderCommandRow renders a single command with label left, keybind right.
func renderCommandRow(label, keybind string, width int, sys styles.System, selected bool) string {
	itemStyle := sys.PaletteItem(selected)
	keybindStyle := sys.PaletteKeybind()
	if selected {
		keybindStyle = keybindStyle.Foreground(lipgloss.Color(sys.Theme.Selection.FG))
	}

	if keybind == "" {
		content := " " + fitWidth(label, maxv(1, width-1))
		return renderStyledRow(itemStyle, width, content)
	}

	keybindRendered := keybindStyle.Render(keybind)
	keybindWidth := lipgloss.Width(keybindRendered)
	sep := 2
	labelWidth := maxv(1, width-keybindWidth-sep-1)
	labelRendered := fitWidth(label, labelWidth)
	gap := maxv(0, width-1-lipgloss.Width(labelRendered)-keybindWidth-sep)
	line := fmt.Sprintf(" %s%s  %s", labelRendered, strings.Repeat(" ", gap), keybindRendered)
	return renderStyledRow(itemStyle, width, line)
}

// renderStyledRow renders content over a full-width styled background using
// a single lipgloss Width().Render() call — no canvas layers.
func renderStyledRow(style lipgloss.Style, width int, content string) string {
	if width <= 0 {
		return ""
	}
	return style.Width(width).Render(content)
}
