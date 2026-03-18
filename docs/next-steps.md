# BentoTUI — Next Steps

## Current state (v0.4.0)

Shipped in this pass:

- Theme interface model — `Theme` is a Go interface, 16 presets as plain structs
- Colors-in architecture — every brick accepts `WithTheme` + `SetTheme`, no forced global
- `card` brick merges `panel` + `elevated-card` — `Flat()` option for panel style
- `Frame` + `FrameMainDrawer` + `FrameTriple` removed — use `Pancake`/`Focus`/`TopbarPancake`
- `theme/styles/System` struct deleted — `Row`, `RowClip`, `ClipANSI` as pure functions
- All docs updated to reflect v0.4.0 reality

---

## Immediate priorities

### 1 — Wave 2 bentos

Build and ship these using existing bricks:

- `registry/bentos/form` — labeled inputs, validation hints
- `registry/bentos/log-viewer` — scrollable filter + output
- `registry/bentos/settings` — left nav + settings content pane
- `registry/bentos/command-view` — command-palette-first screen

Keep each bento runnable with `go run ./registry/bentos/<name>` and
copy-and-own friendly.

### 2 — Brick test coverage

- Add snapshot tests for every brick's `View()` output
- Cover: list delegate rendering, card raised/flat modes, bar card truncation,
  dialog frame dimensions, input style sync

### 3 — `bento init` cleanup

- Simplify generated `main.go` to use `card` instead of old `panel`/`elevated-card`
- Add `// bento add card`, `// bento add bar` comments pointing to next steps
- Remove any old token struct references from scaffold template

### 4 — `bento upgrade`

- Diff local copied component against current registry version
- Print a unified diff — no auto-merge, the user decides

### 5 — CLI + TUI usage examples

The theme refactor unlocked using bricks in CLI/non-TUI contexts. Add an
example showing how to render a `card` or `table` to stdout without a
`tea.Program` — demonstrates the breadth of the "TUI and CLI" positioning.

---

## Non-goals (still true)

- No web renderer or browser output
- No mouse-first interaction model
- No built-in app router or page framework
- No data-fetching abstraction layer
