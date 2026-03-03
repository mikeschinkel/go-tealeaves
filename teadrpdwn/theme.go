package teadrpdwn

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with styles derived from the given theme.
// Direct field assignment takes precedence if done after WithTheme().
func (m DropdownModel) WithTheme(theme teautils.Theme) DropdownModel {
	m.BorderStyle = theme.Dropdown.BorderStyle
	m.ItemStyle = theme.Dropdown.ItemStyle
	m.SelectedStyle = theme.Dropdown.SelectedStyle
	return m
}
