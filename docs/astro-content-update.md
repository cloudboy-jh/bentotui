---
title: "BentoTUI Architecture Freeze: Untouchable Theme Engine + Charm-Backed Bricks"
description: "BentoTUI now locks the core model around bentos, rooms, and bricks, with app-shell as a single-screen UX proving ground and Charm-backed component internals."
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

# BentoTUI Architecture Freeze: Untouchable Theme Engine + Charm-Backed Bricks

BentoTUI now has a locked architecture contract:

- `registry/bentos` = full apps
- `registry/rooms` = geometry/layout composition
- `registry/bricks` = UI components

And the visual system is now officially documented as the **Untouchable Theme Engine**.

## What changed

- Moved shared style helpers from `styles/` to `theme/styles/`.
- Reworked `registry/bentos/app-shell` into a single-screen UX proving ground.
- Added room-level split separation options (`WithGutter`, `WithDivider`).
- Added anchored footer card style modes (`plain`, `chip`, `mixed`).
- Migrated `list`, `table`, `progress`, `select`, `checkbox`, and `tabs` to Charm-backed internals.
- Added `filepicker` brick backed by `bubbles/filepicker`.
- Hardened panel render contract for ANSI-heavy content.

## App-shell role

`registry/bentos/app-shell` is now the canonical proving ground for UX composition
quality, not a scenario harness.

It demonstrates one intentional screen using rooms + bricks:

- left rail
- main table/list/progress workspace
- anchored command footer
- command palette with full theme switching

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
- bricks wrap Charm primitives, then paint with Bento theme tokens

With themes driving visual semantics globally, app teams can compose faster
without per-screen color glue.
