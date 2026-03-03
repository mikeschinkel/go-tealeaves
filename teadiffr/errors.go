package teadiffr

import "errors"

var (
	// ErrDiff is the layer sentinel for all teadiffr errors.
	ErrDiff = errors.New("diff")

	// ErrInvalidFile indicates an invalid file diff (e.g., nil or missing path).
	ErrInvalidFile = errors.New("invalid file")

	// ErrInvalidBlock indicates a malformed condensed block.
	ErrInvalidBlock = errors.New("invalid block")

	// ErrEmptyDiff indicates an empty diff with no files.
	ErrEmptyDiff = errors.New("empty diff")
)
