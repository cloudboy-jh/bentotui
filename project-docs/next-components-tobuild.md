# BentoTUI Next Components To Build

Status: Active
Date: 2026-02-26

Purpose: keep the 1.0 component scope locked and execution-focused.

## 1.0 Locked Set (13)

- [ ] Header (`ui/components/header`) - top bar mirror of footer contract.
- [x] Footer (`ui/components/footer`) - one-row command card bar with deterministic truncation.
- [x] Panel (`ui/components/panel`) - bounded frame container with focus state.
- [x] Dialog Manager (`ui/components/dialog`) - confirm/custom lifecycle and bounded overlay behavior.
- [x] Theme Picker (`ui/components/dialog/theme_picker`) - preview/apply/revert flow.
- [ ] Command Palette (`ui/components/dialog/command_palette`) - slash-first command entry + search.
- [ ] Input Surface (`ui/components/input`) - bounded single-line input wrapper.
- [ ] List (`ui/components/list`) - selectable rows with clipping and paging contract.
- [ ] Table (`ui/components/table`) - bounded columns and deterministic truncation.
- [ ] Tabs (`ui/components/tabs`) - segmented view switcher.
- [ ] Status Banner (`ui/components/status`) - inline info/success/warn/error strip.
- [ ] Toast (`ui/components/toast`) - ephemeral message queue.
- [ ] Progress (`ui/components/progress`) - progress bar/spinner states.

## Acceptance Contract (Each Component)

- [ ] Implements `tea.Model` and `SetSize/GetSize` bounds contract.
- [ ] Uses semantic styles from `ui/styles` only.
- [ ] Maintains deterministic clipping/no-wrap rules under narrow widths.
- [ ] Follows key routing precedence (dialog -> focused component -> shell).
- [ ] Has regression tests for behavior + bounds correctness.
- [ ] Is exercised in `cmd/starter-app` with a visible usage path.

## Release Focus

- P0 for v1.0: all 13 locked components implemented and regression-covered.
- P1 after v1.0: advanced data views, tree controls, and code views.
