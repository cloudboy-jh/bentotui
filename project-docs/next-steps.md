# BentoTUI Next Steps

Status: Active
Date: 2026-02-24

This plan follows ADR-0001 in `project-docs/rendering-system-design.md`.

## Phase 1: Lock Renderer Correctness

- Ensure root renderer remains update-safe (`Update` mutates state, `View` renders state)
- Keep deterministic z-order in app shell (`shell -> body -> status -> scrim -> dialog`)
- Verify split allocations always paint full slot areas
- Add resize stress tests for split layouts and viewport changes

Exit criteria:

- no right/bottom unpainted bands in Windows Terminal with blur/transparency on
- no startup ghost frame
- `go test ./...` and `go vet ./...` clean

## Phase 2: Componentize Surface Rendering

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

## Phase 3: Restore Solid UI Language

- Replace wireframe-heavy borders with tonal surfaces and deliberate hierarchy
- Keep focus accents high-contrast while reducing non-focused border noise
- Normalize spacing rhythm across header, panels, statusbar, and dialogs
- Preserve clear readability on both dark and transparent terminals

Exit criteria:

- harness (`cmd/test-tui`) has clear visual depth and consistent component language
- no regressions in interaction (`1/2/o/q`, tab focus cycle, dialog close)

## Phase 4: Docs and Examples

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
4. Tune panel theme tokens for non-wireframe solid surfaces.
5. Build one polished `cmd/test-tui` theme profile for demos.
