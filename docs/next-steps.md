# Next Steps

Three concrete items that are either partially done or blocked on a decision.
Everything else in the backlog lives in `roadmap.md`.

---

## 1. Finish `bento add <component>` — wire the embed

**Status:** The CLI scaffolding exists (`cmd/bento/add.go`) but currently only
prints shell `cp` commands instead of copying files.

**What it needs:**

```go
//go:embed ../../registry
var registryFS embed.FS
```

Added to `cmd/bento/add.go`, then the `add` command walks `registryFS` for the
requested component directory and writes each file into the user's
`components/<name>/` directory. Import paths stay as-is — they already point at
the real `bentotui` module deps, not local copies.

**Decision needed:** Should `bento add panel` write to `components/panel/` (relative
to `go.mod`) or ask the user for a destination? The registry default is a configurable
`components/ui/` path. Recommend: default to `components/<name>/`, flag to
override.

**File to edit:** `cmd/bento/add.go`

---

## 2. Finish `bento init` — update the generated template

**Status:** `cmd/bento/init.go` generates a starter `main.go` that still imports
the old `bentotui.New()` monolithic API (which no longer exists).

**What it needs:** The template should produce a minimal app that:

1. Imports `github.com/cloudboy-jh/bentotui/layout` and at least one registry
   component.
2. Optionally runs `bento add panel` and `bento add bar` before generating so
   the generated app has actual local copies to import.

The simplest correct template is essentially the starter-app in
`cmd/starter-app/main.go` trimmed down to ~60 lines.

**File to edit:** `cmd/bento/init.go` (the `starterTemplate` const)

---

## 3. `input.View()` calls `SetStyles()` on every frame

**Status:** Working correctly but allocates a `textinput.Styles` struct on every
render. For most apps this is fine. Under heavy update rates (e.g. streaming
output) it could show in a profile.

**Fix (when needed):** Add a `lastTheme string` field to `registry/input/input.go`.
In `View()`, check `theme.CurrentThemeName() != m.lastTheme` before calling
`SetStyles`. Only re-derive styles on theme change.

```go
func (m *Model) View() tea.View {
    if name := theme.CurrentThemeName(); name != m.lastTheme {
        m.input.SetStyles(styles.New(theme.CurrentTheme()).InputStyles())
        m.lastTheme = name
    }
    return tea.NewView(m.input.View())
}
```

This is a micro-optimisation — do it only if a profiler actually points here.

**File to edit:** `registry/input/input.go`
