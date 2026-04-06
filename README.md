# go-tealeaves

Reusable [Bubble Tea](https://github.com/charmbracelet/bubbletea) UI components for building terminal applications in Go.

## Overview

go-tealeaves is a multi-module repository of independently importable Bubble Tea v2 components. Each package is a separate Go module — install only what you need.

## Packages

| Package | Import Path | Description |
|---------|------------|-------------|
| **teacrumbs** | `github.com/mikeschinkel/go-tealeaves/teacrumbs` | Breadcrumb trail with hover/click support |
| **teadiff** | `github.com/mikeschinkel/go-tealeaves/teadiff` | Interactive side-by-side split diff viewer |
| **teafields** | `github.com/mikeschinkel/go-tealeaves/teafields` | Dropdown selection with intelligent positioning |
| **teagrid** | `github.com/mikeschinkel/go-tealeaves/teagrid` | Data grid with sorting, filtering, pagination, panning |
| **teaguide** | `github.com/mikeschinkel/go-tealeaves/teaguide` | Step-by-step guide/wizard component |
| **teahelp** | `github.com/mikeschinkel/go-tealeaves/teahelp` | Help visor overlay with paginated content |
| **teahilite** | `github.com/mikeschinkel/go-tealeaves/teahilite` | Chroma-based syntax highlighting |
| **tealayout** | `github.com/mikeschinkel/go-tealeaves/tealayout` | Constraint-based layout engine (rows/columns) |
| **teamodal** | `github.com/mikeschinkel/go-tealeaves/teamodal` | Modal dialogs (confirm, choice, list, multiselect) |
| **teanotify** | `github.com/mikeschinkel/go-tealeaves/teanotify` | Toast-style notification overlays |
| **teastatus** | `github.com/mikeschinkel/go-tealeaves/teastatus` | Two-zone status bar (menu items + indicators) |
| **teatree** | `github.com/mikeschinkel/go-tealeaves/teatree` | Generic tree view + drilldown path viewer |
| **teatext** | `github.com/mikeschinkel/go-tealeaves/teatext` | Text area with selection and clipboard support |
| **teautils** | `github.com/mikeschinkel/go-tealeaves/teautils` | Theming, key registry, positioning utilities |

## Installation

Each package is an independent module. Install the ones you need:

```bash
go get github.com/mikeschinkel/go-tealeaves/teagrid
go get github.com/mikeschinkel/go-tealeaves/teamodal
go get github.com/mikeschinkel/go-tealeaves/teatree
# etc.
```

There is no root module — do not `go get` the repository root.

## Examples

Each module contains its own examples under `<module>/examples/`:

| Example | Module | Description |
|---------|--------|-------------|
| `teafields/examples/simple` | teafields | Basic dropdown with fruit selection |
| `teafields/examples/demo` | teafields | Comprehensive dropdown with configuration modes |
| `teagrid/examples/simplest` | teagrid | Minimal grid setup |
| `teagrid/examples/filtering` | teagrid | Grid with filter input |
| `teagrid/examples/sorting` | teagrid | Grid with column sorting |
| `teagrid/examples/scrolling` | teagrid | Grid with horizontal scrolling |
| `teagrid/examples/panning` | teagrid | Grid with horizontal panning |
| `tealayout/examples/multipane` | tealayout | Multi-pane layout with nested groups, resizing, and visibility |
| `teamodal/examples/choices` | teamodal | Choice modal dialog |
| `teamodal/examples/editlist` | teamodal | List modal with task management |
| `teamodal/examples/multiselect` | teamodal | Multi-select modal with scrollbar |
| `teamodal/examples/various` | teamodal | Multiple modal dialog types |
| `teamodal/examples/vertical` | teamodal | Vertical modal layout |
| `teanotify/examples/simple` | teanotify | Toast notification demo |
| `teastatus/examples/statusbar` | teastatus | Status bar with menu items and indicators |
| `teadiff/examples/splitdiff` | teadiff | Side-by-side diff viewer |
| `teatree/examples/filetree` | teatree | File tree with navigation |
| `teatree/examples/drilldown` | teatree | Drilldown path viewer |
| `teatext/examples/editor` | teatext | Multi-line text editor with selection |
| `teautils/examples/keyhelp` | teautils | Key registry with help overlay |
| `teautils/examples/theming` | teautils | Theme switching demo |
| `teaguide/example` | teaguide | Step-by-step wizard |

Run any example:

```bash
cd teagrid/examples/simplest && go run .
```

## Architecture Notes

This is a multi-module Go repository. Each package is a separate Go module with its own `go.mod`. Examples are also independent modules that reference their parent package via local `replace` directives.

### Key Patterns

- **ClearPath style**: Named returns, `goto end` pattern, no `else` chains.
- **doterr errors**: Structured error handling with sentinel errors and metadata.
- **Modal message consumption**: Non-nil `tea.Cmd` return from `Update()` signals that the component handled the message.
- **ANSI-aware layout**: All width calculations use `ansi.StringWidth()`, never `len()`.
- **Wither pattern**: Components use `With*()` methods that return updated copies.

### Internal Dependencies

Most packages are standalone. Key dependencies:

```
teatree   --> teafields  (uses dropdown for drilldown alternatives)
teamodal  --> teautils   (uses positioning utilities)
teacrumbs --> teautils
teafields --> teautils
teagrid   --> teautils
teaguide  --> teautils
teahelp   --> teautils
teadiff -> teautils
teanotify --> teautils
teastatus --> teautils
teatext -> teautils
```

### Cross-Cutting Documentation

- `docs/adrs/` — Architecture Decision Records
- `docs/` — Best practices and research

## Requirements

- Go 1.25+
- [Bubble Tea v2](https://charm.land/bubbletea/v2) ecosystem (bubbletea, lipgloss, bubbles)

## Development

```bash
just test            # Run tests across all modules
just tidy            # Run go mod tidy across all modules
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT](LICENSE) — Copyright (c) 2026 Mike Schinkel
