# Changelog

All notable changes to this project will be documented in this file.

The format follows Keep a Changelog style and this project targets Semantic Versioning.

## [Unreleased]

## [0.5.1] - 2026-03-18

### Fixed

- `bento doctor` now checks copied bricks from `bricks/<name>` and derives the
  checklist from the live brick registry instead of stale hardcoded entries.
- CLI/TUI wording is aligned to Bento product language (`bricks` instead of
  `components`) across help text, menus, install logs, and list output.

### Added

- Added `docs/product-direction.md` to lock scope discipline and define when
  custom app-owned bricks are expected.

### Changed

- Updated usage and roadmap docs to clarify that Bento is not a maximal
  component catalog and to document the custom-brick decision path.

## [0.5.0] - 2026-03-18

### Added

- Added high-level room contracts in `registry/rooms`:
  `AppShell`, `SidebarDetail`, `Dashboard`, and `DiffWorkspace`.
- Added guardrail policy tests in `internal/policy/guardrails_test.go` and
  CI coverage in `.github/workflows/ci.yml`.
- Added `docs/usage-guide.md` and `docs/astro-content.md` to lock Bento-first
  positioning and website messaging.

### Changed

- Reframed docs and READMEs around a product workflow:
  bricks (official components), rooms (page layouts), bentos (template apps).
- Updated starter/scaffold copy to guide users toward rooms + bricks page
  composition instead of low-level component exploration.

## [0.4.1] - 2026-03-18

### Added

- **Diff and Syntax tokens on `Theme` interface.** Added 10 diff color methods and
  9 syntax highlight methods to support `bento-diffs` and any diff rendering surface.
  All methods have nil-guarded defaults in `BaseTheme` derived from existing tokens
  (`SuccessColor`, `ErrorColor`, `TextMutedColor`, etc.) — no preset changes required.

  Diff methods: `DiffAddedBG`, `DiffRemovedBG`, `DiffContextBG`, `DiffAddedLineNumBG`,
  `DiffRemovedLineNumBG`, `DiffAdded`, `DiffRemoved`, `DiffLineNum`,
  `DiffHighlightAdded`, `DiffHighlightRemoved`

  Syntax methods (for chroma integration): `SyntaxKeyword`, `SyntaxType`,
  `SyntaxFunction`, `SyntaxVariable`, `SyntaxString`, `SyntaxNumber`,
  `SyntaxComment`, `SyntaxOperator`, `SyntaxPunctuation`

## [0.4.0] - 2026-03-18

### Breaking

- **`Theme` is now a Go interface.** The old struct-based token accessor pattern
  (`t.Surface.Canvas`, `t.Text.Primary`, `t.Selection.BG`, etc.) is removed.
  Components now call interface methods: `t.Background()`, `t.Text()`, `t.SelectionBG()`.
  Any copied brick using the old dot-accessor syntax must be updated.

- **`panel` and `elevated-card` merged into `card`.** The `panel` brick is deleted.
  `elevated-card` is deleted. Both are replaced by a single `card` brick with a
  `Flat()` option for the former panel behavior. `bento add panel` and
  `bento add elevated-card` no longer exist — use `bento add card`.

- **`Frame`, `FrameMainDrawer`, `FrameTriple` rooms removed.** These were thin wrappers
  around `JoinVertical` with named slots. Use `Pancake`, `TopbarPancake`, or `Focus` instead.

- **`theme/styles/System` struct deleted.** `styles.New(t)` and all methods on it
  (`DialogFrame`, `InputColors`, `PaletteItem`, `FooterCardCommand`, etc.) are removed.
  The `styles` package now exports only `Row`, `RowClip`, and `ClipANSI` as package-level
  functions. Any copied brick using `styles.New(t).X()` must be updated.

### Changed

- **Colors-in architecture.** Every brick now accepts `WithTheme(t theme.Theme)` at
  construction and `SetTheme(t theme.Theme)` for live updates. Bricks fall back to
  `theme.CurrentTheme()` only when no explicit theme was provided. No brick calls
  `CurrentTheme()` unconditionally in `View()`.

- **Theme presets are plain Go structs.** The `bubbletint` runtime adapter and all
  contrast validation/quality scoring logic is removed. 16 presets live as hardcoded
  `BaseTheme` structs in `theme/presets.go`. `theme.Preset("name")` returns a value
  with zero global side effects — safe to use in CLI tools, tests, and non-TUI contexts.

- **Global theme manager is now opt-in.** `theme.SetTheme`, `theme.CurrentTheme`,
  `theme.PreviewTheme`, `theme.RegisterTheme`, and `theme.AvailableThemes` all remain
  for apps that want a single active theme. Bricks no longer require it.

- **Bentos hold theme as app state.** `m.theme theme.Theme` is a field on the app model.
  `ThemeChangedMsg` is handled by calling `SetTheme` on each brick. No framework magic —
  the app decides which bricks get the new theme and when.

- **Theme picker UX matches opencode.** Navigate up/down to live-preview themes.
  Enter to confirm. Esc reverts to the theme active when the picker was opened.

- **`docs/astro-content-update.md` deleted.** Stale, no longer relevant.

- Reworked `registry/bricks/list` delegate-driven rendering, explicit focus routing,
  predictable keyboard navigation and `tea.WindowSizeMsg` handling.
- Reworked `registry/bricks/table` with cleaner focus/blur semantics, `VisualClean`
  and `VisualGrid` presets, column priority shrinking.
- Reworked `registry/bricks/filepicker` to align with upstream `DidSelectFile` flow.
- All docs updated to reflect v0.4.0 architecture.
- Root README updated — "shadcn of TUIs" framing, colors-in model, v0.4.0 API examples.

### Added

- `theme.BaseTheme` — embeddable struct for custom theme implementations.
- `theme.Preset(name) Theme` and `theme.Names() []string` — preset access with no global state.
- `card.Flat()` option — flat titled container (former `panel` behavior).
- `card.WithTheme(t)` + `card.SetTheme(t)` on every brick.
- `registry/bricks/package-manager` — spinner + progress sequential install flow.
- `registry/bentos/dashboard-brick-lab` — component showcase bento.

### Removed

- `theme/adapter.go` — bubbletint adapter deleted; presets are hardcoded structs.
- `theme/quality.go` — contrast validation and quality scoring deleted.
- `theme/storage.go` — disk persistence deleted (app concern).
- `theme/bento_rose.go` — merged into `theme/presets.go`.
- `theme/messages.go` (old) — `OpenThemePickerMsg` removed; `ThemeChangedMsg` kept.
- `registry/bricks/panel/` — merged into `registry/bricks/card/`.
- `registry/bricks/elevated-card/` — merged into `registry/bricks/card/`.
- `registry/rooms/frame.go` — `Frame`, `FrameMainDrawer`, `FrameTriple` removed.
- `theme/styles/System` struct and all methods on it.
- `docs/theme-escape-hatch.md` — superseded by the v0.4.0 architecture.
- `docs/astro-content-update.md` — stale.

## [0.3.5] - 2026-03-17

### Changed

- Moved shared style helpers from `styles/` to `theme/styles/` and updated all imports to `github.com/cloudboy-jh/bentotui/theme/styles`.
- Reworked `registry/bentos/app-shell` into a single-screen composition bento with rail + workspace + anchored footer + command palette.
- Retired scenario-driven app-shell runtime paths and removed scenario harness logic from default app-shell behavior.
- Updated docs and README to formalize the **Untouchable Theme Engine** model and lock the architecture language to `bentos + rooms + bricks`.
- Added room-level separation options (`WithGutter`, `WithDivider`) to split/drawer room primitives.
- Added anchored footer card style modes in bar (`plain`, `chip`, `mixed`) and aligned default bento usage around anchored command lanes.
- Expanded list row structure with typed fields (`Primary`, `Secondary`, `RightStat`, `Tone`, `SelectedStyle`) while preserving existing `Label`/`Status`/`Stat` compatibility.
- Added new `elevated-card` brick for raised section containers (`Title` + `Content`) and wired it into dashboard/app-shell composition paths.
- Migrated `list`, `table`, `progress`, `select`, `checkbox`, and `tabs` bricks to Charm-backed internals (`bubbles/*`) while keeping Bento wrapper APIs.
- Added `filepicker` brick backed by `charm.land/bubbles/v2/filepicker` and exposed it via `bento add filepicker`.

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
