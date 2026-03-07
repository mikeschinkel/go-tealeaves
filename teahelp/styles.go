package teahelp

import "charm.land/lipgloss/v2"

// HelpVisorStyles holds the styles for the help visor chrome (border and footer).
// Content styling (title, category, key, description) is handled by
// teautils.HelpVisorStyle, which is passed via WithContentStyle.
type HelpVisorStyles struct {
	BorderStyle     lipgloss.Style // Open-bottom rounded border
	FooterKeyStyle  lipgloss.Style // Style for footer key labels (e.g., "←")
	FooterLabelStyle lipgloss.Style // Style for footer text labels (e.g., "Prev", "Page 1/3")
}

// DefaultHelpVisorStyles returns the default styles for the help visor
func DefaultHelpVisorStyles() HelpVisorStyles {
	return HelpVisorStyles{
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("99")).
			BorderBottom(false).
			PaddingTop(0).
			PaddingRight(3).
			PaddingBottom(1).
			PaddingLeft(0),
		FooterKeyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true),
		FooterLabelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
	}
}
