# Plan: go-tealeaves Major Restructure, Theming, and Code Audit

## Context

go-tealeaves is at v0.2.0 with 9 modules. Before reaching v1.0, we need to:
- Fix package names that are unclear or too abbreviated
- Add a diff rendering package (teadiffr) for TUI diff display
- Implement the theming system (THEMING_PLAN.md) to replace fragmented hardcoded colors
- Comprehensive code audit for ClearPath, API design, godoc, stability annotations
- Document future extraction boundaries in teautils for command layer (post-MVP)
- Restructure documentation site from module-centric to component-centric (discoverability)

All work is in go-tealeaves only. Gomion import fixes are a separate effort.

---

## Phase 1: RENAME — Package Renames (3 commits)

**Order matters**: teadd must rename first because teadep imports it.

### 1a. teadd → teadrpdwn
Files affected:
- `teadd/` → `teadrpdwn/` (all .go files: package declaration change)
- `teadrpdwn/go.mod` — module path update
- `go.work` — `./teadd` → `./teadrpdwn`
- `teadep/go.mod` + `teadep/model.go` — import path + qualifier updates (`teadd.` → `teadrpdwn.`)
- `examples/teadd/` → `examples/teadrpdwn/` — directory, go.mod, main.go imports
- `examples/teadep/treenav/go.mod` — require/replace for teadrpdwn
- `justfile` — modules and examples lists
- Verification: `just tidy && just test`

### 1b. teadep → teadepview
Files affected:
- `teadep/` → `teadepview/` (all .go files)
- `teadepview/go.mod` — module path; replace for `../teadrpdwn`
- `go.work` — `./teadep` → `./teadepview`
- `examples/teadep/` → `examples/teadepview/`
- `justfile`
- Verification: `just tidy && just test`

### 1c. teatextsel → teatxtsnip
Files affected:
- `teatextsel/` → `teatxtsnip/` (all .go files)
- `teatxtsnip/go.mod` — module path
- `go.work` — `./teatextsel` → `./teatxtsnip`
- `examples/teatextsel/` → `examples/teatxtsnip/`
- `justfile`
- Verification: `just tidy && just test`

---

## Phase 2: TEADIFFR — New Diff Rendering Package (1 commit)

Per `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/PLAN-TEADIFF-PACKAGE.md`, using name `teadiffr`.

### Files to create
```
teadiffr/
  doterr.go         — Copy from existing package, change package name
  errors.go         — ErrDiff, ErrInvalidFile, ErrInvalidBlock, ErrEmptyDiff
  types.go          — FileStatus enum, CondensedBlock, FileDiff, DiffRenderer interface
  tui_renderer.go   — TUIRenderer struct (lipgloss), NewTUIRenderer(*TUIRendererArgs)
  render.go         — RenderFileDiffs(files []FileDiff, renderer DiffRenderer, width int) []string
  go.mod            — deps: lipgloss v2, x/ansi only. NO go-dt, gitutils, gomion
```

### Updates
- Add `./teadiffr` to `go.work`
- Add `teadiffr` to `justfile` modules list
- Verification: `go build ./teadiffr/... && go test ./teadiffr/...`

---

## Phase 3: THEMING — Color, Palette, Theme System (4 sub-phases)

Per `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teautils/THEMING_PLAN.md`.

### 3a. COLOR — teacolor subpackage (1 commit)
Create `teautils/teacolor/` subpackage (part of teautils module, no separate go.mod):
- `color.go` — `Color` type alias for `lipgloss.Color`
- `ansi256.go` — Color0 through Color255 constants
- `ansi_names.go` — Black, Red, Green, ..., BrightWhite (0-15)
- `named.go` — ~40-50 curated semantic aliases (Coral, SkyBlue, Gold, DarkGray, etc.)
- Tests verifying constants

### 3b. PALETTE — Palette struct + factories (1 commit)
In `teautils/`:
- `palette.go` — Palette struct (~30 semantic color slots), DarkPalette(), LightPalette(), AdaptivePalette(), DefaultPalette()
- Tests for palette construction and embedding pattern

### 3c. THEME — Theme struct + component themes (1 commit)
In `teautils/`:
- `theme.go` — Theme struct with common styles + component-specific theme structs (StatusBarTheme, HelpVisorTheme, ModalTheme, DropdownTheme, ListTheme, GridTheme)
- NewTheme(Palette) Theme, DefaultTheme()
- Tests

### 3d. INTEGRATE — WithTheme() on all components (1 commit per component, or batched)
Each component gets:
- `WithTheme(Theme)` method
- Internal: if theme set, use themed styles; else use existing defaults
- Existing `Default*Style()` updated to use teacolor constants (non-breaking)
- `With*Style()` methods continue to work as per-instance overrides over theme

Order: teastatus → teamodal → teadrpdwn → teadepview → teatree → teatxtsnip → teanotify → teagrid → teautils (help visor) → teadiffr

---

## Phase 4: AUDIT — Full Code Audit (multiple sub-phases)

### 4a. AUDIT-DOTERR — Missing error infrastructure (1 commit)
Add doterr.go + errors.go to:
- **teagrid** — sentinels: ErrGrid, ErrInvalidColumn, ErrInvalidRow, ErrInvalidData
- **teatxtsnip** — sentinels: ErrTextSnip, ErrClipboard, ErrInvalidPosition
- Audit all packages for `fmt.Errorf` in library code (replace with NewErr)

### 4b. AUDIT-CTOR — Constructor renames (1 commit)
Violations to fix:

| Package | Current | New Name |
|---------|---------|----------|
| teagrid | `New(columns)` | `NewGridModel(columns)` |
| teastatus | `New()` | `NewStatusBarModel()` |
| teatxtsnip | `New()` | `NewTextSnipModel()` |
| teatxtsnip | `NewSingleLine()` | `NewSingleLineTextSnipModel()` |
| teatxtsnip | `NewFromTextarea(ta)` | `NewTextSnipModelFromTextarea(ta)` |
| teatree | `NewModel(tree, height)` | `NewTreeModel(tree, height)` |
| teadrpdwn | `NewModel(...)` | `NewDropdownModel(...)` |

For each: create new function, add `// Deprecated:` alias on old one, update all internal callers + examples.

### 4c. AUDIT-ARGS — *Args struct audit (1 commit)
Check all constructors with 3+ params have *Args structs. Known:
- teadrpdwn.NewDropdownModel has 4 positional params + args — consider folding fieldRow/fieldCol into args
- Rename `ModelArgs` → `DropdownModelArgs` where type name is too generic

### 4d. AUDIT-CLEARPATH — ClearPath convention sweep (1-2 commits)
Per-file sweep of all non-test, non-main .go files:
- Replace `} else {` with early-return or goto-end pattern
- Split compound `if init; cond {` into separate declaration + condition
- Verify single `end:` label, named returns, vars declared before first goto
- doterr.go files: evaluate case-by-case (some `else` may be intrinsic)

### 4e. AUDIT-GODOC — Comments on all public symbols (1-2 commits)
Every exported type, function, method, constant, variable gets proper godoc:
- Comment starts with symbol name
- Priority by API surface: teagrid → teamodal → teadrpdwn → teadepview → teatree → teanotify → teastatus → teatxtsnip → teautils → teadiffr

### 4f. AUDIT-STABILITY — Stability annotations (1 commit)
Every exported symbol gets go-stability contract comment:
```go
// Contract:
// - Stability: provisional
// - Since: v0.3.0 (2026-03-02)
```
- All symbols: `provisional` (pre-v1.0, shaped but may change)
- Deprecated constructor aliases: `deprecated` with UseInstead

### 4g. AUDIT-MISC — Remaining best practices (1 commit)
- No compound if-init statements (33 files identified)
- regexp.MustCompile at package level only (no violations found)
- ANSI-aware width (already done, verify no regressions)
- No ignored errors (`_ =` for error returns)

---

## Phase 5: TEAUTILS-BOUNDARIES — Document Future Extraction Points (1 commit)

Add file-level comments documenting logical boundaries for future extraction:
- `key_registry.go` / `key_identifier.go` — "// EXTRACTION-BOUNDARY: teakeys"
- `render_help_visor.go` — "// EXTRACTION-BOUNDARY: teahelp (depends on teakeys)"
- `positioning.go` / `render_styled.go` — "// EXTRACTION-BOUNDARY: tealayout"
- Action registry / command prompt (future) — "// EXTRACTION-BOUNDARY: teacmds"

Purpose: when teautils grows with command layer, these markers guide the split.

---

## Phase 6: SITE — Restructure Documentation Site (2-3 commits)

### Problem
The current site has a 1:1 mapping between sidebar entries and Go modules. This buries
important features (Key Registry, Help Visor, Theming, Positioning) under a single
"teautils — Utilities" page. A developer searching for TUI key management or theming
will never find them.

### Solution
Restructure the site to focus on **Components** (what you can use) rather than
**Modules** (which Go package to import). Each component gets its own page regardless of
which module it lives in. A separate "Modules" reference section covers Go module
structure and `go get` commands.

### New sidebar structure in `site/astro.config.mjs`

```
Components/
  Overview
  --- UI Components ---
  Dropdown (teadrpdwn)
  Data Grid (teagrid)
  Modal Dialogs (teamodal)
  Notifications (teanotify)
  Tree View (teatree)
  Text Selection (teatxtsnip)
  Status Bar (teastatus)
  Dep Viewer (teadepview)
  Diff Renderer (teadiffr)
  --- Infrastructure ---
  Theming (teautils)
  Key Registry (teautils)
  Help Visor (teautils)
  Positioning & Layout (teautils)

Modules/
  Module Reference (go get commands, dependency graph, versioning)
```

### Files to create/modify
- `site/astro.config.mjs` — restructure sidebar
- `site/src/content/docs/components/index.md` — rewrite overview (component-centric, not module-centric)
- `site/src/content/docs/components/theming.mdx` — NEW: first-class theming page
- `site/src/content/docs/components/key-registry.mdx` — NEW: move from patterns/ to components/
- `site/src/content/docs/components/help-visor.mdx` — NEW: extracted from teautils page
- `site/src/content/docs/components/positioning.mdx` — NEW: extracted from teautils page
- `site/src/content/docs/components/teadiffr.mdx` — NEW: diff renderer page
- `site/src/content/docs/modules/index.mdx` — NEW: module reference with go get commands
- Rename existing component pages for new package names (teadd→teadrpdwn, etc.)
- Each component page gets a consistent "Module" info box noting: `go get github.com/mikeschinkel/go-tealeaves/<module>`

### Create `site/PLAN.md`
Document this restructure rationale for future reference: the shift from module-centric
to component-centric documentation, and the principle that documentation discoverability
should drive site structure while Go module boundaries serve technical concerns.

---

## Phase 7: DOCS — Other Documentation Updates (1 commit)

- `README.md` — update all package names, import paths, code examples
- `API.md` — update package references, constructor names
- `PLAN-TEADIFF-PACKAGE.md` — mark as implemented, note teadiffr name
- Individual package READMEs in renamed packages

---

## Phase 8: VALIDATE — Full Cross-Module Verification

1. `just test` — all tests pass
2. `just tidy` — all modules clean
3. `go vet ./...` per module — no warnings
4. Build all examples
5. Grep for stale references: `teadd\b`, `teadep\b`, `teatextsel\b` in .go/.mod/.work files
6. Verify no `fmt.Errorf` in library code
7. Verify all exported symbols have godoc + stability annotations
8. Tag v0.3.0 for all modules (9 renamed/existing + teadiffr = 10 module tags + root tag)

---

## Future: Pre-v1.0 — Command Layer (plan only, implement later)

Per `/Users/mikeschinkel/Projects/gomion/COMMAND_LAYER_PLAN.md`. This is required for
v1.0.0 but implemented as a separate follow-up after the current restructure.

### Architecture decisions (locked)
- **Action Registry + CommandPromptModel** → live in teautils (opt-in, per COMMAND_LAYER_PLAN.md)
- **tmux-inspired config language parser** → standalone go-pkg: go-cmdlang, usable beyond TUIs
- **Command dispatch + ActionState** → gomion-specific (gomtui)

### Dependency chain for v1.0.0
```
go-cmdlang (new standalone pkg)
  ↑
go-cfgstore (extends to use go-cmdlang parser for file loading)
  ↑
teautils (Action Registry, CommandPromptModel, help/status integration)
  ↑
gomion/gomtui (ActionState, command dispatch, config files)
```

### Implementation phases (from COMMAND_LAYER_PLAN.md)
1. **ACTIONS** — ActionHandler, ActionRegistry, TypedRegistry[S], scope/table concept in teautils
2. **DISPATCH** — Convert gomion CommandInvoker to handlers, refactor Update() methods (gomion)
3. **PROMPT** — CommandPromptModel widget in teautils, `:` prompt overlay
4. **PARSER** — tmux-inspired config language parser in go-cmdlang, grammar design, go-cfgstore integration
5. **HELP** — Help visor reads from action registry, status bar reflects scope
6. **POLISH** — Tab completion, command history, per-command help, macro support

### Prerequisites (must complete before command layer)
- Current restructure (renames, theming, audit) — this plan
- Bubble Tea v2 migration complete (done)
- MVP naming audit (done as part of AUDIT-CTOR in this plan)

---

## Future: Post-v1.0 — Community Contributions

### additional-bubbles PR
Submit a PR to [charm-and-friends/additional-bubbles](https://github.com/charm-and-friends/additional-bubbles)
to make the README genuinely usable for developers searching for Bubble Tea components.

The current format is a flat alphabetical list with one-line descriptions — the same
discoverability problem that cost us significant time when first building a TUI. The PR
should improve the README for **all** listed projects, not just add Tea Leaves:
- Add category headings (Tables/Grids, Overlays/Modals, Navigation, Inputs, Status/Chrome, etc.)
- Organize ALL existing entries into categories
- Improve descriptions where needed for scannability
- Add Tea Leaves components within the appropriate categories alongside everything else

Gated on v1.0.0 release — our components should be `stable` before listing.
Open an issue first to discuss the restructure with maintainers before submitting the PR.

### Bubble Tea Component Discovery Site (vision)
A dedicated website for discovering Bubble Tea components — first-party (Charm) and
third-party alike. Categories, search, screenshots, version compatibility, and
filtering. Fills the gap that additional-bubbles' flat README cannot: a proper discovery
experience for the Bubble Tea ecosystem. This is a separate project, not part of
go-tealeaves, and is gated behind v1.0.0 and the additional-bubbles PR (which tests
community appetite for better organization).

---

## Execution Order Summary

### This effort (implement now)
```
RENAME-DD → RENAME-DEP → RENAME-SEL → TEADIFFR →
COLOR → PALETTE → THEME → INTEGRATE →
AUDIT-DOTERR → AUDIT-CTOR → AUDIT-ARGS → AUDIT-CLEARPATH →
AUDIT-GODOC → AUDIT-STABILITY → AUDIT-MISC →
TEAUTILS-BOUNDARIES → SITE-RESTRUCTURE → DOCS → VALIDATE → tag v0.3.0
```
Estimated: ~20-25 commits.

### Next effort (implement before v1.0.0)
```
go-cmdlang (standalone parser) →
ACTIONS (teautils) → PROMPT (teautils) → HELP (teautils) →
DISPATCH (gomion) → PARSER-INTEGRATE (gomion + go-cfgstore) → POLISH →
tag v1.0.0
```

### Post-v1.0.0
```
additional-bubbles PR → Discovery site (if appetite exists)
```

---

## Critical Files

- `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/go.work` — updated for every rename + new package
- `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/justfile` — module/example lists
- `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teadep/model.go` — cross-package import of teadd
- `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teautils/THEMING_PLAN.md` — theming spec
- `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/PLAN-TEADIFF-PACKAGE.md` — teadiffr spec
- `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/site/astro.config.mjs` — sidebar structure
- `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/site/src/content/docs/components/` — all component pages
