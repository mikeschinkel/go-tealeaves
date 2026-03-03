---
title: "Example Gallery"
description: "All Tea Leaves example programs — ready to run with go run."
---

Every Tea Leaves component ships with runnable example programs. Clone the repository and run any example directly to see the component in action before reading the source.

```bash
git clone https://github.com/mikeschinkel/go-tealeaves.git
cd go-tealeaves
```

All `go run` commands below assume you are in the repository root directory. Each example is a self-contained `main.go` with its own `go.mod`, so you can also `cd` into the example directory and run `go run .` directly.

## Examples by component

| Component | Example | Description | Run Command |
|-----------|---------|-------------|-------------|
| teadrpdwn | simple | Dropdown with fruit selection | `go run ./examples/teadrpdwn/simple` |
| teadrpdwn | demo | Full-featured dropdown demo | `go run ./examples/teadrpdwn/demo` |
| teagrid | simplest | Minimal data grid | `go run ./examples/teagrid/simplest` |
| teagrid | sorting | Grid with column sorting | `go run ./examples/teagrid/sorting` |
| teagrid | filtering | Grid with real-time filtering | `go run ./examples/teagrid/filtering` |
| teagrid | scrolling | Grid with pagination | `go run ./examples/teagrid/scrolling` |
| teamodal | various | Multiple modal types | `go run ./examples/teamodal/various` |
| teamodal | choices | Choice selection dialog | `go run ./examples/teamodal/choices` |
| teamodal | editlist | Editable list dialog | `go run ./examples/teamodal/editlist` |
| teamodal | vertical | Vertical button layout | `go run ./examples/teamodal/vertical` |
| teanotify | simple | Toast notifications | `go run ./examples/teanotify/simple` |
| teastatus | statusbar | Status bar with separators | `go run ./examples/teastatus/statusbar` |
| teatxtsnip | editor | Text editor with selection | `go run ./examples/teatxtsnip/editor` |
| teatree | filetree | File tree navigator | `go run ./examples/teatree/filetree` |
| teadepview | treenav | Dependency tree navigator | `go run ./examples/teadepview/treenav` |
| teautils | keyhelp | Key registry with help modal | `go run ./examples/teautils/keyhelp` |

:::note
Some examples use `tea.WithAltScreen()` which takes over the full terminal. Press `q` or `Ctrl+C` to exit back to your shell.
:::
