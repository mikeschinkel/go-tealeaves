/*
Package teagrid provides a Bubble Tea v2 component for interactive,
customizable data grids.

teagrid supports fixed and flex-width columns with built-in padding and
alignment, per-cell/per-row/per-column styling with a style cascade,
render-time cursor and row highlighting (no row rebuilding), region-based
border configuration with presets, horizontal scrolling with frozen columns,
sorting with separate sort keys, filtering with match highlighting, rich text
cells with Span-based inline styling, and cell cursor navigation with selection
events.

Basic usage:

	columns := []teagrid.Column{
		teagrid.NewColumn("name", "Name", 20),
		teagrid.NewColumn("count", "Count", 10),
	}

	rows := []teagrid.Row{
		teagrid.NewRow(teagrid.RowData{
			"name":  "Cheeseburger",
			"count": 3,
		}),
	}

	grid := teagrid.NewGridModel(columns).WithRows(rows)

	// Use it like any Bubble Tea v2 component in your view
	grid.View()

# Stability

This package is provisional as of v0.3.0. The public API may change in
minor releases until promoted to stable.
*/
package teagrid
