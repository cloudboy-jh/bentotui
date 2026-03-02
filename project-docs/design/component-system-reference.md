# BentoTUI Component System Reference

Status: Canonical
Date: 2026-03-01

Purpose: define the required component contract for all BentoTUI UI work.

## Scope

- Runtime/core orchestration: `core/*`
- UI models: `ui/containers/*`
- Shared render units: `ui/primitives/*`
- Semantic styling: `ui/styles/*`

## Required Component Contract

All UI components must implement:

- `tea.Model` (`Init`, `Update`, `View`)
- `SetSize(width, height int)`
- `GetSize() (width, height int)`

Rules:

- Render only inside assigned bounds.
- Never depend on global viewport state.
- Never emit unbounded rows.

## Key Routing Contract

Routing precedence:

1. active dialog
2. focused component in active page
3. shell/global bindings

Rules:

- text input focus owns character keys first
- global shortcuts are secondary
- footer cards prioritize command discoverability (`/dialog`, `/theme`, `/page`)
- `Confirm` dialogs may use manager-level `enter`
- `Custom` dialogs must receive key events directly

## Dialog Contract

Manager responsibilities:

- open/close lifecycle
- centered placement
- bounded `SetSize` for active dialog

Dialog responsibilities:

- honor assigned bounds
- clip list/input rows to width
- deterministic section layout (header/search/list/hints)

Shell layer order:

- `header -> body -> footer -> scrim -> dialog`

## Style and Token Contract

- Components use semantic styles via `ui/styles` only.
- No hardcoded component color literals for new work.
- Color governance: `project-docs/design/bento-color-design-system.md`.
- Sizing governance: `project-docs/design/component-sizing-contract.md`.

Stitched UI rules:

- Base row/surface paint must happen before content overlay.
- Inputs must stay visually separated from parent containers.
- Footer and header rows must remain a single continuous strip under card/text overlays.
- Newly navigated pages must immediately align to the current active theme.

## Component PR Checklist

- [ ] Implements `SetSize`/`GetSize` correctly
- [ ] No overflow at wide widths
- [ ] No clipping artifacts at narrow widths
- [ ] Focus state is visible
- [ ] Key routing precedence is correct
- [ ] Uses semantic styles only
- [ ] Uses documented theme tokens only
- [ ] Meets sizing contract
- [ ] Preserves stitched surface separation (panel/input/selection/dialog)
- [ ] Includes regression tests

## Reference Index

External research links live in `project-docs/research/external-reference-index.md`.
