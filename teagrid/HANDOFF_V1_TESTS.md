Handoff: Execute teagrid baseline implementation now

Repo:
- /Users/mikeschinkel/Projects/go-pkgs/go-tealeaves

Canonical plan files:
- /Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teagrid/PLAN.md
- /Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md

Important clarification:
- “go-tealeaves/teagrid/PLAN.md” means the repo-local file above (`teagrid/PLAN.md`).

Locked decisions (do not re-debate in this session):
1. v1 baseline teagrid is fork/copy from bubble-table (not placeholder skeleton).
2. v1 baseline implementation must be minimal-touch rename/rebrand/attribution (no feature refactors).
3. v1 baseline requires thorough testing before v0.1.0 tagging (test details come from parallel test-planning session artifact).
4. v2 work is blocked until pre-v2 checkpoint outcomes are complete (fork/rebrand/attribution + gomion v1 adoption + baseline tags
+ post-tag verification).
5. v2 strategy is “Charm v2 best-practices first”: evolve where viable, rewrite where needed.
6. Root-linked blockers should be tracked, but execution should prioritize progress over process overhead.

Primary objective for this implementation session:
- Deliver working v1 baseline `teagrid` in go-tealeaves so gomion can migrate off bubble-table on Charm v1 ASAP.

Implementation scope NOW:
1. Create `teagrid` module by copying/forking from local bubble-table source.
2. Rename package/module/docs from table/bubble-table -> teagrid.
3. Add/verify MIT attribution and license requirements in module docs.
4. Ensure module builds and test suite runs with thorough baseline coverage requirements.
5. Keep code changes minimal-touch beyond rename/fork necessities.
6. Prepare for gomion v1 integration flow (do not start v2 migration here).

Do NOT do in this session:
1. Charm v2 migration implementation.
2. Feature redesign/refactor cleanup unrelated to rename baseline.
3. Process-heavy blocker bookkeeping.

Expected deliverables:
1. `teagrid` module files in this repo with renamed package path.
2. README/license attribution updated for fork origin.
3. Passing baseline build/test results for teagrid.
4. Any required repo wiring updates (module lists/make/test integration) needed for baseline validation.

Current local git state to be aware of:
- `teagrid/PLAN.md` is untracked.
- `UPGRADE_V2_PLAN.md` is untracked.
- `teanotify/PLAN.md` is `AM` (already staged+modified).
- Do not revert unrelated existing work.
