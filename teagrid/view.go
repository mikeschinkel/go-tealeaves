package teagrid

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// View renders the grid as a tea.View (Charm v2).
func (m GridModel) View() tea.View {
	return tea.NewView(m.render())
}

// render produces the complete grid string.
func (m GridModel) render() string {
	m = m.ensureVisibleRowsCached()
	m = m.ensureVisibleColumnsCached()

	var sections []string

	// Header
	header := m.renderHeaders()
	if header != "" {
		sections = append(sections, header)
	}

	// Data rows
	rows := m.renderDataRows()
	if rows != "" {
		sections = append(sections, rows)
	}

	// Bottom border (before footer if footer exists)
	if m.border.HasOuterBorder() && !m.hasFooter() {
		sections = append(sections, m.renderBottomBorder())
	}

	// Footer
	footer := m.renderFooter()
	if footer != "" {
		sections = append(sections, footer)
	}

	// Bottom border (after footer if footer exists)
	if m.border.HasOuterBorder() && m.hasFooter() {
		if m.bottomBorderNeedsJunctions() {
			sections = append(sections, m.renderBottomBorder())
		} else {
			sections = append(sections, m.renderPlainBottomBorder())
		}
	}

	return strings.Join(sections, "\n")
}

// renderDataRows renders all visible data rows for the current page.
func (m GridModel) renderDataRows() string {
	visibleRows := m.cachedVisibleRows
	if len(visibleRows) == 0 {
		return ""
	}

	start, end := m.VisibleIndices()
	if start > len(visibleRows)-1 {
		return ""
	}

	var lines []string
	for i := start; i <= end && i < len(visibleRows); i++ {
		line := m.renderRow(visibleRows[i], i)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m GridModel) leftVertical() string {
	bc := m.border
	if m.overflowIndicator && m.hasHiddenColumnsLeft() &&
		bc.Chars.OverflowVertical != "" && m.horizontalScrollFreezeColumnsCount == 0 {
		return bc.Chars.OverflowVertical
	}
	return bc.Chars.Vertical
}

func (m GridModel) rightVertical() string {
	bc := m.border
	if m.overflowIndicator && m.hasHiddenColumnsRight() && bc.Chars.OverflowVertical != "" {
		return bc.Chars.OverflowVertical
	}
	return bc.Chars.Vertical
}

// renderRow renders a single data row with borders and styling.
// Cursor/highlight are applied at render time — no row rebuilding needed.
func (m GridModel) renderRow(row Row, rowIndex int) string {
	bc := m.border
	outerBorder := bc.HasOuterBorder()

	isHighlightedRow := m.focused && rowIndex == m.rowCursorIndex

	visible := m.visibleColumns()
	cells := make([]string, 0, len(visible))

	// Resolve cursor column key for key-based matching
	cursorKey := ""
	if m.colCursorColumnIndex >= 0 && m.colCursorColumnIndex < len(m.columns) {
		cursorKey = m.columns[m.colCursorColumnIndex].key
	}

	for _, col := range visible {
		origIndex := m.columnOriginalIndex(col.key)
		isColCursor := isHighlightedRow && m.colCursorMode &&
			col.key == cursorKey
		cellStr := m.renderCell(row, col, rowIndex, origIndex, isHighlightedRow, isColCursor)
		cells = append(cells, cellStr)
	}

	var line strings.Builder
	if outerBorder {
		line.WriteString(bc.Outer.Style.Render(m.leftVertical()))
	}
	line.WriteString(m.joinCellsWithDividers(cells, true))
	if outerBorder {
		line.WriteString(bc.Outer.Style.Render(m.rightVertical()))
	}

	return line.String()
}

// renderCell renders a single cell value with padding, alignment, and cursor styling.
//
// Style cascade: base → column → row → rowStyleFunc → cellValue style
// → row highlight overlay → column cursor overlay
func (m GridModel) renderCell(row Row, col Column, rowIndex, colIndex int, isHighlightedRow, isColCursor bool) string {
	// 1. base → column → row
	cellStyle := col.style.Inherit(m.baseStyle)
	cellStyle = row.Style.Inherit(cellStyle)

	// 2. rowStyleFunc
	if m.rowStyleFunc != nil {
		cellStyle = m.rowStyleFunc(RowStyleFuncInput{
			Index:         rowIndex,
			Row:           row,
			IsHighlighted: isHighlightedRow,
		}).Inherit(cellStyle)
	}

	var displayStr string

	// Handle select column
	if col.key == selectColumnKey {
		if row.selected {
			displayStr = m.selectedText
		} else {
			displayStr = m.unselectedText
		}
	} else if data, exists := row.Data[col.key]; !exists {
		if m.missingDataIndicator != nil {
			displayStr = fmt.Sprintf("%v", m.missingDataIndicator)
		}
	} else {
		fmtString := "%v"
		if col.fmtString != "" {
			fmtString = col.fmtString
		}

		switch v := data.(type) {
		case CellValue:
			if v.HasSpans() {
				displayStr = m.renderSpans(v.Spans, col.width)
			} else {
				displayStr = fmt.Sprintf(fmtString, v.Data)
			}

			// 3. cellValue style/styleFunc
			if v.StyleFunc != nil {
				cellStyle = v.StyleFunc(CellStyleInput{
					Data:             v.Data,
					Column:           col,
					Row:              row,
					RowIndex:         rowIndex,
					ColumnIndex:      colIndex,
					IsHighlightedRow: isHighlightedRow,
					IsColCursor:      isColCursor,
					GlobalMetadata:   m.metadata,
				}).Inherit(cellStyle)
			} else {
				cellStyle = v.Style.Inherit(cellStyle)
			}
		default:
			displayStr = fmt.Sprintf(fmtString, data)
		}
	}

	// 4. Row highlight overlay (only when no rowStyleFunc, which handles it itself)
	if isHighlightedRow && m.rowStyleFunc == nil {
		cellStyle = m.highlightStyle.Inherit(cellStyle)
	}

	// 5. Column cursor overlay (strongest)
	if isColCursor {
		cellStyle = m.colCursorStyle.Inherit(cellStyle)
	}

	// Truncate to content width
	displayStr = limitStr(displayStr, col.width)

	// Apply padding and alignment
	padded := strings.Repeat(" ", col.paddingLeft) +
		padOrTruncate(displayStr, col.width, col.alignment) +
		strings.Repeat(" ", col.paddingRight)

	return cellStyle.Render(padded)
}

// renderSpans renders rich text spans, truncating to fit width.
func (m GridModel) renderSpans(spans []Span, maxWidth int) string {
	var result strings.Builder
	remaining := maxWidth

	for _, span := range spans {
		if remaining <= 0 {
			break
		}

		text := span.Text
		textWidth := lipgloss.Width(text)

		if textWidth > remaining {
			text = limitStr(text, remaining)
			textWidth = lipgloss.Width(text)
		}

		result.WriteString(span.Style.Render(text))

		remaining -= textWidth
	}

	return result.String()
}

// bottomBorderNeedsJunctions returns true when the content directly above
// the bottom border has column dividers (data rows or footer rows).
// Returns false when the last section is full-width (info row).
func (m GridModel) bottomBorderNeedsJunctions() bool {
	if !m.hasFooter() {
		// No footer — data rows are directly above, they have column dividers
		return true
	}
	// If the last section in the footer is the info row (full-width), no junctions.
	// If the last section is a footer row (column-aware), junctions needed.
	if m.hasInfoRow() {
		return false
	}
	return len(m.footerRows) > 0
}

// renderPlainBottomBorder renders the bottom border with no inner junctions.
// Used when the content above is full-width (e.g., info row).
func (m GridModel) renderPlainBottomBorder() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Outer.Style

	totalWidth := m.computeVisibleWidth()
	if totalWidth == 0 {
		totalWidth = m.computeNaturalWidth()
	}

	contentWidth := totalWidth - bc.OuterWidth()
	if contentWidth < 1 {
		contentWidth = 1
	}

	var line strings.Builder
	line.WriteString(style.Render(chars.BottomLeft))
	line.WriteString(style.Render(strings.Repeat(chars.Horizontal, contentWidth)))
	line.WriteString(style.Render(chars.BottomRight))
	return line.String()
}

// renderBottomBorder renders the bottom border line with column junctions.
func (m GridModel) renderBottomBorder() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Outer.Style

	visible := m.visibleColumns()

	var line strings.Builder
	line.WriteString(style.Render(chars.BottomLeft))

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

	line.WriteString(style.Render(chars.BottomRight))
	return line.String()
}

// joinCellsWithDividers joins rendered cells with inner dividers, using the
// freeze divider character at the boundary between frozen and scrollable columns.
// When dataRow is true and overflow is active, the freeze divider uses the
// overflow character to indicate hidden columns to the left.
func (m GridModel) joinCellsWithDividers(cells []string, dataRow bool) string {
	bc := m.border
	if !bc.HasInnerDividers() || len(cells) <= 1 {
		return strings.Join(cells, "")
	}

	divider := bc.Inner.Style.Render(bc.Chars.InnerDivider)
	freezeIdx := m.horizontalScrollFreezeColumnsCount - 1

	if freezeIdx < 0 || freezeIdx >= len(cells)-1 || bc.Chars.FreezeDivider == "" {
		return strings.Join(cells, divider)
	}

	freezeChar := bc.Chars.FreezeDivider
	if dataRow && m.overflowIndicator && m.hasHiddenColumnsLeft() && bc.Chars.OverflowVertical != "" {
		freezeChar = bc.Chars.OverflowVertical
	}
	freezeDiv := bc.Inner.Style.Render(freezeChar)
	var sb strings.Builder
	for i, cell := range cells {
		sb.WriteString(cell)
		if i < len(cells)-1 {
			if i == freezeIdx {
				sb.WriteString(freezeDiv)
			} else {
				sb.WriteString(divider)
			}
		}
	}
	return sb.String()
}
