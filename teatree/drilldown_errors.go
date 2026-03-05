package teatree

import "errors"

var (
	ErrDrillDown    = errors.New("drilldown error")
	ErrInvalidNode  = errors.New("invalid node error")
	ErrEmptyPath    = errors.New("empty path error")
	ErrInvalidLevel = errors.New("invalid level error")
)
