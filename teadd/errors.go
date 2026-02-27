package teadd

import "errors"

var (
	ErrDropdown      = errors.New("dropdown error")
	ErrInvalidIndex  = errors.New("invalid index error")
	ErrEmptyOptions  = errors.New("empty options error")
	ErrInvalidBounds = errors.New("invalid bounds error")
	ErrInvalidRow    = errors.New("invalid row error")
	ErrInvalidCol    = errors.New("invalid column error")
)
