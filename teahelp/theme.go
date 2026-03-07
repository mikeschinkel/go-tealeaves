package teahelp

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy of the HelpVisorModel with styles derived from the
// given theme. Individual With*Style() calls take precedence if called after.
func (m HelpVisorModel) WithTheme(theme teautils.Theme) HelpVisorModel {
	m.contentStyle = teautils.ThemedHelpVisorStyle(theme)
	m.Styles = HelpVisorStyles{
		BorderStyle: theme.BorderAccent.
			BorderBottom(false).
			PaddingTop(0).
			PaddingRight(3).
			PaddingBottom(1).
			PaddingLeft(0),
		FooterKeyStyle:   theme.HelpVisor.KeyStyle,
		FooterLabelStyle: theme.HelpVisor.DescStyle,
	}
	return m
}
