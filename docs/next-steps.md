# BentoTUI — Next Steps

## Current state (repo reality)

Completed in this directory today:

- Starter app is present and runnable at `cmd/starter-app/main.go`
- Registry embedding is wired via `registry/embed.go`
- `bento` CLI has working `init`, `add`, `list`, and `doctor` paths under `cmd/bento/`
- Registry components currently shipped: `surface`, `panel`, `bar`, `dialog`, `list`, `table`, `text`, `input`
- No `bentos/` directory exists yet (example screen patterns are still pending)

---

## Immediate priorities

### 1) Ship first bento examples

- Create `bentos/home-screen` first, matching the starter app interaction model
- Add 1-2 follow-up examples (`app-shell`, `dashboard`) to validate composition patterns
- Keep each bento runnable with `go run ./bentos/<name>`

### 2) Expand the component registry (next useful set)

- Add display helpers used by the starter-like screens: `badge`, `kbd`, `wordmark`
- Add common form/feedback controls: `select`, `checkbox`, `textarea`, `spinner`, `progress`
- Keep each new component copy-and-own compatible with `bento add <name>`

### 3) Tighten `bento init` output

- Simplify the generated `main.go` to a shorter starter shell
- Add explicit "next commands" comments (`bento add ...`, `go run .`)
- Ensure the scaffold stays easy to edit and does not feel framework-coupled

### 4) Add confidence checks

- Add `go test ./registry/...` coverage for component rendering behavior
- Add smoke tests for `logic.InstallComponent` and `logic.ScaffoldProject`
- Wire tests into a default local/CI command path

### 5) Start wrap/scaffold foundation

- Define deterministic manifest schema for `bento wrap`
- Implement `--manifest-only` and `--scaffold` before any AI-enhance layer

---

## Non-goals (still true)

- No web renderer or browser output
- No mouse-first interaction model
- No built-in app router or page framework
- No data-fetching abstraction layer
