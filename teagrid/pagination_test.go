package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := New([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	assert.Equal(t, 1, m.CurrentPage())
	assert.Equal(t, 3, m.MaxPages())
	assert.Equal(t, 25, m.TotalRows())
	assert.Equal(t, 10, m.PageSize())
}

func TestVisibleIndices(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := New([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	start, end := m.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 9, end)
}

func TestPageDown(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := New([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		PageDown()

	assert.Equal(t, 2, m.CurrentPage())

	start, end := m.VisibleIndices()
	assert.Equal(t, 10, start)
	assert.Equal(t, 19, end)
}

func TestPageUpWrapping(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := New([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithPaginationWrapping(true).
		PageUp()

	assert.Equal(t, 3, m.CurrentPage(), "should wrap to last page")
}

func TestPageUpNoWrapping(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := New([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithPaginationWrapping(false).
		PageUp()

	assert.Equal(t, 1, m.CurrentPage(), "should stay on first page")
}

func TestPageFirstLast(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := New([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	m = m.PageLast()
	assert.Equal(t, 3, m.CurrentPage())

	m = m.PageFirst()
	assert.Equal(t, 1, m.CurrentPage())
}

func TestNoPagination(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := New([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows)

	assert.Equal(t, 0, m.PageSize())
	assert.Equal(t, 1, m.MaxPages())

	start, end := m.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 24, end)
}
