# Changelog

All notable changes to this project will be documented in this file.

The format follows Keep a Changelog style and this project targets Semantic Versioning.

## [Unreleased]

### Changed

- Moved shared style helpers from `styles/` to `theme/styles/` and updated all imports to `github.com/cloudboy-jh/bentotui/theme/styles`.
- Reworked `registry/bentos/app-shell` into a scenario-driven validation bento with modular `state/`, `ui/`, and `scenarios/` packages.
- Added built-in validation scenarios (`layout`, `hierarchy`, `footer`, `list`, `overlay`, `stress`) with pass/warn/fail checks and diagnostics metrics.
- Updated docs and README to formalize the **Untouchable Theme Engine** model and lock the architecture language to `bentos + rooms + bricks`.
- Added room-level separation options (`WithGutter`, `WithDivider`) to split/drawer room primitives.
- Added anchored footer card style modes in bar (`plain`, `chip`, `mixed`) and scenario coverage in app-shell.
- Expanded list row structure with typed fields (`Primary`, `Secondary`, `RightStat`, `Tone`, `SelectedStyle`) while preserving existing `Label`/`Status`/`Stat` compatibility.
- Added new `elevated-card` brick for raised section containers (`Title` + `Content`) and wired it into dashboard/app-shell composition paths.

## [0.3.4] - 2026-03-14

### Changed

- Renamed `registry/layouts` to `registry/rooms` and updated starter app, shipped bentos, CLI scaffold output, and doctor checks to the new import path.
- Removed legacy `layout/` split package to eliminate overlap with room-based composition.
- Added footer anchored semantic tokens to the theme model and decoupled anchored footer row colors from selection colors when footer tokens are present.
- Extended theme validation with required surface hierarchy checks (`canvas/panel/elevated/interactive`) and added coverage across all built-in presets.
- Added ANSI-safe shared clipping helpers in `styles` (`ClipANSI`, `RowClip`) and switched row-owning component render paths to the shared helpers.
- Renamed `registry/components` to `registry/bricks`, updated imports/CLI copy paths, and aligned docs to the bricks -> rooms -> bentos composition model.

## [0.3.3] - 2026-03-13

### Changed

- Reworked `layouts.Focus` to a native body+footer grammar so footer-anchored screens no longer reserve hidden header rows
- Reworked `layouts.Pancake` to a native header+body+footer grammar and added snapshot coverage to prevent hidden shim-row regressions
- Migrated starter app, shipped bentos, and `bento init` scaffold output to `layouts.Focus(...)` for footer-only chrome by default
- Updated README/docs to document Focus-first home/starter layout usage and added the home-screen demo GIF

## [0.3.2] - 2026-03-13

### Changed

- Removed test/demo header banner content from shipped bentos, starter app, and `bento init` scaffold output
- Kept Frame row structure while defaulting top/subheader rows to minimal app chrome
- Updated docs set under `docs/` to reflect the new minimal header defaults and frame hierarchy state

## [0.3.1] - 2026-03-13

### Changed

- Finalized Frame-first layout grammar and role-aware bar hierarchy (`top`, `subheader`, `footer`)
- Added single muted status pill pattern for header metadata and anchored footer solid-row rendering
- Expanded list/table/bar APIs and starter/bento defaults to the new frame/chrome contract
- Refined theme role mapping/validation and updated docs for the latest layout + styling model

## [0.3.0] - 2026-03-12

### Added

- New registry components: `badge`, `kbd`, `wordmark`, `select`, `checkbox`,
  `progress`, `tabs`, `toast`, `separator`
- New runnable bento examples in `registry/bentos/`: `home-screen`,
  `app-shell`, `dashboard`
- New `registry/layouts/` package with 15 named layout functions and shared sizing engine
- New layout docs at `docs/layouts.md`

### Changed

- `bento add` / `bento list` now include the new component catalog entries
- Release builds now inject CLI version correctly via GoReleaser `-X main.version`
- Component roadmap/docs now treat the component catalog as finalized and shift
  execution focus toward shipping more `registry/bentos/` examples
- Removed spinner-from-registry planning; `spinner` is now a direct
  `charm.land/bubbles/v2/spinner` primitive in app code instead
- Docs now use `registry/bentos/` as the source-of-truth bento path and mark
  the first wave as shipped
- `bento doctor` now checks all shipped copy-and-own components
- README Go badge now matches module Go version floor (`1.25+`)
- Starter app, registry bentos, and `bento init` template now compose screens with `registry/layouts`
- Full-frame rendering contract restored: `surface.Fill(...)` + `surface.Draw(...)` remains the final compositor path
- Docs updated to reflect `registry/layouts` + `surface` composition responsibilities

## [0.2.0] - 2026-03-02

### Changed

- **Breaking — registry model replaces the framework API.**
  `bentotui.New()`, `core/`, `app/`, `ui/`, and `bentotui.go` are deleted.
  BentoTUI is now a registry of copy-and-own components, not an opinionated shell.

- `core/theme/` moved to `theme/` — import path is now `github.com/cloudboy-jh/bentotui/theme`
- `ui/styles/` moved to `styles/` — import path is now `github.com/cloudboy-jh/bentotui/styles`
- `core/layout/` moved to `layout/` — import path is now `github.com/cloudboy-jh/bentotui/layout`

- `theme.CurrentTheme()` is now goroutine-safe (`sync.RWMutex` on every read/write).
  `theme.SetTheme(name)` is safe to call from `main()` before `tea.NewProgram().Run()`.

- All registry components read `theme.CurrentTheme()` in `View()` — no stored theme
  state, no `SetTheme()` propagation chains.

- Every rendered row uses a single `lipgloss.NewStyle().Background().Width().Render(plain)`
  call. ANSI is stripped before rendering. Eliminates background color bleed-through
  where inner escape codes previously fought outer background cells.

- `layout/` exports a conservative public surface:
  `Horizontal`, `Vertical`, `Fixed`, `Flex`, `Split`, `Item`, `Model`, `Sizeable`.
  `Kind` is unexported. `ViewString` and `Fill` are inlined — no `core` dep.

- `bar` no longer stores `m.theme` or handles `focus.FocusChangedMsg`. Row and card
  rendering are inlined. `theme.CurrentTheme()` called in `View()`.

- `dialog.Custom` wraps any `tea.Model` as content — no `panel` import required.
  `dialog.Manager` no longer stores or propagates theme.

- `input.View()` calls `styles.New(theme.CurrentTheme()).InputStyles()` on every
  render so live theme switching works without any `SetTheme()` call.

### Added

- `registry/` — clean rewrites of all components:
  `panel`, `bar`, `dialog` (manager + confirm + custom + theme_picker + command_palette),
  `list`, `table`, `text`, `input`

- `theme/theme_test.go` — ported all tests from `core/theme/`, added goroutine-safety
  smoke test (`TestCurrentThemeConcurrentAccess`)

- `docs/next-steps.md` — 3 concrete known gaps: `bento add` embed wiring,
  `bento init` template, `input.View()` style caching

- `.gitignore`

### Removed

- `app/` — monolithic app shell
- `core/` — shell, router, focus, palette, interfaces, messages, view helpers
- `ui/` — containers, widgets, primitives, styles (all replaced by `registry/`)
- `bentotui.go` — public API facade

### Fixed

- Background bleed-through on panel content rows (canvas Z-layer approach replaced
  with single lipgloss render call)
- Theme picker selection highlight now fills full row width
- `Custom.View()` and `Confirm.View()` now use `theme.CurrentTheme()` instead of
  a stale cached value — dialog frame and content colors stay in sync during live preview

## [0.1.0-initial] - 2026-02-23

### Added

- initial BentoTUI framework foundation
- core app shell, router, layout, focus, and theme modules
- early dialog/footer/panel component set
- canvas-based layout system (Horizontal/Vertical with Fixed/Flex/Min/Max constraints)
- global theme system with 15 professional presets via bubbletint
- complete widget library (Input, List, Text, Card, Table)
- container components (Panel, Bar, Dialog with theme picker)
- reactive theme propagation across all components
