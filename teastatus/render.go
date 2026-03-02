package teastatus

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// renderMenuItems renders the left-side menu items as "[key] label  [key] label  ..."
func (m Model) renderMenuItems() (rendered string) {
	if len(m.menuItems) == 0 {
		goto end
	}

	rendered = m.buildMenuString()

end:
	return rendered
}

// buildMenuString builds the formatted menu string from menu items.
func (m Model) buildMenuString() string {
	return RenderMenuLine(m.menuItems, m.Styles)
}

// RenderMenuLine renders a list of menu items as "[key] label  [key] label  ..."
// using the given styles. Useful for inline key help outside the status bar.
func RenderMenuLine(items []MenuItem, styles Styles) string {
	var parts []string
	var item MenuItem

	for _, item = range items {
		keyPart := styles.MenuKeyStyle.Render("[" + formatKey(item.Key) + "]")
		labelPart := styles.MenuLabelStyle.Render(item.Label)
		parts = append(parts, keyPart+" "+labelPart)
	}

	return strings.Join(parts, styles.MenuSeparator)
}

// renderIndicators renders the right-side indicators according to SeparatorKind.
func (m Model) renderIndicators() (rendered string) {
	if len(m.indicators) == 0 {
		goto end
	}

	switch m.Styles.SeparatorKind {
	case BracketSeparator:
		rendered = m.renderBracketIndicators()
	default:
		rendered = m.renderSeparatedIndicators()
	}

end:
	return rendered
}

// renderSeparatedIndicators renders indicators with Pipe or Space separators.
func (m Model) renderSeparatedIndicators() string {
	var parts []string
	var sep string
	var indicator StatusIndicator

	switch m.Styles.SeparatorKind {
	case PipeSeparator:
		sep = " | "
	default:
		sep = "  "
	}

	for _, indicator = range m.indicators {
		style := m.indicatorStyle(indicator)
		parts = append(parts, style.Render(indicator.Text))
	}

	styledSep := m.Styles.IndicatorSepStyle.Render(sep)
	return strings.Join(parts, styledSep)
}

// renderBracketIndicators renders indicators wrapped as "[text]".
func (m Model) renderBracketIndicators() string {
	var parts []string
	var indicator StatusIndicator

	for _, indicator = range m.indicators {
		style := m.indicatorStyle(indicator)
		parts = append(parts, style.Render("["+indicator.Text+"]"))
	}

	return strings.Join(parts, " ")
}

// indicatorStyle returns the style to use for a given indicator.
// Uses the indicator's own Style if non-zero, otherwise falls back to the default.
func (m Model) indicatorStyle(indicator StatusIndicator) lipgloss.Style {
	// Check if indicator has a custom foreground color set
	_, isNoColor := indicator.Style.GetForeground().(lipgloss.NoColor)
	if !isNoColor {
		return indicator.Style
	}
	return m.Styles.IndicatorStyle
}

// formatKey converts raw key names to display-friendly names.
func formatKey(rawKey string) string {
	switch rawKey {
	case " ":
		return "space"
	case "ctrl+c":
		return "ctrl+c"
	default:
		return rawKey
	}
}
