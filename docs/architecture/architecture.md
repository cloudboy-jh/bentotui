# BentoTUI Architecture

v0.5.0 direction — productized bricks + rooms + bentos

## Overview

BentoTUI is a **product system for full Go TUIs**:

- `bricks`: official copy-and-own UI components
- `rooms`: named page layout patterns
- `bentos`: full app templates for fast shipping

Run `bento add card` and the source lands in your project. You own it.

The three stable shared imports are `theme`, `theme/styles`, and `registry/rooms`.
Everything under `registry/bricks/` is copy-and-own.

---

## Stack

```
theme/ (interface + 16 preset structs)
     │  Theme interface, BaseTheme, Preset("name"), Names()
     │  Manager: CurrentTheme(), SetTheme(), PreviewTheme(), RegisterTheme()
     │  Zero global state required — bricks use whatever Theme you pass them
     ▼
theme/styles/ (pure rendering utilities)
     │  Row(bg, fg, width, content)
     │  RowClip(bg, fg, width, content)
     │  ClipANSI(content, width)
     ▼
registry/bricks/ (copy-and-own UI pieces)
     │  card  bar  dialog  input  list  table  surface  + more
     │  Each brick accepts WithTheme(t) and SetTheme(t)
     │  Falls back to theme.CurrentTheme() if no theme provided
     ▼
registry/rooms/ (named layout patterns)
     │  Geometry only — no theme, no color, no canvas paint
     │  Focus / Pancake / Rail / HolyGrail / HSplit / VSplit / ...
     ▼
surface (Ultraviolet-backed full-terminal cell buffer)
     │  surface.New(w, h) → Fill(bg) → Draw(0,0,layout) → DrawCenter(dialog) → Render()
     ▼
Bubble Tea v2 frame (tea.NewView, AltScreen, BackgroundColor)
```

Registry shape:

```
registry/bentos/  → template apps (state machine + orchestration)
registry/rooms/   → named page layouts + geometry grammar
registry/bricks/  → official UI components (copy-and-own)
```

---

## Rendering Contract

**Every component and bento must follow this. Without it, gaps and color bleed
appear on the Ultraviolet surface.**

### Why it exists

`surface.Fill(bg)` paints every terminal cell with the canvas background.
`surface.Draw(x, y, content)` overlays component strings on top using
Ultraviolet's cell buffer. When a lipgloss-rendered string contains padding
or whitespace cells with `Bg=nil`, UV inherits the already-filled canvas color
for those cells — so the component's intended background only appears on cells
that have an explicit `Bg` set.

Lipgloss does **not** propagate a container's `Background` into padding cells
unless `Width()` is also set. Bare `Render(content)` carries no `Bg` on padding
cells — looks correct in a plain lipgloss render but breaks under UV compositing.

### The rule

> **Every row rendered by a component or bento must have an explicit `Bg` on
> every cell. Use `styles.Row(bg, fg, width, content)` or an equivalent
> `.Background(x).Width(w).Render()` chain. Never use bare `Render(content)`,
> `lipgloss.PlaceHorizontal`, or `lipgloss.Place` for rows that sit on a surface.**

```go
// Correct — every cell has explicit Bg
styles.Row(t.InputBG(), t.InputFG(), width, content)

// Also correct — explicit chain
lipgloss.NewStyle().
    Background(t.BackgroundPanel()).
    Foreground(t.Text()).
    Width(width).
    Render(content)

// Wrong — padding cells have Bg=nil
lipgloss.NewStyle().
    Background(t.BackgroundPanel()).
    Render(content)  // no Width = no Bg on padding cells

// Wrong — leaves unstyled whitespace
lipgloss.PlaceHorizontal(width, lipgloss.Left, styledStr)
```

### For bentos

```go
// In View():
surf := surface.New(m.width, m.height)   // 1. allocate buffer
surf.Fill(t.Background())                 // 2. paint every cell with canvas bg
surf.Draw(x, y, componentStr)            // 3. overlay components
surf.DrawCenter(dialogStr)               // 4. overlay dialogs
return tea.NewView(surf.Render())         // 5. one render
```

---

## Theme Model

`Theme` is a Go **interface**. Every preset implements it. Bricks call methods
on it — `t.Background()`, `t.SelectionBG()`, `t.Text()`, etc.

```go
// Brick reads the theme it was given, falls back to global if none
func (m *Model) activeTheme() theme.Theme {
    if m.theme != nil {
        return m.theme
    }
    return theme.CurrentTheme()
}
```

No dot-accessor token structs. No contrast validation. No builtinThemes map.
The 16 presets are plain Go structs in `theme/presets.go`.

### Passing a theme to a brick

```go
// At construction — functional option
c := card.New(card.Title("My Card"), card.WithTheme(theme.Preset("dracula")))

// After construction — setter for live updates
c.SetTheme(newTheme)

// In your app's Update():
case theme.ThemeChangedMsg:
    m.theme = msg.Theme
    m.card.SetTheme(m.theme)
    m.footer.SetTheme(m.theme)
```

### Global manager (opt-in for apps)

The global manager exists for apps that want a single active theme shared
across all bricks. It is not required. If you never call `theme.SetTheme()`,
bricks that have no explicit theme set will call `theme.CurrentTheme()` and
get the default preset (catppuccin-mocha).

```go
theme.SetTheme("dracula")            // sets global + returns (Theme, error)
theme.PreviewTheme("nord")           // same but no persistence
theme.CurrentTheme() Theme           // read global active theme
theme.CurrentThemeName() string
theme.AvailableThemes() []string     // sorted, default first
theme.RegisterTheme("x", t)         // add custom theme
theme.Preset("tokyo-night")          // get a preset by name, no global state
theme.Names() []string               // all built-in preset names
```

---

## Package Responsibilities

### `theme/`

Go interface + 16 built-in preset structs + optional global manager.

- `Theme` — interface all presets and custom themes implement
- `BaseTheme` — embeddable struct implementing `Theme`, fill color fields
- `Preset(name) Theme` — returns a named preset, falls back to default
- `Names() []string` — all built-in names, default first
- Manager functions (`CurrentTheme`, `SetTheme`, `PreviewTheme`, `RegisterTheme`,
  `AvailableThemes`) — app-level convenience, not required by bricks

### `theme/styles/`

Pure rendering utility functions. No `System` struct. No theme dependency
at the package level — callers pass colors directly.

```go
styles.Row(bg, fg, width, content)      // full-width row with explicit Bg
styles.RowClip(bg, fg, width, content)  // clip ANSI content first, then paint
styles.ClipANSI(content, width)          // ANSI-safe truncation only
```

### `registry/bricks/surface/`

Ultraviolet-backed full-terminal cell buffer. Root canvas for every bento.

```go
surf := surface.New(width, height)
surf.Fill(bg)               // paint every cell — call first
surf.Draw(x, y, str)        // overlay: respects existing Bg
surf.DrawCenter(str)        // centered overlay for dialogs
surf.Render()               // → ANSI string for tea.NewView
```

### `registry/rooms/`

Named geometry patterns. Zero color. Zero theme.

```go
screen := rooms.Focus(w, h, body, footer)
screen := rooms.HolyGrail(w, h, 28, header, sidebar, main, footer)
screen := rooms.Rail(w, h, 9, sidebar, main)
```

Always composite room output through `surface` in app `View()`.

---

## Component Types

### Atomic

Examples: `badge`, `kbd`, `text`, `wordmark`, `separator`

- Returns a styled string — no background region ownership
- Does not need `Width()` — the caller sizes and places it
- Uses foreground color and bold/italic freely

### Container

Examples: `card`, `dialog`, `bar`, `tabs`

- Owns a **width × height region** on the surface
- Must set `Width()` on every row so every cell has explicit `Bg`
- Uses `styles.Row(bg, fg, width, content)` for all body rows

### Surface

`registry/bricks/surface/surface.go` — the full-terminal root canvas.
One per frame, sized to terminal dimensions. Not a UI component.

---

## Component Rules

1. **Colors come in, not out** — bricks use whatever `Theme` they were given via
   `WithTheme()` / `SetTheme()`. Fall back to `theme.CurrentTheme()` only when
   no theme was explicitly provided.
2. **Containers: every row has explicit Bg** — use `styles.Row()` or
   `.Background().Width().Render()`
3. **Atomics: no Width() required** — caller handles placement
4. **No imports between registry components** — each is standalone
5. **Theme tokens via interface methods** — `t.Background()`, `t.Text()`, etc.
   Never dot-accessors on structs (`t.Surface.Canvas` is gone).
6. **No `lipgloss.PlaceHorizontal` for surface-drawn rows** — use `Width()` instead

---

## Why Not a Framework?

Frameworks make the easy cases easy and the hard cases impossible. Every
non-trivial TUI eventually needs to reach inside and change something the
framework didn't anticipate. BentoTUI copies source so that modification is
always zero-friction — there is no "extend" API to learn, no lifecycle hooks,
no middleware. You just edit the file.
