# BentoTUI Next Steps

Status: Active
Date: 2026-02-24

This plan follows ADR-0001 in `project-docs/rendering-system-design.md` and `project-docs/component-system-reference.md`.

## Active Work

1. Footer action model (replace plain hint strings)
- [ ] Add structured footer actions (`key`, `label`, `variant`, `enabled`).
- [ ] Add footer API for explicit action lists.
- [ ] Render action chips/buttons in footer (not freeform text).

2. Footer layout and truncation contract
- [ ] Lock segment roles: `left`, `actions`, `right`.
- [ ] Define deterministic truncation priority at narrow widths.
- [ ] Stabilize behavior for width=0 and resize transitions.

3. Focus manager hardening (`focus/focus.go`)
- [ ] Add explicit ring/index control APIs (`SetRing`, `SetIndex`, `FocusBy`).
- [ ] Add enabled/wrap behavior and safe nil handling.
- [ ] Emit deterministic focus-changed events/messages.

4. Component sizing contract enforcement
- [ ] Ensure all UI components use `SetSize`/`GetSize` consistently.
- [ ] Enforce clipping within assigned bounds.
- [ ] Remove viewport-coupled rendering logic from components.

5. Shared UI primitives
- [ ] Extract shared modal frame primitive.
- [ ] Extract shared input surface primitive.
- [ ] Extract shared list/row selection primitive.
- [ ] Extract shared footer action chip primitive.

6. Visual system normalization
- [ ] Keep solid card surfaces with no border chrome.
- [ ] Keep only subtle section separators where necessary.
- [ ] Ensure focus visibility uses row/text contrast, not border outlines.

7. Dialog and picker regression coverage
- [ ] Add/maintain tests for modal bounds on wide/narrow terminals.
- [ ] Add tests for `enter` apply/close ordering.
- [ ] Add tests for custom-dialog key routing (`enter` not swallowed).

8. Footer regression coverage
- [ ] Add tests for action rendering order.
- [ ] Add tests for truncation behavior and tiny-width fallback.
- [ ] Add tests for disabled/muted action rendering.

9. Harness refinement (`cmd/test-tui`)
- [ ] Keep command-entry flow (`/theme`, `/dialog`, `/confirm`) on Enter.
- [ ] Keep focus behavior deterministic between input/actions.
- [ ] Use harness as canonical visual regression surface across themes.

10. Docs and examples sync
- [ ] Add "Build a real app shell" guide to main spec.
- [ ] Add examples for routing, focus, dialog flows, fullscreen vs inline.
- [ ] Keep README concise; move implementation detail to `project-docs`.
- [ ] Keep `CHANGELOG.md` updated as milestones land.
