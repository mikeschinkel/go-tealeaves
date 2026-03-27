package tealayout

// Element is anything that can be placed in a layout — a widget or a nested *Component.
type Element interface{}

// Component is a sized layout node with optional modifiers. It is the primary
// public type for constructing layout trees. Create with NewRow or NewColumn.
type Component struct {
	dim       Dimension
	elements  []Element
	direction Direction
	minSize   int
	maxSize   int // -1 = unbounded
	optional  bool
	gap       int

	// Lazily built internal container.
	inner container
	built bool
}

// NewRow creates a horizontal Component with the given dimension and elements.
// Elements can be widgets (any type implementing SetSizer/Viewer) or nested *Component.
func NewRow(dim Dimension, elements ...Element) *Component {
	return &Component{
		dim:       dim,
		elements:  elements,
		direction: Horizontal,
		maxSize:   -1,
	}
}

// NewColumn creates a vertical Component with the given dimension and elements.
func NewColumn(dim Dimension, elements ...Element) *Component {
	return &Component{
		dim:       dim,
		elements:  elements,
		direction: Vertical,
		maxSize:   -1,
	}
}

// WithMinSize returns the Component with the minimum size set.
// The component will not be assigned fewer than n cells.
func (c *Component) WithMinSize(n int) *Component {
	c.minSize = n
	return c
}

// WithMaxSize returns the Component with the maximum size set.
// Use -1 for unbounded.
func (c *Component) WithMaxSize(n int) *Component {
	c.maxSize = n
	return c
}

// WithOptional returns the Component with the optional flag set.
// When true and the assigned size falls below MinSize during resolution,
// the component is removed from the layout entirely.
func (c *Component) WithOptional(b bool) *Component {
	c.optional = b
	return c
}

// WithGap sets the spacing (in cells/lines) between children.
func (c *Component) WithGap(n int) *Component {
	c.gap = n
	return c
}

// toConstraint converts this Component's Dimension and modifiers into
// an internal constraint for use in the resolution algorithm.
func (c *Component) toConstraint() constraint {
	cs := constraint{
		minSize:  c.minSize,
		maxSize:  c.maxSize,
		optional: c.optional,
	}
	switch c.dim.kind {
	case dimensionPercent:
		cs.kind = constraintFlex
		cs.flexWeight = c.dim.value
	case dimensionFixed:
		cs.kind = constraintFixed
		cs.fixedSize = int(c.dim.value)
	case dimensionFlex:
		cs.kind = constraintFlex
		cs.flexWeight = c.dim.value
	case dimensionFit:
		cs.kind = constraintFit
	}
	return cs
}

// ensureInner lazily builds the internal container from elements.
func (c *Component) ensureInner() {
	if c.built {
		return
	}
	children := make([]child, 0, len(c.elements))
	for _, elem := range c.elements {
		if elem == nil {
			continue
		}
		switch e := elem.(type) {
		case *Component:
			children = append(children, child{
				Widget:     e,
				Constraint: e.toConstraint(),
			})
		default:
			children = append(children, child{
				Widget:     e,
				Constraint: constraint{kind: constraintFlex, flexWeight: 1.0, maxSize: -1},
			})
		}
	}
	switch c.direction {
	case Horizontal:
		r := newRow(children...)
		r.gap = c.gap
		c.inner = r
	default:
		col := newColumn(children...)
		col.gap = c.gap
		c.inner = col
	}
	c.built = true
}

// SetSize sets the total available dimensions for this component.
func (c *Component) SetSize(width, height int) {
	c.ensureInner()
	c.inner.SetSize(width, height)
}

// View implements Viewer by calling Render and discarding the error.
func (c *Component) View() string {
	c.ensureInner()
	return c.inner.View()
}

// Resolve runs the layout algorithm and returns the assigned sizes.
func (c *Component) Resolve() ([]int, error) {
	c.ensureInner()
	return c.inner.Resolve()
}

// Render resolves the layout and returns the composed output string.
func (c *Component) Render() (string, error) {
	c.ensureInner()
	return c.inner.Render()
}

// MarkDirty forces re-resolution and re-rendering on the next call.
func (c *Component) MarkDirty() {
	if c.inner != nil {
		c.inner.MarkDirty()
	}
}

// Direction returns the component's layout axis.
func (c *Component) Direction() Direction {
	return c.direction
}

// ChildRect returns the positioned rect for the i-th child after resolution.
func (c *Component) ChildRect(i int) Rect {
	c.ensureInner()
	return c.inner.ChildRect(i)
}

// Verify interface compliance at compile time.
var _ container = (*Component)(nil)
