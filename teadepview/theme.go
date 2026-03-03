package teadepview

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with styles derived from the given theme.
// Direct field assignment takes precedence if done after WithTheme().
func (m PathViewerModel) WithTheme(theme teautils.Theme) PathViewerModel {
	m.BorderStyle = theme.Dropdown.BorderStyle
	m.SelectedStyle = theme.Dropdown.SelectedStyle
	return m
}
