# Plan: Implement `teanotify` (BubbleUp Successor)

## Context

`teanotify` replaces `go.dalton.dog/bubbleup` inside go-tealeaves.

This plan must satisfy both tracks:
1. Introduce a Charm v1-compatible baseline module and tag `teanotify/v0.1.0`.
2. Fully migrate `teanotify` to Charm v2 and validate it before synchronized `v0.2.0` release.

This plan is a hard-gated dependency of:
1. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md`
2. `/Users/mikeschinkel/Projects/gomion/CHARM_V1_PLAN.md`
3. `/Users/mikeschinkel/Projects/gomion/CHARM_V2_PLAN.md`

Completion means this plan's acceptance gates are satisfied, not merely that this file exists.

## Locked Requirements

1. Rename `AlertModel` to `NotifyModel`.
2. Rename alert terminology to notify/notice terminology throughout API and docs.
3. Constructor shape is `NewNotifyModel(opts NotifyOpts)` with a single options parameter.
4. Keep relevant immutable `With*()` methods for post-instantiation updates.
5. Follow go-tealeaves house rules for production Go code:
   1. ClearPath style.
   2. doterr-based error handling (no `fmt.Errorf` patterning for production flows).
   3. No ignored errors.
6. Keep phase naming mnemonic (no numeric phase ordering).
7. Preserve MIT attribution from upstream BubbleUp.

## Scope

In scope:
1. Port core BubbleUp behavior and API into `teanotify`.
2. Provide both v1 baseline and v2 implementation tracks.
3. Add comprehensive `_test.go` coverage and example smoke validation.
4. Provide migration docs from BubbleUp to TeaNotify.

Out of scope:
1. New product features beyond parity and required API reshaping.
2. Theming-system integration details owned by `teautils/THEMING_PLAN.md`.
3. Gomion implementation details (owned by gomion plans).

## Source Inventory to Port

Upstream source repository:
1. `~/Projects/go-3rd-party/bubbleup`

Primary file mapping:
1. `alert.go` -> `notice.go`
2. `model.go` -> `model.go`
3. `position.go` -> `position.go`
4. `utils.go` -> `util.go`

## Public API Target

Symbol rename map:

| BubbleUp | TeaNotify |
|---|---|
| `AlertModel` | `NotifyModel` |
| `AlertDefinition` | `NoticeDefinition` |
| `NewAlertModel(width, useNerdFont, duration)` | `NewNotifyModel(opts NotifyOpts)` |
| `NewAlertCmd(alertType, message)` | `NewNotifyCmd(noticeType, message)` |
| `HasActiveAlert()` | `HasActiveNotice()` |
| internal `alertMsg`/`alert` | internal `notifyMsg`/`notice` |

Required options container:
1. `NotifyOpts` must contain all constructor-time options currently provided positionally or by defaults.
2. At minimum, include:
   1. Width controls (`Width`, `MinWidth`).
   2. Lifetime (`Duration`).
   3. Prefix mode (`UseNerdFont`, `UseUnicodePrefix`).
   4. Interaction (`AllowEscToClose`).
   5. Position (`Position`).
3. Constructor validates and normalizes options.
4. `With*()` methods remain for runtime mutation via immutable return values.

Compatibility rule:
1. v0.1.0 and v0.2.0 can break BubbleUp naming.
2. Behavioral parity for notification rendering and placement is required unless explicitly documented.

## PHASE `ALIGN` (Plan + Contract Lock)

1. Finalize exported API contract for `NotifyModel`, `NotifyOpts`, `NoticeDefinition`, and notice keys.
2. Define exact defaults so constructor behavior is deterministic.
3. Define what is parity-critical:
   1. Top/bottom + left/center/right placement.
   2. Dynamic width mode behavior (`MinWidth`).
   3. Prefix modes (ASCII, Unicode, NerdFont).
   4. ESC-to-dismiss semantics.
   5. Single-active-notice lifecycle with timeout ticks.
4. Define implementation policy for color interpolation:
   1. Preserve current behavior first.
   2. Defer optimization/replacement decisions until parity is passing.

Acceptance gate:
1. Public API contract and defaults are documented in this plan and reflected in package doc TODO checklist.

## PHASE `BASELINE_V1` (Charm v1 Baseline Module + `v0.1.0`)

1. Create module skeleton:
   1. `teanotify/go.mod`
   2. package files (`notice.go`, `model.go`, `position.go`, `util.go`)
   3. `README.md`
   4. `LICENSE` attribution notes (and NOTICE text if needed)
2. Port BubbleUp behavior with renamed public symbols.
3. Replace non-compliant patterns during port:
   1. Remove fatal logging pathways from registration/validation.
   2. Return structured errors where operations can fail.
   3. Remove ignored-error declarations.
4. Implement `NewNotifyModel(opts NotifyOpts)` and preserve relevant `With*()` methods.
5. Add a minimal example app under `teanotify/examples/` proving integration flow:
   1. spawn notice command in `Update()`
   2. pass messages through notify model
   3. overlay in `View()`
6. Validate module baseline:
   1. `go test ./...` in `teanotify`
   2. `go vet ./...` in `teanotify`
   3. `go build -o /dev/null ./...` in `teanotify`

Acceptance gate:
1. `teanotify` builds and runs on Charm v1 dependencies.
2. Baseline tests and example smoke checks are green.
3. Module is ready for `teanotify/v0.1.0` tag as part of repo baseline tagging.

## PHASE `TEST_V1` (Behavioral Coverage on Baseline)

Add `_test.go` coverage for parity-critical behavior before v2 migration:
1. Constructor/defaults:
   1. default position and duration behavior
   2. invalid option normalization
2. Notice type registration:
   1. custom notice registration
   2. invalid color handling
3. Update lifecycle:
   1. `NewNotifyCmd` activation
   2. timeout expiration via tick messages
   3. ESC close behavior with and without opt-in flag
4. Render behavior:
   1. overlay at each supported position
   2. dynamic width clamping between min and max
   3. wrapping/hanging-indent behavior
   4. overlay behavior when notice width exceeds content width
5. ANSI-safe cutting helpers:
   1. `cutLeft`
   2. `cutRight`
   3. color reset behavior

Acceptance gate:
1. V1 test suite captures expected behavior and passes reliably.
2. Example smoke test(s) pass.

## PHASE `MIGRATE_V2` (Charm v2 API Migration)

1. Move imports:
   1. `github.com/charmbracelet/bubbletea` -> `charm.land/bubbletea/v2`
   2. `github.com/charmbracelet/lipgloss` -> `charm.land/lipgloss/v2`
2. Update Bubble Tea message handling:
   1. `tea.KeyMsg` -> `tea.KeyPressMsg`
   2. replace `msg.String()=="esc"` checks with v2-equivalent matching strategy
3. Update model interface behavior for v2 (`View() tea.View` where required by interfaces).
4. Confirm timer/tick handling remains semantically equivalent.
5. Re-run full suite and fix regressions.

Acceptance gate:
1. `teanotify` compiles cleanly on Charm v2 APIs.
2. Existing v1 parity tests are migrated and pass on v2.
3. No undocumented behavior regressions.

## PHASE `HARDEN` (API Quality + House-Rule Compliance)

1. Enforce immutable value semantics consistently:
   1. `With*()` methods return updated model value.
   2. No hidden in-place mutation leaks through shared maps/slices.
2. Ensure error flows follow doterr composition.
3. Ensure code style follows ClearPath where applicable.
4. Remove dead compatibility aliases once migration docs are explicit.

Acceptance gate:
1. API is coherent and internally consistent.
2. No remaining known house-rule violations.

## PHASE `DOCS` (Migration and Usage Documentation)

1. Write/refresh `teanotify/README.md` with:
   1. quick start using `NewNotifyModel(opts NotifyOpts)`
   2. notice keys and custom notice registration
   3. `Update()` and overlay integration pattern
   4. position and width behavior
2. Add BubbleUp -> TeaNotify migration section:
   1. import path changes
   2. symbol rename table
   3. constructor migration examples
3. Document known differences from upstream BubbleUp behavior.

Acceptance gate:
1. A consumer can migrate from BubbleUp without tribal knowledge.

## PHASE `RELEASE_GATES` (Tag + Cross-Plan Readiness)

1. Baseline release track:
   1. Ensure `BASELINE_V1` + `TEST_V1` gates are green.
   2. Cut `teanotify/v0.1.0` with repo baseline tagging workflow.
2. Charm v2 release track:
   1. Ensure `MIGRATE_V2`, `HARDEN`, and `DOCS` gates are green.
   2. Ensure linked-plan gates are satisfied per `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md`.
   3. Mark `teanotify` ready for synchronized `v0.2.0` release.

Acceptance gate:
1. `teanotify` has validated v1 baseline and validated v2 implementation status.
2. `teanotify` no longer blocks synchronized go-tealeaves `v0.2.0`.

## Test and Validation Commands

Baseline and v2 tracks both use:
1. `go test ./...` (from `teanotify/`)
2. `go vet ./...` (from `teanotify/`)
3. `go build -o /dev/null ./...` (from `teanotify/`)
4. Example build and run smoke checks under `teanotify/examples/`

## Risks and Mitigations

1. Risk: subtle rendering regressions from ANSI width math.
   Mitigation: targeted util tests and golden-style render assertions.
2. Risk: API churn from constructor redesign.
   Mitigation: lock `NotifyOpts` early and document defaults clearly.
3. Risk: v2 key handling behavior drift (`esc` dismissal).
   Mitigation: explicit key handling tests across v1 and v2.
4. Risk: style-rule cleanup changes behavior while refactoring.
   Mitigation: port for parity first, then harden with tests green.

## Completion Checklist

1. `ALIGN` gate satisfied.
2. `BASELINE_V1` gate satisfied.
3. `TEST_V1` gate satisfied.
4. `MIGRATE_V2` gate satisfied.
5. `HARDEN` gate satisfied.
6. `DOCS` gate satisfied.
7. `RELEASE_GATES` gate satisfied.
