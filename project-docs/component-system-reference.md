# BentoTUI Component System Reference

Status: Active reference
Date: 2026-02-24

Purpose: prevent ad hoc component work by defining one component contract, one dialog contract, and one reference index for implementation decisions.

## Source of Truth

- Runtime/core behavior: `shell`, `router`, `layout`, `focus`, `surface`, `theme`, `core`
- UI layer behavior: `ui/components/*`
- Styling contract: `ui/styles`

Any new UI behavior should be implemented against this document and `project-docs/next-steps.md`.

## Directory Contract

- `ui/components/dialog` - modal manager and modal content components
- `ui/components/footer` - footer layer model and segment behavior
- `ui/components/panel` - bordered/surface panel container
- `ui/styles` - semantic style mapping from `theme.Theme`

No new top-level UI package directories should be introduced.

## Component Contract (Required)

All UI components should implement:

- `tea.Model` (`Init`, `Update`, `View`)
- `SetSize(width, height int)`
- `GetSize() (width, height int)`

Rules:

- Components render only inside assigned bounds.
- Components do not read viewport size from global shell state.
- Components do not emit unbounded rows/lines; width is always clipped.

## Input and Key Routing Contract

Routing precedence:

1. active modal/dialog
2. focused component in active page
3. shell/global bindings

Rules:

- Text input focus receives character keys first.
- Global hotkeys should avoid consuming normal text keys.
- `Confirm` dialogs may use manager-level `enter`; `Custom` dialogs must receive key events directly.

## Dialog Contract (Current Priority)

Manager responsibilities:

- own modal open/close lifecycle
- own centered placement in shell overlay layer
- pass bounded size to active dialog content

Dialog content responsibilities:

- respect `SetSize` bounds
- clip input and list rows to local width
- render deterministic vertical sections (header/search/list/hints)

Target shell layer order:

- `body -> footer -> scrim -> dialog`

## Styling Contract

All component visuals should come from `ui/styles` semantic styles.

Disallowed:

- hardcoded per-component color literals for new work
- ad hoc frame composition duplicated across components

Allowed:

- component-specific behavior with shared frame/list/input primitives

## Done Checklist for Component PRs

- [ ] Uses `SetSize`/`GetSize` correctly
- [ ] No horizontal overflow at wide terminal sizes
- [ ] No clipping artifacts at narrow terminal sizes
- [ ] Focus state is visually obvious
- [ ] Key routing follows precedence rules
- [ ] Uses `ui/styles` semantic styles only
- [ ] Adds or updates regression tests

## External References

## OpenCode

- Theme docs: <https://opencode.ai/docs/themes/>
- Theme schema: <https://opencode.ai/theme.json>
- TUI root: <https://github.com/opencode-ai/opencode/tree/main/internal/tui>
- Dialog components: <https://github.com/opencode-ai/opencode/tree/main/internal/tui/components/dialog>
- Layout primitives: <https://github.com/opencode-ai/opencode/tree/main/internal/tui/layout>
- Style system: <https://github.com/opencode-ai/opencode/tree/main/internal/tui/styles>
- Theme package: <https://github.com/opencode-ai/opencode/tree/main/internal/tui/theme>

Key files:

- <https://github.com/opencode-ai/opencode/blob/main/internal/tui/components/dialog/theme.go>
- <https://github.com/opencode-ai/opencode/blob/main/internal/tui/components/dialog/commands.go>
- <https://github.com/opencode-ai/opencode/blob/main/internal/tui/theme/manager.go>
- <https://github.com/opencode-ai/opencode/blob/main/internal/tui/styles/styles.go>
- <https://github.com/opencode-ai/opencode/blob/main/internal/tui/layout/overlay.go>

## Crush

- UI root: <https://github.com/charmbracelet/crush/tree/main/internal/ui>
- Dialog package: <https://github.com/charmbracelet/crush/tree/main/internal/ui/dialog>
- Style registry: <https://github.com/charmbracelet/crush/tree/main/internal/ui/styles>
- UI model: <https://github.com/charmbracelet/crush/tree/main/internal/ui/model>
- List package: <https://github.com/charmbracelet/crush/tree/main/internal/ui/list>

Key files:

- <https://github.com/charmbracelet/crush/blob/main/internal/ui/model/ui.go>
- <https://github.com/charmbracelet/crush/blob/main/internal/ui/dialog/common.go>
- <https://github.com/charmbracelet/crush/blob/main/internal/ui/dialog/models.go>
- <https://github.com/charmbracelet/crush/blob/main/internal/ui/styles/styles.go>
- <https://github.com/charmbracelet/crush/blob/main/internal/ui/list/list.go>

## Local Files to Anchor Changes

- `shell/model.go`
- `ui/components/dialog/dialog.go`
- `ui/components/dialog/theme_picker.go`
- `ui/components/footer/footer.go`
- `ui/components/panel/panel.go`
- `ui/styles/styles.go`
- `project-docs/next-steps.md`
