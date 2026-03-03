package teadepview

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestPathViewerModel_WithTheme(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkPalette())
	m := NewPathViewer(nil, PathViewerArgs{}).WithTheme(theme)

	if m.SelectedStyle.GetBackground() == nil {
		t.Error("themed SelectedStyle has no background")
	}
}
