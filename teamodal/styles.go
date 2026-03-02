package teamodal

import "charm.land/lipgloss/v2"

// DefaultBorderStyle returns default modal border styling
func DefaultBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("51"))
}

// DefaultTitleStyle returns default title styling
func DefaultTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46"))
}

// DefaultMessageStyle returns default message text styling
func DefaultMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))
}

// DefaultButtonStyle returns default button styling
func DefaultButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Padding(0, 1)
}

// DefaultFocusedButtonStyle returns default focused button styling
func DefaultFocusedButtonStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Bold(true).
		Padding(0, 1)
}

// DefaultCancelKeyStyle returns default styling for "[esc]" in the cancel hint
func DefaultCancelKeyStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("46"))
}

// DefaultCancelTextStyle returns default styling for "Cancel" in the cancel hint
func DefaultCancelTextStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("244"))
}
