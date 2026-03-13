# BentoTUI Layouts

`registry/layouts` provides 15 named layout functions. Each function takes
`(width, height, ...cells)` and returns a rendered `string`.

Contract:

- Layouts do geometry only (allocation + constrain)
- Layouts are theme-agnostic
- Use `surface` as the final compositor in app `View()`

## API Basics

```go
import "github.com/cloudboy-jh/bentotui/registry/layouts"

screen := layouts.Pancake(w, h, header, content, footer)
```

Cells must satisfy:

```go
type Sizable interface {
    SetSize(width, height int)
    View() tea.View
}
```

Helpers:

- `layouts.Static("...")`
- `layouts.RenderFunc(func(width, height int) string { ... })`

## Named Layouts

- `Focus(content, footer)`
- `Pancake(header, content, footer)`
- `TopbarPancake(topbar, header, content, footer)`
- `Sidebar(sideWidth, sidebar, main)`
- `HolyGrail(sideWidth, header, sidebar, main, footer)`
- `HSplit(left, right)`
- `VSplit(top, bottom)`
- `HSplitFooter(left, right, footer)`
- `TripleCol(navW, listW, nav, list, detail)`
- `Dashboard2x2(tl, tr, bl, br)`
- `Dashboard2x2Footer(tl, tr, bl, br, footer)`
- `DrawerRight(drawerW, main, drawer)`
- `DrawerChrome(drawerW, header, main, drawer, footer)`
- `Modal(modalW, modalH, background, modal)`
- `BigTopStrip(stripH, primary, strip)`

## Recommended Render Flow

```go
screen := layouts.TopbarPancake(w, h, topbar, header, content, footer)

surf := surface.New(w, h)
surf.Fill(lipgloss.Color(theme.CurrentTheme().Surface.Canvas))
surf.Draw(0, 0, screen)

if dialogs.IsOpen() {
    surf.DrawCenter(viewString(dialogs.View()))
}

return tea.NewView(surf.Render())
```
