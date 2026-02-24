![Image 1](./bentotui-readme-logo.png)

# BentoTUI

The application framework for Bubble Tea.

BentoTUI sits between low-level Bubble Tea primitives and shipped terminal apps. It gives you a reusable app skeleton for routing, layout, focus, dialogs, status surfaces, and theming.

## Status

BentoTUI is in active v0.1 development.

Implemented today:

- app shell (`app`)
- page routing with lazy page creation (`router`)
- fixed/flex split layouts (`layout`)
- focus ring (`focus`)
- dialog manager (`dialog`)
- status bar (`statusbar`)
- panel surfaces (`panel`)
- semantic theme presets (`theme`)

## Install

```bash
go get github.com/cloudboy-jh/bentotui
```

## Quick Start

```go
package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/cloudboy-jh/bentotui"
	"github.com/cloudboy-jh/bentotui/core"
	"github.com/cloudboy-jh/bentotui/layout"
	"github.com/cloudboy-jh/bentotui/panel"
)

func main() {
	app := bentotui.New(
		bentotui.WithPages(
			bentotui.Page("home", func() core.Page { return newHomePage() }),
		),
		bentotui.WithStatusBar(true),
	)

	p := tea.NewProgram(app)
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

Fullscreen mode is enabled by default. Disable it with:

```go
bentotui.WithFullScreen(false)
```

## Internal Harness

Use the internal harness to validate rendering and interaction behavior:

```bash
go run ./cmd/test-tui
```

## Docs

- Main spec: `project-docs/bentotui-main-spec.md`
- Rendering system design (ADR-0001): `project-docs/rendering-system-design.md`
- Research notes: `project-docs/tui-framework-research.md`

## Development

```bash
go test ./...
```
