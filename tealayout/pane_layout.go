package tealayout

// PaneLayout wraps a Layout and FocusManager into a single convenience type.
// It is not a tea.Model — just a layout helper that centralizes size management,
// rendering, and focus navigation.
type PaneLayout struct {
	layout *Layout
	focus  *FocusManager
}

// NewPaneLayout creates a PaneLayout wrapping the given root pane.
func NewPaneLayout(root *Pane, opts ...Option) *PaneLayout {
	layout := NewLayout(root, opts...)
	focus := NewFocusManager(layout)
	return &PaneLayout{
		layout: layout,
		focus:  focus,
	}
}

// SetSize updates the layout dimensions.
func (pl *PaneLayout) SetSize(w, h int) {
	pl.layout.SetSize(w, h)
}

// Render resolves the layout and returns the composed output.
func (pl *PaneLayout) Render() (string, error) {
	return pl.layout.Render()
}

// MarkDirty forces re-resolution and re-rendering on the next call.
func (pl *PaneLayout) MarkDirty() {
	pl.layout.MarkDirty()
}

// --- Focus delegation ---

// FocusNext advances focus to the next visible pane.
func (pl *PaneLayout) FocusNext() {
	pl.focus.FocusNext()
}

// FocusPrev moves focus to the previous visible pane.
func (pl *PaneLayout) FocusPrev() {
	pl.focus.FocusPrev()
}

// FocusPane focuses a pane by name.
func (pl *PaneLayout) FocusPane(name string) error {
	return pl.focus.FocusPane(name)
}

// FocusedPane returns the currently focused pane, or nil.
func (pl *PaneLayout) FocusedPane() *Pane {
	return pl.focus.FocusedPane()
}

// Focused returns true if the named pane is currently focused.
func (pl *PaneLayout) Focused(name string) bool {
	return pl.focus.Focused(name)
}

// EnsureFocusedVisible checks if the focused pane is still visible
// and advances to the next visible pane if not.
func (pl *PaneLayout) EnsureFocusedVisible() {
	pl.focus.EnsureFocusedVisible()
}

// --- Pane access ---

// Pane returns the named pane, or nil if not found.
func (pl *PaneLayout) Pane(name string) *Pane {
	return pl.layout.Pane(name)
}

// ShowPane makes the named pane visible.
func (pl *PaneLayout) ShowPane(name string) {
	if p := pl.layout.Pane(name); p != nil {
		p.Show()
	}
}

// HidePane makes the named pane hidden.
func (pl *PaneLayout) HidePane(name string) {
	if p := pl.layout.Pane(name); p != nil {
		p.Hide()
	}
}

// SetPaneVisible sets visibility of the named pane.
func (pl *PaneLayout) SetPaneVisible(name string, visible bool) {
	if p := pl.layout.Pane(name); p != nil {
		p.SetVisible(visible)
	}
}

// ResizeFocused adjusts the focused pane's flex weight by delta, trading
// weight with visible flex siblings. The total weight across all visible
// flex panes is preserved: weight gained by the focused pane is taken
// proportionally from siblings, and vice versa. Each pane is clamped to
// its MinFlexWeight floor.
//
// After applying weight changes, the parent is re-resolved to check whether
// the focused pane's cell allocation actually changed. If it didn't (e.g.
// the pane is already at its cell-based minimum from minSize/minSizeFit),
// all weight changes are rolled back so the percentage display stays in
// sync with the visual layout.
//
// No-op if focused pane is nil or not flex/percent.
func (pl *PaneLayout) ResizeFocused(delta float64) {
	fp := pl.focus.FocusedPane()
	if fp == nil {
		return
	}
	pl.resizePane(fp, delta)
}

// ResizePane adjusts the named pane's flex weight by delta, trading weight
// with visible flex siblings. Behaves identically to ResizeFocused but
// operates on a specific pane rather than the focused one.
func (pl *PaneLayout) ResizePane(name string, delta float64) {
	p := pl.layout.Pane(name)
	if p == nil {
		return
	}
	pl.resizePane(p, delta)
}

// ResizePaneByCells adjusts the named pane by the given number of cells
// (columns or rows depending on the parent's axis). The weight delta is
// computed from the parent's total flex weight and available space so that
// each unit corresponds to approximately one cell.
func (pl *PaneLayout) ResizePaneByCells(name string, cells int) {
	p := pl.layout.Pane(name)
	if p == nil || !p.dim.IsFlex() || p.parent == nil {
		return
	}
	delta := pl.cellsToDelta(p, cells)
	if delta == 0 {
		return
	}
	pl.resizePane(p, delta)
}

// ResizeFocusedByCells adjusts the focused pane by the given number of
// cells. See ResizePaneByCells for details.
func (pl *PaneLayout) ResizeFocusedByCells(cells int) {
	fp := pl.focus.FocusedPane()
	if fp == nil || !fp.dim.IsFlex() || fp.parent == nil {
		return
	}
	delta := pl.cellsToDelta(fp, cells)
	if delta == 0 {
		return
	}
	pl.resizePane(fp, delta)
}

// cellsToDelta converts a cell count to a flex weight delta based on the
// pane's parent axis size and total visible flex weight among siblings.
func (pl *PaneLayout) cellsToDelta(p *Pane, cells int) float64 {
	parent := p.parent
	axisSize := parent.axisSize()
	if axisSize == 0 {
		return 0
	}

	// Sum visible flex weights in the parent (including p itself).
	var totalWeight float64
	for _, ch := range parent.children {
		cp, ok := ch.elem.(*Pane)
		if !ok || !cp.visible || !cp.dim.IsFlex() {
			continue
		}
		totalWeight += cp.dim.value
	}
	if totalWeight == 0 {
		return 0
	}

	return float64(cells) * totalWeight / float64(axisSize)
}

// resizePane is the shared implementation for ResizeFocused and ResizePane.
func (pl *PaneLayout) resizePane(fp *Pane, delta float64) {
	if !fp.dim.IsFlex() {
		return
	}

	parent := fp.parent
	if parent == nil {
		return
	}

	// Find the pane's index in the parent's children.
	fpIdx := -1
	for i, ch := range parent.children {
		if ch.elem == fp {
			fpIdx = i
			break
		}
	}
	if fpIdx < 0 {
		return
	}

	// Collect visible flex siblings (excluding the target pane).
	var siblings []*Pane
	var siblingWeight float64
	for _, ch := range parent.children {
		cp, ok := ch.elem.(*Pane)
		if !ok || cp == fp || !cp.visible || !cp.dim.IsFlex() {
			continue
		}
		siblings = append(siblings, cp)
		siblingWeight += cp.dim.value
	}
	if len(siblings) == 0 {
		return
	}

	// Compute effective minimum weight for each pane, considering both
	// the static minFlexWeight and content-derived cell minimums.
	totalWeight := fp.dim.value + siblingWeight
	axisSize := parent.axisSize()
	horizontal := parent.direction == Horizontal

	effectiveMin := func(p *Pane) float64 {
		floor := p.minFlexWeight
		if axisSize <= 0 || totalWeight <= 0 {
			return floor
		}
		minCells := p.minSize
		if p.minSizeFit {
			if hint, ok := p.sizeHint(parent.width, parent.height); ok {
				hintMin := hint.Min.Height
				if horizontal {
					hintMin = hint.Min.Width
				}
				if hintMin > minCells {
					minCells = hintMin
				}
			}
		}
		if minCells > 0 {
			cellWeight := float64(minCells) * totalWeight / float64(axisSize)
			if cellWeight > floor {
				floor = cellWeight
			}
		}
		return floor
	}

	fpMin := effectiveMin(fp)

	// Capture old resolved size for rollback check. If the parent has never
	// been resolved (sizes is nil), skip the rollback check — let weights
	// through since there's nothing to compare against.
	hasOldSize := len(parent.sizes) > fpIdx
	oldSize := 0
	if hasOldSize {
		oldSize = parent.sizes[fpIdx]
	}

	// Snapshot old dimensions for rollback.
	oldFPDim := fp.dim
	oldSibDims := make([]Dimension, len(siblings))
	for i, sib := range siblings {
		oldSibDims[i] = sib.dim
	}

	// Clamp the pane's new weight.
	newWeight := fp.dim.value + delta
	if newWeight < fpMin {
		newWeight = fpMin
	}
	// Cap: siblings can't go below their combined effective minimums.
	maxWeight := fp.dim.value + siblingWeight
	for _, sib := range siblings {
		maxWeight -= effectiveMin(sib)
	}
	if newWeight > maxWeight {
		newWeight = maxWeight
	}

	// Compute the actual delta after clamping.
	actualDelta := newWeight - fp.dim.value
	if actualDelta == 0 {
		return
	}

	// Redistribute -actualDelta proportionally across siblings.
	// Multi-pass: when a sibling hits its floor, freeze it and redistribute
	// the remainder to unfrozen siblings.
	sibMins := make([]float64, len(siblings))
	for i, sib := range siblings {
		sibMins[i] = effectiveMin(sib)
	}

	toDistribute := -actualDelta
	frozen := make([]bool, len(siblings))
	newWeights := make([]float64, len(siblings))
	for i, sib := range siblings {
		newWeights[i] = sib.dim.value
	}

	for range len(siblings) {
		// Sum unfrozen sibling weights.
		var poolWeight float64
		for i := range siblings {
			if !frozen[i] {
				poolWeight += newWeights[i]
			}
		}
		if poolWeight == 0 {
			break
		}

		anyClamped := false
		var leftover float64
		for i := range siblings {
			if frozen[i] {
				continue
			}
			share := toDistribute * (newWeights[i] / poolWeight)
			proposed := newWeights[i] + share
			if proposed < sibMins[i] {
				leftover += proposed - sibMins[i]
				proposed = sibMins[i]
				frozen[i] = true
				anyClamped = true
			}
			newWeights[i] = proposed
		}
		if !anyClamped {
			break
		}
		toDistribute = leftover // redistribute the excess
	}

	// Apply the new weights.
	for i, sib := range siblings {
		sib.SetDimension(Dimension{kind: sib.dim.kind, value: newWeights[i]})
	}
	fp.SetDimension(Dimension{kind: fp.dim.kind, value: newWeight})

	// Resolve-and-check: if the pane's cell size didn't change,
	// roll back all weight changes (the resize had no visual effect).
	if hasOldSize {
		parent.resolved = false
		_, err := parent.Resolve()
		if err != nil || (len(parent.sizes) > fpIdx && parent.sizes[fpIdx] == oldSize) {
			// Resolve failed or no visual change — roll back.
			fp.SetDimension(oldFPDim)
			for i, sib := range siblings {
				sib.SetDimension(oldSibDims[i])
			}
		}
	}
}

// Layout returns the underlying Layout for advanced use.
func (pl *PaneLayout) Layout() *Layout {
	return pl.layout
}
