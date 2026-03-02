# BentoTUI Framework Roadmap

Status: Active
Date: 2026-02-25

This roadmap defines the framework-level plan for BentoTUI. It is separate from `project-docs/roadmap/next-steps.md`, which tracks immediate execution items.
The locked 1.0 component queue lives in `project-docs/roadmap/next-components-tobuild.md`.

## Vision

BentoTUI is a framework layer over Bubble Tea for building complete terminal applications with:

- deterministic shell/layer behavior
- reusable UI component contracts
- theme-driven styling through a dedicated style layer
- strong rendering correctness under resize, overlays, and focus transitions

## Principles

- contract-first implementation (no ad hoc components)
- runtime/core separated from UI layer (`ui/containers/*`, `ui/styles`)
- shared UI primitives centralized in `ui/primitives`
- semantic styles only (no scattered raw color literals in components)
- strict documented theme tokens only (no color derivation)
- explicit component sizing contract for all bounded rendering
- dual distribution model: Go module for library users + CLI binary for copy/scaffold workflows
- test-backed behavior before API expansion

## Status Legend

- `planned`: scoped but not started
- `in-progress`: actively being implemented
- `done`: landed and validated
- `deferred`: intentionally postponed

## v0.1 Foundation (Current)

Status: `in-progress`

Goals:

- shell layering model (`header -> body -> footer -> scrim -> dialog`)
- router, layout, focus baseline
- dialog manager with confirm/custom flows
- theme presets and persistence
- initial UI layer split under `ui/containers/*` and `ui/styles`

Remaining focus to close v0.1:

- header/footer parity and visual normalization across component surfaces
- focus manager hardening and deterministic event contract
- starter-app-first component composition defaults
- regression coverage for dialog bounds and bar truncation behavior

## v0.2 Component System

Status: `planned`

Goals:

- standardized component contracts (`SetSize`/`GetSize`, clipping, focus behavior)
- shared primitives (modal frame, input surface, list row, footer card)
- deterministic key routing precedence across dialogs/pages/shell
- stronger component-level regression suite

## v0.3 Command UX

Status: `planned`

Goals:

- command palette workflow as first-class slash UX
- command/action registry and discoverability
- unified command execution model across harness and real apps

## v0.3 CLI Distribution (MVP)

Status: `planned`

Goals:

- ship `bento` CLI for shadcn-style copy-and-own workflows
- provide `bento init` for project scaffolding
- provide `bento add <component>` for local component source install
- provide `bento doctor` for token/layer/sizing contract checks
- publish CLI binaries via GitHub Actions + GoReleaser while keeping module tags as primary package distribution

## v0.4 Theme Registry and Extensibility

Status: `planned`

Goals:

- theme registry expansion (beyond built-ins)
- external theme loading model and override hierarchy
- schema validation and safer theme ingestion

## v0.5 Stability and API Freeze

Status: `planned`

Goals:

- API review and freeze candidates
- full framework examples and docs alignment
- release hardening with regression matrix and compatibility checks

## Cross-Cutting Tracks

- Rendering correctness: full-frame paint, bounded overlays, resize safety
- Testing: dialog/focus/footer/layout regressions in CI
- Docs: keep spec/roadmap/next-steps synchronized
- Changelog discipline: release notes in `CHANGELOG.md`

## Non-Goals (Near-Term)

- broad plugin architecture
- large public extension API surface before stability gates
- unsupported terminal feature workarounds beyond documented constraints

## Release Gate Checklist (Per Minor)

- [ ] all milestone acceptance items completed
- [ ] `go test ./...` and `go vet ./...` clean
- [ ] docs updated (`README`, spec, roadmap, next-steps)
- [ ] changelog entries updated for shipped scope
