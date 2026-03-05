package teagrid

// recalculateWidth resolves flex column widths and scroll boundaries.
func (m *GridModel) recalculateWidth() {
	targetWidth := m.viewportWidth
	if targetWidth > 0 {
		updateColumnWidths(m.columns, targetWidth, m.border)
	}
	m.recalculateLastHorizontalColumn()
}

// recalculatePageSize computes page size from viewport height minus chrome.
func (m *GridModel) recalculatePageSize() {
	if m.viewportHeight == 0 {
		return
	}

	chrome := m.chromeHeight()
	autoPageSize := m.viewportHeight - chrome

	if autoPageSize < 1 {
		autoPageSize = 1
	}

	m.pageSize = autoPageSize
}

// chromeHeight estimates the non-data height consumed by borders,
// header, and footer. Used for auto page size computation.
func (m *GridModel) chromeHeight() int {
	height := 0

	if m.border.HasOuterBorder() {
		height += 2 // top and bottom border lines
	}

	if m.headerVisible {
		height++ // header row
		if m.border.HasHeaderSeparator() {
			height++ // header separator line
		}
	}

	if m.footerVisible {
		height++ // footer row
		if m.border.HasFooterSeparator() {
			height++ // footer separator line
		}
	}

	return height
}

// computeNaturalWidth returns the minimum width needed to display all columns
// without flex expansion. Flex columns contribute their padding + 1 char minimum.
func (m *GridModel) computeNaturalWidth() int {
	width := m.border.OuterWidth()

	for i, col := range m.columns {
		if col.IsFlex() {
			width += col.paddingLeft + 1 + col.paddingRight
		} else {
			width += col.RenderWidth()
		}

		if i < len(m.columns)-1 {
			width += m.border.InnerDividerWidth()
		}
	}

	return width
}

// computeTotalWidth returns the total rendered width of the grid
// after flex column resolution.
func (m *GridModel) computeTotalWidth() int {
	width := m.border.OuterWidth()

	for i, col := range m.columns {
		width += col.RenderWidth()

		if i < len(m.columns)-1 {
			width += m.border.InnerDividerWidth()
		}
	}

	return width
}

// updateColumnWidths resolves flex column widths to fill totalWidth.
func updateColumnWidths(cols []Column, totalWidth int, border BorderConfig) {
	if totalWidth == 0 || len(cols) == 0 {
		return
	}

	// Compute border overhead
	borderWidth := border.OuterWidth()
	if len(cols) > 1 {
		borderWidth += (len(cols) - 1) * border.InnerDividerWidth()
	}

	availableForFlex := totalWidth - borderWidth
	totalFlexFactor := 0
	flexGCD := 0

	for _, col := range cols {
		if !col.IsFlex() {
			availableForFlex -= col.RenderWidth()
		} else {
			// Padding is always present; only the content width is flexible
			availableForFlex -= col.paddingLeft + col.paddingRight
			totalFlexFactor += col.flexFactor
			flexGCD = gcd(flexGCD, col.flexFactor)
		}
	}

	if totalFlexFactor == 0 || availableForFlex <= 0 {
		return
	}

	totalFlexFactor /= flexGCD
	flexUnit := availableForFlex / totalFlexFactor
	leftoverWidth := availableForFlex % totalFlexFactor

	for i := range cols {
		if !cols[i].IsFlex() {
			continue
		}

		width := flexUnit * (cols[i].flexFactor / flexGCD)

		if leftoverWidth > 0 {
			width++
			leftoverWidth--
		}

		if i == len(cols)-1 {
			width += leftoverWidth
			leftoverWidth = 0
		}

		cols[i].width = max(width, 1)
	}
}
