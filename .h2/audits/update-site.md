# Audit Report: update-site

**Date:** 2026-04-11  
**Auditor:** auditor-tealeaves  
**Recipe:** update-site

---

## Classification Table

### Components

Packages with Bubble Tea v2 component types (Init/Update/View). Source of truth for ComponentCards in index.mdx.

| Package    | Model Type(s)                                                                   |
|------------|---------------------------------------------------------------------------------|
| teacrumbs  | BreadcrumbsModel                                                                |
| teadiff    | SplitDiffModel                                                                  |
| teafields  | DropdownModel                                                                   |
| teagrid    | GridModel                                                                       |
| teaguide   | GuideModel                                                                      |
| teahelp    | HelpVisorModel                                                                  |
| tealayout  | StackLayoutModel                                                                |
| teamodal   | ChoiceModel, ConfirmModel, ListModel[T], MultiSelectModel[T], ProgressModel     |
| teanotify  | NotifyModel                                                                     |
| teastatus  | StatusBarModel                                                                  |
| teatext    | TextSnipModel                                                                   |
| teatree    | DrillDownModel[T], TreeModel[T]                                                 |

**Total component packages: 12** (with 17 model types across them)

All naming conventions pass `tlcli models -check`.

### Foundations

Packages that do NOT contain any Bubble Tea component type. Must appear only as SystemCards in the System section, never as ComponentCards.

| Package   | Purpose                                                                            |
|-----------|------------------------------------------------------------------------------------|
| teacolor  | Named semantic color constants (ANSI/xterm/hex) for consistent terminal UI styling |
| teahilite | Syntax highlighting via Chroma; produces ANSI-colored output                       |
| teapane   | Ready-made tealayout-compatible panes: PlainPane, StyledPane, ScrollPane           |
| teaterm   | **STUB ONLY** — contains only TEATERM_RESEARCH.md, no Go source or go.mod         |
| teautils  | Shared utilities: key registry, theming/palettes, positioning/ANSI helpers         |

---

## Exception Report

### P0 — ComponentCard Misclassifications in index.mdx

The following ComponentCards reference types that are NOT returned by `tlcli models`. They are foundation packages and must be removed from the Components section. Per recipe: foundations appear ONLY as SystemCards in the System section.

| ComponentCard type | Package   | Issue                                                                           |
|--------------------|-----------|---------------------------------------------------------------------------------|
| `Highlighter`      | teahilite | Foundation — no component type. Must be SystemCard only.                        |
| `SemanticColor`    | teacolor  | Foundation — no component type. Must be SystemCard only.                        |
| `Terminal`         | teaterm   | Foundation/stub — no Go source at all. See "teaterm stub" below.                |
| `PlainPane`        | teapane   | Foundation — no component type. Must be SystemCard only.                        |
| `KeyRegistry`      | teautils  | Foundation — no component type. Must be SystemCard only.                        |
| `Theme`            | teautils  | Foundation — no component type. Must be SystemCard only.                        |
| `(functions)`      | teautils  | Foundation — no component type. Must be SystemCard only.                        |
| `Pane`             | tealayout | Wrong type — tealayout IS a component package, but `Pane` is not the component type. Should be `StackLayoutModel`. |

Also: the System section in index.mdx is missing SystemCards for teacolor, teahilite, and teaterm (if it stays). Currently only teapane, tealayout-helpers, key-registry, theming, and positioning have SystemCards.

### P0 — Wrong doc-pages.yaml Mappings

These mapping errors cause the audit tool to attribute doc page staleness to the wrong package and to miss real staleness.

| Issue | Detail |
|-------|--------|
| `teafields → drilldown-view.mdx` is wrong | drilldown-view.mdx documents `DrillDownModel[T]` from **teatree**, not teafields. Remove from teafields entry; add to teatree entry. |
| `teatree` missing `drilldown-view.mdx` | The drilldown page exists and is correct but not mapped to its source package. |
| `tealayout → positioning.mdx` and `tealayout → theming.mdx` are wrong | Both pages document **teautils** types (`Theme`, `Palette`, positioning functions). Must move to teautils entry. |
| `teautils → key-registry.mdx` only | teautils also owns positioning.mdx and theming.mdx. All three should be in the teautils mapping. |
| `teadiffr → diff-renderer.mdx` is stale | Package teadiffr was deleted (commit 4532639, superseded by teadiff). Entry must be removed from doc-pages.yaml. |

### P1 — Stale Doc Pages (source newer than doc)

All of these need doc review/update to match current API. Listed with packages AFTER correcting the mappings above.

| Package   | Doc Page               | Source Timestamp            | Doc Timestamp               |
|-----------|------------------------|-----------------------------|-----------------------------|
| teadiff   | diff-viewer.mdx        | 2026-04-04T10:23:09-04:00   | 2026-04-04T10:14:26-04:00   |
| teafields | dropdown-control.mdx   | 2026-04-04T10:11:40-04:00   | 2026-04-04T04:00:31-04:00   |
| teaguide  | guide-overlay.mdx      | 2026-04-04T11:52:57-04:00   | 2026-04-04T03:52:18-04:00   |
| tealayout | layout-engine.mdx      | 2026-04-07T11:37:00-04:00   | 2026-04-04T04:00:39-04:00   |
| teautils  | theming.mdx            | 2026-04-07T11:37:00-04:00   | 2026-04-04T10:13:57-04:00   |
| teautils  | positioning.mdx        | 2026-04-07T11:37:00-04:00   | 2026-04-04T03:52:26-04:00   |
| teamodal  | progress-dialog.mdx    | 2026-04-10T22:07:24-04:00   | 2026-04-04T03:52:26-04:00   |
| teamodal  | multiselect-dialog.mdx | 2026-04-10T22:07:24-04:00   | 2026-04-04T03:52:23-04:00   |
| teamodal  | list-dialog.mdx        | 2026-04-10T22:07:24-04:00   | 2026-04-04T03:52:22-04:00   |
| teamodal  | confirm-dialog.mdx     | 2026-04-10T22:07:24-04:00   | 2026-04-04T04:00:41-04:00   |
| teamodal  | choice-dialog.mdx      | 2026-04-10T22:07:24-04:00   | 2026-04-04T03:52:11-04:00   |
| teapane   | pane-widgets.mdx       | 2026-04-05T10:23:49-04:00   | 2026-04-04T03:52:24-04:00   |
| teatext   | text-selection.mdx     | 2026-04-04T11:58:51-04:00   | 2026-04-04T03:52:30-04:00   |
| teatree   | tree-view.mdx          | 2026-04-04T12:00:57-04:00   | 2026-04-04T04:00:44-04:00   |
| teautils  | key-registry.mdx       | 2026-04-04T12:01:35-04:00   | 2026-04-04T03:52:20-04:00   |

Note: drilldown-view.mdx is stale per audit but was attributed to teafields; correct package is teatree (src: 2026-04-04T12:00:57-04:00).

### P1 — Broken Icon Reference in PageTitle.astro

`PageTitle.astro` maps `'diff-renderer'` to `/go-tealeaves/icons/icon-teadiffr.svg` but:
- `icon-teadiffr.svg` does not exist in `site/public/icons/`
- The `diff-renderer` page does not exist (teadiffr was deleted)
- This entry should be removed from iconMap.

### P1 — Orphaned Doc Page

`site/src/content/docs/components/index.md` — listed as orphaned by tlcli audit (no package maps to it). This is likely an old overview page superseded by `index.mdx`. Should be reviewed and removed or reassigned.

### P1 — teaterm Is a Stub → MANAGER DECISION: DELETE

`teaterm/` contains only `TEATERM_RESEARCH.md`. There is no `go.mod`, no `.go` source files, and no component type.

**Manager ruling (2026-04-11):** Per the recipe's pre-release policy ("no deprecated or legacy code — only current code or code to be deleted"):
- `site/src/content/docs/components/term-renderer.mdx` must be **DELETED**
- `teaterm` must NOT appear in the sidebar or as a ComponentCard
- The `teaterm/` directory is flagged to the repo owner as source that should be deleted

Action for docs-writer: delete term-renderer.mdx, remove from sidebar config, remove ComponentCard from index.mdx.  
Action for repo owner: delete `teaterm/` directory from the repository.

### P2 — teacolor and teahilite Missing SystemCards

teacolor and teahilite are foundation packages but have NO SystemCard in index.mdx System section. They appear only as ComponentCards (P0 above). Once ComponentCards are removed they will be entirely absent from the home page.

### P2 — Missing teautils SystemCards for Theming and Positioning

Currently `theming.mdx` and `positioning.mdx` are in the System section as separate SystemCards. This is correct, but they need to be re-attributed to teautils (not tealayout) after doc-pages.yaml is fixed.

### P3 — 0 of 176 Code Examples Verified

`tlcli audit` reports 176 fenced Go code blocks, 0 verified, 0 stale.  
`site/examples/` does not exist.  
All code examples in all .mdx files are unverified. Per recipe acceptance criteria this is non-negotiable — coder agent must create `site/examples/` with a single go.mod and verify all examples.

---

## Summary Counts

| Severity | Count | Category |
|----------|-------|----------|
| P0 | 8 | ComponentCard misclassifications (8 types from foundation packages shown as components) |
| P0 | 5 | doc-pages.yaml mapping errors |
| P1 | 15 | Stale doc pages (source newer than doc) |
| P1 | 1 | Broken icon reference (icon-teadiffr.svg missing) |
| P1 | 1 | Orphaned doc page (index.md) |
| P1 | 1 | teaterm is a stub (no Go source, has a doc page) |
| P2 | 2 | Foundation packages with no SystemCard (teacolor, teahilite) |
| P3 | 1 | 176 unverified code examples, site/examples/ does not exist |

---

## Packages With No Doc Page

None. All 17 tea* packages with Go source have at least one mapped doc page.

## Packages Whose Doc Page Is OK (no issues)

| Package   | Doc Page(s)            |
|-----------|------------------------|
| teacrumbs | breadcrumb-nav.mdx     |
| teagrid   | grid-view.mdx          |
| teastatus | statusbar-view.mdx     |
| teanotify | notification-view.mdx  |
| teahilite | syntax-highlighting.mdx |

All others have staleness, mapping, or classification issues noted above.

---

## Action Items for Downstream Agents

### docs-writer
1. Fix 8 ComponentCard misclassifications in index.mdx (move to SystemCards or correct types)
2. Add SystemCards for teacolor and teahilite in System section
3. Fix `Pane` → `StackLayoutModel` for tealayout ComponentCard
4. Update all 15 stale doc pages to match current API (use `doc-go-repo -exclude-file doterr.go ./<package>`)
5. Review/update index page component count subtitle to match actual component count

### coder
1. Create `site/examples/` with single go.mod
2. Verify all 176 code blocks (start with the 15 stale pages, then remaining)
3. Each example needs source comment tracing back to .mdx file and line

### manager / repo owner
1. **teaterm:** delete the `teaterm/` directory from the repo (manager ruling 2026-04-11)
2. Fix doc-pages.yaml: correct all 5 mapping errors listed above
3. Remove `'diff-renderer'` entry from PageTitle.astro iconMap
4. Delete orphaned `site/src/content/docs/components/index.md` or reassign
