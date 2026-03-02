# BentoTUI Layer Architecture

Status: Active reference
Date: 2026-03-02

This document captures the top-down BentoTUI architecture across the four UI/runtime layers:

- `core` - Runtime orchestration and layout
- `ui/containers` - Complex UI components with borders/layout
- `ui/widgets` - Simple content components
- `ui/primitives` - Low-level rendering utilities
- `ui/styles` - Theme mapping

## Top-Down Diagram

```text
App / Feature Code
  |
  v
+------------------------------------------------------+
| core (runtime orchestration)                         |
|------------------------------------------------------|
| shell  router  layout  focus  theme  msgs           |
| - lifecycle + message routing                        |
| - page/dialog/footer composition                     |
| - sizing + layer order: body -> footer -> scrim -> dialog |
| - canvas-based layout (Horizontal/Vertical)          |
+-------------------------+----------------------------+
                          |
                          v
+------------------------------------------------------+
| ui/containers (complex UI components)                |
|------------------------------------------------------|
| panel                bar                dialog       |
| - bounded models implementing tea.Model + SetSize    |
| - manage borders, layout, and child positioning      |
| - panel holds content (widgets) inside bordered area |
+-------------------------+----------------------------+
                          |
                          v
+------------------------------------------------------+
| ui/widgets (simple content components)               |
|------------------------------------------------------|
| card    input    list    table    text               |
| - render content only, no layout management          |
| - used inside panels/containers                      |
| - implement SetSize for content fitting              |
+-------------------------+----------------------------+
                          |
                          v
+------------------------------------------------------+
| ui/primitives (rendering utilities)                  |
|------------------------------------------------------|
| row    frame    surface                              |
| - low-level render helpers (RenderRow, Fill, etc)   |
| - used BY containers, not directly by apps           |
+-------------------------+----------------------------+
                          |
                          v
+------------------------------------------------------+
| ui/styles (semantic style mapping)                   |
|------------------------------------------------------|
| styles.System maps theme tokens -> lipgloss styles   |
| e.g. status bar, footer card command/label, panels   |
+-------------------------+----------------------------+
                          |
                          v
                     theme.Theme tokens
                     (15 presets via bubbletint)
```

## Layer Responsibilities

### 1) `core`

- Owns runtime orchestration and routing
- Composes shell layers and update/view flow
- **Layout package**: Canvas-based positioning (Horizontal/Vertical with Fixed/Flex constraints)
- Defines deterministic layer order and sizing behavior

### 2) `ui/containers`

- Hosts complex UI models (`panel`, `bar`, `dialog`)
- Implements bounded behavior with `SetSize`/`GetSize`
- **Panel**: Bordered box that holds widgets as content
- **Bar**: Header/footer with card system
- **Dialog**: Modal overlays

### 3) `ui/widgets`

- Simple content components (`card`, `input`, `list`, `table`, `text`)
- Render content only - no border/layout management
- Used inside panels/containers
- Implement `SetSize` for proper content fitting

### 4) `ui/primitives`

- Low-level rendering utilities (`RenderRow`, `Fill`, `FitWidth`)
- Used BY containers internally
- Apps should not use these directly

### 5) `ui/styles`

- Centralizes semantic visual rules (`styles.System`)
- Maps `theme.Theme` tokens to reusable Lip Gloss style constructors
- Keeps components free of scattered color literals

## Design Rule of Thumb

- If it controls runtime flow: `core`
- If it manages borders and layout: `ui/containers`
- If it renders content only: `ui/widgets`
- If it's a low-level render helper: `ui/primitives`
- If it defines visual semantics: `ui/styles`

## Usage Pattern

```go
// Layout positions containers
root := layout.Horizontal(
    layout.Flex(1, panel.New(
        panel.Title("Sidebar"),
        panel.Content(widgets.NewList()),  // widget inside container
    )),
    layout.Flex(2, panel.New(
        panel.Title("Main"),
        panel.Content(widgets.NewInput()),
    )),
).WithGutterColor(theme.Border.Subtle)
```

## Key Insight

**Containers hold widgets.** Widgets don't know about borders or layout. The layout system positions containers, containers provide borders, widgets fill the content area.
