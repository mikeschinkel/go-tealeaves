package teatxtsnip

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestModel_WithTheme(t *testing.T) {
	m := NewTextSnipModel(nil)
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m = m.WithTheme(theme)

	// SelectionStyle should have been updated
	if SelectionStyle.GetBackground() == nil {
		t.Error("themed SelectionStyle has no background")
	}
	if SelectionStyle.GetForeground() == nil {
		t.Error("themed SelectionStyle has no foreground")
	}
}
