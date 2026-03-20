package logic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	bentoregistry "github.com/cloudboy-jh/bentotui/registry"
)

// CatalogEntry describes one installable registry item.
type CatalogEntry struct {
	Name  string
	Desc  string
	Files []string
}

// BrickRegistry returns the list of available bricks.
func BrickRegistry() []CatalogEntry {
	return []CatalogEntry{
		{Name: "surface", Desc: "Full-terminal paint surface with UV cell buffer", Files: []string{"surface.go"}},
		{Name: "card", Desc: "Content container — raised (default) or flat via Flat() option", Files: []string{"card.go"}},
		{Name: "bar", Desc: "Header/footer row with keybind cards", Files: []string{"bar.go"}},
		{Name: "dialog", Desc: "Modal manager, Confirm, Custom, ThemePicker, CommandPalette", Files: []string{"dialog.go", "theme_picker.go", "command_palette.go"}},
		{Name: "filepicker", Desc: "File and directory picker wrapping bubbles/filepicker", Files: []string{"filepicker.go"}},
		{Name: "list", Desc: "Scrollable list wrapping bubbles/list", Files: []string{"list.go"}},
		{Name: "table", Desc: "Data table wrapping bubbles/table", Files: []string{"table.go"}},
		{Name: "text", Desc: "Static text label", Files: []string{"text.go"}},
		{Name: "input", Desc: "Single-line text field wrapping bubbles/textinput", Files: []string{"input.go"}},
		{Name: "badge", Desc: "Inline themed label", Files: []string{"badge.go"}},
		{Name: "kbd", Desc: "Keyboard shortcut command+label pair", Files: []string{"kbd.go"}},
		{Name: "wordmark", Desc: "Themed heading/title block", Files: []string{"wordmark.go"}},
		{Name: "select", Desc: "Single-choice picker wrapping bubbles/list", Files: []string{"select.go"}},
		{Name: "checkbox", Desc: "Boolean toggle using bubbles key bindings", Files: []string{"checkbox.go"}},
		{Name: "progress", Desc: "Horizontal progress bar wrapping bubbles/progress", Files: []string{"progress.go"}},
		{Name: "package-manager", Desc: "Sequential install flow with spinner + progress", Files: []string{"package_manager.go"}},
		{Name: "tabs", Desc: "Tab row with bubbles key/paginator input", Files: []string{"tabs.go"}},
		{Name: "toast", Desc: "Stacked notification rows", Files: []string{"toast.go"}},
		{Name: "separator", Desc: "Horizontal or vertical divider", Files: []string{"separator.go"}},
	}
}

// RecipeRegistry returns the list of available recipes.
func RecipeRegistry() []CatalogEntry {
	return []CatalogEntry{
		{Name: "filter-bar", Desc: "Input + status + keybind strip composition", Files: []string{"recipe.go"}},
		{Name: "empty-state-pane", Desc: "Reusable empty-result card content", Files: []string{"recipe.go"}},
		{Name: "command-palette-flow", Desc: "Open command palette and route actions", Files: []string{"recipe.go"}},
		{Name: "vimstatus", Desc: "Vim-style statusline with mode pill, context, and clock", Files: []string{"VimFooter.go"}},
	}
}

// Registry is kept for compatibility and returns bricks.
func Registry() []CatalogEntry {
	return BrickRegistry()
}

// InstallResult holds the result of an install operation.
type InstallResult struct {
	Name    string
	Files   []string
	Skipped []string
	Error   error
}

// InstallComponent copies a brick from the registry to the local project.
// Returns the result including any files written or skipped.
func InstallComponent(name string) InstallResult {
	return installFromCatalog("brick", name, "bricks", BrickRegistry())
}

// InstallRecipe copies a recipe from the registry to the local project.
// Returns the result including any files written or skipped.
func InstallRecipe(name string) InstallResult {
	return installFromCatalog("recipe", name, "recipes", RecipeRegistry())
}

func installFromCatalog(kind, name, dir string, catalog []CatalogEntry) InstallResult {
	result := InstallResult{Name: name}

	var entry *CatalogEntry
	for _, c := range catalog {
		if c.Name == name {
			entry = &c
			break
		}
	}
	if entry == nil {
		result.Error = fmt.Errorf("unknown %s: %s", kind, name)
		return result
	}

	destDir := filepath.Join(dir, entry.Name)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		result.Error = fmt.Errorf("create directory %s: %w", destDir, err)
		return result
	}

	for _, f := range entry.Files {
		srcPath := dir + "/" + entry.Name + "/" + f
		dstPath := filepath.Join(destDir, f)

		if _, err := os.Stat(dstPath); err == nil {
			result.Skipped = append(result.Skipped, dstPath)
			continue
		}

		srcFile, err := bentoregistry.BricksFS.Open(srcPath)
		if err != nil {
			result.Error = fmt.Errorf("%s %q file %q not found: %w", kind, entry.Name, f, err)
			return result
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			srcFile.Close()
			result.Error = fmt.Errorf("create file %s: %w", dstPath, err)
			return result
		}

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
