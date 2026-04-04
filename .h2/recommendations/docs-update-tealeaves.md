# Sprint Contract Refinements & Recommendations
## Pod: docs-update-tealeaves
## Date: 2026-04-04

---

## Criteria Referenced

| ID  | Criterion (original) |
|-----|-----------|
| C1  | Every component package (tea* dir with tea.Model impl) has a dedicated doc page |
| C2  | Each component page references a distinct icon file; no two SVGs identical |
| C3  | Each page: iconMap entry, sidebar config entry, ComponentCard on home page |
| C4  | Home page lists all components with accurate one-line descriptions |
| C5  | Component count in index page subtitle matches actual table row count |
| C6  | Site builds without errors (just site-build) |
| C7  | Navigation sidebar includes every component |
| C8  | Every new page has required structural sections |
| C9  | No broken links or missing references |
| C10 | Component pages document all exported types, fields, and key methods |
| C11 | All code examples exist as compilable source in ./site/examples/ |

---

## Sprint Contract Refinements
### Round 1 — 2026-04-04 (Evaluator observation, pre-work criteria review)

**SCR-1: C1 — Component package definition excludes non-Model packages that need pages**
Source: Evaluator observation.
Packages teadiffr (TUIRenderer/DiffRenderer), teaterm (Terminal/ProcessViewer), and teacolor (SemanticColor) have no tea.Model implementation in production code but each has an existing and expected documentation page. The criterion as written would not require pages for these packages.
Recommended wording: "Every tea* package that exports user-facing types or functions must have a dedicated documentation page under site/src/content/docs/components/."
How to apply: The authoritative list of "component packages" should be derived from the set of tea* directories containing exported non-test symbols — not only those implementing tea.Model.

**SCR-2: C2 — 'References' is undefined**
Source: Evaluator observation.
The word "references" in "each component page references a distinct icon file" does not specify the mechanism. In practice, the reference is an entry in the iconMap in PageTitle.astro.
Recommended wording: "Each component page has an entry in the iconMap in site/src/components/PageTitle.astro pointing to a unique SVG file in site/public/icons/; no two SVG files are byte-for-byte identical."

**SCR-3: C5 — 'Index page subtitle' is ambiguous**
Source: Evaluator observation.
Two index files exist: site/src/content/docs/index.mdx (home page) and site/src/content/docs/components/index.md (component listing). The home page has no component count; the component listing has "25 Tea Leaves components" in its frontmatter description. The word "subtitle" does not map to either file's structure unambiguously.
Recommended wording: "The integer N in the frontmatter description field of site/src/content/docs/components/index.md (e.g. 'Overview of the N Tea Leaves components') must equal the total count of data rows (non-header rows) across all component tables in the body of that file."

**SCR-4: C7 — Redundant with C3**
Source: Evaluator observation.
C3 already requires each page to appear in the navigation sidebar config (astro.config.mjs). C7 repeats this requirement identically. Having duplicate criteria creates ambiguity about whether a single finding constitutes one or two failures.
Recommendation: Merge C7 into C3 or remove C7 entirely.

**SCR-5: C8 — 'Every new page' is not mechanically testable**
Source: Evaluator observation.
There is no defined mechanism to identify which pages are "new" in this sprint vs. pre-existing. This criterion as written only applies to a subset of pages determined by git history, leaving pre-existing broken pages unchecked.
Recommended wording: "Every component page under site/src/content/docs/components/ (excluding index.md) must contain: (1) frontmatter with non-empty title and description fields; (2) a Module or Type line in the body; (3) an ## Install section with a fenced bash code block containing a go get command; (4) at least one Go code example; (5) an ## API Reference section; (6) a ## Related Components section."
How to apply: Apply to ALL component pages, not just new ones.

**SCR-6: C10 — 'Key methods' is subjective and untestable**
Source: Evaluator observation.
"Key methods" requires evaluator judgment and will produce inconsistent results across sprints. The intent is full coverage of the public API.
Recommended wording: "Component pages must document all exported types, their exported fields, and all exported methods, as verifiable by comparing page content against the output of 'go doc -all <package>'."

**SCR-7: C11 — site/examples/ does not exist; criterion introduces out-of-scope infrastructure**
Source: Evaluator observation.
Verified: site/src/content/docs/components/../../../examples/ (./site/examples/) does not exist. Creating this directory, porting all documentation code snippets into compilable Go files, and configuring go build/vet/golangci-lint is significant new build infrastructure — not a documentation update task.
Existing examples live at ./examples/ (teadiffview, teadrpdwn, teagrid, teamodal, teanotify).
Recommendation for this sprint: Replace C11 with "All fenced Go code blocks in documentation pages that contain a complete package declaration (i.e. include 'package main' or 'package <name>') must be syntactically valid Go, verified by running gofmt -e on the extracted blocks."
Recommendation for a future sprint: Create a separate bead for the full verified-examples infrastructure (./site/examples/ directory with go build/vet/lint checks wired into CI).

**SCR-8: C9 — No tooling specified for link checking**
Source: Evaluator observation.
"No broken links or missing references" is testable only if the verification tool and scope (internal-only vs. external) are specified. Astro's built-in link checking catches some issues; a dedicated tool like lychee or broken-link-checker catches more.
Recommendation: Add "Verified using 'astro check' plus manual review of all /go-tealeaves/* internal hrefs. External links are excluded from this criterion."

---

## Evaluation Findings
### Round 1 — 2026-04-04 (Post-evaluation, for future recipe improvement)

**EF-1: C3 — Criterion said 'ComponentCard' but site uses two card types**
Source: Evaluator observation during criteria review.
The home page (index.mdx) uses both ComponentCard (for interactive components) and SystemCard (for utilities and layout helpers). The original criterion required a ComponentCard specifically, which would have caused false failures for pane-widgets, layout-engine, key-registry, theming, and positioning — all correctly represented as SystemCards.
Recipe fix: Acceptance criteria for this project should always say "ComponentCard or SystemCard" and note that utilities/layout pages go into the System section.

**EF-2: C8 — 'Module or Type line' format differs from Markdown definition list**
Source: Evaluator observation.
The actual format used in all pages is a blockquote: `> **Module:** teamodal · **Type:** ConfirmModel`. The criterion said "formatted as a definition list entry" — Markdown definition lists use a different syntax (`term\n: definition`). The examples made intent clear, but the terminology was wrong and could confuse workers.
Recipe fix: Specify the actual format: "a blockquote line starting with > **Module:** or > **Type:**".

**EF-3: C10 — API gaps were widespread and systematic, not isolated**
Source: Evaluator findings across 4 packages.
Multiple packages (teamodal, teagrid, tealayout, teatree) had entire type families missing — not just individual methods. Specifically: ChoiceKeyMap (teamodal), BorderConfig + 4 constructors (teagrid), Direction (tealayout), DefaultTreeKeyMap/DefaultDrillDownKeyMap/DefaultBranchStyle/LoadFileMeta (teatree).
Root cause: workers likely documented the primary model types thoroughly but missed secondary types (KeyMap structs, border configuration, utility functions). All were fixed in one revision pass.
Recipe fix: Add to worker instructions — "for each package, run 'go doc -all .' and verify every non-infrastructure type and function appears in the docs before marking C10 done. KeyMap structs, border/style configuration types, and utility functions are commonly missed."

**EF-4: C9 — astro check not runnable non-interactively in this environment**
Source: Evaluator observation during C9 verification.
Running 'npx astro check' prompted for interactive installation of @astrojs/check — it cannot be used non-interactively without pre-installing the dependency. C9 was verified via build-time link validation and manual slug cross-check instead.
Recipe fix: Add 'bun add @astrojs/check typescript' (or equivalent) to the site setup instructions so astro check can run in CI and evaluation without interactive prompts.

**EF-5: C11 — Script created by coder accurately implements the criterion**
Source: Evaluator observation.
The check-go-syntax.sh script correctly extracts fenced Go blocks containing 'package' declarations and runs gofmt -e. It checked 7 blocks and all passed. The script is well-written and should be retained and wired into CI.
Recipe fix: Add site/scripts/check-go-syntax.sh to the justfile as 'just check-go-syntax' and run it in CI alongside the site build.

---

## Worker Post-Mortems
### Round 1 — 2026-04-04 (Collected by manager after final sign-off)

#### Coder post-mortem

**CP-1: Pre-verify claimed broken links before including them in the brief**
Source: coder-tealeaves.
The sprint brief claimed examples/index.md did not exist; it already existed. This sent the coder on a false trail.
Recipe fix: Before issuing any brief that references a missing file, run `ls` or a quick existence check (e.g. `ls site/src/content/docs/examples*`) to confirm. Add a verification command to the brief itself so the worker can confirm the stated precondition without inferring.

**CP-2: C11 script does not handle fenced blocks with trailing attributes**
Source: coder-tealeaves.
The check-go-syntax.sh script matches opening fences as ` ```go ` with no trailing content. Code blocks written as ` ```go title="example.go" ` or ` ```go filename="..." ` will be silently skipped, leaving those blocks unchecked.
Recipe fix: Update the regex in check-go-syntax.sh to match ` ```go ` followed by any optional trailing attributes (e.g. `` ^(`{3}go[[:space:]]*)`` → `` ^`{3}go ``). Document the intended match pattern in the script header. Add at least one test block with trailing attributes to the script's own test suite.

**CP-3: Concrete machine-verifiable success criteria worked well — use this pattern**
Source: coder-tealeaves (positive confirmation).
Criteria like "just site-build exits 0" required no judgment and could be verified in one command. Workers found these unambiguous and efficient.
Recipe fix: For every criterion in future sprints, include at least one machine-verifiable check command. Prefer exit-code checks over "look at the output and decide."

---

#### Writer post-mortem

**WP-1: Pre-verify symbol lists before issuing the brief**
Source: writer-tealeaves.
Two symbols in the brief (PositionAbove, PositionBelow) do not exist in the codebase. One symbol (EnsureTermGetSize) was attributed to the wrong package (listed under teautils; actually lives in teafields/teamodal). The writer spent time searching for nonexistent symbols.
Recipe fix: Before issuing any brief that names specific Go symbols, run `go doc -all ./pkg/` for each mentioned package and confirm every symbol exists and lives in the stated package. Remove or correct any that do not match.

**WP-2: Symbol lists in briefs should be declared minimums, not exhaustive lists**
Source: writer-tealeaves.
The brief's symbol lists implied completeness, but the evaluator found additional undocumented symbols (ChoiceKeyMap, BorderConfig, Direction, TreeKeyMap, DefaultDrillDownKeyMap) that were not in the brief. Workers stopped at the brief list rather than auditing the full package.
Recipe fix: Add an explicit note to every C10-related brief: "This list is a minimum. Run `go doc -all ./pkg/` and document every non-infrastructure exported symbol — not just those listed here."

**WP-3: Specify the correct go doc invocation form**
Source: writer-tealeaves.
The brief used module path form (`go doc -all github.com/...`) which requires full GOPATH/module setup. The correct form in this workspace is relative path (`go doc -all ./teamodal/` from the repo root, or `go doc -all .` from inside the package directory).
Recipe fix: Standardize all briefs and criteria to use `go doc -all ./pkg/` (relative path from repo root) or `go doc -all .` (from inside the package). Add a note: "Run from the repo root or cd into the package directory first."

**WP-4: Clarify scope when saying 'coverage already good, just add the heading'**
Source: writer-tealeaves.
Some pages were described in the brief as needing only a heading added, but in practice also had missing symbols. Workers took "just add the heading" literally and did not audit for gaps.
Recipe fix: Never use "just add X" shorthand. Either say "heading only — no API audit needed for this page" (and mean it) or "heading plus full API audit." If uncertain, default to full audit.

**WP-5: Break large briefs into focused passes**
Source: writer-tealeaves.
A 17-page brief is difficult to review and course-correct mid-sprint. Errors in early pages set wrong patterns that propagated through later pages.
Recipe fix: Split doc-update briefs into two passes maximum: (1) structural pass — add missing sections (Install, API Reference, Related) to all pages; (2) content pass — fill API coverage gaps per package. Each pass should be a separate bead with its own review checkpoint.

**WP-6: Define 'promote to a proper API entry' with a concrete model**
Source: writer-tealeaves.
The phrase "promote DrillDownKeyMap to a proper API entry" was ambiguous. Workers did not know whether to add a subsection heading, a table, a prose paragraph, or a function signature row.
Recipe fix: When a brief asks to "promote," "expand," or "document" a type, include a concrete model: "Add a subsection heading (### TypeName), a one-line description, a fields table (Field | Type | Description), and a row in the nearest API table for any DefaultX() constructor."

---

#### Auditor post-mortem

**AP-1: Clarify the exact '## API Reference' heading string requirement**
Source: auditor-tealeaves.
C8 requires an "## API Reference" section, but several pages (diff-viewer, layout-engine, dropdown-control) had substantial API coverage under non-standard headings. The auditor flagged these as missing when they had real content.
Recipe fix: Specify in C8 and in briefs whether the criterion requires the exact string "## API Reference" or whether any structured API table qualifies. If the exact heading is required, add a one-time migration task to standardize headings. If equivalent content suffices, define what "equivalent" means (e.g., "a table with columns Name | Signature | Description").

**AP-2: Provide an explicit infrastructure exclusion list for go doc diffs**
Source: auditor-tealeaves.
Running `go doc -all` on each package surfaced 15+ doterr infrastructure functions per package (NewErr, WithErr, AppendErr, ErrMeta, ErrKV variants, IsDotErrEntry, etc.) that are intentionally excluded per C10. Without an explicit list, the auditor had to re-derive exclusions for every package.
Recipe fix: Include a standing exclusion list in the brief and in the recipe: "Exclude from API coverage checks any symbol whose name appears in doterr.go or errors.go, plus: AppendErr, CombineErrs, ErrMeta, ErrValue, Errors, FindErr, MsgErr, NewErr, WithErr, IsDotErrEntry, and all KV-suffixed types and functions." Commit this list to a shared file (e.g., .h2/api-exclusions.txt) so all agents reference the same set.

**AP-3: Separate structural completeness checks from content accuracy checks**
Source: auditor-tealeaves.
Structural checks (does the ## Install section exist?) are fast grep scans that can cover all pages in minutes. Content accuracy checks (does the API table match go doc output?) require line-by-line comparison and take much longer. Mixing them in one pass forces the auditor to switch context constantly.
Recipe fix: Split C8 (structure) and C10 (accuracy) into separate audit passes or assign them to separate agents. Brief each agent with a focused checklist for their pass only.

**AP-4: Include a 'verify prior fixes' step in the audit brief**
Source: auditor-tealeaves.
The project has an existing site/AUDIT.md with prior P0/P1 findings. The auditor had no instruction to cross-reference it, so prior fixes may have been re-audited redundantly or missed if they regressed.
Recipe fix: Add a standard step to all audit briefs: "Read site/AUDIT.md (or the prior audit file). For each P0/P1 item marked fixed, verify the fix is present in the current state. Report any regressions as new findings."

**AP-5: Provide a report template upfront**
Source: auditor-tealeaves.
Without a specified output format, auditor reports varied in structure, making it harder for downstream agents (doc-writers) to parse and act on findings.
Recipe fix: Include a report template in the audit brief. At minimum: "For each package, output: PACKAGE | TYPE | FINDING | SEVERITY (P0/P1/P2) | FILE:LINE." Provide this as a markdown table template. This allows doc-writer agents to ingest audit output directly as a task list.

**AP-6: Weight checklist items by effort and mark P0 vs. nice-to-have**
Source: auditor-tealeaves.
The go doc diff accounted for roughly 80% of audit time. Without priority markers, auditors spent equal time on high-value and low-value checks.
Recipe fix: Mark each checklist item in the brief with a priority tier: P0 (must complete before signoff), P1 (complete if time allows), P2 (record but do not block). The go doc diff is always P0; heading/section checks are P1; style and prose quality are P2.
