package teagrid

import "charm.land/bubbles/v2/key"

// GetVisibleRows returns sorted and filtered rows.
func (m *Model) GetVisibleRows() []Row {
	if m.visibleRowCacheUpdated {
		return m.visibleRowCache
	}

	rows := make([]Row, len(m.rows))
	copy(rows, m.rows)

	if m.filtered {
		rows = m.getFilteredRows(rows)
	}

	rows = getSortedRows(m.sortOrder, rows)

	m.visibleRowCache = rows
	m.visibleRowCacheUpdated = true

	return rows
}

// GetColumnSorting returns the current sort configuration (copy).
func (m *Model) GetColumnSorting() []SortColumn {
	c := make([]SortColumn, len(m.sortOrder))
	copy(c, m.sortOrder)
	return c
}

// GetHighlightedRowIndex returns the index of the highlighted row.
func (m *Model) GetHighlightedRowIndex() int {
	return m.rowCursorIndex
}

// HighlightedRow returns the currently highlighted Row.
func (m Model) HighlightedRow() Row {
	rows := m.GetVisibleRows()
	if len(rows) > 0 && m.rowCursorIndex < len(rows) {
		return rows[m.rowCursorIndex]
	}
	return Row{}
}

// SelectedRows returns all selected rows.
func (m Model) SelectedRows() []Row {
	var selected []Row
	for _, row := range m.GetVisibleRows() {
		if row.selected {
			selected = append(selected, row)
		}
	}
	return selected
}

// GetCellCursorMode returns whether cell cursor mode is enabled.
func (m *Model) GetCellCursorMode() bool {
	return m.cellCursorMode
}

// GetCellCursorColumnIndex returns the current cell cursor column index.
func (m *Model) GetCellCursorColumnIndex() int {
	return m.cellCursorColumnIndex
}

// GetFocused returns whether the grid is focused.
func (m *Model) GetFocused() bool {
	return m.focused
}

// GetCanFilter returns whether filtering is enabled.
func (m *Model) GetCanFilter() bool {
	return m.filtered
}

// GetIsFilterActive returns whether a filter is currently being applied.
func (m *Model) GetIsFilterActive() bool {
	return m.filterTextInput.Value() != ""
}

// GetIsFilterInputFocused returns whether the filter input has focus.
func (m *Model) GetIsFilterInputFocused() bool {
	return m.filterTextInput.Focused()
}

// GetCurrentFilter returns the current filter text.
func (m *Model) GetCurrentFilter() string {
	return m.filterTextInput.Value()
}

// GetHorizontalScrollColumnOffset returns the horizontal scroll offset.
func (m *Model) GetHorizontalScrollColumnOffset() int {
	return m.horizontalScrollOffsetCol
}

// GetHeaderVisibility returns header visibility.
func (m *Model) GetHeaderVisibility() bool {
	return m.headerVisible
}

// GetFooterVisibility returns footer visibility.
func (m *Model) GetFooterVisibility() bool {
	return m.footerVisible
}

// GetPaginationWrapping returns whether pagination wraps.
func (m *Model) GetPaginationWrapping() bool {
	return m.paginationWrapping
}

// NaturalWidth returns the minimum width needed to display all columns
// without flex expansion.
func (m *Model) NaturalWidth() int {
	return m.computeNaturalWidth()
}

// TotalWidth returns the total rendered width after flex column resolution.
func (m *Model) TotalWidth() int {
	return m.computeTotalWidth()
}

// Border returns the current border configuration.
func (m *Model) Border() BorderConfig {
	return m.border
}

// KeyMap returns a copy of the current key map.
func (m Model) KeyMap() KeyMap {
	return m.keyMap
}

// GetLastUpdateUserEvents returns events from the last Update call (copy).
func (m *Model) GetLastUpdateUserEvents() []UserEvent {
	if len(m.lastUpdateUserEvents) == 0 {
		return nil
	}

	returned := make([]UserEvent, len(m.lastUpdateUserEvents))
	copy(returned, m.lastUpdateUserEvents)
	return returned
}

func (m *Model) appendUserEvent(e UserEvent) {
	m.lastUpdateUserEvents = append(m.lastUpdateUserEvents, e)
}

func (m *Model) clearUserEvents() {
	m.lastUpdateUserEvents = nil
}

// hasFooter returns whether the footer should be rendered.
func (m *Model) hasFooter() bool {
	if !m.footerVisible {
		return false
	}
	return m.pageSize > 0 || m.staticFooter != "" || m.filtered
}

// --- Help interface ---

// FullHelp returns grouped key bindings for the full help view.
func (m Model) FullHelp() [][]key.Binding {
	keyBinds := [][]key.Binding{
		{m.keyMap.RowDown, m.keyMap.RowUp, m.keyMap.RowSelectToggle},
		{m.keyMap.PageDown, m.keyMap.PageUp, m.keyMap.PageFirst, m.keyMap.PageLast},
		{m.keyMap.CellLeft, m.keyMap.CellRight, m.keyMap.CellSelect},
		{m.keyMap.Filter, m.keyMap.FilterBlur, m.keyMap.FilterClear, m.keyMap.ScrollRight, m.keyMap.ScrollLeft},
	}
	if m.additionalFullHelpKeys != nil {
		keyBinds = append(keyBinds, m.additionalFullHelpKeys())
	}
	return keyBinds
}

// ShortHelp returns key bindings for the short help view.
func (m Model) ShortHelp() []key.Binding {
	keyBinds := []key.Binding{
		m.keyMap.RowDown,
		m.keyMap.RowUp,
		m.keyMap.RowSelectToggle,
		m.keyMap.PageDown,
		m.keyMap.PageUp,
		m.keyMap.CellLeft,
		m.keyMap.CellRight,
		m.keyMap.Filter,
	}
	if m.additionalShortHelpKeys != nil {
		keyBinds = append(keyBinds, m.additionalShortHelpKeys()...)
	}
	return keyBinds
}
