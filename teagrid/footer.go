package teagrid

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// renderFooter renders the two-zone footer (filter left, pagination right).
// The footer style is independent from baseStyle (fixes v0.1.0 #3, #9).
// A hidden footer produces zero height (fixes v0.1.0 hidden footer space).
func (m GridModel) renderFooter() string {
	if !m.hasFooter() {
		return ""
	}

	bc := m.border
	totalWidth := m.computeTotalWidth()
	if totalWidth == 0 {
		totalWidth = m.computeNaturalWidth()
	}

	var parts []string

	// Footer separator line
	if bc.HasFooterSeparator() {
		parts = append(parts, m.renderFooterSeparator())
	}

	// Compute content width (inside borders)
	contentWidth := totalWidth - bc.OuterWidth()
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Build left zone (filter) and right zone (pagination)
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
	}

	// Compose the footer content
	footerContent := m.composeFooterZones(leftZone, rightZone, contentWidth)

	// Apply footer style (independent from baseStyle)
	footerStyle := m.footerStyle.Width(contentWidth)
	styledFooter := footerStyle.Render(footerContent)

	// Add outer borders
	footerRow := ""
	if bc.HasOuterBorder() {
		footerRow = bc.Outer.Style.Render(bc.Chars.Vertical) +
			styledFooter +
			bc.Outer.Style.Render(bc.Chars.Vertical)
	} else {
		footerRow = styledFooter
	}

	parts = append(parts, footerRow)
	return strings.Join(parts, "\n")
}

// composeFooterZones arranges the left and right zones within the given width.
func (m GridModel) composeFooterZones(left, right string, width int) string {
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)

	if left == "" && right == "" {
		return strings.Repeat(" ", width)
	}

	if left == "" {
		// Right-align the right zone
		padding := width - rightWidth
		if padding < 0 {
			padding = 0
		}
		return strings.Repeat(" ", padding) + right
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

// renderFooterSeparator renders the line between data rows and footer.
func (m GridModel) renderFooterSeparator() string {
	bc := m.border
	chars := bc.Chars
	style := bc.Footer.Style
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
