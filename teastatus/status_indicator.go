package teastatus

import (
	"charm.land/lipgloss/v2"
)

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
