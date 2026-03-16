# app-shell

Framework validation bento for BentoTUI.

This is the reference full bento used to pressure-test how bricks compose together.
It focuses on elevated-card usage with list, table, modal, and footer patterns.

## Run

```bash
go run ./registry/bentos/app-shell
```

## Layout contract

- left rail: scenario selector
- center rail: live scenario canvas (inside an `elevated-card` section)
- footer stack: session `elevated-card` + anchored command `bar`

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

- `cards-list`
- `cards-table`
- `cards-modal`
- `cards-footer`

## Internal shape

- `state/` root model and orchestration
- `ui/` selector/footer helpers
- `scenarios/` scenario contracts and scenario implementations
