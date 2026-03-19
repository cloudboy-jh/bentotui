# Usage Guide: Build with Bento Surface APIs

Use BentoTUI as a product system:

- `bricks` are the official UI components you copy and own
- `rooms` are named layout patterns you choose per page
- `bentos` are full app templates you can remix quickly

If Bento has a primitive for the job, use Bento first.

Do not treat Bento as a promise to ship every possible primitive as an official brick.
For uncovered gaps, create local app-owned bricks and keep shipping.

---

## Default app flow

1. Start from a bento template (`registry/bentos/*`).
2. Copy the bricks you need with `bento add <brick...>`.
3. Pick a room per page (`rooms.AppShell`, `rooms.SidebarDetail`, `rooms.DiffWorkspace`, ...).
4. Keep one root app model that routes pages.
5. Keep theme in model state and propagate with `SetTheme`.

---

## Import policy

- `registry/bricks/*`: copy-and-own components.
- `registry/rooms`: stable layout grammar import.
- `theme` and `theme/styles`: stable theme contracts.

In app composition layers (`registry/bentos/*`, starter apps, generated starter):

- Do not import raw `bubbles/*` directly when a Bento brick exists.
- Keep Charm internals behind brick wrappers.
- Exception: `spinner` is allowed until a Bento spinner strategy exists.

---

## When custom bricks are acceptable

Custom bricks are expected when there is no official brick for your use case.

Use this decision tree:

1. If an official brick exists, use it.
2. If not, compose from existing bricks + rooms.
3. If that still does not fit, create a local app-owned brick.
4. Propose upstream only when the same need repeats across multiple bentos.

This keeps app delivery unblocked while preserving a focused official surface.

---

## Theme ownership contract

Own theme in model state:

```go
type model struct {
    theme theme.Theme
}
```

Propagate explicitly:

```go
case theme.ThemeChangedMsg:
    m.theme = msg.Theme
    m.list.SetTheme(m.theme)
    m.table.SetTheme(m.theme)
    m.footer.SetTheme(m.theme)
```

In `View()`, render from `m.theme`, not from global lookups.

---

## Rooms contract

`registry/rooms` is geometry and composition only:

- no `theme` imports
- no brick imports
- no direct `bubbles/*` dependencies

Choose a room in each page file and compose there.
Go imports are package-level, so room choice is explicit by function call.

```go
import "github.com/cloudboy-jh/bentotui/registry/rooms"

screen := rooms.DiffWorkspace(w, h, 28, header, fileRail, mainDiff, footer)
```

---

## Bricks contract

Bricks are official Bento components:

- standalone (no cross-brick imports)
- ergonomic defaults plus tweakable APIs
- `WithTheme(t)` at creation and `SetTheme(t)` for live updates

Users should be able to ship full apps without thinking about raw bubbles internals.

---

## Enforcement

`internal/policy/guardrails_test.go` enforces:

- rooms import boundaries
- no cross-brick imports
- no raw `bubbles/*` usage in bentos (except spinner)
- no raw `bubbles/*` usage in starter/scaffold composition code
- no `theme.CurrentTheme()` inside bento `View()` methods

Run with:

```bash
go test ./...
```
