# BentoTUI Roadmap

## Current state (v0.6.0 direction)

### What shipped in v0.4.x

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
- **Guardrails landed** — policy tests now enforce critical layer boundaries.

---

## Product direction

Build BentoTUI as a full app product system:

- Bricks are the official component surface
- Rooms are named page-layout contracts
- Bentos are remixable app templates

---

## Backlog

### Rooms first

- [x] Added high-level room contracts (`AppShell`, `SidebarDetail`, `Dashboard`, `DiffWorkspace`)
- [ ] Add room cookbook examples for common page archetypes
- [ ] Add responsive behavior notes per named room

### Bentos (template apps)

- [x] `home-screen` — starter-style entry screen
- [x] `dashboard` — card/table composition
- [x] `app-shell` — rail + workspace + command palette
- [x] `detail-view` — list + detail split
- [ ] `diff-workspace` — file rail + diff main + footer using room contract
- [ ] `settings` — left nav + content pane
- [ ] `log-viewer` — filter + scrollable output

### Bricks

- Keep list/table as flagship polished bricks with stronger snapshots and docs.
- New bricks only when the same gap appears in at least 2 bentos and cannot be
  composed from existing bricks.
- Do not expand toward a maximal component catalog; prefer local app-owned bricks
  for one-off needs.
- No `spinner` brick — use `charm.land/bubbles/v2/spinner` directly.

### Testing

- [ ] Snapshot tests for every brick's rendered output
- [x] Guardrail tests for layering/import policy
- [ ] Smoke tests for `bento add` and `bento init` CLI paths

### CLI (`cmd/bento`) - secondary

- [x] `bento add` — copy brick files from embedded registry
- [x] `bento init <bento>` — initialize a runnable bento template
- [x] `bento list` — show available bentos, bricks, and recipes with descriptions
- [ ] `bento upgrade <component>` — diff local copy vs registry version
- [ ] `bento add --force` — overwrite existing copied component

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

Last updated: 2026-03-20
