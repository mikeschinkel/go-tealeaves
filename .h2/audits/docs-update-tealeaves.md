# Audit Report: go-tealeaves Documentation
**Date:** 2026-04-04
**Recipe:** docs-update-tealeaves

---

## 1. Packages with tea.Model Implementation

The following 12 `tea*` packages contain at least one type implementing the Bubble Tea v2 `tea.Model` interface:

| Package | Model Types |
|---------|-------------|
| `teacrumbs` | `BreadcrumbsModel` |
| `teadiff` | `SplitDiffModel` |
| `teafields` | `DropdownModel` |
| `teagrid` | `GridModel` |
| `teaguide` | `GuideModel` |
| `teahelp` | `HelpVisorModel` |
| `tealayout` | `StackLayoutModel` |
| `teamodal` | `ChoiceModel`, `ConfirmModel`, `ListModel[T]`, `MultiSelectModel[T]`, `ProgressModal` |
| `teanotify` | `NotifyModel` |
| `teastatus` | `StatusBarModel` |
| `teatext` | `TextSnipModel` |
| `teatree` | `TreeModel[T]`, `DrillDownModel[T]` |

**Non-model packages** (no `tea.Model`):
- `teacolor` — color constants and ANSI utilities (no model)
- `teadiffr` — diff renderer (render-only, no model)
- `teahilite` — syntax highlighter (utility, no model)
- `teapane` — pane widgets with `View()` but no `Init()`/`Update()` (not full models)
- `teaterm` — **research only**, no Go source files yet
- `teautils` — utility helpers (no model)

---

## 2. Documentation Pages for Qualifying Packages

All 12 qualifying packages have documentation pages (100% coverage):

| Package | Doc Page |
|---------|----------|
| `teacrumbs` | `components/breadcrumb-nav.mdx` |
| `teadiff` | `components/diff-viewer.mdx` |
| `teafields` | `components/dropdown-control.mdx` |
| `teagrid` | `components/grid-view.mdx` |
| `teaguide` | `components/guide-overlay.mdx` |
| `teahelp` | `components/help-visor.mdx` |
| `tealayout` | `components/layout-engine.mdx` |
| `teamodal` | `components/choice-dialog.mdx`, `components/confirm-dialog.mdx`, `components/list-dialog.mdx`, `components/multiselect-dialog.mdx`, `components/progress-dialog.mdx` |
| `teanotify` | `components/notification-view.mdx` |
| `teastatus` | `components/statusbar-view.mdx` |
| `teatext` | `components/text-selection.mdx` |
| `teatree` | `components/tree-view.mdx`, `components/drilldown-view.mdx` |

---

## 3. Missing Documentation Pages

**None.** All qualifying packages have at least one documentation page.

---

## 4. Documentation Page Section Audit

All 25 MDX pages in `site/src/content/docs/components/`:

| Page | Frontmatter | Install | Bash Block | Code Examples | API Reference | Related Components |
|------|:-----------:|:-------:|:----------:|:-------------:|:-------------:|:-----------------:|
| `breadcrumb-nav.mdx` | ✓ | ✓ | ✓ | 6 blocks | ✓ | **MISSING** |
| `choice-dialog.mdx` | ✓ | ✓ | ✓ | 2 blocks | ✓ | **MISSING** |
| `color-constants.mdx` | ✓ | ✓ | ✓ | 6 blocks | ✓ | **MISSING** |
| `confirm-dialog.mdx` | ✓ | ✓ | ✓ | 3 blocks | ✓ | **MISSING** |
| `diff-renderer.mdx` | ✓ | ✓ | ✓ | 7 blocks | ✓ | **MISSING** |
| `diff-viewer.mdx` | ✓ | ✓ | ✓ | 9 blocks | ✓ | **MISSING** |
| `drilldown-view.mdx` | ✓ | ✓ | ✓ | 2 blocks | ✓ | **MISSING** |
| `dropdown-control.mdx` | ✓ | ✓ | ✓ | 4 blocks | ✓ | **MISSING** |
| `grid-view.mdx` | ✓ | ✓ | ✓ | 7 blocks | ✓ | ✓ |
| `guide-overlay.mdx` | ✓ | ✓ | ✓ | 7 blocks | ✓ | **MISSING** |
| `help-visor.mdx` | ✓ | ✓ | ✓ | 3 blocks | ✓ | **MISSING** |
| `key-registry.mdx` | ✓ | ✓ | ✓ | 4 blocks | ✓ | **MISSING** |
| `layout-engine.mdx` | ✓ | ✓ | ✓ | 19 blocks | ✓ | **MISSING** |
| `list-dialog.mdx` | ✓ | ✓ | ✓ | 1 block ⚠️ | ✓ | **MISSING** |
| `multiselect-dialog.mdx` | ✓ | ✓ | ✓ | 1 block ⚠️ | ✓ | **MISSING** |
| `notification-view.mdx` | ✓ | ✓ | ✓ | 4 blocks | ✓ | **MISSING** |
| `pane-widgets.mdx` | ✓ | ✓ | ✓ | 5 blocks | ✓ | **MISSING** |
| `positioning.mdx` | ✓ | ✓ | ✓ | 2 blocks | ✓ | **MISSING** |
| `progress-dialog.mdx` | ✓ | ✓ | ✓ | 1 block ⚠️ | ✓ | **MISSING** |
| `statusbar-view.mdx` | ✓ | ✓ | ✓ | 5 blocks | ✓ | **MISSING** |
| `syntax-highlighting.mdx` | ✓ | ✓ | ✓ | 8 blocks | ✓ | **MISSING** |
| `term-renderer.mdx` | ✓ | ✓ | ✓ | 4 blocks | ✓ | **MISSING** |
| `text-selection.mdx` | ✓ | ✓ | ✓ | 1 block ⚠️ | ✓ | **MISSING** |
| `theming.mdx` | ✓ | ✓ | ✓ | 5 blocks | ✓ | **MISSING** |
| `tree-view.mdx` | ✓ | ✓ | ✓ | 3 blocks | ✓ | **MISSING** |

**Summary of issues:**
- **24/25 pages missing `## Related Components` section** — only `grid-view.mdx` has it
- **4 pages with only 1 code block (thin examples):** `list-dialog.mdx`, `multiselect-dialog.mdx`, `progress-dialog.mdx`, `text-selection.mdx`

---

## 5. Home Page ComponentCards

File: `site/src/content/docs/index.mdx`
Total ComponentCards: **20**

All 20 ComponentCards are present and descriptions are accurate. Notable items:

- **Terminal Renderer** — description correctly states "Planned component" (no Go source in `teaterm/`, only `TEATERM_RESEARCH.md`)
- **Diff Renderer** — description correctly says "Deprecated — superseded by teadiff"
- **Progress Dialog** — `type="ProgressModal"` matches the actual struct name in source

No inaccurate descriptions detected.

---

## 6. Sidebar Configuration

File: `site/astro.config.mjs`

All 25 component pages are present in the sidebar under the correct categories:

- **Views (9):** Grid View, Tree View, Drilldown View, Status Bar, Notifications, Diff Viewer, Diff Renderer (Legacy), Terminal Renderer, Breadcrumb Nav
- **Dialogs (6):** Confirm Dialog, Choice Dialog, List Dialog, Progress Dialog, MultiSelect Dialog, Guide Overlay
- **Controls (1):** Dropdown Control
- **Text (2):** Text Selection, Syntax Highlighting
- **Layout (2):** Layout Engine, Pane Widgets
- **System (5):** Help Visor, Key Registry, Theming, Color Constants, Positioning

**No missing sidebar entries.**

---

## 7. PageTitle.astro Icon Map

File: `site/src/components/PageTitle.astro`

All 25 component pages have an entry in the `iconMap`. No missing icon mappings.

---

## 8. SVG Icons in site/public/icons/

26 SVG files found:

```
icon-choice.svg        (707B)
icon-confirm.svg       (432B)
icon-helpvisor.svg     (338B)
icon-hilite.svg        (663B)
icon-keyregistry.svg   (370B)
icon-list.svg          (694B)
icon-positioning.svg   (627B)
icon-progress.svg      (479B)
icon-teacolor.svg      (1.0K)
icon-teacrumbs.svg     (491B)
icon-teadiff.svg       (1.7K)
icon-teadiffr.svg      (1.4K)
icon-teadrldwn.svg     (687B)
icon-teadrpdwn.svg     (311B)
icon-teafields.svg     (522B)
icon-teagrid.svg       (952B)
icon-teaguide.svg      (762B)
icon-tealayout.svg     (862B)
icon-teamodal.svg      (1.2K)
icon-teanotify.svg     (482B)
icon-teapane.svg       (911B)
icon-teastatus.svg     (502B)
icon-teaterm.svg       (479B)
icon-teatextsel.svg    (450B)
icon-teatree.svg       (712B)
icon-theming.svg       (313B)
```

**No duplicate icons detected** (all files have unique content).

**Note:** `icon-progress.svg` (479B) and `icon-teaterm.svg` (479B) have the same file size — worth verifying these are visually distinct.

---

## 9. Examples Directories

**Root `examples/` directory:** EXISTS
Contains 6 symlinks pointing to package example subdirectories:
- `teagrid` → `../teagrid/examples`
- `teamodal` → `../teamodal/examples`
- `teanotify` → `../teanotify/examples`
- `teastatus` → `../teastatus/examples`
- `teatree` → `../teatree/examples`
- `teautils` → `../teautils/examples`

Note: 3 symlinks that previously existed were deleted per git status (`teadiffview`, `teadrpdwn`, `teatxtsnip`).

**`site/examples/` directory:** EXISTS but is **EMPTY** — no content.

---

## 10. doc-go-repo Command Check

Binary found at: `/usr/local/bin/doc-go-repo`

Tested against all 12 qualifying packages — all exit 0 successfully:

| Package | Status |
|---------|--------|
| `teacrumbs` | ✓ works |
| `teadiff` | ✓ works |
| `teafields` | ✓ works |
| `teagrid` | ✓ works |
| `teaguide` | ✓ works |
| `teahelp` | ✓ works |
| `tealayout` | ✓ works |
| `teamodal` | ✓ works |
| `teanotify` | ✓ works |
| `teastatus` | ✓ works |
| `teatext` | ✓ works |
| `teatree` | ✓ works |

Command: `doc-go-repo -exclude-file doterr.go ./<package>`

---

## Summary of Actionable Issues

| Priority | Issue | Scope |
|----------|-------|-------|
| P1 | 24/25 doc pages missing `## Related Components` section | All pages except `grid-view.mdx` |
| P2 | 4 doc pages have only 1 code example (thin) | `list-dialog.mdx`, `multiselect-dialog.mdx`, `progress-dialog.mdx`, `text-selection.mdx` |
| P3 | `site/examples/` directory is empty | Should contain site-facing example gallery content |
| P3 | `icon-progress.svg` and `icon-teaterm.svg` same file size (479B) — verify visually distinct | Two icons |
