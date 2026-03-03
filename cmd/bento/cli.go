package main

import (
	"fmt"
	"os"

	"github.com/cloudboy-jh/bentotui/cmd/bento/logic"
)

// runInitCLI runs the init command in CLI mode.
func runInitCLI(args []string) {
	appName := ""
	if len(args) > 0 && args[0] != "--help" && args[0] != "-h" {
		appName = args[0]
	}

	cfg := logic.ProjectConfig{
		AppName: appName,
	}

	fmt.Printf("Creating project: %s\n", cfg.AppName)
	created, err := logic.ScaffoldProject(cfg)
	if err != nil {
		fatal("%v", err)
	}

	for _, f := range created {
		fmt.Printf("  Created: %s\n", f)
	}
	fmt.Println("Done.")
}

// runAddCLI runs the add command in CLI mode.
func runAddCLI(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printAddHelp()
		return
	}

	name := args[0]
	fmt.Printf("Installing component: %s\n", name)

	result := logic.InstallComponent(name)
	if result.Error != nil {
		fatal("%v", result.Error)
	}

	for _, f := range result.Files {
		fmt.Printf("  Created: %s\n", f)
	}
	for _, f := range result.Skipped {
		fmt.Printf("  Skipped (exists): %s\n", f)
	}
	fmt.Println("Done.")
}

// runDoctorCLI runs the doctor command in CLI mode.
func runDoctorCLI(args []string) {
	fmt.Println("Running doctor checks...")
	report := logic.RunDoctor()

	for _, r := range report.Results {
		icon := "✓"
		if !r.Pass {
			icon = "✗"
		}
		if r.Note != "" && !r.Pass {
			fmt.Printf("  [%s] %s - %s\n", icon, r.Label, r.Note)
		} else {
			fmt.Printf("  [%s] %s\n", icon, r.Label)
		}
	}

	if report.AllPass {
		fmt.Println("All checks passed!")
	} else {
		fmt.Println("Some checks failed.")
		os.Exit(1)
	}
}

func printAddHelp() {
	fmt.Print(`Usage: bento add <component>

Copies source into components/<name>/ and are yours to modify.

Available components:
`)
	for _, c := range logic.Registry() {
		fmt.Printf("  %-10s %s\n", c.Name, c.Desc)
	}
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  bento add panel")
	fmt.Println("  bento add dialog")
	fmt.Println()
}
