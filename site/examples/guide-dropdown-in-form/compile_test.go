// Source: site/src/content/docs/cookbook/dropdown-in-form.mdx:203,218
package main_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teafields"
)

// TestCompile_DropdownConsumptionSnippet verifies that the dropdown consumption
// pattern from the cookbook compiles: call Update, check cmd, type-assert.
// Source line 203.
func TestCompile_DropdownConsumptionSnippet(t *testing.T) {
	options := teafields.ToOptions([]string{"A", "B"})
	dropdown := teafields.NewDropdownModel(options, nil)

	var msg tea.Msg = tea.KeyPressMsg{}

	result, cmd := dropdown.Update(msg)
	if cmd != nil {
		dropdown = result.(teafields.DropdownModel)
		_ = dropdown
	}
}

// TestCompile_OverlayDropdownSnippet verifies that OverlayDropdown is callable
// with string background/foreground and integer row/col offsets.
// Source line 218.
func TestCompile_OverlayDropdownSnippet(t *testing.T) {
	options := teafields.ToOptions([]string{"A", "B"})
	dropdown := teafields.NewDropdownModel(options, nil)

	content := "background content\nrow two"
	dropdownView := dropdown.View()

	result := teafields.OverlayDropdown(
		content,              // background: the full form
		dropdownView.Content, // foreground: the dropdown popup
		dropdown.Row+2,       // +2 for border row + padding row
		dropdown.Col+3,       // +3 for border col + 2 padding cols
	)
	_ = result
}
