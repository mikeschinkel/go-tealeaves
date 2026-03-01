package teagrid

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// View renders the grid as a tea.View (Charm v2).
func (m Model) View() tea.View {
	return tea.NewView(m.render())
}

// render produces the complete grid string.
func (m Model) render() string {
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
		sections = append(sections, m.renderBottomBorder())
	}

	return strings.Join(sections, "\n")
}

// renderDataRows renders all visible data rows for the current page.
func (m Model) renderDataRows() string {
	visibleRows := m.GetVisibleRows()
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

// renderRow renders a single data row with borders and styling.
// Cursor/highlight are applied at render time — no row rebuilding needed.
func (m Model) renderRow(row Row, rowIndex int) string {
	bc := m.border
	outerBorder := bc.HasOuterBorder()
	innerDivider := bc.HasInnerDividers()

	isHighlightedRow := m.focused && rowIndex == m.rowCursorIndex

	cells := make([]string, 0, len(m.columns))

	for colIndex, col := range m.columns {
		isCursorCell := isHighlightedRow && m.cellCursorMode &&
			colIndex == m.cellCursorColumnIndex
		cellStr := m.renderCell(row, col, rowIndex, colIndex, isHighlightedRow, isCursorCell)
		cells = append(cells, cellStr)
	}

	divider := ""
	if innerDivider {
		divider = bc.Inner.Style.Render(bc.Chars.InnerDivider)
	}

	var line strings.Builder
	if outerBorder {
		line.WriteString(bc.Outer.Style.Render(bc.Chars.Vertical))
	}
	line.WriteString(strings.Join(cells, divider))
	if outerBorder {
		line.WriteString(bc.Outer.Style.Render(bc.Chars.Vertical))
	}

	return line.String()
}

// renderCell renders a single cell value with padding, alignment, and cursor styling.
//
// Style cascade: base → column → row → rowStyleFunc → cellValue style
// → row highlight overlay → cell cursor overlay
func (m Model) renderCell(row Row, col Column, rowIndex, colIndex int, isHighlightedRow, isCursorCell bool) string {
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
					IsCursorCell:     isCursorCell,
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

	// 5. Cell cursor overlay (strongest)
	if isCursorCell {
		cellStyle = m.cellCursorStyle.Inherit(cellStyle)
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
func (m Model) renderSpans(spans []Span, maxWidth int) string {
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

// renderBottomBorder renders the bottom border line.
func (m Model) renderBottomBorder() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Outer.Style

	var line strings.Builder
	line.WriteString(style.Render(chars.BottomLeft))

	for i, col := range m.columns {
		w := col.RenderWidth()
		line.WriteString(style.Render(strings.Repeat(chars.Horizontal, w)))

		if i < len(m.columns)-1 {
			if bc.HasInnerDividers() {
				line.WriteString(style.Render(chars.BottomJunction))
			} else {
				line.WriteString(style.Render(chars.Horizontal))
			}
		}
	}

	line.WriteString(style.Render(chars.BottomRight))
	return line.String()
}
