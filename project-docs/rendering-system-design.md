# BentoTUI Rendering System Design

Status: Accepted
Date: 2026-02-24
Owner: BentoTUI core
Decision: Adopt a draw-first, area-owned renderer (Crush-style) for framework internals.

---

## 1. Why this document exists

BentoTUI must stop making rendering changes ad hoc.

This document defines the non-negotiable rendering system contract for BentoTUI, including fullscreen ownership, paint guarantees, layout behavior, overlay ordering, and acceptance criteria.

If implementation differs from this document, the implementation is wrong.

---

## 2. Problem statement

Observed failures in BentoTUI:

- right-edge unpainted bands
- content/background anchoring mismatch
- overlay composition inconsistencies under terminal blur/transparency
- mixed rendering semantics (string-first in some paths, layer-first in others)

These are architecture failures, not styling issues.

---

## 3. External proof and references

### 3.1 Crush (reference for selected architecture)

Crush proves a robust draw-first approach in production:

- Draw from live area each frame: `internal/ui/model/ui.go#L1823`
- Explicit full-frame clear before draw: `internal/ui/model/ui.go#L1832`
- Dialog drawn last over full bounds: `internal/ui/model/ui.go#L1916-L1920`
- Fullscreen set in view metadata: `internal/ui/model/ui.go#L1947`
- Background color set in view metadata: `internal/ui/model/ui.go#L1949`

Program init reference (no `WithAltScreen` at startup):

- `internal/cmd/root.go#L94-L99`

### 3.2 OpenCode (reference for string-layout hardening)

OpenCode proves robust width/height wrappers in a string-first renderer:

- Program enables alt screen at startup: `cmd/root.go#L120-L123`
- Root size propagation from `WindowSizeMsg`: `internal/tui/tui.go#L186-L188`
- Split wrapper forces full area + background: `internal/tui/layout/split.go#L118-L120`
- Container wrapper forces area + background: `internal/tui/layout/container.go#L50`, `#L71-L73`
- Overlay merge strategy: `internal/tui/layout/overlay.go#L36`

---

## 4. Architecture options

### Option A: Draw-first renderer (selected)

Model:

- Root `View` returns metadata + drawable root content.
- Root drawable draws from runtime area each frame.
- Layout is rectangle-first and region-anchored.
- Overlays are composited last in deterministic z-order.

Pros:

- Most robust to resize, blur, transparency, and terminal quirks.
- Matches Bubble Tea v2 + Ultraviolet primitives directly.
- Best long-term foundation for framework internals.

Cons:

- Requires strict discipline around region contracts.

### Option B: String-first renderer (not selected)

Model:

- Compose strings, enforce `Width/Height/Background` wrappers everywhere.
- Overlay by string merge.

Pros:

- Easier to prototype.

Cons:

- More fragile under edge cases.
- More difficult to guarantee paint correctness in complex overlays.

Decision:

- Choose Option A as framework baseline.
- Option B patterns are allowed only as compatibility helpers at edges, never as core renderer policy.

---

## 5. Core rendering principles (non-negotiable)

1. Single frame owner
   - Only app root owns final frame composition.

2. Area is source of truth
   - Runtime draw area (`Dx`, `Dy`) is authoritative each frame.

3. Full-frame paint guarantee
   - Every frame paints all cells in viewport bounds.

4. Region anchoring
   - Every allocated layout region paints its own background before content.

5. Deterministic layering
   - Exact order: `shell-bg -> body-bg -> body-content -> footer -> scrim -> dialog`.

6. No lossy conversion in core path
   - Core render path must not rely on `fmt.Sprint(layer)` semantics.

7. Overlay always last
   - Dialog layers render after all base layers and capture input while open.

8. Resize first, draw second
   - Size propagation must occur before drawing dependents.

9. Fullscreen default
   - Framework defaults to fullscreen app mode.

10. Debuggability
    - Render dimensions and major layer bounds must be introspectable.

---

## 6. Root renderer contract

### 6.1 `app.View()` contract

Root view must set:

- `AltScreen = true` by default (configurable)
- `BackgroundColor = theme.Background`
- `Content = root drawable`

### 6.2 Root drawable contract

Each frame:

1. Read `w,h` from draw area.
2. Sync router/footer/dialog sizes from `w,h`.
3. Paint full shell background (`w x h`).
4. Paint body background (`w x bodyH`).
5. Draw body content inside body region.
6. Draw footer region.
7. If dialog open: draw scrim and centered dialog last.

No partial frame paths are allowed.

---

## 7. Layout contract

`layout.Split` must be region-anchored:

- Compute allocations from parent bounds.
- For each child region:
  - paint region background for `allocW x allocH`
  - draw child content clipped/placed in that region
- Distribute integer remainder to last flex region.

Child intrinsic content bounds must never define parent paint coverage.

---

## 8. Component contract

Components provide content; they do not own frame paint.

- `SetSize(w,h)` is authoritative allocation input.
- `View()` returns content for allocated region.
- Panels/footer/dialog may style surfaces, but root/layout guarantee region paint.

---

## 9. Startup and resize contract

- Program startup must receive initial `WindowSizeMsg` from Bubble Tea.
- Root may request size explicitly as a safety net.
- Draw path must still rely on live area for correctness.
- No manual ANSI fullscreen hacks in app code.

---

## 10. Input and overlay contract

- Dialog manager updates before page updates.
- While dialog is open, page input is blocked.
- Overlay layers draw last and are centered in viewport.

---

## 11. Disallowed patterns

1. Mixed core render paradigms in one frame path.
2. Relying on child natural width/height to paint parent regions.
3. Root composition that can skip painting any viewport cells.
4. String fallback as primary core renderer.
5. Manual terminal escape toggles for fullscreen lifecycle.

---

## 12. Acceptance criteria

A rendering change is complete only when all are true:

1. No right-edge or bottom-edge unpainted bands.
2. No startup ghost/bleed frame.
3. Dialogs always overlay (never side-by-side artifacts).
4. Blur/transparency terminals render continuous surfaces.
5. Resize stress loops (grow/shrink/maximize/restore) stay artifact-free.
6. Footer bar spans full width always.
7. `go test ./...` passes.

---

## 13. Validation matrix

Manual checks:

- Windows Terminal (blur/transparency on)
- Non-blur baseline terminal
- At least one alternate terminal implementation if available

Scenarios:

- startup
- route switch
- focus cycle
- dialog open/close
- rapid resize

---

## 14. Implementation phases

Phase 1: Rendering correctness only

- root draw contract
- layout region anchoring
- overlay ordering
- fullscreen defaults and metadata

Phase 2: Visual parity polish

- spacing rhythm
- border density
- tonal hierarchy refinements

No visual polish work should bypass Phase 1 requirements.

---

## 15. Change management

- Any deviation from this design requires an explicit ADR update.
- If a quick fix conflicts with this document, reject the quick fix.
- New components must declare how they satisfy region paint guarantees.
