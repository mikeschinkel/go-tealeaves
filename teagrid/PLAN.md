# Plan: Implement `teagrid` (Bubble Table Successor)

## Context

`teagrid` replaces `github.com/evertras/bubble-table/table` usage in go-tealeaves and gomion.

This plan must satisfy both tracks:
1. Introduce a Charm v1-compatible baseline module and tag `teagrid/v0.1.0`.
2. Fully migrate `teagrid` to Charm v2 and validate it before synchronized `v0.2.0`.

This plan is a required dependency of:
1. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md`
2. `/Users/mikeschinkel/Projects/gomion/CHARM_V1_PLAN.md`
3. `/Users/mikeschinkel/Projects/gomion/CHARM_V2_PLAN.md`
4. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teautils/THEMING_PLAN.md`
5. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teanotify/PLAN.md`
6. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/site/UPGRADE_V2_PLAN.md`

Completion means this plan's acceptance gates are satisfied, not merely that this file exists.

## Hard Pre-v2 Checkpoint

The following outcomes are mandatory before any v2 implementation work starts:
1. Fork `bubble-table` into `teagrid`.
2. Update code/docs/module naming to `teagrid`.
3. Apply MIT attribution and license requirements for forked content.
4. Update `/Users/mikeschinkel/Projects/gomion/CHARM_V1_PLAN.md` to use `teagrid` instead of `bubble-table` and `teanotify` instead of `bubbleup`.
5. Ensure gomion compiles and runs with `teagrid` and `teanotify` baseline integrations.
6. Publish baseline tags `teagrid/v0.1.0` and `teanotify/v0.1.0`.
7. Verify gomion against published baseline tags.

Pragmatic sequencing note:
1. Steps 5 and 6 can iterate in whichever order is fastest while stabilizing baseline.
2. The hard requirement is that both are complete (including post-tag verification) before v2 work begins.

## Execution Order

Phase names are mnemonic and not ordered alphabetically. Execution order is:
1. `ALIGN`
2. `FORK_REBRAND_ATTRIBUTION`
3. `BASELINE_V1`
4. `TEST_V1`
5. `GOMION_V1_PLAN_SYNC`
6. `BASELINE_TAG_V1`
7. `GOMION_V1_EXECUTE`
8. `V2_DECISIONS`
9. `MIGRATE_V2`
10. `HARDEN`
11. `EXAMPLES`
12. `DOCS`
13. `RELEASE_GATES`

## Root-Linked Blockers (Track, Don't Stall)

For synchronized `v0.2.0`, root plan blocks on completion of:
1. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md`
2. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teautils/THEMING_PLAN.md`
3. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teanotify/PLAN.md`
4. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teagrid/PLAN.md`
5. `/Users/mikeschinkel/Projects/gomion/CHARM_V1_PLAN.md`
6. `/Users/mikeschinkel/Projects/gomion/CHARM_V2_PLAN.md`
7. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/site/UPGRADE_V2_PLAN.md`

Execution intent:
1. Keep these visible and linked.
2. Avoid process-heavy blocker bookkeeping that slows delivery.
3. Prioritize steady implementation progress while preserving release gates.

## Root Plan Mapping

Mapping to `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md`:
1. Root `EXPAND` dependency:
   1. `FORK_REBRAND_ATTRIBUTION`
   2. `BASELINE_V1`
   3. `TEST_V1` (module-level validation support)
2. Root `BRANCH` dependency:
   1. `GOMION_V1_PLAN_SYNC`
   2. `GOMION_V1_EXECUTE`
   3. `BASELINE_TAG_V1`
3. Root post-`BRANCH` v2 stream:
   1. `V2_DECISIONS`
   2. `MIGRATE_V2`
   3. `HARDEN`
   4. `EXAMPLES`
   5. `DOCS`
4. Root `RELEASE` dependency:
   1. `RELEASE_GATES` completion
   2. linked-plan acceptance completion
   3. `release/charm-v1` health remains green at release verification

## Locked Requirements

1. Module name is `teagrid` (renamed from the previously proposed `teatable`).
2. Baseline v1 introduction is a fork/rename/attribution flow first, not a greenfield rewrite.
3. Baseline and v2 tracks both preserve functional behavior for gomion's table workflows.
4. `teagrid` must include the fork-specific cell-cursor capabilities gomion relies on.
5. Preserve MIT attribution from upstream `bubble-table`.
6. Use mnemonic phase naming (no numeric phase identifiers).
7. For Charm v2 migration, default to v2-native conventions first.
8. Compatibility fallbacks are allowed only by explicit per-item decisions, not by blanket policy.
9. No Charm v2 implementation starts until the "Hard Pre-v2 Checkpoint" is completed.
10. Immediate priority is rapid delivery of v1 baseline `teagrid` so gomion can migrate off `bubble-table` on Charm v1.
11. v1 baseline implementation is minimal-touch rename/rebrand/attribution work; avoid refactors, feature changes, or cleanup churn.
12. v1 baseline test coverage must be thorough even though implementation changes are minimal.

## Scope

In scope:
1. Create `teagrid` as a standalone go-tealeaves module.
2. Port the current local `bubble-table` fork behavior into `teagrid`.
3. Add/maintain thorough v1 validation and regression coverage for safe gomion adoption and baseline tagging.
4. Add runnable examples and smoke-validation commands.
5. Provide migration docs from `bubble-table/table` to `teagrid`.
6. Include v1 gomion plan update/implementation linkage as hard release gates for baseline tagging.

Out of scope:
1. Replacing `teatree` or changing tree APIs.
2. New table features unrelated to parity and v2 migration.
3. Detailed gomion implementation task breakdown (owned by gomion plan artifact).

## Source Inventory to Port

Primary source of truth for baseline:
1. `/Users/mikeschinkel/Projects/go-pkgs/bubble-table`
2. Package directory: `/Users/mikeschinkel/Projects/go-pkgs/bubble-table/table`

Upstream attribution source:
1. `github.com/evertras/bubble-table` (MIT)

Key v1 fork-only behavior to preserve:
1. `WithCellCursorMode(bool)`
2. `GetCellCursorMode() bool`
3. `GetCellCursorColumnIndex() int`
4. `GetVisibleColumnRange() (start, end int)`
5. Cell-cursor-aware left/right behavior in `Update()`

Baseline file mapping (fork -> module):
1. `table/model.go` -> `teagrid/model.go`
2. `table/update.go` -> `teagrid/update.go`
3. `table/view.go` -> `teagrid/view.go`
4. `table/options.go` -> `teagrid/options.go`
5. `table/query.go` -> `teagrid/query.go`
6. `table/keys.go` -> `teagrid/keys.go`
7. `table/events.go` -> `teagrid/events.go`
8. `table/column.go` -> `teagrid/column.go`
9. `table/row.go` -> `teagrid/row.go`
10. `table/cell.go` -> `teagrid/cell.go`
11. `table/filter.go` -> `teagrid/filter.go`
12. `table/sort.go` -> `teagrid/sort.go`
13. `table/pagination.go` -> `teagrid/pagination.go`
14. `table/scrolling.go` and `table/overflow.go` -> `teagrid/` equivalents
15. `table/header.go` and `table/footer.go` -> `teagrid/` equivalents
16. Utility/support files (`calc.go`, `data.go`, `dimensions.go`, `strlimit.go`, `border.go`, `doc.go`) -> `teagrid/` equivalents

## Gomion Dependency Contract

Current gomion usage baseline (must remain supported through migration):
1. Constructs and data:
   1. `New`, `NewColumn`, `NewFlexColumn`
   2. `Row`, `RowData`, `NewRow`, `StyledCell`, `NewStyledCell`
2. Interactive behavior:
   1. `Focused(true)`
   2. `WithCellCursorMode(true)`
   3. `WithMaxTotalWidth`, `WithHorizontalFreezeColumnCount`
   4. `WithPageSize`, `WithMinimumHeight`, `WithFooterVisibility(false)`
3. Styling and display:
   1. `HeaderStyle`, `HighlightStyle`, `WithBaseStyle`, `BorderRounded`
4. Runtime query hooks:
   1. `GetHighlightedRowIndex()`
   2. `GetCellCursorColumnIndex()`
5. Update cycle:
   1. `Update(msg)` row/cell movement semantics to preserve unless superseded by explicit `V2_DECISIONS`
6. Rendering:
   1. `View()` composability expectations inside parent model views

## Public API Target

API compatibility target for `v0.1.0`:
1. Import-path migration only for typical consumers:
   1. From `github.com/evertras/bubble-table/table`
   2. To `github.com/mikeschinkel/go-tealeaves/teagrid`
2. Preserve core exported surface semantics:
   1. `Model`, `Column`, `Row`, `RowData`, `StyledCell`
   2. `New`, `NewColumn`, `NewFlexColumn`, `NewRow`, `NewStyledCell`, `NewStyledCellWithStyleFunc`
   3. Existing `With*` fluent methods on `Model`, `Column`, `Row`
   4. Keymap, events, sorting, filtering, pagination, scrolling APIs
3. Preserve fork extensions listed above.
4. Prefer zero functional deltas from the current fork during v1 baseline (rename-first objective).

Compatibility note:
1. Pre-v1 release allows breaking changes, but this plan aims for practical compatibility to reduce gomion migration risk.

## Compatibility Decision Register (Case-by-Case)

For each possible fallback to legacy conventions, capture:
1. Topic
2. v2-native default
3. Proposed fallback behavior
4. Reason to include fallback
5. Reason to exclude fallback
6. Decision (`accept fallback`, `reject fallback`, `defer`)
7. Validation impact (tests/examples/docs updates required)

Candidate topics to discuss one by one:
1. `View()` contract details (`tea.View`-first integration boundaries)
2. Key handling compatibility aliases (`tea.KeyPressMsg` vs legacy expectations)
3. Legacy fluent API aliases or renamed methods (if any)
4. Behavior toggles where v2 defaults differ from v1-era expectations
5. Implementation strategy for v2 (`greenfield` vs `evolve-from-fork`)

## PHASE `ALIGN` (Contract and Dependency Alignment)

1. Lock source baseline commit from local fork (`go-pkgs/bubble-table`).
2. Lock package-layout decision for `teagrid`:
   1. package name `teagrid`
   2. import path `github.com/mikeschinkel/go-tealeaves/teagrid`
3. Confirm naming-collision handling with gomion:
   1. gomion currently maps `go-tealeaves/teagrid` to `teatree` via `replace`
   2. gomion v1 migration plan artifact must define the transition step and ordering
4. Define parity-critical behaviors:
   1. Row navigation and highlight
   2. Cell cursor mode and auto horizontal scroll
   3. Filtering with built-in input
   4. Sorting and pagination
   5. Max width overflow and frozen columns

Acceptance gate:
1. API contract, baseline source commit, and parity-critical behavior list are documented and unambiguous.

## PHASE `FORK_REBRAND_ATTRIBUTION` (Mandatory Pre-v2)

1. Fork/copy source from local `bubble-table` baseline into `teagrid`.
2. Rename package/module/docs references from `table`/`bubble-table` to `teagrid` where required.
3. Ensure MIT attribution is preserved and explicit:
   1. `LICENSE` inclusion
   2. upstream reference in README/docs
   3. any additional NOTICE/attribution text needed by repo standards

Acceptance gate:
1. `teagrid` exists as a forked baseline module with naming and attribution completed.

## PHASE `GOMION_V1_PLAN_SYNC` (Mandatory Pre-v2)

1. Update gomion v1 migration plan artifact to include:
   1. replace `bubble-table` usage with `teagrid`
   2. replace `bubbleup` usage with `teanotify`
2. Confirm sequencing in gomion plan matches this repo's baseline tagging and branch strategy.

Acceptance gate:
1. Gomion v1 plan artifact is updated and explicitly references `teagrid` + `teanotify` baseline usage.

## PHASE `GOMION_V1_EXECUTE` (Mandatory Pre-v2)

1. Execute gomion v1 plan in its own session(s).
2. Verify gomion compiles with `teagrid` + `teanotify` baseline modules.
3. Verify gomion runs successfully on the targeted baseline integration path.

Acceptance gate:
1. Gomion compiles and runs with `teagrid` + `teanotify` baseline integrations.

## PHASE `V2_DECISIONS` (Feature Changes and Compatibility Deliberation)

Hard gate:
1. This phase is blocked until `FORK_REBRAND_ATTRIBUTION`, `BASELINE_V1`, `TEST_V1`, `GOMION_V1_PLAN_SYNC`, `GOMION_V1_EXECUTE`, and `BASELINE_TAG_V1` are complete.

Starting position for v2 implementation strategy:
1. Target is Charm v2 best-practices behavior and architecture.
2. Evolve from baseline where that still yields best-practices outcomes.
3. Rewrite greenfield where evolution cannot credibly achieve best-practices quality.
4. Final per-area strategy decision is deferred until after baseline `v0.1.0` delivery and gomion v1 adoption validation.

1. Build the v2 change-set decision log with three buckets:
   1. features to remove
   2. features to add
   3. behaviors to make configurable
2. For each candidate compatibility fallback, evaluate using the register above.
3. Record final scope for v2 implementation:
   1. mandatory v2-native behavior
   2. accepted compatibility fallbacks
   3. rejected compatibility fallbacks
   4. evolve-vs-rewrite decision per subsystem
4. Convert accepted decisions into actionable implementation tasks in this plan.

Acceptance gate:
1. Decision log is complete for all known v2 change topics.
2. No compatibility fallback remains implied or undecided.
3. Evolve-vs-rewrite decisions are explicit for major subsystems.
4. `MIGRATE_V2` scope is explicit and implementation-ready.

## PHASE `BASELINE_V1` (Introduce Module + `v0.1.0`)

1. Create module skeleton:
   1. `teagrid/go.mod`
   2. production files copied/adapted from `bubble-table/table/*.go`
   3. `teagrid/README.md`
   4. `teagrid/LICENSE` and attribution text
2. Rename package declarations:
   1. `package table` -> `package teagrid`
3. Keep public API equivalent to baseline fork where feasible.
4. Ensure local module builds cleanly on Charm v1 dependencies.
5. Ensure repo-level integration updates include `teagrid` where module lists are enumerated (for example `Makefile` module/test loops).
6. Do not perform elective refactors, behavior cleanups, or feature redesign in this phase.
7. Allow test-only changes needed to support thorough baseline verification.

Acceptance gate:
1. `teagrid` module exists, compiles, and is functionally baseline-compatible with local fork behavior.
2. Module is ready for baseline tagging as `teagrid/v0.1.0`.

## PHASE `TEST_V1` (Pin Baseline Behavior)

1. Build and execute a thorough v1 baseline test suite before tagging:
   1. module-level unit and behavior tests for parity-critical table behavior
   2. regression tests for fork-specific cell-cursor functionality
   3. gomion integration-facing checks for the APIs it uses
2. Verify gomion-critical API availability after rename:
   1. `WithCellCursorMode(true)`
   2. `GetHighlightedRowIndex()` and `GetCellCursorColumnIndex()`
   3. `WithMaxTotalWidth` and `WithHorizontalFreezeColumnCount`
3. Treat detailed test case inventory as an input artifact from the parallel test-planning session.
4. Keep production implementation minimal-touch; broaden testing rather than broadening baseline code changes.

Acceptance gate:
1. Thorough v1 tests pass and gomion can compile/run against renamed module.
2. Baseline readiness is blocked until required test artifacts from the parallel test-planning session are incorporated.

## PHASE `BASELINE_TAG_V1` (Mandatory Pre-v2)

1. Cut and push baseline tags after baseline and gomion-v1 validation are complete:
   1. `teagrid/v0.1.0`
   2. `teanotify/v0.1.0`
2. Confirm both tags resolve to validated baseline commits.
3. Record tag references in release notes/plans for downstream dependency pinning.

Acceptance gate:
1. Baseline tags for `teagrid` and `teanotify` exist remotely at `v0.1.0`.

## PHASE `MIGRATE_V2` (Charm v2 Dependencies and API Migration)

Hard gate:
1. This phase starts only after `BASELINE_TAG_V1` and `V2_DECISIONS` are complete.

1. Update dependency imports:
   1. `github.com/charmbracelet/bubbletea` -> `charm.land/bubbletea/v2`
   2. `github.com/charmbracelet/bubbles` -> `charm.land/bubbles/v2`
   3. `github.com/charmbracelet/lipgloss` -> `charm.land/lipgloss/v2`
2. Migrate key handling:
   1. `tea.KeyMsg` -> `tea.KeyPressMsg`
   2. Keep key-map matching semantics equivalent.
3. Migrate textinput integration to v2 API.
4. Implement v2-native rendering and input conventions as default behavior.
5. Apply only those fallback behaviors explicitly accepted in `V2_DECISIONS`.
6. Re-run test suite and resolve regressions before proceeding.

Acceptance gate:
1. `teagrid` compiles cleanly on Charm v2 APIs.
2. v1 parity tests are migrated and passing on v2.
3. No undocumented regressions in parity-critical behavior.

## PHASE `HARDEN` (Behavior and Quality Stabilization)

1. Ensure immutable/value-return APIs remain coherent for all `With*` methods.
2. Validate no hidden mutation hazards in copied slices/maps where immutability is expected.
3. Confirm width, border, and overflow behavior is deterministic (no empirical "magic constants").
4. Keep implementation lean: avoid new features unless required for parity or migration correctness.

Acceptance gate:
1. `teagrid` behavior is stable under tests for layout, scrolling, filtering, sorting, and cell-cursor workflows.

## PHASE `EXAMPLES` (Runnable Validation)

1. Add `teagrid/examples/` with focused runnable apps:
   1. simplest rendering
   2. filtering
   3. pagination
   4. horizontal scrolling + frozen columns
   5. cell cursor mode
2. Ensure examples build and run as smoke tests in CI/local workflow.
3. Prefer example coverage that matches gomion usage patterns and interaction model.

Acceptance gate:
1. All `teagrid` examples build and run under defined smoke criteria.
2. Example behavior reflects v2 APIs and current `teagrid` contracts.

## PHASE `DOCS` (Migration and Usage Documentation)

1. Write `teagrid/README.md` with:
   1. quick-start table example
   2. key APIs (`New`, rows/columns, filtering/sorting/pagination)
   3. cell cursor mode usage
2. Add migration section:
   1. import-path changes (`bubble-table/table` -> `teagrid`)
   2. compatibility notes and any intentional deltas
3. Update repo-level docs references where `teagrid` is listed as a required module.

Acceptance gate:
1. A consumer can migrate from `bubble-table` to `teagrid` with repo docs only.

## PHASE `RELEASE_GATES` (Tag and Cross-Plan Completion)

1. Baseline release track:
   1. `ALIGN`, `FORK_REBRAND_ATTRIBUTION`, `BASELINE_V1`, `TEST_V1`, `GOMION_V1_PLAN_SYNC`, `GOMION_V1_EXECUTE`, `BASELINE_TAG_V1` gates green
   2. `teagrid/v0.1.0` and `teanotify/v0.1.0` tags are published
2. Charm v2 release track:
   1. `V2_DECISIONS`, `MIGRATE_V2`, `HARDEN`, `EXAMPLES`, and `DOCS` gates green
   2. `teagrid` fully implemented and validated on Charm v2
3. Cross-plan completion:
   1. root plan and gomion linked plans satisfy their own acceptance criteria
   2. gomion plan explicitly resolves `teagrid => teatree` replace ambiguity before final migration sign-off
   3. explicitly unblock synchronized `v0.2.0` only after all dependent plans complete
4. Root branch health rule:
   1. verify `release/charm-v1` remains build/test green before synchronized `v0.2.0` publish

Acceptance gate:
1. `teagrid` no longer blocks synchronized go-tealeaves `v0.2.0` release.

## Validation Commands

From `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teagrid`:
1. `go build -o /dev/null ./...`
2. `go test ./...` (including newly-added baseline coverage)
3. `go vet ./...`

Example smoke checks:
1. Build all `teagrid/examples/*`.
2. Run scripted/timeout smoke executions for each example.

## Risks and Mitigations

1. Risk: gomion naming collision (`teagrid` currently aliased to `teatree`).
   Mitigation: treat gomion `CHARM_V1_PLAN.md` and `CHARM_V2_PLAN.md` as explicit sequencing gates.
2. Risk: fork-only cell-cursor behavior regresses during migration.
   Mitigation: keep baseline behavior unchanged in v1 and validate through gomion v1 compile/run gates.
3. Risk: layout/width regressions from lipgloss/bubbles v2 behavior differences.
   Mitigation: preserve render tests and add targeted overflow/frozen-column assertions.
4. Risk: examples compile but interaction behavior drifts.
   Mitigation: include interaction-oriented smoke tests (cursor mode, scrolling, filtering) in validation.
5. Risk: v2 work starts prematurely before baseline tag and gomion v1 validation.
   Mitigation: explicit hard gate in `V2_DECISIONS` and `MIGRATE_V2`.
6. Risk: baseline scope creep delays gomion migration.
   Mitigation: enforce rename-first minimal-touch rule and defer redesign work to v2 phases.
7. Risk: insufficient baseline regression coverage causes hidden breakage in gomion migration.
   Mitigation: require thorough v1 test suite completion as a hard gate before `v0.1.0` tagging.

## Completion Checklist

1. `ALIGN` gate satisfied.
2. `FORK_REBRAND_ATTRIBUTION` gate satisfied.
3. `BASELINE_V1` gate satisfied.
4. `TEST_V1` gate satisfied.
5. `GOMION_V1_PLAN_SYNC` gate satisfied.
6. `GOMION_V1_EXECUTE` gate satisfied.
7. `BASELINE_TAG_V1` gate satisfied.
8. `V2_DECISIONS` gate satisfied.
9. `MIGRATE_V2` gate satisfied.
10. `HARDEN` gate satisfied.
11. `EXAMPLES` gate satisfied.
12. `DOCS` gate satisfied.
13. `RELEASE_GATES` gate satisfied.
