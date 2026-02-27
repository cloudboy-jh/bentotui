# BentoTUI Framework Progress Update

Date: 2026-02-26
Status: Active (early production)

## Current State

BentoTUI has moved from a rough harness into a structured framework baseline with a clear runtime/UI split:

- runtime/core: `shell`, `router`, `layout`, `focus`, `surface`, `theme`, `core`
- UI layer: `ui/components/*`
- shared primitives: `ui/primitives`
- style layer: `ui/styles`

This separation is now documented and enforced through:

- `project-docs/component-system-reference.md`
- `project-docs/next-steps.md`
- `project-docs/framework-roadmap.md`

## Shipped Progress

- Footer-first shell API is in place (`WithFooterBar`, `WithFooter`).
- Dialog system is stabilized for routing and bounds behavior.
- Theme picker flow is functional and deterministic.
- Theme set is locked to:
  - `catppuccin-mocha` (default)
  - `dracula`
  - `osaka-jade`
- Theme persistence and reload behavior are active.
- Starter app command flow is now command-string based:
  - type `/pr`, `/issue`, `/branch` in input, run on Enter
  - legacy aliases (`/dialog`, `/theme`, `/page`) still resolve
  - `/` no longer auto-opens a modal
- Starter app footer cards are command-first (`/pr`, `/issue`, `/branch`) and no longer tied to hotkeys.
- UI naming is now card-first across the stack (`Card`, `Cards`, `LeftCard`, `RightCard`).
- Shared primitive verbs are standardized as `Render*` (`RenderRow`, `RenderFrame`, `RenderInputRow`).
- README and changelog discipline are in place (`CHANGELOG.md`, early-production warning tag).

## Architecture Direction (Locked)

- No ad hoc components.
- Contract-first component behavior (`SetSize`/`GetSize`, bounded rendering, deterministic key routing).
- Visual direction: flat, minimal, solid card surfaces with low chrome.
- Dialog layering remains deterministic:
  - `body -> footer -> scrim -> dialog`

## Current Gaps

- Remaining visual gaps in some component rows still require bounded-row hardening.
- Focus manager needs a clearer API/state contract and stronger test coverage.
- Visual polish is improved but not fully normalized across all components/themes.

## Next Priorities

1. Footer card model and deterministic truncation contract.
2. Focus manager hardening (`SetRing`, `SetIndex`, `FocusBy`, events, wrap/enable behavior).
3. Shared UI primitives (modal frame, input surface, list row, footer card).
4. Broader regression coverage (dialog bounds, footer rendering, key routing).

## Validation Snapshot

Recent framework passes are consistently validated with:

- `go test ./...`
- `go vet ./...`

## Reference Docs

- Main spec: `project-docs/bentotui-main-spec.md`
- Layer architecture: `project-docs/layer-architecture.md`
- Next execution list: `project-docs/next-steps.md`
- Component contract: `project-docs/component-system-reference.md`
- Framework roadmap: `project-docs/framework-roadmap.md`
- Rendering ADR: `project-docs/rendering-system-design.md`
