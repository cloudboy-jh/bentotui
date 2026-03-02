# BentoTUI Architecture

Status: Active reference — updated after v0.2 refactor (2026-03-02)

## Overview

BentoTUI is a **registry of copy-and-own components**, not a framework. The
mental model is ShadCN/UI applied to Go TUIs: you copy the component source into
your project, modify it freely, and the only stable dependency is the small set
of shared packages (`theme`, `styles`, `layout`) that you import as a normal Go
module.

## Layer Diagram

```
your app code
     │
     ▼
┌────────────────────────────────────────────────────────────┐
│  registry components  (copied into your project)           │
│                                                            │
│  panel  bar  dialog  list  table  text  input              │
│                                                            │
│  Rules:                                                    │
│  - theme.CurrentTheme() in View() — no stored theme state  │
│  - lipgloss.NewStyle().Background().Width().Render(plain)  │
│  - widgets return plain text; containers own all color     │
│  - no imports between components in the registry           │
└─────────────────────────┬──────────────────────────────────┘
                          │ imports (real go module deps)
                          ▼
┌────────────────────────────────────────────────────────────┐
│  bentotui module                                           │
│                                                            │
│  theme/   — goroutine-safe global theme store              │
│             SetTheme(), CurrentTheme(), RegisterTheme()    │
│                                                            │
│  styles/  — theme.Theme → lipgloss.Style mapping          │
│             styles.New(t).PanelTitleBar(focused)           │
│                                                            │
│  layout/  — Horizontal/Vertical split with Fixed/Flex      │
│             SetSize() propagates to Sizeable children      │
└────────────────────────────────────────────────────────────┘
```

## Three Core Rules

### 1. Theme is always read at render time

Components never store `m.theme theme.Theme`. They call `theme.CurrentTheme()`
inside `View()`. This makes live theme switching free — no `SetTheme()` calls
needed anywhere in the message loop.

```go
// Correct — reads global at render time
func (m *Model) View() tea.View {
    t := theme.CurrentTheme()
    // ...
}

// Wrong — stale after theme changes
func (m *Model) View() tea.View {
    t := m.theme // don't do this
    // ...
}
```

### 2. One lipgloss call owns every cell

Each rendered row uses a single `lipgloss.NewStyle().Background().Width().Render(plainText)`
call. The plain text is stripped of ANSI before rendering. This prevents
background color bleed-through where inner ANSI escape codes fight the outer
background at the cell level.

```go
// Correct — single call, lipgloss owns every cell
lipgloss.NewStyle().
    Background(lipgloss.Color(bg)).
    Foreground(lipgloss.Color(fg)).
    Width(w).
    Render(ansi.Strip(line))

// Wrong — canvas overlay allows inner codes to bleed through
lipgloss.NewCanvas(
    lipgloss.NewLayer(bgBlock).Z(0),
    lipgloss.NewLayer(styledText).Z(1), // text ANSI fights bg ANSI
)
```

### 3. Widgets return plain text; containers apply color

Widgets (`list`, `table`, `text`) return unstyled strings. The containing `panel`
applies the background/foreground colors via its `contentRow()` function. This
means the panel background is always uniform — it doesn't inherit colors from
child content.

## Package Responsibilities

### `theme/`

Global theme store. Goroutine-safe via `sync.RWMutex`.

- `CurrentTheme() Theme` — read active theme (RLock)
- `SetTheme(name) (Theme, error)` — set + persist (Lock)
- `PreviewTheme(name) (Theme, error)` — set without persisting (Lock)
- `RegisterTheme(name, Theme) error` — add custom theme
- `AvailableThemes() []string` — sorted list, default first

`init()` loads builtins, then attempts to restore a persisted theme from
`$XDG_CONFIG_HOME/bentotui/theme.json`. Safe to call `SetTheme()` from
`main()` before `tea.NewProgram().Run()`.

### `styles/`

Maps `theme.Theme` tokens to `lipgloss.Style` constructors. Centralises
every color decision so components don't scatter hex literals.

```go
sys := styles.New(theme.CurrentTheme())
sys.PanelTitleBar(focused bool) lipgloss.Style
sys.PanelTitleBadge(focused bool) lipgloss.Style
sys.FocusAccent() lipgloss.Style
sys.Divider() lipgloss.Style
sys.DialogFrame() lipgloss.Style
sys.PaletteItem(selected bool) lipgloss.Style
sys.InputStyles() textinput.Styles
// ... etc
```

### `layout/`

Composable split-pane layout. No bentotui-specific deps — any `tea.Model`
works as a child.

```go
root := layout.Horizontal(
    layout.Fixed(30, sidebar),   // exact 30 cols
    layout.Flex(1, main),        // remaining space
    layout.Flex(1, detail),      // equal share of remaining
).WithGutterColor(t.Border.Subtle)

root.SetSize(termW, termH) // propagates to Sizeable children
```

**Exported surface (stable contract):**
`Horizontal`, `Vertical`, `Fixed`, `Flex`, `Split` (type), `Item` (type),
`Model` (interface), `Sizeable` (interface), `SetSize`, `GetSize`,
`WithGutterColor`, `SetGutterColor`, `DebugLayout`.

### `registry/`

Components you copy into your project. Each is a self-contained `.go` file
(or small directory) that imports only the module deps above plus bubbletea
and lipgloss.

| Component | Notes |
|-----------|-------|
| `panel`   | `Title`, `Content`, `Elevated`, `Scrollable` options. Focus-accent left stripe. |
| `bar`     | `Left`, `Right`, `Cards`, `LeftCard`, `RightCard`. Truncates cards to fit width. |
| `dialog`  | `Manager`, `Confirm`, `Custom`, `ThemePicker`, `CommandPalette` |
| `list`    | Plain text output. `Append`, `Prepend`, `Clear`. Shows last N lines. |
| `table`   | Header + rows. Column width is `totalWidth / colCount`. |
| `text`    | Static `string` in `Text.Primary` color. |
| `input`   | Wraps `bubbles/textinput`. Styles updated every `View()`. |

## `dialog` Internals

`Manager` hosts zero or one active `Dialog`. On `OpenMsg` it calls
`dialog.SetSize(width, height)` immediately. On `CloseMsg` it sets
`active = nil`.

`Manager.Update` handles `Confirm` esc/enter automatically. All other dialog
types (`ThemePicker`, `CommandPalette`, `Custom`) handle their own esc.

`Custom` wraps any `tea.Model` as dialog content — it doesn't import `panel`
or any other registry component. Users pass whatever they want.

## Why Not a Framework?

Frameworks make the easy cases easy and the hard cases impossible. Every
non-trivial TUI eventually needs to reach inside and change something the
framework didn't anticipate. BentoTUI copies source so that modification is
always zero-friction — there is no "extend" API to learn, no `Middleware` pattern,
no lifecycle hooks. You just edit the file.

The cost is that updates to the registry don't automatically flow to copies in
user projects. The benefit is that users are never blocked by the registry's
assumptions.
