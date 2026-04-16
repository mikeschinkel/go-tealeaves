// Source: site/src/content/docs/components/progress-dialog.mdx:22#9b99db49
package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// TestCompile_ProgressDialogQuickExample verifies the quick example from progress-dialog.mdx.
func TestCompile_ProgressDialogQuickExample(t *testing.T) {
	modal := teamodal.NewProgressModel(&teamodal.ProgressModelArgs{
		Title:             "Commit Message",
		BackgroundEnabled: true,
	})

	modal = modal.Open()
	_ = modal.IsOpen()
}
