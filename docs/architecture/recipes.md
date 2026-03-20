# BentoTUI Recipes

`registry/recipes/` contains copy-and-own composed patterns.

A **recipe** is not a low-level primitive. It is a reusable app-flow composition
built from bricks, room contracts, and app state wiring.

Install with:

```bash
bento add recipe <name>
```

Copied files land in `recipes/<name>/` inside your project.

---

## Current recipes

| Recipe | Purpose |
|---|---|
| `filter-bar` | Input + footer command strip composition for filter workflows |
| `empty-state-pane` | Empty-result pane composition with card + message body |
| `command-palette-flow` | Command palette open flow helper for dialog manager wiring |
| `vimstatus` | Vim-style statusline recipe with mode badge, context, and clock |

---

## Layering expectations

- Recipes must compose at least one official brick.
- Recipes may use room contracts for page shaping.
- Recipes should avoid raw `bubbles/*` imports directly.
- Recipes are app-facing and expected to be edited after copy.

---

## Brick vs Recipe

- **Brick**: smallest reusable UI piece (button-like, card-like, list-like).
- **Recipe**: workflow-level composition that includes one or more bricks.

Default decision path:

1. Use an official brick when it exists.
2. Use an official recipe when the workflow matches.
3. Build a local app-owned recipe if no official recipe fits.
4. Only propose new bricks for repeated cross-template gaps.
