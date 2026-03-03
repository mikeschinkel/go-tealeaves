---
title: Components
description: Overview of the 10 independent Go modules that make up Tea Leaves — production-ready Bubble Tea components for dropdowns, grids, modals, notifications, trees, text snippets, status bars, dependency viewing, diff rendering, and utilities.
---

Tea Leaves ships 10 components, each as its own Go module with a dedicated `go.mod`. You install only what you need, and each component brings only its own dependencies into your project.

## Component Reference

| Component | Description | Install Command |
|-----------|-------------|-----------------|
| [teadrpdwn](/go-tealeaves/components/teadrpdwn/) | Dropdown selection with intelligent above/below positioning | `go get github.com/mikeschinkel/go-tealeaves/teadrpdwn` |
| [teagrid](/go-tealeaves/components/teagrid/) | Data grid with sorting, filtering, pagination, row selection | `go get github.com/mikeschinkel/go-tealeaves/teagrid` |
| [teamodal](/go-tealeaves/components/teamodal/) | Modal dialogs (Yes/No, OK, Choice, List, Progress) | `go get github.com/mikeschinkel/go-tealeaves/teamodal` |
| [teanotify](/go-tealeaves/components/teanotify/) | Toast notifications with auto-dismiss and color fade | `go get github.com/mikeschinkel/go-tealeaves/teanotify` |
| [teatree](/go-tealeaves/components/teatree/) | Generic tree view with expand/collapse and pluggable providers | `go get github.com/mikeschinkel/go-tealeaves/teatree` |
| [teatxtsnip](/go-tealeaves/components/teatxtsnip/) | Text area with Shift+Arrow selection and clipboard | `go get github.com/mikeschinkel/go-tealeaves/teatxtsnip` |
| [teastatus](/go-tealeaves/components/teastatus/) | Two-zone status bar with menus and indicators | `go get github.com/mikeschinkel/go-tealeaves/teastatus` |
| [teadepview](/go-tealeaves/components/teadepview/) | Dependency path viewer for module graphs | `go get github.com/mikeschinkel/go-tealeaves/teadepview` |
| [teadiffr](/go-tealeaves/components/teadiffr/) | Condensed diff rendering for terminal UIs | `go get github.com/mikeschinkel/go-tealeaves/teadiffr` |
| [teautils](/go-tealeaves/components/teautils/) | Key registry, help visor, positioning, theming | `go get github.com/mikeschinkel/go-tealeaves/teautils` |

## Multi-Module Architecture

:::note
Tea Leaves does **not** have a root Go module. Each component listed above is an independent module with its own `go.mod`, its own version tags, and its own dependency tree. You never `go get` the repository root — you install individual components by their module path.
:::

This multi-module design means adding `teadrpdwn` to your project will never pull in `teagrid`'s dependencies, and upgrading `teamodal` cannot break your `teanotify` import. Components communicate through standard Bubble Tea conventions (`tea.Cmd` and `tea.Msg`), not through shared internal state, so they compose cleanly without tight coupling. See the [Architecture guide](/go-tealeaves/guides/architecture/) for a deeper look at the design philosophy.
