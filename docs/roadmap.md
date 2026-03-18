# BentoTUI Roadmap

## Current state (v0.4.0)

### What shipped in v0.4.0

- **Theme interface model** — `Theme` is a Go interface (opencode-style). 16 presets
  as plain Go structs. No bubbletint runtime dependency. No contrast validation engine.
  No global store required by bricks.
- **Colors-in architecture** — every brick accepts `WithTheme(t)` at construction and
  `SetTheme(t)` for live updates. Bricks fall back to `theme.CurrentTheme()` only when
  no theme was explicitly provided. No brick calls `CurrentTheme()` unconditionally.
- **`panel` + `elevated-card` merged into `card`** — one content-container brick.
  `card.New(...)` = raised (chrome band). `card.New(card.Flat(), ...)` = flat titled
  container. The old naming ambiguity is gone.
- **`Frame` / `FrameMainDrawer` / `FrameTriple` removed** — they were `JoinVertical`
  with named slots. Apps use `Pancake`, `TopbarPancake`, `Focus`, or `JoinVertical` directly.
- **`theme/styles/System` deleted** — `styles.Row`, `RowClip`, `ClipANSI` remain as
  pure package-level functions. No more `styles.New(t)` wiring in bricks.
- **CLI TUI and diff viewer unlocked** — bricks work without any global initialization.
  Use `theme.Preset("name")` + `WithTheme(t)` and render in any context.

---

## Backlog

### CLI (`cmd/bento`)

- [x] `bento add` — copy component files from embedded registry
- [x] `bento init` — generate a runnable starter project
- [x] `bento list` — show available components with descriptions
- [ ] `bento upgrade <component>` — diff local copy vs registry version
- [ ] `bento add --force` — overwrite existing copied component

### New bentos

- [x] `home-screen` — starter-style entry screen
- [x] `dashboard` — card/table composition
- [x] `app-shell` — rail + workspace + command palette
- [x] `detail-view` — list + detail split
- [ ] `form` — labeled inputs + validation flow
- [ ] `log-viewer` — filter + scrollable output
- [ ] `settings` — left nav + content pane

### New bricks

- Current catalog covers all common needs. New bricks only when the same gap
  appears in at least 2 bentos and can't be composed from existing bricks.
- No `spinner` brick — use `charm.land/bubbles/v2/spinner` directly.

### Rooms

- [ ] `rooms.Grid` — fixed-column grid helper for dense dashboards
- [ ] Scrollable region helper — independent body scroll within named layouts

### Testing

- [ ] Snapshot tests for every brick's rendered output
- [ ] Smoke tests for `bento add` and `bento init` CLI paths

### Wrap + scaffold

- [ ] `bento wrap --manifest-only` — parse interface, emit deterministic manifest JSON
- [ ] `bento wrap --scaffold` — generate owned Go scaffold from manifest
- [ ] `bento wrap --enhance` — optional LLM pass after deterministic scaffold
- [ ] `llms.txt` — ship model context for scaffold tooling

---

## Non-goals

- **Mobile / small screens** — assumes a reasonably large terminal
- **Mouse support** — no plans unless a specific component clearly needs it
- **Accessibility** — depends on the terminal emulator, not the TUI library
- **Web renderer** — terminal output only
- **Built-in router** — bentos own their own state machines
- **Data-fetching** — bring your own

---

Last updated: 2026-03-18
