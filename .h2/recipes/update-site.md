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
- **Component pages** — every package needs a documentation page; existing pages must match current API
- **New component pages** — packages without docs need pages created
- **Navigation/sidebar** — must include all components
- **Icons** — each component should have a unique SVG icon
- **Reference pages** — module listing, getting started, etc. must be current
- **Site build** — the site must build without errors when done

## Definitions

- **Component package** — A `tea*` directory is in scope as a component if it contains a Go type that implements the Bubble Tea v2 `tea.Model` interface, i.e. has `Init() tea.Cmd`, `Update(tea.Msg) (tea.Model, tea.Cmd)`, and `View() tea.View` methods.

## Prerequisites

- `doc-go-repo` must be installed and in PATH. It generates structured API documentation from Go source. Install from ~/Projects/go-cli/doc-go-repo.
- Use `doc-go-repo -exclude-file doterr.go ./<package>` to get accurate API inventories. This replaces `go doc -all` for API diffing — it handles generics, embedded types, and monorepo modules correctly.

## Required Reading

The manager MUST instruct ALL agents to read these files before starting any work:

- `docs/BEST_PRACTICES_CHARM_V2.md` — This project uses Bubble Tea v2, released recently after v1 was dominant for years. Claude's training data defaults to v1 patterns which will produce code that compiles but behaves incorrectly. Key differences include key name constants (e.g., `"space"` not `" "`), method signatures, and model interface changes.

## Pre-Sprint Checklist

Before assigning work, run a conformance check against all criteria below to identify pages that already fail. Include remediation of pre-existing failures explicitly in the sprint scope, or explicitly defer them with a note.

## Verification Steps

- **Exported API checklist** — For each package, run `doc-go-repo -exclude-file doterr.go ./<package>` and diff the output against the documentation page. Every exported type, field, function, and method must appear in the page.
- **Type-name cross-check** — For each ComponentCard `type=` prop in index.mdx, verify against `grep '^type' <pkg>/*.go`. Automate with a script rather than relying on memory.
- **Install section format** — Every component page must have a dedicated `## Install` heading with a fenced bash block containing the `go get` command. Do not rely on the header quote block alone.
- **Required sections lint** — Every component page must have: `## Install` (with fenced bash block), at least one usage example section, and `## API Reference`. Consider adding a Remark/MDX lint rule to warn on missing sections.
- **Icon dual-wiring check** — After adding or updating an icon, verify it is wired in both `PageTitle.astro` iconMap AND the `index.mdx` ComponentCard. Consider consolidating icon registration into a single source-of-truth file (e.g. `icons.ts`) that both components import.
- **Code example workflow** — `./site/examples/` (at the repo root, NOT `./site/src/content/docs/examples/`) is where verified Go source code lives. This is a sibling of `./site/src/`, not inside it.
  **STRUCTURE**: There is ONE `go.mod` at `./site/examples/go.mod` — NOT one per subdirectory. Each subdirectory is a Go package under that single module, not a separate module. This avoids bloating `go.work`. Use `package main` files where a complete example is needed, and regular package files or `_test.go` files for snippet verification.
  The coder agent must:
  1. Run `h2pp audit` to get the list of unverified code blocks
  2. For complete examples: create `./site/examples/<name>/main.go`
  3. For incomplete snippets: create test/wrapper files in `./site/examples/<name>/`
  4. Each `.go` file MUST have a Go comment at the top tracing back to the source doc, e.g.:
     `// Source: site/src/content/docs/components/grid-view.mdx:166`
     For files covering multiple snippets, list each source:
     `// Source: site/src/content/docs/components/grid-view.mdx:166,191,201`
  5. Run `go mod tidy`, `go vet ./...`, `go build ./...` from `./site/examples/`
  6. Fix any compilation errors by correcting the code to match the actual API
  7. Update the documentation page with the corrected code and add an HTML comment above each code block pointing to the verified source
  The docs-writer MUST NOT invent code examples directly in documentation. All code in docs must trace back to compiled source in `./site/examples/`.

## Constraints

- The documentation framework is Astro Starlight (MDX format)
- Match the tone and structure of existing documentation pages
- Do not modify Go source code — only documentation and site files
- Verify the site builds (`just site-build`) before reporting done
- This is PRE-RELEASE software. Nothing is "deprecated" or "legacy." If a package has been superseded, do NOT document it, do NOT add it to the sidebar, and do NOT create a doc page for it. Instead, flag it to the manager as source code that should be deleted from the repo. There are no legacy components in pre-release software — only current code and code that should be removed.

## Criteria

**Guideline for writing criteria:** Separate objective (mechanically testable) criteria from subjective ones. For any structural criterion, list the specific required sections or fields. Any task that produces observable, verifiable output must be stated as a positive criterion — never frame verifiable work as "implicit."

- Every component package (see Definitions) has a dedicated documentation page.
- Each component page references a distinct icon file; no two SVG files are identical copies.
- Each component page has an entry in `iconMap` in `PageTitle.astro`, appears in the navigation sidebar config, and has a `ComponentCard` on the home page (`index.mdx`).
- Home page lists all components with accurate one-line descriptions.
- The component count in the index page subtitle matches the actual number of component rows in the table.
- Site builds without errors (`just site-build`).
- Navigation sidebar includes every component.
- Every new page has: frontmatter (title, description), a Module or Type line, an `## Install` section (with fenced bash `go get` block), at least one code example, an `## API Reference` section covering exported types/methods, and a `## Related Components` section.
- No broken links or missing references.
- Component pages document all exported types, their fields, and key methods. Verify using `doc-go-repo -exclude-file doterr.go ./<package>` and diffing against the doc page.
- **NON-NEGOTIABLE**: Every fenced Go code block in `./site/src/content/docs/**/*.mdx` must be verified against the actual package API. All examples live under `./site/examples/` with a SINGLE `go.mod` (not one per subdirectory). Complete examples are `main.go` files; snippets are verified via test/wrapper files. Every `.go` file must have a Go comment tracing back to its source doc file and line. Every example must pass `go build` and `go vet`. If `./site/examples/` does not exist, the coder must create it. The evaluator MUST NOT negotiate this down to syntax-only checking. `./site/examples/` is at repo root (sibling of `./site/src/`), NOT inside `./site/src/content/docs/examples/`.
