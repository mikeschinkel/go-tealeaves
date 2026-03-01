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

// visibleColumns returns the columns that fit in the viewport,
// accounting for frozen columns and horizontal scroll offset.
func (m Model) visibleColumns() []Column {
	if m.viewportWidth == 0 || m.computeTotalWidth() <= m.viewportWidth {
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

	return visible
}
