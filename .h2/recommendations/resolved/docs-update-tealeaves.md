# Resolved Recommendations: docs-update-tealeaves

# Resolved Recommendations: docs-update

## 2026-04-03

## Criteria Results

| ID | Criterion                          |
|----|------------------------------------|
| C8 | All exported APIs documented       |


### Rec #1: Add a package-to-exported-API checklist step — ACCEPTED
Added a "Verification Steps" section to the recipe with a grep-based exported symbol checklist. Ensures every exported symbol in each package is covered by the component page before marking done.

### Rec #3: Add a type-name cross-check step — ACCEPTED
Added to the recipe's Verification Steps. Each ComponentCard `type=` prop in index.mdx must be verified against actual Go type names.

### Rec #4: Standardize the install snippet format — ACCEPTED
Added to the recipe's Verification Steps. Every component page must have a dedicated `## Install` heading with a fenced bash block.

### Rec #5: Add a `doterr.go` exclusion note to C8 — ACCEPTED
Updated the C8 criterion in the recipe to explicitly exclude internal infrastructure files (doterr.go, errors.go) from documentation coverage requirements.

## 2026-04-03

### Criteria Referenced

| ID | Criterion                          |
|----|------------------------------------|
| C2 | Unique icons per component         |

### Rec #2: Separate icon wiring from icon creation in C2 (Unique icons per component) — ACCEPTED WITH MODIFICATION
Original recommendation proposed splitting C2 into two sub-criteria (icon creation vs. icon wiring). User agreed these are different concerns but decided icon wiring (updating iconMap, sidebar config, card components) should be an implicit docs-writer responsibility rather than a separate criterion. Updated C2 in the recipe to note that icon wiring is a docs-writer task.

## 2026-04-03

### Criteria Referenced

| ID | Criterion |
|----|-----------|
| C1 | Every component package has a dedicated documentation page |
| C2 | Every component doc page has a distinct SVG icon; icon references in index page cards and the iconMap are correct |
| C3 | Icon wiring — each component page has an entry in iconMap in PageTitle.astro, appears in the navigation sidebar config, and has a ComponentCard on the home page |
| C6 | New pages have at minimum: an Install section with a fenced bash go get command, at least one usage example, and an API Reference section. Existing pages must also conform |
| C7 | New pages match the tone, structure, and level of detail of existing pages |
| C8 | Component pages mention all exported types, their exported fields, and all exported functions and methods. Exported constants and variables are also in scope |
| C9 | The component count in the index page subtitle matches the actual number of component rows in the table |
| C10 | All Go code examples must compile in Go |

### Accepted Recommendations

**Rec 1 — ACCEPTED: Require `go doc -all <package>` diff for C8 (All exported types/fields/functions/methods documented) verification.**
Instead of writers manually cross-referencing source files, mechanically diff `go doc` output against page content to catch omissions automatically. Added to recipe Verification Steps.

**Rec 2 — ACCEPTED: Add page template or linting rule for C6 (Install section, usage example, API Reference required).**
A Remark/MDX lint rule warns when a component page lacks `## Install` (with fenced bash block), a usage example, and `## API Reference`. Added to recipe Verification Steps.

**Rec 3 — ACCEPTED: Consolidate icon registration for C2 (Unique icons per component).**
Icon wiring is split across `PageTitle.astro` iconMap and `index.mdx` ComponentCard. A single `icons.ts` export used by both would prevent dual-wiring omissions. Added as guidance in recipe Verification Steps; C3 promoted to explicit positive criterion.

**Rec 4 — ACCEPTED: Pre-sprint conformance check.**
Before assigning work, run a conformance check to identify all pages that currently fail criteria. Include remediation of pre-existing failures explicitly in sprint scope, or explicitly defer them. Added as Pre-Sprint Checklist section in recipe.

**Rec 5 — ACCEPTED (modified): Define "component package" explicitly in the recipe for C1 (Every component package has a dedicated documentation page).**
User clarification: a component is defined by implementing the Bubble Tea v2 `tea.Model` interface (`Init() tea.Cmd`, `Update(tea.Msg) (tea.Model, tea.Cmd)`, `View() tea.View`), not merely by exporting any public type. Added as Definitions section in recipe.

**Rec 6 — ACCEPTED: Separate objective criteria from subjective criteria in recipe templates.**
C7 (page tone/structure) originally said "match tone and structure" which was untestable. Recipe template now prompts: "For any structural criterion, list the specific required sections or fields." Added as guidance above the Criteria section.

**Rec 7 — ACCEPTED: Always make implicit tasks explicit positive criteria.**
C3 (icon wiring) was initially framed as "not a separate criterion," which would have left wiring unverified. Any task producing observable, verifiable output is now stated as a positive criterion. C3 rewritten as explicit criterion in recipe.

**Rec 8 — ACCEPTED (modified): Build a code-example compile harness for C10 (All Go code examples must compile).**
User clarification: rather than scoping down to signature-only checks when no harness exists, the harness should be created as a prerequisite. Executable examples should live in `./site/examples/` as independent Go modules (each in its own directory) so they can be compiled and verified. Added to recipe Verification Steps.

## 2026-04-12

### Criteria Referenced

| ID  | Criterion |
|-----|-----------|
| C1  | Every tea* directory with at least one exported Go symbol has a dedicated documentation page |
| C2  | Each component page references a distinct icon file; no two SVG files are identical copies |
| C3  | Each page: iconMap entry, sidebar config entry, ComponentCard/SystemCard on home page |
| C4  | Non-empty, non-placeholder descriptions matching primary exported type or go doc summary |
| C4a | Each ComponentCard has a non-empty `description` prop (objective) |
| C4b | Descriptions reviewed for accuracy against package API/source (subjective) |
| C5  | Component count in index page subtitle matches actual table row count |
| C7  | Navigation sidebar includes every component |
| C8  | Every component page has all required structural sections |
| C9  | No broken links or missing references |
| C10 | All exported symbols documented per `tlcli exports`; doterr.go/errors.go excluded |
| C11 | Each example in site/examples/ has go.mod with replace directives for local packages |
| C13 | NON-NEGOTIABLE: code blocks verified, single go.mod, go build + go vet pass, Source: comments |
| C14 | `tlcli models -check` exits 0 (naming convention) |

### Sprint Contract Refinements — Round 1

**SCR-1 — ACCEPTED: Broaden C1 (Every tea* package with exported symbols has a doc page) to include non-Model packages.**
Recipe already had separate criteria for components and foundations. Updated C1 wording to explicitly cover all tea* packages with exported symbols, with detection via `tlcli models` for components and directory scanning for foundations.

**SCR-2 — ACCEPTED: Clarify C2 (Unique icons per component) mechanism.**
Updated recipe to specify `iconMap` in `PageTitle.astro` as the mechanism and added `md5sum` verification command.

**SCR-3 — ACCEPTED: Clarify C5 (Component count in subtitle) to specify exact file.**
Updated recipe to specify `site/src/content/docs/components/index.md` frontmatter description field and counting method.

**SCR-4 — ACCEPTED (merge into C3): Remove C7 (Navigation sidebar includes every component) as redundant with C3.**
Merged C7 into C3 and removed C7 as a separate criterion.

**SCR-5 — ACCEPTED (already applied): Scope C8 (Page structure) to all component pages, not just new ones.**
Recipe already says "every component page (not just new ones)." Confirmed.

**SCR-6 — SKIPPED: Replace 'key methods' in C10 (API documentation coverage) with `tlcli exports`.**
Already reflected in recipe. User noted `tlcli exports` replaces `doc-diff` (never implemented). Updated all recipe references from `doc-diff` to `tlcli exports`.

**SCR-7 — ACCEPTED (already superseded): Replace C11 (site/examples/) with syntax checking.**
Current recipe has the stronger NON-NEGOTIABLE criterion with compiled examples. No change needed.

**SCR-8 — ACCEPTED (already superseded): Specify tooling for C9 (No broken links).**
Recipe already specifies lychee. No change needed.

### Evaluation Findings — Round 1

**EF-1 — ACCEPTED: C3 uses both ComponentCard and SystemCard.**
Updated recipe to mention both card types. Home page now has Components and Foundations sections; foundations use SystemCard (user prefers FoundationCard naming).

**EF-2 — ACCEPTED (modified): C8 Module/Type info line format.**
User requested borderless table with short module names (e.g. `teamodal` not full path). Full module path belongs on the module page only. Updated recipe.

**EF-3 — ACCEPTED (already addressed): C10 secondary types missed by workers.**
`tlcli exports` catches missing secondary types. Updated all recipe references from `doc-diff` to `tlcli exports`.

**EF-4 — ACCEPTED (superseded): `astro check` needs pre-installed deps.**
Recipe uses lychee instead. No change needed.

**EF-5 — ACCEPTED (superseded): check-go-syntax.sh should be wired into CI.**
Superseded by `tlcli audit` and compiled examples. No change needed.

### Sprint Contract Refinements — Round 2

**SCR-R2-1 through SCR-R2-5 — ACCEPTED (all already applied).**
All 5 refinements (C1 broadened, C4 testable, C8 scoped to all pages, C10 all exported symbols, C11 go.mod approach) were already in the current recipe.

**SCR-R2-6 — ACCEPTED (already applied): Home page must be index.mdx.**
Recipe already specifies index.mdx. No change needed.

### Coder Post-Mortems — Round 1

**CP-1 — ACCEPTED: Pre-verify claimed broken links before including in briefs.**
Added to pre-sprint checklist: verify file existence before issuing briefs.

**CP-2 — ACCEPTED (superseded): check-go-syntax.sh doesn't handle trailing attributes.**
Superseded by `tlcli audit`. No change needed.

**CP-3 — ACCEPTED: Concrete machine-verifiable success criteria work well.**
Added to briefing guidelines: include machine-verifiable check command for every criterion.

### Writer Post-Mortems — Round 1

**WP-1 — ACCEPTED: Pre-verify symbol lists before issuing briefs.**
Added to pre-sprint checklist.

**WP-2 — ACCEPTED: Symbol lists in briefs should be declared minimums.**
Added to briefing guidelines.

**WP-3 — ACCEPTED: Use relative paths for Go tools.**
Added to briefing guidelines: use `tlcli exports ./<package>`.

**WP-4 — ACCEPTED: Don't say 'just add X' shorthand.**
Added to briefing guidelines.

**WP-5 — ACCEPTED: Break large briefs into focused passes.**
Added to briefing guidelines.

**WP-6 — ACCEPTED: Define 'promote' with a concrete model.**
Added to briefing guidelines.

### Evaluation Findings — Round 3

**EF-R3-1 — ACCEPTED: C13 (NON-NEGOTIABLE) — Source: comment at line 1 breaks verified markers.**
Added workflow note to recipe: run `tlcli audit` after adding Source: comments to catch line shifts.

**EF-R3-2 — REJECTED (opposite applied): C2 (Unique icons) — co-package icon sharing.**
User requires each component to have its own unique icon, not shared per-package. Updated C2 to require per-page distinct icons with verification command.

**EF-R3-3 — ACCEPTED: C4b (Description accuracy) — per-model descriptions for multi-model packages.**
Updated recipe: each card's description must describe its specific model type, not the package as a whole.

**EF-R3-4 — ACCEPTED: C13 (NON-NEGOTIABLE) — doc.go infrastructure files exempt from Source: comment.**
Added explicit exemption for doc.go and infrastructure-only files.

### Auditor Post-Mortems — Round 3

**AUD3-1 — SKIPPED: Auditor should fix doc-pages.yaml gaps directly.**
User questions role separation (auditors audit, not fix). Needs further discussion.

**AUD3-2 — SKIPPED: Redefine 'stale' as 'exported API changed.'**
May be addressed by `tlcli audit examples`. Needs verification.

**AUD3-3 — SKIPPED: Pre-sprint verify all doc-pages.yaml packages have live Go source.**
Mappings should be in site/config.yaml, not .h2. May be addressed by `tlcli audit examples`.

**AUD3-4 — SKIPPED: ComponentCard misclassification decision tree.**
May be handled by site/config.yaml. Needs verification.

**AUD3-5 — SKIPPED: Include tlcli exports diff in audit report for stale packages.**
Needs exploration of which tlcli command fits best.

**AUD3-6 — ACCEPTED: Verify foundation packages have a FoundationCard.**
Added to auditor briefing. User prefers FoundationCard naming over SystemCard.

**AUD3-7 — REJECTED: Standardize on 'tlcli audit' not 'h2pp audit'.**
User wants to keep both names.

### Coder Post-Mortems — Round 3

**COD3-1 through COD3-4 — SKIPPED: JSX syntax, audit scope, Source: for both file types, line number semantics.**
User does not have enough context to accept yet.

**COD3-5 — ACCEPTED: Non-compilable blocks need wrapper tests, not skipping.**
Added to coder briefing.

**COD3-6 — ACCEPTED: Test-only packages don't need doc.go.**
Added to coder briefing.

**COD3-7 — ACCEPTED: Add verified markers AFTER finalizing .go source files.**
Added workflow order to coder briefing.

### Icon-Designer Post-Mortems — Round 3

**ICN3-1 — ACCEPTED: Create new icon files, don't rename/delete existing.**
Added icon-designer briefing section to recipe.

**ICN3-2 — SKIPPED: Combine Phase A+B into single briefing.**
Needs role clarification: is icon wiring (PageTitle.astro, index.mdx) the icon designer's or docs-writer's responsibility?

**ICN3-3 — SKIPPED: Add package→slug mapping table.**
User says this belongs in site/config.yaml. Noted in pre-sprint checklist.

**ICN3-4 — ACCEPTED (already applied above in EF-R3-2): Each component gets its own unique icon.**

**ICN3-5 — ACCEPTED (modified): Module exclusion in Constraints.**
Instead of hardcoding teaterm, added generic rule: "If a module has no released Go code and no documentation page, do not create a page, icon, sidebar entry, or card for it."

### Writer Post-Mortems — Round 3

**WRT3-1 — ACCEPTED: Page deletion is a 3-step operation.**
Added page deletion checklist to docs-writer section of recipe.

**WRT3-2 — ACCEPTED: MDX uses JSX comments, not HTML.**
Added JSX comment syntax reminder to docs-writer section of recipe.

**WRT3-3 — ACCEPTED: Bubble Tea v2 migration checklist.**
Added v2 migration checklist to docs-writer section of recipe.

**WRT3-4 — SKIPPED: Document known gaps in tlcli exports.**
User says if gaps are known, fix tlcli exports instead of documenting workarounds.

**WRT3-5 — ACCEPTED: Interface renames require cross-page grep.**
Added type rename checklist to docs-writer section of recipe.

**WRT3-6 — ACCEPTED: Replace 'orphaned page' with explicit instructions.**
Added orphaned-page resolution guidance to docs-writer section of recipe.

### Sprint Contract Refinements — Round 3

**SCR-R3-1 — ACCEPTED: Package→slug mapping should be in site/config.yaml.**
Added note to pre-sprint checklist.

**SCR-R3-2 — ACCEPTED: Add stale entry cleanup to pre-sprint checklist.**
Added: verify every iconMap key has a corresponding .mdx page AND live package.

**SCR-R3-3 — ACCEPTED (already applied): Split C4 (Description accuracy) into objective C4a and subjective C4b.**

**SCR-R3-4 — ACCEPTED (already applied): Scope C8 (Page structure) to all component pages.**

**SCR-R3-5 — ACCEPTED (already applied): C9 (No broken links) specifies lychee.**

**SCR-R3-6 — ACCEPTED: Add C14 (`tlcli models -check` exits 0) as standing criterion.**
Added to recipe Criteria section.

**SCR-R3-7 — ACCEPTED (already applied): Recipe uses `tlcli exports` instead of `go doc -all`.**

### Auditor Post-Mortems — Round 1

**AP-1 — ACCEPTED: Require exact `## API Reference` heading.**
All component pages already use this heading. Added "(exact heading required)" note to C8 criterion.

**AP-2 — ACCEPTED: Infrastructure exclusion list in site/config.yaml.**
Added note to pre-sprint checklist that site/config.yaml should contain the exclusion list.

**AP-3 — ACCEPTED: Two-pass audit (structural then content accuracy).**
Added to auditor briefing: structural pass first (grep), then content accuracy pass (tlcli exports diff).

**AP-4 — ACCEPTED: Verify prior fixes step.**
Added to auditor briefing: read prior audit report and verify P0/P1 fixes before starting new audit.

**AP-5 — ACCEPTED: Provide report template.**
Added to auditor briefing: PACKAGE | TYPE | FINDING | SEVERITY | FILE:LINE markdown table format.

**AP-6 — ACCEPTED: Weight checklist items by priority.**
Added to auditor briefing: P0 (API coverage), P1 (heading/section checks), P2 (style/prose).

### Auditor Post-Mortems — Round 3 (remaining)

**AUD3-1 — ACCEPTED: Auditor fixes config gaps directly (in site/config.yaml, not doc-pages.yaml).**
Auditor has all info needed to add missing entries. Fixes trivial config gaps directly instead of flagging.

**AUD3-2 — REJECTED: Redefine 'stale' as 'exported API changed.'**
Deferred to tlcli enhancement — needs investigation of what "API diff" means and which command implements it.

**AUD3-3 — REJECTED: Pre-sprint verify all config packages have live Go source.**
Deferred to future tlcli enhancement (`tlcli audit` could check this).

**AUD3-4 — REJECTED: ComponentCard misclassification decision tree.**
site/config.yaml already defines classification. No separate decision tree needed.

**AUD3-5 — REJECTED: Include API diff in audit report for stale packages.**
Deferred to tlcli enhancement — needs exploration of what "API diff" means (diff against what?).

### Coder Post-Mortems — Round 3 (remaining)

**COD3-1 — ACCEPTED: Specify JSX comment syntax for verified markers.**
Coder briefing already updated to use `{/* verified: ... */}` JSX syntax.

**COD3-2 — ACCEPTED: Audit scope explicitly includes ALL .mdx files.**
Added explicit scope statement to coder briefing: all directories under site/src/content/docs/.

**COD3-3 — ACCEPTED: Source: comment applies to both main.go and compile_test.go.**
Added explicit mention of both file types in coder briefing.

**COD3-4 — ACCEPTED: Line-shift update instruction after adding Source: comment.**
Added to coder workflow order. Noted future `tlcli audit --fix` could automate marker updates.

### Icon-Designer Post-Mortems — Round 3 (remaining)

**ICN3-2 — ACCEPTED: Icon designer does both creation and wiring.**
Added to icon-designer briefing: create SVGs then wire into PageTitle.astro and index.mdx.

**ICN3-3 — ACCEPTED: Package→slug mapping belongs in site/config.yaml.**
Noted in pre-sprint checklist as a config.yaml enhancement to implement.

### Writer Post-Mortems — Round 3 (remaining)

**WRT3-4 — REJECTED: Document known tlcli exports gaps.**
User directive: fix `tlcli exports` to cover generic constructors, Default*KeyMap(), named return types — don't document workarounds in the recipe.
