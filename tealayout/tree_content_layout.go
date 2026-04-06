package tealayout

// TreeContentLayout is a generic 2-pane layout: a Fit tree pane on the left
// and a Flex(1) content pane on the right. It provides typed access to both
// widgets and composes a PaneLayout internally.
//
// The tree widget should implement SizeHinter (via SizeHint()) so the Fit
// resolver can query its desired width. The content widget takes the
// remaining space.
//
// Self-framed widgets (e.g. models with their own lipgloss borders) must NOT
// implement Styler — tealayout won't add extra borders.
type TreeContentLayout[T, C any] struct {
	pl          *PaneLayout
	treeElem    *Element[T]
	contentElem *Element[C]
}

// NewTreeContentLayout creates a 2-pane layout with a Fit tree on the left
// and a Flex(1) content on the right.
func NewTreeContentLayout[T, C any](tree T, content C) *TreeContentLayout[T, C] {
	treeElem := NewElement(tree)
	contentElem := NewElement(content)

	root := NewRow(Flex(1),
		NewColumn(Fit(), treeElem).WithName("tree").WithFocusable(),
		NewColumn(Flex(1), contentElem).WithName("content").WithFocusable(),
	)

	return &TreeContentLayout[T, C]{
		pl:          NewPaneLayout(root),
		treeElem:    treeElem,
		contentElem: contentElem,
	}
}

// Tree returns the underlying tree widget.
func (tcl *TreeContentLayout[T, C]) Tree() T {
	return tcl.treeElem.Widget()
}

// Content returns the underlying content widget.
func (tcl *TreeContentLayout[T, C]) Content() C {
	return tcl.contentElem.Widget()
}

// --- Delegated to PaneLayout ---

// SetSize updates the layout dimensions.
func (tcl *TreeContentLayout[T, C]) SetSize(w, h int) {
	tcl.pl.SetSize(w, h)
}

// Render resolves the layout and returns the composed output.
func (tcl *TreeContentLayout[T, C]) Render() (string, error) {
	return tcl.pl.Render()
}

// MarkDirty forces re-resolution and re-rendering on the next call.
func (tcl *TreeContentLayout[T, C]) MarkDirty() {
	tcl.pl.MarkDirty()
}

// --- Focus convenience ---

// FocusTree focuses the tree pane.
func (tcl *TreeContentLayout[T, C]) FocusTree() {
	tcl.pl.FocusPane("tree") //nolint:errcheck
}

// FocusContent focuses the content pane.
func (tcl *TreeContentLayout[T, C]) FocusContent() {
	tcl.pl.FocusPane("content") //nolint:errcheck
}

// ToggleFocus switches focus between tree and content panes.
func (tcl *TreeContentLayout[T, C]) ToggleFocus() {
	if tcl.pl.Focused("tree") {
		tcl.FocusContent()
	} else {
		tcl.FocusTree()
	}
}

// TreeFocused returns true if the tree pane is currently focused.
func (tcl *TreeContentLayout[T, C]) TreeFocused() bool {
	return tcl.pl.Focused("tree")
}

// ContentFocused returns true if the content pane is currently focused.
func (tcl *TreeContentLayout[T, C]) ContentFocused() bool {
	return tcl.pl.Focused("content")
}

// --- Visibility ---

// ShowTree makes the tree pane visible.
func (tcl *TreeContentLayout[T, C]) ShowTree() {
	tcl.pl.ShowPane("tree")
}

// HideTree hides the tree pane.
func (tcl *TreeContentLayout[T, C]) HideTree() {
	tcl.pl.HidePane("tree")
}

// ShowContent makes the content pane visible.
func (tcl *TreeContentLayout[T, C]) ShowContent() {
	tcl.pl.ShowPane("content")
}

// HideContent hides the content pane.
func (tcl *TreeContentLayout[T, C]) HideContent() {
	tcl.pl.HidePane("content")
}

// ToggleTree flips the tree pane's visibility.
// No-op if it would hide the last visible pane.
// Calls EnsureFocusedVisible after the change.
func (tcl *TreeContentLayout[T, C]) ToggleTree() {
	tree := tcl.pl.Pane("tree")
	if tree.visible && !tcl.pl.Pane("content").visible {
		return
	}
	tcl.pl.SetPaneVisible("tree", !tree.visible)
	tcl.pl.EnsureFocusedVisible()
}

// ToggleContent flips the content pane's visibility.
// No-op if it would hide the last visible pane.
// Calls EnsureFocusedVisible after the change.
func (tcl *TreeContentLayout[T, C]) ToggleContent() {
	content := tcl.pl.Pane("content")
	if content.visible && !tcl.pl.Pane("tree").visible {
		return
	}
	tcl.pl.SetPaneVisible("content", !content.visible)
	tcl.pl.EnsureFocusedVisible()
}

// SoloTree shows only the tree pane, hides content.
// Focuses tree.
func (tcl *TreeContentLayout[T, C]) SoloTree() {
	tcl.pl.HidePane("content")
	tcl.pl.ShowPane("tree")
	tcl.FocusTree()
}

// SoloContent shows only the content pane, hides tree.
// Focuses content.
func (tcl *TreeContentLayout[T, C]) SoloContent() {
	tcl.pl.HidePane("tree")
	tcl.pl.ShowPane("content")
	tcl.FocusContent()
}

// ShowBoth makes both panes visible.
// Calls EnsureFocusedVisible after the change.
func (tcl *TreeContentLayout[T, C]) ShowBoth() {
	tcl.pl.ShowPane("tree")
	tcl.pl.ShowPane("content")
	tcl.pl.EnsureFocusedVisible()
}

// PaneLayout returns the underlying PaneLayout for advanced use.
func (tcl *TreeContentLayout[T, C]) PaneLayout() *PaneLayout {
	return tcl.pl
}
