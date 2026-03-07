package teacrumbs

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestWithTheme_AppliesStyles(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewBreadcrumbsModel().WithTheme(theme)

	if m.Styles.ParentStyle.GetForeground() == nil {
		t.Error("themed ParentStyle has no foreground")
	}
	if m.Styles.CurrentStyle.GetForeground() == nil {
		t.Error("themed CurrentStyle has no foreground")
	}
	if m.Styles.SeparatorStyle.GetForeground() == nil {
		t.Error("themed SeparatorStyle has no foreground")
	}
}

func TestWithTheme_PreservesSeparator(t *testing.T) {
	m := NewBreadcrumbsModel().WithSeparator(" / ")
	origSep := m.separator

	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m = m.WithTheme(theme)

	if m.separator != origSep {
		t.Errorf("separator changed: %q -> %q", origSep, m.separator)
	}
}

func TestWithTheme_ThenWithStyles_Overrides(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	custom := DefaultStyles()
	m := NewBreadcrumbsModel().WithTheme(theme).WithStyles(custom)

	// WithStyles should override WithTheme
	if m.Styles.ParentStyle.GetForeground() != custom.ParentStyle.GetForeground() {
		t.Error("WithStyles did not override WithTheme")
	}
}

func TestWithStyles_AfterWithTheme_Overrides(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewBreadcrumbsModel().WithTheme(theme)

	custom := DefaultStyles()
	m = m.WithStyles(custom)

	if m.Styles.CurrentStyle.GetForeground() != custom.CurrentStyle.GetForeground() {
		t.Error("WithStyles after WithTheme should override")
	}
}
