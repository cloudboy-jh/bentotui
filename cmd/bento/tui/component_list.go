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
	case " ", "space":
		if a.state.CatalogKind == CatalogBentos {
			return a, nil
		}
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
	if a.state.CatalogKind == CatalogBentos {
		lines = append(lines, "  Choose a bento template (enter to initialize):")
	} else {
		lines = append(lines, "  Select "+label+" to install (space to toggle, enter to install):")
	}
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

		marker := ""
		if a.state.CatalogKind != CatalogBentos {
			marker = "[ ] "
			if selected {
				marker = "[✓] "
			}
		}

		cursor := "  "
		if i == a.state.ComponentCursor {
			cursor = "> "
		}

		nameStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(bentoFrameFG)).
			Background(lipgloss.Color(bentoFrameBG))
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(bentoMuted)).
			Background(lipgloss.Color(bentoFrameBG))

		if i == a.state.ComponentCursor {
			nameStyle = nameStyle.Bold(true).Foreground(lipgloss.Color(bentoAccent))
		}

		line := cursor + marker + nameStyle.Render(comp.Name)
		if len(comp.Desc) > 0 {
			line += " - " + descStyle.Render(comp.Desc)
		}

		lineWidth := clamp(clamp(a.state.Width-2, minFrameWidth, maxFrameWidth)-6, 20, maxFrameWidth-6)
		lineStyle := lipgloss.NewStyle().
			MaxWidth(lineWidth).
			Background(lipgloss.Color(bentoFrameBG))
		lines = append(lines, "  "+lineStyle.Render(line))
	}

	// Footer
	lines = append(lines, "")
	footerText := ""
	if a.state.CatalogKind == CatalogBentos {
		footerText = "  enter initialize • esc go back"
	} else {
		selectedCount := len(a.state.SelectedComponents)
		footerText = fmt.Sprintf("  %d selected %s • enter to install • esc to go back", selectedCount, label)
	}
	lines = append(lines, lipgloss.NewStyle().
		Foreground(lipgloss.Color(bentoMuted)).
		Background(lipgloss.Color(bentoFrameBG)).
		Render(footerText))

	return strings.Join(lines, "\n")
}

func (a *App) installSelectedComponents() (tea.Model, tea.Cmd) {
	if a.state.CatalogKind == CatalogBentos {
		items := a.activeCatalog()
		if len(items) == 0 {
			a.state.AddLog("No bentos available")
			return a, nil
		}
		if a.state.ComponentCursor < 0 || a.state.ComponentCursor >= len(items) {
			a.state.AddLog("Invalid bento selection")
			return a, nil
		}

		name := items[a.state.ComponentCursor].Name
		a.state.AddLog("Initializing bento: " + name)
		result := logic.InstallBento(name)
		if result.Error != nil {
			a.state.AddLog("  Error: " + result.Error.Error())
			return a, nil
		}
		for _, f := range result.Files {
			a.state.AddLog("  Created: " + f)
		}
		a.state.AddLog("Done.")
		return a, nil
	}

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
	if a.state.CatalogKind == CatalogBentos {
		return logic.BentoRegistry()
	}
	if a.state.CatalogKind == CatalogRecipes {
		return logic.RecipeRegistry()
	}
	return logic.BrickRegistry()
}

func (a *App) activeCatalogLabel() string {
	if a.state.CatalogKind == CatalogBentos {
		return "bentos"
	}
	if a.state.CatalogKind == CatalogRecipes {
		return "recipes"
	}
	return "bricks"
}
