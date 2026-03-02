package dialog

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/primitives"
	"github.com/cloudboy-jh/bentotui/ui/styles"
)

// Command is a single entry in the command palette.
type Command struct {
	Label   string // e.g. "Switch theme"
	Group   string // e.g. "System" — empty = ungrouped, shown first
	Keybind string // e.g. "ctrl+x t" — optional, shown right-aligned
	Action  func() tea.Msg
}

// CommandPalette is a searchable, grouped command picker dialog.
// It implements the Dialog interface and renders its own framed output.
type CommandPalette struct {
	commands []Command // full registered set, stable order
	filtered []Command // post-search subset
	groups   []string  // unique group names in stable insertion order
	selected int
	search   textinput.Model
	width    int
	height   int
}

// NewCommandPalette creates a CommandPalette pre-loaded with the given commands.
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
	contentWidth := maxInt(24, p.width)

	rows := make([]string, 0, 16)

	// Search label + input
	rows = append(rows, primitives.RenderRow(contentWidth, "", t.Text.Primary, "Search"))
	rows = append(rows, inputContainer(p.search.View(), contentWidth, t))
	rows = append(rows, primitives.RenderRow(contentWidth, "", "", ""))

	if len(p.filtered) == 0 {
		rows = append(rows, primitives.RenderRow(contentWidth, "", t.Text.Muted, "No matching commands"))
	} else {
		// Viewport: how many item rows we can show
		// header(3) + blank(1) + hint(2) = 6 overhead rows
		maxVisible := maxInt(1, p.height-6)

		// Build a flat indexed list for viewport scrolling
		type entry struct {
			isGroup bool
			label   string // group name or command label
			keybind string
			idx     int // index into p.filtered (-1 for group headers)
		}
		entries := make([]entry, 0, len(p.filtered)+len(p.groups))
		for _, g := range p.groups {
			entries = append(entries, entry{isGroup: true, label: g, idx: -1})
			for i, c := range p.filtered {
				if c.Group == g {
					entries = append(entries, entry{isGroup: false, label: c.Label, keybind: c.Keybind, idx: i})
				}
			}
		}

		// Find the entry index of the selected command so we can scroll to it
		selectedEntry := 0
		for i, e := range entries {
			if !e.isGroup && e.idx == p.selected {
				selectedEntry = i
				break
			}
		}

		// Compute start offset to keep selected entry visible
		start := 0
		if selectedEntry >= maxVisible {
			start = selectedEntry - maxVisible + 1
		}
		end := minInt(len(entries), start+maxVisible)

		for i := start; i < end; i++ {
			e := entries[i]
			if e.isGroup {
				groupLabel := " " + e.label
				rows = append(rows, primitives.RenderStyledRow(sys.PaletteGroupHeader(), contentWidth, groupLabel))
				continue
			}
			selected := e.idx == p.selected
			row := renderCommandRow(e.label, e.keybind, contentWidth, sys, selected)
			rows = append(rows, row)
		}
	}

	// Bottom hint
	rows = append(rows,
		primitives.RenderRow(contentWidth, "", "", ""),
		primitives.RenderRow(contentWidth, "", t.Text.Muted, "↑↓ navigate  enter run  esc close"),
	)

	body := strings.Join(rows, "\n")

	// Render via the shared dialog frame (provides title + esc hint border)
	tt := theme.CurrentTheme()
	view := renderDialogFrame("Commands", body, maxInt(50, p.width+4), maxInt(12, p.height+4), tt)
	return tea.NewView(view)
}

func (p *CommandPalette) SetSize(width, height int) {
	p.width = maxInt(1, width)
	p.height = maxInt(1, height)
	p.syncStyles()
}

func (p *CommandPalette) GetSize() (int, int) { return p.width, p.height }

func (p *CommandPalette) Title() string { return "Commands" }

// syncStyles updates the textinput styles to match the current theme.
func (p *CommandPalette) syncStyles() {
	t := theme.CurrentTheme()
	s := textinput.DefaultStyles(true)
	s.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Muted))
	s.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Text.Primary))
	s.Focused.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(t.Input.Placeholder))
	s.Blurred = s.Focused
	s.Cursor.Color = lipgloss.Color(t.Input.Cursor)
	p.search.SetStyles(s)
	p.search.SetWidth(maxInt(10, p.width-2))
}

// refilter rebuilds p.filtered and p.groups from the current search query.
func (p *CommandPalette) refilter() {
	query := strings.ToLower(strings.TrimSpace(p.search.Value()))

	// Collect matching commands, preserving original order
	next := make([]Command, 0, len(p.commands))
	for _, c := range p.commands {
		if query == "" || strings.Contains(strings.ToLower(c.Label), query) || strings.Contains(strings.ToLower(c.Group), query) {
			next = append(next, c)
		}
	}
	p.filtered = next

	// Rebuild stable group order (insertion order from filtered set)
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

	// Clamp selection
	if p.selected >= len(p.filtered) {
		p.selected = maxInt(0, len(p.filtered)-1)
	}
}

// renderCommandRow renders a single command row with label left and keybind right.
func renderCommandRow(label, keybind string, width int, sys styles.System, selected bool) string {
	itemStyle := sys.PaletteItem(selected)
	keybindStyle := sys.PaletteKeybind()

	if selected {
		// On selected rows, keybind uses the selection foreground but slightly muted
		keybindStyle = keybindStyle.Foreground(lipgloss.Color(sys.Theme.Selection.FG))
	}

	if keybind == "" {
		// No keybind — just render the label full-width
		content := " " + primitives.FitWidth(label, maxInt(1, width-1))
		return primitives.RenderStyledRow(itemStyle, width, content)
	}

	// Label left, keybind right — right-pad label to fill the gap
	keybindRendered := keybindStyle.Render(keybind)
	keybindWidth := lipgloss.Width(keybindRendered)
	sep := 2
	labelWidth := maxInt(1, width-keybindWidth-sep-1) // -1 for leading space
	labelRendered := primitives.FitWidth(label, labelWidth)
	gap := maxInt(0, width-1-lipgloss.Width(labelRendered)-keybindWidth-sep)
	line := fmt.Sprintf(" %s%s  %s", labelRendered, strings.Repeat(" ", gap), keybindRendered)
	return primitives.RenderStyledRow(itemStyle, width, line)
}
