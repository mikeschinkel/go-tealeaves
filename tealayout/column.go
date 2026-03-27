package tealayout

// column is a vertical container that distributes height among children.
type column struct {
	children     []child
	gap          int
	width        int
	height       int
	sizes        []int  // resolved heights per child
	active       []bool // which children are active
	resolved     bool
	dirty        bool
	cachedOutput string
}

// newColumn creates a vertical layout container.
func newColumn(children ...child) *column {
	return &column{
		children: children,
	}
}

// SetSize sets the total available dimensions for this column.
func (c *column) SetSize(width, height int) {
	if c.width != width || c.height != height {
		c.width = width
		c.height = height
		c.resolved = false
		c.dirty = true
		c.cachedOutput = ""
	}
}

// MarkDirty forces re-resolution and re-rendering on the next Render() call.
func (c *column) MarkDirty() {
	c.dirty = true
	c.resolved = false
	c.cachedOutput = ""
}

// Resolve runs the layout algorithm and returns the assigned height for each
// child. Children removed by optional collapse have size 0.
func (c *column) Resolve() ([]int, error) {
	if c.resolved {
		return c.sizes, nil
	}

	n := len(c.children)
	constraints := make([]constraint, n)
	hinters := make([]SizeHinter, n)
	for i, ch := range c.children {
		constraints[i] = ch.Constraint
		if h, ok := ch.Widget.(SizeHinter); ok {
			hinters[i] = h
		}
	}

	sizes, err := resolveLinear(c.height, constraints, c.gap, hinters)
	if err != nil {
		return nil, err
	}

	c.sizes = sizes
	c.active = make([]bool, n)
	for i := range n {
		c.active[i] = sizes[i] > 0 || !constraints[i].optional
	}
	c.resolved = true
	return sizes, nil
}

// View implements Viewer by calling Render and discarding the error.
func (c *column) View() string {
	s, _ := c.Render()
	return s
}

// Direction returns Vertical.
func (c *column) Direction() Direction {
	return Vertical
}

// ChildRect returns the positioned rect for the i-th child after resolution.
// Must call Resolve() first.
func (c *column) ChildRect(i int) Rect {
	if !c.resolved || i < 0 || i >= len(c.sizes) {
		return Rect{}
	}
	y := 0
	for j := range i {
		if c.sizes[j] > 0 {
			y += c.sizes[j]
			if j < i-1 || (j == i-1 && c.sizes[i] > 0) {
				y += c.gap
			}
		}
	}
	return Rect{
		X:      0,
		Y:      y,
		Width:  c.width,
		Height: c.sizes[i],
	}
}
