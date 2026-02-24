![BentoTUI logo](./bentotui-readme-logo.png)

# BentoTUI

BentoTUI is an application framework on top of Bubble Tea for building production-grade TUIs with a structured shell, layered rendering, and reusable UI components.

## What BentoTUI provides

- app shell with deterministic layer order (body -> footer -> overlays)
- page router with lazy page factories
- fixed/flex split layout system
- focus ring manager for keyboard navigation
- modal dialog manager and theme picker flow
- semantic theme system (currently `catppuccin-mocha`, `dracula`, `osaka-jade`)
- UI component layer under `ui/components/*`

## Current package layout

- runtime/framework: `app`, `shell`, `router`, `layout`, `focus`, `surface`, `core`, `theme`
- UI components: `ui/components/dialog`, `ui/components/footer`, `ui/components/panel`
- style layer: `ui/styles`

## Install

```bash
go get github.com/cloudboy-jh/bentotui
```

## Quick start

```go
package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/layout"
	"github.com/cloudboy-jh/bentotui/theme"
	"github.com/cloudboy-jh/bentotui/ui/components/panel"
)

func main() {
	m := bentotui.New(
		bentotui.WithTheme(theme.Preset("catppuccin-mocha")),
		bentotui.WithPages(
			bentotui.Page("home", func() core.Page { return newHomePage() }),
		),
		bentotui.WithFooterBar(true),
	)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("run failed: %v\n", err)
	}
}

type homePage struct {
	root   *layout.Split
	width  int
	height int
}

func newHomePage() *homePage {
	sidebar := panel.New(panel.Title("Sidebar"), panel.Content(staticText("Sessions\nFiles")))
	main := panel.New(panel.Title("Main"), panel.Content(staticText("Welcome to BentoTUI")))
	root := layout.Horizontal(
		layout.Fixed(30, sidebar),
		layout.Flex(1, main),
	)
	return &homePage{root: root}
}

func (p *homePage) Init() tea.Cmd { return nil }
func (p *homePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := p.root.Update(msg)
	return p, cmd
}
func (p *homePage) View() tea.View { return p.root.View() }
func (p *homePage) SetSize(w, h int) {
	p.width = w
	p.height = h
	p.root.SetSize(w, h)
}
func (p *homePage) GetSize() (int, int) { return p.width, p.height }
func (p *homePage) Title() string       { return "Home" }

type staticText string

func (s staticText) Init() tea.Cmd                           { return nil }
func (s staticText) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s staticText) View() tea.View                          { return tea.NewView(string(s)) }
```

Note: full-screen mode is enabled by default. Disable it with `bentotui.WithFullScreen(false)`.

## Internal harness

Run the framework harness:

```bash
go run ./cmd/test-tui
```

It validates shell layering, focus behavior, modal overlays, theme switching, and component rendering.

## Documentation

- spec: `project-docs/bentotui-main-spec.md`
- rendering ADR: `project-docs/rendering-system-design.md`
- implementation roadmap: `project-docs/next-steps.md`
- framework research: `project-docs/tui-framework-research.md`

## Development

```bash
go test ./...
go vet ./...
```
