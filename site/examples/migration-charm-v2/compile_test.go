// Source: site/src/content/docs/migration/charm-v2.mdx:37#3b06ae6f,47#f6c1b482,64#2ea7fc16,73#c00e6560,97#f4a9c7f8,116#7b6632c7,141#8779ae19,152#bd0a2210
package examples_test

import (
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

// Stub types mimicking old charmbracelet/bubbles table API for migration comparison
type tableColumn struct {
	Title string
	Width int
}
type tableRow []string
type tableStyles struct {
	Header   lipgloss.Style
	Selected lipgloss.Style
}
type tableModel struct{}
type tableOption func(*tableModel)

func tableWithColumns(cols []tableColumn) tableOption  { return func(*tableModel) {} }
func tableWithRows(rows []tableRow) tableOption        { return func(*tableModel) {} }
func tableWithStyles(s tableStyles) tableOption        { return func(*tableModel) {} }
func tableWithFocused(b bool) tableOption              { return func(*tableModel) {} }
func tableWithHeight(h int) tableOption                { return func(*tableModel) {} }
func tableDefaultStyles() tableStyles {
	return tableStyles{Header: lipgloss.NewStyle(), Selected: lipgloss.NewStyle()}
}
func tableNew(opts ...tableOption) tableModel { return tableModel{} }

// TestCompile_OldTableColumns verifies the old bubbles table column definition from charm-v2.mdx line 37.
func TestCompile_OldTableColumns(t *testing.T) {
	columns := []tableColumn{
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Age", Width: 5},
	}
	_ = columns
}

// TestCompile_NewTeagridColumns verifies the new teagrid column definition from charm-v2.mdx line 47.
func TestCompile_NewTeagridColumns(t *testing.T) {
	columns := []teagrid.Column{
		teagrid.NewColumn("name", "Name", 20),
		teagrid.NewColumn("status", "Status", 10),
		teagrid.NewColumn("age", "Age", 5),
	}
	_ = columns
}

// TestCompile_OldTableRows verifies the old bubbles table row definition from charm-v2.mdx line 64.
func TestCompile_OldTableRows(t *testing.T) {
	rows := []tableRow{
		{"Alice", "Active", "30"},
		{"Bob", "Inactive", "25"},
	}
	_ = rows
}

// TestCompile_NewTeagridRows verifies the new teagrid row definition from charm-v2.mdx line 73.
func TestCompile_NewTeagridRows(t *testing.T) {
	rows := []teagrid.Row{
		teagrid.NewRow(teagrid.RowData{
			"name":   "Alice",
			"status": "Active",
			"age":    "30",
		}),
		teagrid.NewRow(teagrid.RowData{
			"name":   "Bob",
			"status": "Inactive",
			"age":    "25",
		}),
	}
	_ = rows
}

// TestCompile_OldTableStyling verifies the old bubbles table styling from charm-v2.mdx line 97.
func TestCompile_OldTableStyling(t *testing.T) {
	columns := []tableColumn{
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Age", Width: 5},
	}
	rows := []tableRow{
		{"Alice", "Active", "30"},
		{"Bob", "Inactive", "25"},
	}
	s := tableDefaultStyles()
	s.Header = s.Header.
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))
	tbl := tableNew(
		tableWithColumns(columns),
		tableWithRows(rows),
		tableWithStyles(s),
	)
	_ = tbl
}

// TestCompile_NewTeagridStyling verifies the new teagrid fluent builder styling from charm-v2.mdx line 116.
func TestCompile_NewTeagridStyling(t *testing.T) {
	columns := []teagrid.Column{
		teagrid.NewColumn("name", "Name", 20),
		teagrid.NewColumn("status", "Status", 10),
		teagrid.NewColumn("age", "Age", 5),
	}
	rows := []teagrid.Row{
		teagrid.NewRow(teagrid.RowData{
			"name":   "Alice",
			"status": "Active",
			"age":    "30",
		}),
		teagrid.NewRow(teagrid.RowData{
			"name":   "Bob",
			"status": "Inactive",
			"age":    "25",
		}),
	}
	tbl := teagrid.NewGridModel(columns).
		WithRows(rows).
		WithBaseStyle(
			lipgloss.NewStyle().
				Align(lipgloss.Left).
				Foreground(lipgloss.Color("229")),
		).
		WithHeaderStyle(
			lipgloss.NewStyle().Bold(true),
		).
		WithHighlightStyle(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")),
		)
	_ = tbl
}

// TestCompile_OldTableConstruction verifies the old bubbles table construction from charm-v2.mdx line 141.
func TestCompile_OldTableConstruction(t *testing.T) {
	columns := []tableColumn{
		{Title: "Name", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Age", Width: 5},
	}
	rows := []tableRow{
		{"Alice", "Active", "30"},
		{"Bob", "Inactive", "25"},
	}
	tbl := tableNew(
		tableWithColumns(columns),
		tableWithRows(rows),
		tableWithFocused(true),
		tableWithHeight(10),
	)
	_ = tbl
}

// TestCompile_NewTeagridConstruction verifies the new teagrid construction from charm-v2.mdx line 152.
func TestCompile_NewTeagridConstruction(t *testing.T) {
	columns := []teagrid.Column{
		teagrid.NewColumn("name", "Name", 20),
		teagrid.NewColumn("status", "Status", 10),
		teagrid.NewColumn("age", "Age", 5),
	}
	rows := []teagrid.Row{
		teagrid.NewRow(teagrid.RowData{
			"name":   "Alice",
			"status": "Active",
			"age":    "30",
		}),
		teagrid.NewRow(teagrid.RowData{
			"name":   "Bob",
			"status": "Inactive",
			"age":    "25",
		}),
	}
	tbl := teagrid.NewGridModel(columns).
		WithRows(rows).
		WithPageSize(10).
		WithFocused(true)
	_ = tbl
}
