package teagrid

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with styles derived from the given theme.
// Individual With*Style() calls take precedence over theme if called after.
func (m GridModel) WithTheme(theme teautils.Theme) GridModel {
	m.headerStyle = theme.Grid.HeaderStyle
	m.baseStyle = theme.Grid.BaseStyle
	m.highlightStyle = theme.Grid.HighlightStyle
	return m
}
