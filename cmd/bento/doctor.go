package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type checkResult struct {
	label string
	pass  bool
	note  string // optional — shown on failure
}

func runDoctor(_ []string) {
	fmt.Println("🍱 bento doctor")
	fmt.Println()

	results := []checkResult{
		checkGoMod(),
		checkBentoDep(),
		checkStylesPresent(),
	}

	// Per-component presence checks
	for _, component := range []string{"header", "footer", "panel", "dialog"} {
		results = append(results, checkComponent(component))
	}

	allPass := true
	for _, r := range results {
		icon := "✓"
		if !r.pass {
			icon = "✗"
			allPass = false
		}
		if r.note != "" && !r.pass {
			fmt.Printf("  [%s] %s — %s\n", icon, r.label, r.note)
		} else {
			fmt.Printf("  [%s] %s\n", icon, r.label)
		}
	}

	fmt.Println()
	if allPass {
		fmt.Println("  All checks passed.")
	} else {
		fmt.Println("  Some checks failed. See notes above.")
		os.Exit(1)
	}
}

func checkGoMod() checkResult {
	_, err := os.Stat("go.mod")
	return checkResult{
		label: "go.mod found",
		pass:  err == nil,
		note:  "run 'go mod init <module>' to initialise a Go module",
	}
}

func checkBentoDep() checkResult {
	f, err := os.Open("go.mod")
	if err != nil {
		return checkResult{label: "bentotui dependency declared", pass: false, note: "go.mod not readable"}
	}
	defer f.Close()

	found := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "github.com/cloudboy-jh/bentotui") {
			found = true
			break
		}
	}
	return checkResult{
		label: "bentotui dependency declared",
		pass:  found,
		note:  "add 'require github.com/cloudboy-jh/bentotui latest' to go.mod",
	}
}

func checkStylesPresent() checkResult {
	path := filepath.Join("ui", "styles", "styles.go")
	_, err := os.Stat(path)
	return checkResult{
		label: "ui/styles/styles.go present",
		pass:  err == nil,
		note:  "styles file missing — ensure you are in the project root",
	}
}

func checkComponent(name string) checkResult {
	dir := filepath.Join("ui", "containers", name)
	entries, err := os.ReadDir(dir)
	present := err == nil && len(entries) > 0
	return checkResult{
		label: fmt.Sprintf("ui/containers/%s present", name),
		pass:  present,
		note:  fmt.Sprintf("run: bento add %s", name),
	}
}
