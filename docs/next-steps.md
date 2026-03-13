# BentoTUI — Next Steps

## Current state (repo reality)

Completed in this directory today:

- Starter app is present and runnable at `cmd/starter-app/main.go`
- Registry embedding is wired via `registry/embed.go`
- `bento` CLI has working `init`, `add`, `list`, and `doctor` paths under `cmd/bento/`
- Registry component catalog is finalized and shipped (`surface`, `panel`, `bar`, `dialog`, `list`, `table`, `text`, `input`, `badge`, `kbd`, `wordmark`, `select`, `checkbox`, `progress`, `tabs`, `toast`, `separator`)
- Wave 1 bento examples are now present under `registry/bentos/` (`home-screen`, `app-shell`, `dashboard`)
- Named layout catalog shipped at `registry/layouts/` (15 layouts + shared sizing engine)

---

## Immediate priorities

### 1) Ship Wave 1 bentos (shipped this session)

- Create `registry/bentos/home-screen` first (canonical starter-style composition)
- Add `registry/bentos/app-shell` and `registry/bentos/dashboard` in the same wave
- Keep each bento runnable with `go run ./registry/bentos/<name>` and copy-and-own friendly

### 2) Expand bento catalog breadth (Wave 2)

- Add `detail-view`, `form`, `log-viewer`, `settings`, and `command-view`
- Prefer direct Bubbles primitives (for example `spinner`) inside bentos where appropriate
- Only propose new components if at least 2 bentos need the same missing abstraction

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

### Final component catalog (Bento-owned)

Core layout/container components:

- `surface`
- `panel`
- `bar`
- `dialog`

Display helpers:

- `badge`
- `kbd`
- `wordmark`

Form/feedback components:

- `select`
- `checkbox`
- `progress`

Advanced composition helpers:

- `tabs`
- `toast`
- `separator`

Primitive policy:

- Do not add `spinner` as a Bento component; use `charm.land/bubbles/v2/spinner`
- Default to direct Bubbles primitives unless Bento-specific composition value is clear
- Existing shipped primitive-like components remain supported, but are not the growth focus

### Bento examples to build

- `registry/bentos/home-screen` — mirror starter app pattern
- `registry/bentos/app-shell` — header + sidebar + body + footer/status
- `registry/bentos/dashboard` — cards + table composition
- `registry/bentos/detail-view` — list + detail pane
- `registry/bentos/form` — labeled inputs + validation hints
- `registry/bentos/log-viewer` — filter + scrollable output
- `registry/bentos/settings` — left nav + settings content
- `registry/bentos/command-view` — command-palette-first screen

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

1. Keep `registry/bentos/home-screen` aligned with starter-app behavior
2. Add remaining bentos (`detail-view`, `form`, `log-viewer`, `settings`, `command-view`) under `registry/bentos/`
3. Add tests for registry rendering + CLI logic
4. Improve `bento init` scaffold clarity
5. Implement `bento wrap` deterministic pipeline
7. Layer optional AI enhancement surface

---

## Non-goals (still true)

- No web renderer or browser output
- No mouse-first interaction model
- No built-in app router or page framework
- No data-fetching abstraction layer
