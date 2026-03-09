# BentoTUI — Next Steps

## Current state: v0.2

Core deps (`theme`, `layout`, `panel`, `statusbar`) are working. Starter app in progress. Registry concept established.

---

## Naming

| Term | What it is |
|------|-----------|
| **component** | Atomic UI piece. `bento add input` copies it into your project. |
| **bento** | Pre-built layout composition. A complete screen pattern, ready to run. |

```
registry/components/   ← what bento add copies
bentos/                ← complete screen patterns users copy wholesale
```

---

## 1. Starter App

**File:** `cmd/starter-app/main.go`
**Goal:** Single-screen app that looks like the opencode home screen. Ships with the repo as the first thing anyone runs.

**Components it demonstrates:**
- Wordmark (large ASCII centered text via lipgloss)
- Accented input panel (left-border-only panel variant)
- Badge row (dim inline component names)
- Kbd hints (right-aligned shortcut display)
- Tip line (colored dot + dim text)
- Status bar (bottom, left + right content)

---

## 2. Components

Live in `registry/components/` and are copied into user projects via `bento add <name>`.

Each component:
- Only imports `bubbletea`, `lipgloss`, and `bentotui/theme` + `bentotui/styles`
- Reads `theme.CurrentTheme()` at render time — no stored theme field
- One lipgloss call per row — background + foreground + width in a single chain
- Returns plain string from `View()`

### Tier 1 — Ship with starter app (build these first)

| Component | Description |
|-----------|-------------|
| `input` | Single-line text input with left-border accent. Wraps `bubbles/textinput`. |
| `badge` | Inline colored label. Used in input rows, statusbars, anywhere. |
| `kbd` | Keyboard shortcut display. Styled dim/bright pair — `tab`, `⌘K`. |
| `statusbar` | Bottom bar with left + right slots. Registry version of existing package. |
| `wordmark` | Large centered app name. Lipgloss bold + theme color, responsive to width. |

### Tier 2 — Core interactive components

| Component | Description |
|-----------|-------------|
| `list` | Scrollable selectable list. Keyboard nav. |
| `table` | Data table with headers + row selection. Column widths via layout constraints. |
| `select` | Single-choice picker. Opens inline. |
| `checkbox` | Togglable boolean item. Used in forms and settings. |
| `textarea` | Multi-line text input. Wraps `bubbles/textarea`. |
| `spinner` | Loading indicator. Wraps `bubbles/spinner`, theme colored. |
| `progress` | Progress bar. Percentage + optional label. |

### Tier 3 — Overlay + complex components

| Component | Description |
|-----------|-------------|
| `dialog` | Modal with confirm/cancel. Clean registry version of existing package. |
| `command` | Command palette with fuzzy search. `⌘K` trigger, filterable action list. |
| `toast` | Ephemeral notification. Auto-dismisses, stacks, theme colored. |
| `tabs` | Horizontal tab switcher. Active tab highlighted, keyboard navigable. |
| `separator` | Horizontal or vertical rule. Theme border color. |

---

## 3. Bentos

**Folder:** `bentos/`

Pre-built layout compositions users copy wholesale. Each bento is a single self-contained `.go` file — a complete runnable screen demonstrating real component usage in a named pattern.

```
bentos/
  app-shell/    app_shell.go
  dashboard/    dashboard.go
  detail-view/  detail_view.go
  form/         form.go
  log-viewer/   log_viewer.go
  home-screen/  home_screen.go
  settings/     settings.go
  command-view/ command_view.go
```

### Bento list

| Bento | Description | Components used |
|-------|-------------|-----------------|
| `home-screen` | Wordmark + input + tips (opencode style) | `wordmark`, `input`, `kbd`, `badge`, `statusbar` |
| `app-shell` | Header + sidebar + main + statusbar | `panel`, `layout`, `statusbar`, `tabs` |
| `dashboard` | 2-col stats cards + data table | `panel`, `badge`, `table`, `layout` |
| `detail-view` | Sidebar list + detail pane | `list`, `panel`, `layout` |
| `form` | Centered form with labeled inputs | `input`, `textarea`, `checkbox`, `badge` |
| `log-viewer` | Full-width scrollable output + filter input | `input`, `panel`, `spinner`, `badge` |
| `settings` | Two-col: nav list + settings panel | `list`, `panel`, `checkbox`, `layout` |
| `command-view` | Fullscreen command palette | `command`, `input`, `list` |

Each bento:
- Is a complete runnable app (`go run ./bentos/dashboard`)
- Has a comment header explaining what it demonstrates
- Uses only registry components + bentotui core deps
- Is under 150 lines

---

## 4. CLI — `bento add` embed wiring

`cmd/bento/add.go` needs `//go:embed registry` so all component source is baked into the binary.

```go
//go:embed registry
var registryFS embed.FS

// bento add panel
// → reads registry/components/panel/panel.go from embed
// → writes to ./components/panel/panel.go in user's project
```

Steps:
1. All registry components exist and compile
2. Add `//go:embed registry` to `add.go`
3. Write extraction logic: create `./components/<name>/` dir, write files
4. Test: `bento add input` in a fresh project produces a working file

---

## 5. `bento init` template update

The generated `main.go` from `bento init` should match the starter app pattern — not the old framework pattern.

Generated template should:
- Be a single-screen app (no router, no page system)
- Import only `bubbletea`, `lipgloss`, `bentotui/theme`, `bentotui/layout`
- Show a placeholder panel with "Your app goes here"
- Be under 60 lines
- Have comments pointing to `bento add` for next steps

---

## 6. `bento wrap` foundation

**Goal:** Rapidly bootstrap a BentoTUI app from an existing command interface.

Pipeline:
1. Parse `--help`, flags, and sample output
2. Produce deterministic manifest
3. Generate owned Go scaffold files

Initial deliverables:
- Manifest schema + validator
- `bento wrap --manifest-only`
- `bento wrap --scaffold`

---

## 7. AI integration surface

**Goal:** Add optional enhancement on top of deterministic scaffolding.

Deliverables:
- LLM enhance pass (`bento wrap --enhance`)
- MCP server tools: `bento_wrap`, `bento_scaffold`, `bento_enhance`
- First-class `bento.Enhance()` API
- `llms.txt` for model context

---

## Execution order

1. **Starter app** — opencode prompt running now
2. **Tier 1 components** — `input`, `badge`, `kbd`, `wordmark` (statusbar already exists)
3. **`home-screen` bento** — first bento, mirrors starter app, validates the pattern
4. **Tier 2 components** — `list`, `table`, `select` unblock the most bentos
5. **Remaining bentos** — `app-shell`, `dashboard`, `detail-view` next
6. **Tier 3 components** — `command`, `dialog`, `toast` last
7. **CLI embed wiring** — once all registry components exist
8. **`bento init` template** — final cleanup
9. **`bento wrap` foundation** — deterministic manifest + scaffold
10. **AI integration surface** — MCP tools, `bento.Enhance()`, `llms.txt`

---

## Non-goals (for now)

- No web renderer or browser output
- No mouse-only interactions (keyboard-first always)
- No animation system beyond spinner
- No built-in routing (users write their own)
- No data fetching utilities
- No `bento add bento <name>` CLI command — copy bentos manually for now
