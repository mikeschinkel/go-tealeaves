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
