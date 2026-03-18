# BentoTUI Documentation

v0.4.0

---

- [architecture/architecture.md](./architecture/architecture.md) — Layer diagram, rendering contract, theme model, component rules
- [architecture/bricks.md](./architecture/bricks.md) — Every brick's API with usage examples
- [architecture/rooms.md](./architecture/rooms.md) — Room layout API and render flow
- [architecture/bentos.md](./architecture/bentos.md) — Full app composition contract
- [theme-engine.md](./theme-engine.md) — Theme interface, presets, custom themes, live switching
- [coloring-rules.md](./coloring-rules.md) — Rules for correct color usage in bricks
- [next-steps.md](./next-steps.md) — Immediate priorities
- [roadmap.md](./roadmap.md) — Backlog, non-goals

---

## The short version

BentoTUI is the **shadcn of TUIs** — copy components into your project, own them,
modify them freely. No framework to fight.

**Three stable imports:**

```go
"github.com/cloudboy-jh/bentotui/theme"         // Theme interface + presets
"github.com/cloudboy-jh/bentotui/theme/styles"  // Row/RowClip/ClipANSI utilities
"github.com/cloudboy-jh/bentotui/registry/rooms" // Layout geometry
```

**Everything else is copy-and-own** (`bento add card`, `bento add bar`, etc.).

---

## Core concepts

**Bricks** — UI components. Accept `WithTheme(t)` at construction, `SetTheme(t)` for
live updates. Fall back to `theme.CurrentTheme()` if no theme was provided.

**Rooms** — pure geometry functions. Zero color. Zero theme. Take `(width, height, cells...)`
and return a rendered string.

**Bentos** — full apps. Hold `m.theme theme.Theme` as app state. Call `SetTheme` on bricks
when `theme.ThemeChangedMsg` arrives. Own final frame composition via `surface`.

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
