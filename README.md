# go-tealeaves

Reusable [Bubble Tea](https://github.com/charmbracelet/bubbletea) UI components for building terminal applications in Go.

## Overview

go-tealeaves is a multi-module repository of independently importable Bubble Tea components. Each package is a separate Go module — install only what you need.

## Packages

| Package | Import Path | Description |
|---------|------------|-------------|
| **teadrpdwn** | `github.com/mikeschinkel/go-tealeaves/teadrpdwn` | Dropdown selection with intelligent positioning |
| **teadepview** | `github.com/mikeschinkel/go-tealeaves/teadepview` | Interactive dependency path viewer |
| **teadiffr** | `github.com/mikeschinkel/go-tealeaves/teadiffr` | Diff rendering for terminal UIs |
| **teagrid** | `github.com/mikeschinkel/go-tealeaves/teagrid` | Data grid/table with sorting, filtering, pagination |
| **teamodal** | `github.com/mikeschinkel/go-tealeaves/teamodal` | Modal dialogs (Yes/No, OK, choice lists) |
| **teanotify** | `github.com/mikeschinkel/go-tealeaves/teanotify` | Toast-style notification overlays |
| **teastatus** | `github.com/mikeschinkel/go-tealeaves/teastatus` | Two-zone status bar (menu items + indicators) |
| **teatxtsnip** | `github.com/mikeschinkel/go-tealeaves/teatxtsnip` | Text area with selection and clipboard support |
| **teatree** | `github.com/mikeschinkel/go-tealeaves/teatree` | Generic tree view with customizable node rendering |
| **teautils** | `github.com/mikeschinkel/go-tealeaves/teautils` | Key registry, help visor, positioning utilities |

## Package Details

### teadrpdwn — Dropdown

Full-screen dropdown/popup with automatic positioning. Displays above or below the anchor point based on available terminal space. Supports scrolling, item truncation, and customizable styling.

- `DropdownModel` — main model
- `Option` — display text + value pair
- Sends `OptionSelectedMsg` or `DropdownCancelledMsg`
- Navigation: `Up`/`Down` to move, `Enter` to select, `Esc` to cancel

### teadepview — Dependency Path Viewer

Visualizes a single path through a dependency tree with interactive navigation. At each level, a dropdown shows alternative nodes. Useful for exploring module dependency graphs.

- `PathViewerModel` — main model
- `DependencyNode` — tree node with alternatives and children
- `ChildSelector` — strategy function for choosing the default child at each level
- Navigation: `Up`/`Down` to move levels, `Space`/`Right` to open alternatives

### teadiffr — Diff Renderer

Renders file diffs in the terminal with syntax-highlighted additions, deletions, and context lines. Supports both unified and condensed diff formats.

- `DiffRenderer` — interface for diff rendering strategies
- `TUIRenderer` — terminal UI renderer implementation
- `FileDiff`, `CondensedBlock`, `FileStatus` — diff data types

### teamodal — Modal Dialogs

Centered modal overlays for confirmations, alerts, and choice/list selection. Supports button focus with Tab, mouse clicks, and customizable key bindings. Uses the wither (immutable update) pattern.

- `ModalModel` — Yes/No or OK modal
- `ChoiceModel` — single-selection from a list
- `ListModel` — scrollable list with edit/delete actions
- Non-nil `tea.Cmd` return signals message consumption

### teastatus — Status Bar

Two-zone status bar: left side shows key-action menu items, right side shows text indicators. ANSI-aware width calculations ensure correct layout.

- `Model` — main model
- `MenuItem` — key binding + action label
- `StatusIndicator` — text with optional style
- Separator styles: pipe (`|`), space, or bracket

### teatxtsnip — Text Selection & Clipboard

Wraps Bubble Tea's `textarea.Model` with text selection (Shift+Arrow) and clipboard operations (Ctrl+C/X/V). Includes word, line, and document selection. Falls back to an internal clipboard when system clipboard is unavailable. Supports single-line mode.

- `Model` — wraps textarea with selection state
- `Selection` — start/end positions
- Clipboard: `Ctrl+C` copy, `Ctrl+X` cut, `Ctrl+V` paste

### teatree — Tree View

Generic tree view component parameterized on node data type. Supports expand/collapse, keyboard navigation, viewport scrolling, and pluggable node providers for custom rendering.

- `Model[T]` — generic Bubble Tea model
- `Tree[T]` — tree structure with focus tracking
- `Node[T]` — node with children and expansion state
- `NodeProvider[T]` — interface for custom rendering
- Navigation: `Up`/`Down` to move, `Right` to expand, `Left` to collapse

### teautils — Utilities

Shared utilities for Bubble Tea applications:

- **Key Registry** — centralized key binding management with namespace validation; separates definition (components) from presentation (app)
- **Help Visor** — renders a categorized help overlay from the key registry
- **Status Bar Renderer** — compact `[key] Action` format rendering
- **Positioning** — ANSI-aware overlay compositing functions

## Tools

### color-viewer

A terminal color reference tool that displays all 256 terminal colors as a scrollable grid, showing each background color paired with common foreground colors.

```bash
make build
./bin/color-viewer
```

**Controls:** Arrow keys to scroll, `q`/`Esc` to quit. Each cell shows `bg/fg` (e.g., `53/015` = background color 53, foreground color 15). Useful for choosing lipgloss color values when styling Bubble Tea components.

## Installation

Each package is an independent module. Install the ones you need:

```bash
go get github.com/mikeschinkel/go-tealeaves/teadrpdwn
go get github.com/mikeschinkel/go-tealeaves/teamodal
go get github.com/mikeschinkel/go-tealeaves/teatree
# etc.
```

There is no root module — do not `go get` the repository root.

## Examples

The `examples/` directory contains runnable programs demonstrating each package:

| Example | Package | Description |
|---------|---------|-------------|
| `teadrpdwn/simple` | teadrpdwn | Basic dropdown with fruit selection |
| `teadrpdwn/demo` | teadrpdwn | Comprehensive dropdown with configuration modes |
| `teadepview/treenav` | teadepview | Dependency tree path viewer |
| `teamodal/choices` | teamodal | Choice modal dialog |
| `teamodal/editlist` | teamodal | List modal with task management |
| `teamodal/various` | teamodal | Multiple modal dialog types |
| `teastatus/statusbar` | teastatus | Status bar with menu items and indicators |
| `teatxtsnip/editor` | teatxtsnip | Multi-line text editor with selection |
| `teatree/filetree` | teatree | File tree with navigation |
| `teautils/keyhelp` | teautils | Key registry with help modal overlay |

Build all examples:

```bash
make build-examples
```

Binaries are written to `./bin/examples/`. Run any example directly, e.g.:

```bash
./bin/examples/teadrpdwn-demo
./bin/examples/teatree-filetree
```

Or run one in place:

```bash
cd examples/teastatus/statusbar && go run .
```

## Architecture Notes

This is a multi-module Go repository. Each package is an independent Go module. Examples and tools are also independent modules that reference sibling packages via local `replace` directives.

### Key Patterns

- **Modal message consumption**: Non-nil `tea.Cmd` return from `Update()` signals that the component handled the message. Parent models check this to avoid processing already-consumed messages.
- **ANSI-aware layout**: All width calculations use `ansi.StringWidth()` from `charmbracelet/x/ansi`, never `len()`. Overlay compositing is ANSI-escape aware.
- **Wither/immutable pattern**: Components like `teamodal` use `With*()` methods that return updated copies rather than mutating state.

### Internal Dependencies

Most packages are standalone. Two have internal dependencies:

```
teadepview --> teadrpdwn   (uses dropdown for alternative selection)
teamodal   --> teautils    (uses positioning utilities)
```

### Cross-Cutting Documentation

- `adrs/` — Architecture Decision Records
- `docs/` — Best practices and reference guides

## Requirements

- Go 1.25+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) ecosystem (bubbletea, lipgloss, bubbles)

## Development

```bash
make help            # Show all targets
make build           # Build color-viewer to ./bin/
make test            # Run tests across all modules
make vet             # Run go vet across all modules
make fmt             # Format code with gofmt
make tidy            # Run go mod tidy across all modules and examples
make build-examples  # Build all example programs to ./bin/examples/
make clean           # Clean build artifacts
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE) — Copyright (c) 2026 Mike Schinkel
