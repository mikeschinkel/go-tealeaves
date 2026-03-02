package teadep

import "charm.land/lipgloss/v2"

// DefaultPathStyle returns default path item styling
func DefaultPathStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

// DefaultSelectedStyle returns default selected item styling
func DefaultSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Bold(true)
}

// DefaultBorderStyle returns default border styling
func DefaultBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
}
