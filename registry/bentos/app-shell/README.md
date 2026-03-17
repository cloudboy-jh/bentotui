# app-shell

Reference single-screen app bento for BentoTUI.

This bento shows how rooms + bricks can compose into a complete app shell
without rolling a custom TUI framework.

## Run

```bash
go run ./registry/bentos/app-shell
```

## Layout contract

- left rail: app sections + lightweight status
- main canvas: services table over queue/progress cards
- bottom row: single anchored command bar
- body room: `rooms.Rail(...)` with `rooms.RailFooterStack(...)`

## Controls

- `up/down`: switch active section
- `left/right`: move queue cursor
- `enter`: pulse progress value
- `t`: cycle theme
- `c`: toggle compact table mode
- `ctrl+k`: open command palette
- `1-9`: jump section
- `q`: quit

## Command Palette

- powered by `registry/bricks/dialog/command_palette.go`
- includes full theme list from `theme.AvailableThemes()`
- supports section jumps and view toggles

## Internal shape

- `state/` root model, workspace deck, and palette actions
- `ui/` rail/footer copy helpers
