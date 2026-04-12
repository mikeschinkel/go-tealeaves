// Source: site/src/content/docs/components/grid-view.mdx:21,175,201,212,226,251,265
package examples_test

import (
	"testing"
	"time"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	columnKeyName    = "name"
	columnKeyElement = "element"
	columnKeyStatus  = "status"
)

// TestCompile_GridViewQuickExample verifies the quick example from grid-view.mdx.
func TestCompile_GridViewQuickExample(t *testing.T) {
	table := teagrid.NewGridModel([]teagrid.Column{
		teagrid.NewColumn(columnKeyName, "Name", 13),
		teagrid.NewColumn(columnKeyElement, "Element", 10),
	}).WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Left)).WithRows([]teagrid.Row{
		teagrid.NewRow(teagrid.RowData{columnKeyName: "Pikachu", columnKeyElement: "Electric"}),
		teagrid.NewRow(teagrid.RowData{columnKeyName: "Charmander", columnKeyElement: "Fire"}),
	})
	_ = table.View()
}

// TestCompile_CellStyleFunc verifies CellStyleFunc usage from grid-view.mdx.
func TestCompile_CellStyleFunc(t *testing.T) {
	status := teagrid.NewCellValueWithStyleFunc("ERROR", func(in teagrid.CellStyleInput) lipgloss.Style {
		if in.IsHighlightedRow {
			return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	})
	_ = status
}

// TestCompile_CellWithSpans verifies NewCellValueWithSpans usage from grid-view.mdx.
func TestCompile_CellWithSpans(t *testing.T) {
	cell := teagrid.NewCellValueWithSpans([]teagrid.Span{
		teagrid.NewSpan("✓ ", lipgloss.NewStyle().Foreground(lipgloss.Color("82"))),
		teagrid.NewSpan("Complete", lipgloss.NewStyle().Foreground(lipgloss.Color("246"))),
	}, lipgloss.Style{})
	_ = cell
}

// TestCompile_HorizontalScrolling verifies horizontal scrolling from grid-view.mdx.
func TestCompile_HorizontalScrolling(t *testing.T) {
	grid := teagrid.NewGridModel([]teagrid.Column{
		teagrid.NewColumn("col1", "Col1", 20),
		teagrid.NewColumn("col2", "Col2", 20),
	})

	grid = grid.WithHorizontalFreezeColumnCount(1)
	grid = grid.ScrollLeft()
	grid = grid.ScrollRight()
	_ = grid
}

// TestCompile_RowStyleFunc verifies per-row dynamic styling from grid-view.mdx.
func TestCompile_RowStyleFunc(t *testing.T) {
	grid := teagrid.NewGridModel([]teagrid.Column{
		teagrid.NewColumn("col", "Col", 10),
	})

	grid = grid.WithRowStyleFunc(func(in teagrid.RowStyleFuncInput) lipgloss.Style {
		if in.IsHighlighted {
			return lipgloss.NewStyle().Background(lipgloss.Color("235"))
		}
		return lipgloss.Style{}
	})
	_ = grid
}

// TestCompile_GridMetadata verifies WithMetadata from grid-view.mdx.
func TestCompile_GridMetadata(t *testing.T) {
	grid := teagrid.NewGridModel([]teagrid.Column{
		teagrid.NewColumn("col", "Col", 10),
	})

	selectedSet := map[string]bool{"item1": true}
	grid = grid.WithMetadata(map[string]any{
		"selectedIDs": selectedSet,
	})
	_ = grid
}

// TestCompile_CellWithSortKey verifies NewCellValueWithSortKey from grid-view.mdx.
func TestCompile_CellWithSortKey(t *testing.T) {
	cell := teagrid.NewCellValueWithSortKey(
		"Jan 1",
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		lipgloss.Style{},
	)
	_ = cell
}
