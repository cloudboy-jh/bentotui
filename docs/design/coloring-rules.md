# Bento Brick Coloring Rules

These rules prevent the recurring visual bugs — color bleed, invisible selected
text, mismatched backgrounds — that appear when building or modifying bricks.
Follow them every time.

---

## Rule 1 — SelectionBG is a background, SelectionFG is a foreground. Never swap.

```go
// Wrong — SelectionBG() as a foreground produces near-invisible text
lipgloss.NewStyle().Foreground(t.SelectionBG())

// Correct
lipgloss.NewStyle().
    Background(t.SelectionBG()).
    Foreground(t.SelectionFG())
```

All `BG`/`FG` pairs (`SelectionBG/FG`, `BarBG/FG`, `FooterBG/FG`, `DialogBG/FG`,
`InputBG/FG`) must be applied together. Using one without the other will produce
invisible or low-contrast text on every theme.

---

## Rule 2 — Read theme in View(), never during data mutations

```go
// Wrong — theme read on every AddRow
func (t *Model) AddRow(cells ...string) {
    t.rows = append(t.rows, cells)
    t.syncInner() // ← reads theme — wrong lifecycle
}

// Correct — data sync is theme-free; theme read only at render time
func (t *Model) AddRow(cells ...string) {
    t.rows = append(t.rows, cells)
    t.syncData() // ← no theme
}

func (t *Model) View() tea.View {
    t.applyTheme() // ← theme read here, once per render
    return tea.NewView(t.inner.View())
}
```

---

## Rule 3 — No Background on bubbles/table Cell when using Selected

`bubbles/table` applies `Cell.Render()` per cell, joins them, then applies
`Selected.Render()` to the joined row string. If `Cell` carries a `Background`,
those ANSI codes are already embedded before `Selected` rewraps — causing
padding cells to show `Cell` BG while `Selected` BG only covers the outer wrap.

```go
// Wrong — bleed on selected row padding
styles.Cell = lipgloss.NewStyle().Background(cellBG).Foreground(textFG)

// Correct — Cell has no Background; Selected owns the entire row
styles.Cell = lipgloss.NewStyle().Foreground(textFG)
styles.Selected = lipgloss.NewStyle().Background(selectedBG).Foreground(selectedFG)
```

Header is exempt — headers are never re-wrapped by `Selected`.

---

## Rule 4 — Every row must own its full width with explicit background

```go
// Wrong — transparent row inherits whatever is behind it
return prefix + content

// Correct — full width, explicit Bg on every cell
return styles.Row(bg, fg, width, content)
return styles.RowClip(bg, fg, width, content)
```

Applies to: list row delegates, select delegates, toast rows, custom content
rendered into cards, any row assembled before being painted to a surface.

Exception: content passed immediately into `.Background(x).Width(w).Render()`
is fine — `Width()` expansion fills the row.

---

## Rule 5 — Set all fields of upstream Styles structs

When applying theme colors to a bubbles component's `Styles` struct, set every
field. Skipped fields fall back to upstream hardcoded ANSI256 colors.

```go
// Wrong — Symlink style skipped; shows upstream cyan on all themes
s.Cursor = ...
s.Directory = ...
s.File = ...
// s.Symlink ← missing

// Correct
s.Symlink = lipgloss.NewStyle().Foreground(t.TextAccent())
```

Before shipping a bubbles wrapper, check every field in the upstream `Styles`
struct. If a field is intentionally left at upstream default, add a comment.

---

## Rule 6 — Never scan rendered output to detect selection state

```go
// Wrong — fragile, breaks when content contains "> "
menu := m.inner.View()
for _, line := range strings.Split(menu, "\n") {
    if strings.HasPrefix(line, "> ") { ... }
}

// Correct — query state in the delegate
func (d myDelegate) Render(w io.Writer, m bubbleslist.Model, index int, item bubbleslist.Item) {
    isSelected := index == m.Index()
    if isSelected {
        // paint with selection colors
    }
}
```

---

## Rule 7 — Focused and blurred states must look different

Blurred does not mean invisible — it means reduced contrast.

```go
// Focused: full selection contrast
styles.Selected = lipgloss.NewStyle().
    Bold(true).
    Background(t.SelectionBG()).
    Foreground(t.SelectionFG())

// Blurred: position visible but clearly not active
styles.Selected = lipgloss.NewStyle().
    Foreground(t.TextAccent())
```

Standard state-to-token mapping:

| State | Background | Foreground |
|---|---|---|
| Selected + Focused | `SelectionBG()` | `SelectionFG()` |
| Selected + Blurred | `BackgroundInteractive()` or none | `TextAccent()` |
| Normal row | `BackgroundPanel()` | `Text()` |
| Section header | `BackgroundPanel()` | `TextMuted()` |
| Disabled / muted | `Background()` | `TextMuted()` |
| Error / danger | `Error()` | `TextInverse()` |

---

## Rule 8 — Use theme interface methods, not struct field accessors

```go
// Wrong — old struct-based token access (removed in v0.4.0)
lipgloss.Color(t.Surface.Canvas)
lipgloss.Color(t.Text.Muted)
lipgloss.Color(t.Selection.BG)

// Correct — interface methods, return color.Color directly
t.Background()
t.TextMuted()
t.SelectionBG()
```

Theme methods return `color.Color`. Pass them directly to lipgloss — do not
wrap in `lipgloss.Color()` (that takes a string, not a `color.Color`).

```go
// Correct
lipgloss.NewStyle().Background(t.BackgroundPanel()).Foreground(t.Text())

// Wrong — double-wrapping
lipgloss.NewStyle().Background(lipgloss.Color(t.BackgroundPanel()))
```

---

## Token quick-reference

| Method | Correct use |
|---|---|
| `Background()` | Full-screen canvas (`surf.Fill`) |
| `BackgroundPanel()` | Panel, card, list row background |
| `BackgroundInteractive()` | Hovered / blurred-selected background |
| `CardChrome()` | Raised card header band |
| `CardBody()` | Raised card content slab |
| `SelectionBG()` | Selected row background only |
| `SelectionFG()` | Selected row foreground only |
| `Text()` | Normal body text |
| `TextMuted()` | Secondary text, section headers, hints |
| `TextAccent()` | Blurred selection indicator, links |
| `TextInverse()` | Text on high-contrast selection backgrounds |
| `BorderFocus()` | Focused border / fallback selection |
| `BorderNormal()` | Default unfocused border |
| `BorderSubtle()` | Dividers, cell separators |
| `Error()` | Error state background |
| `Warning()` | Warning state |
| `Success()` | Success state |
| `Info()` | Info state |

---

## Checklist before shipping a new brick

- [ ] `SelectionBG()` used only as `Background()`, `SelectionFG()` only as `Foreground()`
- [ ] Theme read only in `View()` (or `applyTheme()` called from `View()`), never in mutation paths
- [ ] `Cell` style has no `Background` if the component uses `Selected` to wrap full rows
- [ ] All fields of upstream `Styles` structs are set (or explicitly noted as intentional defaults)
- [ ] All rendered rows go through `styles.Row()` / `styles.RowClip()` or `.Background().Width().Render()`
- [ ] No string scanning of `View()` output to detect selection state
- [ ] Focused and blurred states produce visually distinct output
- [ ] Brick accepts `WithTheme(t theme.Theme)` option and `SetTheme(t theme.Theme)` setter
- [ ] No `lipgloss.Color(t.Method())` wrapping — pass `color.Color` values directly
