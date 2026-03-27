package tealayout

// Size represents a width/height pair in terminal cells.
type Size struct {
	Width  int
	Height int
}

// SizeHint describes a widget's size preferences for a given available space.
type SizeHint struct {
	Min     Size // Minimum usable size
	Desired Size // Preferred size
	Max     Size // Maximum useful size; -1 per dimension = unbounded
}
