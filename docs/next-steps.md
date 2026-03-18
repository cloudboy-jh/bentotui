# BentoTUI — Next Steps

## Current focus

Lock BentoTUI as the fastest way to ship full Go TUIs:

- rooms-first page composition
- official brick components with polished defaults
- bento templates users can remix in a day

---

## Immediate priorities

### 1 — Rooms as first-class page API

- Ship and document high-level room contracts (`AppShell`, `SidebarDetail`, `Dashboard`, `DiffWorkspace`)
- Add room cookbook examples for common app page shapes
- Keep lower-level split/layout APIs as advanced escape hatches

### 2 — Flagship brick polish

- Finalize `list` and `table` as reference-grade bricks
- Add richer snapshot/state tests (focus/blur/theme/resize variants)
- Publish recommended presets and usage patterns

### 3 — Template-grade bentos

- Promote `home-screen`, `app-shell`, and `detail-view` as remixable templates
- Add extension-point docs for each bento
- Add `diff-workspace` template using mock diff DTOs

### 4 — One-day OpenCode-style path

- Document a concrete flow: pick bento, add pages, wire data, ship
- Ensure the path never requires raw bubbles imports in app composition code

### 5 — CLI remains secondary

- Keep `bento add` and `bento init` healthy
- Treat CLI as convenience, not the core product surface

---

## Non-goals (still true)

- No web renderer or browser output
- No mouse-first interaction model
- No built-in app router framework
- No data-fetching abstraction layer
