package logic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	bentoregistry "github.com/cloudboy-jh/bentotui/registry"
)

// ComponentInfo describes a registry component.
type ComponentInfo struct {
	Name  string
	Desc  string
	Files []string
}

// Registry returns the list of available components.
func Registry() []ComponentInfo {
	return []ComponentInfo{
		{Name: "surface", Desc: "Full-terminal paint surface with UV cell buffer", Files: []string{"surface.go"}},
		{Name: "panel", Desc: "Titled, focusable content container", Files: []string{"panel.go"}},
		{Name: "elevated-card", Desc: "Raised section container with title + content", Files: []string{"elevated_card.go"}},
		{Name: "bar", Desc: "Header/footer row with keybind cards", Files: []string{"bar.go"}},
		{Name: "dialog", Desc: "Modal manager, Confirm, Custom, ThemePicker, CommandPalette", Files: []string{"dialog.go", "theme_picker.go", "command_palette.go"}},
		{Name: "list", Desc: "Scrollable log-style list (plain text output)", Files: []string{"list.go"}},
		{Name: "table", Desc: "Header row + data rows", Files: []string{"table.go"}},
		{Name: "text", Desc: "Static text label", Files: []string{"text.go"}},
		{Name: "input", Desc: "Single-line text field wrapping bubbles/textinput", Files: []string{"input.go"}},
		{Name: "badge", Desc: "Inline themed label", Files: []string{"badge.go"}},
		{Name: "kbd", Desc: "Keyboard shortcut command+label pair", Files: []string{"kbd.go"}},
		{Name: "wordmark", Desc: "Themed heading/title block", Files: []string{"wordmark.go"}},
		{Name: "select", Desc: "Single-choice inline picker", Files: []string{"select.go"}},
		{Name: "checkbox", Desc: "Boolean toggle input", Files: []string{"checkbox.go"}},
		{Name: "progress", Desc: "Horizontal progress bar", Files: []string{"progress.go"}},
		{Name: "tabs", Desc: "Keyboard-navigable tab row", Files: []string{"tabs.go"}},
		{Name: "toast", Desc: "Stacked notification rows", Files: []string{"toast.go"}},
		{Name: "separator", Desc: "Horizontal or vertical divider", Files: []string{"separator.go"}},
	}
}

// InstallResult holds the result of installing a component.
type InstallResult struct {
	Component string
	Files     []string
	Skipped   []string
	Error     error
}

// InstallComponent copies a component from the registry to the local project.
// Returns the result including any files written or skipped.
func InstallComponent(name string) InstallResult {
	result := InstallResult{Component: name}

	// Find component in registry
	var comp *ComponentInfo
	for _, c := range Registry() {
		if c.Name == name {
			comp = &c
			break
		}
	}
	if comp == nil {
		result.Error = fmt.Errorf("unknown component: %s", name)
		return result
	}

	// Create destination directory
	destDir := filepath.Join("bricks", comp.Name)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		result.Error = fmt.Errorf("create directory %s: %w", destDir, err)
		return result
	}

	// Copy each file
	for _, f := range comp.Files {
		srcPath := "bricks/" + comp.Name + "/" + f
		dstPath := filepath.Join(destDir, f)

		// Check if file already exists
		if _, err := os.Stat(dstPath); err == nil {
			result.Skipped = append(result.Skipped, dstPath)
			continue
		}

		// Read from embedded FS
		srcFile, err := bentoregistry.BricksFS.Open(srcPath)
		if err != nil {
			result.Error = fmt.Errorf("component %q file %q not found: %w", comp.Name, f, err)
			return result
		}

		// Create destination file
		dstFile, err := os.Create(dstPath)
		if err != nil {
			srcFile.Close()
			result.Error = fmt.Errorf("create file %s: %w", dstPath, err)
			return result
		}

		// Copy content
		_, copyErr := io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()

		if copyErr != nil {
			result.Error = fmt.Errorf("write file %s: %w", dstPath, copyErr)
			return result
		}

		result.Files = append(result.Files, dstPath)
	}

	return result
}
