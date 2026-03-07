package teautils

import (
	"os"

	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils/teacolor"
)

// SystemPalette holds semantic color slots for UI chrome. All fields use
// SemanticColor which wraps color.Color with pre-built lipgloss styles.
// Components use these colors via Theme; apps compose SystemPalette with
// domain-specific colors via Palette[T].
type SystemPalette struct {
	// Text hierarchy (4 levels)
	TextPrimary   SemanticColor `theme:"text.primary"`   // Main content text
	TextSecondary SemanticColor `theme:"text.secondary"` // Labels, descriptions
	TextMuted     SemanticColor `theme:"text.muted"`     // Metadata, hints
	TextDim       SemanticColor `theme:"text.dim"`       // Disabled, placeholders

	// Accent colors (for highlights, titles, active elements)
	Accent       SemanticColor `theme:"accent.primary"` // Primary accent (titles, active items)
	AccentAlt    SemanticColor `theme:"accent.alt"`     // Secondary accent (category headers, links)
	AccentSubtle SemanticColor `theme:"accent.subtle"`  // Subtle accent (key display, help text)

	// Selection and focus
	SelectionBg SemanticColor `theme:"selection.bg"` // Selected item background
	SelectionFg SemanticColor `theme:"selection.fg"` // Selected item foreground
	FocusBorder SemanticColor `theme:"focus.border"` // Focused pane/component border
	FocusBg     SemanticColor `theme:"focus.bg"`     // Focused row/cell background

	// Borders and chrome
	Border       SemanticColor `theme:"border.default"` // Default border color
	BorderAccent SemanticColor `theme:"border.accent"`  // Emphasized borders (modals, active panes)
	Separator    SemanticColor `theme:"separator"`      // Dividers, scrollbar tracks

	// Status indicators
	StatusSuccess SemanticColor `theme:"status.success"` // Success, confirmations
	StatusInfo    SemanticColor `theme:"status.info"`    // Informational
	StatusWarn    SemanticColor `theme:"status.warn"`    // Warnings, caution
	StatusError   SemanticColor `theme:"status.error"`   // Errors, destructive

	// Interactive elements
	ButtonBg      SemanticColor `theme:"button.bg"`       // Unfocused button background
	ButtonFg      SemanticColor `theme:"button.fg"`       // Unfocused button foreground
	ButtonFocusBg SemanticColor `theme:"button.focus.bg"` // Focused button background
	ButtonFocusFg SemanticColor `theme:"button.focus.fg"` // Focused button foreground

	// Edit mode
	EditBg SemanticColor `theme:"edit.bg"` // Inline edit background
	EditFg SemanticColor `theme:"edit.fg"` // Inline edit foreground

	// Diff context (tints for code/content backgrounds)
	TintPositive SemanticColor `theme:"tint.positive"` // Added content background tint
	TintNegative SemanticColor `theme:"tint.negative"` // Removed content background tint
	TintNeutral  SemanticColor `theme:"tint.neutral"`  // Unchanged content background tint

	// Scrollbar
	ScrollTrack SemanticColor `theme:"scroll.track"` // Scrollbar track
	ScrollThumb SemanticColor `theme:"scroll.thumb"` // Scrollbar thumb

	// Syntax highlighting
	HighlightStyle string `theme:"highlight.style"` // Chroma style name (e.g. "monokai", "github")
}

// Palette combines system UI chrome colors with app-specific domain colors.
// T is the app's domain color struct (e.g., GomionColors).
type Palette[T any] struct {
	System SystemPalette
	App    T
}

// PaletteOpts controls palette construction behavior.
// Pass nil for default behavior (fixed colors).
type PaletteOpts struct {
	// Adaptive selects ANSI 0-15 colors for text, accents, and status
	// indicators. These colors respect the user's terminal theme.
	// When false (default), all colors use fixed ANSI 16-255 values.
	Adaptive bool
}

var (
	// ANSI 0-15 semantic color aliases — terminal-theme-adaptive.
	// These colors respect the user's terminal theme configuration.

	SemanticBlack         = NewSemanticColor(teacolor.Black)         // 0
	SemanticRed           = NewSemanticColor(teacolor.Red)           // 1
	SemanticGreen         = NewSemanticColor(teacolor.Green)         // 2
	SemanticYellow        = NewSemanticColor(teacolor.Yellow)        // 3
	SemanticBlue          = NewSemanticColor(teacolor.Blue)          // 4
	SemanticMagenta       = NewSemanticColor(teacolor.Magenta)       // 5
	SemanticCyan          = NewSemanticColor(teacolor.Cyan)          // 6
	SemanticWhite         = NewSemanticColor(teacolor.White)         // 7
	SemanticBrightBlack   = NewSemanticColor(teacolor.BrightBlack)   // 8
	SemanticBrightRed     = NewSemanticColor(teacolor.BrightRed)     // 9
	SemanticBrightGreen   = NewSemanticColor(teacolor.BrightGreen)   // 10
	SemanticBrightYellow  = NewSemanticColor(teacolor.BrightYellow)  // 11
	SemanticBrightBlue    = NewSemanticColor(teacolor.BrightBlue)    // 12
	SemanticBrightMagenta = NewSemanticColor(teacolor.BrightMagenta) // 13
	SemanticBrightCyan    = NewSemanticColor(teacolor.BrightCyan)    // 14
	SemanticBrightWhite   = NewSemanticColor(teacolor.BrightWhite)   // 15

	SemanticColor0   = NewSemanticColor(teacolor.Color0)   // black
	SemanticColor15  = NewSemanticColor(teacolor.Color15)  // white — items, buttons
	SemanticColor22  = NewSemanticColor(teacolor.Color22)  // dark green
	SemanticColor28  = NewSemanticColor(teacolor.Color28)  // dark green
	SemanticColor30  = NewSemanticColor(teacolor.Color30)  // dark cyan
	SemanticColor33  = NewSemanticColor(teacolor.Color33)  // medium blue
	SemanticColor39  = NewSemanticColor(teacolor.Color39)  // light blue — focused row
	SemanticColor46  = NewSemanticColor(teacolor.Color46)  // bright green — titles, active items
	SemanticColor51  = NewSemanticColor(teacolor.Color51)  // cyan — modal border
	SemanticColor52  = NewSemanticColor(teacolor.Color52)  // dark red
	SemanticColor62  = NewSemanticColor(teacolor.Color62)  // purple
	SemanticColor86  = NewSemanticColor(teacolor.Color86)  // cyan — key display, help text
	SemanticColor124 = NewSemanticColor(teacolor.Color124) // dark red
	SemanticColor130 = NewSemanticColor(teacolor.Color130) // dark gold
	SemanticColor153 = NewSemanticColor(teacolor.Color153) // light blue
	SemanticColor157 = NewSemanticColor(teacolor.Color157) // light green
	SemanticColor160 = NewSemanticColor(teacolor.Color160) // red
	SemanticColor166 = NewSemanticColor(teacolor.Color166) // dark orange
	SemanticColor178 = NewSemanticColor(teacolor.Color178) // gold — category headers
	SemanticColor214 = NewSemanticColor(teacolor.Color214) // orange
	SemanticColor217 = NewSemanticColor(teacolor.Color217) // light red/pink
	SemanticColor226 = NewSemanticColor(teacolor.Color226) // bright yellow
	SemanticColor230 = NewSemanticColor(teacolor.Color230) // bright yellow — selected text
	SemanticColor232 = NewSemanticColor(teacolor.Color232) // near-black
	SemanticColor238 = NewSemanticColor(teacolor.Color238) // dark gray
	SemanticColor240 = NewSemanticColor(teacolor.Color240) // dark gray — dividers
	SemanticColor244 = NewSemanticColor(teacolor.Color244) // medium gray — hints, cancel text
	SemanticColor248 = NewSemanticColor(teacolor.Color248) // lighter gray
	SemanticColor250 = NewSemanticColor(teacolor.Color250) // light gray
	SemanticColor252 = NewSemanticColor(teacolor.Color252) // very light gray
	SemanticColorNil = NewSemanticColor(nil)               // no color
)

// DarkSystemPalette returns a SystemPalette with colors suited for dark terminal
// backgrounds. When opts is nil or opts.Adaptive is false, all colors use fixed
// ANSI 16-255 values. When opts.Adaptive is true, text, accents, and status
// indicators use ANSI 0-15 colors that respect the user's terminal theme.
func DarkSystemPalette(opts *PaletteOpts) SystemPalette {
	if opts != nil && opts.Adaptive {
		return darkAdaptivePalette()
	}
	return darkFixedPalette()
}

func darkFixedPalette() SystemPalette {
	return SystemPalette{
		// Text hierarchy
		TextPrimary:   SemanticColor15,  // white — items, buttons
		TextSecondary: SemanticColor252, // light gray — labels, descriptions
		TextMuted:     SemanticColor244, // medium gray — hints, cancel text
		TextDim:       SemanticColor240, // dark gray — disabled, separators

		// Accent colors
		Accent:       SemanticColor46,  // bright green — titles, active items
		AccentAlt:    SemanticColor178, // gold — category headers
		AccentSubtle: SemanticColor86,  // cyan — key display, help text

		// Selection and focus
		SelectionBg: SemanticColor62,  // purple — selected background
		SelectionFg: SemanticColor230, // bright yellow — selected text
		FocusBorder: SemanticColor51,  // cyan — focused modal border
		FocusBg:     SemanticColor39,  // light blue — focused row

		// Borders and chrome
		Border:       SemanticColor240, // dark gray — default borders
		BorderAccent: SemanticColor51,  // cyan — modal borders
		Separator:    SemanticColor240, // dark gray — dividers

		// Status indicators
		StatusSuccess: SemanticColor46,  // bright green
		StatusInfo:    SemanticColor86,  // cyan
		StatusWarn:    SemanticColor214, // orange
		StatusError:   SemanticColor160, // red

		// Interactive elements
		ButtonBg:      SemanticColorNil, // no bg for unfocused
		ButtonFg:      SemanticColor15,  // white
		ButtonFocusBg: SemanticColor62,  // purple
		ButtonFocusFg: SemanticColor230, // bright yellow

		// Edit mode
		EditBg: SemanticColor226, // bright yellow
		EditFg: SemanticColor232, // near-black

		// Diff context
		TintPositive: SemanticColor22,  // dark green
		TintNegative: SemanticColor52,  // dark red
		TintNeutral:  SemanticColorNil, // no tint

		// Scrollbar
		ScrollTrack: SemanticColor240, // dark gray
		ScrollThumb: SemanticColor248, // lighter gray

		// Syntax highlighting
		HighlightStyle: "monokai",
	}
}

// LightSystemPalette returns a SystemPalette with colors suited for light
// terminal backgrounds. When opts is nil or opts.Adaptive is false, all colors
// use fixed ANSI 16-255 values. When opts.Adaptive is true, text, accents, and
// status indicators use ANSI 0-15 colors that respect the user's terminal theme.
func LightSystemPalette(opts *PaletteOpts) SystemPalette {
	if opts != nil && opts.Adaptive {
		return lightAdaptivePalette()
	}
	return lightFixedPalette()
}

func lightFixedPalette() SystemPalette {
	return SystemPalette{
		// Text hierarchy
		TextPrimary:   SemanticColor0,   // black
		TextSecondary: SemanticColor238, // dark gray
		TextMuted:     SemanticColor244, // medium gray
		TextDim:       SemanticColor250, // light gray

		// Accent colors
		Accent:       SemanticColor28,  // dark green
		AccentAlt:    SemanticColor130, // dark gold
		AccentSubtle: SemanticColor30,  // dark cyan

		// Selection and focus
		SelectionBg: SemanticColor62,  // purple
		SelectionFg: SemanticColor230, // bright yellow
		FocusBorder: SemanticColor33,  // medium blue
		FocusBg:     SemanticColor153, // light blue

		// Borders and chrome
		Border:       SemanticColor250, // light gray
		BorderAccent: SemanticColor33,  // medium blue
		Separator:    SemanticColor250, // light gray

		// Status indicators
		StatusSuccess: SemanticColor28,  // dark green
		StatusInfo:    SemanticColor30,  // dark cyan
		StatusWarn:    SemanticColor166, // dark orange
		StatusError:   SemanticColor124, // dark red

		// Interactive elements
		ButtonBg:      SemanticColorNil, // no bg for unfocused
		ButtonFg:      SemanticColor0,   // black
		ButtonFocusBg: SemanticColor62,  // purple
		ButtonFocusFg: SemanticColor230, // bright yellow

		// Edit mode
		EditBg: SemanticColor226, // bright yellow
		EditFg: SemanticColor0,   // black

		// Diff context
		TintPositive: SemanticColor157, // light green
		TintNegative: SemanticColor217, // light red/pink
		TintNeutral:  SemanticColorNil, // no tint

		// Scrollbar
		ScrollTrack: SemanticColor252, // very light gray
		ScrollThumb: SemanticColor244, // medium gray

		// Syntax highlighting
		HighlightStyle: "github",
	}
}

// AdaptiveSystemPalette detects the terminal background and returns the
// appropriate palette. Falls back to DarkSystemPalette if detection fails.
func AdaptiveSystemPalette(opts *PaletteOpts) SystemPalette {
	if lipgloss.HasDarkBackground(os.Stdin, os.Stderr) {
		return DarkSystemPalette(opts)
	}
	return LightSystemPalette(opts)
}

// DefaultSystemPalette returns AdaptiveSystemPalette(opts).
func DefaultSystemPalette(opts *PaletteOpts) SystemPalette {
	return AdaptiveSystemPalette(opts)
}

func darkAdaptivePalette() SystemPalette {
	p := darkFixedPalette()

	// 12 fields → ANSI 0-15 (adaptive)
	p.TextPrimary = SemanticBrightWhite   // 15
	p.TextSecondary = SemanticWhite       // 7
	p.TextMuted = SemanticBrightBlack     // 8
	p.TextDim = SemanticBrightBlack       // 8
	p.Accent = SemanticBrightGreen        // 10
	p.AccentAlt = SemanticBrightYellow    // 11
	p.AccentSubtle = SemanticBrightCyan   // 14
	p.StatusSuccess = SemanticBrightGreen // 10
	p.StatusInfo = SemanticBrightCyan     // 14
	p.StatusWarn = SemanticBrightYellow   // 11
	p.StatusError = SemanticBrightRed     // 9
	p.ButtonFg = SemanticBrightWhite      // 15

	return p
}

func lightAdaptivePalette() SystemPalette {
	p := lightFixedPalette()

	// 12 fields → ANSI 0-15 (adaptive)
	p.TextPrimary = SemanticBlack         // 0
	p.TextSecondary = SemanticBrightBlack // 8
	p.TextMuted = SemanticBrightBlack     // 8
	p.TextDim = SemanticWhite             // 7
	p.Accent = SemanticGreen              // 2
	p.AccentAlt = SemanticYellow          // 3
	p.AccentSubtle = SemanticCyan         // 6
	p.StatusSuccess = SemanticGreen       // 2
	p.StatusInfo = SemanticCyan           // 6
	p.StatusWarn = SemanticYellow         // 3
	p.StatusError = SemanticRed           // 1
	p.ButtonFg = SemanticBlack            // 0

	return p
}
