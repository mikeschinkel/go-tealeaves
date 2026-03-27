package teagrid

import "charm.land/bubbles/v2/key"

// VisibleRows returns sorted and filtered rows.
func (m GridModel) VisibleRows() []Row {
	if !m.visibleRowsDirty && m.cachedVisibleRows != nil {
		return m.cachedVisibleRows
	}

	rows := make([]Row, len(m.rows))
	copy(rows, m.rows)

	if m.filtered {
		rows = m.getFilteredRows(rows)
	}

	rows = getSortedRows(m.sortOrder, rows)

	return rows
}

func (m GridModel) ensureVisibleRowsCached() GridModel {
	if !m.visibleRowsDirty && m.cachedVisibleRows != nil {
		return m
	}
	m.cachedVisibleRows = m.VisibleRows()
	m.visibleRowsDirty = false
	return m
}

// cursorRowBound returns the number of rows the cursor can visit.
// When dataRowCount is set, cursor stops at the data boundary.
func (m GridModel) cursorRowBound() int {
	total := m.visibleRowCount()
	if m.dataRowCount > 0 && m.dataRowCount < total {
		return m.dataRowCount
	}
	return total
}

func (m GridModel) visibleRowCount() int {
	if !m.visibleRowsDirty && m.cachedVisibleRows != nil {
		return len(m.cachedVisibleRows)
	}
	return len(m.VisibleRows())
}

// ColumnSorting returns the current sort configuration (copy).
func (m GridModel) ColumnSorting() []SortColumn {
	c := make([]SortColumn, len(m.sortOrder))
	copy(c, m.sortOrder)
	return c
}

// HighlightedRowIndex returns the index of the highlighted row.
func (m GridModel) HighlightedRowIndex() int {
	return m.rowCursorIndex
}

// HighlightedRow returns the currently highlighted Row.
func (m GridModel) HighlightedRow() Row {
	rows := m.cachedVisibleRows
	if rows == nil {
		rows = m.VisibleRows()
	}
	if len(rows) > 0 && m.rowCursorIndex < len(rows) {
		return rows[m.rowCursorIndex]
	}
	return Row{}
}

// SelectedRows returns all selected rows.
func (m GridModel) SelectedRows() []Row {
	rows := m.cachedVisibleRows
	if rows == nil {
		rows = m.VisibleRows()
	}
	var selected []Row
	for _, row := range rows {
		if row.selected {
			selected = append(selected, row)
		}
	}
	return selected
}

// ColCursorMode returns whether column cursor mode is enabled.
func (m GridModel) ColCursorMode() bool {
	return m.colCursorMode
}

// ColCursorWrapping returns whether the column cursor wraps around columns.
func (m GridModel) ColCursorWrapping() bool {
	return m.colCursorWrapping
}

// ColCursorColumnIndex returns the current column cursor column index.
func (m GridModel) ColCursorColumnIndex() int {
	return m.colCursorColumnIndex
}

// RowCursorWrapping returns whether the row cursor wraps.
func (m GridModel) RowCursorWrapping() bool {
	return m.rowCursorWrapping
}

// IsFocused returns whether the grid is focused.
func (m GridModel) IsFocused() bool {
	return m.focused
}

// CanFilter returns whether filtering is enabled.
func (m GridModel) CanFilter() bool {
	return m.filtered
}

// IsFilterActive returns whether a filter is currently being applied.
func (m GridModel) IsFilterActive() bool {
	return m.filterTextInput.Value() != ""
}

// IsFilterInputFocused returns whether the filter input has focus.
func (m GridModel) IsFilterInputFocused() bool {
	return m.filterTextInput.Focused()
}

// CurrentFilter returns the current filter text.
func (m GridModel) CurrentFilter() string {
	return m.filterTextInput.Value()
}

// FillWidth returns whether fill-width mode is enabled.
func (m GridModel) FillWidth() bool {
	return m.fillWidth
}

// HorizontalScrollColumnOffset returns the horizontal scroll offset.
func (m GridModel) HorizontalScrollColumnOffset() int {
	return m.horizontalScrollOffsetCol
}

// IsHeaderVisible returns header visibility.
func (m GridModel) IsHeaderVisible() bool {
	return m.headerVisible
}

// IsFooterVisible returns footer visibility.
func (m GridModel) IsFooterVisible() bool {
	return m.footerVisible
}

// IsPaginationWrapping returns whether pagination wraps.
func (m GridModel) IsPaginationWrapping() bool {
	return m.paginationWrapping
}

// ScrollOffset returns the current vertical scroll offset.
func (m GridModel) ScrollOffset() int {
	return m.scrollOffset
}

// NaturalWidth returns the minimum width needed to display all columns
// without flex expansion.
func (m GridModel) NaturalWidth() int {
	return m.computeNaturalWidth()
}

// TotalWidth returns the total rendered width after flex column resolution.
func (m GridModel) TotalWidth() int {
	return m.computeTotalWidth()
}

// Border returns the current border configuration.
func (m GridModel) Border() BorderConfig {
	return m.border
}

// KeyMap returns a copy of the current key map.
func (m GridModel) KeyMap() KeyMap {
	return m.keyMap
}

// LastUpdateUserEvents returns events from the last Update call (copy).
func (m GridModel) LastUpdateUserEvents() []UserEvent {
	if len(m.lastUpdateUserEvents) == 0 {
		return nil
	}

	returned := make([]UserEvent, len(m.lastUpdateUserEvents))
	copy(returned, m.lastUpdateUserEvents)
	return returned
}

func (m GridModel) appendUserEvent(e UserEvent) GridModel {
	m.lastUpdateUserEvents = append(m.lastUpdateUserEvents, e)
	return m
}

func (m GridModel) clearUserEvents() GridModel {
	m.lastUpdateUserEvents = nil
	return m
}

// hasFooter returns whether the footer should be rendered.
func (m GridModel) hasFooter() bool {
	if !m.footerVisible {
		return false
	}
	return m.pageSize > 0 || m.staticFooter != "" || m.filtered || len(m.footerRows) > 0
}

// hasInfoRow returns whether the info row (static text + pagination + filter) should render.
func (m GridModel) hasInfoRow() bool {
	return m.pageSize > 0 || m.staticFooter != "" || m.filtered
}

// --- Help interface ---

// FullHelp returns grouped key bindings for the full help view.
func (m GridModel) FullHelp() [][]key.Binding {
	keyBinds := [][]key.Binding{
		{m.keyMap.RowDown, m.keyMap.RowUp, m.keyMap.RowSelectToggle},
		{m.keyMap.PageDown, m.keyMap.PageUp, m.keyMap.PageFirst, m.keyMap.PageLast},
		{m.keyMap.ColLeft, m.keyMap.ColRight, m.keyMap.ColSelect},
		{m.keyMap.Filter, m.keyMap.FilterBlur, m.keyMap.FilterClear, m.keyMap.ScrollRight, m.keyMap.ScrollLeft},
	}
	if m.additionalFullHelpKeys != nil {
		keyBinds = append(keyBinds, m.additionalFullHelpKeys())
	}
	return keyBinds
}

// ShortHelp returns key bindings for the short help view.
func (m GridModel) ShortHelp() []key.Binding {
	keyBinds := []key.Binding{
		m.keyMap.RowDown,
		m.keyMap.RowUp,
		m.keyMap.RowSelectToggle,
		m.keyMap.PageDown,
		m.keyMap.PageUp,
		m.keyMap.ColLeft,
		m.keyMap.ColRight,
		m.keyMap.Filter,
	}
	if m.additionalShortHelpKeys != nil {
		keyBinds = append(keyBinds, m.additionalShortHelpKeys()...)
	}
	return keyBinds
}
