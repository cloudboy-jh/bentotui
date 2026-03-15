# app-shell

Framework validation bento for BentoTUI.

This is not a product mock. It is a scenario runner used to pressure-test
layout layering, hierarchy, clipping, footer behavior, and theme stability.

## Run

```bash
go run ./registry/bentos/app-shell
```

## Layout contract

- left rail: scenario selector
- center rail: live scenario canvas
- right rail: diagnostics
- anchored footer: controls + state tuple + pass/warn/fail counts

Responsive body modes:

- wide (`>=120`): left + center + right
- medium (`84..119`): left + center (diagnostics collapsed)
- narrow (`<84`): vertical stack (left over center)

## Controls

- `j/k` or up/down: scenario
- `h/l` or left/right: viewport preset (`80x24`, `100x30`, `140x42`)
- `1-6`: jump scenario
- `[` / `]` (or shift+tab/tab): focus walker
- `t`: cycle theme
- `d`: toggle paint debug ruler
- `s`: snapshot mode
- `m`: show/hide keymap in diagnostics
- `r`: increment stress step
- `q`: quit

## Scenarios

- `layout`
- `hierarchy`
- `footer`
- `list`
- `overlay`
- `stress`

## Internal shape

- `state/` root model and orchestration
- `ui/` selector/diagnostics/footer helpers
- `scenarios/` scenario contracts and scenario implementations
