package teamodal

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// OverlayModal overlays modal view on background view at specified position.
// Uses ANSI-aware string operations to correctly handle styled text.
// This follows the proven pattern from teadrpdwn (same as OverlayDropdown).
//
// Parameters:
//   - background: The base view (fully rendered string with ANSI codes)
//   - foreground: The modal view (fully rendered string with ANSI codes)
//   - row: Line number in background where foreground row 0 should appear (0-indexed)
//   - col: Display column in background where foreground col 0 should appear (0-indexed)
//
// Returns:
//   - Composited view with foreground overlaid on background
func OverlayModal(background, foreground string, row, col int) string {
	var result strings.Builder

	bgLines := strings.Split(background, "\n")
	fgLines := strings.Split(foreground, "\n")

	for i, bgLine := range bgLines {
		fgRow := i - row

		// This line has no foreground overlay
		if fgRow < 0 || fgRow >= len(fgLines) {
			result.WriteString(bgLine)
			result.WriteString("\n")
			continue
		}

		// Overlay foreground line onto background line
		fgLine := fgLines[fgRow]
		composited := overlayLine(bgLine, fgLine, col)
		result.WriteString(composited)
		result.WriteString("\n")
	}

	// Remove trailing newline
	output := result.String()
	if len(output) > 0 && output[len(output)-1] == '\n' {
		output = output[:len(output)-1]
	}

	return output
}

// overlayLine overlays foreground onto background at column position (ANSI-aware).
// This follows the pattern from teadrpdwn: split into left + overlay + right.
//
// The key insight: Standard Go string operations (len, slicing) count ANSI escape
// codes as characters, which breaks positioning. We use ansi.StringWidth() for
// visual width and ansi.Truncate/TruncateLeft for ANSI-safe string cutting.
func overlayLine(background, foreground string, col int) string {
	if col < 0 {
		col = 0
	}

	bgWidth := ansi.StringWidth(background)
	fgWidth := ansi.StringWidth(foreground)

	// Build: left part of background + foreground + right part of background
	// This is the same pattern as teadrpdwn's overlayLine()
	var result strings.Builder

	// Left part: truncate background to col width
	if col > 0 {
		if col <= bgWidth {
			left := ansi.Truncate(background, col, "")
			result.WriteString(left)
		} else {
			// Need padding beyond background width
			result.WriteString(background)
			result.WriteString(strings.Repeat(" ", col-bgWidth))
		}
	}

	// Middle part: foreground content (the overlay)
	result.WriteString(foreground)

	// Right part: remainder of background after foreground
	endCol := col + fgWidth
	if endCol < bgWidth {
		// TruncateLeft(s, n) skips the first n display columns
		// We want to skip everything up to endCol
		remaining := ansi.TruncateLeft(background, endCol, "")
		result.WriteString(remaining)
	}

	return result.String()
}
