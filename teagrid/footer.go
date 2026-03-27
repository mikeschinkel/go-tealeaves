package teagrid

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// FooterCell represents a single cell within a column-aware footer row.
type FooterCell struct {
	ColumnKey string          // anchor to column by key; "" = next available
	ColSpan   int             // columns to span (default 1)
	Value     string          // display text
	Alignment lipgloss.Position
	Style     lipgloss.Style
}

// FooterRow represents a column-aware footer row containing cells that
// can span multiple columns, similar to HTML <tfoot> with colspan.
type FooterRow struct {
	Cells []FooterCell
	Style lipgloss.Style
}

// NewFooterCell creates a footer cell anchored to the given column key.
func NewFooterCell(key, value string) FooterCell {
	return FooterCell{
		ColumnKey: key,
		ColSpan:   1,
		Value:     value,
		Alignment: lipgloss.Left,
	}
}

// NewFooterCellSpan creates a footer cell spanning multiple columns.
func NewFooterCellSpan(key, value string, span int) FooterCell {
	if span < 1 {
		span = 1
	}
	return FooterCell{
		ColumnKey: key,
		ColSpan:   span,
		Value:     value,
		Alignment: lipgloss.Left,
	}
}

// NewFooterRow creates a footer row from the given cells.
func NewFooterRow(cells ...FooterCell) FooterRow {
	return FooterRow{Cells: cells}
}

// WithAlignment returns a copy of the cell with the given alignment.
func (c FooterCell) WithAlignment(a lipgloss.Position) FooterCell {
	c.Alignment = a
	return c
}

// WithStyle returns a copy of the cell with the given style.
func (c FooterCell) WithStyle(s lipgloss.Style) FooterCell {
	c.Style = s
	return c
}

// WithStyle returns a copy of the row with the given style.
func (r FooterRow) WithStyle(s lipgloss.Style) FooterRow {
	r.Style = s
	return r
}

// renderFooter renders the complete footer section:
// 1. Footer separator (column junctions if footer rows follow, plain if not)
// 2. Column-aware footer rows (if any)
// 3. Plain separator between footer rows and info row (if both exist)
// 4. Full-width info row (static text + pagination + filter)
func (m GridModel) renderFooter() string {
	if !m.hasFooter() {
		return ""
	}

	bc := m.border
	hasRows := len(m.footerRows) > 0
	hasInfo := m.hasInfoRow()

	var parts []string

	// Footer separator line
	if bc.HasFooterSeparator() {
		if hasRows {
			// Column junctions connect to footer row dividers below
			parts = append(parts, m.renderFooterSeparator())
		} else {
			// Plain separator — no column junctions orphaned
			parts = append(parts, m.renderPlainFooterSeparator())
		}
	}

	// Column-aware footer rows
	for _, row := range m.footerRows {
		parts = append(parts, m.renderFooterRow(row))
	}

	// Separator between footer rows and info row
	if hasRows && hasInfo && bc.HasFooterSeparator() {
		parts = append(parts, m.renderPlainFooterSeparator())
	}

	// Info row (static text + pagination + filter)
	if hasInfo {
		parts = append(parts, m.renderInfoRow())
	}

	return strings.Join(parts, "\n")
}

// renderInfoRow renders the full-width info row with filter, static text, and pagination.
func (m GridModel) renderInfoRow() string {
	bc := m.border
	totalWidth := m.computeVisibleWidth()
	if totalWidth == 0 {
		totalWidth = m.computeNaturalWidth()
	}

	contentWidth := totalWidth - bc.OuterWidth()
	if contentWidth < 1 {
		contentWidth = 1
	}

	leftZone := ""
	rightZone := ""

	if m.staticFooter != "" {
		rightZone = m.staticFooter
	}

	if m.filtered && m.filterTextInput.Value() != "" {
		leftZone = m.filterTextInput.View()
	} else if m.filtered && m.filterTextInput.Focused() {
		leftZone = m.filterTextInput.View()
	}

	if m.pageSize > 0 {
		pagination := fmt.Sprintf("%d/%d", m.CurrentPage(), m.MaxPages())
		if rightZone != "" {
			rightZone = rightZone + " " + pagination
		} else {
			rightZone = pagination
		}
		// Override alignment for this render — value receiver, original unchanged
		m.footerAlignment = lipgloss.Right
	}

	// Apply cell-matching padding (1 space on each side)
	paddedWidth := contentWidth - 2*defaultPaddingRight
	if paddedWidth < 1 {
		paddedWidth = 1
	}
	footerContent := " " + m.composeFooterZones(leftZone, rightZone, paddedWidth) + " "

	footerStyle := m.footerStyle.Width(contentWidth)
	styledFooter := footerStyle.Render(footerContent)

	if bc.HasOuterBorder() {
		return bc.Outer.Style.Render(bc.Chars.Vertical) +
			styledFooter +
			bc.Outer.Style.Render(bc.Chars.Vertical)
	}
	return styledFooter
}

// composeFooterZones arranges the left and right zones within the given width.
// When only the right zone is set (no filter active), footerAlignment controls placement.
func (m GridModel) composeFooterZones(left, right string, width int) string {
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)

	if left == "" && right == "" {
		return strings.Repeat(" ", width)
	}

	if left == "" {
		// Use footerAlignment to position the right zone
		padding := width - rightWidth
		if padding < 0 {
			padding = 0
		}
		switch m.footerAlignment {
		case lipgloss.Left:
			return right + strings.Repeat(" ", padding)
		case lipgloss.Center:
			leftPad := padding / 2
			rightPad := padding - leftPad
			return strings.Repeat(" ", leftPad) + right + strings.Repeat(" ", rightPad)
		default: // lipgloss.Right
			return strings.Repeat(" ", padding) + right
		}
	}

	if right == "" {
		// Left-align the left zone
		padding := width - leftWidth
		if padding < 0 {
			padding = 0
		}
		return left + strings.Repeat(" ", padding)
	}

	// Both zones: left-align left, right-align right
	gap := width - leftWidth - rightWidth
	if gap < 1 {
		gap = 1
	}

	return left + strings.Repeat(" ", gap) + right
}

// renderFooterSeparator renders the line between data rows and footer
// with column junctions matching the inner dividers above.
func (m GridModel) renderFooterSeparator() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Footer.Style
	outerStyle := bc.Outer.Style

	var line strings.Builder

	if bc.HasOuterBorder() {
		line.WriteString(outerStyle.Render(chars.LeftJunction))
	}

	visible := m.visibleColumns()
	freezeIdx := m.horizontalScrollFreezeColumnsCount - 1

	for i, col := range visible {
		w := col.RenderWidth()
		line.WriteString(style.Render(strings.Repeat(chars.Horizontal, w)))

		if i < len(visible)-1 {
			if bc.HasInnerDividers() {
				junction := chars.InnerJunction
				if i == freezeIdx && chars.FreezeInnerJunction != "" {
					junction = chars.FreezeInnerJunction
				}
				line.WriteString(style.Render(junction))
			} else {
				line.WriteString(style.Render(chars.Horizontal))
			}
		}
	}

	if bc.HasOuterBorder() {
		line.WriteString(outerStyle.Render(chars.RightJunction))
	}

	return line.String()
}

// renderPlainFooterSeparator renders a horizontal line with upward-connecting
// junctions (BottomJunction) at column positions. The junctions connect UP to
// the column dividers in the content above but are flat on the bottom since the
// content below (info row) is full-width with no column dividers.
func (m GridModel) renderPlainFooterSeparator() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Footer.Style
	outerStyle := bc.Outer.Style

	var line strings.Builder

	if bc.HasOuterBorder() {
		line.WriteString(outerStyle.Render(chars.LeftJunction))
	}

	visible := m.visibleColumns()
	freezeIdx := m.horizontalScrollFreezeColumnsCount - 1

	for i, col := range visible {
		w := col.RenderWidth()
		line.WriteString(style.Render(strings.Repeat(chars.Horizontal, w)))

		if i < len(visible)-1 {
			if bc.HasInnerDividers() {
				junction := chars.BottomJunction
				if i == freezeIdx && chars.FreezeBottomJunction != "" {
					junction = chars.FreezeBottomJunction
				}
				line.WriteString(style.Render(junction))
			} else {
				line.WriteString(style.Render(chars.Horizontal))
			}
		}
	}

	if bc.HasOuterBorder() {
		line.WriteString(outerStyle.Render(chars.RightJunction))
	}

	return line.String()
}

// renderFooterRow renders a single column-aware footer row.
// Cells are mapped to visible columns; spans merge multiple column widths.
func (m GridModel) renderFooterRow(row FooterRow) string {
	bc := m.border
	visible := m.visibleColumns()
	if len(visible) == 0 {
		return ""
	}

	// Build a key→visibleIndex map
	keyToIdx := make(map[string]int, len(visible))
	for i, col := range visible {
		keyToIdx[col.key] = i
	}

	// Track which visible columns are consumed
	consumed := make([]bool, len(visible))

	// Segment: rendered content for a contiguous span of columns
	type segment struct {
		startIdx int
		endIdx   int // inclusive
		content  string
	}
	var segments []segment

	for _, cell := range row.Cells {
		startIdx, ok := keyToIdx[cell.ColumnKey]
		if !ok {
			continue
		}

		span := cell.ColSpan
		if span < 1 {
			span = 1
		}
		endIdx := startIdx + span - 1
		if endIdx >= len(visible) {
			endIdx = len(visible) - 1
		}

		// Calculate total width for this span
		spanWidth := 0
		for j := startIdx; j <= endIdx; j++ {
			spanWidth += visible[j].RenderWidth()
			if j > startIdx {
				spanWidth += bc.InnerDividerWidth()
			}
			consumed[j] = true
		}

		// Render cell content within span width
		cellStyle := cell.Style.Inherit(row.Style)
		content := padOrTruncate(cell.Value, spanWidth, cell.Alignment)
		segments = append(segments, segment{
			startIdx: startIdx,
			endIdx:   endIdx,
			content:  cellStyle.Render(content),
		})
	}

	// Fill gaps with empty segments
	var allSegments []segment
	nextCol := 0
	for _, seg := range segments {
		if seg.startIdx > nextCol {
			// Gap before this segment
			gapWidth := 0
			for j := nextCol; j < seg.startIdx; j++ {
				gapWidth += visible[j].RenderWidth()
				if j > nextCol {
					gapWidth += bc.InnerDividerWidth()
				}
			}
			allSegments = append(allSegments, segment{
				startIdx: nextCol,
				endIdx:   seg.startIdx - 1,
				content:  strings.Repeat(" ", gapWidth),
			})
		}
		allSegments = append(allSegments, seg)
		nextCol = seg.endIdx + 1
	}

	// Trailing gap
	if nextCol < len(visible) {
		gapWidth := 0
		for j := nextCol; j < len(visible); j++ {
			gapWidth += visible[j].RenderWidth()
			if j > nextCol {
				gapWidth += bc.InnerDividerWidth()
			}
		}
		allSegments = append(allSegments, segment{
			startIdx: nextCol,
			endIdx:   len(visible) - 1,
			content:  strings.Repeat(" ", gapWidth),
		})
	}

	// Join segments with dividers between them
	freezeIdx := m.horizontalScrollFreezeColumnsCount - 1
	var line strings.Builder

	if bc.HasOuterBorder() {
		line.WriteString(bc.Outer.Style.Render(bc.Chars.Vertical))
	}

	for i, seg := range allSegments {
		line.WriteString(seg.content)
		if i < len(allSegments)-1 {
			if bc.HasInnerDividers() {
				divChar := bc.Chars.InnerDivider
				if seg.endIdx == freezeIdx && bc.Chars.FreezeDivider != "" {
					divChar = bc.Chars.FreezeDivider
				}
				line.WriteString(bc.Inner.Style.Render(divChar))
			}
		}
	}

	if bc.HasOuterBorder() {
		line.WriteString(bc.Outer.Style.Render(bc.Chars.Vertical))
	}

	return line.String()
}
