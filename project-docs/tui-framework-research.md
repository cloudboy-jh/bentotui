# TUI Application Framework Research
## Building the missing layer between Bubble Tea and shipped apps

**Author:** Jack Horton  
**Date:** February 23, 2026  
**Status:** Archived research baseline (superseded by implementation docs)

This document is a pre-implementation research snapshot.
For current framework behavior and roadmap, use:
- `project-docs/bentotui-main-spec.md`
- `project-docs/next-steps.md`

Implementation note (2026-02-26):
- shared UI render primitives now live in `ui/primitives`
- harness/footer behavior is card-first with command cards (`/pr`, `/issue`, `/branch`)

---

## 1. Executive Summary

Every polished TUI application built on Charm's Bubble Tea framework (Crush, OpenCode, lazygit, k9s) independently reinvents the same architectural patterns: layout systems, focus management, dialog overlays, command palettes, and component composition. None of these patterns are provided by Bubble Tea or its companion libraries. This framework fills that gap.

**The thesis:** There is a well-defined architectural layer between "Bubble Tea primitives" and "shipped application" that nobody has packaged into a reusable, importable framework.

---

## 2. The Charm Ecosystem — What Exists Today

### 2.1 Bubble Tea (39.7k stars)
The core TUI framework. Implements the Elm Architecture (Init/Update/View) in Go. Provides:
- Program lifecycle management
- Message-based event loop
- Framerate-based renderer
- Mouse support, focus reporting
- Alternate screen buffer
- Inline and fullscreen modes

**What it does NOT provide:** Any opinion on how to compose components, manage focus between them, handle layouts, or build multi-page applications.

### 2.2 Bubbles (companion library)
Pre-built components implementing `tea.Model`:
- `textinput` — single-line text input
- `textarea` — multi-line text input
- `viewport` — scrollable content area
- `list` — filterable, paginated item list
- `table` — tabular data display
- `spinner` — loading indicator
- `progress` — progress bar
- `paginator` — pagination logic
- `filepicker` — file system browser
- `help` — auto-generated keybinding help
- `key` — keybinding management
- `timer` / `stopwatch` — time tracking
- `cursor` — cursor management

**Gap:** These are individual widgets. There is no system for composing them into application layouts, managing focus between them, or layering them (e.g., dialog over content).

### 2.3 Lip Gloss (10.6k stars)
Declarative styling library (CSS-like for terminals):
- Text styling: colors, bold, italic, underline, strikethrough
- Box model: padding, margins, borders
- Alignment: horizontal and vertical
- Joining: `JoinHorizontal()`, `JoinVertical()`
- Placement: `Place()` for positioning in whitespace
- **v2 beta additions:** `Layer`, `Canvas` for compositing overlays
- Tables, Trees, Lists as rendering components

**Gap:** No flexbox layout (requested in issue #166, still open). No responsive layout system. `JoinHorizontal`/`JoinVertical` are manual — the developer calculates all sizes. The `stickers` third-party library attempted CSS flexbox for lipgloss but is unmaintained.

### 2.4 Other Charm Libraries
| Library | Purpose | Relevance |
|---------|---------|-----------|
| Harmonica | Spring animation | Nice-to-have for transitions |
| BubbleZone | Mouse event tracking per component | Essential for click handling |
| Gum | Shell script TUI helpers | Different use case |
| VHS | Terminal recording | Tooling, not framework |
| Log | Pretty logging | Tangential |

---

## 3. Crush (formerly OpenCode) — The Reference Implementation

### 3.1 Background
Crush is Charm's AI coding agent (20.3k stars). Originally built as "OpenCode" by Kujtim Hoxha. Acquired by Charm, leading to a fork war with Dax (SST) and Adam. Charm's version was renamed to Crush; the original continues as OpenCode under SST.

Both repos share a common ancestor and independently evolved the same TUI patterns. This makes them the two best reference implementations for what the framework needs to extract.

### 3.2 Crush TUI Architecture (from source analysis)

**Root:** `internal/tui/tui.go` — `appModel` struct
```
appModel
├── pages[]         — page system (currently chat only, extensible)
├── dialog          — modal dialog overlay manager
├── completions     — contextual autocomplete popup
├── status          — status bar / help text
├── currentPage     — page routing
├── width/height    — terminal dimensions
└── app             — reference to backend services
```

**Page System:** `internal/tui/page/chat/chat.go` — `chatPage` struct
```
chatPage
├── editor          — multi-line input with completions
├── messages        — scrollable message viewport
├── sidebar         — session list / file tracker
├── header          — compact mode header
├── focusedPane     — focus state (editor | messages | sidebar)
└── compactMode     — responsive layout flag
```

**Layout Constants:**
| Constant | Value | Purpose |
|----------|-------|---------|
| CompactModeWidthBreakpoint | 120 | Activates compact layout |
| CompactModeHeightBreakpoint | 30 | Activates compact layout |
| EditorHeight | 5 | Fixed editor area |
| SideBarWidth | 31 | Fixed sidebar width |
| HeaderHeight | 1 | Compact mode header |

### 3.3 Component Interface System
Crush defines several interfaces that components implement:

| Interface | Methods | Purpose |
|-----------|---------|---------|
| `util.Model` | `Init()`, `Update()`, `View()` | Standard Bubble Tea lifecycle |
| `layout.Sizeable` | `SetSize(w, h)`, `GetSize()` | Responsive resize handling |
| `layout.Focusable` | `Focus()`, `Blur()`, `IsFocused()` | Input focus management |
| `layout.Positional` | `SetPosition(x, y)` | Absolute positioning |
| `layout.Help` | `Bindings()` | Keybinding registration |
| `util.Cursor` | `Cursor()` | Cursor position reporting |

**This interface system is the skeleton of the framework.** It's how Crush solves component composition — but it's locked in `internal/`.

### 3.4 Message Routing Pattern
Messages flow hierarchically: `appModel` → `page` → `component`

Each level decides to:
1. **Handle** — consume the message
2. **Forward** — pass to child
3. **Ignore** — return unchanged

Key message categories:
- **Input:** `tea.KeyPressMsg`, `tea.MouseClickMsg`, `tea.PasteMsg`
- **Window:** `tea.WindowSizeMsg` (broadcast to all)
- **Navigation:** `PageChangeMsg`
- **Dialog:** `OpenDialogMsg`, `CloseDialogMsg`
- **Completions:** `OpenCompletionsMsg`, `CloseCompletionsMsg`
- **Status:** `InfoMsg`, `ClearStatusMsg`
- **App events:** `pubsub.Event[T]` from backend services

### 3.5 Rendering Pipeline
```
1. Page renders content → pages[current].View()
2. Status bar appends → status.View()
3. Base layer created → lipgloss.NewLayer(page + status)
4. Dialog layer added → if dialog active, lipgloss.NewLayer(dialog)
5. Completions layer → if completions open, lipgloss.NewLayer(popup)
6. Canvas composited → lipgloss.NewCanvas(layers...)
7. Cursor positioned → from active focused component
8. tea.View returned
```

### 3.6 Focus Management
Three focus states in chat page: `editor`, `messages`, `sidebar`
- Tab/Shift+Tab cycles focus
- Focused component receives keyboard input
- Blurred components ignore input
- Focus affects rendering (visual indicators)

### 3.7 Dialog System
Modal overlays that:
- Capture all input when active
- Render above base content via lipgloss layers
- Implement `dialogs.DialogModel` interface
- Support open/close via message passing
- Include: permission dialogs, model picker, command palette, session picker

---

## 4. OpenTUI — The TypeScript Alternative

### 4.1 Architecture
- **Core:** Zig-based native rendering engine with TypeScript bindings (`@opentui/core`)
- **Reconcilers:** SolidJS (`@opentui/solid`) and React (`@opentui/react`)
- **Approach:** Declarative/reactive UI, similar to React Native for terminals
- **Stars:** 6.2k
- **Status:** "Not ready for production use"
- **Maintained by:** SST (Dax's company) — same people who forked OpenCode

### 4.2 Why It's Too Niche
1. **TypeScript/Zig dependency** — alienates the Go TUI community (which is the largest)
2. **Requires Zig installed** — heavy build dependency
3. **React/Solid paradigm** — doesn't match Go developers' mental model
4. **Early stage** — 504 commits, 40 open issues, not production-ready
5. **Tight coupling to SST ecosystem** — feels like an internal tool

### 4.3 What It Gets Right
- The declarative component model is the right abstraction level
- The idea of providing both primitives AND composed patterns
- `bun create tui` — scaffolding for new projects
- AI coding assistant integration (opentui-skill)

### 4.4 Takeaways for Our Framework
- Ship in Go, for Go developers — that's where the ecosystem lives
- Provide scaffolding/templates for quick starts
- Make it AI-friendly (good docs, clear patterns)
- Don't require exotic build dependencies

---

## 5. Shipped TUI Apps — Pattern Analysis

### 5.1 Common Patterns Across Apps

Analyzed: Crush, OpenCode, lazygit, k9s, PUG (Terraform TUI)

| Pattern | Crush | lazygit | k9s | PUG |
|---------|-------|---------|-----|-----|
| Multi-panel layout | ✓ | ✓ | ✓ | ✓ |
| Focus management | ✓ | ✓ | ✓ | ✓ |
| Command palette / slash commands | ✓ | ✓ | ✓ | - |
| Modal dialogs | ✓ | ✓ | ✓ | - |
| Searchable list/picker | ✓ | ✓ | ✓ | ✓ |
| Status bar | ✓ | ✓ | ✓ | ✓ |
| Responsive/compact mode | ✓ | - | ✓ | ✓ |
| Keybinding help | ✓ | ✓ | ✓ | ✓ |
| Theming/color system | ✓ | ✓ | ✓ | - |
| Split/resizable panes | - | ✓ | ✓ | ✓ |

### 5.2 What Every App Rebuilds From Scratch
1. **Application shell** — root model with page/view routing
2. **Layout engine** — responsive panel composition with size calculation
3. **Focus system** — cycling, visual indicators, input routing
4. **Dialog/modal overlays** — layered rendering, input capture
5. **Command palette** — slash commands, fuzzy search, keyboard navigation
6. **Searchable picker** — grouped items, sections, selection state
7. **Status bar** — context-aware help, status messages, keybinding hints
8. **Theme system** — coordinated colors, dark/light, accent colors

---

## 6. Gap Analysis — Framework Surface

### 6.1 What Charm Provides vs. What Apps Need

```
CHARM PROVIDES              APPS NEED (THE GAP)              APPS BUILD
─────────────               ──────────────────               ──────────
tea.Model interface    →    Component composition system  →   Custom interfaces
tea.Msg routing        →    Hierarchical message routing  →   Manual switch statements
tea.WindowSizeMsg      →    Responsive layout engine      →   Manual size math
lipgloss.Style         →    Theme system                  →   Global style vars
lipgloss.Join*         →    Panel/grid layout             →   Manual join calls
lipgloss.Layer (v2)    →    Dialog/overlay system         →   Custom layer management
bubbles/list           →    Searchable grouped picker     →   Custom list wrapper
bubbles/textarea       →    Rich editor with completions  →   Extended textarea
bubbles/help           →    Context-aware keybinding help →   Custom help views
bubbles/viewport       →    Scrollable content panel      →   Extended viewport
(nothing)              →    Focus management system       →   Custom focus state
(nothing)              →    Command palette               →   Built from scratch
(nothing)              →    Application shell/router      →   Custom appModel
(nothing)              →    Page/view system              →   Custom page switching
```

### 6.2 Proposed Framework Modules

**Core:**
- `app` — Application shell, lifecycle, root model
- `layout` — Responsive panel system (horizontal/vertical splits, flex ratios)
- `focus` — Focus management, cycling, visual indicators
- `router` — Page/view routing with lazy loading
- `theme` — Coordinated color system with presets

**Components:**
- `dialog` — Modal overlay system with input capture
- `palette` — Command palette with fuzzy search
- `picker` — Searchable grouped item picker
- `statusbar` — Context-aware status bar with help hints
- `editor` — Enhanced textarea with completions support
- `panel` — Bordered content panel with title and controls

**Utilities:**
- `keys` — Keybinding registration and conflict detection
- `events` — Typed event bus for component communication
- `size` — Terminal size utilities, breakpoints, responsive helpers

---

## 7. Competitive Landscape

| Project | Language | Stars | Level | Status |
|---------|----------|-------|-------|--------|
| Bubble Tea | Go | 39.7k | Low-level framework | Active, stable |
| Bubbles | Go | ~8k | Widget library | Active, stable |
| Lip Gloss | Go | 10.6k | Styling library | Active, v2 beta |
| OpenTUI | TS/Zig | 6.2k | Full framework | Early, not production |
| Ratatui | Rust | ~12k | Widget library | Active, stable |
| Textual | Python | ~25k | Full framework | Active, production |
| tview | Go | ~11k | Widget library | Active but dated |
| **This framework** | **Go** | **-** | **Application framework** | **Planned** |

**Key insight:** Textual (Python) is the closest analogue to what we're building — it provides application-level patterns on top of rendering primitives. But nothing like it exists in the Go ecosystem, which is where the most TUI development is happening.

---

## 8. Design Principles

1. **Additive, not replacement** — sits on top of Bubble Tea, doesn't fork or wrap it
2. **Opt-in complexity** — use the full shell or pick individual components
3. **Zero magic** — all patterns are explicit, debuggable, traceable
4. **Minimal dependencies** — Bubble Tea, Bubbles, Lip Gloss, nothing else
5. **Dogfood-driven** — built for and validated by real tools (Churn, Pact, Porter CLI)
6. **Convention over configuration** — sensible defaults, override everything

---

## 9. Open Questions

1. **Bubble Tea v2 timing** — v2 is in beta. Do we build on v1 stable or v2 beta? Crush is on v2.
2. **Lip Gloss v2 Layer/Canvas** — critical for dialog system. Beta-only. Risk?
3. **Naming** — needs a name that signals "application framework for Bubble Tea"
4. **Scope of v0.1** — what's the minimum viable surface to ship?
5. **Flexbox** — do we build our own or wait for lipgloss to ship it?

---

## 10. Recommended Next Steps

1. **Clone and audit** Crush's `internal/tui/` directory line by line — map every component to generic vs. domain-specific
2. **Clone and audit** OpenCode's `internal/tui/` for divergence since the fork
3. **Prototype** the layout system and focus manager as standalone packages
4. **Build a demo app** — simple multi-panel TUI using only the framework
5. **Define the public API** — what does `import "github.com/cloudboy-jh/[framework]"` look like?

---

## Sources

- Crush repo: https://github.com/charmbracelet/crush
- Crush TUI architecture: https://deepwiki.com/charmbracelet/crush/5.1-tui-architecture
- OpenCode repo: https://github.com/opencode-ai/opencode (Dax/Adam fork)
- OpenTUI repo: https://github.com/sst/opentui
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Bubbles: https://github.com/charmbracelet/bubbles
- Lip Gloss: https://github.com/charmbracelet/lipgloss
- Lip Gloss flexbox issue: https://github.com/charmbracelet/lipgloss/issues/166
- Stickers (flexbox attempt): https://github.com/76creates/stickers
- PUG tips: https://leg100.github.io/en/posts/building-bubbletea-programs/
