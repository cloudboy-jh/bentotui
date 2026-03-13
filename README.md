![BentoTUI logo](./bentotui-readme-logo.png)

# BentoTUI

> [!WARNING]
> Early development — APIs and registry paths will change.

[![Go Version](https://img.shields.io/badge/go-1.25%2B-00ADD8?logo=go)](https://go.dev/)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v2-FF5F87?logo=charm&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Status](https://img.shields.io/badge/status-v0.3.3%20active-6D5EF3)](#status)
[![Changelog](https://img.shields.io/badge/changelog-keep%20a%20changelog-2EA043)](./CHANGELOG.md)

A registry of copy-and-own terminal UI components built on
[Bubble Tea v2](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss v2](https://github.com/charmbracelet/lipgloss).

Run `bento add input` and the source lands in your project. You own it — read it,
modify it, delete what you don't need. No framework lock-in, no lifecycle hooks,
no "extend" API to learn.

Layout composition is handled by `registry/layouts`. Home/starter screens use
`Focus(...)` (body + anchored footer), while multi-row app shells can use
`Frame(...)` (top row, subheader row, body, subfooter). Final frame painting
and overlays are handled by `surface`.

Bar rows support role-aware rendering (`top`, `subheader`, `footer`) with
optional anchored footer mode for command focus.

## How it works

Two things live in this repo:

| Thing | What it is |
|-------|-----------|
| **component** | Atomic UI piece. `bento add input` copies it into your project. |
| **bentos** | Pre-built layout composition. A complete screen pattern you copy wholesale. |

`registry/layouts` is also shipped in this module as importable layout primitives
for composing screens (not copied by `bento add`).

```
registry/components/   ← copied into your project by `bento add`
registry/bentos/       ← complete runnable screen patterns
registry/layouts/      ← imported layout primitives (not copied)
```

## Install

```bash
# Core deps — these you import, not copy
go get github.com/cloudboy-jh/bentotui

# CLI to copy components into your project
go install github.com/cloudboy-jh/bentotui/cmd/bento@latest
```

## Quick Start

```bash
# See the home screen
go run ./cmd/starter-app

# Run bento examples
go run ./registry/bentos/home-screen
go run ./registry/bentos/app-shell
go run ./registry/bentos/dashboard

# Copy components into your project
bento add input bar surface
```

Home-screen demo:

![BentoTUI home-screen demo](./demo.gif)

## Components

All components live in `registry/components/` and are copied into your project by `bento add`.
Once copied they live at `yourmodule/components/<name>` — you own the source.

### Available now

| Component | Description |
|-----------|-------------|
| `surface` | Full-screen cell buffer backed by Ultraviolet. Deterministic background paint — no ANSI whitespace bleed. Used by every full-screen layout. |
| `bar` | Role-aware status/nav row with `StatusPill`, compact cards, anchored footer mode, and priority-aware overflow. |
| `input` | Single-line text input with left-border accent. Wraps `bubbles/textinput`. |
| `panel` | Titled, focusable content container. |
| `dialog` | Modal manager — `Confirm`, `Custom`, `ThemePicker`. |
| `list` | Scrollable list with optional sections and row formatting hooks. |
| `table` | Header + data rows with compact/borderless and per-column width/align. |
| `text` | Static styled label. |
| `badge` | Inline themed label. |
| `kbd` | Keyboard shortcut command + label pair. |
| `wordmark` | Themed heading/title block. |
| `select` | Single-choice inline picker. |
| `checkbox` | Boolean toggle input. |
| `progress` | Horizontal progress bar. |
| `tabs` | Keyboard-navigable tab row. |
| `toast` | Stacked notification rows. |
| `separator` | Horizontal or vertical divider. |

`bento add` currently supports: `surface`, `panel`, `bar`, `dialog`, `list`, `table`, `text`, `input`, `badge`, `kbd`, `wordmark`, `select`, `checkbox`, `progress`, `tabs`, `toast`, `separator`.

Primitive policy: Bento does not ship a `spinner` registry component. Use
`charm.land/bubbles/v2/spinner` directly.

## Bentos

Bentos are complete runnable screen patterns you copy wholesale. Each bento in
`registry/bentos/` is a self-contained `main.go` demonstrating real component usage.

```
registry/bentos/
  home-screen/
  app-shell/
  dashboard/
```

### Current bentos

| Bento | Components used |
|-------|----------------|
| `home-screen` | `wordmark`, `input`, `kbd`, `badge`, `bar`, `surface` |
| `app-shell` | `panel`, `bar`, `tabs`, `surface` |
| `dashboard` | `panel`, `badge`, `table`, `bar`, `surface` |

Future wave examples (`detail-view`, `form`, `log-viewer`, `settings`, `command-view`) remain planned.

## Core packages (real imports, not copied)

These are stable module deps your project imports directly:

| Package | Import path | What it is |
|---------|-------------|------------|
| `theme` | `github.com/cloudboy-jh/bentotui/theme` | Global theme store, 16 presets, goroutine-safe |
| `styles` | `github.com/cloudboy-jh/bentotui/styles` | Theme → Lip Gloss style mapping |
| `layouts` | `github.com/cloudboy-jh/bentotui/registry/layouts` | Named visual layout grammar (`Focus`, `Frame`, splits, dashboard, modal, and more) |

`surface` remains a copy-and-own registry component (`bento add surface`) and is
the recommended final compositor for full-frame paint (`Fill`) and overlays (`DrawCenter`).

## Theme System

16 built-in presets:

```go
// Set once at startup
theme.SetTheme("tokyo-night")

// Components always call this in View() — never cache it
t := theme.CurrentTheme()
```

Available themes: `catppuccin-mocha` (default), `catppuccin-macchiato`,
`catppuccin-frappe`, `dracula`, `tokyo-night`, `tokyo-night-storm`, `nord`,
`bento-rose`, `gruvbox-dark`, `monokai-pro`, `kanagawa`, `rose-pine`, `ayu-mirage`,
`one-dark`, `material-ocean`, `github-dark`.

Switch themes at runtime via the `/theme` command in the starter app.

## Architecture

```
your app code
     │
     ▼
┌─────────────────────────────────────────────────────┐
│  registry/components/  (you own the source)          │
│  surface  panel  bar  dialog  input  list  table      │
│  text  badge  kbd  wordmark  select  checkbox         │
│  progress  tabs  toast  separator                     │
└─────────────────┬───────────────────────────────────┘
                  │ imports
                  ▼
┌─────────────────────────────────────────────────────┐
│  bentotui module deps  (real go imports)             │
│  theme   styles   registry/layouts                   │
└─────────────────────────────────────────────────────┘
                  │ built on
                  ▼
┌─────────────────────────────────────────────────────┐
│  Charm stack                                         │
│  Bubble Tea v2   Lip Gloss v2   Ultraviolet          │
└─────────────────────────────────────────────────────┘
```

Render contract:

1. Build structure with `registry/layouts` (typically `Focus` for home/starter screens).
2. Paint the full frame with `surface.Fill(theme.CurrentTheme().Surface.Canvas)`.
3. Draw layout output with `surface.Draw(0, 0, screen)`.
4. Draw overlays/dialogs with `surface.DrawCenter(...)`.

This avoids ANSI whitespace gaps and keeps theme canvas rendering deterministic.

## Run the starter app

```bash
go run ./cmd/starter-app
```

Type `/theme` to switch themes live. Type `/dialog` to open a sample dialog.

## CLI status

- `bento init` scaffolds a runnable project
- `bento add <component...>` copies registry component source into `components/<name>/`
- `bento list` shows available components
- `bento doctor` runs environment/project checks

## Next roadmap slice

1. Expand `registry/bentos/` beyond the shipped first wave (`home-screen`, `app-shell`, `dashboard`)
2. Expand bento catalog breadth (`detail-view`, `form`, `log-viewer`, `settings`, `command-view`)
3. Keep primitive policy strict (use established Bubbles primitives like `spinner` directly)
4. Tighten `bento init` template output and guidance comments
5. Add component + CLI logic tests (`go test ./registry/...` and command smoke coverage)
6. Start deterministic `bento wrap` manifest/scaffold pipeline

## License

MIT — see [LICENSE](./LICENSE)
