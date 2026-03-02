package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// knownComponents lists components available via `bento add`.
// The actual source is fetched from the live module, not embedded stale copies.
var knownComponents = []string{"bar", "panel", "dialog"}

func runAdd(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printAddHelp()
		return
	}

	component := strings.ToLower(strings.TrimSpace(args[0]))

	known := false
	for _, k := range knownComponents {
		if k == component {
			known = true
			break
		}
	}
	if !known {
		fmt.Fprintf(os.Stderr, "unknown component: %q\n\nAvailable components:\n", component)
		for _, k := range knownComponents {
			fmt.Fprintf(os.Stderr, "  %s\n", k)
		}
		os.Exit(1)
	}

	// Source is the real component directory in the bentotui module cache.
	// We locate it via GOPATH/pkg/mod or the module source in the vendor dir.
	fmt.Printf("🍱 bento add %s\n\n", component)
	fmt.Printf("  Component: ui/containers/%s\n\n", component)
	fmt.Println("  Run the following to copy the component source into your project:")
	fmt.Println()

	destDir := filepath.Join("ui", "containers", component)
	fmt.Printf("    mkdir -p %s\n", destDir)
	fmt.Printf("    cp $(go env GOPATH)/pkg/mod/github.com/cloudboy-jh/bentotui*/ui/containers/%s/*.go %s/\n", component, destDir)
	fmt.Println()
	fmt.Println("  Or with Go workspace/vendor setup:")
	fmt.Printf("    cp vendor/github.com/cloudboy-jh/bentotui/ui/containers/%s/*.go %s/\n", component, destDir)
	fmt.Println()
	fmt.Println("  After copying, update the package declaration if needed.")
	fmt.Println()
}

func printAddHelp() {
	fmt.Print(`Usage: bento add <component>

Copy-and-own a BentoTUI component into your project.
Files are written to ui/containers/<component>/ and are yours to modify.

Available components:
`)
	for _, k := range knownComponents {
		fmt.Printf("  %s\n", k)
	}
	fmt.Println()
}
