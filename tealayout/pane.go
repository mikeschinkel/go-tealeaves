package tealayout

import lipgloss "charm.land/lipgloss/v2"

// paneChild pairs an element with a layout constraint.
type paneChild struct {
	elem       element
	constraint constraint
}

// Pane is a sized layout node with optional modifiers. It is the primary
// public type for constructing layout trees. Create with NewRow or NewColumn.
//
// Pane implements the unexported element interface, so panes nest directly
// inside other panes without wrapping.
type Pane struct {
	dim       Dimension
	direction Direction
	minSize      int
	maxSize      int // -1 = unbounded
	optional     bool
	gap          int
	minFlexWeight float64 // floor for ResizeFocused (prevents weight-crushing)

	// Visibility, naming, alignment, focus, and tree structure.
	name       string
	visible    bool      // default true; user-controlled show/hide
	focusable  bool      // when true, FocusManager includes this pane in its order
	minSizeFit bool      // when true, use SizeHinter.Min as effective minimum
	alignment  Alignment // default TopLeft
	focused    bool      // set by FocusManager
	parent     *Pane

	// Flattened container state (absorbed from old row/column types).
	children     []paneChild
	width        int
	height       int
	sizes        []int  // resolved sizes along the pane's axis
	active       []bool // which children survived optional removal
	resolved     bool
	dirty        bool
	cachedOutput string
}

// NewRow creates a horizontal Pane with the given dimension and child elements.
// Elements can be *Element[T] (from NewElement) or nested *Pane.
func NewRow(dim Dimension, elements ...element) *Pane {
	p := &Pane{
		dim:       dim,
		direction: Horizontal,
		maxSize:   -1,
		visible:   true,
	}
	p.buildChildren(elements)
	return p
}

// NewColumn creates a vertical Pane with the given dimension and child elements.
func NewColumn(dim Dimension, elements ...element) *Pane {
	p := &Pane{
		dim:       dim,
		direction: Vertical,
		maxSize:   -1,
		visible:   true,
	}
	p.buildChildren(elements)
	return p
}

// buildChildren converts element arguments into paneChild entries and
// sets parent pointers on child Panes.
func (p *Pane) buildChildren(elements []element) {
	children := make([]paneChild, 0, len(elements))
	for _, elem := range elements {
		if elem == nil {
			continue
		}
		switch e := elem.(type) {
		case *Pane:
			e.parent = p
			children = append(children, paneChild{
				elem:       e,
				constraint: e.toConstraint(),
			})
		default:
			children = append(children, paneChild{
				elem:       e,
				constraint: constraint{kind: constraintFlex, flexWeight: 1.0, maxSize: -1},
			})
		}
	}
	p.children = children
}

// --- Builder methods (chainable) ---

// WithMinSize returns the Pane with the minimum size set.
// The pane will not be assigned fewer than n cells.
func (p *Pane) WithMinSize(n int) *Pane {
	p.minSize = n
	return p
}

// WithMaxSize returns the Pane with the maximum size set.
// Use -1 for unbounded.
func (p *Pane) WithMaxSize(n int) *Pane {
	p.maxSize = n
	return p
}

// WithOptional returns the Pane with the optional flag set.
// When true and the assigned size falls below MinSize during resolution,
// the pane is removed from the layout entirely.
func (p *Pane) WithOptional(b bool) *Pane {
	p.optional = b
	return p
}

// WithGap sets the spacing (in cells/lines) between children.
func (p *Pane) WithGap(n int) *Pane {
	p.gap = n
	return p
}

// WithMinFlexWeight sets the minimum flex weight for ResizeFocused.
// This prevents the pane from being weight-crushed even when the resolver's
// cell-based minSize would still technically allow it.
func (p *Pane) WithMinFlexWeight(w float64) *Pane {
	p.minFlexWeight = w
	return p
}

// MinFlexWeight returns the minimum flex weight floor for ResizeFocused.
func (p *Pane) MinFlexWeight() float64 {
	return p.minFlexWeight
}

// Dimension returns the pane's current dimension.
func (p *Pane) Dimension() Dimension {
	return p.dim
}

// WithName sets a name for programmatic access via Layout.Pane(name).
func (p *Pane) WithName(name string) *Pane {
	p.name = name
	return p
}

// Name returns the pane's name, or empty string if unnamed.
func (p *Pane) Name() string {
	return p.name
}

// WithFocusable marks this pane as focusable by the FocusManager.
// Only focusable panes are included in the focus order.
func (p *Pane) WithFocusable() *Pane {
	p.focusable = true
	return p
}

// WithoutFocusable removes this pane from the FocusManager's focus order.
func (p *Pane) WithoutFocusable() *Pane {
	p.focusable = false
	return p
}

// IsFocusable returns whether this pane participates in focus navigation.
func (p *Pane) IsFocusable() bool {
	return p.focusable
}

// WithMinSizeFit tells the layout engine to query the element's SizeHint.Min
// as the effective minimum size. The pane will not be allocated fewer cells
// than the widget's reported minimum, preventing content overflow.
func (p *Pane) WithMinSizeFit() *Pane {
	p.minSizeFit = true
	return p
}

// WithoutMinSizeFit removes the dynamic minimum size constraint.
func (p *Pane) WithoutMinSizeFit() *Pane {
	p.minSizeFit = false
	return p
}

// WithAlignment sets alignment for content within this pane.
// Composable: WithAlignment(Top) sets vertical without clearing horizontal.
// WithAlignment(TopLeft) sets both axes.
func (p *Pane) WithAlignment(a Alignment) *Pane {
	p.alignment = mergeAlignment(p.alignment, a)
	return p
}

// Visible returns whether this pane is visible.
func (p *Pane) Visible() bool {
	return p.visible
}

// SetVisible sets visibility and marks the tree dirty.
func (p *Pane) SetVisible(v bool) {
	if p.visible == v {
		return
	}
	p.visible = v
	p.markDirtyUp()
}

// Show is shorthand for SetVisible(true).
func (p *Pane) Show() {
	p.SetVisible(true)
}

// Hide is shorthand for SetVisible(false).
func (p *Pane) Hide() {
	p.SetVisible(false)
}

// Focused returns whether this pane currently has focus (set by FocusManager).
func (p *Pane) Focused() bool {
	return p.focused
}

// markDirtyUp marks this pane and all ancestors as dirty.
func (p *Pane) markDirtyUp() {
	p.MarkDirty()
	if p.parent != nil {
		p.parent.markDirtyUp()
	}
}

// --- Dimension mutation ---

// SetDimension updates the pane's dimension and rebuilds its constraint,
// marking the tree dirty. Use for runtime resize without tree rebuild.
func (p *Pane) SetDimension(dim Dimension) {
	p.dim = dim
	p.rebuildConstraintInParent()
	p.markDirtyUp()
}

// SetMinSize updates the minimum size at runtime.
func (p *Pane) SetMinSize(n int) {
	p.minSize = n
	p.rebuildConstraintInParent()
	p.markDirtyUp()
}

// SetMaxSize updates the maximum size at runtime. Use -1 for unbounded.
func (p *Pane) SetMaxSize(n int) {
	p.maxSize = n
	p.rebuildConstraintInParent()
	p.markDirtyUp()
}

// rebuildConstraintInParent updates this pane's constraint entry in its
// parent's children slice, so the next resolve picks up the change.
func (p *Pane) rebuildConstraintInParent() {
	if p.parent == nil {
		return
	}
	for i, ch := range p.parent.children {
		if ch.elem == p {
			p.parent.children[i].constraint = p.toConstraint()
			return
		}
	}
}

// --- Layout operations ---

// SetSize sets the total available dimensions for this pane.
func (p *Pane) SetSize(width, height int) {
	if p.width != width || p.height != height {
		p.width = width
		p.height = height
		p.resolved = false
		p.dirty = true
		p.cachedOutput = ""
	}
}

// MarkDirty forces re-resolution and re-rendering on the next call.
// It propagates downward to all children so nested caches are invalidated.
func (p *Pane) MarkDirty() {
	p.dirty = true
	p.resolved = false
	p.cachedOutput = ""
	for _, ch := range p.children {
		ch.elem.markDirty()
	}
}

// markDirty satisfies the element interface, delegating to MarkDirty.
func (p *Pane) markDirty() {
	p.MarkDirty()
}

// Direction returns the pane's layout axis.
func (p *Pane) Direction() Direction {
	return p.direction
}

// Resolve runs the layout algorithm and returns the assigned sizes.
// Hidden panes (visible == false) are excluded from resolution entirely
// so they don't consume gap space.
func (p *Pane) Resolve() ([]int, error) {
	if p.resolved {
		return p.sizes, nil
	}

	n := len(p.children)

	// Build a mapping of visible children only.
	visibleIdx := make([]int, 0, n)       // indices into p.children
	visConstraints := make([]constraint, 0, n)
	visHinters := make([]SizeHinter, 0, n)

	for i, ch := range p.children {
		if cp, ok := ch.elem.(*Pane); ok && !cp.visible {
			continue
		}
		visibleIdx = append(visibleIdx, i)
		visConstraints = append(visConstraints, ch.constraint)
		if _, ok := ch.elem.sizeHint(0, 0); ok {
			elem := ch.elem // capture for closure
			visHinters = append(visHinters, sizeHinterFunc(func(availW, availH int) SizeHint {
				h, _ := elem.sizeHint(availW, availH)
				return h
			}))
		} else {
			visHinters = append(visHinters, nil)
		}
	}

	available := p.axisSize()
	visSizes, err := resolveLinear(available, visConstraints, p.gap, visHinters, p.direction == Horizontal)
	if err != nil {
		return nil, err
	}

	// Map visible sizes back to full children array.
	sizes := make([]int, n)
	for j, idx := range visibleIdx {
		sizes[idx] = visSizes[j]
	}

	p.sizes = sizes
	p.active = make([]bool, n)
	for i := range n {
		if cp, ok := p.children[i].elem.(*Pane); ok && !cp.visible {
			p.active[i] = false
			continue
		}
		p.active[i] = sizes[i] > 0 || !p.children[i].constraint.optional
	}
	p.resolved = true
	return sizes, nil
}

// Render resolves the layout and returns the composed output string.
func (p *Pane) Render() (string, error) {
	if p.dirty || !p.resolved {
		p.resolved = false
	}
	if p.cachedOutput != "" && !p.dirty {
		return p.cachedOutput, nil
	}

	sizes, err := p.Resolve()
	if err != nil {
		return "", err
	}

	views := make([]string, 0, len(p.children))
	for i, ch := range p.children {
		if sizes[i] <= 0 {
			continue
		}
		childW, childH := p.childDimensions(sizes[i])
		setChildSizeViaElement(ch.elem, childW, childH)
		v := contentChildElement(ch.elem, childW, childH)
		// Apply alignment if the child pane has one set.
		if cp, ok := ch.elem.(*Pane); ok && cp.alignment != 0 {
			v = alignContent(v, childW, childH, cp.alignment)
		}
		views = append(views, v)
	}

	var output string
	if p.direction == Horizontal {
		output = joinHorizontal(views, p.gap, p.height)
	} else {
		output = joinVertical(views, p.gap)
	}
	p.cachedOutput = output
	p.dirty = false
	return output, nil
}

// Content implements ContentProvider by calling Render and discarding the error.
func (p *Pane) Content() string {
	s, _ := p.Render()
	return s
}

// ChildRect returns the positioned rect for the i-th child after resolution.
func (p *Pane) ChildRect(i int) Rect {
	if !p.resolved || i < 0 || i >= len(p.sizes) {
		return Rect{}
	}
	if p.direction == Horizontal {
		return p.childRectHorizontal(i)
	}
	return p.childRectVertical(i)
}

// --- element interface (allows Pane to nest inside other Panes) ---

func (p *Pane) content() string {
	return p.Content()
}

func (p *Pane) setSize(width, height int) {
	p.SetSize(width, height)
}

func (p *Pane) style() (lipgloss.Style, bool) {
	return lipgloss.Style{}, false
}

func (p *Pane) sizeHint(availW, availH int) (SizeHint, bool) {
	// Delegate to single child element if it provides a size hint.
	// This allows Fit() panes wrapping a SizeHinter widget to work.
	if len(p.children) == 1 {
		return p.children[0].elem.sizeHint(availW, availH)
	}
	return SizeHint{}, false
}

// Compile-time check: Pane satisfies element.
var _ element = (*Pane)(nil)

// --- Internal helpers ---

// toConstraint converts this Pane's Dimension and modifiers into
// an internal constraint for use in the resolution algorithm.
func (p *Pane) toConstraint() constraint {
	cs := constraint{
		minSize:    p.minSize,
		maxSize:    p.maxSize,
		optional:   p.optional,
		minSizeFit: p.minSizeFit,
	}
	switch p.dim.kind {
	case dimensionPercent:
		cs.kind = constraintFlex
		cs.flexWeight = p.dim.value
	case dimensionFixed:
		cs.kind = constraintFixed
		cs.fixedSize = int(p.dim.value)
	case dimensionFlex:
		cs.kind = constraintFlex
		cs.flexWeight = p.dim.value
	case dimensionFit:
		cs.kind = constraintFit
	}
	return cs
}

// axisSize returns the available space along the pane's primary axis.
func (p *Pane) axisSize() int {
	if p.direction == Horizontal {
		return p.width
	}
	return p.height
}

// childDimensions returns (width, height) for a child given its resolved
// size along the pane's axis.
func (p *Pane) childDimensions(axisSize int) (int, int) {
	if p.direction == Horizontal {
		return axisSize, p.height
	}
	return p.width, axisSize
}

// childRectHorizontal computes the rect for the i-th child in a horizontal pane.
func (p *Pane) childRectHorizontal(i int) Rect {
	x := 0
	for j := range i {
		if p.sizes[j] > 0 {
			x += p.sizes[j]
			if j < i-1 || (j == i-1 && p.sizes[i] > 0) {
				x += p.gap
			}
		}
	}
	return Rect{
		X:      x,
		Y:      0,
		Width:  p.sizes[i],
		Height: p.height,
	}
}

// childRectVertical computes the rect for the i-th child in a vertical pane.
func (p *Pane) childRectVertical(i int) Rect {
	y := 0
	for j := range i {
		if p.sizes[j] > 0 {
			y += p.sizes[j]
			if j < i-1 || (j == i-1 && p.sizes[i] > 0) {
				y += p.gap
			}
		}
	}
	return Rect{
		X:      0,
		Y:      y,
		Width:  p.width,
		Height: p.sizes[i],
	}
}

// sizeHinterFunc adapts a function to the SizeHinter interface.
type sizeHinterFunc func(availWidth, availHeight int) SizeHint

func (f sizeHinterFunc) SizeHint(availWidth, availHeight int) SizeHint {
	return f(availWidth, availHeight)
}
