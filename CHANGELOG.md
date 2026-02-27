# Changelog

All notable changes to this project will be documented in this file.

The format follows Keep a Changelog style and this project targets Semantic Versioning.

## [Unreleased]

### Added

- introduced `ui/styles` as a dedicated style layer for semantic UI styling
- added searchable theme picker dialog with selection highlight and current-theme marker
- added explicit theme tokens for layered surfaces and interaction states
- expanded theme tests for token completeness and preset stability
- added shared UI render primitives under `ui/primitives` (`chip`, `row`, `frame`, `inputrow`)
- added `ui/components/header` as a top statusline twin of footer card behavior
- added focus manager event contract with `FocusChangedMsg {From, To}`
- added non-persistent theme preview flow with enter commit and esc revert

### Changed

- moved UI packages into `ui/components/*` (`dialog`, `footer`, `panel`)
- promoted footer-first shell API (`WithFooterBar`, `WithFooter`)
- updated shell composition to treat footer as a first-class layer
- refreshed `cmd/starter-app` starter app around header/footer, dialogs, and theme interactions
- rewrote README to reflect the new architecture and usage paths
- finalized footer statusline behavior with deterministic truncation and one-row chip rendering
- updated starter app shell layering to `header -> body -> footer -> scrim -> dialog`

### Removed

- removed top-level UI package paths (`dialog`, `statusbar`, `panel`, `styles`)

## [0.1.0-initial] - 2026-02-23

### Added

- initial BentoTUI framework foundation
- core app shell, router, layout, focus, and theme modules
- early dialog/footer/panel component set
