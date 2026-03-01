package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestSortByAsc(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"name": "Charlie"}),
		NewRow(RowData{"name": "Alice"}),
		NewRow(RowData{"name": "Bob"}),
	}
	m := New([]Column{NewColumn("name", "Name", 10)}).
		WithRows(rows).
		SortByAsc("name")

	visible := m.GetVisibleRows()
	assert.Equal(t, "Alice", visible[0].Data["name"])
	assert.Equal(t, "Bob", visible[1].Data["name"])
	assert.Equal(t, "Charlie", visible[2].Data["name"])
}

func TestSortByDesc(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"name": "Alice"}),
		NewRow(RowData{"name": "Charlie"}),
		NewRow(RowData{"name": "Bob"}),
	}
	m := New([]Column{NewColumn("name", "Name", 10)}).
		WithRows(rows).
		SortByDesc("name")

	visible := m.GetVisibleRows()
	assert.Equal(t, "Charlie", visible[0].Data["name"])
	assert.Equal(t, "Bob", visible[1].Data["name"])
	assert.Equal(t, "Alice", visible[2].Data["name"])
}

func TestSortNumeric(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"score": 100}),
		NewRow(RowData{"score": 3}),
		NewRow(RowData{"score": 42}),
	}
	m := New([]Column{NewColumn("score", "Score", 10)}).
		WithRows(rows).
		SortByAsc("score")

	visible := m.GetVisibleRows()
	assert.Equal(t, 3, visible[0].Data["score"])
	assert.Equal(t, 42, visible[1].Data["score"])
	assert.Equal(t, 100, visible[2].Data["score"])
}

func TestSortWithSortValue(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"date": NewCellValueWithSortKey("Jan 1", 1, lipgloss.NewStyle())}),
		NewRow(RowData{"date": NewCellValueWithSortKey("Mar 15", 3, lipgloss.NewStyle())}),
		NewRow(RowData{"date": NewCellValueWithSortKey("Feb 14", 2, lipgloss.NewStyle())}),
	}
	m := New([]Column{NewColumn("date", "Date", 10)}).
		WithRows(rows).
		SortByAsc("date")

	visible := m.GetVisibleRows()
	cv0 := visible[0].Data["date"].(CellValue)
	cv1 := visible[1].Data["date"].(CellValue)
	cv2 := visible[2].Data["date"].(CellValue)
	assert.Equal(t, "Jan 1", cv0.Data)
	assert.Equal(t, "Feb 14", cv1.Data)
	assert.Equal(t, "Mar 15", cv2.Data)
}

func TestThenSortBy(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"dept": "B", "name": "Charlie"}),
		NewRow(RowData{"dept": "A", "name": "Bob"}),
		NewRow(RowData{"dept": "A", "name": "Alice"}),
	}
	m := New([]Column{
		NewColumn("dept", "Dept", 10),
		NewColumn("name", "Name", 10),
	}).WithRows(rows).
		SortByAsc("dept").
		ThenSortByAsc("name")

	visible := m.GetVisibleRows()
	assert.Equal(t, "A", visible[0].Data["dept"])
	assert.Equal(t, "Alice", visible[0].Data["name"])
	assert.Equal(t, "A", visible[1].Data["dept"])
	assert.Equal(t, "Bob", visible[1].Data["name"])
}
