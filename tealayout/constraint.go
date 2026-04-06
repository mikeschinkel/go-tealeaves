package tealayout

// constraintKind identifies which constraint type is active.
type constraintKind int

const (
	// constraintFixed allocates exactly fixedSize cells.
	constraintFixed constraintKind = iota

	// constraintFlex distributes remaining space proportionally by weight.
	constraintFlex

	// constraintFit measures the widget via SizeHinter and uses its desired size.
	constraintFit
)

// constraint describes how a child should be sized along a container's axis.
type constraint struct {
	kind       constraintKind
	fixedSize  int
	flexWeight float64
	minSize    int  // 0 = no minimum
	maxSize    int  // -1 = unbounded (default)
	optional   bool
	minSizeFit bool // when true, query SizeHinter.Min as effective minimum
}
