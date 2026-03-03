package teadrpdwn

// Position represents field position for testing permutations
type Position int

const (
	TopLeft Position = iota
	TopMiddle
	TopRight
	MiddleLeft
	Middle
	MiddleRight
	BottomLeft
	BottomMiddle
	BottomRight
)
