package teamodal

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy of the ConfirmModel with styles derived from the
// given theme. Individual With*Style() calls take precedence if called after.
func (m ConfirmModel) WithTheme(theme teautils.Theme) ConfirmModel {
	m.borderStyle = theme.Modal.BorderStyle
	m.titleStyle = theme.Modal.TitleStyle
	m.messageStyle = theme.Modal.MessageStyle
	m.buttonStyle = theme.Modal.ButtonStyle
	m.focusedButtonStyle = theme.Modal.FocusedButtonStyle
	return m
}
