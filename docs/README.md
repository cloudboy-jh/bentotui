# BentoTUI Documentation

- [architecture/architecture.md](./architecture/architecture.md) — Layer diagram, design principles, and render rules
- [architecture/bricks.md](./architecture/bricks.md) — Brick API reference and module deps
- [architecture/bentos.md](./architecture/bentos.md) — Full app composition contract and validation bento role
- [architecture/rooms.md](./architecture/rooms.md) — Room API and layout composition contract
- [theme-engine.md](./theme-engine.md) — Untouchable Theme Engine contract and workflow
- [astro-content-update.md](./astro-content-update.md) — Astro-ready architecture freeze update
- [next-steps.md](./next-steps.md) — Immediate priorities based on current repo state
- [roadmap.md](./roadmap.md) — Backlog and non-goals

Rooms define frame grammar + geometry; `surface` handles final frame paint and overlays.

Current default bentos/starter templates use a footer-anchored `Focus` room
(no top/subheader rows by default).
