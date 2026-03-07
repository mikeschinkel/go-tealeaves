package teacrumbs

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// renderTrail renders the breadcrumb trail as a styled string.
// If width > 0 and the rendered trail exceeds width, it is truncated.
func renderTrail(trail []Crumb, separator string, width int, styles Styles) string {
	if len(trail) == 0 {
		return ""
	}

	parts := make([]string, len(trail))
	for i, crumb := range trail {
		switch {
		case crumb.PreStyled:
			parts[i] = crumb.Text
		case i == len(trail)-1:
			parts[i] = styles.CurrentStyle.Render(crumb.Text)
		default:
			parts[i] = styles.ParentStyle.Render(crumb.Text)
		}
	}

	styledSep := styles.SeparatorStyle.Render(separator)
	rendered := strings.Join(parts, styledSep)

	if width <= 0 {
		return rendered
	}

	// Check if truncation is needed using ANSI-aware width
	if ansi.StringWidth(rendered) <= width {
		return rendered
	}

	return truncateTrail(trail, separator, width, styles)
}

// truncateTrail shortens the breadcrumb trail to fit within width.
// Keeps first and last items, replaces middle with "...".
func truncateTrail(trail []Crumb, separator string, width int, styles Styles) string {
	styledSep := styles.SeparatorStyle.Render(separator)

	if len(trail) <= 2 {
		if len(trail) == 1 {
			return styleCrumb(trail[0], true, styles)
		}
		return styleCrumb(trail[0], false, styles) + styledSep + styleCrumb(trail[1], true, styles)
	}

	// Build: first > ... > last
	first := styleCrumb(trail[0], false, styles)
	last := styleCrumb(trail[len(trail)-1], true, styles)
	ellipsis := styles.SeparatorStyle.Render("...")

	result := first + styledSep + ellipsis + styledSep + last

	// Check if even this is too long
	plainWidth := ansi.StringWidth(trail[0].Text) +
		ansi.StringWidth(separator) +
		3 + // "..."
		ansi.StringWidth(separator) +
		ansi.StringWidth(trail[len(trail)-1].Text)

	if plainWidth <= width {
		return result
	}

	// Progressive truncation: shorten the last crumb text
	lastText := trail[len(trail)-1].Text
	overhead := ansi.StringWidth(trail[0].Text) +
		ansi.StringWidth(separator) +
		3 + // "..."
		ansi.StringWidth(separator)
	available := width - overhead
	if available > 3 {
		truncated := ansi.Truncate(lastText, available-3, "")
		last = styles.CurrentStyle.Render(truncated + "...")
		return first + styledSep + ellipsis + styledSep + last
	}

	// Extreme truncation: just show first crumb truncated
	if width > 3 {
		firstText := trail[0].Text
		truncated := ansi.Truncate(firstText, width-3, "")
		return styles.ParentStyle.Render(truncated + "...")
	}

	return styles.ParentStyle.Render("...")
}

// styleCrumb applies the appropriate style to a crumb.
func styleCrumb(crumb Crumb, isCurrent bool, styles Styles) string {
	if crumb.PreStyled {
		return crumb.Text
	}
	if isCurrent {
		return styles.CurrentStyle.Render(crumb.Text)
	}
	return styles.ParentStyle.Render(crumb.Text)
}
