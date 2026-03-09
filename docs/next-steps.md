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

## Detailed build backlog (restored)

### Components to build

Tier 1 (starter-facing):

- `badge` — inline colored label
- `kbd` — keyboard shortcut pair (`tab`, `cmd+k` style)
- `wordmark` — large centered title component

Tier 2 (core form/feedback):

- `select` — single-choice picker
- `checkbox` — toggle boolean input
- `textarea` — multiline text input
- `spinner` — loading indicator
- `progress` — horizontal progress bar

Tier 3 (advanced/overlay):

- `command` — command palette with fuzzy search
- `toast` — ephemeral stacked notifications
- `tabs` — keyboard-navigable tab switcher
- `separator` — horizontal/vertical divider

### Bento examples to build

- `home-screen` — mirror starter app pattern
- `app-shell` — header + sidebar + body + footer/status
- `dashboard` — cards + table composition
- `detail-view` — list + detail pane
- `form` — labeled inputs + validation hints
- `log-viewer` — filter + scrollable output
- `settings` — left nav + settings content
- `command-view` — command-palette-first screen

### CLI and platform items

- `bento upgrade <component>` — diff local copy vs registry version
- Improve `bento init` template toward a smaller single-screen scaffold
- Add overwrite strategy for `bento add` (`--force` or guided mode)

### Wrap + AI integration items

- `bento wrap --manifest-only`
- `bento wrap --scaffold`
- `bento wrap --enhance` (optional LLM pass after deterministic scaffold)
- MCP tools: `bento_wrap`, `bento_scaffold`, `bento_enhance`
- Public enhancement API: `bento.Enhance()`
- `llms.txt` context for enhancement/scaffold tooling

### Suggested execution order

1. Create `bentos/home-screen`
2. Build Tier 1 components (`badge`, `kbd`, `wordmark`)
3. Add `app-shell` and `dashboard` bentos
4. Build Tier 2 components
5. Add remaining bentos (`detail-view`, `form`, `log-viewer`, `settings`, `command-view`)
6. Build Tier 3 components
7. Add tests for registry rendering + CLI logic
8. Implement `bento wrap` deterministic pipeline
9. Layer optional AI enhancement surface

---

## Non-goals (still true)

- No web renderer or browser output
- No mouse-first interaction model
- No built-in app router or page framework
- No data-fetching abstraction layer
