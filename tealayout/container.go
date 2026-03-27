package tealayout

// container is the interface implemented by row, column, and Component.
// It allows the Layout type to work with any container type polymorphically.
type container interface {
	SetSizer
	Viewer

	// Resolve runs the layout algorithm and returns sizes along the
	// container's axis.
	Resolve() ([]int, error)

	// Render resolves, sets child sizes, renders children, and composes output.
	Render() (string, error)

	// MarkDirty forces re-resolution and re-rendering.
	MarkDirty()

	// Direction returns the container's layout axis.
	Direction() Direction

	// ChildRect returns the positioned rect for the i-th child.
	ChildRect(i int) Rect
}

// Verify interface compliance at compile time.
var (
	_ container = (*row)(nil)
	_ container = (*column)(nil)
)
