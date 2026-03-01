package teagrid

import "charm.land/bubbles/v2/key"

// KeyMap defines the key bindings for the grid.
type KeyMap struct {
	RowDown key.Binding
	RowUp   key.Binding

	RowSelectToggle key.Binding

	PageDown  key.Binding
	PageUp    key.Binding
	PageFirst key.Binding
	PageLast  key.Binding

	// CellLeft moves the cell cursor left.
	CellLeft key.Binding

	// CellRight moves the cell cursor right.
	CellRight key.Binding

	// CellSelect emits a UserEventCellSelected event.
	CellSelect key.Binding

	// Filter starts the filter input.
	Filter key.Binding

	// FilterBlur exits filter input.
	FilterBlur key.Binding

	// FilterClear clears the filter while blurred.
	FilterClear key.Binding

	// ScrollRight scrolls the viewport one column right.
	ScrollRight key.Binding

	// ScrollLeft scrolls the viewport one column left.
	ScrollLeft key.Binding
}

// DefaultKeyMap returns sensible default key bindings for Charm v2.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		RowDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		RowUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		RowSelectToggle: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "select row"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdn", "next page"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "previous page"),
		),
		PageFirst: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "first page"),
		),
		PageLast: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "last page"),
		),
		CellLeft: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "cell left"),
		),
		CellRight: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "cell right"),
		),
		CellSelect: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select cell"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		FilterBlur: key.NewBinding(
			key.WithKeys("enter", "escape"),
			key.WithHelp("enter/esc", "stop filtering"),
		),
		FilterClear: key.NewBinding(
			key.WithKeys("escape"),
			key.WithHelp("esc", "clear filter"),
		),
		ScrollRight: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+→", "scroll right"),
		),
		ScrollLeft: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+←", "scroll left"),
		),
	}
}
