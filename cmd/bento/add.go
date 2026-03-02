package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// component describes a registry entry available via `bento add`.
type component struct {
	name  string
	desc  string
	files []string // file names inside registry/<name>/
}

var registry = []component{
	{name: "panel", desc: "Titled, focusable content container", files: []string{"panel.go"}},
	{name: "bar", desc: "Header/footer row with keybind cards", files: []string{"bar.go"}},
	{name: "dialog", desc: "Modal manager, Confirm, Custom, ThemePicker, CommandPalette", files: []string{"dialog.go", "theme_picker.go", "command_palette.go"}},
	{name: "list", desc: "Scrollable log-style list (plain text output)", files: []string{"list.go"}},
	{name: "table", desc: "Header row + data rows", files: []string{"table.go"}},
	{name: "text", desc: "Static text label", files: []string{"text.go"}},
	{name: "input", desc: "Single-line text field wrapping bubbles/textinput", files: []string{"input.go"}},
}

func runAdd(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printAddHelp()
		return
	}

	name := strings.ToLower(strings.TrimSpace(args[0]))

	var comp *component
	for i := range registry {
		if registry[i].name == name {
			comp = &registry[i]
			break
		}
	}
	if comp == nil {
		fmt.Fprintf(os.Stderr, "unknown component: %q\n\nAvailable components:\n", name)
		for _, c := range registry {
			fmt.Fprintf(os.Stderr, "  %-10s %s\n", c.name, c.desc)
		}
		os.Exit(1)
	}

	// Destination: components/<name>/ relative to cwd (where go.mod lives).
	dest := filepath.Join("components", comp.name)
	fmt.Printf("bento add %s\n\n", comp.name)

	// Locate the bentotui module source in the module cache.
	// GOPATH/pkg/mod/github.com/cloudboy-jh/bentotui@<version>/registry/<name>/
	// We print the commands rather than executing them so the user can inspect
	// and redirect to a different destination if needed.
	//
	// NOTE: Once `//go:embed registry` is wired into this binary (see
	// docs/next-steps.md) this will write files directly. For now it prints
	// the equivalent shell commands.
	fmt.Printf("  Writing to %s/\n\n", dest)
	fmt.Println("  Run these commands to copy the component source:")
	fmt.Println()
	fmt.Printf("    mkdir -p %s\n", dest)
	for _, f := range comp.files {
		src := fmt.Sprintf("$(go env GOPATH)/pkg/mod/github.com/cloudboy-jh/bentotui*/registry/%s/%s", comp.name, f)
		fmt.Printf("    cp %s %s/\n", src, dest)
	}
	fmt.Println()
	fmt.Printf("  Then import from your module: \"yourmodule/%s\"\n", dest)
	fmt.Println()
	fmt.Println("  Required module deps (already in your go.mod if you ran bento init):")
	fmt.Println("    charm.land/bubbletea/v2")
	fmt.Println("    charm.land/lipgloss/v2")
	fmt.Println("    github.com/cloudboy-jh/bentotui  (theme, styles, layout)")
	fmt.Println()
}

func printAddHelp() {
	fmt.Print(`Usage: bento add <component>

Copy-and-own a registry component into your project.
Files are written to components/<name>/ and are yours to modify.

Available components:
`)
	for _, c := range registry {
		fmt.Printf("  %-10s %s\n", c.name, c.desc)
	}
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  bento add panel")
	fmt.Println("  bento add dialog")
	fmt.Println()
}
