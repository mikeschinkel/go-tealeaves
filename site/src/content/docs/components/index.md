---
title: Components
description: Overview of the 16 Tea Leaves components — production-ready Bubble Tea building blocks for dialogs, views, app chrome, and utilities.
---

Tea Leaves ships 16 components across 10 Go modules. Each module has its own `go.mod` — you install only what you need.

## Dialogs

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Confirmation Modal](/go-tealeaves/components/confirmation-modal/) | `ModalModel` | teamodal | Yes/No and OK confirmation dialogs |
| [Choice Modal](/go-tealeaves/components/choice-modal/) | `ChoiceModel` | teamodal | Multi-option selection with hotkeys |
| [List Modal](/go-tealeaves/components/list-modal/) | `ListModel[T]` | teamodal | Editable list with inline editing and CRUD |
| [Progress Modal](/go-tealeaves/components/progress-modal/) | `ProgressModal` | teamodal | Progress indicator with cancel/background |
| [Dropdown](/go-tealeaves/components/dropdown/) | `DropdownModel` | teadrpdwn | Dropdown selection with smart positioning |

## Views

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Data Grid](/go-tealeaves/components/data-grid/) | `Model` | teagrid | Sorting, filtering, pagination, row selection |
| [Tree View](/go-tealeaves/components/tree-view/) | `Model[T]` | teatree | Expand/collapse, pluggable node providers |
| [Diff Renderer](/go-tealeaves/components/diff-renderer/) | `TUIRenderer` | teadiffr | Condensed diff rendering with lipgloss |
| [Dep Viewer](/go-tealeaves/components/dep-viewer/) | `PathViewerModel` | teadepview | Interactive dependency path viewer |

## App Chrome

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Status Bar](/go-tealeaves/components/status-bar/) | `Model` | teastatus | Two-zone status bar with menus and indicators |
| [Notifications](/go-tealeaves/components/notifications/) | `NotifyModel` | teanotify | Toast notifications with auto-dismiss and color fade |
| [Help Visor](/go-tealeaves/components/help-visor/) | `HelpVisorStyle` | teautils | Help overlay styling for categorized key bindings |

## Utilities

| Component | Go Type | Module | Description |
|-----------|---------|--------|-------------|
| [Key Registry](/go-tealeaves/components/key-registry/) | `KeyRegistry` | teautils | Centralized key binding management |
| [Theming](/go-tealeaves/components/theming/) | `Theme` / `Palette` | teautils | Consistent colors across all components |
| [Positioning](/go-tealeaves/components/positioning/) | (functions) | teautils | ANSI-aware centering and measurement |
| [Text Selection](/go-tealeaves/components/text-selection/) | `Model` | teatxtsnip | Textarea with Shift+Arrow selection and clipboard |

## Multi-Module Architecture

:::note
Tea Leaves does **not** have a root Go module. Each module listed above is independent with its own `go.mod`, version tags, and dependency tree. You never `go get` the repository root — you install individual modules by their path. See the [Module Reference](/go-tealeaves/reference/modules/) for the full module-to-component mapping.
:::

Components communicate through standard Bubble Tea conventions (`tea.Cmd` and `tea.Msg`), not through shared internal state, so they compose cleanly without tight coupling. See the [Architecture guide](/go-tealeaves/guides/architecture/) for a deeper look at the design philosophy.
