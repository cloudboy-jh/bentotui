# BentoTUI Next Steps

Status: Active
Date: 2026-02-25

This list is intentionally component-focused and execution-oriented.
Use this with `project-docs/component-system-reference.md`.

## 1. Footer Statusline Contract

- [x] Keep footer as one continuous strip (height = 1 always).
- [x] Enforce segment ownership: `left | actions | right`.
- [x] Remove freeform footer hint strings as primary API.

## 2. Footer Action Model

- [x] Add structured action schema: `key`, `label`, `variant`, `enabled`.
- [x] Add explicit footer API for action arrays.
- [x] Render actions as first-class inline chips/buttons.

## 3. Footer Truncation Rules

- [x] Lock deterministic truncation priority: `right > left > actions`.
- [x] Add action collapse policy: full chip -> key only -> drop from end.
- [x] Guarantee no wrapping in narrow widths.

## 4. Focus Manager API Hardening (`core/focus/focus.go`)

- [ ] Add explicit APIs: `SetRing`, `SetIndex`, `FocusBy`.
- [ ] Add `enabled` and `wrap` controls.
- [ ] Handle empty/nil ring entries safely.

## 5. Focus Event Contract

- [ ] Add `FocusChangedMsg {from, to}`.
- [ ] Emit focus changes deterministically from manager.
- [ ] Wire footer to consume focus change state.

## 6. Shared UI Primitives

- [x] Extract footer action chip primitive.
- [x] Extract reusable list row primitive.
- [x] Extract reusable input surface primitive.
- [x] Keep modal frame primitive shared and bounded.

## 7. Theme Picker UX Refinement

- [ ] Add hover preview on selection move (no commit yet).
- [ ] Commit on `enter`, revert preview on `esc`.
- [ ] Keep picker bounded and clipping-safe in all terminal sizes.

## 8. Command Palette Component

- [ ] Add command palette dialog component.
- [ ] Route `/` to command-entry/palette workflow (as finalized behavior).
- [ ] Keep `/theme`, `/dialog`, `/page` command paths consistent.

## 9. Component Regression Coverage

- [x] Footer tests: layout/order/truncation/no-wrap.
- [ ] Focus tests: ring updates/wrap/index/events.
- [ ] Dialog tests: custom enter routing + bounds stability.
- [ ] Theme picker tests: preview/apply/revert behavior.

## 10. Harness + Docs Sync

- [x] Update `cmd/test-tui` to consume structured footer actions.
- [x] Keep `project-docs/component-system-reference.md` aligned.
- [x] Keep `project-docs/framework-roadmap.md` and `CHANGELOG.md` updated.
