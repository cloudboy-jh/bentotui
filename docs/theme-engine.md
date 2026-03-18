# BentoTUI Theme System

v0.4.0 — `Theme` is an interface. Bricks accept themes as inputs. No mandatory global store.

---

## The model

`Theme` is a Go interface. Every preset and custom theme implements it.
Components call methods on it — `t.Background()`, `t.SelectionBG()`, `t.Text()`.

```go
// Any of these work
t := theme.Preset("dracula")           // named preset, no global
t := theme.CurrentTheme()              // global active theme
t := &MyCustomTheme{...}               // custom implementation of Theme interface
```

Pass it to a brick:

```go
// At construction
c := card.New(card.Title("file.go"), card.WithTheme(t))

// After construction (live update)
c.SetTheme(newTheme)
```

The brick uses whatever theme it was given. If none was given, it falls back
to `theme.CurrentTheme()`.

---

## Presets

16 built-in presets, no external dependencies:

| Name | Style |
|---|---|
| `catppuccin-mocha` | Default — soft dark purple |
| `catppuccin-macchiato` | Slightly cooler Catppuccin |
| `catppuccin-frappe` | Lighter Catppuccin |
| `dracula` | Classic dark purple |
| `tokyo-night` | Deep blue-grey |
| `tokyo-night-storm` | Tokyo Night, lighter base |
| `nord` | Arctic grey-blue |
| `bento-rose` | Warm rose-dark |
| `gruvbox-dark` | Earth tones |
| `monokai-pro` | Vibrant warm dark |
| `kanagawa` | Soft Japanese ink |
| `rose-pine` | Pastel dark |
| `ayu-mirage` | Blue-grey mirage |
| `one-dark` | Atom One Dark |
| `material-ocean` | Deep ocean blue |
| `github-dark` | GitHub dark mode |

```go
t := theme.Preset("tokyo-night")
names := theme.Names()           // all names, default first
```

---

## Custom themes

Embed `BaseTheme`, fill the color fields:

```go
type MyTheme struct {
    theme.BaseTheme
}

func NewMyTheme() *MyTheme {
    m := &MyTheme{}
    m.ThemeName = "my-theme"
    m.BackgroundColor = lipgloss.Color("#1a1a2e")
    m.TextColor = lipgloss.Color("#e0e0e0")
    m.SelectionBGColor = lipgloss.Color("#7c4dff")
    m.SelectionFGColor = lipgloss.Color("#ffffff")
    // ... fill all fields
    return m
}

// Register for use with theme picker
theme.RegisterTheme("my-theme", NewMyTheme())
```

All `color.Color` fields — use `lipgloss.Color("hex")` to create them.

---

## Global manager

The global manager is optional app-level infrastructure. Bricks do not require it.

```go
theme.SetTheme("dracula")              // set global + returns (Theme, error)
theme.PreviewTheme("nord")             // live preview, no persist
theme.CurrentTheme() Theme             // read global (fallback used by bricks)
theme.CurrentThemeName() string
theme.AvailableThemes() []string       // registered names, default first
theme.RegisterTheme("x", t)           // add to global registry
```

---

## Live theme switching

The theme picker dispatches `theme.ThemeChangedMsg` on every cursor move
(preview) and on enter (confirm). ESC reverts to the theme active when the
picker opened.

Your app handles it:

```go
case theme.ThemeChangedMsg:
    m.theme = msg.Theme
    m.footer.SetTheme(m.theme)
    m.card.SetTheme(m.theme)
    m.input.SetTheme(m.theme)
    // any other stateful bricks
```

No framework magic. You decide which bricks get the new theme. This is the
entire mechanism — it's just a message and a setter.

---

## Using themes without the global

For CLI tools, tests, or any context where you don't want the global manager:

```go
t := theme.Preset("dracula")

// Render a card with no global state
c := card.New(
    card.Title("Results"),
    card.WithTheme(t),
)
c.SetSize(80, 10)
// use c.View() to render
```

No `init()`, no mutex, no global store touched.
