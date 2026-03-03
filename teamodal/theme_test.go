package teamodal

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestModalModel_WithTheme(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkPalette())
	m := NewOKModal("test", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	}).WithTheme(theme)

	if m.titleStyle.GetForeground() == nil {
		t.Error("themed titleStyle has no foreground")
	}
	if m.messageStyle.GetForeground() == nil {
		t.Error("themed messageStyle has no foreground")
	}
	if m.focusedButtonStyle.GetBackground() == nil {
		t.Error("themed focusedButtonStyle has no background")
	}
}

func TestModalModel_WithTheme_PreservesContent(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkPalette())
	m := NewOKModal("hello", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Test",
	}).WithTheme(theme)

	if m.message != "hello" {
		t.Errorf("message changed: %q", m.message)
	}
	if m.title != "Test" {
		t.Errorf("title changed: %q", m.title)
	}
}
