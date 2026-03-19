package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cloudboy-jh/bentotui/cmd/bento/logic"
)

func (a *App) handleComponentListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if a.state.ComponentCursor > 0 {
			a.state.ComponentCursor--
		}
	case "down", "j":
		items := a.activeCatalog()
		if a.state.ComponentCursor < len(items)-1 {
			a.state.ComponentCursor++
		}
	case " ":
		// Toggle selection
		items := a.activeCatalog()
		if len(items) == 0 {
			return a, nil
		}
		comp := items[a.state.ComponentCursor]
		if a.state.SelectedComponents[comp.Name] {
			delete(a.state.SelectedComponents, comp.Name)
		} else {
			a.state.SelectedComponents[comp.Name] = true
		}
	case "enter":
		// Install selected bricks
		return a.installSelectedComponents()
	}
	return a, nil
}

func (a *App) renderComponentList(height int) string {
	_ = height
	items := a.activeCatalog()
	lines := make([]string, 0, len(items)+5)

	label := a.activeCatalogLabel()
	lines = append(lines, "  Select "+label+" to install (space to toggle, enter to install):")
	lines = append(lines, "")

	// Component list
	visibleCount := 8
	startIdx := 0
	if a.state.ComponentCursor >= visibleCount {
		startIdx = a.state.ComponentCursor - visibleCount + 1
	}

	for i := startIdx; i < len(items) && i < startIdx+visibleCount; i++ {
		comp := items[i]
		selected := a.state.SelectedComponents[comp.Name]

		checkbox := "[ ]"
		if selected {
			checkbox = "[✓]"
		}

		cursor := "  "
		if i == a.state.ComponentCursor {
			cursor = "> "
		}

		nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoFrameFG))
		descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(bentoMuted))

		if i == a.state.ComponentCursor {
			nameStyle = nameStyle.Bold(true).Foreground(lipgloss.Color(bentoAccent))
		}

		line := cursor + checkbox + " " + nameStyle.Render(comp.Name)
		if len(comp.Desc) > 0 {
			line += " - " + descStyle.Render(comp.Desc)
		}
		lines = append(lines, "  "+lipgloss.NewStyle().MaxWidth(80).Render(line))
	}

	// Footer
	lines = append(lines, "")
	selectedCount := len(a.state.SelectedComponents)
	footerText := fmt.Sprintf("  %d selected %s • enter to install • esc to go back", selectedCount, label)
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color(bentoMuted)).Render(footerText))

	return strings.Join(lines, "\n")
}

func (a *App) installSelectedComponents() (tea.Model, tea.Cmd) {
	kindLabel := a.activeCatalogLabel()
	if len(a.state.SelectedComponents) == 0 {
		a.state.AddLog("No " + kindLabel + " selected")
		return a, nil
	}

	a.state.AddLog("Installing " + kindLabel + "...")

	for name := range a.state.SelectedComponents {
		a.state.AddLog("Installing: " + name)
		result := logic.InstallComponent(name)
		if a.state.CatalogKind == CatalogRecipes {
			result = logic.InstallRecipe(name)
		}

		if result.Error != nil {
			a.state.AddLog("  Error: " + result.Error.Error())
		} else {
			for _, f := range result.Files {
				a.state.AddLog("  Created: " + f)
			}
			for _, f := range result.Skipped {
				a.state.AddLog("  Skipped (exists): " + f)
			}
		}
	}

	a.state.AddLog("Done.")
	a.state.SelectedComponents = make(map[string]bool) // Clear selection

	return a, nil
}

func (a *App) activeCatalog() []logic.CatalogEntry {
	if a.state.CatalogKind == CatalogRecipes {
		return logic.RecipeRegistry()
	}
	return logic.BrickRegistry()
}

func (a *App) activeCatalogLabel() string {
	if a.state.CatalogKind == CatalogRecipes {
		return "recipes"
	}
	return "bricks"
}
