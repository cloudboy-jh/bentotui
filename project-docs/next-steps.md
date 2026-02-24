# BentoTUI Next Steps

Status: Active
Date: 2026-02-24

This plan follows ADR-0001 in `project-docs/rendering-system-design.md`.

## Phase A: Paint Baseline (Color Blocking)

- Ensure every shell/layout/component slot paints a full rectangular area each frame
- Remove rendering paths that flatten child ANSI styling into a single foreground pass
- Verify no transparent seams remain during route switch, resize, or overlay open/close

Exit criteria:

- no visible bleed-through of previous terminal buffer content
- panel/dialog/status surfaces read as solid blocks, not wireframes

## Phase B: Theme and Surface Language (Primary)

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

## Phase C: Focus Visibility Pass

- Increase focused title/border delta while keeping non-focused panels quiet
- Validate focus cues in both wide and compact layouts
- Add tests that assert focus-cycling causes visible state changes

Exit criteria:

- no ambiguity about focused target after a single `tab` press

## Phase D: Lock Renderer Correctness

- Ensure root renderer remains update-safe (`Update` mutates state, `View` renders state)
- Keep deterministic z-order in app shell (`shell -> body -> status -> scrim -> dialog`)
- Verify split allocations always paint full slot areas
- Add resize stress tests for split layouts and viewport changes

Exit criteria:

- no right/bottom unpainted bands in Windows Terminal with blur/transparency on
- no startup ghost frame
- `go test ./...` and `go vet ./...` clean

## Phase E: Componentize Surface Rendering

- Expand `surface` package as shared primitive layer:
  - fill blocks
  - fixed-size regions
  - width fitting and clipping helpers
- Move panel frame/title/body rendering into composable internal units
- Unify dialog frame rendering through shared modal frame helpers
- Standardize statusbar segment layout with shared fit policies

Exit criteria:

- no duplicated low-level width/fill helpers across modules
- panel/dialog/status all consume shared surface utilities

## Phase F: Docs and Examples

- Add a dedicated "Build a real app shell" guide to main spec
- Add examples for:
  - multi-page routing
  - focus manager integration
  - dialog flows
  - fullscreen vs inline mode
- Keep README concise and professional; push deep detail to project docs

Exit criteria:

- README stays short and stable
- project docs cover operational details and architecture contracts

## Milestone Work Items (Near-Term)

1. Add renderer regression tests for full-frame paint guarantees.
2. Add layout slot paint tests for fixed/flex remainder behavior.
3. Add statusbar truncation tests for narrow widths.
4. Finalize default theme set (`catppuccin-mocha`, `dracula`, `osaka-jade`) and startup/persistence policy.
5. Keep `cmd/test-tui` as the canonical color/focus validation harness.
