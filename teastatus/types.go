package teastatus

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a single key-action entry in the left menu area.
// Rendered as: "[key] label"
type MenuItem struct {
	Key     string      // Display text for the key, e.g. "esc", "tab", "?"
	Label   string      // Short description, e.g. "Back", "Menu", "Quit"
	Binding key.Binding // Original binding (preserved for potential future use)
}

// NewMenuItemFromBinding creates a MenuItem from a key.Binding and label.
// Uses the first key from the binding as the display key.
func NewMenuItemFromBinding(binding key.Binding, label string) MenuItem {
	keyStr := ""
	keys := binding.Keys()
	if len(keys) > 0 {
		keyStr = keys[0]
	}
	return MenuItem{
		Key:     keyStr,
		Label:   label,
		Binding: binding,
	}
}

// StatusIndicator represents a single text indicator in the right status area.
// Examples: "Verified", "3 batches", "DEPS IN-FLUX"
type StatusIndicator struct {
	Text  string         // Display text
	Style lipgloss.Style // Optional per-indicator style override (zero value = use default)
}

// NewStatusIndicator creates a StatusIndicator with the given text.
func NewStatusIndicator(text string) StatusIndicator {
	return StatusIndicator{Text: text}
}

// WithStyle returns a copy with the given style override.
func (si StatusIndicator) WithStyle(style lipgloss.Style) StatusIndicator {
	si.Style = style
	return si
}

// SeparatorKind identifies which separator style to use for indicators.
type SeparatorKind int

const (
	// PipeSeparator separates indicators with " | "
	PipeSeparator SeparatorKind = iota
	// SpaceSeparator separates indicators with "  "
	SpaceSeparator
	// BracketSeparator wraps each indicator as "[text]"
	BracketSeparator
)
