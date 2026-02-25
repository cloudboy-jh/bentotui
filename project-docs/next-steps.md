# BentoTUI Next Steps

Status: Active
Date: 2026-02-24

This plan follows ADR-0001 in `project-docs/rendering-system-design.md`, reflects the current repository structure (`ui/components/*`, `ui/styles`, footer-first shell API), and should be executed alongside `project-docs/component-system-reference.md`.

## Recently Completed

- [x] Moved UI layer to `ui/components/*` and `ui/styles`.
- [x] Promoted footer-first shell API (`WithFooterBar`, `WithFooter`).
- [x] Finalized built-in theme IDs and default behavior:
  - [x] `catppuccin-mocha` (default)
  - [x] `dracula`
  - [x] `osaka-jade`
- [x] Added `CHANGELOG.md` and adopted changelog tracking.
- [x] Fixed harness text-input hotkey conflict (`d`/`q` can be typed when input is focused).

## Current Sprint (Priority)

## Phase 7A: Theme Dialog Bounds and Layout

Scope: fix the current theme dialog rendering artifacts and enforce modal bounds.

- [x] Clamp theme picker input/list rows to dialog content width.
- [x] Prevent dialog content from stretching to viewport width.
- [x] Ensure modal remains centered with stable dimensions on wide terminals.
- [x] Remove duplicate/stacked visual artifacts in search + list areas.

Exit criteria:

- [x] no horizontal overflow or stretched theme modal frames
- [x] list and search surfaces render once, cleanly
- [x] dialog remains centered and visually stable under resize

## Phase 7B: Theme Apply Routing Hardening

Scope: ensure theme switching is deterministic from dialog interaction.

- [x] Guarantee `enter` applies the selected theme before dialog closes.
- [x] Ensure theme updates propagate to shell, footer, dialogs, and active page.
- [x] Confirm selected theme persists and restores correctly on restart.

Exit criteria:

- [x] theme apply on `enter` is reliable
- [x] no race between close/apply paths
- [x] persisted theme loads correctly on relaunch

## Phase 7C: Footer Contract Standardization

Scope: formalize footer behavior and truncation rules.

- [ ] Lock segment roles: `left`, `help`, `right`.
- [ ] Define deterministic truncation priority for narrow widths.
- [ ] Stabilize rendering for width=0, tiny widths, and resize transitions.
- [ ] Add structured footer actions (button/chip model) instead of plain hint strings.
- [ ] Define footer action schema (`key`, `label`, `variant`, `enabled`).
- [ ] Add footer API for explicit action lists and phased migration from ad hoc text.

Exit criteria:

- [ ] footer stays readable and deterministic at constrained widths
- [ ] no segment ordering drift during resize
- [ ] footer actions render as first-class UI elements, not freeform text

## Next Sprint

## Phase 8: Component Standardization Contract

Scope: make component behavior consistent across the UI layer.

- [ ] Standardize `SetSize` ownership and bounds clipping across components.
- [ ] Enforce keybinding precedence rules by focus context.
- [ ] Ensure focused state is visually explicit and consistent.
- [ ] Keep components free of ad-hoc color literals (style-layer only).

Exit criteria:

- [ ] panel/dialog/footer/picker follow one sizing and focus contract
- [ ] no component can render outside assigned bounds

## Phase 9: Shared UI Primitives

Scope: remove duplicated rendering patterns across components.

- [ ] Extract shared modal frame primitive.
- [ ] Extract shared input surface primitive.
- [ ] Extract shared selectable row/list primitive.
- [ ] Extract shared footer action chip primitive.
- [ ] Keep component-specific behavior, share visual frame mechanics.

Exit criteria:

- [ ] reduced duplication in panel/dialog/picker/footer rendering
- [ ] consistent spacing and hierarchy across components

## Stability and Verification

## Phase 10: Regression Coverage

- [ ] Add renderer regression tests for full-frame paint guarantees.
- [ ] Add layout slot paint tests for fixed/flex remainder behavior.
- [ ] Add theme dialog bound/overflow tests.
- [ ] Add footer truncation and segment-priority tests.
- [ ] Add footer action rendering/truncation tests.
- [ ] Add dialog apply/close routing tests.

Exit criteria:

- [ ] `go test ./...` and `go vet ./...` clean
- [ ] no known dialog overflow or footer truncation regressions

## Docs and Examples

## Phase 11: Spec and Example Sync

- [ ] Add a dedicated "Build a real app shell" guide to main spec.
- [ ] Add examples for:
  - [ ] multi-page routing
  - [ ] focus manager integration
  - [ ] dialog flows
  - [ ] fullscreen vs inline mode
- [ ] Keep README concise; keep deeper architecture detail in `project-docs`.
- [ ] Update `CHANGELOG.md` continuously as phases land.

Exit criteria:

- [ ] README remains high-signal and current
- [ ] spec/doc examples match real package paths and APIs

## Deferred

## Slash Command Palette (`/`)

- [ ] Full command palette UX remains deferred.
- [ ] Current slash behavior stays explicit (`/theme`, `/dialog`, `/confirm`) until palette design is implemented.
