# BentoTUI Documentation

## Quick Start

- [Components Guide](./components.md) - Complete API reference for layouts, containers, and widgets
- [Architecture](./architecture.md) - Layer architecture and design principles
- [Roadmap](./roadmap.md) - Future plans and next steps

## Overview

BentoTUI provides four layers:

1. **Layout** (`core/layout/`) - Canvas-based positioning
2. **Containers** (`ui/containers/`) - Complex components (Panel, Bar, Dialog)
3. **Widgets** (`ui/widgets/`) - Simple content (Card, Input, List, Table, Text)
4. **Styles** (`ui/styles/`) - Theme mapping

## Usage Example

```go
package main

import (
    "github.com/cloudboy-jh/bentotui/core/layout"
    "github.com/cloudboy-jh/bentotui/core/theme"
    "github.com/cloudboy-jh/bentotui/ui/containers/panel"
    "github.com/cloudboy-jh/bentotui/ui/widgets"
)

func main() {
    // Build 3-column layout
    root := layout.Horizontal(
        layout.Flex(1, panel.New(
            panel.Title("Sidebar"),
            panel.Content(widgets.NewList()),
        )),
        layout.Flex(2, panel.New(
            panel.Title("Main"),
            panel.Content(widgets.NewInput()),
        )),
    ).WithGutterColor(theme.CurrentTheme().Border.Subtle)
}
```

See [components.md](./components.md) for complete API documentation.
