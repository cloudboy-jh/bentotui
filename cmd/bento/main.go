package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		runInit(os.Args[2:])
	case "add":
		runAdd(os.Args[2:])
	case "doctor":
		runDoctor(os.Args[2:])
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

func printHelp() {
	fmt.Print(`🍱 bento — BentoTUI project CLI

Usage:
  bento <command> [arguments]

Commands:
  init [name]        Scaffold a new BentoTUI app
  add <component>    Copy-and-own a component into your project
  doctor             Check your project for common issues
  version            Print the bento version

Run 'bento <command> --help' for details on a specific command.
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
