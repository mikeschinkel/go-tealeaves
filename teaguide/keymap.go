package teaguide

import "charm.land/bubbles/v2/key"

// GuideKeyMap defines the key bindings for the guide overlay.
type GuideKeyMap struct {
	Close       key.Binding // esc — close guide
	ScrollUp    key.Binding // up, k — scroll viewport up
	ScrollDown  key.Binding // down, j — scroll viewport down
	ToggleBlock key.Binding // space, enter — toggle blocked section
}

// DefaultGuideKeyMap returns the default key bindings for the guide.
func DefaultGuideKeyMap() GuideKeyMap {
	return GuideKeyMap{
		Close: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("Esc", "close"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down", "scroll down"),
		),
		ToggleBlock: key.NewBinding(
			key.WithKeys("space", "enter"),
			key.WithHelp("space", "expand/collapse"),
		),
	}
}
