package teatxtsnip

import "charm.land/bubbles/v2/key"

// SelectionKeyMap defines key bindings for text selection and clipboard operations
type SelectionKeyMap struct {
	// Selection by character
	SelectLeft  key.Binding
	SelectRight key.Binding
	SelectUp    key.Binding
	SelectDown  key.Binding

	// Selection by word
	SelectWordLeft  key.Binding
	SelectWordRight key.Binding

	// Selection to line boundaries
	SelectToLineStart key.Binding
	SelectToLineEnd   key.Binding

	// Selection to document boundaries
	SelectToStart key.Binding
	SelectToEnd   key.Binding

	// Select all
	SelectAll key.Binding

	// Clipboard operations
	Copy  key.Binding
	Cut   key.Binding
	Paste key.Binding

	// Clear selection
	ClearSelection key.Binding
}

// DefaultSelectionKeyMap returns the default key bindings for selection
func DefaultSelectionKeyMap() SelectionKeyMap {
	return SelectionKeyMap{
		// Selection by character
		SelectLeft: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+left", "select left"),
		),
		SelectRight: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+right", "select right"),
		),
		SelectUp: key.NewBinding(
			key.WithKeys("shift+up"),
			key.WithHelp("shift+up", "select up"),
		),
		SelectDown: key.NewBinding(
			key.WithKeys("shift+down"),
			key.WithHelp("shift+down", "select down"),
		),

		// Selection by word
		SelectWordLeft: key.NewBinding(
			key.WithKeys("ctrl+shift+left", "alt+shift+left"),
			key.WithHelp("ctrl+shift+left", "select word left"),
		),
		SelectWordRight: key.NewBinding(
			key.WithKeys("ctrl+shift+right", "alt+shift+right"),
			key.WithHelp("ctrl+shift+right", "select word right"),
		),

		// Selection to line boundaries
		SelectToLineStart: key.NewBinding(
			key.WithKeys("shift+home"),
			key.WithHelp("shift+home", "select to line start"),
		),
		SelectToLineEnd: key.NewBinding(
			key.WithKeys("shift+end"),
			key.WithHelp("shift+end", "select to line end"),
		),

		// Selection to document boundaries
		SelectToStart: key.NewBinding(
			key.WithKeys("ctrl+shift+home"),
			key.WithHelp("ctrl+shift+home", "select to start"),
		),
		SelectToEnd: key.NewBinding(
			key.WithKeys("ctrl+shift+end"),
			key.WithHelp("ctrl+shift+end", "select to end"),
		),

		// Select all
		SelectAll: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
		),

		// Clipboard operations
		// Note: We use tea.WithoutSignalHandler() to disable SIGINT, allowing Ctrl+C
		// to be used for copy. This is standard for terminal text editors.
		Copy: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "copy"),
		),
		Cut: key.NewBinding(
			key.WithKeys("ctrl+x"),
			key.WithHelp("ctrl+x", "cut"),
		),
		Paste: key.NewBinding(
			key.WithKeys("ctrl+v"),
			key.WithHelp("ctrl+v", "paste"),
		),

		// Clear selection
		ClearSelection: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear selection"),
		),
	}
}
