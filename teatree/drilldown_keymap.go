package teatree

import "charm.land/bubbles/v2/key"

// DrillDownKeyMap defines the key bindings for drill-down path navigation
type DrillDownKeyMap struct {
	Up           key.Binding // up/k - move up in path
	Down         key.Binding // down/j - move down in path
	OpenDropdown key.Binding // space/right - open dropdown with alternatives
	Select       key.Binding // enter - confirm selection on leaf nodes
}

// DefaultDrillDownKeyMap returns the default key bindings for drill-down navigation
func DefaultDrillDownKeyMap() DrillDownKeyMap {
	return DrillDownKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		OpenDropdown: key.NewBinding(
			key.WithKeys("space", "right"),
			key.WithHelp("space/→", "show alternatives"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}
