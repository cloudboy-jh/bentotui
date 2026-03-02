# BentoTUI — Main Spec

Status: Active
Date: 2026-02-25

> 🍱 The application framework for Bubble Tea. Compartmentalized layouts, composable components, shipped apps.

---

## Overview

BentoTUI is a Go framework that sits between Charm's Bubble Tea primitives and production TUI applications. Every polished terminal app built today — Crush, OpenCode, lazygit, k9s — independently reinvents the same architectural patterns: layout systems, focus management, dialog overlays, command palettes, and component composition. BentoTUI extracts these patterns into importable, composable building blocks.

Charm gives you bricks. BentoTUI gives you rooms.

```
import "github.com/cloudboy-jh/bentotui"
```

### Implementation Update (Current)

- v0.1 foundation is implemented in code (`core/shell`, `core/router`, `core/layout`, `core/focus`, `core/theme`, `core/surface`, `ui/containers/dialog`, `ui/containers/footer`, `ui/containers/panel`, `ui/primitives`)
- Public facade package `bentotui` remains the recommended app entrypoint; `app` currently aliases `core/shell` for compatibility
- Rendering moved from plain string concatenation to styled surfaces with Lip Gloss v2
- Horizontal composition now uses ANSI-aware joining to avoid escape-sequence width drift
- Dialogs are rendered through a layer/canvas composition path and centered in the app shell
- Internal starter app added at `cmd/starter-app` for daily framework regression checks

---

## What This Is

- An application-level framework built on top of Bubble Tea
- A layout and composition system for terminal UIs
- A set of higher-order components that every TUI app needs
- The missing layer between `bubbletea.Model` and a shipped product

## What This Is Not

- A replacement for Bubble Tea, Bubbles, or Lip Gloss
- A widget library (Bubbles already does that)
- A styling library (Lip Gloss already does that)
- A terminal rendering engine
- A React/declarative UI system (that's OpenTUI's lane)

---

## Architecture

### The Gap

```
BUBBLE TEA              BENTOTUI                    YOUR APP
(primitives)            (application patterns)      (domain logic)
─────────────           ──────────────────          ──────────────
tea.Model          →    App shell + router      →   Pages
tea.Msg            →    Message routing tree    →   Business events  
tea.WindowSizeMsg  →    Responsive layout       →   Panel config
lipgloss.Style     →    Theme system            →   Brand colors
lipgloss.Join*     →    Panel composition       →   Layout definition
lipgloss.Layer     →    Dialog/overlay system   →   Modals
bubbles/*          →    Enhanced components     →   Domain widgets
(nothing)          →    Focus management        →   Navigation
(nothing)          →    Command palette         →   App commands
(nothing)          →    Footer bar              →   Context help
```

### Component Model

Every BentoTUI component implements a layered interface system:

```go
// Core — every component
type Component interface {
    tea.Model
}

// Sizeable — responds to terminal resize
type Sizeable interface {
    Component
    SetSize(width, height int)
    GetSize() (width, height int)
}

// Focusable — participates in focus system
type Focusable interface {
    Component
    Focus()
    Blur()
    IsFocused() bool
}

// Positional — can be placed at coordinates
type Positional interface {
    Component
    SetPosition(x, y int)
}

// Bindable — registers keybindings
type Bindable interface {
    Component
    Bindings() []key.Binding
}
```

Components opt into interfaces as needed. A simple label implements `Component`. A panel implements `Component + Sizeable`. An input implements all of them.

### Message Routing

Messages flow top-down through the component tree:

```
App
├── Router (page switching)
│   ├── Page A
│   │   ├── Panel (focused) ← receives input
│   │   ├── Panel
│   └── Page B
│       └── ...
├── DialogManager ← captures input when active
│   └── Active Dialog
└── FooterBar
```

Each level decides to handle, forward, or ignore. When a dialog is open, it captures all input. When closed, input flows to the focused component in the active page.

---

## Modules

### Core

#### `shell` — Application Shell (Primary)
The root model that bootstraps everything. Manages lifecycle, router, dialog manager, and footer bar. The shell renders full-size themed surfaces and uses Lip Gloss v2 canvas layers for overlay composition.

`app` remains available as a compatibility alias over `shell`.

```go
app := bentotui.New(
    bentotui.WithTheme(myTheme),
    bentotui.WithPages(
        bentotui.Page("home", newHomePage),
        bentotui.Page("settings", newSettingsPage),
    ),
    bentotui.WithFooterBar(true),
)

p := tea.NewProgram(app)

// Full-screen app mode is enabled by default.
// Opt out for inline mode:
// bentotui.WithFullScreen(false)
```

#### `router` — Page System
Manages page navigation with lazy loading. Pages are created on first visit and cached.

```go
// Navigate via message
func switchPage() tea.Msg {
    return router.Navigate("settings")
}

// Pages implement the Component interface
type Page interface {
    Component
    Sizeable
    Title() string
}
```

#### `layout` — Responsive Panel System
Compartmentalized layouts — the bento box itself. Define panels with flex ratios that respond to terminal size.

```go
// Horizontal split: sidebar (fixed 30) | main (flex)
layout := layout.Horizontal(
    layout.Fixed(30, sidebar),
    layout.Flex(1, mainContent),
)

// Vertical split: header (fixed 1) | body (flex) | editor (fixed 5)
layout := layout.Vertical(
    layout.Fixed(1, header),
    layout.Flex(1, body),
    layout.Fixed(5, editor),
)

// Nested — bento grid
layout := layout.Horizontal(
    layout.Fixed(30, sidebar),
    layout.Vertical(
        layout.Flex(1, messages),
        layout.Fixed(5, editor),
    ),
)
```

Responsive helpers are still planned at framework level. Current compact-mode behavior is exercised in `cmd/starter-app`.

#### `focus` — Focus Management
Handles focus cycling between components, visual indicators, and input routing.

```go
focus := focus.New(
    focus.Ring(editor, messages, sidebar),
    focus.Keys(
        key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next panel")),
        key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev panel")),
    ),
)

// Focus state affects rendering
if panel.IsFocused() {
    style = style.BorderForeground(appTheme.Border.Focus)
}
```

#### `theme` — Color System
Coordinated color definitions with presets and semantic surface tokens for shell/component rendering.

```go
// Use a registered preset
appTheme := theme.Preset("catppuccin-mocha")

// Or switch active theme from the registry
applied, err := theme.SetTheme("osaka-jade")
if err != nil {
    // unknown or invalid theme
}

// Themes use strict semantic token groups:
// surface, text, border, state, selection, input, bar, dialog
_ = applied
```

### Components

#### `ui/containers/dialog` — Modal Overlay System
Modal dialogs that capture input and render above content via lipgloss layers.

```go
// Open a dialog via message
func confirmDelete() tea.Msg {
    return dialog.Open(dialog.Confirm{
        DialogTitle: "Delete secret?",
        Message: "This cannot be undone.",
        OnConfirm: func() tea.Msg { return deleteMsg{} },
    })
}

// Custom dialog content
func modelPicker() tea.Msg {
    return dialog.Open(dialog.Custom{
        DialogTitle: "Select model",
        Content:     myPickerComponent,
        Width:       60,
        Height:      20,
    })
}
```

#### `palette` — Command Palette (Planned, v0.2)
Slash-command / fuzzy-search command palette overlay.

```go
palette := palette.New(
    palette.Command("/new", "New session", newSessionCmd),
    palette.Command("/models", "Switch model", openModelPickerCmd),
    palette.Command("/help", "Help", openHelpCmd),
    palette.Command("/exit", "Exit", tea.Quit),
)

// Trigger with a keybinding
// User types "/" → palette opens → fuzzy filter → select → execute
```

#### `picker` — Searchable Grouped Picker (Planned, v0.2)
A list picker with sections, search, selection highlighting, and keyboard hints. The model picker from OpenCode/Crush, extracted.

```go
picker := picker.New(
    picker.Group("Recent",
        picker.Item("GPT-5.3 Codex", "OpenAI"),
        picker.Item("Claude Sonnet 4.5", "Anthropic"),
    ),
    picker.Group("Anthropic",
        picker.Item("Claude Haiku 4.5", ""),
        picker.Item("Claude Opus 4.5", ""),
    ),
    picker.WithSearch(true),
    picker.WithKeys("ctrl+m"),
)
```

#### `ui/containers/footer` — Context-Aware Footer Bar
Bottom layer showing command cards, status messages, and contextual info.

```go
footer := footer.New(
    footer.LeftCard(footer.Card{Command: "/dialog", Label: "open dialog", Variant: footer.CardNormal, Enabled: true}),
    footer.Cards(footer.Card{Command: "/theme", Label: "switch theme", Variant: footer.CardPrimary, Enabled: true}),
    footer.RightCard(footer.Card{Command: "/page", Label: "next page", Variant: footer.CardMuted, Enabled: true}),
)
```

#### `ui/containers/panel` — Bordered Content Panel
A themed surface container with optional title, focused border state, and content sizing.

```go
panel := panel.New(
    panel.Theme(theme.Preset("catppuccin-mocha")),
    panel.Title("Messages"),
    panel.Content(viewport),
    panel.Scrollable(true),
)
```

### Utilities (Planned)

#### `keys` — Keybinding Management (Planned)
Registration, conflict detection, and help generation.

#### `events` — Typed Event Bus (Planned)
Generic typed pub/sub for component communication beyond tea.Msg.

#### `size` — Terminal Size Utilities (Planned)
Breakpoint helpers, responsive size calculation, compact mode detection.

---

## Design Principles

1. **Additive** — sits on top of Bubble Tea, never wraps or replaces it
2. **Opt-in** — use the full app shell or cherry-pick individual components
3. **Explicit** — no magic, all patterns are visible and debuggable
4. **Minimal deps** — Bubble Tea + Bubbles + Lip Gloss, nothing else
5. **Dogfood-driven** — validated by real apps (Veil, Churn, Porter CLI)
6. **Convention over configuration** — sensible defaults, override everything

---

## Build Targets

BentoTUI builds on Bubble Tea v2 (beta) and Lip Gloss v2 (beta). This is a forward bet — Crush is already on v2, and the Layer/Canvas API in Lip Gloss v2 is required for proper dialog compositing.

```
go 1.23+

require (
    charm.land/bubbletea/v2
    charm.land/bubbles/v2
    charm.land/lipgloss/v2
)
```

---

## Validation Plan

### First App: Veil

Veil (encrypted secrets manager TUI) is built on BentoTUI as the first real consumer. It exercises:

| Veil Feature | BentoTUI Module |
|---|---|
| Home / Projects / Settings pages | `router` |
| Sidebar + main content layout | `layout` |
| Secrets table with grouped sections | `layout` + `ui/containers/panel` |
| Add/edit/import overlays | `ui/containers/dialog` |
| Init wizard | `ui/containers/dialog` (multi-step) |
| Project tab navigation | `focus` |
| Command card bar | `ui/containers/footer` |
| Catppuccin/Dracula/Osaka Jade themes | `theme` |

### Internal Harness (Current)

`cmd/starter-app` is the active internal validation surface for framework behavior and rendering quality.

Run it with:

```bash
go run ./cmd/starter-app
```

It currently validates:

- single-page shell composition (`header` + `main input` + `footer`)
- theme switching via `/theme` (legacy aliases still accepted)
- dialog overlays via `/dialog`
- focus handling between input and action controls (`tab`, `shift+tab`)
- primitive-first rendering behavior on fullscreen alt-screen

### Ecosystem Apps (Planned)

| Tool | Description | Key BentoTUI Modules |
|---|---|---|
| Pretty Log | Clean log output viewer | `layout`, `ui/containers/panel`, `theme` |
| File Tree | Terminal file navigator | `layout`, `focus`, `keys` |
| Diff Viewer | Side-by-side diff display | `layout`, `ui/containers/panel`, `theme` |

---

## Scope — v0.1

The minimum surface to ship and build Veil on:

- [x] `app` — application shell with lifecycle
- [x] `router` — page switching with lazy load
- [x] `layout` — horizontal/vertical splits with fixed/flex
- [x] `focus` — focus ring with tab cycling
- [x] `theme` — color system with presets + semantic surface tokens
- [x] `ui/containers/dialog` — modal overlay with confirm/custom
- [x] `ui/containers/footer` — command cards + themed footer surface
- [x] `ui/containers/panel` — themed bordered content container with focus state

**Not in v0.1:**
- Command palette (v0.2)
- Picker component (v0.2)
- Responsive breakpoints (v0.2)
- Event bus (v0.2)

---

## Name & Branding

- **Name:** BentoTUI
- **Emoji:** 🍱
- **Repo:** `github.com/cloudboy-jh/bentotui`
- **Tagline:** The application framework for Bubble Tea
- **Logo:** Bento box icon — playful, warm tones, recognizable at small sizes

---

## References

- [Rendering System Design (ADR-0001)](../architecture/rendering-system-design.md)
- [Component System Reference](../architecture/component-system-reference.md)
- [Component Sizing Contract](../architecture/component-sizing-contract.md)
- [Bento Color Design System](../design/bento-color-design-system.md)
- [Implementation Next Steps](../roadmap/next-steps.md)
- [TUI Framework Research Doc](../research/tui-framework-research.md)
- [Crush TUI Architecture (DeepWiki)](https://deepwiki.com/charmbracelet/crush/5.1-tui-architecture)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Bubbles](https://github.com/charmbracelet/bubbles)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- [OpenTUI](https://github.com/sst/opentui)
