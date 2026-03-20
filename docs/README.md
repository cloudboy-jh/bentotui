# BentoTUI Documentation

v0.5.4

---

- [usage-guide.md](./usage-guide.md) — Build apps with bricks + recipes + rooms + bentos

### Architecture (`docs/architecture/`)

- [architecture/architecture.md](./architecture/architecture.md) — Layer diagram, rendering contract, theme model, component rules
- [architecture/bricks.md](./architecture/bricks.md) — Official brick APIs and conventions
- [architecture/recipes.md](./architecture/recipes.md) — Composed recipe APIs and flow patterns
- [architecture/rooms.md](./architecture/rooms.md) — Named room layouts and page composition patterns
- [architecture/bentos.md](./architecture/bentos.md) — Template app contract and extension model
- [architecture/astro-content-update.md](./architecture/astro-content-update.md) — Current-state Astro/docs content update source

### Design (`docs/design/`)

- [design/theme-engine.md](./design/theme-engine.md) — Theme interface, presets, custom themes, live switching
- [design/coloring-rules.md](./design/coloring-rules.md) — Rules for correct color usage in bricks

### Development (`docs/development/`)

- [development/product-direction.md](./development/product-direction.md) — Scope and contract to reduce framework churn
- [development/next-steps.md](./development/next-steps.md) — Immediate priorities
- [development/roadmap.md](./development/roadmap.md) — Backlog, non-goals

---

## The short version

BentoTUI is the **best way to build full Go TUIs quickly**:

- copy-and-own bricks
- copy-and-own recipes for composed flows
- named rooms for layout
- bento templates for full-screen apps

**Three stable imports:**

```go
"github.com/cloudboy-jh/bentotui/theme"         // Theme interface + presets
"github.com/cloudboy-jh/bentotui/theme/styles"  // Row/RowClip/ClipANSI utilities
"github.com/cloudboy-jh/bentotui/registry/rooms" // Layout geometry
```

**Everything else is copy-and-own** (`bento add card`, `bento add recipe filter-bar`, etc.).

---

## Core concepts

**Bricks** — UI components. Accept `WithTheme(t)` at construction, `SetTheme(t)` for
live updates. Fall back to `theme.CurrentTheme()` if no theme was provided.

**Recipes** — composed app-facing patterns built from bricks. Copy with
`bento add recipe <name>` and adapt freely.

**Rooms** — named page layouts. Zero color, zero theme. Choose one per page and
compose your app shape there.

**Bentos** — template-grade full apps. Start here, replace data/domain logic,
ship quickly.

**Surface** — Ultraviolet-backed cell buffer. Every bento uses:
`surface.New` → `Fill(bg)` → `Draw(layout)` → `DrawCenter(dialog)` → `Render()`.

---

## Quick start

```go
t := theme.CurrentTheme()

footer := bar.New(bar.FooterAnchored(), bar.WithTheme(t),
    bar.Left("my-app"),
    bar.Cards(bar.Card{Command: "ctrl+c", Label: "quit", Enabled: true}),
)

content := card.New(card.Title("Main"), card.WithTheme(t), card.Content(myList))

// In View():
surf := surface.New(w, h)
surf.Fill(t.Background())
surf.Draw(0, 0, rooms.Focus(w, h, content, footer))
v := tea.NewView(surf.Render())
v.AltScreen = true
v.BackgroundColor = t.Background()
return v
```
