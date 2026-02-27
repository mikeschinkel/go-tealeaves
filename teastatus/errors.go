package teastatus

import "errors"

var (
	// ErrStatusBar represents status bar errors
	ErrStatusBar = errors.New("status bar")

	// ErrInvalidWidth represents invalid width errors
	ErrInvalidWidth = errors.New("invalid width")
)
