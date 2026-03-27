package teagrid

// PageSize returns the current page size, or 0 if pagination is disabled.
func (m GridModel) PageSize() int {
	return m.pageSize
}

// CurrentPage returns the current page number (1-indexed).
func (m GridModel) CurrentPage() int {
	return m.currentPage + 1
}

// MaxPages returns the total number of pages.
func (m GridModel) MaxPages() int {
	totalRows := m.visibleRowCount()
	if m.pageSize == 0 || totalRows == 0 {
		return 1
	}
	return (totalRows-1)/m.pageSize + 1
}

// TotalRows returns the total row count across all pages.
func (m GridModel) TotalRows() int {
	return m.visibleRowCount()
}

// VisibleIndices returns the start and end indices (0-based, inclusive)
// of the currently visible page of rows.
func (m GridModel) VisibleIndices() (start, end int) {
	totalRows := m.visibleRowCount()

	if m.pageSize == 0 {
		return 0, totalRows - 1
	}

	start = m.scrollOffset
	end = start + m.pageSize - 1

	if end >= totalRows {
		end = totalRows - 1
	}

	return start, end
}

func (m GridModel) pageDown() GridModel {
	if m.pageSize == 0 {
		return m
	}
	totalRows := m.visibleRowCount()
	if totalRows == 0 {
		return m
	}

	visibleEnd := m.scrollOffset + m.pageSize - 1
	if visibleEnd >= totalRows {
		visibleEnd = totalRows - 1
	}

	if m.rowCursorIndex < visibleEnd {
		// Phase 1: jump cursor to bottom of viewport
		m.rowCursorIndex = visibleEnd
	} else {
		// Phase 2: scroll viewport forward by pageSize
		maxOffset := totalRows - m.pageSize
		if maxOffset < 0 {
			maxOffset = 0
		}

		// Already at the end — wrap or no-op
		if m.scrollOffset >= maxOffset {
			if m.paginationWrapping {
				m.rowCursorIndex = 0
				m.scrollOffset = 0
				m.currentPage = 0
				return m
			}
			return m
		}

		// Scroll forward, clamping to maxOffset
		newOffset := m.scrollOffset + m.pageSize
		if newOffset > maxOffset {
			newOffset = maxOffset
		}
		m.scrollOffset = newOffset
		m.rowCursorIndex = m.scrollOffset + m.pageSize - 1
		if m.rowCursorIndex >= totalRows {
			m.rowCursorIndex = totalRows - 1
		}
	}
	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	return m
}

func (m GridModel) pageUp() GridModel {
	if m.pageSize == 0 {
		return m
	}
	totalRows := m.visibleRowCount()
	if totalRows == 0 {
		return m
	}

	if m.rowCursorIndex > m.scrollOffset {
		// Phase 1: jump cursor to top of viewport
		m.rowCursorIndex = m.scrollOffset
	} else {
		// Phase 2: scroll viewport backward by pageSize

		// Already at the beginning — wrap or no-op
		if m.scrollOffset == 0 {
			if m.paginationWrapping {
				maxOffset := totalRows - m.pageSize
				if maxOffset < 0 {
					maxOffset = 0
				}
				m.scrollOffset = maxOffset
				m.rowCursorIndex = totalRows - 1
				m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
				return m
			}
			return m
		}

		// Scroll backward, clamping to 0
		newOffset := m.scrollOffset - m.pageSize
		if newOffset < 0 {
			newOffset = 0
		}
		m.scrollOffset = newOffset
		m.rowCursorIndex = newOffset
	}
	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	return m
}

func (m GridModel) pageFirst() GridModel {
	m.scrollOffset = 0
	m.currentPage = 0
	m.rowCursorIndex = 0
	return m
}

func (m GridModel) pageLast() GridModel {
	totalRows := m.visibleRowCount()
	if totalRows == 0 {
		return m
	}
	maxOffset := totalRows - m.pageSize
	if maxOffset < 0 {
		maxOffset = 0
	}
	m.scrollOffset = maxOffset
	m.rowCursorIndex = totalRows - 1
	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
	return m
}

func (m GridModel) expectedPageForRowIndex(rowIndex int) int {
	if m.pageSize == 0 {
		return 0
	}
	return rowIndex / m.pageSize
}
