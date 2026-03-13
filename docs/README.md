# BentoTUI Documentation

- [components.md](./components.md) — API reference for every component and module dep
- [layouts.md](./layouts.md) — Named layout API and render contract
- [astro-content-update.md](./astro-content-update.md) — Astro-ready changelog content for the layouts rollout
- [astro-content-frame-hierarchy-update.md](./astro-content-frame-hierarchy-update.md) — Astro-ready content for frame hierarchy and solid-row rendering cleanup
- [architecture.md](./architecture.md) — Layer diagram, design principles, and rules
- [next-steps.md](./next-steps.md) — Immediate priorities based on current repo state
- [roadmap.md](./roadmap.md) — Backlog and non-goals

Layouts define frame grammar + geometry; `surface` handles final frame paint and overlays.

Current default bentos/starter templates use a footer-anchored `Focus` layout
(no top/subheader rows by default).
