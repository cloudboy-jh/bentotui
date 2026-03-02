# BentoTUI Roadmap

## Current State (v0.2)

The repository completed a full structural refactor in v0.2:

- Deleted the monolithic `bentotui.New()` framework API
- Deleted `app/`, `core/`, `ui/`, `bentotui.go`
- Moved `core/theme/` → `theme/`, `ui/styles/` → `styles/`, `core/layout/` → `layout/`
- Created `registry/` with clean rewrites of every component
- Every component reads `theme.CurrentTheme()` in `View()` — no stored theme state
- Every row uses a single `lipgloss.NewStyle().Background().Width().Render()` call
- Starter-app rewritten as a full component showcase with live theme switching

See [next-steps.md](./next-steps.md) for the three immediate known gaps.

## Backlog

### CLI (`cmd/bento`)

- [ ] `bento add` — wire `//go:embed registry` so files are actually copied
- [ ] `bento init` — update generated template to use new package paths
- [ ] `bento list` — show available components with one-line descriptions
- [ ] `bento upgrade <component>` — diff registry version against copied version

### New Components

- [ ] `tabs` — keyboard-navigable tab bar with panel content area
- [ ] `select` — dropdown/select widget (wraps or reimplements bubbles/list)
- [ ] `spinner` — animated loading indicator
- [ ] `progress` — horizontal progress bar with theme colors
- [ ] `checkbox` — boolean toggle
- [ ] `badge` — inline colored label (useful inside panel titles)

### Layout Enhancements

- [ ] `layout.Grid` — fixed-column grid (simpler than manual Horizontal+Vertical nesting)
- [ ] Scrollable `Split` — allow body regions to scroll independently

### Developer Experience

- [ ] Layout debugger mode — render allocation boundaries as colored overlays
- [ ] `go test ./registry/...` — snapshot tests for every component's rendered output
- [ ] More starter-app variants — e.g. a file browser, a log viewer

### Theme

- [ ] Light theme support — adapter needs a `fromLightTint` variant with inverted surface mapping
- [ ] Theme validation CLI — `bento validate-theme mytheme.json`

## Non-Goals

- **Mobile / small screens** — BentoTUI assumes a reasonably large terminal.
  Components have minimum widths but no responsive breakpoint system.
- **Mouse support** — Bubble Tea v2 has mouse events; none of the current
  components handle them. Not planned unless a component clearly needs it (e.g. scrollbars).
- **Accessibility** — Terminal screen reader support depends on the terminal
  emulator, not the TUI library. No plans here.

---

Last updated: 2026-03-02
