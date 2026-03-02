# Bento Color Design System

Status: Active reference
Date: 2026-03-01

This document defines BentoTUI color system rules and required theme tokens.

## Design Principles

- Theme values come from documented token definitions only.
- No token derivation, blending, or fallback synthesis is allowed.
- Components must consume semantic roles via `ui/styles`.
- New themes must be complete and valid before registration.

## Required Token Groups

Every theme must define all tokens below.

### surface

- `surface.canvas`
- `surface.panel`
- `surface.elevated`
- `surface.overlay`
- `surface.interactive`

### text

- `text.primary`
- `text.muted`
- `text.inverse`
- `text.accent`

### border

- `border.normal`
- `border.subtle`
- `border.focus`

### state

- `state.info`
- `state.success`
- `state.warning`
- `state.danger`

### selection

- `selection.bg`
- `selection.fg`

### input

- `input.bg`
- `input.fg`
- `input.placeholder`
- `input.cursor`
- `input.border`

### bar

- `bar.bg`
- `bar.fg`

### dialog

- `dialog.bg`
- `dialog.fg`
- `dialog.border`
- `dialog.scrim`

## Hierarchy Rules

- `surface.canvas` is the shell baseline.
- `surface.panel` and `surface.elevated` define normal depth stacking.
- `surface.overlay` is reserved for modal/detached surfaces.
- `surface.interactive` is used for active or editable regions.
- `text.inverse` is used only where foreground sits on high-contrast accent/state surfaces.

## Stitched Surface Ladder

All composed pages must preserve this order:

- `surface.canvas` -> `surface.panel` -> `surface.elevated` -> `surface.interactive` -> `dialog.bg`

Rules:

- Inputs always use `input.bg`, not parent panel surface.
- Selected rows always use `selection.bg`/`selection.fg`.
- Footer/base bars always use `bar.bg` with card overlays on top.
- Primary action cards must not share the same background as `bar.bg`.

## Component Mapping Rules

- Header/footer bars: `bar.*`
- Panel body/title strips: `surface.*`, `text.*`, `border.*`
- List selection and active controls: `selection.*`
- Input rows: `input.*`
- Dialog frame/scrim: `dialog.*`
- Error/warning/success/attention accents: `state.*`

## Validation Rules

- Theme registry rejects any theme with missing required tokens.
- Theme registry rejects empty token values.
- `AvailableThemes()` returns only validated themes.
- Unknown theme names must fail on `SetTheme`.
- Stitched surfaces must preserve visible separation between panel/input/selection/dialog layers.

## Theme Authoring Checklist

- [ ] All required token groups are fully defined.
- [ ] Values are explicit and documented for this theme.
- [ ] No generated token values are introduced.
- [ ] Preview/apply/revert flows work in theme picker.
- [ ] Starter app remains visually consistent across bars, panels, input, and dialog.
