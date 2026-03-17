# Bento Brick Coloring Rules

These rules exist to stop recurring visual bugs — color bleed, invisible selected text,
mismatched backgrounds — from being introduced when building or modifying bricks.
Follow them every time. They are not optional.

---

## Rule 1 — Always use `Selection.FG` on foreground, `Selection.BG` on background

**Never pass a `BG` token as a `Foreground()` argument, or a `FG` token as a `Background()` argument.**

```go
// WRONG — Selection.BG is a background color; as FG it becomes near-invisible text
lipgloss.NewStyle().Foreground(lipgloss.Color(t.Selection.BG))

// CORRECT
lipgloss.NewStyle().
    Background(lipgloss.Color(t.Selection.BG)).
    Foreground(lipgloss.Color(t.Selection.FG))
```

The `Selection.*`, `Bar.*`, `Footer.*`, `Card.*`, `Dialog.*`, `Input.*` token pairs
always come in `BG`/`FG` pairs. Both must be applied together. Using one without the
other will produce invisible or low-contrast text on every theme.

**Quick check:** if you see `t.Selection.BG` inside a `Foreground()` call, it is a bug.

---

## Rule 2 — Apply theme tokens only in `View()`, never during data mutations

Theme must be read once per render, not during `AddRow`, `SetSize`, `Focus`, `Blur`,
`SetItems`, or any other state-mutation path.

```go
// WRONG — theme read on every data mutation
func (t *Model) AddRow(cells ...string) {
    t.rows = append(t.rows, cells)
    t.syncInner() // ← syncInner reads theme.CurrentTheme() — wrong lifecycle
}

// CORRECT — data sync is theme-free; theme read only at render time
func (t *Model) AddRow(cells ...string) {
    t.rows = append(t.rows, cells)
    t.syncData() // ← no theme access here
}

func (t *Model) View() tea.View {
    t.applyTheme() // ← theme read happens here, once per render
    return tea.NewView(t.inner.View())
}
```

**Why:** theme switches (`theme.SetTheme(...)`) take effect on the next `View()` call.
If theme is baked into data structures during mutations, a theme change before `View()`
will not take effect until the next mutation triggers a re-sync — which is unpredictable
and means live theme switching is broken.

---

## Rule 3 — Never set `Background` on Bubbles `Cell` style when also using `Selected`

When wrapping `bubbles/table` (or similar row-based Bubbles components), the upstream
`renderRow()` function applies `Cell.Render()` per-cell first, joins the cells into a
string, then applies `Selected.Render()` to the joined string for the selected row.

If `Cell` carries a `Background`, those ANSI escape sequences are embedded into the
joined string *before* `Selected` repaints. The result: padding characters inside each
cell carry the `Cell` BG color, while `Selected` wraps the whole row with `Selected` BG.
The terminal sees conflicting background ANSI codes within the same row — color bleed.

```go
// WRONG — Cell.Background causes bleed inside Selected rows
styles.Cell = lipgloss.NewStyle().
    Background(lipgloss.Color(cellBG)). // ← ANSI codes now embedded in joined row string
    Foreground(lipgloss.Color(textFG))

styles.Selected = lipgloss.NewStyle().
    Background(lipgloss.Color(selectedBG)). // ← rewraps already-colored string: bleed
    Foreground(lipgloss.Color(selectedFG))

// CORRECT — Cell has no Background; Selected owns the entire row background
styles.Cell = lipgloss.NewStyle().
    Foreground(lipgloss.Color(textFG)) // No Background

styles.Selected = lipgloss.NewStyle().
    Background(lipgloss.Color(selectedBG)). // Full control — no conflict
    Foreground(lipgloss.Color(selectedFG))
```

**Header is exempt:** headers are never re-wrapped by `Selected`, so `Header.Background`
is safe.

---

## Rule 4 — Every rendered row must own its full width with explicit background

Any row emitted to the terminal must carry an explicit background color for its full
width. Rows without explicit backgrounds are transparent — they inherit whatever color
is behind them, which changes unpredictably across parent panels, surfaces, and themes.

Use `styles.Row()` or `styles.RowClip()` for all surface-facing row output:

```go
// WRONG — bare string, transparent background
return prefix + content

// CORRECT — explicit bg/fg, full width owned
return styles.Row(bg, fg, width, content)
// or for content that may overflow:
return styles.RowClip(bg, fg, width, content)
```

This applies to:
- List row delegates
- Select item delegates
- Toast rows
- Custom content rendered into panels or cards
- Any row assembled with string concatenation before being painted to a surface

**Exception:** content that is immediately passed into a `lipgloss.Style.Width(w).Render()`
chain that carries a `Background` is fine — the `Width()` expansion will fill the row.

---

## Rule 5 — Set all fields of upstream `Styles` structs, not a partial subset

When applying theme colors to a Bubbles component's `Styles` struct, set every field
that struct exposes. Skipped fields fall back to upstream defaults — which are always
hardcoded ANSI256 colors, not theme colors.

```go
// WRONG — Symlink style skipped; symlinks show upstream ANSI color "36" (cyan) on all themes
s.Cursor = ...
s.Directory = ...
s.File = ...
// s.Symlink ← missing

// CORRECT
s.Cursor = ...
s.Directory = ...
s.File = ...
s.Symlink = lipgloss.NewStyle().Foreground(lipgloss.Color(pick(t.Text.Accent, t.Text.Primary)))
```

Before shipping a new Bubbles wrapper, look at the upstream `Styles` struct definition
and check every field is covered. If a field is intentionally left at upstream default,
add a comment saying so.

---

## Rule 6 — Never scan rendered output to detect selection state

Do not call `m.inner.View()` and then scan the returned string for cursor markers like
`"> "` to determine which row is selected. This is fragile and breaks whenever content
contains the marker string.

Instead, query selection state from the model directly — either via `m.Index()` in a
delegate `Render()` method, or by tracking cursor state in your wrapper model.

```go
// WRONG — string scanning to detect selected row
menu := m.inner.View()
for _, line := range strings.Split(menu, "\n") {
    if strings.HasPrefix(strings.TrimLeft(line, " "), "> ") { // ← fragile
        // paint as selected
    }
}

// CORRECT — query state in the delegate where you have direct access to it
func (d myDelegate) Render(w io.Writer, m bubbleslist.Model, index int, item bubbleslist.Item) {
    isSelected := index == m.Index() // ← source of truth
    if isSelected {
        // paint with selection colors
    }
}
```

---

## Rule 7 — Set `focus` state visually in both focused and blurred cases

Every interactive brick must render a distinct visual difference between focused and
blurred states. Blurred does not mean invisible — it means reduced contrast, not removed.

```go
// Focus: full selection contrast
styles.Selected = lipgloss.NewStyle().
    Bold(true).
    Background(lipgloss.Color(t.Selection.BG)).
    Foreground(lipgloss.Color(t.Selection.FG))

// Blur: position indicator without high contrast
styles.Selected = lipgloss.NewStyle().
    Foreground(lipgloss.Color(t.Text.Accent))
    // No bold, no Background — position is visible but clearly not active
```

The standard token mapping for states:

| State | Background token | Foreground token |
|---|---|---|
| Selected + Focused | `Selection.BG` | `Selection.FG` |
| Selected + Blurred | none or `Surface.Interactive` | `Text.Accent` |
| Normal row | `Surface.Panel` | `Text.Primary` |
| Section header | `Surface.Panel` | `Text.Muted` |
| Disabled / muted | `Surface.Canvas` | `Text.Muted` |
| Error / danger | `State.Danger` | `Text.Inverse` |

---

## Rule 8 — Read `theme.CurrentTheme()` at most once per `View()` call

Call `theme.CurrentTheme()` once at the top of `View()` (or the method it calls) and
pass the result down. Do not call it inside loops, per-row, or per-delegate-render.

```go
// WRONG — called per row (inside delegate Render, which is called per item)
func (d delegate) Render(...) {
    t := theme.CurrentTheme() // ← called N times per View()
}

// CORRECT — called once, passed via the delegate's owner pointer
func (d delegate) Render(...) {
    t := theme.CurrentTheme() // acceptable in delegate since it runs once per item;
    // but prefer: store the theme snapshot on the owner model in View() and read from there
}
```

If a delegate is used (list, select), either:
- Read `theme.CurrentTheme()` once in the delegate `Render()` (acceptable — it's fast),
- Or snapshot the theme on the owner model in `View()` and read it from the owner pointer
  in the delegate (preferred for consistency).

---

## Token quick-reference

| Token | Correct use |
|---|---|
| `Surface.Canvas` | Full-screen base background |
| `Surface.Panel` | Panel, card, list row background |
| `Surface.Interactive` | Hovered/blurred-selected background |
| `Selection.BG` | Selected row background only |
| `Selection.FG` | Selected row foreground only |
| `Text.Primary` | Normal body text |
| `Text.Muted` | Secondary, section headers, hints |
| `Text.Accent` | Blurred-selected indicator, links, icons |
| `Text.Inverse` | Text on high-contrast selection backgrounds |
| `Border.Focus` | Focused panel border, fallback selection BG |
| `Border.Normal` | Default unfocused border |
| `Border.Subtle` | Dividers, cell separators |
| `State.Danger` | Error state background |
| `State.Warn` | Warning state indicator |
| `State.Success` | Success state indicator |

---

## Checklist before shipping any new brick

- [ ] `Selection.BG` used only as `Background`, `Selection.FG` only as `Foreground`
- [ ] Theme read only in `View()` (or `applyTheme()` called from `View()`), never in mutation paths
- [ ] `Cell` style has no `Background` if the component uses `Selected` to wrap full rows
- [ ] All fields of upstream `Styles` structs are set (or explicitly noted as intentionally defaulted)
- [ ] All rendered rows go through `styles.Row()` or `styles.RowClip()` or equivalent `Background().Width().Render()` chain
- [ ] No string scanning of `View()` output to detect selection state
- [ ] Focused and blurred states produce visually distinct (not identical) output
- [ ] `theme.CurrentTheme()` called once per render, not per row or per mutation
