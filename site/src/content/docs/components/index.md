---
title: Components
description: Overview of the 25 Tea Leaves components — production-ready Bubble Tea building blocks for dialogs, views, app chrome, and utilities.
---

Tea Leaves ships components across Go modules. Each module has its own `go.mod` — you install only what you need.

## Dialogs & Overlays

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Confirm Dialog](/go-tealeaves/components/confirm-dialog/) | `ConfirmModel` | teamodal | Yes/No and OK confirmation dialogs |
| [Choice Dialog](/go-tealeaves/components/choice-dialog/) | `ChoiceModel` | teamodal | Multi-option selection with hotkeys |
| [List Dialog](/go-tealeaves/components/list-dialog/) | `ListModel[T]` | teamodal | Editable list with inline editing and CRUD |
| [Progress Dialog](/go-tealeaves/components/progress-dialog/) | `ProgressModal` | teamodal | Progress indicator with cancel/background |
| [MultiSelect Dialog](/go-tealeaves/components/multiselect-dialog/) | `MultiSelectModel[T]` | teamodal | Generic checkbox list modal with scrolling and configurable buttons |
| [Guide Overlay](/go-tealeaves/components/guide-overlay/) | `GuideModel` | teaguide | Context-aware workflow guidance with direct key dispatch |

## Views

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Data Grid](/go-tealeaves/components/grid-view/) | `GridModel` | teagrid | Sorting, filtering, pagination, row selection |
| [Tree View](/go-tealeaves/components/tree-view/) | `TreeModel[T]` | teatree | Expand/collapse, pluggable node providers |
| [Drilldown View](/go-tealeaves/components/drilldown-view/) | `DrillDownModel[T]` | teatree | Interactive dependency path viewer |
| [Diff Viewer](/go-tealeaves/components/diff-viewer/) | `SplitDiffModel` / `TUIRenderer` | teadiff | Side-by-side interactive diff and condensed diff rendering |
| [Terminal Renderer](/go-tealeaves/components/term-renderer/) | `Terminal`, `ProcessViewer` | teaterm | Full PTY emulator and piped process viewer — planned, not yet released |
| [Breadcrumb Nav](/go-tealeaves/components/breadcrumb-nav/) | `BreadcrumbsModel` | teacrumbs | Mouse-aware breadcrumb trail with click and hover events |

## Controls

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Dropdown](/go-tealeaves/components/dropdown-control/) | `DropdownModel` | teafields | Dropdown selection with smart positioning |

## App Chrome

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Status Bar](/go-tealeaves/components/statusbar-view/) | `StatusBarModel` | teastatus | Two-zone status bar with menus and indicators |
| [Notifications](/go-tealeaves/components/notification-view/) | `NotifyModel` | teanotify | Toast notifications with auto-dismiss and color fade |
| [Help Visor](/go-tealeaves/components/help-visor/) | `HelpVisorModel` | teahelp | Paginated keyboard shortcuts overlay with category ordering |

## Layout

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Layout Engine](/go-tealeaves/components/layout-engine/) | `Pane`, `Element[T]`, `Layout` | tealayout | Declarative Percent/Fixed/Flex/Fit layout with focus management |
| [Pane Widgets](/go-tealeaves/components/pane-widgets/) | `StyledPane`, `ScrollPane`, `PlainPane` | teapane | Ready-made tealayout-compatible pane widgets |

## Utilities

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Key Registry](/go-tealeaves/components/key-registry/) | `KeyRegistry` | teautils | Centralized key binding management |
| [Theming](/go-tealeaves/components/theming/) | `Theme` / `Palette` | teautils | Consistent colors across all components |
| [Color Constants](/go-tealeaves/components/color-constants/) | `SemanticColor` | teacolor | Named ANSI 256 colors and zero-allocation style caching |
| [Positioning](/go-tealeaves/components/positioning/) | (functions) | teautils | ANSI-aware centering and measurement |
| [Text Selection](/go-tealeaves/components/text-selection/) | `TextSnipModel` | teatext | Textarea with Shift+Arrow selection and clipboard |
| [Syntax Highlighting](/go-tealeaves/components/syntax-highlighting/) | `Highlighter` | teahilite | Chroma-based syntax highlighting with configurable themes and language detection |

## Multi-Module Architecture

:::note
Tea Leaves does **not** have a root Go module. Each module listed above is independent with its own `go.mod`, version tags, and dependency tree. You never `go get` the repository root — you install individual modules by their path. See the [Module Reference](/go-tealeaves/reference/modules/) for the full module-to-component mapping.
:::

Components communicate through standard Bubble Tea conventions (`tea.Cmd` and `tea.Msg`), not through shared internal state, so they compose cleanly without tight coupling. See the [Architecture guide](/go-tealeaves/guides/architecture/) for a deeper look at the design philosophy.
