# BentoTUI Roadmap

## Current State (v0.2)

The repository completed a full structural refactor in v0.2 and has a working
baseline CLI + starter flow:

- Deleted the monolithic `bentotui.New()` framework API
- Deleted `app/`, `core/`, `ui/`, `bentotui.go`
- Moved `core/theme/` → `theme/`, `ui/styles/` → `styles/`, `core/layout/` → `layout/`
- Created `registry/` with clean rewrites of every component
- Every component reads `theme.CurrentTheme()` in `View()` — no stored theme state
- Every row uses a single `lipgloss.NewStyle().Background().Width().Render()` call
- Starter app rewritten as a full component showcase with live theme switching
- `registry/embed.go` now embeds registry component files for installer logic
- `bento add`, `bento init`, `bento list`, and `bento doctor` are implemented

See [next-steps.md](./next-steps.md) for the current priorities.

## Backlog

### CLI (`cmd/bento`)

- [x] `bento add` — copy component files from embedded registry
- [x] `bento init` — generate a runnable starter project
- [x] `bento list` — show available components with one-line descriptions
- [ ] `bento upgrade <component>` — diff registry version against copied version
- [ ] `bento wrap` command group — deterministic scaffold pipeline

### Wrap + AI Scaffold

- [ ] Deterministic manifest spec — define and validate a stable wrap manifest format used by scaffold generation
- [ ] `bento wrap --manifest-only` — parse interface data and emit deterministic manifest JSON
- [ ] `bento wrap --scaffold` — generate owned Go scaffold files from a manifest
- [ ] Optional enhance pass — `bento wrap --enhance` applies an LLM refinement pass after deterministic scaffold
- [ ] MCP server tools — expose `bento_wrap`, `bento_scaffold`, `bento_enhance`
- [ ] `bento.Enhance()` API — first-class API for post-scaffold enhancement workflows
- [ ] `llms.txt` — ship model context for enhancement and scaffold tooling

### New Components

- [x] `tabs` — keyboard-navigable tab bar with panel content area
- [x] `select` — dropdown/select widget (wraps or reimplements bubbles/list)
- [x] `progress` — horizontal progress bar with theme colors
- [x] `checkbox` — boolean toggle
- [x] `badge` — inline colored label (useful inside panel titles)
- [x] `kbd` — keyboard shortcut display pair
- [x] `wordmark` — large title/header display helper
- [x] `toast` — ephemeral stacked notifications
- [x] `separator` — horizontal/vertical divider

Already shipped in registry: `surface`, `panel`, `bar`, `dialog`, `list`, `table`, `text`, `input`.

Primitive policy: Bento does not plan a `spinner` registry component; use
`charm.land/bubbles/v2/spinner` directly.

### Bento Examples

- [ ] `bentos/home-screen` — canonical copy-and-own starter screen
- [ ] `bentos/app-shell` — sidebar + body + status layout pattern
- [ ] `bentos/dashboard` — cards + table composition pattern
- [ ] `bentos/form` — form controls and validation flow pattern

### Layout Enhancements

- [ ] `layout.Grid` — fixed-column grid (simpler than manual Horizontal+Vertical nesting)
- [ ] Scrollable `Split` — allow body regions to scroll independently

### Developer Experience

- [ ] Layout debugger mode — render allocation boundaries as colored overlays
- [ ] `go test ./registry/...` — snapshot tests for every component's rendered output
- [ ] More starter-app variants — e.g. file browser, log viewer
- [ ] `bento add --force` or guided overwrite mode for existing copied components

### Theme

- [ ] Light theme support — adapter needs a `fromLightTint` variant with inverted surface mapping
- [ ] Theme validation CLI — `bento validate-theme mytheme.json`

## Non-Goals

- **Mobile / small screens** — BentoTUI assumes a reasonably large terminal.
  Components have minimum widths but no responsive breakpoint system.
- **Mouse support** — Bubble Tea v2 has mouse events; none of the current
  components handle them. Not planned unless a component clearly needs it (e.g. scrollbars).
- **Accessibility** — Terminal screen reader support depends on the terminal
  emulator, not the TUI library. No plans here.

---

Last updated: 2026-03-09
