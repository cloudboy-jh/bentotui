# Product Direction: Ship Faster, Churn Less

This document is the product contract for BentoTUI's near-term direction.

## Product model

- `bricks`: copy-and-own UI building blocks (`bento add <brick>`)
- `recipes`: copy-and-own composed flow patterns (`bento add recipe <name>`)
- `rooms`: stable page layout contracts (import-only)
- `bentos`: full app templates (state + routing + composition)

Low-level UI is bricks. Recipes are app-flow composition. Bentos are app templates, not a giant component catalog.

## Scope discipline

BentoTUI does not try to ship every possible UI primitive as an official brick.

New official bricks are added only when all conditions are true:

1. The same gap appears in at least two bentos.
2. The gap cannot be cleanly composed from existing bricks.
3. The proposed brick has docs, tests, and an upgrade story.

If these are not true, build a local app-owned brick in the user project.

## Bento-first decision policy

When building app features:

1. Use an official Bento brick if it exists.
2. If a matching recipe exists, use it before inventing bespoke orchestration.
3. If missing, compose from existing bricks and rooms.
4. If still missing, create a local custom recipe in app code.
5. Upstream as official only after repeated cross-template demand.

This keeps shipping velocity high without forcing framework churn.

## Stable contracts

The stable shared imports are:

- `github.com/cloudboy-jh/bentotui/theme`
- `github.com/cloudboy-jh/bentotui/theme/styles`
- `github.com/cloudboy-jh/bentotui/registry/rooms`

Official bricks and recipes are stable by name through `bento add` and `bento list`.

## v1 core surface

Core bricks:

- `surface`, `card`, `bar`, `input`, `list`, `table`, `dialog`, `select`, `checkbox`, `progress`

Core rooms:

- `Focus`, `AppShell`, `SidebarDetail`, `DiffWorkspace`

Canonical bentos:

- `app-shell`, `home-screen`, `detail-view`

Everything else is optional and must justify its maintenance cost.

## Enforcement

Guardrails should fail CI when:

- room/theme/brick layering rules are broken
- starter/bentos bypass official brick boundaries without allowed exceptions
- docs and CLI wording drift from the bricks/rooms/bentos model

Product trust is a release feature.
