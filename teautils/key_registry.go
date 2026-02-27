package teautils

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/mikeschinkel/go-dt/dtx"
)

// KeyMeta holds app-specific display metadata for a key binding
type KeyMeta struct {
	ID      KeyIdentifier // Unique identifier (e.g., "file-intent.commit")
	Binding key.Binding   // Standard Bubble Tea binding

	// App-level display decisions
	StatusBar      bool     // Show in status bar?
	StatusBarLabel string   // Short label for status bar (e.g., "Menu", "Back") - optional, defaults to binding.Help().Desc
	HelpModal      bool     // Show in help modal?
	Category       string   // Help modal category ("Navigation", "File Actions", etc.)
	HelpText       string   // Extended description for help modal (optional, defaults to binding.Help().Desc)
	DisplayKeys    []string // Custom display names for keys (e.g., []string{"space"} for " ")
}

// KeyRegistry manages key bindings and their display metadata
// Allows the app to control presentation while components use standard key.Binding
type KeyRegistry struct {
	keys *dtx.OrderedMap[KeyIdentifier, *KeyMeta]
}

// NewKeyRegistry creates a new key registry
func NewKeyRegistry() *KeyRegistry {
	return &KeyRegistry{
		keys: dtx.NewOrderedMap[KeyIdentifier, *KeyMeta](16),
	}
}

// Register adds or updates a key binding with metadata
// If HelpText is empty, defaults to the binding's Help().Desc
func (r *KeyRegistry) Register(meta KeyMeta) (err error) {
	if meta.ID == "" {
		err = ErrEmptyKeyID
		goto end
	}

	// Default HelpText to binding's help description if not provided
	if meta.HelpText == "" {
		help := meta.Binding.Help()
		if help.Desc != "" {
			meta.HelpText = help.Desc
		}
	}

	r.keys.Set(meta.ID, &meta)

end:
	return err
}

// MustRegister is like Register but panics on error
func (r *KeyRegistry) MustRegister(meta KeyMeta) {
	err := r.Register(meta)
	if err != nil {
		panic(err)
	}
}

// RegisterMany adds multiple key bindings at once
// Returns error on first failure
func (r *KeyRegistry) RegisterMany(metas []KeyMeta) (err error) {
	for _, meta := range metas {
		err = r.Register(meta)
		if err != nil {
			goto end
		}
	}

end:
	return err
}

// MustRegisterMany is like RegisterMany but panics on error
func (r *KeyRegistry) MustRegisterMany(metas []KeyMeta) {
	err := r.RegisterMany(metas)
	if err != nil {
		panic(err)
	}
}

// Get retrieves a key by its identifier
func (r *KeyRegistry) Get(id KeyIdentifier) (meta *KeyMeta, err error) {
	var exists bool
	meta, exists = r.keys.Get(id)
	if !exists {
		err = NewErr(ErrKeyNotFound, "id", id)
		goto end
	}

end:
	return meta, err
}

// Clear removes all keys from the registry
// Useful when switching views/contexts
func (r *KeyRegistry) Clear() {
	r.keys.Clear()
}

// ForStatusBar returns keys marked for status bar display, in registration order
func (r *KeyRegistry) ForStatusBar() []KeyMeta {
	var result []KeyMeta
	for _, meta := range r.keys.Iterator() {
		if meta.StatusBar {
			result = append(result, *meta)
		}
	}
	return result
}

// ForHelpModal returns all keys marked for help modal display
func (r *KeyRegistry) ForHelpModal() []KeyMeta {
	var result []KeyMeta
	for _, meta := range r.keys.Iterator() {
		if meta.HelpModal {
			result = append(result, *meta)
		}
	}
	return result
}

// ByCategory returns keys grouped by category for help modal
// Keys without a category are grouped under "Other"
// Preserves registration order within each category
func (r *KeyRegistry) ByCategory() map[string][]KeyMeta {
	result := make(map[string][]KeyMeta)

	// Iterate in registration order (OrderedMap preserves insertion order)
	for _, meta := range r.keys.Iterator() {
		if !meta.HelpModal {
			continue
		}

		category := meta.Category
		if category == "" {
			category = "Other"
		}

		result[category] = append(result[category], *meta)
	}

	return result
}
