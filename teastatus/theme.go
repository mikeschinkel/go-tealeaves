package teastatus

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with styles derived from the given theme.
// Individual WithStyles() calls take precedence over theme if called after.
func (m StatusBarModel) WithTheme(theme teautils.Theme) StatusBarModel {
	m.Styles = Styles{
		MenuKeyStyle:      theme.StatusBar.MenuKeyStyle,
		MenuLabelStyle:    theme.StatusBar.MenuLabelStyle,
		MenuSeparator:     m.Styles.MenuSeparator,
		IndicatorStyle:    theme.StatusBar.IndicatorStyle,
		IndicatorSepStyle: theme.StatusBar.IndicatorSepStyle,
		SeparatorKind:     m.Styles.SeparatorKind,
		BarStyle:          theme.StatusBar.BarStyle,
	}
	return m
}
