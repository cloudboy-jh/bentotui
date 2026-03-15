---
title: "BentoTUI Architecture Freeze: Untouchable Theme Engine + Validation Bento"
description: "BentoTUI now locks the core model around bentos, rooms, and bricks, with app-shell as a scenario-driven framework validation bento."
pubDate: 2026-03-15
tags:
  - bentotui
  - go
  - bubbletea
  - terminal-ui
  - architecture
  - themes
draft: false
---

# BentoTUI Architecture Freeze: Untouchable Theme Engine + Validation Bento

BentoTUI now has a locked architecture contract:

- `registry/bentos` = full apps
- `registry/rooms` = geometry/layout composition
- `registry/bricks` = UI components

And the visual system is now officially documented as the **Untouchable Theme Engine**.

## What changed

- Moved shared style helpers from `styles/` to `theme/styles/`.
- Reworked `registry/bentos/app-shell` into a scenario-driven validation bento.
- Added room-level split separation options (`WithGutter`, `WithDivider`).
- Added anchored footer card style modes (`plain`, `chip`, `mixed`).
- Expanded structured list rows with typed fields while keeping compatibility.
- Hardened panel render contract for ANSI-heavy content.

## Validation bento role

`registry/bentos/app-shell` is now the canonical proving ground for framework quality,
not a static demo.

It ships scenario groups for:

- layout
- hierarchy
- footer
- list
- overlay
- stress

This makes regressions reproducible via a deterministic tuple:

`scenario + viewport + theme + focus + snapshot`

## Docs reshaped

Architecture docs now live under `docs/architecture/`:

- `architecture.md`
- `bentos.md`
- `bricks.md`
- `rooms.md`

Top-level docs now focus on strategy and direction:

- `theme-engine.md`
- `roadmap.md`
- `next-steps.md`

## Why this matters

This removes layout/color guesswork by making responsibilities explicit:

- bentos orchestrate
- rooms allocate
- bricks paint

With themes driving visual semantics globally, app teams can compose faster
without per-screen color glue.
