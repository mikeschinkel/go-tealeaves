package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisibleRows(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 3}),
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		SortByAsc("x")

	visible := m.VisibleRows()
	assert.Len(t, visible, 3)
	assert.Equal(t, 1, visible[0].Data["x"])
}

func TestVisibleRowsCache(t *testing.T) {
	rows := []Row{NewRow(RowData{"x": 1})}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithRows(rows)

	// First call populates cache
	v1 := m.VisibleRows()
	// Second call uses cache
	v2 := m.VisibleRows()

	assert.Equal(t, len(v1), len(v2))
}

func TestHighlightedRow(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": "a"}),
		NewRow(RowData{"x": "b"}),
	}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		WithHighlightedRow(1)

	row := m.HighlightedRow()
	assert.Equal(t, "b", row.Data["x"])
}

func TestHighlightedRowEmpty(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)})
	row := m.HighlightedRow()
	assert.Nil(t, row.Data)
}

func TestSelectedRows(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}).Selected(true),
		NewRow(RowData{"x": 2}),
		NewRow(RowData{"x": 3}).Selected(true),
	}
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithRows(rows)

	selected := m.SelectedRows()
	assert.Len(t, selected, 2)
}

func TestNaturalWidth(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 20),
	}
	m := NewGridModel(cols)

	assert.Equal(t, 37, m.NaturalWidth())
}

func TestTotalWidth(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewColumn("b", "B", 20),
	}
	m := NewGridModel(cols)

	assert.Equal(t, 37, m.TotalWidth())
}

func TestTotalWidthWithFlex(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 10),
		NewFlexColumn("b", "B", 1),
	}
	m := NewGridModel(cols).WithSize(50, 24)

	// TotalWidth should match viewport after flex resolution
	assert.Equal(t, 50, m.TotalWidth())
}

func TestHasFooter(t *testing.T) {
	t.Run("no footer when hidden", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithFooterVisibility(false)
		assert.False(t, m.hasFooter())
	})

	t.Run("footer with pagination", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithPageSize(10)
		assert.True(t, m.hasFooter())
	})

	t.Run("footer with static text", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithStaticFooter("info")
		assert.True(t, m.hasFooter())
	})

	t.Run("footer with filter", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithFiltered(true)
		assert.True(t, m.hasFooter())
	})
}

func TestColumnSorting(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		SortByAsc("x").
		ThenSortByDesc("y")

	sorting := m.ColumnSorting()
	assert.Len(t, sorting, 2)

	// Mutation of returned slice should not affect model
	sorting[0].ColumnKey = "mutated"
	assert.NotEqual(t, "mutated", m.ColumnSorting()[0].ColumnKey)
}

func TestLastUpdateUserEvents(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)})

	assert.Nil(t, m.LastUpdateUserEvents())

	m = m.appendUserEvent(UserEventFilterInputFocused{})
	events := m.LastUpdateUserEvents()
	assert.Len(t, events, 1)

	m = m.clearUserEvents()
	assert.Nil(t, m.LastUpdateUserEvents())
}
