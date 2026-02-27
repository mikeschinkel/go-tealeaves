package teamodal

import "errors"

var (
	ErrModal         = errors.New("modal error")
	ErrInvalidBounds = errors.New("invalid bounds error")
	ErrCancelled     = errors.New("cancelled error")
)
