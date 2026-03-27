package tealayout

// dimensionKind identifies which dimension type is active.
type dimensionKind int

const (
	dimensionPercent dimensionKind = iota
	dimensionFixed
	dimensionFlex
	dimensionFit
)

// Dimension specifies how a component is sized within its parent.
type Dimension struct {
	kind  dimensionKind
	value float64 // percentage (0-100), fixed cells, or flex weight
}

// Percent creates a Dimension sized as a percentage of the parent.
// Siblings' percentages are treated as flex weights — they don't need
// to sum to 100, but doing so makes the layout intuitive.
func Percent(n float64) Dimension {
	return Dimension{kind: dimensionPercent, value: n}
}

// Fixed creates a Dimension for exactly n cells along the parent's axis.
func Fixed(n int) Dimension {
	return Dimension{kind: dimensionFixed, value: float64(n)}
}

// Flex creates a Dimension that distributes remaining space proportionally
// by the given weight. For power users who need non-percentage ratios
// (e.g., golden ratio 1.0 : 1.618).
func Flex(weight float64) Dimension {
	return Dimension{kind: dimensionFlex, value: weight}
}

// Fit creates a Dimension that measures the widget via SizeHinter and
// uses its desired size.
func Fit() Dimension {
	return Dimension{kind: dimensionFit}
}

// Common Dimension constants for frequently used percentages.
var (
	Percent100 = Percent(100)
	Percent75  = Percent(75)
	Percent50  = Percent(50)
	Percent33  = Percent(33)
	Percent25  = Percent(25)
	Percent20  = Percent(20)
)
