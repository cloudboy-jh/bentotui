# BentoTUI Documentation

v0.5.0 (in progress)

---

- [usage-guide.md](./usage-guide.md) — Build apps with bricks + rooms + bentos
- [architecture/rooms.md](./architecture/rooms.md) — Named room layouts and page composition patterns
- [architecture/bricks.md](./architecture/bricks.md) — Official brick APIs and conventions
- [architecture/bentos.md](./architecture/bentos.md) — Template app contract and extension model
- [architecture/architecture.md](./architecture/architecture.md) — Layer diagram, rendering contract, theme model, component rules
- [theme-engine.md](./theme-engine.md) — Theme interface, presets, custom themes, live switching
- [coloring-rules.md](./coloring-rules.md) — Rules for correct color usage in bricks
- [astro-content.md](./astro-content.md) — Marketing/site copy source for docs and landing pages
- [product-direction.md](./product-direction.md) — Scope and contract to reduce framework churn
- [next-steps.md](./next-steps.md) — Immediate priorities
- [roadmap.md](./roadmap.md) — Backlog, non-goals

---

## The short version

BentoTUI is the **best way to build full Go TUIs quickly**:

- copy-and-own bricks
- named rooms for layout
- bento templates for full-screen apps

**Three stable imports:**

```go
"github.com/cloudboy-jh/bentotui/theme"         // Theme interface + presets
"github.com/cloudboy-jh/bentotui/theme/styles"  // Row/RowClip/ClipANSI utilities
"github.com/cloudboy-jh/bentotui/registry/rooms" // Layout geometry
```

**Everything else is copy-and-own** (`bento add card`, `bento add list`, etc.).

---

## Core concepts

**Bricks** — UI components. Accept `WithTheme(t)` at construction, `SetTheme(t)` for
live updates. Fall back to `theme.CurrentTheme()` if no theme was provided.

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
