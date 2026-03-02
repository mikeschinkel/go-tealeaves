package teadd

import "charm.land/bubbles/v2/key"

// DropdownKeyMap defines the key bindings for dropdown menus
type DropdownKeyMap struct {
	Up     key.Binding // up/k - move selection up
	Down   key.Binding // down/j - move selection down
	Select key.Binding // enter - select item
	Cancel key.Binding // esc - cancel/close dropdown
}

// DefaultDropdownKeyMap returns the default key bindings for dropdowns
func DefaultDropdownKeyMap() DropdownKeyMap {
	return DropdownKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}
