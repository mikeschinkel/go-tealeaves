package teastatus

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestWithTheme_AppliesStyles(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkPalette())
	m := NewStatusBarModel().WithTheme(theme)

	if m.Styles.MenuKeyStyle.GetForeground() == nil {
		t.Error("themed MenuKeyStyle has no foreground")
	}
	if m.Styles.MenuLabelStyle.GetForeground() == nil {
		t.Error("themed MenuLabelStyle has no foreground")
	}
	if m.Styles.IndicatorSepStyle.GetForeground() == nil {
		t.Error("themed IndicatorSepStyle has no foreground")
	}
}

func TestWithTheme_PreservesNonStyleFields(t *testing.T) {
	m := NewStatusBarModel()
	origSep := m.Styles.MenuSeparator
	origKind := m.Styles.SeparatorKind

	theme := teautils.NewTheme(teautils.DarkPalette())
	m = m.WithTheme(theme)

	if m.Styles.MenuSeparator != origSep {
		t.Errorf("MenuSeparator changed: %q -> %q", origSep, m.Styles.MenuSeparator)
	}
	if m.Styles.SeparatorKind != origKind {
		t.Errorf("SeparatorKind changed: %v -> %v", origKind, m.Styles.SeparatorKind)
	}
}

func TestWithTheme_ThenWithStyles_Overrides(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkPalette())
	custom := DefaultStyles()
	m := NewStatusBarModel().WithTheme(theme).WithStyles(custom)

	// WithStyles should override WithTheme
	if m.Styles.MenuKeyStyle.GetForeground() != custom.MenuKeyStyle.GetForeground() {
		t.Error("WithStyles did not override WithTheme")
	}
}
