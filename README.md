![BentoTUI logo](./bentotui-readme-logo.png)

# BentoTUI

> [!WARNING]
> In early production.

[![Go Version](https://img.shields.io/badge/go-1.23%2B-00ADD8?logo=go)](https://go.dev/)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v2-FF5F87?logo=charm&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Status](https://img.shields.io/badge/status-v0.2%20active-6D5EF3)](#status)
[![Changelog](https://img.shields.io/badge/changelog-keep%20a%20changelog-2EA043)](./CHANGELOG.md)

**ShadCN for Go TUIs.** Copy components into your project and own them completely.

BentoTUI is a registry of well-crafted, copy-and-own terminal UI components built
on [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss v2](https://github.com/charmbracelet/lipgloss). No framework lock-in —
you copy the source, modify it freely, and keep the real Charm packages as deps.

## Install

```bash
# Real module deps (theme, styles, layout — these you import, not copy)
go get github.com/cloudboy-jh/bentotui

# CLI to copy components into your project
go install github.com/cloudboy-jh/bentotui/cmd/bento@latest
```

## Quick Start

```bash
bento add panel bar dialog
```

Then in your app:

```go
package main

import (
    "fmt"

    tea "charm.land/bubbletea/v2"
    "github.com/cloudboy-jh/bentotui/layout"
    "github.com/cloudboy-jh/bentotui/theme"

    // Components copied into your project by `bento add`
    "yourmodule/components/bar"
    "yourmodule/components/panel"
)

func main() {
    m := newModel()
    if _, err := tea.NewProgram(m).Run(); err != nil {
        fmt.Printf("error: %v\n", err)
    }
}

type model struct {
    root   *layout.Split
    header *bar.Model
}

func newModel() *model {
    hdr := bar.New(
        bar.Left("my app"),
        bar.Cards(
            bar.Card{Command: "ctrl+c", Label: "quit", Enabled: true},
        ),
    )

    content := panel.New(
        panel.Title("Main"),
        panel.Content(nil), // your widget here
    )

    root := layout.Vertical(
        layout.Fixed(1, hdr),
        layout.Flex(1, content),
    )

    return &model{root: root, header: hdr}
}

func (m *model) Init() tea.Cmd                          { return m.root.Init() }
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if _, ok := msg.(tea.KeyMsg); ok {
        if msg.(tea.KeyMsg).String() == "ctrl+c" {
            return m, tea.Quit
        }
    }
    updated, cmd := m.root.Update(msg)
    m.root = updated.(*layout.Split)
    return m, cmd
}
func (m *model) View() tea.View { return m.root.View() }
```

## Components

| Component | Import (after copy) | What it does |
|-----------|--------------------|-|
| `panel`   | `components/panel` | Titled, focusable content container |
| `bar`     | `components/bar`   | Header/footer row with keybind cards |
| `dialog`  | `components/dialog`| Modal manager + Confirm + Custom + ThemePicker + CommandPalette |
| `list`    | `components/list`  | Scrollable log-style list (plain text) |
| `table`   | `components/table` | Header + data rows |
| `text`    | `components/text`  | Static label |
| `input`   | `components/input` | Single-line text field |

## Real Deps (not copied)

These you import directly from the module — they are the stable contract:

| Package | Import path | What it is |
|---------|-------------|------------|
| `theme` | `github.com/cloudboy-jh/bentotui/theme` | Global theme store, 15 presets, goroutine-safe |
| `styles` | `github.com/cloudboy-jh/bentotui/styles` | theme → Lip Gloss style mapping |
| `layout` | `github.com/cloudboy-jh/bentotui/layout` | `Horizontal` / `Vertical` split layout |

## Theme System

15 professional presets via [bubbletint](https://github.com/lrstanley/bubbletint):

```go
// In main(), before tea.NewProgram().Run()
theme.SetTheme("dracula")

// Components always call this in View() — never cache it
t := theme.CurrentTheme()
```

Available themes: `catppuccin-mocha` (default), `catppuccin-macchiato`,
`catppuccin-frappe`, `dracula`, `tokyo-night`, `tokyo-night-storm`, `nord`,
`gruvbox-dark`, `monokai-pro`, `kanagawa`, `rose-pine`, `ayu-mirage`,
`one-dark`, `material-ocean`, `github-dark`.

Custom themes:

```go
theme.RegisterTheme("my-theme", theme.Theme{
    Surface: theme.SurfaceTokens{Canvas: "#0d0d0d", Panel: "#1a1a1a", ...},
    // ... all tokens required
})
```

## Architecture

Three layers — you only touch the middle one:

```
your app code
     │
     ▼
┌────────────────────────────────────────────┐
│  registry components  (you own the source) │
│  panel  bar  dialog  list  table  input    │
└────────────┬───────────────────────────────┘
             │ imports
             ▼
┌────────────────────────────────────────────┐
│  bentotui module deps  (real go imports)   │
│  theme   styles   layout                   │
└────────────────────────────────────────────┘
```

See [docs/architecture.md](./docs/architecture.md) for more detail.

## Run the Showcase

```bash
go run ./cmd/starter-app
```

Shows every component with live theme switching (`ctrl+t`), command palette
(`ctrl+p`), and dialogs (`ctrl+d`).

## Documentation

- [docs/components.md](./docs/components.md) — API reference for every component
- [docs/architecture.md](./docs/architecture.md) — Design principles and layer diagram
- [docs/next-steps.md](./docs/next-steps.md) — Known gaps and what to build next
- [docs/roadmap.md](./docs/roadmap.md) — Longer-term plans

## License

MIT — see [LICENSE](./LICENSE)
