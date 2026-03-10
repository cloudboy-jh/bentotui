package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui/cmd/bento/tui"
)

var version = "dev"

func main() {
	// No args = TUI mode
	if len(os.Args) < 2 {
		runTUI()
		return
	}

	// CLI mode
	switch os.Args[1] {
	case "init":
		runInitCLI(os.Args[2:])
	case "add":
		runAddCLI(os.Args[2:])
	case "list":
		runListCLI(os.Args[2:])
	case "doctor":
		runDoctorCLI(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Printf("🍱 bento v%s\n", version)
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func runTUI() {
	app := tui.NewApp()
	if _, err := tea.NewProgram(app).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(`🍱 bento — BentoTUI project CLI

Usage:
  bento                Launch interactive TUI
  bento <command>      Run CLI command

Commands:
  init [name]          Scaffold a new BentoTUI app
  add <component>      Copy-and-own a component into your project
  list                 Show available registry components
  doctor               Check your project for common issues
  version              Print the bento version
  help                 Show this help message

Run 'bento' with no arguments to use the interactive TUI.
`)
}

// fatal prints an error message and exits with code 1.
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}

// check exits with a fatal error if err is non-nil.
func check(err error, context string) {
	if err != nil {
		fatal("%s: %v", context, err)
	}
}
