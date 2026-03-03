package teagrid

import "errors"

var (
	// ErrGrid represents grid errors
	ErrGrid = errors.New("grid")

	// ErrInvalidColumn represents invalid column configuration
	ErrInvalidColumn = errors.New("invalid column")

	// ErrInvalidRow represents invalid row data
	ErrInvalidRow = errors.New("invalid row")
)
