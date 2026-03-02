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
	note  string
}

func runDoctor(_ []string) {
	fmt.Println("bento doctor")
	fmt.Println()

	results := []checkResult{
		checkGoMod(),
		checkBentoDep(),
	}

	// Check that the three stable module deps are importable (present in go.mod).
	for _, pkg := range []string{
		"github.com/cloudboy-jh/bentotui/theme",
		"github.com/cloudboy-jh/bentotui/styles",
		"github.com/cloudboy-jh/bentotui/layout",
	} {
		results = append(results, checkModDep(pkg))
	}

	// Check for any copied registry components.
	for _, name := range []string{"panel", "bar", "dialog", "list", "table", "text", "input"} {
		results = append(results, checkCopiedComponent(name))
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
		fmt.Println("  Registry components that are missing have not been added yet — that is fine.")
		os.Exit(1)
	}
}

func checkGoMod() checkResult {
	_, err := os.Stat("go.mod")
	return checkResult{
		label: "go.mod present",
		pass:  err == nil,
		note:  "run 'go mod init <module>' to initialise a Go module",
	}
}

func checkBentoDep() checkResult {
	f, err := os.Open("go.mod")
	if err != nil {
		return checkResult{label: "bentotui declared in go.mod", pass: false, note: "go.mod not readable"}
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
		label: "bentotui declared in go.mod",
		pass:  found,
		note:  "run: go get github.com/cloudboy-jh/bentotui",
	}
}

// checkModDep verifies that a package path is resolvable by looking for it in
// the module cache or vendor directory. This is a lightweight proxy check —
// it does not invoke the compiler.
func checkModDep(pkg string) checkResult {
	label := fmt.Sprintf("  %s importable", strings.TrimPrefix(pkg, "github.com/cloudboy-jh/bentotui/"))

	// Check vendor/ first (common in CI).
	vendorPath := filepath.Join("vendor", filepath.FromSlash(pkg))
	if info, err := os.Stat(vendorPath); err == nil && info.IsDir() {
		return checkResult{label: label, pass: true}
	}

	// Fall back to GOPATH module cache.
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}
	cachePath := filepath.Join(gopath, "pkg", "mod", filepath.FromSlash("github.com/cloudboy-jh"))
	entries, err := os.ReadDir(cachePath)
	if err == nil {
		suffix := strings.TrimPrefix(pkg, "github.com/cloudboy-jh/bentotui")
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "bentotui@") {
				candidate := filepath.Join(cachePath, e.Name(), filepath.FromSlash(suffix))
				if _, err := os.Stat(candidate); err == nil {
					return checkResult{label: label, pass: true}
				}
			}
		}
	}

	return checkResult{
		label: label,
		pass:  false,
		note:  "run: go mod tidy",
	}
}

// checkCopiedComponent checks whether a registry component has been copied into
// the project's components/ directory. A missing component is not an error —
// users only copy what they need.
func checkCopiedComponent(name string) checkResult {
	dir := filepath.Join("components", name)
	entries, err := os.ReadDir(dir)
	present := err == nil && len(entries) > 0
	label := fmt.Sprintf("  components/%s copied", name)
	note := fmt.Sprintf("run: bento add %s  (optional)", name)
	return checkResult{label: label, pass: present, note: note}
}
