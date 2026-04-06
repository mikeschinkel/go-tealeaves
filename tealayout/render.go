package tealayout

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

// setChildSizeViaElement calls setSize on an element with content dimensions.
// If the element implements Styler (detected via style()), border+padding are
// subtracted first.
func setChildSizeViaElement(elem element, totalW, totalH int) {
	contentW, contentH := totalW, totalH
	if style, ok := elem.style(); ok {
		contentW -= style.GetHorizontalPadding() + style.GetHorizontalBorderSize()
		contentH -= style.GetVerticalPadding() + style.GetVerticalBorderSize()
		if contentW < 0 {
			contentW = 0
		}
		if contentH < 0 {
			contentH = 0
		}
	}
	elem.setSize(contentW, contentH)
}

// contentChildElement renders an element. If it has a content function,
// content() is called. Otherwise returns empty space of the given dimensions.
func contentChildElement(elem element, width, height int) string {
	c := elem.content()
	if c != "" {
		return c
	}
	return emptyBlock(width, height)
}

// emptyBlock returns a block of spaces with the given dimensions.
func emptyBlock(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	line := strings.Repeat(" ", width)
	lines := make([]string, height)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

// joinHorizontal joins rendered blocks side by side, padding shorter blocks
// to the target height, with gap spaces between them.
func joinHorizontal(views []string, gap, height int) string {
	if len(views) == 0 {
		return ""
	}

	gapStr := ""
	if gap > 0 {
		gapStr = strings.Repeat(" ", gap)
	}

	// Split each view into lines and pad to height
	allLines := make([][]string, len(views))
	for i, v := range views {
		lines := strings.Split(v, "\n")
		// Determine width from first line
		w := 0
		if len(lines) > 0 {
			w = lipgloss.Width(lines[0])
		}
		// Pad to height
		for len(lines) < height {
			lines = append(lines, strings.Repeat(" ", w))
		}
		allLines[i] = lines
	}

	// Join line by line
	result := make([]string, height)
	for r := range height {
		parts := make([]string, len(allLines))
		for col, lines := range allLines {
			if r < len(lines) {
				parts[col] = lines[r]
			}
		}
		result[r] = strings.Join(parts, gapStr)
	}

	return strings.Join(result, "\n")
}

// joinVertical joins rendered blocks vertically with gap empty lines.
func joinVertical(views []string, gap int) string {
	if len(views) == 0 {
		return ""
	}

	gapStr := ""
	if gap > 0 {
		gapStr = strings.Repeat("\n", gap)
	}

	parts := make([]string, 0, len(views)*2)
	for i, v := range views {
		if i > 0 && gapStr != "" {
			parts = append(parts, gapStr)
		}
		parts = append(parts, v)
	}
	return strings.Join(parts, "\n")
}
