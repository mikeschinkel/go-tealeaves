# ADR-2026-04-12: tlcli exports — API Extraction Without MDX Diffing

## Status
**Accepted**

## Date
2026-04-12

## Context

The documentation update recipe needed a way to verify that MDX doc pages cover
all exported Go symbols. An initial design (`PROMPT-doc-diff.md`) proposed a
`tlcli doc-diff` command that would extract the Go API, parse MDX pages, and
report gaps — all in one tool.

Two problems emerged during implementation:

1. **MDX parsing required fragile heuristics.** Extracting "documented symbols"
   from MDX is inherently fuzzy — backtick spans contain external types
   (`tea.Cmd`, `lipgloss.Style`), code blocks contain stdlib usage, and there
   is no reliable way to distinguish package-local symbols from external ones
   without maintaining a brittle exclusion list.

2. **The recipe workflow already has an agent doing the comparison.** The update
   recipe instructs agents to "run the tool and diff the output against the
   documentation page." The agent reads both the tool output and the MDX — it
   does not need the tool to do both.

Separately, the initial implementation only supported single-package extraction.
The source tool (`doc-go-repo`) scans all modules in a repo, and the recipe
needs full-repo output.

## Decision

1. **Name the command `exports`, not `doc-diff`.** It extracts and outputs the
   exported API. It does not read MDX or compute diffs. The agent performs the
   comparison.

2. **Port the full doc-go-repo logic** — module discovery (`go mod edit -json`),
   package listing (`go list -json`), AST parsing (`go/parser`, `go/doc`) —
   so the command scans all modules, not just one package.

3. **Collect `doc.Type.Funcs`** (type-associated constructor functions) in
   addition to `doc.Type.Methods`. Go's `go/doc` package places functions like
   `NewGridModel()` and `DefaultTreeKeyMap()` in `Type.Funcs`, not in the
   top-level `Funcs` list or `Type.Methods`. The original doc-go-repo had the
   same gap. Also collect `Type.Consts` and `Type.Vars` for completeness.

4. **Add `--api-only` flag** that excludes `cmd/*`, `site/*`, and `*/examples*`
   in one shot — the common non-library paths. This replaces four separate
   `--exclude-path` flags in the recipe.

5. **Add `--exclude-path` and `--include-path`** repeatable glob flags for
   custom filtering. `--exclude-path` takes precedence over `--include-path`.
   Matching is against repo-relative paths (e.g. `teagrid`, `cmd/tlcli`).

## Consequences

- The recipe uses `bin/tlcli exports -api-only` for bulk API extraction and
  `bin/tlcli exports -include-path=<pkg>` for single-package spot checks.
- `tlcli seed-hashes` remains available but is no longer referenced by the
  recipe (superseded by `audit examples --fix`).
- MDX diffing is the agent's responsibility, keeping the tool simple and the
  heuristics where they can be adjusted per-context.
