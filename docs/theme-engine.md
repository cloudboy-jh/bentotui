# Untouchable Theme Engine

Official rule: BentoTUI ships a theme-first workflow. You should not hand-wire
per-component color systems for normal usage.

## What it means

- Choose a preset with `theme.SetTheme(name)`.
- Compose layout with `registry/rooms`.
- Render components from `registry/bricks`.
- Let `theme/styles` map semantic tokens to concrete `lipgloss.Style` values.

## Why this exists

- Prevent color drift between screens and apps.
- Keep hierarchy consistent (`canvas -> panel -> elevated -> interactive`).
- Make theme switching global and deterministic.
- Avoid app-level CSS/Tailwind-like color glue in terminal UIs.

## User paths

1. Copy a prebuilt bento (`registry/bentos/home-screen`, `dashboard`) and change behavior/content.
2. Start from `registry/bentos/app-shell` to validate UX composition, layering, and theme switching in a single-screen app.
3. Build custom app from `rooms + bricks` and rely on theme tokens for visual system defaults.

## Non-goals right now

- Per-example custom color overrides.
- Multiple competing style systems inside one app.
