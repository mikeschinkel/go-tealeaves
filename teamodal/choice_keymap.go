package teamodal

import "charm.land/bubbles/v2/key"

// ChoiceKeyMap defines key bindings for ChoiceModel
type ChoiceKeyMap struct {
	NextButton key.Binding // Tab, Right arrow
	PrevButton key.Binding // Shift+Tab, Left arrow
	Confirm    key.Binding // Enter
	Cancel     key.Binding // Esc
}

// DefaultChoiceKeyMap returns default key bindings for ChoiceModel
func DefaultChoiceKeyMap() ChoiceKeyMap {
	return ChoiceKeyMap{
		NextButton: key.NewBinding(
			key.WithKeys("tab", "right", "down"),
			key.WithHelp("tab/→/↓", "next"),
		),
		PrevButton: key.NewBinding(
			key.WithKeys("shift+tab", "left", "up"),
			key.WithHelp("shift+tab/←/↑", "prev"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}
