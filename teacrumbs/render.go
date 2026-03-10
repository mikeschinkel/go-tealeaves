package teacrumbs

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// renderResult holds rendered content and crumb bounds for hit testing.
type renderResult struct {
	content string
	bounds  []crumbBound
}

// render renders the breadcrumb crumbs, orchestrating the pipeline:
// 1. Try default text (last=Text, others=Short) → if fits, return
// 2. Try compact (all=Short) → if fits, return
// 3. Fall through to truncate for ellipsis
func (m BreadcrumbsModel) render() renderResult {
	if len(m.crumbs) == 0 {
		return renderResult{}
	}

	// Try default: last=Text, others=Short (fallback to Text)
	result := m.renderParts(false)
	if m.width <= 0 || ansi.StringWidth(result.content) <= m.width {
		return result
	}

	// Try compact: all=Short (fallback to Text)
	compact := m.renderParts(true)
	if ansi.StringWidth(compact.content) <= m.width {
		return compact
	}

	// Truncate with ellipsis
	return m.truncate(true)
}

// renderParts renders all crumbs, joins them with separator, and computes bounds.
// When compact is true, all crumbs use Short text. Otherwise, last crumb uses Text.
func (m BreadcrumbsModel) renderParts(compact bool) renderResult {
	parts := make([]string, len(m.crumbs))
	for i := range m.crumbs {
		parts[i] = m.renderCrumb(i, compact)
	}

	styledSep := m.Styles.SeparatorStyle.Render(m.separator)
	rendered := strings.Join(parts, styledSep)
	bounds := computeBounds(parts, styledSep)

	return renderResult{content: rendered, bounds: bounds}
}

// truncate produces a "first > ... > last" layout when crumbs exceed width.
func (m BreadcrumbsModel) truncate(compact bool) renderResult {
	styledSep := m.Styles.SeparatorStyle.Render(m.separator)

	if len(m.crumbs) <= 2 {
		if len(m.crumbs) == 1 {
			content := m.renderCrumb(0, compact)
			bounds := computeBounds([]string{content}, styledSep)
			return renderResult{content: content, bounds: bounds}
		}
		first := m.renderCrumb(0, compact)
		last := m.renderCrumb(1, compact)
		content := first + styledSep + last
		bounds := computeBounds([]string{first, last}, styledSep)
		return renderResult{content: content, bounds: bounds}
	}

	// Build: first > ... > last
	first := m.renderCrumb(0, compact)
	lastIdx := len(m.crumbs) - 1
	last := m.renderCrumb(lastIdx, compact)
	ellipsis := m.Styles.SeparatorStyle.Render("...")

	result := first + styledSep + ellipsis + styledSep + last

	firstText := m.resolveText(0, compact)
	lastText := m.resolveText(lastIdx, compact)

	// Check if even this is too long
	plainWidth := ansi.StringWidth(firstText) +
		ansi.StringWidth(m.separator) +
		3 + // "..."
		ansi.StringWidth(m.separator) +
		ansi.StringWidth(lastText)

	if plainWidth <= m.width {
		bounds := computeTruncatedBounds(first, last, styledSep, ellipsis)
		return renderResult{content: result, bounds: bounds}
	}

	// Progressive truncation: shorten the last crumb text
	overhead := ansi.StringWidth(firstText) +
		ansi.StringWidth(m.separator) +
		3 + // "..."
		ansi.StringWidth(m.separator)
	available := m.width - overhead
	if available > 3 {
		truncated := ansi.Truncate(lastText, available-3, "")
		last = m.Styles.CurrentStyle.Render(truncated + "...")
		content := first + styledSep + ellipsis + styledSep + last
		bounds := computeTruncatedBounds(first, last, styledSep, ellipsis)
		return renderResult{content: content, bounds: bounds}
	}

	// Extreme truncation: just show first crumb truncated
	if m.width > 3 {
		truncated := ansi.Truncate(firstText, m.width-3, "")
		content := m.Styles.ParentStyle.Render(truncated + "...")
		bounds := []crumbBound{{startX: 0, endX: ansi.StringWidth(content)}}
		return renderResult{content: content, bounds: bounds}
	}

	content := m.Styles.ParentStyle.Render("...")
	return renderResult{content: content, bounds: nil}
}

// renderCrumb renders a single crumb including text resolution, styling, and hover.
func (m BreadcrumbsModel) renderCrumb(index int, compact bool) string {
	crumb := m.crumbs[index]

	// Custom renderer bypasses everything
	if crumb.Renderer != nil {
		return crumb.Renderer.Render(index, m)
	}

	text := m.resolveText(index, compact)

	// Hover style
	if index == m.mouse.hoveredIdx {
		return m.applyHoverStyle(index, text)
	}

	// Normal styling
	return m.styleCrumb(index, text)
}

// resolveText picks Text vs Short for the given crumb.
// When compact is true, always uses Short (falling back to Text).
// When compact is false, last crumb uses Text, others use Short (falling back to Text).
func (m BreadcrumbsModel) resolveText(index int, compact bool) string {
	crumb := m.crumbs[index]
	isCurrent := index == len(m.crumbs)-1

	if !compact && isCurrent {
		return crumb.Text
	}
	if crumb.Short != "" {
		return crumb.Short
	}
	return crumb.Text
}

// styleCrumb applies per-crumb, parent, or current style to text.
func (m BreadcrumbsModel) styleCrumb(index int, text string) string {
	crumb := m.crumbs[index]
	if crumb.Style != nil {
		return crumb.Style.Render(text)
	}
	if index == len(m.crumbs)-1 {
		return m.Styles.CurrentStyle.Render(text)
	}
	return m.Styles.ParentStyle.Render(text)
}

// applyHoverStyle applies hover styling to a crumb's text.
func (m BreadcrumbsModel) applyHoverStyle(index int, text string) string {
	crumb := m.crumbs[index]
	if crumb.Style != nil {
		return crumb.Style.Underline(true).Render(text)
	}
	return m.Styles.HoverStyle.Render(text)
}

// computeBounds calculates crumb bounds from rendered parts and separator.
func computeBounds(parts []string, styledSep string) []crumbBound {
	if len(parts) == 0 {
		return nil
	}
	bounds := make([]crumbBound, len(parts))
	sepWidth := ansi.StringWidth(styledSep)
	x := 0
	for i, part := range parts {
		w := ansi.StringWidth(part)
		bounds[i] = crumbBound{startX: x, endX: x + w}
		x += w
		if i < len(parts)-1 {
			x += sepWidth
		}
	}
	return bounds
}

// computeTruncatedBounds computes bounds for "first > ... > last" layout.
// Only first and last crumbs get bounds; the ellipsis is not clickable.
func computeTruncatedBounds(first, last, styledSep, ellipsis string) []crumbBound {
	sepWidth := ansi.StringWidth(styledSep)
	firstWidth := ansi.StringWidth(first)
	ellipsisWidth := ansi.StringWidth(ellipsis)
	lastStart := firstWidth + sepWidth + ellipsisWidth + sepWidth
	lastWidth := ansi.StringWidth(last)
	return []crumbBound{
		{startX: 0, endX: firstWidth},
		{startX: lastStart, endX: lastStart + lastWidth},
	}
}
