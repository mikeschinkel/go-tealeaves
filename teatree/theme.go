package teatree

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with the given theme stored for use by node
// providers that are theme-aware. The default CompactNodeProvider does not
// use color (it uses Reverse for focus), so theming primarily benefits
// custom NodeProvider implementations that read the theme from the model.
func (m TreeModel[T]) WithTheme(theme teautils.Theme) TreeModel[T] {
	m.theme = &theme
	return m
}

// Theme returns the stored theme, or nil if no theme has been set.
func (m TreeModel[T]) Theme() *teautils.Theme {
	return m.theme
}
