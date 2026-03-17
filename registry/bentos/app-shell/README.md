# app-shell

Reference UX/UI sandbox bento for BentoTUI.

This is the first full bento used to validate how rooms and bricks compose into
an intentional interface. It focuses on elevated-card usage with list, table,
modal, footer, and theme-audit patterns.

## Run

```bash
go run ./registry/bentos/app-shell
```

## Layout contract

- left rail: scenario selector
- center rail: live scenario canvas (inside an `elevated-card` section)
- footer stack: session `elevated-card` + anchored command `bar`
- body room: `rooms.Rail(...)` (fixed rail + flexible main)

## Controls

- `up/down`: scenario
- `left/right`: viewport preset (`80x24`, `100x30`, `140x42`)
- `1-9`: jump scenario
- `t`: cycle theme
- `d`: toggle paint debug ruler
- `s`: snapshot mode
- `q`: quit

Theme behavior:

- app-shell defaults to `AvailableStableThemes()`
- falls back to `AvailableThemes()` if stable list is empty

## Scenarios

- `cards-list`
- `cards-table`
- `cards-modal`
- `cards-footer`
- `cards-theme-audit`

## Internal shape

- `state/` root model and orchestration
- `ui/` selector/footer helpers
- `scenarios/` scenario contracts and scenario implementations
