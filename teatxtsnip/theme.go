package teatxtsnip

import (
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with the selection style derived from the theme.
// The theme's FocusBg and TextPrimary palette colors are used for the
// selection highlight.
func (m TextSnipModel) WithTheme(theme teautils.Theme) TextSnipModel {
	SelectionStyle = lipgloss.NewStyle().
		Background(theme.System.FocusBg).
		Foreground(theme.System.TextPrimary)
	return m
}
