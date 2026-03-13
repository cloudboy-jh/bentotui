---
title: "BentoTUI Frame Hierarchy + Solid Row Rendering Update"
description: "Frame roles, anchored footer behavior, status pill semantics, and theme hierarchy tuning for cleaner cross-theme visuals."
pubDate: 2026-03-13
tags:
  - bentotui
  - go
  - bubbletea
  - themes
  - layouts
  - rendering
draft: false
---

# BentoTUI Frame Hierarchy + Solid Row Rendering Update

This update finishes the frame-level visual cleanup pass across components,
layouts, themes, and starter/bento defaults.

## What Shipped

- `Frame(...)` remains the multi-row grammar (`top`, `subheader`, `body`, `subfooter`).
- Bar rows now support explicit roles and footer modes:
  - top/subheader/footer row roles
  - anchored footer mode for strong command focus
- Anchored footer rows now render as one continuous strip (no segmented chip backgrounds).
- Top-row metadata moved to a single muted status pill pattern (`StatusPill("LIVE")`).
- Starter and shipped bentos now default to `Focus(...)` (body + anchored footer).
- Panel title/focus treatment reduced to avoid competing with anchored footer emphasis.
- Theme mapping/validation updated so frame hierarchy behaves more consistently across presets.

## Why This Matters

The previous state mixed multiple accent-heavy surfaces in the same frame,
especially in top bars, focused panel title badges, and footer actions. That
made some themes feel patchy and visually noisy.

The new hierarchy keeps the screen readable by assigning clear visual roles:

- top/subheader = metadata and context
- body surfaces = primary content
- anchored footer = command focus

## Usage Pattern

```go
top := bar.New(
    bar.RoleTopBar(),
    bar.StatusPill("LIVE"),
    bar.Left("bento app-shell"),
    bar.Right("workspace: demo"),
)

sub := bar.New(
    bar.RoleSubBar(),
    bar.Left("scope: nav"),
    bar.Right("nav active"),
)

foot := bar.New(
    bar.FooterAnchored(),
    bar.Left("scope: nav"),
    bar.Cards(
        bar.Card{Command: "j/k", Label: "move", Priority: 4, Enabled: true},
        bar.Card{Command: "tab", Label: "focus tabs", Priority: 3, Enabled: true},
        bar.Card{Command: "q", Label: "quit", Priority: 2, Enabled: true},
    ),
)

screen := layouts.Frame(w, h, top, sub, body, foot)
```

## Reference Files

- `registry/components/bar/bar.go`
- `registry/components/panel/panel.go`
- `styles/styles.go`
- `theme/adapter.go`
- `theme/theme.go`
- `registry/bentos/app-shell/main.go`
