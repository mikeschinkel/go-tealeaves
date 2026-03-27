package teaguide

import "charm.land/lipgloss/v2"

// GuideStyles configures the visual appearance of the guide overlay.
type GuideStyles struct {
	Border         lipgloss.Style
	Title          lipgloss.Style
	SectionHeading lipgloss.Style
	ItemKey        lipgloss.Style
	ItemLabel      lipgloss.Style
	ItemProse      lipgloss.Style
	BlockedHeading lipgloss.Style
	BlockedItem    lipgloss.Style
	BlockReason    lipgloss.Style
	Footer         lipgloss.Style
}

// DefaultGuideStyles returns sensible default styles for the guide overlay.
func DefaultGuideStyles() GuideStyles {
	return GuideStyles{
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("219")),
		SectionHeading: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("75")),
		ItemKey: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("114")),
		ItemLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		ItemProse: lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			Italic(true),
		BlockedHeading: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		BlockedItem: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		BlockReason: lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true),
		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
	}
}
