# Layout Hierarchy Research: Crush & OpenCode

**Date:** 2026-03-01  
**Status:** Findings applied — fix implemented in `core/layout` and `ui/containers/panel`

---

## Problem

The BentoTUI starter-app displayed blurry, undefined panel boundaries. Panels appeared to bleed into one another with no clear region ownership. The root cause was traced to `layout.Split.View()` not forcing a background fill across each allocated region before compositing child content.

---

## How Crush Does It

**Source:** `internal/ui/model/ui.go` (charm-community/crush)

Crush uses `ultraviolet.ScreenBuffer` as the root draw loop:

```go
screen.Clear(scr)  // explicit full-screen clear every frame
```

Layout is computed as typed `image.Rectangle` regions:

```go
type uiLayout struct {
    sidebar image.Rectangle
    main    image.Rectangle
    editor  image.Rectangle
}
```

Each component receives pixel-exact bounds via `image.Rectangle`. Mouse hit-testing uses:

```go
image.Pt(msg.X, msg.Y).In(m.layout.main)
```

Compositing is done via `lipgloss.Canvas` + `lipgloss.Layer` with explicit X/Y/Z:

```go
lipgloss.NewLayer(view).X(x).Y(y).Z(z)
```

**Key insight:** Crush clears the entire screen buffer before every frame. Each component's `View()` is placed into a layer at an absolute pixel offset. Because the screen is pre-cleared, any region not painted by a component remains blank — there is no bleed.  
Crush does **not** rely on each component painting its own background; it relies on the root clear + absolute positioning.

Compact mode breakpoints:
```go
compactModeWidthBreakpoint  = 120
compactModeHeightBreakpoint = 30
```

Dialog is drawn last over the full canvas via `dialog.Overlay`.

---

## How OpenCode Does It

**Source:** `internal/tui/layout/container.go`, `split.go`, `overlay.go` (opencode-ai/opencode)

OpenCode is a string-first renderer (Bubble Tea v1 style). Its key pattern is in `container.go`:

```go
lipgloss.NewStyle().
    Width(w).
    Height(h).
    Background(t.Background()).
    Render(content)
```

**Every container forces its full region background** before rendering content. This is the canonical fix pattern.

`split.go` (`SplitPaneLayout`) uses float64 ratio (not integer weight):

```go
leftWidth  = int(float64(width) * s.ratio)
rightWidth = width - leftWidth
```

`View()` wraps the final composed string:

```go
lipgloss.NewStyle().
    Width(s.width).
    Height(s.height).
    Background(t.Background()).
    Render(finalView)
```

This guarantees the full `w × h` region is painted regardless of what the child renders.

`overlay.go` is a string-based `PlaceOverlay(x, y, fg, bg)` that merges strings line-by-line preserving ANSI, with optional shadow rendering.

---

## Root Cause in BentoTUI

`layout.Split.View()` was compositing child views via `lipgloss.Canvas` + `lipgloss.Layer`:

```go
layer := lipgloss.NewLayer(core.ViewLayer(item.child.View())).X(x).Y(0).Z(i)
```

The problem: `lipgloss.Canvas` composites layers over **whatever the terminal already shows**. Without an explicit background fill for each allocated region, the canvas background bleeds through — child panels paint only their content, not their surrounding empty space.

Similarly, `panel.Model.View()` built rows manually but wrapped them in a `PanelFrame` style that did not guarantee a solid `Width × Height` background fill.

---

## The Fix

### `layout.Split` — OpenCode container pattern

Before compositing child layers onto the canvas, each allocated region is filled with the canvas background color using a solid background rectangle:

```go
// Fill allocated region with canvas background before placing child content
bgFill := strings.Repeat(" ", w) // repeated for each of h rows
bgLayer := lipgloss.NewLayer(
    lipgloss.NewStyle().
        Width(w).Height(h).
        Background(lipgloss.Color(canvasBG)).
        Render(bgFill),
).X(x).Y(0).Z(0)
```

This mirrors OpenCode's `Width + Height + Background` wrapper pattern but applied per-region in the canvas compositor.

### `panel.Model` — explicit background wrap

`panel.Model.View()` now wraps its final rendered string in:

```go
lipgloss.NewStyle().
    Width(outerWidth).
    Height(outerHeight).
    Background(lipgloss.Color(bg))
```

This guarantees the panel paints every cell in its allocated `w × h` region, even if content is shorter than the allocated height.

---

## Z-order Reminder

Per `project-docs/design/rendering-system-design.md`:

```
shell-bg(0) → body(1) → header(2) → footer(3) → scrim(4) → dialog(5)
```

Within `Split`, child panel background fills are Z(i*2) and child content is Z(i*2+1) to preserve ordering. The shell background remains the lowest layer.

---

## References

- Crush source: `https://github.com/charm-community/crush/blob/main/internal/ui/model/ui.go`
- OpenCode split: `https://github.com/opencode-ai/opencode/blob/main/internal/tui/layout/split.go`
- OpenCode container: `https://github.com/opencode-ai/opencode/blob/main/internal/tui/layout/container.go`
- OpenCode overlay: `https://github.com/opencode-ai/opencode/blob/main/internal/tui/layout/overlay.go`
- BentoTUI rendering ADR: `project-docs/design/rendering-system-design.md`
