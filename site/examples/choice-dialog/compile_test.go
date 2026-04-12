// Source: site/src/content/docs/components/choice-dialog.mdx:21,139
package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teamodal"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TestCompile_ChoiceQuickExample verifies the quick example from choice-dialog.mdx.
func TestCompile_ChoiceQuickExample(t *testing.T) {
	modal := teamodal.NewChoiceModel(&teamodal.ChoiceModelArgs{
		Title:   "How should we proceed?",
		Message: "There are unsaved changes in 3 files.",
		Options: []teamodal.ChoiceOption{
			{Label: "Save & Exit", Hotkey: 's', ID: "save"},
			{Label: "Discard", Hotkey: 'd', ID: "discard"},
			{Label: "Continue", Hotkey: 'c', ID: "continue"},
		},
		DefaultIndex: 0,
		Orientation:  teamodal.Horizontal,
	})

	modal, _ = modal.Open()
	_ = modal.IsOpen()
	_ = modal.FocusButton()
}

// TestCompile_ChoiceConfiguration verifies configuration example from choice-dialog.mdx.
func TestCompile_ChoiceConfiguration(t *testing.T) {
	options := []teamodal.ChoiceOption{
		{Label: "Option A", Hotkey: 'a', ID: "a"},
		{Label: "Option B", Hotkey: 'b', ID: "b"},
	}

	allowCancel := true
	modal := teamodal.NewChoiceModel(&teamodal.ChoiceModelArgs{
		Title:       "How should we proceed?",
		Message:     "Select an action.",
		Options:     options,
		Orientation: teamodal.Vertical,
		AllowCancel: &allowCancel,
	})

	theme := teautils.DefaultTheme()
	modal = modal.WithTheme(theme)
	_ = modal
}
