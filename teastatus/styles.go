package teastatus

import "charm.land/lipgloss/v2"

// Styles holds all styling for the status bar.
type Styles struct {
	// Left side (menu items)
	MenuKeyStyle   lipgloss.Style // Style for "[key]"
	MenuLabelStyle lipgloss.Style // Style for "label"
	MenuSeparator  string         // Between menu items

	// Right side (indicators)
	IndicatorStyle    lipgloss.Style // Default indicator text style
	IndicatorSepStyle lipgloss.Style // Separator style
	SeparatorKind     SeparatorKind  // Which separator pattern

	// Bar-level
	BarStyle lipgloss.Style // Overall bar style (background, padding, etc.)
}

// DefaultStyles returns defaults matching the current cyan/gray scheme.
func DefaultStyles() Styles {
	return Styles{
		MenuKeyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Cyan
			Bold(true),
		MenuLabelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")), // Gray
		MenuSeparator: "  ",

		IndicatorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")), // Gray
		IndicatorSepStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")), // Darker gray
		SeparatorKind: PipeSeparator,

		BarStyle: lipgloss.NewStyle(),
	}
}
