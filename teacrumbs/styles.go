package teacrumbs

import "charm.land/lipgloss/v2"

// Styles holds all styling for the breadcrumb crumbs.
type Styles struct {
	// ParentStyle is applied to all breadcrumbs except the last.
	ParentStyle lipgloss.Style

	// CurrentStyle is applied to the last (current) breadcrumb.
	CurrentStyle lipgloss.Style

	// SeparatorStyle is applied to the separator between breadcrumbs.
	SeparatorStyle lipgloss.Style

	// HoverStyle is applied to the breadcrumb under the mouse cursor.
	HoverStyle lipgloss.Style
}

// DefaultStyles returns defaults matching the current cyan/gray scheme.
func DefaultStyles() Styles {
	return Styles{
		ParentStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")), // Gray
		CurrentStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Cyan
			Bold(true),
		SeparatorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")), // Darker gray
		HoverStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")). // Light blue
			Underline(true),
	}
}
