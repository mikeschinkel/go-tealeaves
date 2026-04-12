// Source: site/src/content/docs/components/confirm-dialog.mdx:21,50,205
package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teamodal"
	lipgloss "charm.land/lipgloss/v2"
)

// TestCompile_ConfirmQuickExample verifies the quick example from confirm-dialog.mdx.
func TestCompile_ConfirmQuickExample(t *testing.T) {
	modal := teamodal.NewYesNoModal("Save changes before exiting?", &teamodal.ConfirmModelArgs{
		Title:      "Unsaved Changes",
		DefaultYes: true,
		YesLabel:   "Save",
		NoLabel:    "Discard",
	})

	modal, _ = modal.Open()
	_ = modal.IsOpen()
}

// TestCompile_OKModal verifies the OK modal from confirm-dialog.mdx.
func TestCompile_OKModal(t *testing.T) {
	modal := teamodal.NewOKModal("Operation completed successfully.", &teamodal.ConfirmModelArgs{
		Title:   "Done",
		OKLabel: "Got it",
	})
	_ = modal
}

// TestCompile_ConfirmConfiguration verifies configuration from confirm-dialog.mdx.
func TestCompile_ConfirmConfiguration(t *testing.T) {
	modal := teamodal.NewYesNoModal("message", &teamodal.ConfirmModelArgs{
		Title:              "Title",
		DefaultYes:         true,
		YesLabel:           "Confirm",
		NoLabel:            "Cancel",
		TextAlign:          lipgloss.Left,
		ButtonAlign:        lipgloss.Center,
		TitleStyle:         lipgloss.NewStyle().Bold(true),
		MessageStyle:       lipgloss.NewStyle(),
		ButtonStyle:        lipgloss.NewStyle(),
		FocusedButtonStyle: lipgloss.NewStyle().Bold(true),
		BorderStyle:        lipgloss.NewStyle(),
	})

	modal = modal.
		WithTitle("New Title").
		WithMessage("New message").
		WithTextAlign(lipgloss.Left).
		SetSize(80, 24)
	_ = modal
}
