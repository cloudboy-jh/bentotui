# Bentos

`registry/bentos/` contains full runnable apps.

In BentoTUI terms, a **bento** is not a widget and not a layout helper. A bento
is the full app composition layer:

- state machine
- focus ownership
- keymap and interaction flow
- draw/layer order
- scenario orchestration and diagnostics

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
- `dashboard` — dense panel/table composition
- `app-shell` — validation bento for framework QA

## app-shell role

`registry/bentos/app-shell` is the canonical framework validation bento.

It exists to pressure-test rooms and bricks under repeatable scenarios:

- layout
- hierarchy
- footer
- list
- overlay
- stress

Use `app-shell` to reproduce visual regressions with a deterministic tuple:

`scenario + viewport + theme + focus + snapshot`
