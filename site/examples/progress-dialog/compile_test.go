package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// TestCompile_ProgressDialogQuickExample verifies the quick example from progress-dialog.mdx.
func TestCompile_ProgressDialogQuickExample(t *testing.T) {
	modal := teamodal.NewProgressModal(&teamodal.ProgressModalArgs{
		Title:             "Commit Message",
		BackgroundEnabled: true,
	})

	modal = modal.Open()
	_ = modal.IsOpen()
}
