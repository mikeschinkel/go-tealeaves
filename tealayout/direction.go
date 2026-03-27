package tealayout

// Direction represents the axis along which a container distributes space.
type Direction int

const (
	// Horizontal distributes width (Row).
	Horizontal Direction = iota
	// Vertical distributes height (Column).
	Vertical
)
