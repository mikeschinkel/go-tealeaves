package tealayout

import "errors"

var (
	// ErrZeroDimensions represents a layout with zero width or height.
	ErrZeroDimensions = errors.New("zero dimensions")

	// ErrDuplicatePaneName is returned when two panes share the same name.
	ErrDuplicatePaneName = errors.New("duplicate pane name")

	// ErrPaneNotFound is returned when a named pane lookup fails.
	ErrPaneNotFound = errors.New("pane not found")

	// ErrStackEmpty is returned when popping from an empty stack.
	ErrStackEmpty = errors.New("stack empty")

	// ErrStackUnderflow is returned when popping below the minimum stack depth.
	ErrStackUnderflow = errors.New("stack underflow")
)
