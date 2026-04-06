package teagrid

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateNotFocused(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows([]Row{NewRow(RowData{"x": 1})})

	updated, cmd := m.Update(tea.KeyPressMsg{})
	assert.Nil(t, cmd)
	assert.NotNil(t, updated)
}

func TestMoveDown(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
		NewRow(RowData{"x": 3}),
	}

	t.Run("clamps by default", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithRows(rows).
			WithFocused(true)

		m = m.moveDown()
		m = m.moveDown()
		assert.Equal(t, 2, m.rowCursorIndex)

		m = m.moveDown()
		assert.Equal(t, 2, m.rowCursorIndex, "should clamp at last row")
	})

	t.Run("wraps when wrapping enabled", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithRows(rows).
			WithRowCursorWrapping(true).
			WithFocused(true)

		assert.Equal(t, 0, m.rowCursorIndex)

		m = m.moveDown()
		assert.Equal(t, 1, m.rowCursorIndex)

		m = m.moveDown()
		assert.Equal(t, 2, m.rowCursorIndex)

		m = m.moveDown()
		assert.Equal(t, 0, m.rowCursorIndex, "should wrap to start")
	})
}

func TestMoveUp(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}

	t.Run("clamps by default", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithRows(rows).
			WithFocused(true)

		m = m.moveUp()
		assert.Equal(t, 0, m.rowCursorIndex, "should clamp at first row")
	})

	t.Run("wraps when wrapping enabled", func(t *testing.T) {
		m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
			WithRows(rows).
			WithRowCursorWrapping(true).
			WithFocused(true)

		m = m.moveUp()
		assert.Equal(t, 1, m.rowCursorIndex, "should wrap to end")

		m = m.moveUp()
		assert.Equal(t, 0, m.rowCursorIndex)
	})
}

func TestMoveDownEmitsEvent(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}

	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		WithFocused(true)

	m = m.clearUserEvents()
	m = m.moveDown()

	events := m.LastUpdateUserEvents()
	require.Len(t, events, 1)

	event, ok := events[0].(UserEventHighlightedIndexChanged)
	require.True(t, ok)
	assert.Equal(t, 0, event.PreviousRowIndex)
	assert.Equal(t, 1, event.SelectedRowIndex)
}

func TestMoveColRight(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
		NewColumn("c", "C", 5),
	}

	t.Run("clamps at last column by default", func(t *testing.T) {
		m := NewGridModel(cols).
			WithColCursorMode(true).
			WithFocused(true)

		assert.Equal(t, 0, m.colCursorColumnIndex)

		m = m.moveColRight()
		assert.Equal(t, 1, m.colCursorColumnIndex)

		m = m.moveColRight()
		assert.Equal(t, 2, m.colCursorColumnIndex)

		m = m.moveColRight()
		assert.Equal(t, 2, m.colCursorColumnIndex, "should clamp at last column")
	})

	t.Run("wraps when wrapping enabled", func(t *testing.T) {
		m := NewGridModel(cols).
			WithColCursorMode(true).
			WithColCursorWrapping(true).
			WithFocused(true)

		m = m.moveColRight()
		m = m.moveColRight()
		assert.Equal(t, 2, m.colCursorColumnIndex)

		m = m.moveColRight()
		assert.Equal(t, 0, m.colCursorColumnIndex, "should wrap to first")
	})
}

func TestMoveColLeft(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
	}

	t.Run("clamps at first column by default", func(t *testing.T) {
		m := NewGridModel(cols).
			WithColCursorMode(true).
			WithFocused(true)

		m = m.moveColLeft()
		assert.Equal(t, 0, m.colCursorColumnIndex, "should clamp at first column")
	})

	t.Run("wraps when wrapping enabled", func(t *testing.T) {
		m := NewGridModel(cols).
			WithColCursorMode(true).
			WithColCursorWrapping(true).
			WithFocused(true)

		m = m.moveColLeft()
		assert.Equal(t, 1, m.colCursorColumnIndex, "should wrap to last")
	})
}

func TestSelectCol(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"name": "Alice", "age": 30}),
	}

	m := NewGridModel([]Column{
		NewColumn("name", "Name", 10),
		NewColumn("age", "Age", 5),
	}).WithRows(rows).
		WithColCursorMode(true).
		WithFocused(true)

	m.colCursorColumnIndex = 1
	m = m.clearUserEvents()
	m = m.selectCol()

	events := m.LastUpdateUserEvents()
	require.Len(t, events, 1)

	event, ok := events[0].(UserEventCellSelected)
	require.True(t, ok)
	assert.Equal(t, 0, event.RowIndex)
	assert.Equal(t, 1, event.ColumnIndex)
	assert.Equal(t, "age", event.ColumnKey)
	assert.Equal(t, 30, event.Data)
}

func TestToggleRowSelection(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}

	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		WithSelectableRows(true).
		WithFocused(true)

	m = m.clearUserEvents()
	m = m.toggleRowSelection()

	events := m.LastUpdateUserEvents()
	require.Len(t, events, 1)

	event, ok := events[0].(UserEventRowSelectToggled)
	require.True(t, ok)
	assert.Equal(t, 0, event.RowIndex)
	assert.True(t, event.IsSelected)

	// Toggle again
	m = m.clearUserEvents()
	m = m.toggleRowSelection()

	events = m.LastUpdateUserEvents()
	require.Len(t, events, 1)
	event = events[0].(UserEventRowSelectToggled)
	assert.False(t, event.IsSelected)
}

func TestWindowSizeMsg(t *testing.T) {
	m := NewGridModel([]Column{NewColumn("x", "X", 5)}).WithFocused(true)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	assert.Equal(t, 80, updated.viewportWidth)
	assert.Equal(t, 24, updated.viewportHeight)
}
