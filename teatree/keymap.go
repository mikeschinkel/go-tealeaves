package teatree

import "charm.land/bubbles/v2/key"

// TreeKeyMap defines the key bindings for tree navigation
type TreeKeyMap struct {
	Up            key.Binding // up/k - move up
	Down          key.Binding // down/j - move down
	ExpandOrEnter key.Binding // right/l - expand or move to first child
	CollapseOrUp  key.Binding // left/h - collapse or move to parent
	Toggle        key.Binding // enter/space - toggle expansion
}

// DefaultTreeKeyMap returns the default key bindings for tree navigation
func DefaultTreeKeyMap() TreeKeyMap {
	return TreeKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		ExpandOrEnter: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "expand/enter"),
		),
		CollapseOrUp: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "collapse/up to parent"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("enter", "space"),
			key.WithHelp("enter/space", "toggle expand/collapse"),
		),
	}
}
