package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cloudboy-jh/bentotui/registry/bentos/app-shell/scenarios"
)

type DiagnosticsInput struct {
	TerminalW     int
	TerminalH     int
	BodyH         int
	Viewport      scenarios.Viewport
	ThemeName     string
	ScenarioID    string
	FocusOwner    string
	PaintDebug    bool
	Snapshot      bool
	ShowKeymap    bool
	Status        string
	ContrastScore string
	Checks        []scenarios.Check
	Metrics       map[string]string
}

func DiagnosticsText(in DiagnosticsInput) string {
	lines := []string{
		"Session",
		fmt.Sprintf("- terminal: %dx%d", in.TerminalW, in.TerminalH),
		fmt.Sprintf("- body: %dx%d", in.TerminalW, in.BodyH),
		fmt.Sprintf("- viewport: %s (%dx%d)", in.Viewport.Name, in.Viewport.Width, in.Viewport.Height),
		fmt.Sprintf("- theme: %s", in.ThemeName),
		fmt.Sprintf("- contrast score: %s", in.ContrastScore),
		fmt.Sprintf("- scenario: %s", in.ScenarioID),
		fmt.Sprintf("- focus owner: %s", in.FocusOwner),
		fmt.Sprintf("- paint debug: %t", in.PaintDebug),
		fmt.Sprintf("- snapshot mode: %t", in.Snapshot),
		"",
		"Checks",
	}
	for _, c := range in.Checks {
		lines = append(lines, fmt.Sprintf("- [%s] %s: %s", c.Level, c.Name, c.Detail))
	}

	if len(in.Metrics) > 0 {
		lines = append(lines, "", "Metrics")
		keys := make([]string, 0, len(in.Metrics))
		for k := range in.Metrics {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			lines = append(lines, fmt.Sprintf("- %s: %s", k, in.Metrics[k]))
		}
	}

	if in.ShowKeymap {
		lines = append(lines, "", "Keymap", "- j/k scenario", "- h/l viewport", "- t theme", "- d paint-debug", "- [/] focus pane", "- s snapshot", "- m hide/show keymap", "- r stress step", "- q quit")
	}

	lines = append(lines, "", "Status", "- "+in.Status)
	return strings.Join(lines, "\n")
}
