// Package teadrpdwn provides a Bubble Tea v2 component for popup dropdown
// selection menus.
//
// A dropdown appears at a specified screen position (relative to a field),
// supports keyboard navigation, and emits selection messages.
//
// Usage:
//
//	dd := teadrpdwn.NewDropdownModel(options, &teadrpdwn.DropdownModelArgs{
//	    FieldRow: 5, FieldCol: 10,
//	    ScreenWidth: 80, ScreenHeight: 24,
//	})
//	dd, cmd = dd.Open()
//
// # Stability
//
// This package is provisional as of v0.3.0. The public API may change in
// minor releases until promoted to stable.
package teadrpdwn
