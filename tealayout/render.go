package tealayout

import (
	"strings"

	lipgloss "charm.land/lipgloss/v2"
)

// Render resolves the row layout, calls SetSize on children that implement
// SetSizer, calls View on children that implement Viewer, and composes
// the output by joining horizontally with gap spacing.
func (r *row) Render() (string, error) {
	if r.dirty || !r.resolved {
		r.resolved = false // force re-resolve
	}
	if r.cachedOutput != "" && !r.dirty {
		return r.cachedOutput, nil
	}

	sizes, err := r.Resolve()
	if err != nil {
		return "", err
	}

	// Phase 5: SetSize on children, then View
	views := make([]string, 0, len(r.children))
	for i, ch := range r.children {
		if sizes[i] <= 0 {
			continue
		}
		setChildSize(ch.Widget, sizes[i], r.height)
		views = append(views, viewChild(ch.Widget, sizes[i], r.height))
	}

	output := joinHorizontal(views, r.gap, r.height)
	r.cachedOutput = output
	r.dirty = false
	return output, nil
}

// Render resolves the column layout, calls SetSize on children, and composes
// output by joining vertically with gap spacing.
func (c *column) Render() (string, error) {
	if c.dirty || !c.resolved {
		c.resolved = false
	}
	if c.cachedOutput != "" && !c.dirty {
		return c.cachedOutput, nil
	}

	sizes, err := c.Resolve()
	if err != nil {
		return "", err
	}

	views := make([]string, 0, len(c.children))
	for i, ch := range c.children {
		if sizes[i] <= 0 {
			continue
		}
		setChildSize(ch.Widget, c.width, sizes[i])
		views = append(views, viewChild(ch.Widget, c.width, sizes[i]))
	}

	output := joinVertical(views, c.gap)
	c.cachedOutput = output
	c.dirty = false
	return output, nil
}

// setChildSize calls SetSize on a widget with content dimensions. If the
// widget implements Styler, border+padding are subtracted first.
func setChildSize(widget any, totalW, totalH int) {
	ss, ok := widget.(SetSizer)
	if !ok {
		return
	}
	contentW, contentH := totalW, totalH
	if styler, ok := widget.(Styler); ok {
		style := styler.Style()
		contentW -= style.GetHorizontalPadding() + style.GetHorizontalBorderSize()
		contentH -= style.GetVerticalPadding() + style.GetVerticalBorderSize()
		if contentW < 0 {
			contentW = 0
		}
		if contentH < 0 {
			contentH = 0
		}
	}
	ss.SetSize(contentW, contentH)
}

// viewChild renders a widget. If it implements Viewer, View() is called.
// Otherwise returns empty space of the given dimensions.
func viewChild(widget any, width, height int) string {
	if v, ok := widget.(Viewer); ok {
		return v.View()
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
