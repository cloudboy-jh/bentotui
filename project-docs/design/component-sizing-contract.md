# Bento Component Sizing Contract

Status: Active reference
Date: 2026-03-01

This document defines non-negotiable sizing behavior for BentoTUI components.

## Core Rules

- `SetSize(width, height)` is authoritative allocation input.
- Components render only inside assigned bounds.
- Parent layout area is never inferred from child intrinsic content.
- In bounded views, each rendered row width must equal assigned width.

## Shell-Level Sizing

- Shell receives viewport size from `tea.WindowSizeMsg`.
- Body height is computed as `viewport - header - footer`.
- Dialog manager receives full viewport size and centers bounded dialog content.

## Split/Layout Sizing

- Fixed and flex allocations come from parent bounds.
- Integer remainder goes to the last flex item.
- Child `SetSize` is called with allocated region dimensions.

## Panel Sizing

- Panel outer size equals assigned bounds.
- Content width is panel width minus frame/border overhead.
- Title row consumes one row when present.
- Content rows are clipped/padded to fit panel content width.

## Row Primitive Contract

- `RenderRow` and `RenderStyledRow` must return exact target width.
- Styled content may be clipped, but background coverage cannot shrink.
- Row background must be painted as a full-width base before overlay content.
- Input rows use the same width guarantee as normal rows.

## Dialog Sizing

- Dialogs must be compact and viewport-bounded.
- Width and height must be clamped before rendering content.
- Suggested compact policy:
  - width range: 52 to 88 (then clamped by viewport)
  - height range: 14 to 24 (then clamped by viewport)
- Dialog content receives bounded inner size from dialog frame.

## Theme Picker Sizing

- Picker width is derived from assigned dialog content width.
- List/search/help rows are clipped and padded to exact width.
- Theme dialog must not expand to near-fullscreen on ultrawide terminals.

## Acceptance Criteria

- No right-edge render gaps in bars, panel rows, or input rows.
- Header and footer remain single-line and full-width.
- Dialog remains compact in large terminals and usable in small terminals.
- Input rows remain visually distinct from panel background across all themes.
- `go test ./...` includes width regression checks for these contracts.
