# app-shell

Framework validation bento for BentoTUI.

This is not a product mock. It is a minimal scenario runner used to pressure-test
layout layering, clipping, footer behavior, and theme stability.

## Run

```bash
go run ./registry/bentos/app-shell
```

## Layout contract

- left rail: scenario selector
- center rail: live scenario canvas
- anchored footer: controls + state tuple + pass/warn/fail counts
- diagnostics: compact inline summary under the canvas

Responsive body modes:

- wide (`>=84`): left + center
- narrow (`<84`): vertical stack (left over center)

## Controls

- `up/down`: scenario
- `left/right`: viewport preset (`80x24`, `100x30`, `140x42`)
- `1-3`: jump scenario
- `t`: cycle theme
- `d`: toggle paint debug ruler
- `s`: snapshot mode
- `q`: quit

## Scenarios

- `layout`
- `footer`
- `stress`

## Internal shape

- `state/` root model and orchestration
- `ui/` selector/footer helpers
- `scenarios/` scenario contracts and scenario implementations
