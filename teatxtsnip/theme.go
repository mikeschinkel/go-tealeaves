package teatxtsnip

import (
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with the selection style derived from the theme.
// The theme's FocusBg and TextPrimary palette colors are used for the
// selection highlight.
func (m Model) WithTheme(theme teautils.Theme) Model {
	SelectionStyle = lipgloss.NewStyle().
		Background(theme.Palette.FocusBg).
		Foreground(theme.Palette.TextPrimary)
	return m
}
