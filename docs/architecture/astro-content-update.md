# BentoTUI Astro Content Update

Last verified: 2026-03-20
Repository: `github.com/cloudboy-jh/bentotui`
Current docs/app snapshot target: `v0.5.4`

---

## 1) Canonical one-liner

BentoTUI is the fastest way to build full Go TUIs using copy-and-own bricks and recipes, import-only room layouts, and runnable bento templates.

---

## 2) Product framing (use this language in docs/site)

- **Bricks**: official UI building blocks copied into app code with `bento add <brick>`.
- **Recipes**: copy-and-own composed flow patterns copied with `bento add recipe <name>`.
- **Rooms**: stable import-only page layout contracts from `registry/rooms`.
- **Bentos**: runnable template apps in `registry/bentos/*` for rapid remix.

Positioning guardrails:

- Do not describe BentoTUI as a framework shell users must conform to.
- Do not lead with low-level Bubble Tea/Bubbles internals.
- Lead with outcomes: ship a serious app shell in one day.
- Use Bento vocabulary consistently: bricks, recipes, rooms, bentos.

---

## 3) Current release snapshot

### Current version

- Latest tagged docs state: `v0.5.4`.
- `CHANGELOG.md` top release: `0.5.4` dated `2026-03-19`.

### What v0.5.4 fixed (high-level)

- Command palette action ordering is deterministic.
- Dialog lifecycle handling (`dialog.OpenMsg` and `dialog.CloseMsg`) is now explicit in bento update flows.
- `app-shell` theme propagation now updates footer, center deck, and dialog manager consistently.

### Runtime/tooling baseline

- Go module: `go 1.25.2`.
- Core dependencies:
  - `charm.land/bubbletea/v2 v2.0.0-rc.2`
  - `charm.land/bubbles/v2 v2.0.0-rc.1`
  - `charm.land/lipgloss/v2 v2.0.0-beta.3...`
  - `github.com/charmbracelet/ultraviolet`

---

## 4) Installation and CLI commands (current)

### Install

```bash
go get github.com/cloudboy-jh/bentotui
go install github.com/cloudboy-jh/bentotui/cmd/bento@latest
```

### CLI surface

- `bento` (no args): launches interactive TUI installer/launcher.
- `bento init [name]`: scaffold a starter app.
- `bento add <brick...>`: copy brick source into `bricks/<name>/`.
- `bento add recipe <name...>`: copy recipe source into `recipes/<name>/`.
- `bento list`: prints installable bricks and recipes.
- `bento doctor`: checks project health and optional copied bricks.
- `bento version`: prints CLI version.

---

## 5) Current installable catalog

### Bricks (copy-and-own via `bento add`)

1. `surface` - full-terminal paint surface with UV cell buffer
2. `card` - content container (raised default or flat with `Flat()`)
3. `bar` - header/footer row with keybind cards
4. `dialog` - modal manager + confirm/custom/theme picker/command palette
5. `filepicker` - file and directory picker wrapper
6. `list` - scrollable list wrapper
7. `table` - data table wrapper
8. `text` - static text label
9. `input` - single-line text field
10. `badge` - inline status label
11. `kbd` - command/label shortcut pair
12. `wordmark` - themed heading block
13. `select` - single-choice picker
14. `checkbox` - boolean toggle
15. `progress` - horizontal progress bar
16. `package-manager` - sequential install flow with spinner + progress
17. `tabs` - tab row with keyboard input
18. `toast` - stacked notifications
19. `separator` - horizontal/vertical divider

### Recipes (copy-and-own via `bento add recipe`)

1. `filter-bar` - input + status + keybind strip composition
2. `empty-state-pane` - reusable empty-result card content
3. `command-palette-flow` - open palette and route actions
4. `vimstatus` - vim-style statusline with mode/context/clock

Notes for docs sync:

- Some docs pages currently mention only the first 3 recipes; CLI registry includes `vimstatus` and should be treated as source of truth for installable recipes.

---

## 6) Current bento templates and demos

Shipped runnable entries under `registry/bentos/*`:

- `home-screen` - starter entry screen with theme picker/dialog examples.
- `dashboard` - 2x2 metrics/table composition with anchored footer.
- `app-shell` - canonical workspace shell with command palette and theme flow.
- `detail-view` - list/detail split and session card.
- `dashboard-brick-lab` - component showcase and layout test bed.
- `vimstatus-demo` - recipe-driven demo for vim-style status line.

Canonical docs emphasis should remain on `home-screen`, `app-shell`, and `detail-view` as primary remix templates.

---

## 7) Rooms API state (import-only)

Recommended product room contracts:

- `rooms.AppShell(w, h, content, footer)`
- `rooms.SidebarDetail(w, h, sidebarWidth, sidebar, detail, footer)`
- `rooms.Dashboard(w, h, topLeft, topRight, bottomLeft, bottomRight, footer)`
- `rooms.DiffWorkspace(w, h, railWidth, header, fileRail, diffMain, footer)`

Core room primitives include:

- Focus family: `Focus`, `Pancake`, `TopbarPancake`
- Rail/split: `Rail`, `RailFooterStack`, `HSplit`, `VSplit`, `HSplitFooter`
- Multi-pane: `HolyGrail`, `TripleCol`, `DrawerRight`, `DrawerChrome`
- Dashboard: `Dashboard2x2`, `Dashboard2x2Footer`
- Overlay/strip: `Modal`, `BigTopStrip`

Rules to keep in docs:

- Rooms are geometry-only (no theme, no color decisions).
- Compose room output through `surface` in app `View()`.
- `Frame` era APIs are removed; use current room contracts.

---

## 8) Theme engine state

### Model

- `Theme` is an interface (not a struct token bag).
- Bricks accept `WithTheme(t)` and runtime `SetTheme(t)` updates.
- If no explicit theme is set, bricks may fall back to `theme.CurrentTheme()`.

### Built-in presets

16 built-in presets:

`catppuccin-mocha` (default), `catppuccin-macchiato`, `catppuccin-frappe`, `dracula`, `tokyo-night`, `tokyo-night-storm`, `nord`, `bento-rose`, `gruvbox-dark`, `monokai-pro`, `kanagawa`, `rose-pine`, `ayu-mirage`, `one-dark`, `material-ocean`, `github-dark`.

### Global manager (optional)

- `theme.CurrentTheme()`
- `theme.CurrentThemeName()`
- `theme.SetTheme(name)`
- `theme.PreviewTheme(name)`
- `theme.AvailableThemes()`
- `theme.RegisterTheme(name, t)`

### Extended token surface

Theme includes both app UI tokens and newer diff/syntax tokens:

- Diff tokens (added/removed/context backgrounds and intraline highlights)
- Syntax tokens (`SyntaxKeyword`, `SyntaxType`, `SyntaxFunction`, etc.)

This should be reflected in API docs wherever the full theme contract is listed.

---

## 9) Architecture contract to keep explicit

Recommended docs diagram order:

1. `theme` (interface + presets + optional global manager)
2. `theme/styles` (`Row`, `RowClip`, `ClipANSI`)
3. `registry/bricks` (copy-and-own components)
4. `registry/recipes` (copy-and-own compositions)
5. `registry/rooms` (import-only geometry)
6. `surface` (UV buffer compositor)
7. Bubble Tea frame (`tea.NewView`, `AltScreen`, `BackgroundColor`)

Core rendering rule:

- Every row drawn to surface-backed layouts must have explicit background ownership (via `styles.Row(...)` or equivalent width+background rendering) to avoid bleed/gaps.

---

## 10) Policy and guardrails (CI-enforced)

Guardrail tests currently enforce:

- Rooms do not import theme, bricks, or raw bubbles.
- Bricks do not import other bricks.
- Bentos/starter/scaffold/recipes avoid raw `bubbles/*` (except spinner exception).
- Bento `View()` methods do not call `theme.CurrentTheme()` directly; theme should be model-owned state.

Docs implication:

- Present guardrails as a product reliability feature, not internal trivia.

---

## 11) Messaging blocks for Astro pages

### Hero title options

- Build full Go TUIs fast.
- Ship terminal apps in days, not weeks.
- Compose serious TUIs with bricks, rooms, and bentos.

### Hero subtitle

BentoTUI gives Go teams a production-oriented TUI system: copy-and-own bricks and recipes, stable room layout contracts, and runnable bento templates you can remix quickly.

### Primary CTA options

- Run `go run ./registry/bentos/home-screen`
- Start with `app-shell`
- Install a brick with `bento add card`

### Secondary CTA options

- Browse bricks and recipes
- Explore room contracts
- Read the architecture guide

### Positioning bullets

- Not a framework shell you have to fight.
- Not a low-level primitives tutorial.
- A fast path to shipping full Go TUIs with ownership.

---

## 12) Known docs sync gaps to correct in your Astro app

1. Include `vimstatus` in recipe catalog where installable recipes are listed.
2. Keep current release references at `v0.5.4`.
3. Mention diff/syntax methods in the theme interface docs.
4. Keep command palette/theme-picker behavior aligned with latest `app-shell` fixes.

---

## 13) Suggested docs IA (if you want parity with repo docs)

- Overview / Getting Started
- Bricks API reference
- Recipes API reference
- Rooms API reference
- Bentos templates
- Theme engine
- Rendering and coloring rules
- CLI commands (`init`, `add`, `add recipe`, `list`, `doctor`)
- Product direction and roadmap

---

## 14) Source-of-truth files for future updates

- Product overview: `README.md`
- Release notes: `CHANGELOG.md`
- Install catalog: `cmd/bento/logic/add.go`
- CLI behavior: `cmd/bento/main.go`, `cmd/bento/cli.go`
- Theme contract: `theme/theme.go`, `theme/manager.go`, `theme/presets.go`
- Rooms API: `registry/rooms/*.go`, `docs/architecture/rooms.md`
- Guardrails: `internal/policy/guardrails_test.go`
- Existing marketing source: `docs/astro-content.md`
