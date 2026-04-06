---
pod: commit-and-tag
pod_vars:
  codename:
  working_dir: "."
  language: "Go"
---
# Commit and Tag Multi-Module Release

## Goal

Commit and tag all changed modules in this multi-module Go repository
with appropriate SemVer version bumps. Each module gets its own commit
and tag. Root non-module files are committed and tagged separately.

## Scope

- All Go modules with changes since their last tag (directories
  containing go.mod, excluding the root which has no go.mod)
- Example modules under `*/examples/` — these are part of their
  parent module, not tagged independently
- Root non-module files: README.md, LICENSE, justfile, go.work,
  go.work.sum, CONTRIBUTING.md, ROADMAP.md, TODO.md, CLAUDE.md,
  docs/, adrs/, site/, cmd/, test/, assets/, .h2/

## Instructions

1. The change-analyst examines every module for changes since its
   last tag. For each changed module, the analyst classifies changes
   as breaking or non-breaking by inspecting the Go API surface
   (exported types, functions, methods, constants).

2. The analyst reports findings to the manager with proposed version
   bumps for each module and for root files.

3. The manager reviews the analyst's report, adjusts if needed,
   and approves the version plan.

4. The releaser executes the approved plan: for each module in
   dependency order, stage files, commit, and tag. Then handle
   root non-module files the same way.

5. The release-evaluator independently verifies all commits and
   tags are correct, complete, and follow conventions.

6. If the evaluator finds issues, the manager directs fixes and
   re-verification until all P0/P1 issues are resolved.

## Constraints

- **No pushing**: All commits and tags are LOCAL only. The user
  will push when ready.
- **Pre-v1.0 SemVer**: Non-breaking changes bump patch (0.0.+1).
  Breaking changes bump minor (0.+1.0). New modules start at v0.1.0.
- **One commit per module**: Each module's changes go in exactly one
  commit. Do not mix files from different modules in a single commit.
- **Dependency order**: Modules that are dependencies of other modules
  must be committed and tagged first (e.g., teautils before teagrid).
- **Explicit staging**: Never use `git add .` or `git add -A`. Always
  stage files explicitly by path or by module directory.
- **Module tags**: Use `<module>/v0.X.Y` format (e.g., `teagrid/v0.7.0`).
- **Root tags**: Use `v0.X.Y` format for non-module files.
- **Skip unchanged**: Do not commit or tag modules with no changes.

## Criteria

- Every module with changes since its last tag has a new commit and tag
- No uncommitted changes remain after all commits (`git status` is clean)
- All tags follow the correct SemVer format (`<module>/vX.Y.Z` or `vX.Y.Z`)
- Breaking changes result in minor version bumps; non-breaking in patch bumps
- Each commit contains files from exactly one module (or root only)
- Modules are committed in dependency order (dependencies before dependents)
- Root non-module files are committed and tagged separately from modules
- No files are pushed to the remote
