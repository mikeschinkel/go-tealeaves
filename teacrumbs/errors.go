package teacrumbs

import "errors"

var (
	// ErrBreadcrumbs represents breadcrumb errors.
	ErrBreadcrumbs = errors.New("breadcrumbs")

	// ErrIndexOutOfRange represents an index out of range error.
	ErrIndexOutOfRange = errors.New("index out of range")
)
