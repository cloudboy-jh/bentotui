---
title: "BentoTUI Rooms System + Surface Compositor Update"
description: "Docs and examples now reflect the named rooms system with surface-based final rendering."
pubDate: 2026-03-12
tags:
  - bentotui
  - go
  - bubbletea
  - rooms
  - docs
draft: false
---

# BentoTUI Rooms System + Surface Compositor Update

This update ships the new `registry/rooms` package and aligns docs, starter
flows, and bento examples with the current rendering contract.

Note: current starter and shipped bento defaults now use `Focus` layout
(body + anchored footer, no top/subheader rows).

## What Changed

- Added `registry/rooms` with named room functions
- Added a shared sizing/constrain engine under `registry/rooms/internal/engine`
- Added ASCII previews in room family files so users can quickly see each shape
- Updated starter app and shipped bentos to compose structure with rooms
- Restored `surface` as the final compositor for full-frame canvas paint and overlays
- Updated `bento init` scaffold template to follow the same render flow
- Updated docs to reflect current responsibilities for `theme`, `rooms`, and `surface`

## Render Contract (Canonical)

Use rooms for geometry, and `surface` for final paint/composition:

```go
screen := rooms.Focus(width, height, content, footer)

surf := surface.New(width, height)
surf.Fill(lipgloss.Color(theme.CurrentTheme().Surface.Canvas))
surf.Draw(0, 0, screen)

if dialogs.IsOpen() {
    surf.DrawCenter(viewString(dialogs.View()))
}

return tea.NewView(surf.Render())
```

## Why This Matters

Using `surface.Fill(...)` as the final pass guarantees deterministic background
paint for every frame and avoids whitespace/canvas gaps that can appear when a
screen is rendered as raw composed strings only.

## Where to Look

- Room API: `registry/rooms/`
- Starter reference: `cmd/starter-app/main.go`
- Bento references:
  - `registry/bentos/home-screen/main.go`
  - `registry/bentos/app-shell/main.go`
  - `registry/bentos/dashboard/main.go`
- Docs:
  - `docs/rooms.md`
  - `docs/architecture.md`
  - `docs/components.md`
