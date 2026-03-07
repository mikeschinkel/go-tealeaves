package teahelp

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestHelpVisorModel_WithTheme(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewHelpVisorModel().WithTheme(theme)

	if m.Styles.BorderStyle.GetBorderBottomSize() != 0 {
		t.Error("themed BorderStyle should have no bottom border")
	}
	if m.Styles.FooterKeyStyle.GetForeground() == nil {
		t.Error("themed FooterKeyStyle has no foreground")
	}
	if m.Styles.FooterLabelStyle.GetForeground() == nil {
		t.Error("themed FooterLabelStyle has no foreground")
	}
	if m.contentStyle.TitleStyle.GetForeground() == nil {
		t.Error("themed contentStyle TitleStyle has no foreground")
	}
	if m.contentStyle.KeyStyle.GetForeground() == nil {
		t.Error("themed contentStyle KeyStyle has no foreground")
	}
}

func TestHelpVisorModel_WithTheme_PreservesKeys(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewHelpVisorModel().WithTheme(theme)

	if len(m.Keys.Close.Keys()) == 0 {
		t.Error("WithTheme should not clear key bindings")
	}
}
