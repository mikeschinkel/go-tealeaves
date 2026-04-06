package tealayout

// PaneDef describes a content pane for MultiPaneLayout.
type PaneDef struct {
	Name          string    // required: lookup key, must be unique
	Element       any       // required when Children is empty: *Element[T] or *Pane
	Children      []PaneDef // vertical stack of sub-panes; Element is ignored when set
	Dim           Dimension // sizing: Percent(25), Flex(1.618), Fixed(30), Fit(); zero defaults to Flex(1)
	MinFlexWeight float64   // floor for ResizeFocused (0 = no floor)
	MinSize       int       // minimum cells (0 = none)
	MaxSize       int       // maximum cells (0 = unbounded)
	MinSizeFit    bool      // use SizeHinter.Min as effective minimum
	Optional      bool      // auto-remove when too small
}

// MultiPaneLayoutOption configures a MultiPaneLayout.
type MultiPaneLayoutOption func(*multiPaneConfig)

type multiPaneConfig struct {
	header    element
	footer    element
	headerDim Dimension
	footerDim Dimension
	contentGap int
}

// WithHeader adds a header element (fixed 1-line by default).
func WithHeader(elem any) MultiPaneLayoutOption {
	return func(c *multiPaneConfig) {
		c.header = toElement(elem)
	}
}

// WithFooter adds a footer element (fixed 1-line by default).
func WithFooter(elem any) MultiPaneLayoutOption {
	return func(c *multiPaneConfig) {
		c.footer = toElement(elem)
	}
}

// WithHeaderDim overrides the header dimension (default Fixed(1)).
func WithHeaderDim(d Dimension) MultiPaneLayoutOption {
	return func(c *multiPaneConfig) {
		c.headerDim = d
	}
}

// WithFooterDim overrides the footer dimension (default Fixed(1)).
func WithFooterDim(d Dimension) MultiPaneLayoutOption {
	return func(c *multiPaneConfig) {
		c.footerDim = d
	}
}

// WithContentGap sets the gap between content columns.
func WithContentGap(n int) MultiPaneLayoutOption {
	return func(c *multiPaneConfig) {
		c.contentGap = n
	}
}

// MultiPaneLayout provides a header + N content columns + footer layout
// with built-in resize, focus, and visibility management.
type MultiPaneLayout struct {
	pl        *PaneLayout
	paneNames []string // ordered content pane names
}

// NewMultiPaneLayout builds a layout from PaneDefs with optional header/footer.
func NewMultiPaneLayout(panes []PaneDef, opts ...MultiPaneLayoutOption) *MultiPaneLayout {
	cfg := multiPaneConfig{
		headerDim: Fixed(1),
		footerDim: Fixed(1),
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	// Build content columns from PaneDefs.
	columns := make([]element, 0, len(panes))
	names := make([]string, 0, len(panes))
	for _, pd := range panes {
		dim := defaultDim(pd.Dim)

		var col *Pane
		if len(pd.Children) > 0 {
			// Nested group: build child Rows inside a Column.
			// Children within children is not supported.
			childRows := make([]element, 0, len(pd.Children))
			for _, child := range pd.Children {
				if len(child.Children) > 0 {
					panic("MultiPaneLayout: nested Children within Children is not supported")
				}
				row := NewRow(defaultDim(child.Dim), toElement(child.Element)).
					WithName(child.Name).
					WithFocusable().
					WithMinFlexWeight(child.MinFlexWeight)
				if child.MinSize > 0 {
					row.WithMinSize(child.MinSize)
				}
				if child.MaxSize > 0 {
					row.WithMaxSize(child.MaxSize)
				}
				if child.MinSizeFit {
					row.WithMinSizeFit()
				}
				if child.Optional {
					row.WithOptional(true)
				}
				childRows = append(childRows, row)
			}
			// Group Column: NOT focusable (children handle focus).
			col = NewColumn(dim, childRows...).
				WithName(pd.Name).
				WithMinFlexWeight(pd.MinFlexWeight)
		} else {
			col = NewColumn(dim, toElement(pd.Element)).
				WithName(pd.Name).
				WithFocusable().
				WithMinFlexWeight(pd.MinFlexWeight)
		}

		if pd.MinSize > 0 {
			col.WithMinSize(pd.MinSize)
		}
		if pd.MaxSize > 0 {
			col.WithMaxSize(pd.MaxSize)
		}
		if pd.MinSizeFit {
			col.WithMinSizeFit()
		}
		if pd.Optional {
			col.WithOptional(true)
		}
		names = append(names, pd.Name)
		columns = append(columns, col)
	}

	contentRow := NewRow(Flex(1), columns...)
	if cfg.contentGap > 0 {
		contentRow.WithGap(cfg.contentGap)
	}
	contentRow.WithName("content")

	// Build the root column: [header], contentRow, [footer].
	rootChildren := make([]element, 0, 3)
	if cfg.header != nil {
		headerPane := NewRow(cfg.headerDim, cfg.header).WithName("header")
		rootChildren = append(rootChildren, headerPane)
	}
	rootChildren = append(rootChildren, contentRow)
	if cfg.footer != nil {
		footerPane := NewRow(cfg.footerDim, cfg.footer).WithName("footer")
		rootChildren = append(rootChildren, footerPane)
	}

	root := NewColumn(Percent100, rootChildren...)
	pl := NewPaneLayout(root)

	return &MultiPaneLayout{
		pl:        pl,
		paneNames: names,
	}
}

// --- Delegated methods ---

// SetSize updates the layout dimensions.
func (m *MultiPaneLayout) SetSize(w, h int) { m.pl.SetSize(w, h) }

// Render resolves the layout and returns the composed output.
func (m *MultiPaneLayout) Render() (string, error) { return m.pl.Render() }

// MarkDirty forces re-resolution and re-rendering.
func (m *MultiPaneLayout) MarkDirty() { m.pl.MarkDirty() }

// FocusNext advances focus to the next visible pane.
func (m *MultiPaneLayout) FocusNext() { m.pl.FocusNext() }

// FocusPrev moves focus to the previous visible pane.
func (m *MultiPaneLayout) FocusPrev() { m.pl.FocusPrev() }

// FocusPane focuses a pane by name.
func (m *MultiPaneLayout) FocusPane(name string) error { return m.pl.FocusPane(name) }

// FocusedPane returns the currently focused pane, or nil.
func (m *MultiPaneLayout) FocusedPane() *Pane { return m.pl.FocusedPane() }

// Focused returns true if the named pane is currently focused.
func (m *MultiPaneLayout) Focused(name string) bool { return m.pl.Focused(name) }

// EnsureFocusedVisible checks if the focused pane is still visible
// and advances to the next visible pane if not.
func (m *MultiPaneLayout) EnsureFocusedVisible() { m.pl.EnsureFocusedVisible() }

// Pane returns the named pane, or nil if not found.
func (m *MultiPaneLayout) Pane(name string) *Pane { return m.pl.Pane(name) }

// ShowPane makes the named pane visible.
func (m *MultiPaneLayout) ShowPane(name string) { m.pl.ShowPane(name) }

// HidePane makes the named pane hidden.
func (m *MultiPaneLayout) HidePane(name string) { m.pl.HidePane(name) }

// SetPaneVisible sets visibility of the named pane.
func (m *MultiPaneLayout) SetPaneVisible(name string, visible bool) {
	m.pl.SetPaneVisible(name, visible)
}

// ResizeFocused adjusts the focused pane's flex weight by delta.
func (m *MultiPaneLayout) ResizeFocused(delta float64) { m.pl.ResizeFocused(delta) }

// ResizeFocusedColumn adjusts the top-level content column containing the
// focused pane. If the focused pane is a child inside a group, the group
// column is resized among its top-level siblings. If the focused pane is
// itself a top-level column, it behaves identically to ResizeFocused.
func (m *MultiPaneLayout) ResizeFocusedColumn(delta float64) {
	fp := m.pl.FocusedPane()
	if fp == nil {
		return
	}
	// Walk up from the focused pane to find the top-level content column.
	// A top-level column's parent has name "content" (the content row).
	col := fp
	for col.parent != nil && col.parent.name != "content" {
		col = col.parent
	}
	if col.name == "" {
		return
	}
	m.pl.ResizePane(col.name, delta)
}

// ResizeFocusedColumnByCells adjusts the top-level content column containing
// the focused pane by the given number of cells. See ResizeFocusedColumn for
// the column-finding logic and ResizePaneByCells for the cell calculation.
func (m *MultiPaneLayout) ResizeFocusedColumnByCells(cells int) {
	fp := m.pl.FocusedPane()
	if fp == nil {
		return
	}
	col := fp
	for col.parent != nil && col.parent.name != "content" {
		col = col.parent
	}
	if col.name == "" {
		return
	}
	m.pl.ResizePaneByCells(col.name, cells)
}

// PaneLayout returns the underlying PaneLayout for advanced use.
func (m *MultiPaneLayout) PaneLayout() *PaneLayout { return m.pl }

// --- Convenience methods ---

// PaneNames returns the ordered content pane names.
func (m *MultiPaneLayout) PaneNames() []string {
	out := make([]string, len(m.paneNames))
	copy(out, m.paneNames)
	return out
}

// VisiblePaneNames returns the currently visible content pane names.
func (m *MultiPaneLayout) VisiblePaneNames() []string {
	var out []string
	for _, name := range m.paneNames {
		if p := m.pl.Pane(name); p != nil && p.visible {
			out = append(out, name)
		}
	}
	return out
}

// visibleFlexTotal returns the sum of flex weights across visible content panes.
func (m *MultiPaneLayout) visibleFlexTotal() float64 {
	total := 0.0
	for _, name := range m.paneNames {
		if p := m.pl.Pane(name); p != nil && p.visible && p.dim.IsFlex() {
			total += p.dim.value
		}
	}
	return total
}

// VisibleFlexPercent returns the named pane's flex weight as a percentage
// of the total visible flex weight (0–100). Returns 0 if the pane is not
// found, hidden, or not a flex/percent dimension.
func (m *MultiPaneLayout) VisibleFlexPercent(name string) float64 {
	p := m.pl.Pane(name)
	if p == nil || !p.visible || !p.dim.IsFlex() {
		return 0
	}
	total := m.visibleFlexTotal()
	if total == 0 {
		return 0
	}
	return p.dim.value / total * 100
}

// VisibleFlexPercents returns every visible content pane's flex weight as a
// percentage of the total visible flex weight (0–100). Hidden and non-flex
// panes are omitted from the map.
func (m *MultiPaneLayout) VisibleFlexPercents() map[string]float64 {
	total := m.visibleFlexTotal()
	if total == 0 {
		return nil
	}
	out := make(map[string]float64)
	for _, name := range m.paneNames {
		if p := m.pl.Pane(name); p != nil && p.visible && p.dim.IsFlex() {
			out[name] = p.dim.value / total * 100
		}
	}
	return out
}

// TogglePane flips visibility of the named content pane.
// No-op if toggling would hide the last visible content pane.
// Calls EnsureFocusedVisible after the change.
func (m *MultiPaneLayout) TogglePane(name string) {
	p := m.pl.Pane(name)
	if p == nil {
		return
	}
	if p.visible && len(m.VisiblePaneNames()) <= 1 {
		return
	}
	m.pl.SetPaneVisible(name, !p.visible)
	m.pl.EnsureFocusedVisible()
}

// SoloPane shows only the named content pane, hiding all others.
// Focuses the solo pane. No-op if name not found.
func (m *MultiPaneLayout) SoloPane(name string) {
	if m.pl.Pane(name) == nil {
		return
	}
	for _, n := range m.paneNames {
		m.pl.SetPaneVisible(n, n == name)
	}
	m.pl.FocusPane(name) //nolint:errcheck
}

// ShowAllPanes makes all content panes visible.
// Calls EnsureFocusedVisible after the change.
func (m *MultiPaneLayout) ShowAllPanes() {
	for _, n := range m.paneNames {
		m.pl.ShowPane(n)
	}
	m.pl.EnsureFocusedVisible()
}

// defaultDim returns the given Dimension, or Flex(1) if it is the zero value.
func defaultDim(d Dimension) Dimension {
	if d == (Dimension{}) {
		return Flex(1)
	}
	return d
}

// toElement converts an any value to an element. It handles *Pane and
// *Element[T] (which both satisfy element) and panics on invalid types.
func toElement(v any) element {
	if v == nil {
		return nil
	}
	if e, ok := v.(element); ok {
		return e
	}
	panic("MultiPaneLayout: Element must be *Element[T] or *Pane")
}
