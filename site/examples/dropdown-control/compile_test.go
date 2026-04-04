package examples_test

import (
	"testing"

	key "charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teafields"
)

// TestCompile_DropdownQuickExample verifies the quick example from dropdown-control.mdx.
func TestCompile_DropdownQuickExample(t *testing.T) {
	items := teafields.ToOptions([]string{"Apple", "Banana", "Cherry", "Date", "Elderberry"})
	dropdown := teafields.NewDropdownModel(items, &teafields.DropdownModelArgs{
		FieldRow: 3,
		FieldCol: 18,
	})
	_ = dropdown
}

// TestCompile_DropdownKeyMap verifies key map configuration from dropdown-control.mdx.
func TestCompile_DropdownKeyMap(t *testing.T) {
	items := teafields.ToOptions([]string{"A", "B", "C"})
	dropdown := teafields.NewDropdownModel(items, &teafields.DropdownModelArgs{})

	km := teafields.DefaultDropdownKeyMap()
	km.Up = key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "move up"))
	dropdown.Keys = km
	_ = dropdown
}

// TestCompile_DropdownWithPosition verifies wither method usage from dropdown-control.mdx.
func TestCompile_DropdownWithPosition(t *testing.T) {
	items := teafields.ToOptions([]string{"A", "B"})
	dropdown := teafields.NewDropdownModel(items, &teafields.DropdownModelArgs{})

	dropdown = dropdown.
		WithScreenSize(80, 24).
		WithPosition(3, 10)
	_ = dropdown
}

// TestCompile_DropdownMessages verifies message types from dropdown-control.mdx.
func TestCompile_DropdownMessages(t *testing.T) {
	var msg tea.Msg
	switch m := msg.(type) {
	case teafields.OptionSelectedMsg:
		_ = m.Text
		_ = m.Value
	case teafields.DropdownCancelledMsg:
		// no fields
	}
}
