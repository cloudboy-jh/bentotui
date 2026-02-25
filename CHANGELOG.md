# Changelog

All notable changes to this project will be documented in this file.

The format follows Keep a Changelog style and this project targets Semantic Versioning.

## [Unreleased]

### Added

- introduced `ui/styles` as a dedicated style layer for semantic UI styling
- added searchable theme picker dialog with selection highlight and current-theme marker
- added explicit theme tokens for layered surfaces and interaction states
- expanded theme tests for token completeness and preset stability

### Changed

- moved UI packages into `ui/components/*` (`dialog`, `footer`, `panel`)
- promoted footer-first shell API (`WithFooterBar`, `WithFooter`)
- updated shell composition to treat footer as a first-class layer
- refreshed `cmd/test-tui` harness around footer, dialogs, and theme interactions
- rewrote README to reflect the new architecture and usage paths

### Removed

- removed top-level UI package paths (`dialog`, `statusbar`, `panel`, `styles`)

## [0.1.0-initial] - 2026-02-23

### Added

- initial BentoTUI framework foundation
- core app shell, router, layout, focus, and theme modules
- early dialog/footer/panel component set
