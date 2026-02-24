[Image 1]
![Image 1](./bentotui-readme-logo.png)

# BentoTUI

The application framework for Bubble Tea.

BentoTUI sits between low-level Bubble Tea primitives and shipped terminal apps. It gives you reusable app patterns: page routing, panel layouts, focus orchestration, dialog handling, and status surfaces.

## Current Status

This repository is in early v0.1 scaffolding.

Implemented today:

- `app` shell for root composition and message flow
- `router` with lazy page factories and page caching
- `layout` with fixed/flex horizontal and vertical splits
- `focus` ring manager foundations
- `theme` presets and option-based overrides
- `dialog` manager with `Confirm` and `Custom` models
- `statusbar` and `panel` primitives
- `cmd/demo` starter app

## Why BentoTUI

Bubble Tea gives you great primitives. Production apps still have to rebuild the same architectural layer every time:

- multi-panel layout math
- focus routing across components
- page switching and app shell composition
- modal/dialog input capture patterns
- status/help surface wiring

BentoTUI packages those patterns in composable Go packages so your app code can stay focused on domain logic.

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

## Run the Demo

```bash
go run ./cmd/demo
```

## Package Overview

- `app`: root shell and top-level routing/status/dialog coordination
- `router`: page registration, lazy creation, active page switching
- `layout`: fixed/flex split containers for horizontal/vertical composition
- `focus`: focus ring helper and focus cycling model
- `dialog`: modal manager and dialog message contracts
- `statusbar`: contextual status line and key help surface
- `panel`: bordered container for child content
- `theme`: theme presets and custom token options
- `core`: shared interfaces and common messages

## Roadmap

Near-term goals:

- command palette
- searchable grouped picker
- richer dialog overlay compositing
- responsive breakpoint helpers
- expanded docs and examples

## Development

```bash
go test ./...
```

## License

TBD
