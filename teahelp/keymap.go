package teahelp

import "charm.land/bubbles/v2/key"

// HelpVisorKeyMap defines the key bindings for the help visor
type HelpVisorKeyMap struct {
	Close    key.Binding // esc/? — close the visor
	PrevPage key.Binding // left — previous page
	NextPage key.Binding // right — next page
}

// DefaultHelpVisorKeyMap returns the default key bindings for the help visor
func DefaultHelpVisorKeyMap() HelpVisorKeyMap {
	return HelpVisorKeyMap{
		Close: key.NewBinding(
			key.WithKeys("esc", "?"),
			key.WithHelp("esc/?", "close help"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("\u2190", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("\u2192", "next page"),
		),
	}
}
