// Package teacolor provides named color constants for terminal UI components.
//
// It wraps lipgloss.Color to provide ANSI 256, standard ANSI names, and
// curated semantic aliases. All constants are of type color.Color, compatible
// with lipgloss v2's Foreground/Background methods.
package teacolor

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Color creates a color.Color from an ANSI color string.
// This is a convenience wrapper around lipgloss.Color.
func Color(s string) color.Color {
	return lipgloss.Color(s)
}
