# app-shell

Template-grade workspace bento for BentoTUI.

This bento shows how rooms + bricks compose into a full app shell you can clone
and adapt in a day.

## Run

```bash
go run ./registry/bentos/app-shell
```

## Layout contract

- main canvas: services table over queue/progress cards
- bottom row: single anchored command bar
- room contract: `rooms.AppShell(...)`

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
- `ui/` footer copy helpers

## Remix flow

1. Keep the room contract and footer controls.
2. Replace list/table/progress data with your domain data.
3. Add pages behind the same shell as your app grows.
