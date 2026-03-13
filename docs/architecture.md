# BentoTUI Architecture

Status: Active reference — updated after frame hierarchy cleanup (2026-03-13)

## Overview

BentoTUI is a **registry of copy-and-own components**, not a framework. Run
`bento add input` and the source lands in your project. The stable shared
imports are `theme`, `styles`, and `registry/layouts`.

---

## Stack

```
bubbletint (palette registry)
     │
     ▼
theme/ (semantic token adapter + runtime store)
     │  CurrentTheme(), SetTheme(), PreviewTheme()
     │  16 built-in palettes (incl. bento-rose)
     │  Contrast validation: key layer pairs guaranteed distinct
     ▼
styles/ (token → lipgloss.Style mapping)
     │  styles.New(t).DialogFrame(), InputStyles(), Row(), etc.
     ▼
registry/components/ (copy-and-own UI pieces)
     │  input  bar  dialog  panel  list  table  text  surface
     │  Each component renders to a string via lipgloss
     ▼
registry/layouts (named layout patterns)
     │  Frame / FrameMainDrawer / FrameTriple + compatibility layouts
     │  screen grammar + geometry only, no theme or canvas paint
     ▼
surface (Ultraviolet-backed full-terminal cell buffer)
     │  surface.New(w, h) → Fill(bg) → Draw(0,0,layout) → DrawCenter(dialog) → Render()
     │  Root canvas for every full-screen frame
     ▼
Bubble Tea v2 frame (tea.NewView, AltScreen, BackgroundColor)
```

---

## Rendering Contract

**This is the most important rule in the codebase. Every component and bento
must follow it — without it, gaps and color bleed appear on the surface.**

### Why it exists

`surface.Fill(bg)` paints every terminal cell with the canvas background.
`surface.Draw(x, y, content)` overlays component strings on top using
Ultraviolet's cell buffer. When a lipgloss-rendered string contains padding
or whitespace cells with `Bg=nil`, UV's overlay inherits the already-filled
canvas color for those cells — so the component's intended background
only appears on cells that have an explicit `Bg` set.

Lipgloss **does not** propagate a container's `Background` into the `Bg` field
of individual padding cells unless `Width()` is also set. Padding cells
emitted by bare `Render(content)` carry no `Bg` — they look correct in a
full-terminal lipgloss render but break under UV cell-level compositing.

### The rule

> **Every row rendered by a component or bento must have an explicit `Bg` on
> every cell. Use `styles.Row(bg, fg, width, content)` or an equivalent
> `.Background().Width(w).Render()` chain. Never use bare `Render(content)`,
> `lipgloss.PlaceHorizontal`, or `lipgloss.Place` for rows that sit on a
> surface.**

```go
// ✅ Correct — every cell has explicit Bg
styles.Row(t.Input.BG, t.Input.FG, width, content)

// ✅ Also correct — explicit chain
lipgloss.NewStyle().
    Background(lipgloss.Color(bg)).
    Foreground(lipgloss.Color(fg)).
    Width(width).
    Render(content)

// ❌ Wrong — padding cells have Bg=nil, canvas color bleeds through
lipgloss.NewStyle().
    Background(lipgloss.Color(bg)).
    Render(content)                   // no Width = no Bg on padding cells

// ❌ Wrong — PlaceHorizontal leaves unstyled whitespace
lipgloss.PlaceHorizontal(width, lipgloss.Left, styledStr)

// ❌ Wrong — bare render, no bg
lipgloss.NewStyle().Foreground(fg).Render(content)
```

### For containers with padding/borders

Set `Width()` on the container **and** on every inner row:

```go
innerW := containerW - borderCells - paddingCells
rowStyle := lipgloss.NewStyle().Background(lipgloss.Color(bg)).Width(innerW)

inner := lipgloss.JoinVertical(lipgloss.Left,
    rowStyle.Render(line1),
    rowStyle.Render(line2),
)
lipgloss.NewStyle().
    Background(lipgloss.Color(bg)).
    Width(containerW).
    Padding(1, 2).
    Render(inner)
```

### For bentos

```go
// In View():
surf := surface.New(m.width, m.height)   // 1. allocate full-terminal buffer
surf.Fill(canvasColor)                    // 2. paint every cell with canvas bg
surf.Draw(x, y, componentStr)            // 3. overlay components (Bg-safe strings only)
surf.DrawCenter(dialogStr)               // 4. overlay dialogs
surf.Draw(0, m.height-1, statusStr)      // 5. status bar last row
return tea.NewView(surf.Render())         // 6. one render, nothing appended outside
```

---

## Theme Contract

Bubbletint provides palette slots. BentoTUI's `theme/adapter.go` maps those
slots to semantic layer tokens with guaranteed visual separation:

| Token pair | Min luminance delta |
|------------|-------------------|
| `input.bg` vs `surface.canvas` | 0.03 |
| `selection.bg` vs `surface.canvas` | 0.05 |
| `selection.bg` vs `input.bg` | 0.05 |
| `dialog.bg` vs `surface.canvas` | 0.03 |

Layer hierarchy (darkest → lightest for dark themes):

```
Surface.Canvas      ← terminal root (surface.Fill)
Surface.Panel       ← first raised layer (panels, containers)
Surface.Elevated    ← secondary depth (nested panels, sidebars)
Input.BG            ← text field background (contrasts canvas)
Dialog.BG           ← modal body (contrasts canvas)
Selection.BG        ← always brightest accent (contrasts everything)
```

Components use semantic tokens — never raw palette slots or hardcoded hex.

---

## Package Responsibilities

### `theme/`

Goroutine-safe global theme store backed by bubbletint palettes.

- `CurrentTheme() Theme` — read active theme in `View()`, never cache
- `SetTheme(name) (Theme, error)` — set + persist to disk
- `PreviewTheme(name) (Theme, error)` — live preview, no persistence
- `RegisterTheme(name, Theme) error` — add custom theme (validated)
- `AvailableThemes() []string` — sorted list, default first

### `styles/`

Maps `theme.Theme` tokens to `lipgloss.Style` constructors. All color
decisions live here — components never scatter hex literals.

Key helpers:

```go
sys := styles.New(theme.CurrentTheme())
styles.Row(bg, fg, width, content)   // ← use this for all rows
sys.DialogFrame()
sys.DialogSearchRow()
sys.DialogListRow()
sys.DialogListRowSelected()
sys.InputStyles()                    // textinput.Styles for bubbles/textinput
```

### `registry/components/surface/`

Ultraviolet-backed full-terminal cell buffer. The root canvas for every
bento and layout. Copy into your project with `bento add surface`.

```go
surf := surface.New(width, height)
surf.Fill(bg)               // paint every cell — call first
surf.Draw(x, y, str)        // overlay: drops nil pre-clear cells, inherits bg
surf.DrawCenter(str)        // centered overlay for dialogs
surf.Render()               // → ANSI string for tea.NewView
```

### `registry/layouts/`

Named layout patterns with strict cell constraints.

- `registry/layouts` does geometry only (allocation, constrain, join, overlay math)
- It intentionally does not import `theme` or `surface`
- Always composite layout output through `surface` in app `View()`

```go
screen := layouts.TopbarPancake(w, h,
    topbar,
    header,
    content,
    footer,
)

surf := surface.New(w, h)
surf.Fill(lipgloss.Color(theme.CurrentTheme().Surface.Canvas))
surf.Draw(0, 0, screen)
return tea.NewView(surf.Render())
```

---

## Component Types

Components in `registry/components/` have three distinct roles. Knowing which
type you are building determines what rendering rules apply.

### Atomic

Examples: `input`, `badge`, `kbd`, `text`, `wordmark`

- Returns a **styled string** — no awareness of surface position or terminal width
- Does **not** need `Width()` — the caller sizes and places it
- Does **not** own a background region — bg is provided by the container it sits in
- May use foreground color and bold/italic freely

```go
// Atomic — returns a string, caller draws it
badge := lipgloss.NewStyle().
    Foreground(lipgloss.Color(t.Text.Accent)).
    Render("v0.2.0")
surf.Draw(x, y, badge)
```

### Container

Examples: `panel`, `dialog`, `bar`

- Owns a **width × height region** on the surface
- Must set `Width()` on every row so every cell has explicit `Bg`
- Uses `styles.Row(bg, fg, width, content)` for all body rows
- Never uses bare `Render(content)` or `PlaceHorizontal` for rows

```go
// Container — explicit width on every row
row := styles.Row(t.Dialog.BG, t.Text.Primary, width, content)
```

### Surface

`registry/components/surface/surface.go`

- The **full-terminal root canvas** — one per frame, sized to terminal dimensions
- Not a UI component — it is the compositor that everything else draws onto
- `Fill(bg)` first, then `Draw(x, y, ...)` for all atomics and containers
- One `Render()` call at the end — nothing appended outside it

```go
surf := surface.New(m.width, m.height)
surf.Fill(canvasColor)          // root layer
surf.Draw(x, y, containerStr)  // containers
surf.Draw(x, y, atomicStr)     // atomics inside containers
surf.DrawCenter(dialogStr)      // overlays
surf.Draw(0, m.height-1, bar)  // status bar last
return tea.NewView(surf.Render())
```

---

## Component Rules

1. **Read theme at render time** — `t := theme.CurrentTheme()` in `View()`, never stored
2. **Containers: every row has explicit Bg** — use `styles.Row()` or `.Background().Width().Render()`
3. **Atomics: no Width() required** — caller handles placement and sizing
4. **No imports between registry components** — each is standalone
5. **No raw palette slots** — use semantic tokens from `theme.Theme`
6. **No `lipgloss.PlaceHorizontal` for surface-drawn rows** — use `Width()` instead

---

## Why Not a Framework?

Frameworks make the easy cases easy and the hard cases impossible. Every
non-trivial TUI eventually needs to reach inside and change something the
framework didn't anticipate. BentoTUI copies source so that modification is
always zero-friction — there is no "extend" API to learn, no lifecycle hooks,
no middleware pattern. You just edit the file.
