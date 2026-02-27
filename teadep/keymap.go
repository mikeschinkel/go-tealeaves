package teadep

import "github.com/charmbracelet/bubbles/key"

// PathViewerKeyMap defines the key bindings for dependency path navigation
type PathViewerKeyMap struct {
	Up           key.Binding // up/k - move up in path
	Down         key.Binding // down/j - move down in path
	OpenDropdown key.Binding // space/right - open dropdown with alternatives
	Select       key.Binding // enter - confirm selection on leaf nodes
}

// DefaultPathViewerKeyMap returns the default key bindings for path viewer
func DefaultPathViewerKeyMap() PathViewerKeyMap {
	return PathViewerKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		OpenDropdown: key.NewBinding(
			key.WithKeys(" ", "right"),
			key.WithHelp("space/→", "show alternatives"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}
