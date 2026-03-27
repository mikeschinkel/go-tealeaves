package teagrid

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Update handles messages per the Bubble Tea architecture.
// Uses tea.KeyPressMsg (Charm v2) instead of tea.KeyMsg (v1).
func (m GridModel) Update(msg tea.Msg) (GridModel, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	m = m.clearUserEvents()
	m = m.ensureVisibleRowsCached()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m = m.WithSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m GridModel) handleKeyPress(msg tea.KeyPressMsg) (GridModel, tea.Cmd) {
	// Filter input handling — when filter is focused, keys go to the text input
	if m.filterTextInput.Focused() {
		return m.handleFilterInput(msg)
	}

	switch {
	case key.Matches(msg, m.keyMap.RowDown):
		m = m.moveDown()

	case key.Matches(msg, m.keyMap.RowUp):
		m = m.moveUp()

	case key.Matches(msg, m.keyMap.ColRight):
		if m.colCursorMode {
			m = m.moveColRight()
		}

	case key.Matches(msg, m.keyMap.ColLeft):
		if m.colCursorMode {
			m = m.moveColLeft()
		}

	case key.Matches(msg, m.keyMap.ColSelect):
		if m.colCursorMode {
			m = m.selectCol()
		} else if m.selectableRows {
			m = m.toggleRowSelection()
		}

	case key.Matches(msg, m.keyMap.RowSelectToggle):
		if m.selectableRows {
			m = m.toggleRowSelection()
		}

	case key.Matches(msg, m.keyMap.PageDown):
		m = m.pageDown()

	case key.Matches(msg, m.keyMap.PageUp):
		m = m.pageUp()

	case key.Matches(msg, m.keyMap.PageFirst):
		m = m.pageFirst()

	case key.Matches(msg, m.keyMap.PageLast):
		m = m.pageLast()

	case key.Matches(msg, m.keyMap.Filter):
		if m.filtered {
			m.filterTextInput.Focus()
			m = m.appendUserEvent(UserEventFilterInputFocused{})
		}

	case key.Matches(msg, m.keyMap.FilterClear):
		if m.filtered && m.filterTextInput.Value() != "" {
			m.filterTextInput.SetValue("")
			m.visibleRowsDirty = true
			m = m.pageFirst()
		}

	case key.Matches(msg, m.keyMap.ScrollRight):
		m = m.scrollRight()

	case key.Matches(msg, m.keyMap.ScrollLeft):
		m = m.scrollLeft()
	}

	return m, nil
}

func (m GridModel) handleFilterInput(msg tea.KeyPressMsg) (GridModel, tea.Cmd) {
	if key.Matches(msg, m.keyMap.FilterBlur) {
		m.filterTextInput.Blur()
		m = m.appendUserEvent(UserEventFilterInputUnfocused{})
		return m, nil
	}

	prevValue := m.filterTextInput.Value()
	var cmd tea.Cmd
	m.filterTextInput, cmd = m.filterTextInput.Update(msg)

	if m.filterTextInput.Value() != prevValue {
		m.visibleRowsDirty = true
		m = m.pageFirst()
	}

	return m, cmd
}

func (m GridModel) moveDown() GridModel {
	previousIndex := m.rowCursorIndex
	totalRows := m.cursorRowBound()

	if totalRows == 0 {
		return m
	}

	m.rowCursorIndex++

	if m.rowCursorIndex >= totalRows {
		if m.rowCursorWrapping {
			m.rowCursorIndex = 0
		} else {
			m.rowCursorIndex = totalRows - 1
		}
	}

	if m.pageSize > 0 {
		m = m.ensureRowCursorVisible()
	}

	if m.rowCursorIndex != previousIndex {
		m = m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}

	return m
}

func (m GridModel) moveUp() GridModel {
	previousIndex := m.rowCursorIndex
	totalRows := m.cursorRowBound()

	if totalRows == 0 {
		return m
	}

	m.rowCursorIndex--

	if m.rowCursorIndex < 0 {
		if m.rowCursorWrapping {
			m.rowCursorIndex = totalRows - 1
		} else {
			m.rowCursorIndex = 0
		}
	}

	if m.pageSize > 0 {
		m = m.ensureRowCursorVisible()
	}

	if m.rowCursorIndex != previousIndex {
		m = m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}

	return m
}

func (m GridModel) moveColRight() GridModel {
	if len(m.columns) == 0 {
		return m
	}

	m.colCursorColumnIndex++
	if m.colCursorColumnIndex >= len(m.columns) {
		if m.colCursorWrapping {
			m.colCursorColumnIndex = 0
		} else {
			m.colCursorColumnIndex = len(m.columns) - 1
		}
	}

	m = m.ensureColCursorVisible()
	return m
}

func (m GridModel) moveColLeft() GridModel {
	if len(m.columns) == 0 {
		return m
	}

	m.colCursorColumnIndex--
	if m.colCursorColumnIndex < 0 {
		if m.colCursorWrapping {
			m.colCursorColumnIndex = len(m.columns) - 1
		} else {
			m.colCursorColumnIndex = 0
		}
	}

	m = m.ensureColCursorVisible()
	return m
}

func (m GridModel) selectCol() GridModel {
	rows := m.cachedVisibleRows
	if rows == nil {
		rows = m.VisibleRows()
	}
	if len(rows) == 0 || m.rowCursorIndex >= len(rows) {
		return m
	}

	row := rows[m.rowCursorIndex]
	colIndex := m.colCursorColumnIndex
	if colIndex >= len(m.columns) {
		return m
	}

	col := m.columns[colIndex]
	data := row.Data[col.key]

	// Extract CellValue data for the event
	if cv, ok := data.(CellValue); ok {
		data = cv.Data
	}

	m = m.appendUserEvent(UserEventCellSelected{
		RowIndex:    m.rowCursorIndex,
		ColumnIndex: colIndex,
		ColumnKey:   col.key,
		Data:        data,
	})

	return m
}

func (m GridModel) toggleRowSelection() GridModel {
	rows := m.cachedVisibleRows
	if rows == nil {
		rows = m.VisibleRows()
	}
	if len(rows) == 0 || m.rowCursorIndex >= len(rows) {
		return m
	}

	row := rows[m.rowCursorIndex]
	newSelected := !row.selected

	// Update the row in the source data
	for i, r := range m.rows {
		if r.id == row.id {
			m.rows[i] = r.Selected(newSelected)
			m.visibleRowsDirty = true
			break
		}
	}

	m = m.appendUserEvent(UserEventRowSelectToggled{
		RowIndex:   m.rowCursorIndex,
		IsSelected: newSelected,
	})

	return m
}
