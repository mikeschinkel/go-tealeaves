package teanotify

import "errors"

var (
	// ErrNotify is the layer sentinel for all teanotify errors.
	ErrNotify = errors.New("notify")

	// ErrInvalidColor indicates a bad hex color in a notice definition.
	ErrInvalidColor = errors.New("invalid color")

	// ErrInvalidNoticeKey indicates an empty or duplicate notice key.
	ErrInvalidNoticeKey = errors.New("invalid notice key")

	// ErrInvalidWidth indicates a width <= 0.
	ErrInvalidWidth = errors.New("invalid width")

	// ErrInvalidDuration indicates a duration <= 0.
	ErrInvalidDuration = errors.New("invalid duration")
)
