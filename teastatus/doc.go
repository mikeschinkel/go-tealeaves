// Package teastatus provides a Bubble Tea v2 component for rendering
// a two-zone status bar at the bottom of the terminal.
//
// The left zone displays key-action menu items (e.g., "[?] Help  [tab] Switch"),
// and the right zone displays text indicators (e.g., "DEPS IN-FLUX | 3 batches").
//
// Usage:
//
//	sb := teastatus.NewStatusBarModel().
//	    SetMenuItems(items).
//	    SetIndicators(indicators).
//	    SetSize(width)
//
//	// In View():
//	sb.View()
//
// # Stability
//
// This package is provisional as of v0.3.0. The public API may change in
// minor releases until promoted to stable.
package teastatus
