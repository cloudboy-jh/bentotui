# BentoTUI OpenCode Rules

This file defines repository-level contribution rules for BentoTUI.

## Architecture priority

1. Bentos and recipes are the primary product surface.
2. Bricks are foundational UI primitives and should evolve more conservatively.
3. Rooms are stable layout contracts and stay import-only.

## Layering contracts

- `registry/rooms/*` must remain geometry-only:
  - no `theme` imports
  - no `registry/bricks/*` imports
  - no raw `bubbles/*` imports
- `registry/bricks/*` must stay standalone:
  - no cross-brick imports
- `registry/recipes/*` are app-flow compositions:
  - must import at least one official brick from `registry/bricks/*`
  - must not import raw `bubbles/*` directly (spinner exception remains allowed)
  - are copy-and-own and expected to be edited downstream
- `registry/bentos/*` are template-grade full apps:
  - should compose rooms + bricks (+ recipes when useful)
  - should not import raw `bubbles/*` directly (spinner exception remains allowed)

## Theme ownership contract

- Keep theme in model-owned state (`m.theme`).
- Propagate theme changes via `SetTheme(...)` on composed bricks/recipes.
- In bento `View()` methods, render from model-owned theme state, not global lookups.

## Recipe definition (strict)

A BentoTUI recipe is valid only if it composes at least one official brick.

Examples:

- valid: recipe composes `bar`, `input`, `dialog`, `card`, `badge`, etc.
- invalid: recipe uses only `theme` + `theme/styles` + `lipgloss` with no bricks.

## Docs and catalog sync

When adding or changing official bricks/recipes/bentos:

- update CLI catalog entries in `cmd/bento/logic/add.go` if needed
- update relevant docs in `docs/architecture/*` and `docs/usage-guide.md`
- keep names and descriptions aligned between registry and docs

## Validation

Run before merging:

```bash
go test ./...
```

Policy tests in `internal/policy/guardrails_test.go` are source-of-truth for enforced rules.
