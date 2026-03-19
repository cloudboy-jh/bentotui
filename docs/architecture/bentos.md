# Bentos

`registry/bentos/` contains template-grade runnable apps.

A **bento** is the full app composition layer you clone/remix to ship quickly:

- state machine and focus ownership
- keymap and interaction flow
- draw/layer order and dialog routing
- product-facing UX composition

Bentos are not toy demos. They are the fastest starting point for real apps.

---

## Contract

- Bentos compose layout with `registry/rooms`
- Bentos render UI with `registry/bricks`
- Bentos may use `registry/recipes` for reusable app-flow composition
- Bentos hold `m.theme theme.Theme` as app state
- Bentos call `SetTheme()` on bricks when `theme.ThemeChangedMsg` arrives
- Bentos own final frame composition (`surface.Fill` + `surface.Draw`)
- Bentos avoid direct raw `bubbles/*` imports unless a gap is explicitly allowed

---

## Theme wiring

```go
type model struct {
    theme   theme.Theme  // app owns this — not a global call in View()
    footer  *bar.Model
    content *card.Model
    ...
}

func newModel() *model {
    t := theme.CurrentTheme()
    return &model{
        theme:   t,
        footer:  bar.New(bar.FooterAnchored(), bar.WithTheme(t), ...),
        content: card.New(card.Title("Main"), card.WithTheme(t), ...),
    }
}

case theme.ThemeChangedMsg:
    m.theme = msg.Theme
    m.footer.SetTheme(m.theme)
    m.content.SetTheme(m.theme)
```

---

## Layering model

Recommended frame stack:

1. `surf.Fill(t.Background())` — canvas (z0)
2. `surf.Draw(0, 0, rooms.Focus(...))` — body + footer (z1)
3. `surf.DrawCenter(dm.View())` — dialogs last (z2)

---

## Shipped bentos

| Bento | Description |
|---|---|
| `home-screen` | Starter entry shell with theme picker and command hint |
| `dashboard` | Dense card/table composition — 2×2 grid of metric cards |
| `app-shell` | Single-screen composition bento: rail + table + list + progress + command palette |
| `detail-view` | List + detail pane split view |
| `dashboard-brick-lab` | Component showcase and layout test surface |

---

## app-shell role

`registry/bentos/app-shell` is the canonical workspace template.

It demonstrates the full Bento contract:

- rooms provide page contracts (`AppShell`) plus lower-level geometry when needed
- bricks provide primitives (card/table/list/progress/dialog/bar)
- theme switching is explicit: `ThemeChangedMsg` → `SetTheme` on bricks
- command palette is wired to the footer keybind

Use app-shell to validate a new brick/room contract before shipping it.

---

## One-day build target

The intended workflow:

1. Start from `home-screen` or `app-shell`
2. Replace fake data with your domain
3. Add pages using room patterns
4. Ship a serious app shell in one day
