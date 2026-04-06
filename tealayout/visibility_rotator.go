package tealayout

// VisibilityRotator cycles a MultiPaneLayout through predefined visibility
// combinations. Each combo is a slice of pane names that should be visible;
// all other content panes are hidden.
type VisibilityRotator struct {
	mpl    *MultiPaneLayout
	combos [][]string
	index  int
}

// NewVisibilityRotator creates a rotator for the given layout and combos.
// Panics if combos is empty (programming error).
func NewVisibilityRotator(mpl *MultiPaneLayout, combos [][]string) *VisibilityRotator {
	if len(combos) == 0 {
		panic("NewVisibilityRotator: combos must not be empty")
	}
	return &VisibilityRotator{
		mpl:    mpl,
		combos: combos,
	}
}

// Next advances to the next combo (wrapping) and applies it.
func (r *VisibilityRotator) Next() {
	r.index = (r.index + 1) % len(r.combos)
	r.apply()
}

// Prev moves to the previous combo (wrapping) and applies it.
func (r *VisibilityRotator) Prev() {
	r.index = (r.index - 1 + len(r.combos)) % len(r.combos)
	r.apply()
}

// Apply applies the current combo without changing the index.
func (r *VisibilityRotator) Apply() {
	r.apply()
}

// SetIndex jumps to combo i and applies it. Panics if i is out of range.
func (r *VisibilityRotator) SetIndex(i int) {
	if i < 0 || i >= len(r.combos) {
		panic("VisibilityRotator.SetIndex: index out of range")
	}
	r.index = i
	r.apply()
}

// Index returns the current combo index.
func (r *VisibilityRotator) Index() int {
	return r.index
}

// Current returns a copy of the current combo.
func (r *VisibilityRotator) Current() []string {
	out := make([]string, len(r.combos[r.index]))
	copy(out, r.combos[r.index])
	return out
}

// Len returns the number of combos.
func (r *VisibilityRotator) Len() int {
	return len(r.combos)
}

func (r *VisibilityRotator) apply() {
	visSet := make(map[string]bool, len(r.combos[r.index]))
	for _, name := range r.combos[r.index] {
		visSet[name] = true
	}
	for _, name := range r.mpl.paneNames {
		r.mpl.pl.SetPaneVisible(name, visSet[name])
	}
	r.mpl.pl.EnsureFocusedVisible()
}
