package teatxtsnip

import "errors"

var (
	// ErrTextSnip represents text snippet errors
	ErrTextSnip = errors.New("text snip")

	// ErrInvalidSelection represents invalid selection range
	ErrInvalidSelection = errors.New("invalid selection")
)
