# BentoTUI Components

API reference for every registry component and module dep.
All components are copy-and-own — run `bento add <name>` to copy source into your project.

Primitive policy: Bento does not ship a `spinner` component. Use
`charm.land/bubbles/v2/spinner` directly when you need loading indicators.

---

## Module deps (import, don't copy)

### `theme`

```go
import "github.com/cloudboy-jh/bentotui/theme"

theme.SetTheme("bento-rose")         // persist choice
theme.PreviewTheme("nord")           // live preview, no persist
theme.CurrentTheme() theme.Theme     // always read in View()
theme.CurrentThemeName() string
theme.AvailableThemes() []string     // sorted, default first
theme.RegisterTheme("x", t)         // add custom preset
```

`Theme` has tokens for every visual slot:

```go
t.Surface.{Canvas, Panel, Elevated, Overlay, Interactive}
t.Text.{Primary, Muted, Inverse, Accent}
t.Border.{Normal, Subtle, Focus}
t.State.{Info, Success, Warning, Danger}
t.Selection.{BG, FG}
t.Input.{BG, FG, Placeholder, Cursor, Border}
t.Bar.{BG, FG}
t.Dialog.{BG, FG, Border, Scrim}
```

### `registry/layouts`

```go
import (
    "charm.land/lipgloss/v2"
    "github.com/cloudboy-jh/bentotui/registry/components/surface"
    "github.com/cloudboy-jh/bentotui/registry/layouts"
    "github.com/cloudboy-jh/bentotui/theme"
)

screen := layouts.HolyGrail(width, height,
    28,       // sidebar width
    header,
    sidebar,
    main,
    footer,
)

overlay := layouts.Modal(width, height, 56, 16,
    layouts.Static(screen),
    dialog,
)
```

All 15 layouts call `SetSize` on each child, constrain each cell to its exact
allocation, and return a final `string`.

Use `surface` as the final compositor in full-screen apps:

```go
screen := layouts.Pancake(width, height, header, body, footer)
surf := surface.New(width, height)
surf.Fill(lipgloss.Color(theme.CurrentTheme().Surface.Canvas))
surf.Draw(0, 0, screen)
```

---

## Registry components

### `panel`

Titled, focusable content container with left-edge focus stripe.

```go
import "yourmodule/components/panel"

p := panel.New(
    panel.Title("Sidebar"),
    panel.Content(myWidget),   // any tea.Model
    panel.Elevated(),          // Surface.Elevated bg instead of Panel
)

p.SetSize(width, height)
p.Focus()
p.Blur()
p.IsFocused() bool
```

Layout inside panel:
- Row 0: title badge + title bar (1 row, only when Title set)
- Row 1: `───` separator (1 row, only when Title set)
- Rows 2…n: content lines with left-edge focus stripe when focused

Content receives `SetSize(width-2, height-titleRows)` if it implements `Sizeable`.

---

### `bar`

Single-row header or footer bar. Truncates cards gracefully when width is tight.

```go
import "yourmodule/components/bar"

b := bar.New(
    bar.RoleTopBar(),
    bar.StatusPill("LIVE"),
    bar.Left("my app"),
    bar.Right("v1.0"),
)

footer := bar.New(
    bar.FooterAnchored(),
    bar.Left("scope: nav"),
    bar.Cards(
        bar.Card{Command: "j/k", Label: "move", Enabled: true, Priority: 4},
        bar.Card{Command: "tab", Label: "focus tabs", Enabled: true, Priority: 3},
        bar.Card{Command: "q", Label: "quit", Enabled: true, Priority: 2},
    ),
)

b.SetSize(width, 1)
b.SetLeft(s string)
b.SetRight(s string)
b.SetStatusPill("LIVE")
b.SetCards([]bar.Card)
b.SetCompactCards(true)
b.SetRole(bar.RoleFooter)
b.SetAnchored(true)
```

**Card variants:** `CardNormal`, `CardPrimary`, `CardMuted`, `CardDanger`

Cards render as `command label` pairs. In compact mode they render denser. When
width is tight, labels drop first, then lower-priority cards drop before
higher-priority cards.

Row roles: `RoleTop`, `RoleSubheader`, `RoleFooter`.
Footer modes: `FooterModeNormal`, `FooterModeAnchored`.
Use `bar.FooterAnchored()` for vim-style focused command rows.
Use `StatusPill(...)` only when you have real runtime status metadata.
Starter and shipped bentos now default to minimal top/subheader content.

---

### `dialog`

Modal overlay manager plus built-in dialog types.

```go
import "yourmodule/components/dialog"

// In your root model
dm := dialog.New()

// Open from any Update() return
return m, func() tea.Msg { return dialog.Open(dialog.Confirm{
    DialogTitle: "Delete?",
    Message:     "This cannot be undone.",
    OnConfirm:   func() tea.Msg { return deleteMsg{} },
}) }

// Route messages through the manager
updated, cmd := dm.Update(msg)
dm = updated.(*dialog.Manager)

// Render — composite with surface in your View()
if dm.IsOpen() {
    surf.DrawCenter(viewString(dm.View()))
}

dm.SetSize(width, height)
dm.IsOpen() bool
```

#### `Confirm`

Simple yes/no. Manager handles `enter` (fires `OnConfirm`) and `esc` (closes) automatically.

```go
dialog.Confirm{
    DialogTitle: "Confirm",
    Message:     "Are you sure?",
    OnConfirm:   func() tea.Msg { return doSomethingMsg{} },
}
```

#### `Custom`

Wraps any `tea.Model` as dialog content. The dialog title frame is provided by
`Custom` — your content model renders inside it.

```go
dialog.Custom{
    DialogTitle: "Settings",
    Content:     mySettingsModel,
    Width:       60,
    Height:      20,
}
```

#### `ThemePicker`

Live-previewing theme switcher. Broadcasts `theme.ThemeChangedMsg` on navigation
and on confirm. ESC reverts to the theme active when the picker was opened.

```go
dialog.Open(dialog.Custom{
    DialogTitle: "Themes",
    Content:     dialog.NewThemePicker(),
})
```

#### `CommandPalette`

Searchable, grouped command list.

```go
palette := dialog.NewCommandPalette([]dialog.Command{
    {Label: "Switch theme", Group: "App", Keybind: "ctrl+t", Action: func() tea.Msg {
        return dialog.Open(dialog.Custom{Content: dialog.NewThemePicker()})
    }},
    {Label: "Quit", Group: "App", Keybind: "ctrl+c", Action: func() tea.Msg {
        return tea.Quit()
    }},
})

dialog.Open(dialog.Custom{
    DialogTitle: "Commands",
    Content:     palette,
})
```

---

### `list`

Scrollable log-style list. Returns **plain text** — the parent panel applies color.

```go
import "yourmodule/components/list"

l := list.New(200)  // max 200 items stored
l.Append("line added to bottom")
l.Prepend("line added to top")
l.AppendSection("Today")
l.AppendRow(list.Row{Label: "api", Status: "ok", Stat: "36ms"})
l.Clear()
l.Items() []string  // copy of current items

l.SetFormatter(func(row list.Row, selected bool, width int) string {
    return row.Label
})

l.SetSize(width, height)  // shows last N lines that fit
```

---

### `table`

Header row + data rows with optional compact/borderless modes and per-column
width/alignment control.

```go
import "yourmodule/components/table"

t := table.New("Name", "Status", "Size")
t.SetCompact(true)
t.SetBorderless(true)
t.SetColumnAlign(2, table.AlignRight)
t.SetColumnWidth(0, 18)
t.AddRow("main.go", "ok", "4.2 KB")
t.AddRow("go.mod", "ok", "1.1 KB")
t.Clear()

t.SetSize(width, height)
```

---

### `text`

Static string in `Text.Primary` color. For styled text, build a `lipgloss.Style`
directly in your app instead of using this component.

```go
import "yourmodule/components/text"

t := text.New("All systems operational")
t.SetText("Updated message")

t.SetSize(width, height)
```

---

### `input`

Single-line text field wrapping `bubbles/textinput`. Styles update from
`theme.CurrentTheme()` on every `View()` call — live theme switching works
without any extra wiring.

```go
import "yourmodule/components/input"

i := input.New()
i.SetPlaceholder("Search…")
i.Focus()   // returns tea.Cmd — include in Init() or batch
i.Blur()
i.SetValue("initial text")
i.Value() string

i.SetSize(width, height)
```

---

### `badge`

Inline themed label for compact status/state chips.

```go
import "yourmodule/components/badge"

b := badge.New("beta")
b.SetVariant(badge.VariantInfo)
b.SetBold(true)
```

Variants: `VariantNeutral`, `VariantInfo`, `VariantSuccess`, `VariantWarning`, `VariantDanger`, `VariantAccent`.

---

### `kbd`

Keyboard shortcut pair (`command label`) matching bar-card visual language.

```go
import "yourmodule/components/kbd"

k := kbd.New("ctrl+k", "commands")
k.SetVariant("primary") // normal, primary, muted, danger
k.SetActive(true)
```

---

### `wordmark`

Themed title/heading block.

```go
import "yourmodule/components/wordmark"

w := wordmark.New("BentoTUI")
w.SetBold(true)
```

---

### `select`

Single-choice inline picker with open/close and keyboard navigation.

```go
import selectx "yourmodule/components/select"

s := selectx.New(
    selectx.Item{Label: "Tokyo Night", Value: "tokyo-night"},
    selectx.Item{Label: "Nord", Value: "nord"},
)
s.Focus()
s.SetPlaceholder("Choose theme")

// keyboard:
// - enter/space: open or select current
// - up/down (or k/j): move cursor
// - esc: close
```

---

### `checkbox`

Boolean toggle input with optional focus behavior.

```go
import "yourmodule/components/checkbox"

c := checkbox.New("Enable live preview")
c.Focus()
c.Toggle()
c.Checked() bool
```

---

### `progress`

Horizontal progress bar with optional label and percentage text.

```go
import "yourmodule/components/progress"

p := progress.New(30)   // bar width in cells
p.SetLabel("Sync")
p.SetValue(0.42)        // clamped to [0,1]
p.SetShowPercent(true)
```

---

### `tabs`

Keyboard-navigable tab row.

```go
import "yourmodule/components/tabs"

t := tabs.New(
    tabs.Tab{ID: "overview", Label: "Overview"},
    tabs.Tab{ID: "logs", Label: "Logs"},
)
t.Focus()
t.SetActive(1)
t.Active() int
```

---

### `toast`

Stacked notification rows for non-blocking feedback.

```go
import "yourmodule/components/toast"

toasts := toast.New(3) // max visible
id := toasts.Push("Saved settings", toast.Success)
toasts.Dismiss(id)
toasts.Clear()
```

---

### `separator`

Horizontal or vertical divider.

```go
import "yourmodule/components/separator"

h := separator.New(separator.Horizontal, 40)
v := separator.New(separator.Vertical, 8)
```

---

## Complete Example

```go
package main

import (
    "fmt"

    tea "charm.land/bubbletea/v2"
    "charm.land/lipgloss/v2"
    "github.com/cloudboy-jh/bentotui/registry/components/surface"
    "github.com/cloudboy-jh/bentotui/registry/layouts"
    "github.com/cloudboy-jh/bentotui/theme"
    "yourmodule/components/bar"
    "yourmodule/components/dialog"
    "yourmodule/components/list"
    "yourmodule/components/panel"
)

func main() {
    if _, err := tea.NewProgram(newModel()).Run(); err != nil {
        fmt.Printf("error: %v\n", err)
    }
}

type model struct {
    header  *bar.Model
    footer  *bar.Model
    body    *panel.Model
    log     *list.Model
    dialogs *dialog.Manager
    w, h    int
}

func newModel() *model {
    log := list.New(100)
    log.Append("ready")

    body := panel.New(panel.Title("Main"), panel.Content(log))
    hdr := bar.New(
        bar.Left("my app"),
        bar.Cards(bar.Card{Command: "ctrl+t", Label: "theme", Enabled: true}),
    )

    ftr := bar.New(bar.Left("ctrl+c quit"))

    return &model{header: hdr, footer: ftr, body: body, log: log, dialogs: dialog.New()}
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if m.dialogs.IsOpen() {
        u, cmd := m.dialogs.Update(msg)
        m.dialogs = u.(*dialog.Manager)
        if tc, ok := msg.(theme.ThemeChangedMsg); ok {
            m.header.SetRight(tc.Name)
        }
        return m, cmd
    }
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.w, m.h = msg.Width, msg.Height
        m.dialogs.SetSize(m.w, m.h)
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "ctrl+t":
            return m, func() tea.Msg {
                return dialog.Open(dialog.Custom{
                    DialogTitle: "Themes",
                    Content:     dialog.NewThemePicker(),
                })
            }
        }
    }
    return m, nil
}

func (m *model) View() tea.View {
    t := theme.CurrentTheme()

    screen := layouts.Pancake(m.w, m.h,
        m.header,
        m.body,
        m.footer,
    )

    surf := surface.New(m.w, m.h)
    surf.Fill(lipgloss.Color(t.Surface.Canvas))
    surf.Draw(0, 0, screen)

    if m.dialogs.IsOpen() {
        surf.DrawCenter(viewString(m.dialogs.View()))
    }

    v := tea.NewView(surf.Render())
    v.AltScreen = true
    v.BackgroundColor = lipgloss.Color(t.Surface.Canvas)
    return v
}

func viewString(v tea.View) string {
    if v.Content == nil { return "" }
    if r, ok := v.Content.(interface{ Render() string }); ok { return r.Render() }
    if s, ok := v.Content.(interface{ String() string }); ok { return s.String() }
    return fmt.Sprint(v.Content)
}
```
