package tealayout

// FocusManager tracks focus state across named panes in a layout.
// It navigates in depth-first order, skipping hidden panes.
type FocusManager struct {
	layout  *Layout
	focused *Pane
	order   []*Pane // depth-first leaf/named panes
}

// NewFocusManager creates a FocusManager for the given layout.
// The focus order is built by walking the tree depth-first, collecting
// leaf panes (no children) or named panes.
func NewFocusManager(layout *Layout) *FocusManager {
	fm := &FocusManager{layout: layout}
	fm.buildOrder(layout.Root())
	if len(fm.order) > 0 {
		fm.focusPane(fm.order[0])
	}
	return fm
}

// buildOrder walks the tree depth-first, collecting focusable panes.
func (fm *FocusManager) buildOrder(p *Pane) {
	isLeaf := true
	for _, ch := range p.children {
		if cp, ok := ch.elem.(*Pane); ok {
			isLeaf = false
			fm.buildOrder(cp)
		}
	}
	// A pane is focusable if it's a leaf (has no child panes) or is named.
	if isLeaf || p.name != "" {
		fm.order = append(fm.order, p)
	}
}

// FocusNext advances focus to the next visible pane.
func (fm *FocusManager) FocusNext() {
	fm.advance(1)
}

// FocusPrev moves focus to the previous visible pane.
func (fm *FocusManager) FocusPrev() {
	fm.advance(-1)
}

// FocusPane focuses a pane by name. Returns ErrPaneNotFound if not found.
func (fm *FocusManager) FocusPane(name string) error {
	p := fm.layout.Pane(name)
	if p == nil {
		return NewErr(ErrPaneNotFound, "name", name)
	}
	fm.focusPane(p)
	return nil
}

// FocusedPane returns the currently focused pane, or nil.
func (fm *FocusManager) FocusedPane() *Pane {
	return fm.focused
}

// Focused returns true if the named pane is currently focused.
func (fm *FocusManager) Focused(name string) bool {
	return fm.focused != nil && fm.focused.name == name
}

// advance moves focus by delta (+1 or -1), skipping hidden panes.
func (fm *FocusManager) advance(delta int) {
	n := len(fm.order)
	if n == 0 {
		return
	}

	// Find current index
	cur := 0
	for i, p := range fm.order {
		if p == fm.focused {
			cur = i
			break
		}
	}

	// Walk in direction, skipping hidden
	for step := range n {
		_ = step
		cur = (cur + delta + n) % n
		if fm.order[cur].visible {
			fm.focusPane(fm.order[cur])
			return
		}
	}
}

// focusPane blurs the old pane and focuses the new one.
func (fm *FocusManager) focusPane(p *Pane) {
	if fm.focused == p {
		return
	}
	if fm.focused != nil {
		fm.focused.focused = false
		fm.blurWidget(fm.focused)
		fm.focused.MarkDirty()
	}
	fm.focused = p
	p.focused = true
	fm.focusWidget(p)
	p.MarkDirty()
}

// focusWidget calls Focus() on the pane's child elements if they implement Focusable.
func (fm *FocusManager) focusWidget(p *Pane) {
	for _, ch := range p.children {
		if f, ok := ch.elem.(focusableElement); ok {
			f.focus()
		}
	}
}

// blurWidget calls Blur() on the pane's child elements if they implement Focusable.
func (fm *FocusManager) blurWidget(p *Pane) {
	for _, ch := range p.children {
		if f, ok := ch.elem.(focusableElement); ok {
			f.blur()
		}
	}
}

// focusableElement extends element with focus capability detection.
type focusableElement interface {
	element
	focus()
	blur()
}

// EnsureFocusedVisible checks if the focused pane is still visible.
// If not, advances to the next visible pane.
func (fm *FocusManager) EnsureFocusedVisible() {
	if fm.focused != nil && !fm.focused.visible {
		fm.FocusNext()
	}
}
