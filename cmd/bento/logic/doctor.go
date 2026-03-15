package logic

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CheckResult represents the result of a single doctor check.
type CheckResult struct {
	Label string
	Pass  bool
	Note  string
}

// DoctorReport contains all check results.
type DoctorReport struct {
	Results []CheckResult
	AllPass bool
}

// RunDoctor performs all checks and returns the report.
func RunDoctor() DoctorReport {
	results := []CheckResult{
		checkGoMod(),
		checkBentoDep(),
	}

	// Check that the three stable module deps are importable
	for _, pkg := range []string{
		"github.com/cloudboy-jh/bentotui/theme",
		"github.com/cloudboy-jh/bentotui/theme/styles",
		"github.com/cloudboy-jh/bentotui/registry/rooms",
	} {
		results = append(results, checkModDep(pkg))
	}

	// Check for any copied registry components
	for _, name := range []string{
		"surface",
		"panel",
		"bar",
		"dialog",
		"list",
		"table",
		"text",
		"input",
		"badge",
		"kbd",
		"wordmark",
		"select",
		"checkbox",
		"progress",
		"tabs",
		"toast",
		"separator",
	} {
		results = append(results, checkCopiedComponent(name))
	}

	allPass := true
	for _, r := range results {
		if !r.Pass {
			allPass = false
			break
		}
	}

	return DoctorReport{
		Results: results,
		AllPass: allPass,
	}
}

func checkGoMod() CheckResult {
	_, err := os.Stat("go.mod")
	return CheckResult{
		Label: "go.mod present",
		Pass:  err == nil,
		Note:  "run 'go mod init <module>' to initialise a Go module",
	}
}

func checkBentoDep() CheckResult {
	f, err := os.Open("go.mod")
	if err != nil {
		return CheckResult{Label: "bentotui declared in go.mod", Pass: false, Note: "go.mod not readable"}
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
	return CheckResult{
		Label: "bentotui declared in go.mod",
		Pass:  found,
		Note:  "run: go get github.com/cloudboy-jh/bentotui",
	}
}

func checkModDep(pkg string) CheckResult {
	label := fmt.Sprintf("%s importable", strings.TrimPrefix(pkg, "github.com/cloudboy-jh/bentotui/"))

	// Check vendor/ first (common in CI)
	vendorPath := filepath.Join("vendor", filepath.FromSlash(pkg))
	if info, err := os.Stat(vendorPath); err == nil && info.IsDir() {
		return CheckResult{Label: label, Pass: true}
	}

	// Fall back to GOPATH module cache
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
					return CheckResult{Label: label, Pass: true}
				}
			}
		}
	}

	return CheckResult{
		Label: label,
		Pass:  false,
		Note:  "run: go mod tidy",
	}
}

func checkCopiedComponent(name string) CheckResult {
	dir := filepath.Join("components", name)
	entries, err := os.ReadDir(dir)
	present := err == nil && len(entries) > 0
	label := fmt.Sprintf("bricks/%s copied", name)
	note := fmt.Sprintf("run: bento add %s  (optional)", name)
	return CheckResult{Label: label, Pass: present, Note: note}
}
