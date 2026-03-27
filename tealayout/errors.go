package tealayout

import "errors"

var (
	// ErrZeroDimensions represents a layout with zero width or height.
	ErrZeroDimensions = errors.New("zero dimensions")
)
