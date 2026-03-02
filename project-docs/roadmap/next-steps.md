# BentoTUI Next Steps

Status: Active
Date: 2026-03-01

This list is intentionally component-focused and execution-oriented.
Use this with `project-docs/design/component-system-reference.md`.

## 1. Command Palette Component

- [ ] Add command palette dialog component.
- [ ] Route `/` to command-entry/palette workflow (as finalized behavior).
- [ ] Keep canonical slash commands coherent (`/dialog`, `/theme`, `/page`) while supporting legacy aliases.

## 2. Component Regression Coverage

- [ ] Dialog tests: custom enter routing + bounds stability.

## 3. CLI + Distribution MVP

- [ ] Add `cmd/bento` CLI entrypoint.
- [ ] Implement `bento init` (scaffold minimal Bento app shell).
- [ ] Implement `bento add <component>` for copy-and-own component installs.
- [ ] Implement `bento doctor` for theme token + layering + sizing contract checks.
- [ ] Add GitHub Actions tag workflow for GoReleaser CLI releases.
- [ ] Keep module semver tagging flow as canonical package distribution path.
