# BentoTUI Next Steps

Status: Active
Date: 2026-02-25

This list is intentionally component-focused and execution-oriented.
Use this with `project-docs/component-system-reference.md`.

## 1. Footer Statusline Contract

- [x] Keep footer as one continuous strip (height = 1 always).
- [x] Enforce segment ownership: `left | cards | right`.
- [x] Remove freeform footer hint strings as primary API.

## 2. Footer Card Model

- [x] Add structured card schema: `command`, `label`, `variant`, `enabled`.
- [x] Add explicit footer API for card arrays.
- [x] Render cards as first-class inline UI surfaces.

## 3. Footer Truncation Rules

- [x] Lock deterministic truncation priority: `right > left > cards`.
- [x] Add card collapse policy: full card -> command only -> drop from end.
- [x] Guarantee no wrapping in narrow widths.

## 4. Focus Manager API Hardening (`core/focus/focus.go`)

- [x] Add explicit APIs: `SetRing`, `SetIndex`, `FocusBy`.
- [x] Add `enabled` and `wrap` controls.
- [x] Handle empty/nil ring entries safely.

## 5. Focus Event Contract

- [x] Add `FocusChangedMsg {from, to}`.
- [x] Emit focus changes deterministically from manager.
- [x] Wire footer to consume focus change state.

## 6. Shared UI Primitives

- [x] Extract reusable footer card primitive.
- [x] Extract reusable list row primitive.
- [x] Extract reusable input surface primitive.
- [x] Keep modal frame primitive shared and bounded.

## 7. Theme Picker UX Refinement

- [x] Add hover preview on selection move (no commit yet).
- [x] Commit on `enter`, revert preview on `esc`.
- [x] Keep picker bounded and clipping-safe in all terminal sizes.

## 8. Command Palette Component

- [ ] Add command palette dialog component.
- [ ] Route `/` to command-entry/palette workflow (as finalized behavior).
- [ ] Keep slash command aliases coherent (`/pr`, `/issue`, `/branch` + legacy aliases).

## 9. Component Regression Coverage

- [x] Footer tests: layout/order/truncation/no-wrap.
- [x] Focus tests: ring updates/wrap/index/events.
- [ ] Dialog tests: custom enter routing + bounds stability.
- [x] Theme picker tests: preview/apply/revert behavior.

## 10. Harness + Docs Sync

- [x] Update `cmd/starter-app` to consume structured footer cards.
- [x] Keep `project-docs/component-system-reference.md` aligned.
- [x] Keep `project-docs/framework-roadmap.md` and `CHANGELOG.md` updated.

## 11. Footer Design Language

- [x] Rename footer model from action/hotkey semantics to card semantics.
- [x] Rename shared primitive from legacy naming to card language.
- [x] Align harness footer copy to command-first GitHub-style card examples.
