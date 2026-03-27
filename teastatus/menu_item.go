package teastatus

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// MenuItem represents a single key-action entry in the left menu area.
// Rendered as: "[key] label"
type MenuItem struct {
	KeyText string      // Display text for the key, e.g. "esc", "tab", "?"
	Label   string      // Short description, e.g. "Back", "Menu", "Quit"
	Binding key.Binding // Original binding (preserved for potential future use)
}

type MenuItemOpts struct {
	Label string // Short description, e.g. "Back", "Menu", "Quit"
}

// NewMenuItem creates a MenuItem
// Uses the first key from the binding as the display key.
func NewMenuItem(binding key.Binding, opts *MenuItemOpts) MenuItem {
	var o MenuItemOpts
	if opts == nil {
		opts = &o
	}
	var keyText string
	keys := binding.Keys()
	if len(keys) > 0 {
		keyText = keys[0]
	}
	if keyText == "" {
		panic(fmt.Sprintf("No key binding %v found for %v", binding, opts))
	}
	if opts.Label == "" {
		opts.Label = keyText
	}
	return MenuItem{
		KeyText: keyText,
		Label:   opts.Label,
		Binding: binding,
	}
}

// MenuItemsFromMetas creates MenuItems from a slice of KeyMeta
// (e.g., from KeyRegistry.ForStatusBar()). Uses StatusBarLabel as Label.
// TODO: I do not like this func name, but not yet sure what to rename it,
func MenuItemsFromMetas(metas []teautils.KeyMeta) []MenuItem {
	items := make([]MenuItem, len(metas))
	for i, meta := range metas {
		items[i] = NewMenuItem(meta.Binding, &MenuItemOpts{
			Label: meta.StatusBarLabel,
		})
	}
	return items
}
