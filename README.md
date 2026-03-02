![BentoTUI logo](./bentotui-readme-logo.png)

# BentoTUI

> [!WARNING]
> In early production.

[![Go Version](https://img.shields.io/badge/go-1.23%2B-00ADD8?logo=go)](https://go.dev/)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v2-FF5F87?logo=charm&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Status](https://img.shields.io/badge/status-v0.1%20active-6D5EF3)](#status)
[![Changelog](https://img.shields.io/badge/changelog-keep%20a%20changelog-2EA043)](./CHANGELOG.md)

BentoTUI is an application framework on top of Bubble Tea for building production-grade terminal apps with a structured shell, deterministic layering, and reusable UI components.

Charm gives you bricks. BentoTUI gives you rooms.

## Status

BentoTUI is in active `v0.1` development with a stable core architecture.

**Completed:**

- Canvas-based layout system (Horizontal/Vertical with Fixed/Flex/Min/Max constraints)
- Global theme system with 15 professional presets via bubbletint
- Complete widget library (Input, List, Text, Card, Table)
- Container components (Panel, Bar, Dialog with theme picker)
- Reactive theme propagation across all components

**Current focus:**

- Additional container components (Tabs, Sidebar, Split)
- Widget enhancements (ScrollableList, RichText)
- Performance optimizations
- Documentation and examples

## Feature Snapshot

- shell model with explicit layer order (`header -> body -> footer -> scrim -> dialog`)
- lazy page router and page factories
- **canvas-based layout** (Horizontal/Vertical with Fixed/Flex constraints)
- focus ring and keyboard routing
- modal dialog manager (confirm, custom, theme picker)
- semantic theme presets (15 themes via bubbletint)
- structured UI layer:
  - `core/layout` - Canvas-based positioning (Horizontal, Vertical)
  - `ui/containers` - Complex components (Panel, Bar, Dialog)
  - `ui/widgets` - Simple content components (Card, Input, List, Table, Text)
  - `ui/styles` - Theme-to-lipgloss mapping

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
	"github.com/cloudboy-jh/bentotui/core/layout"
	"github.com/cloudboy-jh/bentotui/core/theme"
	"github.com/cloudboy-jh/bentotui/ui/containers/panel"
	"github.com/cloudboy-jh/bentotui/ui/widgets"
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
	root *layout.Split
}

func newHomePage() *homePage {
	// Create content widgets
	infoList := widgets.NewList(50)
	infoList.Append("Page: home")
	infoList.Append("Status: Ready")

	// Create panels (containers with borders)
	leftPanel := panel.New(
		panel.Title("Sidebar"),
		panel.Content(infoList),
	)

	centerPanel := panel.New(
		panel.Title("Main"),
		panel.Content(widgets.NewText("Center content")),
	)

	rightPanel := panel.New(
		panel.Title("Details"),
		panel.Content(widgets.NewText("Right content")),
	)

	// Arrange horizontally with gutter spacing
	root := layout.Horizontal(
		layout.Flex(1, leftPanel),
		layout.Flex(2, centerPanel),
		layout.Flex(1, rightPanel),
	).WithGutterColor(theme.CurrentTheme().Border.Subtle)

	return &homePage{root: root}
}

func (p *homePage) Init() tea.Cmd {
	return p.root.Init()
}

func (p *homePage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return p.root.Update(msg)
}

func (p *homePage) View() tea.View {
	return p.root.View()
}

func (p *homePage) SetSize(w, h int) {
	p.root.SetSize(w, h)
}

func (p *homePage) GetSize() (int, int) {
	return p.root.GetSize()
}

func (p *homePage) Title() string {
	return "Home"
}
```

## Architecture

Four clear layers:

1. **Layout** (`core/layout/`) - Canvas-based positioning using Horizontal/Vertical
2. **Containers** (`ui/containers/`) - Complex components with borders (Panel, Bar, Dialog)
3. **Widgets** (`ui/widgets/`) - Simple content components (Card, Input, List, Table, Text)
4. **Your App** - Composes layouts, containers, and widgets

```go
// Layout positions containers (which hold widgets)
root := layout.Horizontal(
    layout.Flex(1, panel.New(
        panel.Title("Info"),
        panel.Content(widgets.NewList()),
    )),
    layout.Flex(2, panel.New(
        panel.Title("Main"),
        panel.Content(widgets.NewInput()),
    )),
).WithGutterColor(theme.Border.Subtle)
```

See [docs/components.md](./docs/components.md) for complete component documentation and [docs/architecture.md](./docs/architecture.md) for system design.

## License

MIT - see [LICENSE](./LICENSE)
