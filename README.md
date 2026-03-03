![BentoTUI logo](./bentotui-readme-logo.png)

# BentoTUI

> [!WARNING]
> Early development — APIs and registry paths will change.

[![Go Version](https://img.shields.io/badge/go-1.23%2B-00ADD8?logo=go)](https://go.dev/)
[![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v2-FF5F87?logo=charm&logoColor=white)](https://github.com/charmbracelet/bubbletea)
[![Status](https://img.shields.io/badge/status-v0.2%20active-6D5EF3)](#status)
[![Changelog](https://img.shields.io/badge/changelog-keep%20a%20changelog-2EA043)](./CHANGELOG.md)

A registry of copy-and-own terminal UI components built on
[Bubble Tea v2](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss v2](https://github.com/charmbracelet/lipgloss).

Run `bento add input` and the source lands in your project. You own it — read it,
modify it, delete what you don't need. No framework lock-in, no lifecycle hooks,
no "extend" API to learn.

## How it works

Two things live in this repo:

| Thing | What it is |
|-------|-----------|
| **component** | Atomic UI piece. `bento add input` copies it into your project. |
| **bentos** | Pre-built layout composition. A complete screen pattern you copy wholesale. |

```
registry/components/   ← copied into your project by `bento add`
bentos/                ← complete runnable screen patterns
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

# Copy components into your project
bento add input bar surface
```

## Components

All components live in `registry/components/` and are copied into your project by `bento add`.
Once copied they live at `yourmodule/components/<name>` — you own the source.

### Available now

| Component | Description |
|-----------|-------------|
| `surface` | Full-screen cell buffer backed by Ultraviolet. Deterministic background paint — no ANSI whitespace bleed. Used by every full-screen layout. |
| `bar` | Status/nav bar with left + right slots. |
| `input` | Single-line text input with left-border accent. Wraps `bubbles/textinput`. |
| `panel` | Titled, focusable content container. |
| `dialog` | Modal manager — `Confirm`, `Custom`, `ThemePicker`. |
| `list` | Scrollable log-style list. |
| `table` | Header + data rows. |
| `text` | Static styled label. |

### Coming next — Tier 1

| Component | Description |
|-----------|-------------|
| `badge` | Inline colored label. |
| `kbd` | Keyboard shortcut display. Dim/bright pair. |
| `wordmark` | Large centered app name. Theme-colored, responsive. |
| `statusbar` | Standalone registry version of `bar`. |

### Tier 2

| Component | Description |
|-----------|-------------|
| `select` | Single-choice picker. Opens inline. |
| `checkbox` | Togglable boolean. |
| `textarea` | Multi-line input. Wraps `bubbles/textarea`. |
| `spinner` | Loading indicator. Wraps `bubbles/spinner`. |
| `progress` | Progress bar with optional label. |

### Tier 3

| Component | Description |
|-----------|-------------|
| `command` | Command palette with fuzzy search. |
| `toast` | Ephemeral notification. Auto-dismisses, stacks. |
| `tabs` | Horizontal tab switcher. |
| `separator` | Horizontal or vertical rule. |

## Bentos

Bentos are complete runnable screen patterns you copy wholesale. Each is a single
self-contained `.go` file demonstrating real component usage.

```
bentos/
  home-screen/   ← mirrors starter app (coming soon)
  app-shell/
  dashboard/
  detail-view/
  form/
  log-viewer/
  settings/
  command-view/
```

### Planned bentos

| Bento | Components used |
|-------|----------------|
| `home-screen` | `wordmark`, `input`, `kbd`, `badge`, `bar`, `surface` |
| `app-shell` | `panel`, `layout`, `bar`, `tabs`, `surface` |
| `dashboard` | `panel`, `badge`, `table`, `layout`, `surface` |
| `detail-view` | `list`, `panel`, `layout`, `surface` |
| `form` | `input`, `textarea`, `checkbox`, `badge`, `surface` |
| `log-viewer` | `input`, `panel`, `spinner`, `badge`, `surface` |
| `settings` | `list`, `panel`, `checkbox`, `layout`, `surface` |
| `command-view` | `command`, `input`, `list`, `surface` |

## Core packages (real imports, not copied)

These are stable module deps your project imports directly:

| Package | Import path | What it is |
|---------|-------------|------------|
| `theme` | `github.com/cloudboy-jh/bentotui/theme` | Global theme store, 16 presets, goroutine-safe |
| `styles` | `github.com/cloudboy-jh/bentotui/styles` | Theme → Lip Gloss style mapping |
| `layout` | `github.com/cloudboy-jh/bentotui/layout` | `Horizontal` / `Vertical` split layout |

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
│  surface  input  bar  dialog  panel  list  table     │
└─────────────────┬───────────────────────────────────┘
                  │ imports
                  ▼
┌─────────────────────────────────────────────────────┐
│  bentotui module deps  (real go imports)             │
│  theme   styles   layout                             │
└─────────────────────────────────────────────────────┘
                  │ built on
                  ▼
┌─────────────────────────────────────────────────────┐
│  Charm stack                                         │
│  Bubble Tea v2   Lip Gloss v2   Ultraviolet          │
└─────────────────────────────────────────────────────┘
```

`surface` is the key rendering primitive. It wraps Ultraviolet's cell buffer
so every cell is explicitly painted before Bubble Tea flushes the frame —
eliminating ANSI whitespace-reset bleed that occurs with pure string composition.

## Run the starter app

```bash
go run ./cmd/starter-app
```

Type `/theme` to switch themes live. Type `/dialog` to open a sample dialog.

## Execution order

1. ✅ **Starter app** — home screen with surface-backed rendering
2. **Tier 1 components** — `badge`, `kbd`, `wordmark`, `statusbar`
3. **`home-screen` bento** — first bento, mirrors starter app
4. **Tier 2 components** — `select`, `checkbox`, `textarea`, `spinner`, `progress`
5. **Remaining bentos** — `app-shell`, `dashboard`, `detail-view`, `form`
6. **Tier 3 components** — `command`, `toast`, `tabs`, `separator`
7. **CLI embed wiring** — `bento add` copies from embedded registry fs
8. **`bento init` template** — single-screen starter template

## License

MIT — see [LICENSE](./LICENSE)
