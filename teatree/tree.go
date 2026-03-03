package teatree

// Tree represents a collection of root nodes with focus management
// Simplified from treeview's Tree[T] - single focus, visible node tracking
type Tree[T any] struct {
	nodes       []*Node[T]
	focusedNode *Node[T]
	provider    NodeProvider[T]
}

type TreeArgs[T any] struct {
	ExpanderControls *ExpanderControls
	NodeProvider     NodeProvider[T]
	FocusedNode      *Node[T]
}

// NewTree creates a new tree with the given root nodes
func NewTree[T any](nodes []*Node[T], args *TreeArgs[T]) *Tree[T] {
	if args == nil {
		args = &TreeArgs[T]{}
	}
	if args.ExpanderControls == nil {
		args.ExpanderControls = &TriangleExpanderControls
	}
	if args.NodeProvider == nil {
		args.NodeProvider = NewCompactNodeProvider[T](*args.ExpanderControls)
	}
	t := &Tree[T]{
		nodes:       nodes,
		focusedNode: args.FocusedNode,
		provider:    args.NodeProvider,
	}

	// Set initial focus to first visible node if not already set
	if t.focusedNode == nil && len(nodes) > 0 {
		visible := t.VisibleNodes()
		if len(visible) > 0 {
			t.focusedNode = visible[0]
		}
	}
	return t
}

// TreeOption is a functional option for configuring a Tree
type TreeOption[T any] func(*Tree[T])

// WithProvider sets the node provider
func WithProvider[T any](provider NodeProvider[T]) TreeOption[T] {
	return func(t *Tree[T]) {
		t.provider = provider
	}
}

// WithInitialFocus sets the initial focused node by ID
func WithInitialFocus[T any](nodeID string) TreeOption[T] {
	return func(t *Tree[T]) {
		for _, root := range t.nodes {
			node := root.FindByID(nodeID)
			if node != nil {
				t.focusedNode = node
				return
			}
		}
	}
}

// Nodes returns the root nodes
func (t *Tree[T]) Nodes() []*Node[T] {
	return t.nodes
}

// SetNodes replaces the root nodes
func (t *Tree[T]) SetNodes(nodes []*Node[T]) {
	t.nodes = nodes

	// If focused node is no longer in tree, reset focus
	if t.focusedNode != nil {
		found := false
		for _, root := range nodes {
			if root.FindByID(t.focusedNode.ID()) != nil {
				found = true
				break
			}
		}
		if !found {
			t.focusedNode = nil
		}
	}

	// Set focus to first visible node if no focus
	if t.focusedNode == nil && len(nodes) > 0 {
		visible := t.VisibleNodes()
		if len(visible) > 0 {
			t.focusedNode = visible[0]
		}
	}
}

// Provider returns the node provider
func (t *Tree[T]) Provider() NodeProvider[T] {
	return t.provider
}

// SetProvider sets the node provider
func (t *Tree[T]) SetProvider(provider NodeProvider[T]) {
	t.provider = provider
}

// FocusedNode returns the currently focused node (nil if none)
func (t *Tree[T]) FocusedNode() (node *Node[T]) {
	if t == nil {
		return nil
	}
	node = t.focusedNode
	if node == nil {
		node = t.FirstNode()
	}
	if node != nil {
		t.focusedNode = node
	}
	return node
}

// FirstNode returns the first node of the tree
func (t *Tree[T]) FirstNode() *Node[T] {
	if t == nil {
		return nil
	}
	if len(t.nodes) == 0 {
		return nil
	}
	return t.nodes[0]
}

// IsFocusedNode returns the currently focused node (nil if none)
func (t *Tree[T]) IsFocusedNode(node *Node[T]) bool {
	return t.focusedNode == node
}

// SetFocusedNode sets the focused node by ID
func (t *Tree[T]) SetFocusedNode(nodeID string) bool {
	for _, root := range t.nodes {
		node := root.FindByID(nodeID)
		if node != nil {
			t.focusedNode = node
			return true
		}
	}
	return false
}

// FindByID finds a node by ID in the tree
func (t *Tree[T]) FindByID(nodeID string) *Node[T] {
	for _, root := range t.nodes {
		node := root.FindByID(nodeID)
		if node != nil {
			return node
		}
	}
	return nil
}

// VisibleNodes returns all currently visible nodes in tree order
func (t *Tree[T]) VisibleNodes() []*Node[T] {
	var visible []*Node[T]
	for _, root := range t.nodes {
		t.collectVisibleNodes(root, &visible)
	}
	return visible
}

// collectVisibleNodes recursively collects visible nodes
func (t *Tree[T]) collectVisibleNodes(node *Node[T], result *[]*Node[T]) {
	if !node.IsVisible() {
		return
	}

	*result = append(*result, node)

	// Only traverse children if node is expanded
	if node.IsExpanded() {
		for _, child := range node.Children() {
			t.collectVisibleNodes(child, result)
		}
	}
}

// MoveUp moves focus to the previous visible node
func (t *Tree[T]) MoveUp() bool {
	if t.focusedNode == nil {
		return false
	}

	visible := t.VisibleNodes()
	for i, node := range visible {
		if node == t.focusedNode && i > 0 {
			t.focusedNode = visible[i-1]
			return true
		}
	}

	return false
}

// MoveDown moves focus to the next visible node
func (t *Tree[T]) MoveDown() bool {
	if t.focusedNode == nil {
		return false
	}

	visible := t.VisibleNodes()
	for i, node := range visible {
		if node == t.focusedNode && i < len(visible)-1 {
			t.focusedNode = visible[i+1]
			return true
		}
	}

	return false
}

// ExpandFocused expands the currently focused node
func (t *Tree[T]) ExpandFocused() bool {
	if !t.focusedNode.HasChildren() {
		return false
	}

	if !t.focusedNode.IsExpanded() {
		t.focusedNode.Expand()
		return true
	}

	return false
}

// CollapseFocused collapses the currently focused node
func (t *Tree[T]) CollapseFocused() bool {
	if !t.focusedNode.HasChildren() {
		return false
	}

	if t.focusedNode.IsExpanded() {
		t.focusedNode.Collapse()
		return true
	}

	return false
}

// ToggleFocused toggles the expansion state of the focused node
func (t *Tree[T]) ToggleFocused() bool {
	if !t.focusedNode.HasChildren() {
		return false
	}

	t.focusedNode.Toggle()
	return true
}

// ExpandAll expands all nodes in the tree
func (t *Tree[T]) ExpandAll() {
	for _, root := range t.nodes {
		t.expandAll(root)
	}
}

// expandAll recursively expands all nodes
func (t *Tree[T]) expandAll(node *Node[T]) {
	if node.HasChildren() {
		node.Expand()
		for _, child := range node.Children() {
			t.expandAll(child)
		}
	}
}

// CollapseAll collapses all nodes in the tree
func (t *Tree[T]) CollapseAll() {
	for _, root := range t.nodes {
		t.collapseAll(root)
	}
}

// collapseAll recursively collapses all nodes
func (t *Tree[T]) collapseAll(node *Node[T]) {
	if node.HasChildren() {
		node.Collapse()
		for _, child := range node.Children() {
			t.collapseAll(child)
		}
	}
}
