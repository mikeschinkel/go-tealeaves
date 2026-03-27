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

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
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

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
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

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// PgDn #1: Phase 1 — cursor jumps to bottom of viewport, no scroll
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex, "phase 1: cursor at bottom of viewport")
	assert.Equal(t, 0, m.scrollOffset, "phase 1: no scroll")
	start, end := m.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 9, end)
	assert.Equal(t, 1, m.CurrentPage())

	// PgDn #2: Phase 2 — scroll viewport forward, cursor at bottom of new viewport
	m = m.PageDown()
	assert.Equal(t, 19, m.rowCursorIndex, "phase 2: cursor at bottom of new viewport")
	assert.Equal(t, 10, m.scrollOffset, "phase 2: scroll offset advanced by pageSize")
	start, end = m.VisibleIndices()
	assert.Equal(t, 10, start)
	assert.Equal(t, 19, end)
	assert.Equal(t, 2, m.CurrentPage())
}

func TestPageUpWrapping(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	// Cursor at row 0, scroll at 0 — Phase 2 wrap
	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithPaginationWrapping(true).
		PageUp()

	assert.Equal(t, 24, m.rowCursorIndex, "cursor wraps to last row")
	assert.Equal(t, 15, m.scrollOffset, "scrollOffset positions viewport at end")
	assert.Equal(t, 3, m.CurrentPage(), "should wrap to last page")
}

func TestPageUpNoWrapping(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	// Cursor at row 0, scroll at 0 — Phase 2 clamp (no change)
	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithPaginationWrapping(false).
		PageUp()

	assert.Equal(t, 0, m.rowCursorIndex, "cursor clamped to 0")
	assert.Equal(t, 0, m.scrollOffset, "scroll stays at 0")
	assert.Equal(t, 1, m.CurrentPage(), "should stay on first page")
}

func TestPageFirstLast(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	m = m.PageLast()
	assert.Equal(t, 24, m.rowCursorIndex, "cursor at last row")
	assert.Equal(t, 15, m.scrollOffset, "scrollOffset positions viewport at end")
	assert.Equal(t, 3, m.CurrentPage())

	m = m.PageFirst()
	assert.Equal(t, 0, m.rowCursorIndex, "cursor at first row")
	assert.Equal(t, 0, m.scrollOffset, "scrollOffset reset to 0")
	assert.Equal(t, 1, m.CurrentPage())
}

func TestNoPagination(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows)

	assert.Equal(t, 0, m.PageSize())
	assert.Equal(t, 1, m.MaxPages())

	start, end := m.VisibleIndices()
	assert.Equal(t, 0, start)
	assert.Equal(t, 24, end)
}

func TestPageDownFullCycle(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// PgDn #1: Phase 1 — cursor to bottom of viewport
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex)
	assert.Equal(t, 0, m.scrollOffset)

	// PgDn #2: Phase 2 — scroll to rows 10-19
	m = m.PageDown()
	assert.Equal(t, 19, m.rowCursorIndex)
	assert.Equal(t, 10, m.scrollOffset)

	// PgDn #3: Phase 2 — scroll to rows 15-24 (clamped)
	m = m.PageDown()
	assert.Equal(t, 24, m.rowCursorIndex)
	assert.Equal(t, 15, m.scrollOffset)

	// PgDn #4: Phase 2 — already at end, no-op
	m = m.PageDown()
	assert.Equal(t, 24, m.rowCursorIndex, "no-op at end")
	assert.Equal(t, 15, m.scrollOffset, "no-op at end")
}

func TestPageUpFullCycle(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	// Start at the end
	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		PageLast()

	assert.Equal(t, 24, m.rowCursorIndex)
	assert.Equal(t, 15, m.scrollOffset)

	// PgUp #1: Phase 1 — cursor to top of viewport
	m = m.PageUp()
	assert.Equal(t, 15, m.rowCursorIndex)
	assert.Equal(t, 15, m.scrollOffset)

	// PgUp #2: Phase 2 — scroll to rows 5-14
	m = m.PageUp()
	assert.Equal(t, 5, m.rowCursorIndex)
	assert.Equal(t, 5, m.scrollOffset)

	// PgUp #3: Phase 2 — scroll to rows 0-9 (clamped)
	m = m.PageUp()
	assert.Equal(t, 0, m.rowCursorIndex)
	assert.Equal(t, 0, m.scrollOffset)

	// PgUp #4: Phase 2 — already at start, no-op
	m = m.PageUp()
	assert.Equal(t, 0, m.rowCursorIndex, "no-op at start")
	assert.Equal(t, 0, m.scrollOffset, "no-op at start")
}

func TestPageDownSinglePage(t *testing.T) {
	rows := make([]Row, 5)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// PgDn #1: Phase 1 — cursor to last row (only 5 rows)
	m = m.PageDown()
	assert.Equal(t, 4, m.rowCursorIndex)
	assert.Equal(t, 0, m.scrollOffset)

	// PgDn #2: Phase 2 — already at end, no-op
	m = m.PageDown()
	assert.Equal(t, 4, m.rowCursorIndex, "no-op: single page")
	assert.Equal(t, 0, m.scrollOffset, "no-op: single page")
}

func TestPageUpSinglePage(t *testing.T) {
	rows := make([]Row, 5)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// PgUp #1: Phase 1 — cursor already at 0, so Phase 2
	// Already at start, no-op
	m = m.PageUp()
	assert.Equal(t, 0, m.rowCursorIndex, "no-op at start")
	assert.Equal(t, 0, m.scrollOffset, "no-op at start")
}

func TestPageDownWrappingFullCycle(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithPaginationWrapping(true)

	// Phase 1
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex)

	// Phase 2 — scroll to rows 10-19
	m = m.PageDown()
	assert.Equal(t, 19, m.rowCursorIndex)
	assert.Equal(t, 10, m.scrollOffset)

	// Phase 2 — scroll to rows 15-24 (clamped)
	m = m.PageDown()
	assert.Equal(t, 24, m.rowCursorIndex)
	assert.Equal(t, 15, m.scrollOffset)

	// Phase 2 — at end, wraps to beginning
	m = m.PageDown()
	assert.Equal(t, 0, m.rowCursorIndex, "wraps to first row")
	assert.Equal(t, 0, m.scrollOffset, "wraps to first offset")
	assert.Equal(t, 1, m.CurrentPage(), "wraps to first page")
}

func TestPageDownWrappingSinglePage(t *testing.T) {
	rows := make([]Row, 5)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithPaginationWrapping(true)

	// Phase 1 — cursor to last row
	m = m.PageDown()
	assert.Equal(t, 4, m.rowCursorIndex)

	// Phase 2 — at end, wraps
	m = m.PageDown()
	assert.Equal(t, 0, m.rowCursorIndex, "wraps on single page")
	assert.Equal(t, 0, m.scrollOffset, "wraps on single page")
}

func TestPageDownUnevenRows(t *testing.T) {
	rows := make([]Row, 23)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// Phase 1
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex)

	// Phase 2 — scroll to rows 10-19
	m = m.PageDown()
	assert.Equal(t, 19, m.rowCursorIndex)
	assert.Equal(t, 10, m.scrollOffset)

	// Phase 2 — scroll to rows 13-22 (clamped maxOffset=13)
	m = m.PageDown()
	assert.Equal(t, 22, m.rowCursorIndex, "cursor at last row")
	assert.Equal(t, 13, m.scrollOffset, "clamped to maxOffset")
}

func TestPageDownEmptyGrid(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithPageSize(10)

	m = m.PageDown()
	assert.Equal(t, 0, m.rowCursorIndex)
	assert.Equal(t, 0, m.scrollOffset)
}

func TestPageDownOneRow(t *testing.T) {
	rows := []Row{NewRow(RowData{"i": 0})}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// Phase 1 — cursor already at row 0 which is also visibleEnd
	m = m.PageDown()
	assert.Equal(t, 0, m.rowCursorIndex, "single row no-op")
	assert.Equal(t, 0, m.scrollOffset, "single row no-op")
}

func TestPageDownAfterArrowKeys(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithFocused(true)

	// Arrow down 3 times — cursor at row 3
	for i := 0; i < 3; i++ {
		m = m.moveDown()
	}
	assert.Equal(t, 3, m.rowCursorIndex)

	// PgDn: Phase 1 — cursor jumps to bottom of viewport (row 9)
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex, "phase 1 from mid-page")
	assert.Equal(t, 0, m.scrollOffset, "no scroll yet")
}

func TestPageUpAfterArrowKeys(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	// Start at page 2 (rows 10-19), cursor at row 15
	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithFocused(true)

	// Navigate to page 2
	m = m.PageDown() // Phase 1: cursor to 9
	m = m.PageDown() // Phase 2: scroll to 10, cursor to 19

	// Move cursor up a few rows within viewport
	for i := 0; i < 4; i++ {
		m = m.moveUp()
	}
	assert.Equal(t, 15, m.rowCursorIndex)
	assert.Equal(t, 10, m.scrollOffset)

	// PgUp: Phase 1 — cursor jumps to top of viewport (row 10)
	m = m.PageUp()
	assert.Equal(t, 10, m.rowCursorIndex, "phase 1: cursor to top of viewport")
	assert.Equal(t, 10, m.scrollOffset, "no scroll yet")
}

func TestAlternatingPageDownPageUp(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// PgDn Phase 1
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex)

	// PgDn Phase 2 — scroll to 10-19
	m = m.PageDown()
	assert.Equal(t, 19, m.rowCursorIndex)
	assert.Equal(t, 10, m.scrollOffset)

	// PgUp Phase 1 — cursor to top of viewport (10)
	m = m.PageUp()
	assert.Equal(t, 10, m.rowCursorIndex)
	assert.Equal(t, 10, m.scrollOffset)

	// PgUp Phase 2 — scroll back to 0-9
	m = m.PageUp()
	assert.Equal(t, 0, m.rowCursorIndex)
	assert.Equal(t, 0, m.scrollOffset)
}

func TestPageSizeEqualsRows(t *testing.T) {
	rows := make([]Row, 10)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10)

	// Phase 1 — cursor to last visible row
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex)
	assert.Equal(t, 0, m.scrollOffset)

	// Phase 2 — exact fit, already at end, no-op
	m = m.PageDown()
	assert.Equal(t, 9, m.rowCursorIndex, "no-op: exact fit")
	assert.Equal(t, 0, m.scrollOffset, "no-op: exact fit")
}

func TestScrollModeDefault(t *testing.T) {
	rows := make([]Row, 25)
	for i := range rows {
		rows[i] = NewRow(RowData{"i": i})
	}

	m := NewGridModel([]Column{NewColumn("i", "I", 5)}).
		WithRows(rows).
		WithPageSize(10).
		WithFocused(true)

	// Move cursor past the last visible row (index 9)
	for i := 0; i < 10; i++ {
		m = m.moveDown()
	}
	assert.Equal(t, 10, m.rowCursorIndex)
	assert.Equal(t, 1, m.scrollOffset, "scroll mode: viewport shifts by 1")

	start, end := m.VisibleIndices()
	assert.Equal(t, 1, start)
	assert.Equal(t, 10, end)
}

