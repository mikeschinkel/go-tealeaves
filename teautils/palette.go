package teautils

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils/teacolor"
)

// Palette holds semantic color slots for UI chrome. All fields use color.Color
// (compatible with lipgloss v2). Designed for embedding — apps add domain-specific
// color slots by embedding Palette in their own struct.
type Palette struct {
	// Text hierarchy (4 levels)
	TextPrimary   color.Color // Main content text
	TextSecondary color.Color // Labels, descriptions
	TextMuted     color.Color // Metadata, hints
	TextDim       color.Color // Disabled, placeholders

	// Accent colors (for highlights, titles, active elements)
	Accent       color.Color // Primary accent (titles, active items)
	AccentAlt    color.Color // Secondary accent (category headers, links)
	AccentSubtle color.Color // Subtle accent (key display, help text)

	// Selection and focus
	SelectionBg color.Color // Selected item background
	SelectionFg color.Color // Selected item foreground
	FocusBorder color.Color // Focused pane/component border
	FocusBg     color.Color // Focused row/cell background

	// Borders and chrome
	Border       color.Color // Default border color
	BorderAccent color.Color // Emphasized borders (modals, active panes)
	Separator    color.Color // Dividers, scrollbar tracks

	// Status indicators
	StatusSuccess color.Color // Success, confirmations
	StatusInfo    color.Color // Informational
	StatusWarn    color.Color // Warnings, caution
	StatusError   color.Color // Errors, destructive

	// Interactive elements
	ButtonBg      color.Color // Unfocused button background
	ButtonFg      color.Color // Unfocused button foreground
	ButtonFocusBg color.Color // Focused button background
	ButtonFocusFg color.Color // Focused button foreground

	// Edit mode
	EditBg color.Color // Inline edit background
	EditFg color.Color // Inline edit foreground

	// Diff context (tints for code/content backgrounds)
	TintPositive color.Color // Added content background tint
	TintNegative color.Color // Removed content background tint
	TintNeutral  color.Color // Unchanged content background tint

	// Scrollbar
	ScrollTrack color.Color // Scrollbar track
	ScrollThumb color.Color // Scrollbar thumb
}

// DarkPalette returns a Palette with colors suited for dark terminal backgrounds.
// The colors are mapped from the existing hardcoded values across go-tealeaves
// components to provide visual continuity.
func DarkPalette() Palette {
	return Palette{
		// Text hierarchy
		TextPrimary:   teacolor.Color15,  // white — items, buttons
		TextSecondary: teacolor.Color252, // light gray — labels, descriptions
		TextMuted:     teacolor.Color244, // medium gray — hints, cancel text
		TextDim:       teacolor.Color240, // dark gray — disabled, separators

		// Accent colors
		Accent:       teacolor.Color46,  // bright green — titles, active items
		AccentAlt:    teacolor.Color178, // gold — category headers
		AccentSubtle: teacolor.Color86,  // cyan — key display, help text

		// Selection and focus
		SelectionBg: teacolor.Color62,  // purple — selected background
		SelectionFg: teacolor.Color230, // bright yellow — selected text
		FocusBorder: teacolor.Color51,  // cyan — focused modal border
		FocusBg:     teacolor.Color39,  // light blue — focused row

		// Borders and chrome
		Border:       teacolor.Color240, // dark gray — default borders
		BorderAccent: teacolor.Color51,  // cyan — modal borders
		Separator:    teacolor.Color240, // dark gray — dividers

		// Status indicators
		StatusSuccess: teacolor.Color46,  // bright green
		StatusInfo:    teacolor.Color86,  // cyan
		StatusWarn:    teacolor.Color214, // orange
		StatusError:   teacolor.Color160, // red

		// Interactive elements
		ButtonBg:      nil,               // no bg for unfocused
		ButtonFg:      teacolor.Color15,  // white
		ButtonFocusBg: teacolor.Color62,  // purple
		ButtonFocusFg: teacolor.Color230, // bright yellow

		// Edit mode
		EditBg: teacolor.Color226, // bright yellow
		EditFg: teacolor.Color232, // near-black

		// Diff context
		TintPositive: teacolor.Color22, // dark green
		TintNegative: teacolor.Color52, // dark red
		TintNeutral:  nil,              // no tint

		// Scrollbar
		ScrollTrack: teacolor.Color240, // dark gray
		ScrollThumb: teacolor.Color248, // lighter gray
	}
}

// LightPalette returns a Palette with colors suited for light terminal backgrounds.
func LightPalette() Palette {
	return Palette{
		// Text hierarchy
		TextPrimary:   teacolor.Color0,   // black
		TextSecondary: teacolor.Color238, // dark gray
		TextMuted:     teacolor.Color244, // medium gray
		TextDim:       teacolor.Color250, // light gray

		// Accent colors
		Accent:       teacolor.Color28,  // dark green
		AccentAlt:    teacolor.Color130, // dark gold
		AccentSubtle: teacolor.Color30,  // dark cyan

		// Selection and focus
		SelectionBg: teacolor.Color62,  // purple
		SelectionFg: teacolor.Color230, // bright yellow
		FocusBorder: teacolor.Color33,  // medium blue
		FocusBg:     teacolor.Color153, // light blue

		// Borders and chrome
		Border:       teacolor.Color250, // light gray
		BorderAccent: teacolor.Color33,  // medium blue
		Separator:    teacolor.Color250, // light gray

		// Status indicators
		StatusSuccess: teacolor.Color28,  // dark green
		StatusInfo:    teacolor.Color30,  // dark cyan
		StatusWarn:    teacolor.Color166, // dark orange
		StatusError:   teacolor.Color124, // dark red

		// Interactive elements
		ButtonBg:      nil,               // no bg for unfocused
		ButtonFg:      teacolor.Color0,   // black
		ButtonFocusBg: teacolor.Color62,  // purple
		ButtonFocusFg: teacolor.Color230, // bright yellow

		// Edit mode
		EditBg: teacolor.Color226, // bright yellow
		EditFg: teacolor.Color0,   // black

		// Diff context
		TintPositive: teacolor.Color157, // light green
		TintNegative: teacolor.Color217, // light red/pink
		TintNeutral:  nil,               // no tint

		// Scrollbar
		ScrollTrack: teacolor.Color252, // very light gray
		ScrollThumb: teacolor.Color244, // medium gray
	}
}

// AdaptivePalette detects the terminal background and returns the appropriate
// palette. Falls back to DarkPalette if detection fails.
func AdaptivePalette() Palette {
	if lipgloss.HasDarkBackground(os.Stdin, os.Stderr) {
		return DarkPalette()
	}
	return LightPalette()
}

// DefaultPalette returns AdaptivePalette().
func DefaultPalette() Palette {
	return AdaptivePalette()
}
