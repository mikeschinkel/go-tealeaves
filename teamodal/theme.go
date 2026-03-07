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

// WithTheme returns a copy of the ChoiceModel with styles derived from the
// given theme. Individual With*Style() calls take precedence if called after.
func (m ChoiceModel) WithTheme(theme teautils.Theme) ChoiceModel {
	m.borderStyle = theme.Modal.BorderStyle
	m.titleStyle = theme.Modal.TitleStyle
	m.messageStyle = theme.Modal.MessageStyle
	m.buttonStyle = theme.Modal.ButtonStyle
	m.focusedButtonStyle = theme.Modal.FocusedButtonStyle
	m.cancelKeyStyle = theme.Modal.CancelKeyStyle
	m.cancelTextStyle = theme.Modal.CancelTextStyle
	return m
}

// WithTheme returns a copy of the ListModel with styles derived from the
// given theme. Individual With*Style() calls take precedence if called after.
func (m ListModel[T]) WithTheme(theme teautils.Theme) ListModel[T] {
	m.borderStyle = theme.Modal.BorderStyle
	m.titleStyle = theme.Modal.TitleStyle
	m.itemStyle = theme.List.ItemStyle
	m.selectedItemStyle = theme.List.SelectedItemStyle
	m.activeItemStyle = theme.List.ActiveItemStyle
	m.footerStyle = theme.List.FooterStyle
	m.statusStyle = theme.List.StatusStyle
	m.editItemStyle = theme.List.EditItemStyle
	return m
}
