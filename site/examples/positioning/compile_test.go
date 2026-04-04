package examples_test

import (
	"os"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teafields"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TestCompile_PositioningQuickExample verifies the quick example from positioning.mdx.
func TestCompile_PositioningQuickExample(t *testing.T) {
	renderedModal := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render("Modal content")

	// Measure a rendered view
	width, height := teautils.MeasureRenderedView(renderedModal)
	_, _ = width, height

	// Calculate centered position
	row, col := teautils.CalculateCenter(80, 24, width, height)
	_, _ = row, col

	// Combined
	w, h, r, c := teautils.CenterModal(renderedModal, 80, 24)
	_, _, _, _ = w, h, r, c
}

// TestCompile_EnsureTermGetSize verifies EnsureTermGetSize from positioning.mdx.
func TestCompile_EnsureTermGetSize(t *testing.T) {
	w, h, ok := teafields.EnsureTermGetSize(os.Stdout.Fd())
	_, _, _ = w, h, ok
}
