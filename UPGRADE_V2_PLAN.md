# Upgrading go-tealeaves to Charm v2

## Context

Charm Bracelet released v2.0.0 of Bubble Tea, Lip Gloss, and Bubbles — the first major breaking change in the project's history. The new version brings a ground-up rewritten renderer (Cursed Renderer), declarative View model, enhanced keyboard support, native cursor control, clipboard over SSH, and Lip Gloss as a pure styling library (no more I/O contention with Bubble Tea).

go-tealeaves has 7 modules (~80 .go files) all on Charm v1. This plan upgrades everything to v2 and produces an `UPGRADE_GUIDE_V2.md`.

This is a planning artifact. Implementation happens in separate execution sessions.

**Project location:** `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/`

## Support Policy

### Charm v1 compatible (`release/charm-v1`)

- Best-effort maintenance only; no guaranteed response time.
- No new features.
- Intended only to give existing users time to migrate.

### Charm v2 compatible (`main`)

- Primary active development line.
- New features and fixes land here.
- Still pre-v1.0 — APIs may change before go-tealeaves v1.0 release.

## Final Decisions Locked In

1. No `/v2` module path suffixes for go-tealeaves modules in this effort.
2. Module versions remain `v0.x.y` (pre-v1 strategy).
3. Compatibility target is Charm v2, but go-tealeaves versions do not track Charm major versions.
4. Branch name for Charm v1 compatibility line is `release/charm-v1`.
5. External migration guide filename is `UPGRADE_GUIDE_V2.md`.
6. Example validation is strict: examples must build and run (smoke-validated).
7. `release/charm-v1` health is a hard release gate for publishing Charm v2 tags.
8. Keep `github.com/charmbracelet/x/ansi` unless a concrete dependency break requires change.
9. `teatable` is renamed to `teagrid` everywhere.
10. `v0.2.0` is a synchronized release cut, not staggered per module.

## Versioning and Tagging Policy

1. Root tags are plain: `v0.1.0`, `v0.2.0`.
2. Submodule tags use Go multi-module format, e.g. `teadd/v0.1.0`, `teadd/v0.2.0`.
3. Breaking pre-v1 migration uses minor bump (`v0.Y+1.0`), so baseline `v0.1.0` -> migrated `v0.2.0`.
4. No standalone tags for `examples/*` modules.
5. Baseline `v0.1.0` tags are created only for modules that exist at tagging time.
6. `teanotify` and `teagrid` are introduced in baseline Charm v1 form and receive `v0.1.0` tags in that baseline.

## Required Linked Plans and Hard Gates

`v0.2.0` is blocked on completion of all linked plans below.

1. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md`
2. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teautils/THEMING_PLAN.md`
3. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teanotify/PLAN.md`
4. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teagrid/PLAN.md`
5. `/Users/mikeschinkel/Projects/gomion/CHARM_V1_PLAN.md`
6. `/Users/mikeschinkel/Projects/gomion/CHARM_V2_PLAN.md`
7. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/site/UPGRADE_V2_PLAN.md`

For dependent plans, completion means their own acceptance gates are satisfied.

## Scope

In scope:
1. Full Charm v2 migration for existing go-tealeaves modules, examples, and docs.
2. Introduce Charm v1 baseline modules for `teanotify` and `teagrid` (fork/rename/license/readme/setup), then tag baseline `v0.1.0`.
3. Define and enforce cross-plan hard gates for synchronized `v0.2.0` release.

Out of scope for this plan's implementation steps:
1. Detailed implementation steps for `teanotify` Charm v2 (owned by `teanotify/PLAN.md`).
2. Detailed implementation steps for `teagrid` Charm v2 (owned by `teagrid/PLAN.md`).
3. Detailed implementation steps for theming system (owned by `teautils/THEMING_PLAN.md`).
4. Detailed implementation steps for gomion migrations (owned by gomion plans).
5. Dual-support code paths for Charm v1 and v2 in the same branch.

## Measured Baseline (Pre-Expansion)

1. Go files (excluding `site`): **83**
2. `go.mod` files across repo: **19**
3. Files with key-message usage (`tea.KeyMsg`/`tea.KeyPressMsg`): **21**
4. Files with `View() string`: **10**
5. Files using removed v1 program options (`WithAltScreen`/mouse options): **12**
6. Baseline `make test`: green

---

## BASELINE — Baseline, Inventory, and Gate Definitions

1. Confirm current baseline on `main`:
   1. `make test`
   2. `make vet`
   3. `make build-examples` (if present)
2. Record module inventory and current versions.
3. Record current example run commands and expected smoke behavior.
4. Confirm and pin all linked plan artifact paths listed in "Required Linked Plans and Hard Gates."
5. Confirm all linked plans define explicit acceptance gates.
6. Record hard gate rule: no Charm v2 code changes start until baseline-v1 prerequisites are complete (BRANCH gate).
7. Record hard gate rule: no `v0.2.0` publish until all linked plans are completed.

**Acceptance gate:** Baseline outputs are captured and reproducible; linked plan artifact paths and gate semantics are unambiguous.

---

## EXPAND — Expand Charm v1 Baseline (`teanotify` and `teagrid`)

Introduce `teanotify` and `teagrid` as Charm v1-compatible baseline modules. These modules are independent of each other — their introduction, testing, and v2 migration can happen in parallel sessions.

1. Introduce `teanotify` as a Charm v1 baseline module:
   - Fork/rename from bubbleup source
   - License attribution, README, module setup (`teanotify/go.mod`)
   - Rename `AlertModel` -> `NotifyModel`, all "alert" -> "notify"/"notice" per `teanotify/PLAN.md`
2. Introduce `teagrid` as a Charm v1 baseline module:
   - Module skeleton (`teagrid/go.mod`, README)
   - Placeholder types following go-tealeaves patterns (ClearPath style, doterr errors)
3. Verify expanded baseline:
   1. `make test`
   2. `make vet`
   3. `make build-examples` (if present)
4. Refresh repo inventory counts post-expansion.

**Acceptance gate:** `teanotify` and `teagrid` exist and build as Charm v1 baselines; baseline validation checks are green.

---

## BRANCH — Baseline Tagging, Gomion Charm v1 Migration, and Branch Preservation

**Hard sequence rule:** No Charm v2 code changes before this phase's gate is satisfied.

1. Create baseline tags from the post-expansion baseline commit:
   1. Root `v0.1.0`
   2. Submodule `v0.1.0` tags for all existing modules and newly introduced modules (`teanotify`, `teagrid`).
2. Complete and validate `/Users/mikeschinkel/Projects/gomion/CHARM_V1_PLAN.md` using baseline tags.
3. Create `release/charm-v1` from the same baseline commit used for `v0.1.0` tagging.
4. Update README on `release/charm-v1` noting maintenance-only status.
5. Push `release/charm-v1` and tags before starting migration commits.
6. Keep `main` as Charm v2 migration target.

**Acceptance gate:** Baseline tags exist (`v0.1.0` root and submodule tags); `CHARM_V1_PLAN.md` acceptance gates are satisfied; `release/charm-v1` exists remotely and is documented as maintenance-only.

---

## TESTV1 — Comprehensive v1 Test Suite

Write tests against current v1 behavior **before** starting migration. These become the regression suite for v2 migration.

1. Identify all modules lacking test coverage.
2. Write unit tests exercising core public APIs for each module:
   - **teadd** — Model Init, Update (key handling), View output
   - **teastatus** — Model lifecycle and view rendering
   - **teadep** — Dependency dropdown interactions
   - **teatree** — Tree navigation, expand/collapse, viewport integration
   - **teamodal** — Modal display, choice selection, list selection, mouse handling
   - **teatextsel** — Text selection, shift-key handling, cursor movement
   - **teautils** — Helper function coverage
3. Focus on behaviors that will be affected by v2 migration:
   - Key message handling (all modules)
   - View output correctness (all modules)
   - Mouse interaction (teamodal)
   - Component integration (teatextsel textarea, teatree viewport)
4. Add/extend example smoke harnesses so examples are validated by build and run checks.
5. Ensure all tests pass on the v1 codebase before proceeding.

**Acceptance gate:** Meaningful test coverage exists for all modules; all tests pass on v1 codebase; example smoke validation is green; tests are committed to `main` (and cherry-picked to `release/charm-v1` if appropriate).

---

## IMPORTS — Mechanical Dependency Migration

For each module (`teadd`, `teautils`, `teastatus`, `teatextsel`, `teadep`, `teatree`, `teamodal`):

### go.mod import path changes

| Old Import | New Import |
|---|---|
| `github.com/charmbracelet/bubbletea` v1.3.x | `charm.land/bubbletea/v2` |
| `github.com/charmbracelet/bubbles` v0.21.0 | `charm.land/bubbles/v2` |
| `github.com/charmbracelet/lipgloss` v1.1.0 | `charm.land/lipgloss/v2` |
| `github.com/charmbracelet/x/ansi` v0.8.0 | **No change** — upstream v2 still uses this path |

### Module paths

No `/v2` suffix needed — modules are v0.x. Module paths stay as-is (e.g. `github.com/mikeschinkel/go-tealeaves/teadd`).

### Inter-module references

- `teadep` imports `teadd` (has `replace` directive)
- `teamodal` imports `teautils` (has `replace` directive)
- All examples have `replace` directives pointing to relative parent paths
- Update these only if module paths change (they shouldn't in this phase).

### Steps

1. Update all `go.mod` files with new Charm v2 dependency paths.
2. Replace import paths in all `.go` files (core + examples + tools).
3. `go get` to pull v2 deps; `go mod tidy` per module.
4. Regenerate `go.sum` files.

**Files:** `*/go.mod`, `*/go.sum`, all `.go` files with charm imports.

**Acceptance gate:** All modules resolve dependencies cleanly; `go mod tidy` succeeds per module; no module path changes to `/v2` were introduced.

---

## BTEA — Bubble Tea v2 API Migration

### VIEW — Adopt declarative tea.View (10+ files)

`View() string` -> `View() tea.View`. Return `tea.NewView(content)`. Named returns `(view string)` -> `(view tea.View)` for ClearPath style.

Where modules set alt screen or mouse mode imperatively via commands (`WithAltScreen`, `WithMouseCellMotion`), move to View struct fields:
```go
v.AltScreen = true
v.MouseMode = tea.MouseModeCellMotion
```

**Per-module:**

| Module | File | Current Signature |
|---|---|---|
| teadd | `model.go:175` | `View() (view string)` |
| teastatus | `model.go:68` | `View() (view string)` |
| teadep | `model.go:190` | `View() (view string)` |
| teatree | `model.go:111` | `View() string` |
| teamodal | `model.go:297` | `View() (view string)` |
| teamodal | `choice_model.go` | `View() (view string)` |
| teamodal | `list_model.go` | `View() (view string)` |
| teatextsel | `view.go` | View called from model |

### KEY — Migrate tea.KeyMsg to tea.KeyPressMsg (21+ files)

This is the highest-touch change. Field-level migration:

| v1 | v2 |
|---|---|
| `tea.KeyMsg` | `tea.KeyPressMsg` |
| `msg.Type` | `msg.Code` (rune) |
| `msg.Runes` | `msg.Text` (string, not []rune) |
| `msg.Alt` | `msg.Mod.Contains(tea.ModAlt)` |
| `msg.String()` returns `" "` for space | Returns `"space"` |
| `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}` | New v2 key construction |

**Per-module impact:**

- **teadd/model.go** — `msg.(tea.KeyMsg)` type assertion + `key.Matches()`. Low complexity.
- **teatextsel/model.go** — **Heaviest migration.** Uses `keyMsg.Type`, `keyMsg.Runes`, `keyMsg.String()`, `len(msg.Runes)`, `msg.Type == tea.KeySpace/KeyEnter`. Shift detection via `msg.String()` contains `"shift"` -> use `msg.Mod.Contains(tea.ModShift)` instead.
- **teastatus/model.go** — `tea.KeyMsg` assertion. Low complexity.
- **teadep/model.go** — `msg.(tea.KeyMsg)` + `key.Matches()`. Low complexity.
- **teatree/model.go** — `tea.KeyMsg` switch + `key.Matches()`. Medium complexity.
- **teamodal/model.go** — `tea.KeyMsg` in Update. Low complexity.
- **teamodal/choice_model.go** — `keyMsg.Type == tea.KeyRunes && len(keyMsg.Runes) == 1` + `keyMsg.Runes[0]` -> migrate to `msg.Text` and rune access via `[]rune(msg.Text)[0]` or `msg.Code`.
- **teamodal/list_model.go** — `keyMsg.Type` switch + `keyMsg.Runes` access. Medium complexity.
- **teamodal/choice_model_test.go** — Test key construction: `tea.KeyMsg{Type: tea.KeyTab}` -> v2 equivalent.
- **All example `main.go` files** (8+) — `msg.String()` patterns.
- **cmd/color-viewer/main.go** — `msg.String()` pattern.

### MOUSE — Migrate tea.MouseMsg (2 files)

| v1 | v2 |
|---|---|
| `tea.MouseMsg` (struct) | Interface; use `msg.Mouse()` for coords |
| `switch mouseMsg.Type` | Type-switch: `tea.MouseClickMsg`, `tea.MouseReleaseMsg`, `tea.MouseWheelMsg`, `tea.MouseMotionMsg` |
| `MouseButtonLeft` | `MouseLeft` |
| Imperative mouse mode commands | `v.MouseMode = tea.MouseModeCellMotion` in View |

**Files:**
- **teamodal/model.go:269** — `switch mouseMsg.Type` -> type-switch on specific msg types.
- **teamodal/choice_model.go** — if mouse handling exists.

**Acceptance gate:** All Bubble Tea API usage is v2-correct; core modules compile on Charm v2 APIs.

---

## BUBLIP — Bubbles and Lip Gloss v2 Migration

### Lip Gloss v2 changes (17 files)

| v1 | v2 |
|---|---|
| `lipgloss.Color` (type) | `lipgloss.Color()` (function -> `color.Color`) |
| `TerminalColor` interface | `color.Color` (from `image/color`) |
| `lipgloss.NewStyle()` tied to renderer | Pure value; no renderer |
| Renderer type exists | **Removed entirely** |
| `WithWhitespaceForeground` + `WithWhitespaceBackground` | `WithWhitespaceStyle()` |
| `AdaptiveColor` in root package | Moved to `compat` package (NOT USED -- no impact) |

Note: `lipgloss.Color("86")` syntax is the same in both versions, but the return type changes from `lipgloss.Color` to `color.Color`. Usage in `.Foreground()` etc. is compatible since those now accept `color.Color`.

Add `import "image/color"` where color values are stored as typed fields.

**Files with lipgloss.Color() calls:**
- `teadd/styles.go` — colors "240", "15", "62", "230"
- `teautils/render_help_visor.go` — colors "99", "178", "86", "252"
- `teastatus/styles.go` — colors "86", "246", "240"
- `teatextsel/view.go` — colors "39", "232"

### Bubbles v2 component changes

| v1 | v2 |
|---|---|
| `DefaultKeyMap` (variable) | `DefaultKeyMap()` (function) |
| `m.Width` (field) | `m.SetWidth()` / `m.Width()` (methods) |
| `NewModel()` | `New()` |
| `FocusedStyle` / `BlurredStyle` | `Styles.Focused` / `Styles.Blurred` |
| `model.Cursor` (virtual) | `model.Cursor()` (real `*tea.Cursor`) |
| `SetCursor()` | `SetCursorColumn()` |

**Per-module:**
- **teatextsel** — Wraps `textarea.Model`. Uses `.FocusedStyle.CursorLine.Render()`, `.BlurredStyle.Base.Render()` -> migrate to `Styles.Focused` / `Styles.Blurred`. Cursor handling changes.
- **teatree** — Uses `viewport` -> width/height now methods.
- **teamodal** — May use viewport or other bubbles components.
- **teadep** — Uses dropdown which wraps bubbles.

**Acceptance gate:** Behavior parity for existing components; core component packages compile and behave equivalently on v2 APIs. No known regressions in core module interaction patterns.

---

## VALIDATE — Test Migration and Strict Example Validation

### Tests

1. Migrate tests that construct key messages explicitly (notably `teamodal/choice_model_test.go`).
2. Update v1 test suite (from TESTV1) for v2 API changes.
3. Run all tests after migration: `go test ./...` in each module.
4. Build check: `go build -o /dev/null ./...` in each module.

### Examples

1. Migrate all example apps in `examples/*/main.go` to v2 APIs.
2. Migrate `cmd/color-viewer/main.go` to v2 APIs.
3. Ensure all examples compile.
4. Run smoke checks: non-interactive where possible; scripted interaction where required.

**Acceptance gate:** All tests pass; all examples build; all examples run successfully under defined smoke criteria.

---

## DOCS — Documentation

1. Create `UPGRADE_GUIDE_V2.md` for external users covering:
   - Support policy and branch expectations (`release/charm-v1` best-effort, `main` primary)
   - Old/new import paths (no module path changes, only Charm dependency paths)
   - Package-by-package breaking changes with before/after code examples
   - Migration checklist
   - Explicit note that module paths remain unsuffixed while project is pre-v1
2. Update package READMEs and root README for Charm v2 usage.

**Acceptance gate:** A user can migrate from Charm v1 to Charm v2 usage of go-tealeaves using repo docs alone; compatibility/versioning policy is unambiguous.

---

## SITE — Documentation Site Update

Implement the plan defined in `site/UPGRADE_V2_PLAN.md`. This upgrades the Astro + Starlight documentation site from placeholder content to real go-tealeaves documentation with v2-correct API references and examples.

Key activities:
1. Remove Starlight starter placeholder content (SITEBASE).
2. Update Astro/Starlight config and sidebar structure (SITECONFIG).
3. Write module reference pages, getting started guide, and migration guide (SITECONTENT).
4. Finalize sidebar navigation (SITESIDEBAR).
5. Verify build, deploy, and live site (SITEVERIFY).

SITEBASE and SITECONFIG can begin during earlier phases. SITECONTENT is blocked on DOCS (requires `UPGRADE_GUIDE_V2.md` and finalized v2 APIs). SITEVERIFY is the final step before RELEASE.

**Acceptance gate:** `site/UPGRADE_V2_PLAN.md` acceptance gates are all satisfied; live site is deployed and fully navigable.

---

## RELEASE — Release Readiness and Synchronized `v0.2.0` Publish

1. Run final validation suite:
   1. `make test`
   2. `make vet`
   3. `make build-examples`
   4. Example run smoke suite
2. **Hard gate:** Verify `release/charm-v1` still builds/tests cleanly (no v2 changes leaked). Publishing is blocked if `release/charm-v1` is broken.
3. Verify all linked plans are completed per their own acceptance gates:
   1. `teautils/THEMING_PLAN.md`
   2. `teanotify/PLAN.md`
   3. `teagrid/PLAN.md`
   4. `gomion/CHARM_V2_PLAN.md`
   5. `site/UPGRADE_V2_PLAN.md`
4. Confirm `teanotify` and `teagrid` are fully implemented and validated on Charm v2 before release.
5. Cut synchronized tags:
   1. Root `v0.2.0`
   2. Submodule `v0.2.0` tags for all Charm v2-ready modules, including `teanotify` and `teagrid`.
6. Publish concise release notes referencing `UPGRADE_GUIDE_V2.md`.

**Acceptance gate:** Validation is fully green; `release/charm-v1` is green; all linked plan gates are satisfied; synchronized `v0.2.0` tags are published.

---

## Parallelism Notes

The following work streams can run in **independent parallel sessions** after BRANCH's gate is satisfied:

- **Core migration** (IMPORTS through DOCS of this plan) — existing module v2 migration
- **`teanotify` v2 implementation** — per `teanotify/PLAN.md`
- **`teagrid` v2 implementation** — per `teagrid/PLAN.md`
- **Theming system** — per `teautils/THEMING_PLAN.md`
- **Gomion Charm v2 migration** — per `gomion/CHARM_V2_PLAN.md`
- **Documentation site** (SITEBASE/SITECONFIG of `site/UPGRADE_V2_PLAN.md`) — cleanup and config can start early; SITECONTENT blocked on DOCS

All streams must complete before RELEASE (synchronized `v0.2.0` publish).

---

## Commit Strategy

Use small, reviewable commits by phase:

- **BASELINE** — Planning/gate updates and baseline expansion prep
- **EXPAND** — `teanotify`/`teagrid` baseline introduction
- **BRANCH** — Baseline tagging and branch preservation (`release/charm-v1`)
- **TESTV1** — Comprehensive v1 test suite
- **IMPORTS** — Mechanical dependency migration (go.mod + .go imports)
- **BTEA** — Bubble Tea API migration (View, KeyMsg, MouseMsg)
- **BUBLIP** — Bubbles/Lip Gloss component fixes
- **VALIDATE** — Tests and examples migration + smoke validation
- **DOCS** — Migration docs (`UPGRADE_GUIDE_V2.md` + README updates)
- **SITE** — Documentation site update (per `site/UPGRADE_V2_PLAN.md`)
- **RELEASE** — Final gate verification and synchronized release tagging

This keeps regressions isolated and makes rollback/cherry-pick straightforward.

---

## Risk Notes

1. **`teatextsel` is highest risk** due to textarea API changes (style restructuring, cursor API, keymap) plus custom selection rendering and heavy KeyMsg usage with shift detection.
2. **`teamodal` is second highest risk** due to key and mouse semantics changes across model.go, choice_model.go, and list_model.go.
3. **`teatree` and `cmd/color-viewer`** are likely straightforward but touch viewport API changes.
4. **Multi-module dependency updates** must be done carefully to avoid broken intra-repo `replace` directives. Process all 19 go.mod files systematically.
5. **Inconsistent per-module version bumps or mixed naming** (`teatable` vs `teagrid`) — mitigate with explicit tag matrix and consistent naming enforcement before tagging.
6. **Example runtime behavior drift** despite compile success — mitigate with strict smoke testing in VALIDATE.
7. **Cross-plan dependency risk** — `v0.2.0` blocked by unfinished theming or gomion migration plans. Mitigate by tracking linked plan progress and identifying blockers early.

---

## Explicit Deferrals

- **`teanotify` v2 implementation** — owned by `teanotify/PLAN.md`. Baseline introduction is in EXPAND; v2 migration details are in the linked plan.
- **`teagrid` v2 implementation** — owned by `teagrid/PLAN.md`. Baseline introduction is in EXPAND; v2 migration details are in the linked plan. Will replace `bubble-table` (`github.com/evertras/bubble-table`) and be built from scratch using lipgloss v2.
- **Theming system** — owned by `teautils/THEMING_PLAN.md`. Out of scope for this plan's implementation steps.
- **Gomion migrations** — owned by `gomion/CHARM_V1_PLAN.md` and `gomion/CHARM_V2_PLAN.md`. Gomion v1 is a hard gate for BRANCH; gomion v2 is a hard gate for RELEASE.

---

## Reference: Upstream Upgrade Guides

- [Bubble Tea v2: What's New](https://github.com/charmbracelet/bubbletea/discussions/1374)
- [Bubble Tea v2.0.0 Release](https://github.com/charmbracelet/bubbletea/releases/tag/v2.0.0)
- [Bubble Tea v2 Upgrade Guide](https://github.com/charmbracelet/bubbletea/blob/v2.0.0/UPGRADE_GUIDE_V2.md)
- [Lip Gloss v2: What's New](https://github.com/charmbracelet/lipgloss/discussions/506)
- [Bubbles v2 Upgrade Guide](https://github.com/charmbracelet/bubbles/blob/v2.0.0/UPGRADE_GUIDE_V2.md)
- [Charm v2 Blog Post](https://charm.land/blog/v2/)
- [Bubble Tea v2 go.mod](https://raw.githubusercontent.com/charmbracelet/bubbletea/v2.0.0/go.mod) (confirms x/ansi path)
- [Bubbles v2 go.mod](https://raw.githubusercontent.com/charmbracelet/bubbles/v2.0.0/go.mod) (confirms x/ansi path)
