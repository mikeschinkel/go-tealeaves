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
- `cmd/tlcli` is an independent module — version and tag it
  separately (`cmd/tlcli/vX.Y.Z`), not as part of the root commit
- Root non-module files: README.md, LICENSE, justfile, go.work,
  go.work.sum, CONTRIBUTING.md, ROADMAP.md, TODO.md, CLAUDE.md,
  docs/, adrs/, site/, test/, assets/, .h2/
- `.h2/` session files created during the release run itself ARE
  included in the root commit

## Instructions

0. **Pre-flight check**: Before staging anything, run
   `git diff --cached --name-only`. If any files are already staged,
   run `git restore --staged .` to unstage them, then verify with
   `git status --short`. This prevents previously staged files from
   contaminating the first commit attempt.

1. The change-analyst examines every module for changes since its
   last tag. The analyst MUST use the **working-tree diff form**
   (`git diff <tag> -- <module>/`, without `..HEAD`) to compare
   the tag directly against the current working tree. This captures
   committed, staged, AND unstaged changes — not just committed
   history. The `..HEAD` form (`git diff <tag>..HEAD`) only shows
   committed changes and will miss uncommitted work.

   Additionally, the analyst MUST run `git status --short` to
   identify ALL changed files — this shows modified (`M`), deleted
   (`D`), and untracked (`??`) files, not just untracked ones.
   Both commands are required together for a complete picture.

   For each changed module, the analyst classifies changes as
   breaking or non-breaking by inspecting the Go API surface
   (exported types, functions, methods, constants).

   **Deleted module detection**: Compare the current module list
   (directories containing go.mod) against all modules with prior
   tags:
   ```
   git tag --list '*/v*' --sort=-v:refname | sed 's|/v[0-9].*||' | sort -u
   ```
   Any name appearing in tag history but having no current go.mod
   directory is a deleted module. Classify it as DELETED and
   recommend committing the deletion without creating a new tag.

   **Examples-only changes**: Changes exclusively in `examples/`
   subdirectories (binary rebuilds, import path updates, go.mod/
   go.sum updates, main.go tweaks) are always NON-BREAKING and
   warrant a patch bump.

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

   **Late-discovered modules**: If some modules are already tagged
   mid-run when new changed modules are discovered, provide only
   the remaining untagged modules in dependency order. Note which
   already-committed modules each new module depends on so the
   releaser can verify ordering constraints are still satisfied.
   Do NOT simply append late discoveries to the end of the list —
   slot them into the correct dependency position.

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
