// Package list provides a scrollable log-style list widget.
// Copy this file into your project: bento add list
//
// The widget returns plain text — the containing panel applies all color.
// Dependencies: charm.land/bubbletea/v2 only.
package list

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Model is a scrollable list that shows the last N items that fit in height.
// It produces plain text — no ANSI styling — so the parent panel can paint
// the background without bleed-through.
type Model struct {
	width  int
	height int
	items  []string
	max    int
}

// New creates a list with an optional cap on stored items.
// maxItems <= 0 defaults to 200.
func New(maxItems int) *Model {
	if maxItems <= 0 {
		maxItems = 200
	}
	return &Model{max: maxItems}
}

// Append adds an item to the bottom of the list.
func (l *Model) Append(item string) {
	l.items = append(l.items, item)
	if len(l.items) > l.max {
		l.items = l.items[1:]
	}
}

// Prepend adds an item to the top of the list.
func (l *Model) Prepend(item string) {
	l.items = append([]string{item}, l.items...)
	if len(l.items) > l.max {
		l.items = l.items[:l.max]
	}
}

// Clear removes all items.
func (l *Model) Clear() { l.items = nil }

// Items returns a copy of the current item list.
func (l *Model) Items() []string { return append([]string(nil), l.items...) }

func (l *Model) Init() tea.Cmd                           { return nil }
func (l *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return l, nil }

func (l *Model) View() tea.View {
	if len(l.items) == 0 {
		return tea.NewView("")
	}
	if l.width <= 0 || l.height <= 0 {
		return tea.NewView(strings.Join(l.items, "\n"))
	}

	start := len(l.items) - l.height
	if start < 0 {
		start = 0
	}
	lines := make([]string, 0, l.height)
	for i := start; i < len(l.items); i++ {
		line := l.items[i]
		// Clip to width (rune-safe enough for log lines).
		if len([]rune(line)) > l.width {
			line = string([]rune(line)[:l.width])
		}
		lines = append(lines, line)
	}
	return tea.NewView(strings.Join(lines, "\n"))
}

func (l *Model) SetSize(width, height int) {
	l.width = width
	l.height = height
}

func (l *Model) GetSize() (int, int) { return l.width, l.height }
