![BentoTUI logo](./bentotui-readme-logo.png)

# BentoTUI

> [!WARNING]
> Early development — APIs and registry paths will change.

[![Go Version](https://img.shields.io/badge/go-1.25%2B-00ADD8?logo=go)](https://go.dev/)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v2-FF5F87?logo=charm&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Status](https://img.shields.io/badge/status-main%20active-6D5EF3)](#status)
[![Changelog](https://img.shields.io/badge/changelog-keep%20a%20changelog-2EA043)](./CHANGELOG.md)

The **best way to build full Go TUIs fast** — opinionated components, named
layout rooms, and template apps you can remix quickly.

Run `bento add card` and the source lands in your project. You own it — read it,
modify it, delete what you do not need.

Build apps with:

- **bricks** for UI components
- **rooms** for page layout contracts
- **bentos** for full app templates

## How it works

Three building blocks:

| Thing | What it is |
|---|---|
| **bricks** | UI components — copied into your project via `bento add`. You own them. |
| **rooms** | Named layout contracts for pages. Imported, not copied. Zero color, zero theme. |
| **bentos** | Complete app templates you can clone/remix for production. |

```
registry/bricks/    ← copy-and-own via `bento add`
registry/rooms/     ← import directly from the module
registry/bentos/    ← runnable full-screen app patterns
```

## Install

```bash
go get github.com/cloudboy-jh/bentotui

go install github.com/cloudboy-jh/bentotui/cmd/bento@latest
```

## Quick start

```bash
# Run the home screen
go run ./registry/bentos/home-screen

# Other bentos
go run ./registry/bentos/dashboard
go run ./registry/bentos/app-shell
go run ./registry/bentos/detail-view

# Copy bricks into your project
bento add card bar input dialog list table surface
```

Home-screen demo:

![BentoTUI home-screen demo](./demo.gif)

## Bricks

All bricks live in `registry/bricks/` and are copied by `bento add`.
Once copied they live at `yourmodule/bricks/<name>` — you own the source.

Every brick accepts a theme at construction and supports live updates:

```go
// Pass a theme at construction (optional — falls back to global if omitted)
c := card.New(card.Title("file.go"), card.WithTheme(theme.Preset("dracula")))

// Update on theme change
c.SetTheme(newTheme)
```

| Brick | Description |
|---|---|
| `surface` | Full-screen Ultraviolet cell buffer. Deterministic background paint. |
| `card` | Content container — raised (chrome band) or flat (titled pane) via `Flat()`. Replaces `panel` + `elevated-card`. |
| `bar` | Header/footer row with keybind cards, status pill, priority-aware overflow. |
| `input` | Single-line text field for command bars and forms. |
| `dialog` | Modal manager — `Confirm`, `Custom`, `ThemePicker`, `CommandPalette`. |
| `list` | Scrollable list with sections and structured rows. |
| `table` | Header + data rows with compact/borderless/grid modes. |
| `badge` | Inline status label — neutral, info, success, warning, danger, accent. |
| `tabs` | Keyboard-navigable tab row. |
| `kbd` | Keyboard shortcut pair (`command label`). |
| `select` | Single-choice inline picker backed by `bubbles/list`. |
| `checkbox` | Boolean toggle with `bubbles/key` bindings. |
| `progress` | Horizontal progress bar backed by `bubbles/progress`. |
| `filepicker` | File/directory picker wrapping `bubbles/filepicker`. |
| `toast` | Stacked transient notifications. |
| `separator` | Horizontal or vertical divider. |
| `text` | Static themed label. |
| `wordmark` | Themed heading/title block. |
| `package-manager` | Sequential install flow with spinner + progress. |

Primitive policy: Bento does not ship a `spinner` brick. Use `charm.land/bubbles/v2/spinner` directly.

## Bentos

Complete runnable screen patterns in `registry/bentos/`. Copy and own wholesale.

| Bento | Description |
|---|---|
| `home-screen` | Starter-style entry screen — wordmark, input, theme picker |
| `dashboard` | Dense card/table composition — 2×2 metric grid |
| `app-shell` | Rail + workspace + command palette + theme switching |
| `detail-view` | List + detail split pane |
| `dashboard-brick-lab` | Component showcase — list/table/filepicker/progress in cards |

Use these as template baselines: keep the room contract, replace data and
interactions with your own domain.

## Rooms

Import once per page file and choose a room function for that page.

```go
import "github.com/cloudboy-jh/bentotui/registry/rooms"

screen := rooms.AppShell(w, h, content, footer)
screen := rooms.SidebarDetail(w, h, 26, nav, detail, footer)
screen := rooms.DiffWorkspace(w, h, 28, header, files, diff, footer)
```

Lower-level composition functions (`HSplit`, `VSplit`, `HolyGrail`, etc.) remain
available for advanced layouts.

## Stable imports

Three packages you import directly (not copied):

```go
// Theme interface + 16 presets + optional global manager
import "github.com/cloudboy-jh/bentotui/theme"

// Row/RowClip/ClipANSI rendering utilities
import "github.com/cloudboy-jh/bentotui/theme/styles"

// Layout contracts and geometry — AppShell, SidebarDetail, DiffWorkspace, ...
import "github.com/cloudboy-jh/bentotui/registry/rooms"
```

Everything else is copy-and-own.

## Usage policy

Use Bento defaults first:

- Use `bento add` bricks for UI primitives you want to own
- Use `registry/rooms` to pick room contracts per page
- Use `theme` tokens for colors and state
- Keep raw `bubbles/*` usage out of bentos unless there is a documented gap (current exception: spinner)

## Theme system

`Theme` is a Go interface. 16 built-in presets. No mandatory global store.
Pass themes as inputs to bricks — or use the global manager if you prefer:

```go
// Named preset — no global state
t := theme.Preset("tokyo-night")
card := card.New(card.Title("file.go"), card.WithTheme(t))

// Global manager (optional)
theme.SetTheme("dracula")
t := theme.CurrentTheme()
```

Available presets: `catppuccin-mocha` (default), `catppuccin-macchiato`,
`catppuccin-frappe`, `dracula`, `tokyo-night`, `tokyo-night-storm`, `nord`,
`bento-rose`, `gruvbox-dark`, `monokai-pro`, `kanagawa`, `rose-pine`,
`ayu-mirage`, `one-dark`, `material-ocean`, `github-dark`.

Custom themes: embed `theme.BaseTheme`, fill the color fields, implement
`theme.Theme`. Register with `theme.RegisterTheme("name", t)`.

## Minimal app example

```go
package main

import (
    tea "charm.land/bubbletea/v2"
    "github.com/cloudboy-jh/bentotui/registry/bricks/surface"
    "github.com/cloudboy-jh/bentotui/registry/rooms"
    "github.com/cloudboy-jh/bentotui/theme"
    "yourmodule/bricks/bar"
    "yourmodule/bricks/card"
    "yourmodule/bricks/list"
)

type model struct {
    theme   theme.Theme
    footer  *bar.Model
    content *card.Model
    log     *list.Model
    w, h    int
}

func newModel() *model {
    t := theme.CurrentTheme()
    l := list.New(100)
    l.Append("ready")
    return &model{
        theme:   t,
        log:     l,
        content: card.New(card.Title("Output"), card.Content(l), card.WithTheme(t)),
        footer:  bar.New(
            bar.FooterAnchored(),
            bar.Left("my-app"),
            bar.Cards(bar.Card{Command: "ctrl+c", Label: "quit", Enabled: true}),
            bar.WithTheme(t),
        ),
    }
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.w, m.h = msg.Width, msg.Height
    case theme.ThemeChangedMsg:
        m.theme = msg.Theme
        m.content.SetTheme(m.theme)
        m.footer.SetTheme(m.theme)
    case tea.KeyMsg:
        if msg.String() == "ctrl+c" {
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m *model) View() tea.View {
    t := m.theme
    surf := surface.New(m.w, m.h)
    surf.Fill(t.Background())
    surf.Draw(0, 0, rooms.Focus(m.w, m.h, m.content, m.footer))
    v := tea.NewView(surf.Render())
    v.AltScreen = true
    v.BackgroundColor = t.Background()
    return v
}
```

## CLI

- `bento init` — scaffold a runnable starter project
- `bento add <brick...>` — copy brick source into `bricks/<name>/`
- `bento list` — list available bricks with descriptions
- `bento doctor` — environment and project checks

## Architecture

```
theme/ (interface + 16 preset structs)
     │  Preset("name"), CurrentTheme(), SetTheme(), RegisterTheme()
     ▼
theme/styles/ (Row / RowClip / ClipANSI)
     ▼
registry/bricks/ (copy-and-own)
     │  card  bar  dialog  input  list  table  surface  + more
     │  Each brick: WithTheme(t) + SetTheme(t)
     ▼
registry/rooms/ (named geometry — Focus, Rail, HolyGrail, ...)
     ▼
surface (Ultraviolet cell buffer — Fill → Draw → DrawCenter → Render)
     ▼
Bubble Tea v2 (tea.NewView, AltScreen, BackgroundColor)
```

## Docs

- [docs/README.md](./docs/README.md) — index
- [docs/architecture/architecture.md](./docs/architecture/architecture.md) — rendering contract, theme model, component rules
- [docs/architecture/bricks.md](./docs/architecture/bricks.md) — brick API reference
- [docs/architecture/rooms.md](./docs/architecture/rooms.md) — room layout API
- [docs/architecture/bentos.md](./docs/architecture/bentos.md) — full app composition
- [docs/theme-engine.md](./docs/theme-engine.md) — theme interface, presets, custom themes
- [docs/coloring-rules.md](./docs/coloring-rules.md) — rules for correct color usage in bricks
- [docs/usage-guide.md](./docs/usage-guide.md) — use Bento defaults first, layering/import rules
- [docs/astro-content.md](./docs/astro-content.md) — website/marketing copy source

## License

MIT — see [LICENSE](./LICENSE)
