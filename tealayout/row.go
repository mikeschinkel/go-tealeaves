package tealayout

// row is a horizontal container that distributes width among children.
type row struct {
	children     []child
	gap          int
	width        int
	height       int
	sizes        []int  // resolved widths per child
	active       []bool // which children are active (not removed by optional)
	resolved     bool
	dirty        bool
	cachedOutput string
}

// newRow creates a horizontal layout container.
func newRow(children ...child) *row {
	return &row{
		children: children,
	}
}

// SetSize sets the total available dimensions for this row.
func (r *row) SetSize(width, height int) {
	if r.width != width || r.height != height {
		r.width = width
		r.height = height
		r.resolved = false
		r.dirty = true
		r.cachedOutput = ""
	}
}

// MarkDirty forces re-resolution and re-rendering on the next Render() call.
func (r *row) MarkDirty() {
	r.dirty = true
	r.resolved = false
	r.cachedOutput = ""
}

// Resolve runs the layout algorithm and returns the assigned width for each
// child. Children removed by optional collapse have size 0.
func (r *row) Resolve() ([]int, error) {
	if r.resolved {
		return r.sizes, nil
	}

	n := len(r.children)
	constraints := make([]constraint, n)
	hinters := make([]SizeHinter, n)
	for i, ch := range r.children {
		constraints[i] = ch.Constraint
		if h, ok := ch.Widget.(SizeHinter); ok {
			hinters[i] = h
		}
	}

	sizes, err := resolveLinear(r.width, constraints, r.gap, hinters)
	if err != nil {
		return nil, err
	}

	r.sizes = sizes
	r.active = make([]bool, n)
	for i := range n {
		r.active[i] = sizes[i] > 0 || !constraints[i].optional
	}
	r.resolved = true
	return sizes, nil
}

// View implements Viewer by calling Render and discarding the error.
func (r *row) View() string {
	s, _ := r.Render()
	return s
}

// Direction returns Horizontal.
func (r *row) Direction() Direction {
	return Horizontal
}

// ChildRect returns the positioned rect for the i-th child after resolution.
// Must call Resolve() first.
func (r *row) ChildRect(i int) Rect {
	if !r.resolved || i < 0 || i >= len(r.sizes) {
		return Rect{}
	}
	x := 0
	for j := range i {
		if r.sizes[j] > 0 {
			x += r.sizes[j]
			if j < i-1 || (j == i-1 && r.sizes[i] > 0) {
				x += r.gap
			}
		}
	}
	return Rect{
		X:      x,
		Y:      0,
		Width:  r.sizes[i],
		Height: r.height,
	}
}
