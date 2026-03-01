# teagrid

A data grid component for [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Charm v2). Forked from [bubble-table](https://github.com/evertras/bubble-table) by Brandon Fulljames -- see the [ADR](adrs/adr-2025-03-01-fork-over-upstream-prs.md) for why a fork was chosen over upstream PRs.

## Installation

```bash
go get github.com/mikeschinkel/go-tealeaves/teagrid
```

## Quick Start

```go
columns := []teagrid.Column{
    teagrid.NewColumn("name", "Name", 15),
    teagrid.NewColumn("email", "Email", 25),
}

rows := []teagrid.Row{
    teagrid.NewRow(teagrid.RowData{"name": "Alice", "email": "alice@example.com"}),
    teagrid.NewRow(teagrid.RowData{"name": "Bob", "email": "bob@example.com"}),
}

table := teagrid.New(columns).WithRows(rows).Focused(true)
```

## Key Features

- **Left-aligned text by default** with per-column padding (`paddingLeft=1`)
- **Region-based borders** (Outer, Header, Inner, Footer) with presets: `BorderRounded()`, `BorderDefault()`, `Borderless()`, `BorderMinimal()`
- **Render-time cursor/highlight** -- no O(n) row rebuilding
- **CellValue** with separate `SortValue`, `StyleFunc`, and rich text `Spans`
- **Cell cursor mode** with per-cell highlighting
- **Independent footer** (never inherits baseStyle) with filter + pagination zones
- **SetSize(w, h)** auto fill/scroll (replaces manual `WithTargetWidth`/`WithMaxTotalWidth`)
- **Sorting, filtering** (contains + fuzzy), **pagination**
- **Horizontal scrolling** with frozen columns
- **Selectable rows** (no auto-added checkbox column)

## Charm v2

v0.2.0 targets Charm v2 exclusively:

- `charm.land/bubbletea/v2`
- `charm.land/lipgloss/v2`
- `charm.land/bubbles/v2`
- `View()` returns `tea.View`
- `Update()` returns `(Model, tea.Cmd)` (concrete type, standard bubbles pattern)
- Key handling uses `tea.KeyPressMsg`

## Migration from v0.1.0

| v0.1.0 | v0.2.0 |
|--------|--------|
| `github.com/charmbracelet/*` imports | `charm.land/*/v2` |
| `tea.KeyMsg` | `tea.KeyPressMsg` |
| `View() string` | `View() tea.View` (use `.Content` to extract string) |
| `WithTargetWidth(w)` / `WithMaxTotalWidth(w)` | `SetSize(w, h)` |
| `WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Left))` | Not needed (left-align is the default) |
| Manual `" "` padding | Not needed (`paddingLeft=1` default) |
| `StyledCell` | `CellValue` (type alias provided for compatibility) |
| `tea.WithAltScreen()` | Set `v.AltScreen = true` on `tea.View` |

## Migration from bubble-table

If you are migrating directly from `github.com/evertras/bubble-table/table`:

- Change imports from `github.com/evertras/bubble-table/table` to `github.com/mikeschinkel/go-tealeaves/teagrid`
- `table.New(columns)` becomes `teagrid.New(columns)`
- `StyledCell` becomes `CellValue`
- All of the v0.1.0 migration changes listed above also apply
- See the [ADR](adrs/adr-2025-03-01-fork-over-upstream-prs.md) for detailed rationale

## Examples

| Example | Description |
|---------|-------------|
| [`examples/teagrid/simplest/`](../examples/teagrid/simplest/) | Minimal non-interactive table |
| [`examples/teagrid/sorting/`](../examples/teagrid/sorting/) | Sorting by different columns |
| [`examples/teagrid/filtering/`](../examples/teagrid/filtering/) | Filter with `/` key, flex columns, pagination |
| [`examples/teagrid/scrolling/`](../examples/teagrid/scrolling/) | Horizontal scrolling with frozen columns |

## Attribution

This package is a fork of [`github.com/evertras/bubble-table/table`](https://github.com/evertras/bubble-table) (MIT License, Copyright (c) 2022 Brandon Fulljames). See [ADR: Fork from bubble-table](adrs/adr-2025-03-01-fork-over-upstream-prs.md) for detailed rationale.

## License

MIT -- see [LICENSE](LICENSE) for details.
