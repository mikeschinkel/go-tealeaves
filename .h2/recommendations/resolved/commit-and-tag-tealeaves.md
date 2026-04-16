# Resolved Recommendations: commit-and-tag-tealeaves

## 2026-04-15

### Criteria Referenced

| ID  | Criterion |
|-----|-----------|
| C1  | Every module with file changes (within its module directory) since its last tag has a new commit and tag |
| C2  | No uncommitted changes remain after all commits (git status is clean) |
| C3  | Module tags follow `<module>/vX.Y.Z` format; root tag follows `vX.Y.Z` format |
| C4  | Breaking changes result in minor version bumps with patch reset to 0; non-breaking in patch bumps |
| C5  | Each commit contains files from exactly one module directory (or root non-module files only) |
| C6  | Modules are committed in dependency order (dependencies before dependents) |
| C7  | Root non-module files are committed and tagged separately from all modules |
| C8  | No files are pushed to the remote |
| C9  | Each tag points to the correct commit (the one that staged that module's changes) |
| C10 | Version numbers are continuous — no skipped or backwards versions relative to last tag |
| C11 | No duplicate tags exist for the same module |
| C12 | Commit messages follow project convention |

### R1 — ACCEPTED: Add pre-flight staged-file check
Before staging anything, the recipe now requires: run `git diff --cached --name-only`. If any files are already staged, run `git restore --staged .` to unstage them, then verify with `git status --short`. Prevents previously staged files from contaminating the first commit attempt.
Source: releaser-tealeaves post-mortem

### R2 — ACCEPTED: Analyst scope must cover working-tree changes, not only committed history
The analyst must also run `git status --short` and account for every modified, added, or deleted file in the working tree. Recipe updated to require both `git diff <tag>` (working-tree form) and `git status --short` together.
Source: releaser-tealeaves post-mortem

### R3+R12 — ACCEPTED (revised): cmd/tlcli versioned independently
Original R3 proposed adding cmd/ to root-commit staging. Original R12 (evaluator finding, C5 — Each commit contains files from exactly one module directory) noted cmd/ sub-modules were bundled with root. Revised: user deleted all cmd/ modules except cmd/tlcli. Recipe now states cmd/tlcli is an independent module versioned separately (`cmd/tlcli/vX.Y.Z`), not part of the root commit.
Source: releaser-tealeaves post-mortem + evaluator finding + user decision

### R5+R6+R8 — ACCEPTED (merged): Fix diff command, clarify git status, explain diff forms
R5: Changed recipe from `git diff <tag>..HEAD -- <module>/` (committed history only) to `git diff <tag> -- <module>/` (working-tree form capturing all changes). R6: Clarified that `git status --short` shows modified (M) and deleted (D) tracked files, not just untracked. R8: Added explanatory text about why the working-tree diff form is always correct for release analysis. All three address the same root cause: the analyst not seeing uncommitted changes.
Source: analyst-tealeaves post-mortem

### R7 — ACCEPTED: Add step to detect deleted modules
Added detection step: compare current module list (go.mod directories) against tag history (`git tag --list '*/v*'`). Any module with prior tags but no current go.mod is classified as DELETED — commit the deletion without a new tag.
Source: analyst-tealeaves post-mortem

### R9 — ACCEPTED: Add guidance for slotting new findings into a partially-completed release
Added guidance: if modules are already tagged mid-run when new changed modules are discovered, slot them into the correct dependency position among remaining untagged modules. Do not append to the end. Note which already-committed modules each depends on.
Source: analyst-tealeaves post-mortem

### R10 — ACCEPTED: Add classification rule for examples-only changes
Added rule: changes exclusively in `examples/` subdirectories (binary rebuilds, import updates, go.mod/go.sum updates, main.go tweaks) are always NON-BREAKING and warrant a patch bump.
Source: analyst-tealeaves post-mortem

### R11 — ACCEPTED (logged): C6 (Modules committed in dependency order) — Three dependency-order violations
tealayout and teapane were committed before teacrumbs; teatree before teafields. Root cause was the incomplete first analysis pass (working-tree diff issue). Addressed by R2 and R5+R6+R8 — no separate recipe change needed.
Source: rel-eval-tealeaves evaluator finding

### R4 — ACCEPTED: .h2/ session files created during the release run are included in root commit
Session files written by agents during the release run are included in the root commit. This keeps git status clean after the release and preserves the agent state that produced it.
Source: releaser-tealeaves post-mortem
