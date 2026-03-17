# Bentos

`registry/bentos/` contains full runnable apps.

In BentoTUI terms, a **bento** is not a widget and not a layout helper. A bento
is the full app composition layer:

- state machine
- focus ownership
- keymap and interaction flow
- draw/layer order
- product-facing UX composition

## Contract

- Bentos compose with `registry/rooms`.
- Bentos render with `registry/bricks`.
- Bentos read `theme.CurrentTheme()` at render time.
- Bentos own final frame composition (`surface.Fill` + `surface.Draw`).

## Layering model

Recommended frame stack:

1. canvas fill (z0)
2. body frame card (z1)
3. anchored footer row (z2)
4. overlays/dialogs last (z3)

## Current shipped bentos

- `home-screen` — starter-style entry screen
- `dashboard` — dense elevated-card/table composition
- `app-shell` — single-screen composition bento (rail + table + list + progress + palette)

## app-shell role

`registry/bentos/app-shell` is the canonical UX proving ground.

It demonstrates the Bento contract directly:

- rooms provide grammar (`RailFooterStack`)
- bricks provide UI primitives (table/list/progress/dialog/footer)
- theme switching is global (footer + command palette)

It is intentionally a minimal app surface, not a scenario harness.
