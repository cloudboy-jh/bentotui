# BentoTUI Next Steps

Status: Active
Date: 2026-02-24

This plan follows ADR-0001 in `project-docs/rendering-system-design.md`.

## Recent Progress

- 2026-02-24: README was rewritten to align package layout and UI component framing (`567cda6`).

## Phase 1: Paint Baseline (Color Blocking)

- Ensure every shell/layout/component slot paints a full rectangular area each frame
- Remove rendering paths that flatten child ANSI styling into a single foreground pass
- Verify no transparent seams remain during route switch, resize, or overlay open/close

Exit criteria:

- no visible bleed-through of previous terminal buffer content
- panel/dialog/footer surfaces read as solid blocks, not wireframes

## Phase 2: Theme and Surface Language (Primary)

- Finalize v0.1 built-in presets and startup default policy:
  - `catppuccin-mocha` (default)
  - `dracula`
  - `osaka-jade`
- Tune semantic token ladder (`Background`, `Surface`, `SurfaceMuted`) for clear depth
- Strengthen component token contrast (`Border`, `BorderFocused`, `TitleBG`, `StatusBG`, `DialogBG`)

Exit criteria:

- clear visual hierarchy in harness at a glance
- focused panels are unmistakable during tab cycling
- theme switching and restart persistence stable

## Phase 3: Focus Visibility Pass

- Increase focused title/border delta while keeping non-focused panels quiet
- Validate focus cues in both wide and compact layouts
- Add tests that assert focus-cycling causes visible state changes

Exit criteria:

- no ambiguity about focused target after a single `tab` press

## Phase 4: Lock Renderer Correctness

- Ensure root renderer remains update-safe (`Update` mutates state, `View` renders state)
- Keep deterministic z-order in app shell (`shell -> body -> footer -> scrim -> dialog`)
- Verify split allocations always paint full slot areas
- Add resize stress tests for split layouts and viewport changes

Exit criteria:

- no right/bottom unpainted bands in Windows Terminal with blur/transparency on
- no startup ghost frame
- `go test ./...` and `go vet ./...` clean

## Phase 5: Componentize Surface Rendering

- Expand `surface` package as shared primitive layer:
  - fill blocks
  - fixed-size regions
  - width fitting and clipping helpers
- Move panel frame/title/body rendering into composable internal units
- Unify dialog frame rendering through shared modal frame helpers
- Standardize footer segment layout with shared fit policies

Exit criteria:

- no duplicated low-level width/fill helpers across modules
- panel/dialog/footer all consume shared surface utilities

## Phase 6: Docs and Examples

- Add a dedicated "Build a real app shell" guide to main spec
- Add examples for:
  - multi-page routing
  - focus manager integration
  - dialog flows
  - fullscreen vs inline mode
- Keep README concise and professional; push deep detail to project docs (done)

Exit criteria:

- README stays short and stable
- project docs cover operational details and architecture contracts

## Phase 7: Dialog UX and Footer Contract (Next)

Scope: stabilize theme application flow, clean theme picker rendering, and formalize footer behavior.
Out of scope: slash command palette (`/`) and global command execution UX.

- Fix dialog event routing so picker dialogs receive `enter` and can apply selection before close
- Ensure theme apply from Theme Picker updates shell, footer, dialog manager, and active page reliably
- Clean Theme Picker layout:
  - remove duplicate search surface artifacts
  - normalize spacing and row widths
  - tighten keyboard hint line and current-selection marker behavior
- Formalize `ui/components/footer/footer.go` contract:
  - explicit left/help/right segment roles
  - deterministic truncation priority under narrow widths
  - stable rendering behavior for width=0, tiny widths, and resize transitions
- Add focused tests for:
  - theme apply-on-select behavior from dialog
  - dialog manager enter/esc behavior by dialog type
  - footer truncation and segment ordering at constrained widths

Exit criteria:

- selecting a theme in dialog applies immediately and persists across restart
- no visual duplication/glitch in theme picker search/list surfaces
- footer remains readable and deterministic across terminal widths
- `go test ./...` and `go vet ./...` clean

## Deferred: Slash Command Palette

- `/` as a command index/palette is deferred to a later phase
- current slash behavior remains explicit command handling (`/theme`, `/dialog`, `/confirm`) until palette design is implemented

## Milestone Work Items (Near-Term)

1. Add renderer regression tests for full-frame paint guarantees.
2. Add layout slot paint tests for fixed/flex remainder behavior.
3. Add footer truncation tests for narrow widths.
4. Finalize default theme set (`catppuccin-mocha`, `dracula`, `osaka-jade`) and startup/persistence policy.
5. Keep `cmd/test-tui` as the canonical color/focus validation harness.
6. Add the "Build a real app shell" guide to `project-docs/bentotui-main-spec.md`.
7. Implement Phase 7 dialog event routing and theme-apply flow hardening.
8. Add Theme Picker and footer rendering cleanup plus regression tests.
