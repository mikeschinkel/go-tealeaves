package teagrid

// PageSize returns the current page size, or 0 if pagination is disabled.
func (m *GridModel) PageSize() int {
	return m.pageSize
}

// CurrentPage returns the current page number (1-indexed).
func (m *GridModel) CurrentPage() int {
	return m.currentPage + 1
}

// MaxPages returns the total number of pages.
func (m *GridModel) MaxPages() int {
	totalRows := len(m.GetVisibleRows())
	if m.pageSize == 0 || totalRows == 0 {
		return 1
	}
	return (totalRows-1)/m.pageSize + 1
}

// TotalRows returns the total row count across all pages.
func (m *GridModel) TotalRows() int {
	return len(m.GetVisibleRows())
}

// VisibleIndices returns the start and end indices (0-based, inclusive)
// of the currently visible page of rows.
func (m *GridModel) VisibleIndices() (start, end int) {
	totalRows := len(m.GetVisibleRows())

	if m.pageSize == 0 {
		return 0, totalRows - 1
	}

	start = m.pageSize * m.currentPage
	end = start + m.pageSize - 1

	if end >= totalRows {
		end = totalRows - 1
	}

	return start, end
}

func (m *GridModel) pageDown() {
	if m.pageSize == 0 || len(m.GetVisibleRows()) <= m.pageSize {
		return
	}

	m.currentPage++
	maxPageIndex := m.MaxPages() - 1

	if m.currentPage > maxPageIndex {
		if m.paginationWrapping {
			m.currentPage = 0
		} else {
			m.currentPage = maxPageIndex
		}
	}

	m.rowCursorIndex = m.currentPage * m.pageSize
}

func (m *GridModel) pageUp() {
	if m.pageSize == 0 || len(m.GetVisibleRows()) <= m.pageSize {
		return
	}

	m.currentPage--
	maxPageIndex := m.MaxPages() - 1

	if m.currentPage < 0 {
		if m.paginationWrapping {
			m.currentPage = maxPageIndex
		} else {
			m.currentPage = 0
		}
	}

	m.rowCursorIndex = m.currentPage * m.pageSize
}

func (m *GridModel) pageFirst() {
	m.currentPage = 0
	m.rowCursorIndex = 0
}

func (m *GridModel) pageLast() {
	m.currentPage = m.MaxPages() - 1
	m.rowCursorIndex = m.currentPage * m.pageSize
}

func (m *GridModel) expectedPageForRowIndex(rowIndex int) int {
	if m.pageSize == 0 {
		return 0
	}
	return rowIndex / m.pageSize
}
