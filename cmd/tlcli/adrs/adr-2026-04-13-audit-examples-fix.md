# ADR-2026-04-13: Unified --fix for tlcli audit examples

## Status
**Accepted**

## Date
2026-04-13

## Context

The example verification system uses bidirectional references between MDX doc
pages and Go example files:

- MDX: `{/* verified: site/examples/grid-view/compile_test.go:22 */}`
- Go:  `// Source: site/src/content/docs/components/grid-view.mdx:22#a53b0433`

When someone edits an MDX file (adding or removing lines above a code block),
the line numbers in both references become stale. CRC32 content hashes in the
Source comment enable detection of line shifts (content matches but line number
differs).

Three problems existed:

1. **`tlcli audit` had vestigial stale-block detection** based on file
   modification timestamps. This was unreliable (any `touch` triggered it)
   and redundant after `tlcli audit examples` was added with content-hash-based
   detection. The mtime-based section reported hundreds of false positives.

2. **Hash seeding required a separate `tlcli seed-hashes` command.** Operators
   had to know to run it before lines shifted, creating a chicken-and-egg
   problem: if you forgot to seed hashes and then lines shifted, `seed-hashes`
   could no longer find the blocks to hash (it required exact line matches).

3. **Source comments without hashes fell into "Needs Example Written"** — a
   misleading category. The example file existed; only the line number was
   stale with no hash to confirm the match.

## Decision

1. **Remove the Code Example Verification section from `tlcli audit`.** All
   example verification is handled by `tlcli audit examples`.

2. **Add fuzzy nearby-line matching (plus/minus 3 lines)** to the detection pass.
   When no exact line match and no hash match exist, search Source entries at
   adjacent line numbers. This handles the common case of small line shifts
   without hashes.

3. **Add a "Line Numbers Stale" category** distinct from "Needs Example Written"
   for cases where the example file exists and the Source comment references
   the right MDX file, but line numbers don't match.

4. **Make `--fix` handle all fixable categories in one pass:**
   - **Line Numbers Stale (no hash):** Update line number + seed hash.
   - **Line Numbers Stale (has hash):** Update line number, preserve hash.
   - **Needs Hash Seeding:** Seed hash at exact matching line.

   This eliminates the need to ever run `tlcli seed-hashes` separately.

5. **Report fixable vs unfixable counts** in the summary line. When fixable
   issues exist, the summary says `N issues (M fixable with --fix)`.

## Consequences

- The operator workflow is: run `bin/tlcli audit examples`, and if there are
  fixable issues, re-run with `--fix`. One command handles everything.
- `tlcli seed-hashes` still exists for backward compatibility but is no longer
  referenced in the recipe.
- A prior bug in the fix logic (stripping newlines when rewriting Source
  comments via `strings.TrimSpace`) was caught and fixed. The rewrite now
  preserves leading/trailing whitespace including newlines on each
  comma-separated entry.
