package teamodal

import "github.com/charmbracelet/bubbles/key"

// ModalKeyMap defines the key bindings for modal dialogs
type ModalKeyMap struct {
	Confirm     key.Binding // enter - confirm selection
	Cancel      key.Binding // esc - cancel/close
	NextButton  key.Binding // tab - move to next button (YesNo only)
	PrevButton  key.Binding // shift+tab - move to previous button (YesNo only)
	SelectLeft  key.Binding // left - select left button (Yes) (YesNo only)
	SelectRight key.Binding // right - select right button (No) (YesNo only)
}

// DefaultModalKeyMap returns the default key bindings for modals
func DefaultModalKeyMap() ModalKeyMap {
	return ModalKeyMap{
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		NextButton: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next button"),
		),
		PrevButton: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous button"),
		),
		SelectLeft: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "select Yes"),
		),
		SelectRight: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "select No"),
		),
	}
}
