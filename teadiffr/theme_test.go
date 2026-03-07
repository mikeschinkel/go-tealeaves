package teadiffr

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestNewThemedTUIRenderer(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	r := NewThemedTUIRenderer(theme)

	if r.AddedColor == nil {
		t.Error("themed AddedColor is nil")
	}
	if r.DeletedColor == nil {
		t.Error("themed DeletedColor is nil")
	}
	if r.FileHeaderColor == nil {
		t.Error("themed FileHeaderColor is nil")
	}
	if r.NewBgColor == nil {
		t.Error("themed NewBgColor is nil")
	}
	if r.DeletedBgColor == nil {
		t.Error("themed DeletedBgColor is nil")
	}
}
