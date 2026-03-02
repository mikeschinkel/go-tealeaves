package teamodal

import "charm.land/lipgloss/v2"

// DefaultListItemStyle returns default styling for list items
func DefaultListItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))
}

// DefaultSelectedItemStyle returns default styling for the item at cursor position
func DefaultSelectedItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Reverse(true)
}

// DefaultActiveItemStyle returns default styling for the active/in-use item badge
func DefaultActiveItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("43")).
		Bold(true)
}

// DefaultListFooterStyle returns default styling for the footer with key hints
func DefaultListFooterStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)
}

// DefaultListScrollbarStyle returns default styling for the scrollbar
func DefaultListScrollbarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
}

// DefaultListScrollbarThumbStyle returns default styling for the scrollbar thumb
func DefaultListScrollbarThumbStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("248"))
}

// DefaultStatusStyle returns default styling for the status/feedback line
func DefaultStatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")). // Orange/yellow for attention
		MarginTop(1)
}

// DefaultEditItemStyle returns default styling for inline editing
func DefaultEditItemStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("11")). // Bright yellow background
		Foreground(lipgloss.Color("226")) // Yellow text
}
