package teamodal

import "charm.land/bubbles/v2/key"

// MultiSelectKeyMap defines key bindings for MultiSelectModel
type MultiSelectKeyMap struct {
	Up         key.Binding // up, k — cursor up (list focus)
	Down       key.Binding // down, j — cursor down (list focus)
	Toggle     key.Binding // space — toggle checkbox (list focus)
	Confirm    key.Binding // enter — toggle (list focus) or activate (button focus)
	Cancel     key.Binding // esc — cancel
	NextFocus  key.Binding // tab — cycle: list → button0 → ... → list
	PrevFocus  key.Binding // shift+tab — reverse cycle
	NextButton key.Binding // right, l — between buttons only
	PrevButton key.Binding // left, h — between buttons only
}

// DefaultMultiSelectKeyMap returns default key bindings for MultiSelectModel
func DefaultMultiSelectKeyMap() MultiSelectKeyMap {
	return MultiSelectKeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "Move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "Move down"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "Toggle checkbox"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Toggle / Activate"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Cancel"),
		),
		NextFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "Next focus"),
		),
		PrevFocus: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "Prev focus"),
		),
		NextButton: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "Next button"),
		),
		PrevButton: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "Prev button"),
		),
	}
}
