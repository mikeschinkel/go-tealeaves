package teadrpdwn

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestDropdownModel_WithTheme(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkPalette())
	opts := []Option{{Text: "A"}, {Text: "B"}}
	m := NewDropdownModel(opts, 0, 0, nil).WithTheme(theme)

	if m.SelectedStyle.GetBackground() == nil {
		t.Error("themed SelectedStyle has no background")
	}
	if m.ItemStyle.GetForeground() == nil {
		t.Error("themed ItemStyle has no foreground")
	}
}
