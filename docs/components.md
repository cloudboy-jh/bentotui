---
title: "BentoTUI Components"
description: "Complete guide to BentoTUI components and layouts"
---

# BentoTUI Components

BentoTUI provides a clear separation between layout, containers, and widgets.

## Architecture

Four layers:

1. **Layout** (`core/layout/`) - Canvas-based positioning
2. **Containers** (`ui/containers/`) - Complex components (Panel, Bar, Dialog)
3. **Widgets** (`ui/widgets/`) - Simple content (Card, Input, List, Table, Text)
4. **Styles** (`ui/styles/`) - Theme mapping

```go
// Layout positions containers (which hold widgets)
root := layout.Horizontal(
    layout.Flex(1, panel.New(
        panel.Title("Info"),
        panel.Content(widgets.NewList()),
    )),
).WithGutterColor(theme.Border.Subtle)
```

## Layout (`core/layout`)

Canvas-based positioning system.

### Horizontal

Arranges children horizontally with proper canvas layering.

```go
layout.Horizontal(
    layout.Fixed(30, sidebar),    // exactly 30 cells wide
    layout.Flex(1, mainContent),  // takes remaining space
    layout.Flex(1, rightPanel),   // takes half of remaining
).WithGutterColor("#585B70")      // gutter between panels
```

### Vertical

Arranges children vertically.

```go
layout.Vertical(
    layout.Fixed(3, header),   // exactly 3 rows
    layout.Flex(1, body),      // takes remaining height
    layout.Fixed(1, footer),   // exactly 1 row
)
```

### Constraints

- `layout.Fixed(n)` - Exactly n cells
- `layout.Flex(w)` - Proportional share (weight w)

### Gutter

Use `WithGutterColor(color)` to add spacing between children:

```go
layout.Horizontal(...).WithGutterColor(theme.Border.Subtle)
```

## Containers (`ui/containers`)

### Panel

Bordered container that holds widgets as content.

```go
infoPanel := panel.New(
    panel.Title("Info"),
    panel.Content(widgets.NewList()),
)

// With options
panel.New(
    panel.Title("My Panel"),
    panel.Content(myWidget),
    panel.Elevated(),  // uses elevated surface color
)
```

### Bar

Header or footer bar with card system.

```go
footer := bar.New(
    bar.LeftCard(bar.Card{
        Command: "bento",
        Variant: bar.CardMuted,
    }),
    bar.Cards(
        bar.Card{Command: "/", Label: "commands"},
        bar.Card{Command: "tab", Label: "focus"},
    ),
    bar.RightCard(bar.Card{
        Command: "Status",
        Variant: bar.CardPrimary,
    }),
)
```

### Dialog

Modal dialog system.

```go
// Show a dialog
dialog.Open(dialog.Custom{
    DialogTitle: "Confirm",
    Content:     myContentWidget,
    Width:       60,
    Height:      10,
})
```

## Widgets (`ui/widgets`)

Simple content components that render inside containers.

### Card

Key/value badge for keybindings.

```go
card := widgets.NewCard("/", "commands")
// Renders: "/ commands" with accent styling
```

### Input

Text input field.

```go
input := widgets.NewInput()
input.SetValue("Type here...")
input.Focus()  // returns tea.Cmd
```

### Text

Static text display.

```go
text := widgets.NewText("Hello World")
text.SetText("Updated text")
```

### List

Scrollable list of items.

```go
list := widgets.NewList(100)  // max 100 items
list.Append("Item 1")
list.Append("Item 2")
list.Prepend("Item 0")  // adds to beginning
list.Clear()
```

### Table

Data table with headers.

```go
table := widgets.NewTable("Name", "Status", "Date")
table.AddRow("Task 1", "Done", "2024-01-01")
table.AddRow("Task 2", "Pending", "2024-01-02")
```

## Complete Example: 3-Column Layout

```go
package main

import (
    "github.com/cloudboy-jh/bentotui/core/layout"
    "github.com/cloudboy-jh/bentotui/core/theme"
    "github.com/cloudboy-jh/bentotui/ui/containers/panel"
    "github.com/cloudboy-jh/bentotui/ui/widgets"
)

func buildLayout() *layout.Split {
    // Left panel with list
    leftPanel := panel.New(
        panel.Title("Sidebar"),
        panel.Content(widgets.NewList()),
    )

    // Center panel with input
    centerPanel := panel.New(
        panel.Title("Main"),
        panel.Content(widgets.NewInput()),
    )

    // Right panel with text
    rightPanel := panel.New(
        panel.Title("Details"),
        panel.Content(widgets.NewText("Info...")),
    )

    // Arrange with gutter
    return layout.Horizontal(
        layout.Flex(1, leftPanel),
        layout.Flex(2, centerPanel),
        layout.Flex(1, rightPanel),
    ).WithGutterColor(theme.CurrentTheme().Border.Subtle)
}
```

## Theming

All containers and widgets respect the current theme.

```go
// Set theme on containers that support it
panel.SetTheme(myTheme)
input.SetTheme(myTheme)

// Use theme tokens for gutters
layout.Horizontal(...).WithGutterColor(theme.Border.Subtle)
```

Available themes (via bubbletint):
- `catppuccin-mocha`, `catppuccin-macchiato`, `catppuccin-frappe`
- `dracula`
- `tokyo-night`, `tokyo-night-storm`
- `nord`
- `gruvbox-dark`
- `kanagawa`
- `rose-pine`
- `one-dark`
- `monokai-pro`
- `material-ocean`
- `ayu-mirage`
- `github-dark`

## Key Principles

1. **Layout positions containers** - Use `layout.Horizontal()` / `layout.Vertical()`
2. **Containers hold widgets** - Panels have borders, widgets don't
3. **Use constraints** - `Fixed(n)` or `Flex(w)`, no manual positioning
4. **Gutter for spacing** - `WithGutterColor()` not hardcoded gaps
5. **Canvas-based rendering** - Layout uses proper layer compositing

## Import Summary

```go
// Layout (positioning)
"github.com/cloudboy-jh/bentotui/core/layout"

// Containers (complex components)
"github.com/cloudboy-jh/bentotui/ui/containers/panel"
"github.com/cloudboy-jh/bentotui/ui/containers/bar"
"github.com/cloudboy-jh/bentotui/ui/containers/dialog"

// Widgets (content)
"github.com/cloudboy-jh/bentotui/ui/widgets"

// Styles (theming)
"github.com/cloudboy-jh/bentotui/ui/styles"
```
