package teagrid

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// renderHeaders renders the header row with border characters inserted as
// literal strings rather than via lipgloss border styles.
func (m GridModel) renderHeaders() string {
	if !m.headerVisible {
		return ""
	}

	var parts []string

	bc := m.border
	outerBorder := bc.HasOuterBorder()
	innerDivider := bc.HasInnerDividers()
	headerSep := bc.HasHeaderSeparator()

	// Top border line
	if outerBorder {
		parts = append(parts, m.renderTopBorder())
	}

	// Header content
	headerCells := make([]string, 0, len(m.columns))
	for _, col := range m.columns {
		title := col.title
		if col.fmtString != "" {
			title = fmt.Sprintf(col.fmtString, title)
		}
		title = limitStr(title, col.width)

		// Apply padding
		padded := strings.Repeat(" ", col.paddingLeft) +
			padOrTruncate(title, col.width, col.alignment) +
			strings.Repeat(" ", col.paddingRight)

		cellStyle := m.headerStyle.Inherit(m.baseStyle)
		headerCells = append(headerCells, cellStyle.Render(padded))
	}

	divider := ""
	if innerDivider {
		divStyle := bc.Inner.Style
		divider = divStyle.Render(bc.Chars.InnerDivider)
	}

	headerRow := ""
	if outerBorder {
		headerRow += bc.Outer.Style.Render(bc.Chars.Vertical)
	}
	headerRow += strings.Join(headerCells, divider)
	if outerBorder {
		headerRow += bc.Outer.Style.Render(bc.Chars.Vertical)
	}

	parts = append(parts, headerRow)

	// Header separator line
	if headerSep {
		parts = append(parts, m.renderHeaderSeparator())
	}

	return strings.Join(parts, "\n")
}

// renderTopBorder renders the top border line with junctions.
func (m GridModel) renderTopBorder() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Outer.Style

	var line strings.Builder
	line.WriteString(style.Render(chars.TopLeft))

	for i, col := range m.columns {
		w := col.RenderWidth()
		line.WriteString(style.Render(strings.Repeat(chars.Horizontal, w)))

		if i < len(m.columns)-1 {
			if bc.HasInnerDividers() {
				line.WriteString(style.Render(chars.TopJunction))
			} else {
				line.WriteString(style.Render(chars.Horizontal))
			}
		}
	}

	line.WriteString(style.Render(chars.TopRight))
	return line.String()
}

// renderHeaderSeparator renders the line between header and data rows.
func (m GridModel) renderHeaderSeparator() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Header.Style
	outerStyle := bc.Outer.Style

	var line strings.Builder

	if bc.HasOuterBorder() {
		line.WriteString(outerStyle.Render(chars.LeftJunction))
	}

	for i, col := range m.columns {
		w := col.RenderWidth()
		line.WriteString(style.Render(strings.Repeat(chars.Horizontal, w)))

		if i < len(m.columns)-1 {
			if bc.HasInnerDividers() {
				line.WriteString(style.Render(chars.InnerJunction))
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

// padOrTruncate pads or truncates text to exactly the given width,
// respecting alignment.
func padOrTruncate(text string, width int, align lipgloss.Position) string {
	textWidth := lipgloss.Width(text)

	if textWidth > width {
		return limitStr(text, width)
	}

	padding := width - textWidth
	if padding == 0 {
		return text
	}

	switch align {
	case lipgloss.Right:
		return strings.Repeat(" ", padding) + text
	case lipgloss.Center:
		left := padding / 2
		right := padding - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	default: // lipgloss.Left
		return text + strings.Repeat(" ", padding)
	}
}
