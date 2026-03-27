package teagrid

// clampCursors ensures the row and column cursors stay within bounds
// after a viewport resize reduces visible rows or columns.
func (m GridModel) clampCursors() GridModel {
	// Clamp row cursor
	totalRows := m.cursorRowBound()
	if totalRows > 0 && m.rowCursorIndex >= totalRows {
		m.rowCursorIndex = totalRows - 1
	}

	// Clamp column cursor against total columns (not visible columns).
	// The column cursor is an absolute index into m.columns; horizontal
	// scrolling handles which columns are visible in the viewport.
	if len(m.columns) > 0 && m.colCursorColumnIndex >= len(m.columns) {
		m.colCursorColumnIndex = len(m.columns) - 1
	}

	// Ensure scroll offsets keep cursors visible
	if m.pageSize > 0 {
		m = m.ensureRowCursorVisible()
	}
	m = m.ensureColCursorVisible()

	return m
}

func (m GridModel) hasHiddenColumnsLeft() bool {
	return m.horizontalScrollOffsetCol > 0
}

func (m GridModel) hasHiddenColumnsRight() bool {
	return m.maxHorizontalColumnIndex > 0 &&
		m.horizontalScrollOffsetCol < m.maxHorizontalColumnIndex
}

func (m GridModel) scrollRight() GridModel {
	if m.horizontalScrollOffsetCol < m.maxHorizontalColumnIndex {
		m.horizontalScrollOffsetCol++
		m.visibleColumnsDirty = true
	}
	return m
}

func (m GridModel) scrollLeft() GridModel {
	if m.horizontalScrollOffsetCol > 0 {
		m.horizontalScrollOffsetCol--
		m.visibleColumnsDirty = true
	}
	return m
}

// ensureRowCursorVisible adjusts scrollOffset so the row cursor is visible.
// The viewport shifts by the minimum amount needed (1 row).
func (m GridModel) ensureRowCursorVisible() GridModel {
	if m.pageSize == 0 {
		return m
	}

	visibleEnd := m.scrollOffset + m.pageSize - 1

	if m.rowCursorIndex > visibleEnd {
		m.scrollOffset = m.rowCursorIndex - m.pageSize + 1
	} else if m.rowCursorIndex < m.scrollOffset {
		m.scrollOffset = m.rowCursorIndex
	}

	// Clamp scrollOffset
	totalRows := m.visibleRowCount()
	maxOffset := totalRows - m.pageSize
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	return m
}

// ensureColCursorVisible adjusts horizontalScrollOffsetCol so the column
// cursor column is within the visible set. Frozen columns need no scrolling.
func (m GridModel) ensureColCursorVisible() GridModel {
	if m.viewportWidth == 0 || len(m.columns) == 0 {
		return m
	}

	cursorCol := m.colCursorColumnIndex
	freezeCount := m.horizontalScrollFreezeColumnsCount

	// Frozen columns are always visible — no scroll adjustment needed
	if cursorCol < freezeCount {
		return m
	}

	// scrollableIndex is the cursor's position within the scrollable columns
	scrollableIndex := cursorCol - freezeCount

	// If cursor is before the scroll window, scroll left to reveal it
	if scrollableIndex < m.horizontalScrollOffsetCol {
		m.horizontalScrollOffsetCol = scrollableIndex
		m.visibleColumnsDirty = true
		return m
	}

	// Scroll right until cursor column fits in viewport
	for !m.isColumnVisible(cursorCol) {
		if m.horizontalScrollOffsetCol < m.maxHorizontalColumnIndex {
			m.horizontalScrollOffsetCol++
		} else {
			break
		}
	}
	m.visibleColumnsDirty = true

	return m
}

// isColumnVisible checks whether the given column index fits in the viewport
// without allocating a visible column slice.
func (m GridModel) isColumnVisible(colIndex int) bool {
	freezeCount := m.horizontalScrollFreezeColumnsCount
	usedWidth := m.border.OuterWidth()

	for i := 0; i < freezeCount && i < len(m.columns); i++ {
		usedWidth += m.columns[i].RenderWidth()
		if i < len(m.columns)-1 {
			usedWidth += m.border.InnerDividerWidth()
		}
	}

	scrollStart := freezeCount + m.horizontalScrollOffsetCol
	for i := scrollStart; i < len(m.columns); i++ {
		colWidth := m.columns[i].RenderWidth()
		if usedWidth > m.border.OuterWidth() {
			colWidth += m.border.InnerDividerWidth()
		}
		if usedWidth+colWidth > m.viewportWidth {
			return false
		}
		usedWidth += colWidth
		if i == colIndex {
			return true
		}
	}
	return false
}

func (m GridModel) recalculateLastHorizontalColumn() GridModel {
	if m.viewportWidth == 0 {
		m.maxHorizontalColumnIndex = 0
		return m
	}

	totalWidth := m.computeTotalWidth()
	if totalWidth <= m.viewportWidth {
		m.maxHorizontalColumnIndex = 0
		return m
	}

	if m.horizontalScrollFreezeColumnsCount >= len(m.columns) {
		m.maxHorizontalColumnIndex = 0
		return m
	}

	// Compute width consumed by frozen columns + outer border
	visibleWidth := m.border.OuterWidth()

	for i := 0; i < m.horizontalScrollFreezeColumnsCount && i < len(m.columns); i++ {
		visibleWidth += m.columns[i].RenderWidth()
		if i < len(m.columns)-1 {
			visibleWidth += m.border.InnerDividerWidth()
		}
	}

	m.maxHorizontalColumnIndex = len(m.columns) - 1

	// Work backwards from the right to find the maximum scroll offset
	for i := len(m.columns) - 1; i >= m.horizontalScrollFreezeColumnsCount && visibleWidth <= m.viewportWidth; i-- {
		visibleWidth += m.columns[i].RenderWidth()
		if i > m.horizontalScrollFreezeColumnsCount {
			visibleWidth += m.border.InnerDividerWidth()
		}

		if visibleWidth <= m.viewportWidth {
			m.maxHorizontalColumnIndex = i - m.horizontalScrollFreezeColumnsCount
		}
	}

	return m
}
