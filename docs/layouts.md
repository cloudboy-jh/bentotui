# BentoTUI Layouts

`registry/layouts` provides named layout functions. Each function takes
`(width, height, ...cells)` and returns a rendered `string`.

Contract:

- Layouts define screen grammar + geometry (allocation + constrain)
- Layouts are theme-agnostic
- Use `surface` as the final compositor in app `View()`
- Prefer `Focus` for body + anchored-footer screens; use `Frame` when you explicitly need top/subheader rows.
- Bar roles for `Frame` rows: top (`RoleTopBar`), subheader (`RoleSubBar`), footer (`RoleFooterBar` / `FooterAnchored`).

## API Basics

```go
import "github.com/cloudboy-jh/bentotui/registry/layouts"

screen := layouts.Focus(w, h, content, footer)
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

Frame grammar:

- `Frame(top, subheader, body, subfooter)`
- `FrameMainDrawer(drawerW, top, subheader, main, drawer, subfooter)`
- `FrameTriple(navW, listW, top, subheader, nav, list, detail, subfooter)`

Compatibility layouts:

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
screen := layouts.Focus(w, h, content, footer)

surf := surface.New(w, h)
surf.Fill(lipgloss.Color(theme.CurrentTheme().Surface.Canvas))
surf.Draw(0, 0, screen)

if dialogs.IsOpen() {
    surf.DrawCenter(viewString(dialogs.View()))
}

return tea.NewView(surf.Render())
```
