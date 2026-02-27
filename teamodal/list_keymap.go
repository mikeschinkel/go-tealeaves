package teamodal

import "github.com/charmbracelet/bubbles/key"

// ListKeyMap defines the key bindings for list modal dialogs
type ListKeyMap struct {
	Up      key.Binding // Up arrow, k - move cursor up
	Down    key.Binding // Down arrow, j - move cursor down
	Preview key.Binding // Space - preview select (mark active, keep browsing)
	Accept  key.Binding // Enter - commit selection and close
	New     key.Binding // a - request new item creation
	Edit    key.Binding // e, F2 - request edit of item at cursor
	Delete  key.Binding // d - request deletion of item at cursor
	Help    key.Binding // ? - toggle help visor
	Cancel  key.Binding // Esc - close without selection
}

// DefaultListKeyMap returns the default key bindings for list modals
func DefaultListKeyMap() ListKeyMap {
	return ListKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "Move cursor up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "Move cursor down"),
		),
		Preview: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "Preview select (mark active)"),
		),
		Accept: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Accept selection and close"),
		),
		New: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Add new item"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e", "f2"),
			key.WithHelp("e/F2", "Edit item inline"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Delete item"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Toggle help"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Cancel and close"),
		),
	}
}
