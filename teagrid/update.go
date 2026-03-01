package teagrid

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// Update handles messages per the Bubble Tea architecture.
// Uses tea.KeyPressMsg (Charm v2) instead of tea.KeyMsg (v1).
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	m.clearUserEvents()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	// Filter input handling — when filter is focused, keys go to the text input
	if m.filterTextInput.Focused() {
		return m.handleFilterInput(msg)
	}

	switch {
	case key.Matches(msg, m.keyMap.RowDown):
		m = m.moveDown()

	case key.Matches(msg, m.keyMap.RowUp):
		m = m.moveUp()

	case key.Matches(msg, m.keyMap.CellRight):
		if m.cellCursorMode {
			m = m.moveCellRight()
		}

	case key.Matches(msg, m.keyMap.CellLeft):
		if m.cellCursorMode {
			m = m.moveCellLeft()
		}

	case key.Matches(msg, m.keyMap.CellSelect):
		if m.cellCursorMode {
			m = m.selectCell()
		} else if m.selectableRows {
			m = m.toggleRowSelection()
		}

	case key.Matches(msg, m.keyMap.RowSelectToggle):
		if m.selectableRows {
			m = m.toggleRowSelection()
		}

	case key.Matches(msg, m.keyMap.PageDown):
		m.pageDown()

	case key.Matches(msg, m.keyMap.PageUp):
		m.pageUp()

	case key.Matches(msg, m.keyMap.PageFirst):
		m.pageFirst()

	case key.Matches(msg, m.keyMap.PageLast):
		m.pageLast()

	case key.Matches(msg, m.keyMap.Filter):
		if m.filtered {
			m.filterTextInput.Focus()
			m.appendUserEvent(UserEventFilterInputFocused{})
		}

	case key.Matches(msg, m.keyMap.FilterClear):
		if m.filtered && m.filterTextInput.Value() != "" {
			m.filterTextInput.SetValue("")
			m.visibleRowCacheUpdated = false
			m.pageFirst()
		}

	case key.Matches(msg, m.keyMap.ScrollRight):
		m.scrollRight()

	case key.Matches(msg, m.keyMap.ScrollLeft):
		m.scrollLeft()
	}

	return m, nil
}

func (m Model) handleFilterInput(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	if key.Matches(msg, m.keyMap.FilterBlur) {
		m.filterTextInput.Blur()
		m.appendUserEvent(UserEventFilterInputUnfocused{})
		return m, nil
	}

	var cmd tea.Cmd
	m.filterTextInput, cmd = m.filterTextInput.Update(msg)
	m.visibleRowCacheUpdated = false
	m.pageFirst()

	return m, cmd
}

func (m Model) moveDown() Model {
	previousIndex := m.rowCursorIndex
	totalRows := len(m.GetVisibleRows())

	if totalRows == 0 {
		return m
	}

	m.rowCursorIndex++

	if m.rowCursorIndex >= totalRows {
		m.rowCursorIndex = 0
	}

	if m.pageSize > 0 {
		m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	}

	if m.rowCursorIndex != previousIndex {
		m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}

	return m
}

func (m Model) moveUp() Model {
	previousIndex := m.rowCursorIndex
	totalRows := len(m.GetVisibleRows())

	if totalRows == 0 {
		return m
	}

	m.rowCursorIndex--

	if m.rowCursorIndex < 0 {
		m.rowCursorIndex = totalRows - 1
	}

	if m.pageSize > 0 {
		m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	}

	if m.rowCursorIndex != previousIndex {
		m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}

	return m
}

func (m Model) moveCellRight() Model {
	if len(m.columns) == 0 {
		return m
	}

	m.cellCursorColumnIndex++
	if m.cellCursorColumnIndex >= len(m.columns) {
		m.cellCursorColumnIndex = 0
	}

	return m
}

func (m Model) moveCellLeft() Model {
	if len(m.columns) == 0 {
		return m
	}

	m.cellCursorColumnIndex--
	if m.cellCursorColumnIndex < 0 {
		m.cellCursorColumnIndex = len(m.columns) - 1
	}

	return m
}

func (m Model) selectCell() Model {
	rows := m.GetVisibleRows()
	if len(rows) == 0 || m.rowCursorIndex >= len(rows) {
		return m
	}

	row := rows[m.rowCursorIndex]
	colIndex := m.cellCursorColumnIndex
	if colIndex >= len(m.columns) {
		return m
	}

	col := m.columns[colIndex]
	data := row.Data[col.key]

	// Extract CellValue data for the event
	if cv, ok := data.(CellValue); ok {
		data = cv.Data
	}

	m.appendUserEvent(UserEventCellSelected{
		RowIndex:    m.rowCursorIndex,
		ColumnIndex: colIndex,
		ColumnKey:   col.key,
		Data:        data,
	})

	return m
}

func (m Model) toggleRowSelection() Model {
	rows := m.GetVisibleRows()
	if len(rows) == 0 || m.rowCursorIndex >= len(rows) {
		return m
	}

	row := rows[m.rowCursorIndex]
	newSelected := !row.selected

	// Update the row in the source data
	for i, r := range m.rows {
		if r.id == row.id {
			m.rows[i] = r.Selected(newSelected)
			break
		}
	}

	m.visibleRowCacheUpdated = false

	m.appendUserEvent(UserEventRowSelectToggled{
		RowIndex:   m.rowCursorIndex,
		IsSelected: newSelected,
	})

	return m
}
