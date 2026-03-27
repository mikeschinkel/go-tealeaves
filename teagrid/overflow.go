package teagrid

// OverflowConfig controls how horizontal overflow is indicated.
// Default is no indicator — no width is stolen from columns.
type OverflowConfig struct {
	// LeftIndicator is shown when content extends left of the viewport.
	LeftIndicator string

	// RightIndicator is shown when content extends right of the viewport.
	RightIndicator string
}

// DefaultOverflowConfig returns the default overflow config with no indicators.
func DefaultOverflowConfig() OverflowConfig {
	return OverflowConfig{}
}

// WithIndicators returns an overflow config with simple arrow indicators.
func WithIndicators() OverflowConfig {
	return OverflowConfig{
		LeftIndicator:  "<",
		RightIndicator: ">",
	}
}

// buildColumnKeyIndex builds a map from column key to index for O(1) lookup.
func (m GridModel) buildColumnKeyIndex() GridModel {
	m.columnKeyIndex = make(map[string]int, len(m.columns))
	for i, col := range m.columns {
		m.columnKeyIndex[col.key] = i
	}
	return m
}

// columnOriginalIndex resolves the original index of a column in m.columns
// from its key. Returns -1 if the key is not found.
func (m GridModel) columnOriginalIndex(key string) int {
	if idx, ok := m.columnKeyIndex[key]; ok {
		return idx
	}
	return -1
}

// visibleColumns returns the columns that fit in the viewport,
// accounting for frozen columns and horizontal scroll offset.
// Uses a cache that is invalidated when scroll position or column widths change.
func (m GridModel) visibleColumns() []Column {
	if !m.visibleColumnsDirty && m.cachedVisibleColumns != nil {
		return m.cachedVisibleColumns
	}
	return m.computeVisibleColumns()
}

func (m GridModel) ensureVisibleColumnsCached() GridModel {
	if !m.visibleColumnsDirty && m.cachedVisibleColumns != nil {
		return m
	}
	m.cachedVisibleColumns = m.computeVisibleColumns()
	m.visibleColumnsDirty = false
	return m
}

// computeVisibleColumns calculates which columns fit in the viewport.
func (m GridModel) computeVisibleColumns() []Column {
	if m.viewportWidth == 0 || m.computeTotalWidth() <= m.viewportWidth {
		if m.fillWidth && m.viewportWidth > 0 {
			visible := make([]Column, len(m.columns))
			copy(visible, m.columns)
			return m.applyFillWidth(visible)
		}
		return m.columns
	}

	var visible []Column

	// Frozen columns always visible
	frozenCount := m.horizontalScrollFreezeColumnsCount
	if frozenCount > len(m.columns) {
		frozenCount = len(m.columns)
	}

	for i := 0; i < frozenCount; i++ {
		visible = append(visible, m.columns[i])
	}

	// Scrollable columns from offset
	scrollStart := frozenCount + m.horizontalScrollOffsetCol
	if scrollStart >= len(m.columns) {
		if m.fillWidth {
			return m.applyFillWidth(visible)
		}
		return visible
	}

	// Add columns that fit within viewport
	usedWidth := m.border.OuterWidth()
	for _, col := range visible {
		usedWidth += col.RenderWidth() + m.border.InnerDividerWidth()
	}

	for i := scrollStart; i < len(m.columns); i++ {
		colWidth := m.columns[i].RenderWidth()
		if len(visible) > 0 {
			colWidth += m.border.InnerDividerWidth()
		}

		if usedWidth+colWidth > m.viewportWidth {
			break
		}

		usedWidth += colWidth
		visible = append(visible, m.columns[i])
	}

	if m.fillWidth {
		return m.applyFillWidth(visible)
	}
	return visible
}

// applyFillWidth pads the last visible column's right padding to fill
// the viewport width, eliminating the gap on the right edge.
func (m GridModel) applyFillWidth(visible []Column) []Column {
	if len(visible) == 0 || m.viewportWidth == 0 {
		return visible
	}

	usedWidth := m.border.OuterWidth()
	for i, col := range visible {
		usedWidth += col.RenderWidth()
		if i < len(visible)-1 {
			usedWidth += m.border.InnerDividerWidth()
		}
	}

	gap := m.viewportWidth - usedWidth
	if gap > 0 {
		visible[len(visible)-1].paddingRight += gap
	}

	return visible
}
