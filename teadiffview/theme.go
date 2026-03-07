package teadiffview

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// NewThemedTUIRenderer creates a TUIRenderer with colors derived from the
// given theme's palette. This is an alternative to NewTUIRenderer when using
// the theming system.
func NewThemedTUIRenderer(theme teautils.Theme) *TUIRenderer {
	p := theme.System
	return &TUIRenderer{
		FileHeaderColor:    p.TextMuted,
		BlockHeaderColor:   p.AccentAlt,
		ContextColor:       p.TextDim,
		AddedColor:         p.StatusSuccess,
		DeletedColor:       p.StatusError,
		NewStatusColor:     p.StatusSuccess,
		DeletedStatusColor: p.StatusError,
		NewBgColor:         p.TintPositive,
		DeletedBgColor:     p.TintNegative,
	}
}

// CommitGroupColors defines colors for commit group gutter indicators.
var CommitGroupColors = []color.Color{
	lipgloss.Color("#4CAF50"), // Group 1: Green
	lipgloss.Color("#2196F3"), // Group 2: Blue
	lipgloss.Color("#FF9800"), // Group 3: Orange
	lipgloss.Color("#9C27B0"), // Group 4: Purple
	lipgloss.Color("#00BCD4"), // Group 5: Cyan
	lipgloss.Color("#F44336"), // Group 6: Red
	lipgloss.Color("#FFEB3B"), // Group 7: Yellow
	lipgloss.Color("#795548"), // Group 8: Brown
}

// SingleCommitColor is used for 'c' gutter indicator.
var SingleCommitColor = lipgloss.Color("#4CAF50") // Green

// SelectionBgColor is the background color for selected rows (256-palette).
const SelectionBgColor = "7" // Light gray

// ChangeBlockBgColor is the background color for change block rows (256-palette).
const ChangeBlockBgColor = "24" // Dim cyan

// ANSI escape sequences for background colors (used in string replacements).
var (
	SelectionBgANSI   = "\x1b[48;5;" + SelectionBgColor + "m"
	ChangeBlockBgANSI = "\x1b[48;5;" + ChangeBlockBgColor + "m"
	ANSIReset         = "\x1b[0m"
)
