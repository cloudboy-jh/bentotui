# BentoTUI Layer Architecture

Status: Active reference
Date: 2026-02-26

This document captures the top-down BentoTUI architecture across the four UI/runtime layers:

- `core`
- `ui/components`
- `ui/primitives`
- `ui/styles`

## Top-Down Diagram

```text
App / Feature Code
  |
  v
+------------------------------------------------------+
| core (runtime orchestration)                         |
|------------------------------------------------------|
| shell  router  layout  focus  theme  surface  msgs  |
| - lifecycle + message routing                        |
| - page/dialog/footer composition                     |
| - sizing + layer order: body -> footer -> scrim -> dialog |
+-------------------------+----------------------------+
                          |
                          v
+------------------------------------------------------+
| ui/components (composed UI units)                    |
|------------------------------------------------------|
| panel                footer                dialog    |
| - bounded models implementing tea.Model + SetSize    |
| - component behavior and layout-level semantics       |
+-------------------------+----------------------------+
                          |
                          v
+------------------------------------------------------+
| ui/primitives (reusable visual building blocks)      |
|------------------------------------------------------|
| card                  row/frame/inputrow             |
| - small render-focused units                          |
| - no app workflow knowledge                           |
+-------------------------+----------------------------+
                          |
                          v
+------------------------------------------------------+
| ui/styles (semantic style mapping)                   |
|------------------------------------------------------|
| styles.System maps theme tokens -> lipgloss styles   |
| e.g. status bar, footer card command/label, panels   |
+------------------------------------------------------+
                          |
                          v
                    theme.Theme tokens
```

## Layer Responsibilities

### 1) `core`

- Owns runtime orchestration and routing.
- Composes shell layers and update/view flow.
- Defines deterministic layer order and sizing behavior.

### 2) `ui/components`

- Hosts composed UI models (`footer`, `panel`, `dialog`).
- Implements bounded behavior with `SetSize`/`GetSize`.
- Uses primitives/styles instead of ad hoc rendering.

### 3) `ui/primitives`

- Provides reusable visual building blocks (`card`, `row`, `frame`, `inputrow`).
- Keeps APIs minimal and render-oriented.
- Avoids business logic and command routing concerns.

### 4) `ui/styles`

- Centralizes semantic visual rules (`styles.System`).
- Maps `theme.Theme` tokens to reusable Lip Gloss style constructors.
- Keeps component/primitives free of scattered color literals.

## Design Rule of Thumb

- If it controls runtime flow: `core`.
- If it is a composed UI model: `ui/components`.
- If it is a small reusable render piece: `ui/primitives`.
- If it defines visual semantics: `ui/styles`.
