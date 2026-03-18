# BentoTUI Rooms

`registry/rooms` provides named layout functions. Each function takes
`(width, height, ...cells)` and returns a rendered `string`.

Rules:

- Rooms define geometry only — allocation, constrain, join
- Rooms are completely theme-agnostic and do not import `theme` or `surface`
- Composite room output through `surface` in your app `View()`
- `Frame` and its variants have been removed — use `Pancake`, `TopbarPancake`,
  or `Focus` instead

---

## Sizable interface

All room cells must satisfy:

```go
type Sizable interface {
    SetSize(width, height int)
    View() tea.View
}
```

Helpers for ad-hoc content:

```go
rooms.Static("static string")
rooms.RenderFunc(func(width, height int) string { ... })
```

---

## Named Rooms

### Focus layouts

- `Focus(w, h, content, footer)` — body + anchored footer, no top rows
- `Pancake(w, h, header, content, footer)` — header + body + footer
- `TopbarPancake(w, h, topbar, header, content, footer)` — topbar + header + body + footer

### Rail layouts

- `Rail(w, h, railWidth, rail, main)`
- `RailFooterStack(w, h, railWidth, footerCardRows, rail, main, footerCard, footer)`

### Split layouts

- `HSplit(w, h, left, right)` — equal horizontal halves
- `VSplit(w, h, top, bottom)` — equal vertical halves
- `HSplitFooter(w, h, left, right, footer)`

### Multi-pane layouts

- `HolyGrail(w, h, railWidth, header, rail, main, footer)`
- `TripleCol(w, h, navW, listW, nav, list, detail)`
- `DrawerRight(w, h, drawerWidth, main, drawer)`
- `DrawerChrome(w, h, drawerWidth, header, main, drawer, footer)`

### Dashboard layouts

- `Dashboard2x2(w, h, topLeft, topRight, bottomLeft, bottomRight)`
- `Dashboard2x2Footer(w, h, topLeft, topRight, bottomLeft, bottomRight, footer)`

### Overlay layout

- `Modal(w, h, modalWidth, modalHeight, background, modal)`

### Strip layout

- `BigTopStrip(w, h, stripHeight, primary, strip)`

---

## Separation options

```go
rooms.HSplit(w, h, left, right, rooms.WithGutter(1))
rooms.DrawerRight(w, h, 28, main, drawer, rooms.WithGutter(1), rooms.WithDivider("subtle"))
```

`rooms.WithGutter(n)` adds an explicit spacer column/row between panes.
`rooms.WithDivider("subtle")` fills the gutter with `.` characters;
`rooms.WithDivider("normal")` fills it with `|` characters.

Rooms are theme-agnostic — dividers are plain ASCII, no ANSI color applied.
If you want a styled divider, pass a `Static(styledString)` or use a `separator`
brick as a gutter cell.

---

## Render flow

```go
func (m *model) View() tea.View {
    t := m.theme

    screen := rooms.Focus(m.w, m.h, m.body, m.footer)

    surf := surface.New(m.w, m.h)
    surf.Fill(t.Background())
    surf.Draw(0, 0, screen)

    if m.dialogs.IsOpen() {
        surf.DrawCenter(viewString(m.dialogs.View()))
    }

    v := tea.NewView(surf.Render())
    v.AltScreen = true
    v.BackgroundColor = t.Background()
    return v
}
```
