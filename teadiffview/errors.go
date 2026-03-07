package teadiffview

import "errors"

var (
	// ErrDiffView is the layer sentinel for all teadiffview errors.
	ErrDiffView = errors.New("diffview")

	// ErrInvalidFile indicates an invalid file diff (e.g., nil or missing path).
	ErrInvalidFile = errors.New("invalid file")

	// ErrInvalidBlock indicates a malformed condensed block.
	ErrInvalidBlock = errors.New("invalid block")

	// ErrEmptyDiff indicates an empty diff with no files.
	ErrEmptyDiff = errors.New("empty diff")

	// ErrInvalidContent indicates invalid DiffContent data.
	ErrInvalidContent = errors.New("invalid content")
)
