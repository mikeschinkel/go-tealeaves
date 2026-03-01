package teagrid

func (m *Model) scrollRight() {
	if m.horizontalScrollOffsetCol < m.maxHorizontalColumnIndex {
		m.horizontalScrollOffsetCol++
	}
}

func (m *Model) scrollLeft() {
	if m.horizontalScrollOffsetCol > 0 {
		m.horizontalScrollOffsetCol--
	}
}

func (m *Model) recalculateLastHorizontalColumn() {
	if m.viewportWidth == 0 {
		m.maxHorizontalColumnIndex = 0
		return
	}

	totalWidth := m.computeTotalWidth()
	if totalWidth <= m.viewportWidth {
		m.maxHorizontalColumnIndex = 0
		return
	}

	if m.horizontalScrollFreezeColumnsCount >= len(m.columns) {
		m.maxHorizontalColumnIndex = 0
		return
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
}
