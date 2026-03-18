# BentoTUI Bricks

API reference for every registry brick.
All bricks are copy-and-own — run `bento add <name>` to copy source into your project.

Bento-first policy: bricks are the official app-facing API. Users should be
able to build full apps from bricks without touching raw bubbles internals.

Some bricks use Charm primitives internally, but that remains an implementation
detail hidden behind Bento APIs.

Primitive policy: Bento does not ship a `spinner` brick. Use
`charm.land/bubbles/v2/spinner` directly.

---

## Shared module imports (import, don't copy)

### `theme`

```go
import "github.com/cloudboy-jh/bentotui/theme"

// Presets — no global state required
t := theme.Preset("dracula")         // returns a Theme interface value
names := theme.Names()               // all built-in preset names

// Global manager — app-level convenience
theme.SetTheme("tokyo-night")        // (Theme, error)
theme.PreviewTheme("nord")           // live preview, no persist
t := theme.CurrentTheme()            // fallback used by bricks with no explicit theme
theme.CurrentThemeName() string
theme.AvailableThemes() []string     // sorted, default first
theme.RegisterTheme("x", t)          // register a custom Theme
```

`Theme` is an interface. Methods return `color.Color`:

```go
t.Background()          t.BackgroundPanel()       t.BackgroundOverlay()
t.BackgroundInteractive()
t.CardChrome()          t.CardBody()              t.CardFrameFG()
t.CardFocusEdge()
t.Text()                t.TextMuted()             t.TextInverse()
t.TextAccent()
t.BorderNormal()        t.BorderSubtle()          t.BorderFocus()
t.Success()             t.Warning()               t.Error()
t.Info()
t.SelectionBG()         t.SelectionFG()
t.InputBG()             t.InputFG()               t.InputPlaceholder()
t.InputCursor()         t.InputBorder()
t.BarBG()               t.BarFG()
t.FooterBG()            t.FooterFG()              t.FooterMuted()
t.DialogBG()            t.DialogFG()              t.DialogBorder()
t.DialogScrim()
t.Name() string
```

Custom themes: embed `theme.BaseTheme`, fill the color fields, return a pointer.

### `theme/styles`

```go
import "github.com/cloudboy-jh/bentotui/theme/styles"

styles.Row(bg, fg, width, content)      // full-width row, explicit Bg on every cell
styles.RowClip(bg, fg, width, content)  // ANSI-safe clip then paint
styles.ClipANSI(content, width)          // truncation only, no painting
```

### `registry/rooms`

```go
import "github.com/cloudboy-jh/bentotui/registry/rooms"

screen := rooms.Focus(w, h, body, footer)
screen := rooms.Pancake(w, h, header, body, footer)
screen := rooms.Rail(w, h, 24, sidebar, main)
```

See `rooms.md` for the full API.

---

## Bricks

### `card`

The one content-container brick in BentoTUI. Replaces the old `panel` and
`elevated-card` — merged because they were the same component at different
elevation levels.

```go
import "yourmodule/bricks/card"

// Raised (default) — chrome header band + body slab
c := card.New(
    card.Title("Service Health"),
    card.Content(myWidget),
    card.WithTheme(t),          // optional — falls back to theme.CurrentTheme()
)

// Flat — plain titled container with separator line
c := card.New(
    card.Title("Sidebar"),
    card.Content(myWidget),
    card.Flat(),
)

c.SetSize(width, height)
c.Focus()
c.Blur()
c.IsFocused() bool
c.SetTheme(t theme.Theme)   // live update on ThemeChangedMsg
c.SetTitle(s string)
c.SetMeta(s string)
c.SetFooter(s string)
```

**Raised layout:** left accent edge (1 cell) + title/meta chrome band + body slab
rows + optional footer band + bottom chrome lane. Every row is explicitly painted.

**Flat layout:** title row + `───` separator + content rows with optional left
focus edge. Used for sidebars, panes, split-view regions.

Content receives `SetSize(width-2, height-reserved)` if it implements `Sizeable`.

Options:
- `card.Title(s)` — card title
- `card.Meta(s)` — secondary title line (raised only)
- `card.Footer(s)` — bottom chrome text (raised only)
- `card.Content(model)` — any `tea.Model`
- `card.Inset(n)` — outer margin cells (raised only, 0–4)
- `card.Flat()` — flat elevation
- `card.Raised()` — raised elevation (default, rarely needed)
- `card.WithTheme(t)` — explicit theme

---

### `bar`

Single-row header or footer strip. Cards truncate gracefully when width is tight.

```go
import "yourmodule/bricks/bar"

// Top bar
b := bar.New(
    bar.RoleTopBar(),
    bar.Left("my app"),
    bar.Right("v0.4.0"),
    bar.WithTheme(t),
)

// Anchored footer with keybind cards
footer := bar.New(
    bar.FooterAnchored(),
    bar.Left("~ project/path"),
    bar.Cards(
        bar.Card{Command: "enter", Label: "run",  Variant: bar.CardPrimary, Enabled: true, Priority: 3},
        bar.Card{Command: "ctrl+c", Label: "quit", Variant: bar.CardMuted,  Enabled: true, Priority: 2},
    ),
    bar.CompactCards(),
    bar.WithTheme(t),
)

b.SetSize(width, 1)
b.SetLeft(s string)
b.SetRight(s string)
b.SetStatusPill(s string)
b.SetCards([]bar.Card)
b.SetTheme(t theme.Theme)
```

Card variants: `CardNormal`, `CardPrimary`, `CardMuted`, `CardDanger`

Row roles: `RoleTop`, `RoleSubheader`, `RoleFooter`

`FooterAnchored()` renders in opencode-style: command keys bold, labels muted,
no chip backgrounds — clean anchored command row.

---

### `dialog`

Modal overlay manager plus built-in dialog types.

```go
import "yourmodule/bricks/dialog"

dm := dialog.New()
dm.SetTheme(t)

// Open
return m, func() tea.Msg {
    return dialog.Open(dialog.Confirm{
        DialogTitle: "Delete?",
        Message:     "This cannot be undone.",
        OnConfirm:   func() tea.Msg { return deleteMsg{} },
    })
}

// Route
updated, cmd := dm.Update(msg)
dm = updated.(*dialog.Manager)

// Render (composite in View())
if dm.IsOpen() {
    surf.DrawCenter(viewString(dm.View()))
}

dm.SetSize(width, height)
dm.IsOpen() bool
```

**`Confirm`** — yes/no. Manager handles `enter` (fires OnConfirm) and `esc` auto.

**`Custom`** — wraps any `tea.Model` as dialog content. Frame provided by Custom.

```go
dialog.Custom{
    DialogTitle: "Settings",
    Content:     mySettingsModel,
    Width:       60,
    Height:      20,
}
```

**`ThemePicker`** — live-previewing theme switcher. Broadcasts `theme.ThemeChangedMsg`
on cursor movement (preview) and on enter (confirm). ESC reverts.

```go
dialog.Open(dialog.Custom{
    DialogTitle: "Themes",
    Content:     dialog.NewThemePicker(),
    Width:       44,
    Height:      len(theme.AvailableThemes()) + 8,
})
```

**`CommandPalette`** — searchable grouped action list.

```go
palette := dialog.NewCommandPalette([]dialog.Command{
    {Label: "Switch theme", Group: "App", Keybind: "ctrl+t", Action: func() tea.Msg {
        return dialog.Open(dialog.Custom{Content: dialog.NewThemePicker()})
    }},
})
dialog.Open(dialog.Custom{DialogTitle: "Commands", Content: palette})
```

---

### `input`

Single-line text field wrapping `bubbles/textinput`.

```go
import "yourmodule/bricks/input"

i := input.New()
i.SetPlaceholder("Search…")
i.SetTheme(t)          // syncs styles immediately; also called on every View()
cmd := i.Focus()       // returns tea.Cmd — include in Init() or batch
i.Blur()
i.SetValue("text")
v := i.Value() string
i.SetSize(width, height)
```

---

### `list`

Scrollable list backed by `bubbles/list`. Delegate-driven row rendering with
explicit Bg on every cell.

```go
import "yourmodule/bricks/list"

l := list.New(200)           // max 200 items
l.Append("line")
l.Prepend("line")
l.AppendSection("Today")
l.AppendRow(list.Row{
    Primary:   "api",
    Secondary: "healthy",
    Tone:      list.ToneSuccess,
    RightStat: "36ms",
})
l.Clear()
l.Items() []string

l.SetFormatter(func(row list.Row, selected bool, width int) string {
    return row.Label  // plain content — delegate applies bg/fg
})
l.SetDensity(list.DensityCompact)
l.SetTheme(t)
l.SetSize(width, height)
l.Focus()
l.Blur()
```

---

### `table`

Data table backed by `bubbles/table`. Compact/borderless modes, optional grid
borders, column priority shrinking, per-column alignment.

```go
import "yourmodule/bricks/table"

t := table.New("Name", "Status", "Size")
t.SetVisualStyle(table.VisualGrid)
t.SetCompact(true)
t.SetBorderless(true)
t.SetColumnAlign(2, table.AlignRight)
t.SetColumnWidth(0, 18)
t.SetColumnPriority(0, 2)  // higher = shrinks last
t.AddRow("main.go", "ok", "4.2 KB")
t.Clear()
t.SetTheme(t)
t.SetSize(width, height)
t.Focus()
t.Blur()
```

---

### `badge`

Inline themed label for status/state chips.

```go
import "yourmodule/bricks/badge"

b := badge.New("beta")
b.SetVariant(badge.VariantInfo)
b.SetBold(true)
b.SetTheme(t)
```

Variants: `VariantNeutral`, `VariantInfo`, `VariantSuccess`, `VariantWarning`,
`VariantDanger`, `VariantAccent`

---

### `tabs`

Keyboard-navigable tab row.

```go
import "yourmodule/bricks/tabs"

t := tabs.New(
    tabs.Tab{ID: "overview", Label: "Overview"},
    tabs.Tab{ID: "logs",     Label: "Logs"},
)
t.Focus()
t.SetActive(1)
idx := t.Active()
t.SetTheme(theme)
t.SetSize(width, 1)
```

---

### `kbd`

Keyboard shortcut pair (`command label`) matching the bar-card visual language.

```go
import "yourmodule/bricks/kbd"

k := kbd.New("ctrl+k", "commands")
k.SetVariant("primary") // normal | primary | muted | danger
k.SetActive(true)
k.SetTheme(t)
```

---

### `select`

Single-choice inline picker backed by `bubbles/list`.

```go
import selectx "yourmodule/bricks/select"

s := selectx.New(
    selectx.Item{Label: "Tokyo Night", Value: "tokyo-night"},
    selectx.Item{Label: "Nord",        Value: "nord"},
)
s.Focus()
s.SetPlaceholder("Choose theme")
// enter/space: open or confirm; up/down (j/k): move; esc: close
```

---

### `checkbox`

Boolean toggle.

```go
import "yourmodule/bricks/checkbox"

c := checkbox.New("Enable live preview")
c.Focus()
c.Toggle()
checked := c.Checked()
```

---

### `progress`

Horizontal progress bar backed by `bubbles/progress`.

```go
import "yourmodule/bricks/progress"

p := progress.New(30)  // bar width in cells
p.SetLabel("Sync")
p.SetValue(0.42)       // [0,1]
p.SetShowPercent(true)
```

---

### `filepicker`

File and directory picker backed by `bubbles/filepicker`.

```go
import "yourmodule/bricks/filepicker"

fp := filepicker.New(".")
fp.SetAllowFiles(true)
fp.SetAllowDirectories(false)
fp.SetAllowedTypes(".go", ".md")
fp.SetShowHidden(false)
fp.SelectedPath() string
fp.HighlightedPath() string
```

---

### `package-manager`

Sequential install flow with spinner + progress bar.

```go
import packagemanager "yourmodule/bricks/package-manager"

pm := packagemanager.New([]string{"bubbles", "bubbletea", "lipgloss"})
pm.SetSize(width, 1)
pm.SetInstaller(func(pkg string) tea.Cmd { ... })
```

---

### `toast`

Stacked transient notifications.

```go
import "yourmodule/bricks/toast"

toasts := toast.New(3)  // max visible
id := toasts.Push("Saved settings", toast.Success)
toasts.Dismiss(id)
toasts.Clear()
```

---

### `text`

Static string in `Text()` color.

```go
import "yourmodule/bricks/text"

t := text.New("All systems operational")
t.SetText("Updated message")
t.SetSize(width, height)
```

---

### `wordmark`

Themed title/heading block.

```go
import "yourmodule/bricks/wordmark"

w := wordmark.New("BentoTUI")
w.SetBold(true)
```

---

### `separator`

Horizontal or vertical divider.

```go
import "yourmodule/bricks/separator"

h := separator.New(separator.Horizontal, 40)
v := separator.New(separator.Vertical, 8)
```

---

## Theme wiring in an app

```go
type model struct {
    theme   theme.Theme
    footer  *bar.Model
    content *card.Model
    dialogs *dialog.Manager
    w, h    int
}

func newModel() *model {
    t := theme.CurrentTheme()
    return &model{
        theme:   t,
        footer:  bar.New(bar.FooterAnchored(), bar.WithTheme(t), ...),
        content: card.New(card.Title("Main"), card.WithTheme(t), ...),
        dialogs: dialog.New(),
    }
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    ...
    case theme.ThemeChangedMsg:
        m.theme = msg.Theme
        m.footer.SetTheme(m.theme)
        m.content.SetTheme(m.theme)
    ...
}

func (m *model) View() tea.View {
    t := m.theme
    surf := surface.New(m.w, m.h)
    surf.Fill(t.Background())
    surf.Draw(0, 0, rooms.Focus(m.w, m.h, m.content, m.footer))
    if m.dialogs.IsOpen() {
        surf.DrawCenter(viewString(m.dialogs.View()))
    }
    v := tea.NewView(surf.Render())
    v.AltScreen = true
    v.BackgroundColor = t.Background()
    return v
}
```
