---
pod: docs-update
pod_vars:
  codename: "tealeaves"
  working_dir: "."
  language: "Go"
  doc_framework: "Astro Starlight (MDX)"
---
# Update go-tealeaves Documentation Site

## Goal

The Astro Starlight documentation site at `site/` should accurately reflect
every component package in this repo. It is currently behind the code.

## Scope

- **Home page** — must list all current components with accurate descriptions
- **Home page classification** — `site/config.yaml` is the source of truth
  for which components appear in which category and which packages are
  foundations. Components are listed as `package.Type` entries under their
  category; foundations are listed by package name.
- **Component pages** — every package needs a documentation page; existing pages must match current API
- **New component pages** — packages without docs need pages created
- **Navigation/sidebar** — must include all components; foundations get their own top-level "Foundations" section (not nested under Components)
- **Icons** — each component should have a unique SVG icon
- **Module pages** — each Go module needs its own page under `/modules/<name>/` with full exported API; a module index page links to all of them
- **Reference pages** — getting started, etc. must be current
- **Site build** — the site must build without errors when done

## Definitions

- **Component package** — A `tea*` directory containing a Go type that
  implements the Bubble Tea v2 component pattern (i.e. has `Init() tea.Cmd`,
  `Update(tea.Msg) (Self, tea.Cmd)`, and `View() tea.View` methods).
  Detection command (uses Go type checker, not grep):
  ```bash
  bin/tlcli models
  ```
  Each type listed is a **component**. Packages not in the output are
  **foundations** (utilities, color constants, pane primitives, etc.).

  By convention, component types MUST be named `*Model`. To verify:
  ```bash
  bin/tlcli models -check
  ```
  This exits non-zero if any component type violates the naming convention.

  `bin/tlcli` is built from `cmd/tlcli/` in this repo (`go build -o bin/tlcli ./cmd/tlcli`).

- **Foundation package** — A `tea*` directory that does NOT contain any
  component type (i.e. not listed by `bin/tlcli models`). Foundations are
  documented in the Foundations section of the home page and in a
  top-level "Foundations" sidebar section, NOT nested under Components.

## Prerequisites

- `bin/tlcli` must be built: `go build -o bin/tlcli ./cmd/tlcli` from the repo root. It detects tea.Model component types using the Go type checker.
- `bin/tlcli exports -api-only` extracts exported API from all library packages (excluding cmd/, site/, and examples/). It uses the Go AST parser directly, handles generics, embedded types, and monorepo modules correctly. Use `-include-path=<pkg>` for a single package. This replaces `go doc -all` for API verification.
- `lychee` must be installed for link checking. Install via `brew install lychee` or `cargo install lychee`. Verify it works against the built site before starting work.
- **Home page config** (`site/config.yaml`) is the single source of truth for home page classification — which components go in which categories and which packages are foundations. Components are identified by `package.Type` (e.g., `teagrid.GridModel`); foundations by package name. Agents must read this file before modifying the home page.

## Required Reading

The manager MUST instruct ALL agents to read these files before starting any work:

- `docs/BEST_PRACTICES_CHARM_V2.md` — This project uses Bubble Tea v2, released recently after v1 was dominant for years. Claude's training data defaults to v1 patterns which will produce code that compiles but behaves incorrectly. Key differences include key name constants (e.g., `"space"` not `" "`), method signatures, and model interface changes.

## Pre-Sprint Checklist

Before assigning work, run a conformance check against all criteria below to identify pages that already fail. Include remediation of pre-existing failures explicitly in the sprint scope, or explicitly defer them with a note.

- **Verify file existence** — Before issuing any brief that references a missing file, run `ls` or a quick existence check to confirm. Add a verification command to the brief so the worker can confirm the stated precondition.
- **Verify symbol lists** — Before issuing any brief that names specific Go symbols, run `bin/tlcli exports -include-path=<package>` for each mentioned package and confirm every symbol exists and lives in the stated package.
- **Stale entry cleanup** — Verify every `iconMap` key in `PageTitle.astro` has a corresponding `.mdx` page AND a corresponding live package (directory with non-test `.go` files). Remove stale entries for deleted or superseded packages before issuing work.
- **Package→slug mapping** — The canonical mapping from package name to doc slug should be defined in `site/config.yaml`. If it is not yet there, derive it from `PageTitle.astro` iconMap keys and add it.
- **Infrastructure exclusion list** — `site/config.yaml` should contain a list of infrastructure symbols excluded from API coverage checks (AppendErr, CombineErrs, ErrMeta, ErrValue, NewErr, WithErr, IsDotErrEntry, etc.). `bin/tlcli exports -api-only` and auditors should reference this single source.

## Briefing Guidelines

- **Machine-verifiable criteria** — For every criterion, include at least one machine-verifiable check command. Prefer exit-code checks over "look at the output and decide."
- **Symbol lists are minimums** — Add to every API-coverage brief: "This list is a minimum. Run `bin/tlcli exports -include-path=<package>` and document every non-infrastructure exported symbol — not just those listed here."
- **No 'just add X' shorthand** — Either say "heading only — no API audit needed for this page" or "heading plus full API audit." If uncertain, default to full audit.
- **Break large briefs into focused passes** — Split doc-update briefs into two passes max: (1) structural pass — add missing sections; (2) content pass — fill API coverage gaps. Each pass gets its own review checkpoint.
- **Define 'promote' concretely** — When asking to "promote," "expand," or "document" a type, include a concrete model: subsection heading, one-line description, fields table, and constructor row.

## Verification Steps

- **Exported API checklist** — Run `bin/tlcli exports -api-only` and diff the output against the documentation pages. Every exported type, field, function, and method must appear in the corresponding page. Use `-include-path=<package>` for a single package.
- **Classification check** — Run `bin/tlcli models` and cross-reference against
  `site/config.yaml`. Every `package.Type` in config.yaml must appear in
  `bin/tlcli models` output. Any mismatch is a misclassification. Foundation
  packages in config.yaml must NOT appear in `bin/tlcli models` output.
  Also run `bin/tlcli models -check` to verify naming conventions.
- **Install section format** — Every component page must have a dedicated `## Install` heading with a fenced bash block containing the `go get` command. Do not rely on the header quote block alone.
- **Required sections lint** — Every component page must have: `## Install` (with fenced bash block), at least one usage example section, and `## API Reference`. Consider adding a Remark/MDX lint rule to warn on missing sections.
- **Icon dual-wiring check** — After adding or updating an icon, verify it is wired in both `PageTitle.astro` iconMap AND the home page card entry in `index.mdx`. Consider consolidating icon registration into a single source-of-truth file (e.g. `icons.ts`) that both components import.
- **Code example verification** — All code examples must live in `./site/examples/` (repo root, sibling of `./site/src/`). Run `bin/tlcli audit` to identify unverified and stale code blocks. See the coder agent briefing for the full workflow. The docs-writer MUST NOT invent code examples directly in documentation — all code in docs must trace back to compiled source in `./site/examples/`.
- **Link checking** — After building the site, run `lychee ./site/dist/` against the built output to verify no broken links. Configure with `lychee.toml` if needed to exclude external URLs or known false positives.
- **API completeness artifacts** — The coder must produce a `bin/tlcli exports -api-only` artifact at submission (not just a spot-check) and diff against the corresponding doc pages.

## Constraints

- The documentation framework is Astro Starlight (MDX format)
- Match the tone and structure of existing documentation pages
- Do not modify Go source code — only documentation and site files
- Verify the site builds (`just site-build`) before reporting done
- This is PRE-RELEASE software. Nothing is "deprecated" or "legacy." If a package has been superseded, do NOT document it, do NOT add it to the sidebar, and do NOT create a doc page for it. Instead, flag it to the manager as source code that should be deleted from the repo. There are no legacy components in pre-release software — only current code and code that should be removed.
- If a module has no released Go code and no documentation page, do not create a page, icon, sidebar entry, or card for it.

## Agent Briefings

Project-specific instructions delivered to each agent via `h2pp agent briefing`.

### code-auditor

- Run `bin/tlcli models` and cross-reference against `site/config.yaml`
- Verify every `package.Type` in config.yaml appears in `bin/tlcli models` output
- Verify foundation packages in config.yaml do NOT appear in `bin/tlcli models`
- `site/config.yaml` is the source of truth for classification — flag discrepancies to the manager
- Run `bin/tlcli models -check` to verify naming conventions
- Run `bin/tlcli audit` (or `h2pp audit`) for structural gap analysis
- Verify every foundation package (per `bin/tlcli models` exclusion + `site/config.yaml`) has a corresponding FoundationCard in the Foundations section of `index.mdx`
- **Two-pass audit**: (1) Structural pass — verify heading/section existence via grep (fast); (2) Content accuracy pass — verify API coverage via `bin/tlcli exports -api-only` diff (slow). Run structural pass first.
- **Verify prior fixes** — Before starting a new audit, read the prior audit report. For each P0/P1 item marked fixed, verify the fix is present in the current state. Report any regressions as new findings.
- **Report template** — Output findings as a markdown table: `PACKAGE | TYPE | FINDING | SEVERITY (P0/P1/P2) | FILE:LINE`. This allows downstream agents to ingest audit output directly as a task list.
- **Priority tiers** — Mark each checklist item: P0 (must complete before signoff), P1 (complete if time allows), P2 (record but do not block). API coverage diff (`bin/tlcli exports -api-only`) is always P0; heading/section checks are P1; style and prose quality are P2.

### evaluator

- Run `bin/tlcli audit` (or `h2pp audit`) for structural discovery instead
  of manually exploring directories
- For classification verification: independently run `bin/tlcli models` and
  cross-reference against `site/config.yaml` AND the rendered home page
- Any mismatch between `bin/tlcli models`, config.yaml, and the home page is a P0 failure
- Also run `bin/tlcli models -check` to verify naming conventions

### coder

- `./site/examples/` is at repo root (sibling of `./site/src/`), NOT
  inside `./site/src/content/docs/examples/`
- ONE `go.mod` at `./site/examples/go.mod` — not one per subdirectory
- Each subdirectory is a Go package under that single module
- Complete examples: `./site/examples/<name>/main.go`
- Snippet verification: `./site/examples/<name>/verify_test.go`
- Every `.go` file needs a source comment tracing to the doc:
  `// Source: site/src/content/docs/components/grid-view.mdx:166`
  This applies to BOTH `main.go` (complete examples) AND `compile_test.go`
  (snippet verification) files.
- **Audit scope is ALL `.mdx` files** under `site/src/content/docs/`, including
  components/, guides/, cookbook/, reference/, contributing/, and migration/.
  Run `bin/tlcli audit` which covers all pages automatically.
- Verify with: `go mod tidy`, `go vet ./...`, `go build ./...`
  from `./site/examples/`
- Every fenced Go code block in docs must have a JSX comment above it
  pointing to the verified source: `{/* verified: site/examples/path/file.go:LINE */}`
- Non-compilable or pseudo-code blocks must be wrapped in a test function
  that adapts the snippet (stub missing types) so it compiles. There are
  no skip-eligible code blocks.
- Do NOT create `doc.go` in `site/examples/<name>/` directories that only
  contain `compile_test.go` files — the test package declaration suffices.
- **Workflow order**: (1) Write and finalize the .go source file;
  (2) run `go build ./...` and `go vet ./...`; (3) THEN add the
  `{/* verified: */}` marker to the MDX; (4) run `bin/tlcli audit` to
  confirm 0 stale. Do not add MDX markers until the .go file is finalized.
  Note: adding a Source: comment at line 1 shifts code down — update MDX
  markers from `:N` to `:N+1`. Run `bin/tlcli audit examples --fix` to update stale line numbers and seed hashes automatically.

### icon-designer

- If an existing icon has a non-standard name, do NOT rename or delete it — other live references may depend on it. Create a new correctly-named file. Only delete the old file after confirming no doc pages or config files reference it.
- Icon designer is responsible for both creating SVG icons AND wiring them into `PageTitle.astro` iconMap and `index.mdx` cards. Complete icon creation first, then proceed to wiring without waiting for manager approval.

### docs-writer

- Read `site/config.yaml` first — it defines which components go in which
  categories and which packages are foundations
- Components (per config.yaml) → Cards in their category section on the home page
- Foundation packages (per config.yaml) → Foundations section only
- Do NOT place foundation packages in the Components section
- Foundations get a top-level sidebar section, not nested under Components
- If config.yaml is unclear or missing an entry, ask the manager
- **MDX comment syntax** is `{/* comment */}`. Do NOT use `<!-- -->` — it is a syntax error in MDX.
- **Bubble Tea v2 checklist** — Scan all code examples for v1 patterns: (1) `tea.KeyMsg` → `tea.KeyPressMsg`; (2) `View() string` → `View() tea.View`; (3) old import paths `github.com/charmbracelet/bubbletea` → `charm.land/bubbletea/v2`. See `docs/BEST_PRACTICES_CHARM_V2.md` for the full list.
- **Page deletion checklist** — Deleting a doc page is a 3-step operation: (1) delete the `.mdx` file; (2) remove its entry from `site/astro.config.mjs` sidebar config; (3) remove its `{/* verified: ... */}` markers and any config entries. A page deletion is not complete until all three steps are done and `just site-build` passes.
- **Type rename checklist** — When renaming a type, interface, or function: (1) update the primary page; (2) `grep -r 'OldName' site/src/content/docs/` and fix all occurrences across all pages; (3) verify `just site-build` passes. Never assume a rename is a single-page change.
- **Orphaned pages** — Do not use the term "orphaned page" without explicit instructions. For each page without a package mapping, the briefing must specify one of: "Delete X.mdx and remove its sidebar entry," "Re-map X.mdx to package Y in config," or "This is standalone cross-cutting content — add as a standalone entry."

## Criteria

**Guideline for writing criteria:** Separate objective (mechanically testable) criteria from subjective ones. For any structural criterion, list the specific required sections or fields. Any task that produces observable, verifiable output must be stated as a positive criterion — never frame verifiable work as "implicit."

- Every tea* package that exports user-facing types or functions has a dedicated documentation page — components under `site/src/content/docs/components/`, foundations in their own section. `site/config.yaml` defines the home page classification. Detection: `bin/tlcli models` for components; any tea* directory with exported symbols not in that list is a foundation.
- Each component page has an entry in the `iconMap` in `site/src/components/PageTitle.astro` pointing to a unique SVG file in `site/public/icons/`. Every component gets its own distinct icon — co-package components (e.g. multiple teamodal models) each need a unique icon, not a shared one. No two SVG files on disk are byte-for-byte identical (verify: `md5sum site/public/icons/*.svg | sort | uniq -d` must return empty). Additionally, no two pages may reference the same icon path (verify: `grep -h 'icon=' site/src/content/docs/index.mdx | sort | uniq -d` must return empty). All referenced icon paths must resolve to existing files.
- Each component page has an entry in `iconMap` in `PageTitle.astro`, appears in the navigation sidebar config (`astro.config.mjs`), and appears on the home page (`index.mdx`) as a `ComponentCard` (for components) or `SystemCard` (for foundations/utilities) in the section dictated by `site/config.yaml`. Each entry appears exactly once on the home page — no duplicates across sections. The sidebar must include every component (this subsumes the former C7 requirement).
- Home page one-line descriptions must match the package-level Go doc comment (or first sentence thereof). For multi-page packages (e.g., teamodal → multiple dialog pages), each card's description must describe its specific model type, not the package as a whole.
- The integer N in the frontmatter `description` field of `site/src/content/docs/components/index.md` (e.g. "Overview of the N Tea Leaves components") must equal the total count of data rows (non-header rows) across all component tables in that file.
- Site builds without errors (`just site-build`).
- Every component page (not just new ones) has: frontmatter (title, description), a Module/Type info line formatted as a borderless table with short module names (e.g. `teamodal` not the full `github.com/mikeschinkel/go-tealeaves/teamodal` — the full path belongs on the module page only), an `## Install` section (with fenced bash `go get` block), at least one code example, an `## API Reference` section (exact heading required for mechanical testability) covering exported types/methods, and a `## Related Components` section.
- No broken links or missing references. Verified by running `lychee` against the built `./site/dist/` output.
- Component pages document all exported types, their fields, and all exported functions and methods. Verify using `bin/tlcli exports -api-only` and diffing against the doc pages. The coder must produce a diff artifact for ALL packages at submission.
- Every component entry on the home page references a `package.Type` from `site/config.yaml` components section. No foundation packages appear in the Components area.
- Every foundation package (per `site/config.yaml`) appears in the Foundations section of the home page and has its own top-level sidebar section, not nested under Components.
- Running `bin/tlcli models -check` exits with code 0 (all component types follow the `*Model` naming convention). This is a standing criterion for all sprints involving component packages.
- **NON-NEGOTIABLE**: Every fenced Go code block in `./site/src/content/docs/**/*.mdx` must be verified against the actual package API. The scope includes ALL MDX files (components/, guides/, cookbook/, reference/, contributing/, etc.) — not just component pages. Before starting, the coder must produce a count of all matching Go blocks to agree on scope. Only fenced ` ```go ` blocks require `{/* verified: ... */}` or `{/* non-compilable: ... */}` markers — other language blocks (bash, text, diagrams) do not. All examples live under `./site/examples/` with a SINGLE `go.mod` (not one per subdirectory); this module structure requirement is part of the criterion, not just the agent briefing. Snippet wrappers follow the naming convention `./site/examples/<page-slug>/` to enable mechanical verification of coverage. Complete examples are `main.go` files; snippets are verified via test/wrapper files. Every `.go` file containing extracted documentation code must have a Go comment tracing back to its source doc file and line. Exemption: `doc.go` package-declaration files and other infrastructure-only files with no extracted documentation code are exempt from the Source: comment requirement. After adding Source: comments to line 1 of any `.go` file, run `bin/tlcli audit` immediately to catch stale verified markers caused by line shifts — update MDX markers before committing. Every example must pass `go build` and `go vet`. If `./site/examples/` does not exist, the coder must create it. The evaluator MUST NOT negotiate this down to syntax-only checking. `./site/examples/` is at repo root (sibling of `./site/src/`), NOT inside `./site/src/content/docs/examples/`.
