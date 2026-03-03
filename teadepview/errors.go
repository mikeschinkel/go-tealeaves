package teadepview

import "errors"

var (
	ErrDependency   = errors.New("dependency error")
	ErrInvalidNode  = errors.New("invalid node error")
	ErrEmptyPath    = errors.New("empty path error")
	ErrInvalidLevel = errors.New("invalid level error")
)
