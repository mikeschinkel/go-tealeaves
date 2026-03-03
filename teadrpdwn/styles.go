package teadrpdwn

import "charm.land/lipgloss/v2"

// DefaultBorderStyle returns default border styling
func DefaultBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
}

// DefaultItemStyle returns default item styling
func DefaultItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")). // Explicit white color
		Underline(false)
}

// DefaultSelectedStyle returns default selected item styling
func DefaultSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Underline(false)
}
