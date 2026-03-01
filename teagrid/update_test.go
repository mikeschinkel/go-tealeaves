package teagrid

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newKeyPress(keyStr string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: rune(keyStr[0])}
}

func TestUpdateNotFocused(t *testing.T) {
	m := New([]Column{NewColumn("x", "X", 5)}).
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

	m := New([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		Focused(true)

	assert.Equal(t, 0, m.rowCursorIndex)

	m = m.moveDown()
	assert.Equal(t, 1, m.rowCursorIndex)

	m = m.moveDown()
	assert.Equal(t, 2, m.rowCursorIndex)

	m = m.moveDown()
	assert.Equal(t, 0, m.rowCursorIndex, "should wrap to start")
}

func TestMoveUp(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}

	m := New([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		Focused(true)

	m = m.moveUp()
	assert.Equal(t, 1, m.rowCursorIndex, "should wrap to end")

	m = m.moveUp()
	assert.Equal(t, 0, m.rowCursorIndex)
}

func TestMoveDownEmitsEvent(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"x": 1}),
		NewRow(RowData{"x": 2}),
	}

	m := New([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		Focused(true)

	m.clearUserEvents()
	m = m.moveDown()

	events := m.GetLastUpdateUserEvents()
	require.Len(t, events, 1)

	event, ok := events[0].(UserEventHighlightedIndexChanged)
	require.True(t, ok)
	assert.Equal(t, 0, event.PreviousRowIndex)
	assert.Equal(t, 1, event.SelectedRowIndex)
}

func TestMoveCellRight(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
		NewColumn("c", "C", 5),
	}

	m := New(cols).
		WithCellCursorMode(true).
		Focused(true)

	assert.Equal(t, 0, m.cellCursorColumnIndex)

	m = m.moveCellRight()
	assert.Equal(t, 1, m.cellCursorColumnIndex)

	m = m.moveCellRight()
	assert.Equal(t, 2, m.cellCursorColumnIndex)

	m = m.moveCellRight()
	assert.Equal(t, 0, m.cellCursorColumnIndex, "should wrap")
}

func TestMoveCellLeft(t *testing.T) {
	cols := []Column{
		NewColumn("a", "A", 5),
		NewColumn("b", "B", 5),
	}

	m := New(cols).
		WithCellCursorMode(true).
		Focused(true)

	m = m.moveCellLeft()
	assert.Equal(t, 1, m.cellCursorColumnIndex, "should wrap to last")
}

func TestSelectCell(t *testing.T) {
	rows := []Row{
		NewRow(RowData{"name": "Alice", "age": 30}),
	}

	m := New([]Column{
		NewColumn("name", "Name", 10),
		NewColumn("age", "Age", 5),
	}).WithRows(rows).
		WithCellCursorMode(true).
		Focused(true)

	m.cellCursorColumnIndex = 1
	m.clearUserEvents()
	m = m.selectCell()

	events := m.GetLastUpdateUserEvents()
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

	m := New([]Column{NewColumn("x", "X", 5)}).
		WithRows(rows).
		WithSelectableRows(true).
		Focused(true)

	m.clearUserEvents()
	m = m.toggleRowSelection()

	events := m.GetLastUpdateUserEvents()
	require.Len(t, events, 1)

	event, ok := events[0].(UserEventRowSelectToggled)
	require.True(t, ok)
	assert.Equal(t, 0, event.RowIndex)
	assert.True(t, event.IsSelected)

	// Toggle again
	m.clearUserEvents()
	m = m.toggleRowSelection()

	events = m.GetLastUpdateUserEvents()
	require.Len(t, events, 1)
	event = events[0].(UserEventRowSelectToggled)
	assert.False(t, event.IsSelected)
}

func TestWindowSizeMsg(t *testing.T) {
	m := New([]Column{NewColumn("x", "X", 5)}).Focused(true)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	assert.Equal(t, 80, updated.viewportWidth)
	assert.Equal(t, 24, updated.viewportHeight)
}
