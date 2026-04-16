# ADR-2026-04-13: Config Package Verification in tlcli audit

## Status
**Accepted**

## Date
2026-04-13

## Context

`site/config.yaml` is the source of truth for which packages appear on the
documentation home page — components listed as `package.Type` entries (e.g.
`teagrid.GridModel`) and foundations listed by package name (e.g. `teapane`).

If a package is removed or renamed in the Go source but config.yaml is not
updated, the home page references a phantom package. This was a manual check
with no tooling to catch it.

## Decision

Add a "Config Package Verification" section to `tlcli audit` that:

1. Parses `site/config.yaml` to extract all package names (from both
   `components` entries and `foundations` entries).
2. For each package name, verifies that the corresponding directory exists
   under the repo root and contains at least one non-test `.go` file.
3. Reports any missing or empty packages and adds them to the issue count.

## Consequences

- `bin/tlcli audit` now catches stale config.yaml entries automatically.
- The check is fast (directory stat + readdir) and adds no external
  dependencies — config.yaml is already parsed with `gopkg.in/yaml.v3`
  which is in go.mod.
