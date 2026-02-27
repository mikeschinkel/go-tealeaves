package teautils

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// CalculateCenter computes the row and column to center a modal of given dimensions
// within a screen of given dimensions. Results are clamped to 0.
func CalculateCenter(screenW, screenH, modalW, modalH int) (row, col int) {
	row = (screenH - modalH) / 2
	if row < 0 {
		row = 0
	}

	col = (screenW - modalW) / 2
	if col < 0 {
		col = 0
	}

	return row, col
}

// MeasureRenderedView measures the width and height of a rendered view string.
// Width is the maximum ANSI-aware line width, height is the number of lines.
func MeasureRenderedView(renderedView string) (width, height int) {
	lines := strings.Split(renderedView, "\n")
	height = len(lines)

	for _, line := range lines {
		lineWidth := ansi.StringWidth(line)
		if lineWidth > width {
			width = lineWidth
		}
	}

	return width, height
}

// CenterModal measures a rendered view and computes its centered position.
// Returns the modal dimensions and the row/col for overlay positioning.
// The row is shifted up by 1 (slightly above center looks better).
func CenterModal(renderedView string, screenW, screenH int) (width, height, row, col int) {
	width, height = MeasureRenderedView(renderedView)
	row, col = CalculateCenter(screenW, screenH, width, height)

	// Shift up one line (slightly above center looks better)
	if row > 0 {
		row--
	}

	return width, height, row, col
}
