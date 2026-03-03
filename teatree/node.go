package teatree

// Node represents a single node in the tree with generic data
// Design adopted from github.com/Digital-Shane/treeview with simplifications
type Node[T any] struct {
	id               string
	name             string
	text             string // Text for display/data access (may differ from name for tree building)
	data             T
	children         []*Node[T]
	parent           *Node[T]
	expanded         bool
	visible          bool
	hasGrandChildren *bool
}

// NewNode creates a new node with the given id, name, and data
func NewNode[T any](id, name string, data T) *Node[T] {
	return &Node[T]{
		id:       id,
		name:     name,
		data:     data,
		children: make([]*Node[T], 0),
		parent:   nil,
		expanded: false,
		visible:  true,
	}
}

// ID returns the node's unique identifier
func (n *Node[T]) ID() string {
	return n.id
}

// IsRoot returns true if no parent
func (n *Node[T]) IsRoot() bool {
	if n != nil {
		return n.parent == nil
	}
	return true
}

// Name returns the node's display name
func (n *Node[T]) Name() string {
	return n.name
}

// SetName sets the node's display name
func (n *Node[T]) SetName(name string) {
	n.name = name
}

// Text returns the node's text (used for data access)
// Falls back to name if text is not set
func (n *Node[T]) Text() string {
	if n.text != "" {
		return n.text
	}
	return n.name
}

// SetText sets the node's text
func (n *Node[T]) SetText(text string) {
	n.text = text
}

// Data returns the node's payload
func (n *Node[T]) Data() *T {
	return &n.data
}

// Children returns the node's child nodes
func (n *Node[T]) Children() []*Node[T] {
	return n.children
}

// HasGrandChildren returns true if any children have children
func (n *Node[T]) HasGrandChildren() bool {
	if n.hasGrandChildren != nil {
		goto end
	}
	n.hasGrandChildren = new(bool)
	for _, child := range n.children {
		if len(child.children) == 0 {
			continue
		}
		*n.hasGrandChildren = true
		goto end
	}
	*n.hasGrandChildren = false
end:
	return *n.hasGrandChildren
}

// Parent returns the node's parent (nil for root nodes)
func (n *Node[T]) Parent() *Node[T] {
	return n.parent
}

// HasChildren returns true if the node has children
func (n *Node[T]) HasChildren() bool {
	return n != nil && len(n.children) > 0
}

// IsExpanded returns true if the node is expanded
func (n *Node[T]) IsExpanded() bool {
	return n.expanded
}

// IsVisible returns true if the node is visible
func (n *Node[T]) IsVisible() bool {
	return n.visible
}

// SetExpanded sets the node's expansion state
func (n *Node[T]) SetExpanded(expanded bool) {
	n.expanded = expanded
}

// SetVisible sets the node's visibility state
func (n *Node[T]) SetVisible(visible bool) {
	n.visible = visible
}

// Expand expands the node
func (n *Node[T]) Expand() {
	n.expanded = true
}

// Collapse collapses the node
func (n *Node[T]) Collapse() {
	n.expanded = false
}

// Toggle toggles the node's expansion state
func (n *Node[T]) Toggle() {
	n.expanded = !n.expanded
}

// AddChild adds a child node and sets the reciprocal parent pointer
func (n *Node[T]) AddChild(child *Node[T]) {
	child.parent = n
	n.children = append(n.children, child)
}

// SetChildren replaces all children and wires up parent pointers
func (n *Node[T]) SetChildren(children []*Node[T]) {
	n.children = children
	n.hasGrandChildren = nil
	if n.parent != nil {
		n.parent.hasGrandChildren = nil
	}
	for _, child := range children {
		child.parent = n
	}
}

// RemoveChild removes a child node by ID
func (n *Node[T]) RemoveChild(id string) bool {
	for i, child := range n.children {
		if child.id == id {
			// Remove by slicing
			n.children = append(n.children[:i], n.children[i+1:]...)
			child.parent = nil
			// Invalidate cached hasGrandChildren
			n.hasGrandChildren = nil
			return true
		}
	}
	return false
}

// InsertChildSorted inserts a child node in sorted order by name.
// Uses the provided less function to determine order.
// If less is nil, sorts alphabetically by name.
func (n *Node[T]) InsertChildSorted(child *Node[T], less func(a, b *Node[T]) bool) {
	child.parent = n

	if less == nil {
		less = func(a, b *Node[T]) bool {
			return a.name < b.name
		}
	}

	// Find insertion point
	insertIdx := len(n.children)
	for i, existing := range n.children {
		if less(child, existing) {
			insertIdx = i
			break
		}
	}

	// Insert at position
	n.children = append(n.children[:insertIdx], append([]*Node[T]{child}, n.children[insertIdx:]...)...)

	// Invalidate cached hasGrandChildren
	n.hasGrandChildren = nil
}

// FindByID recursively searches for a node by ID in this subtree
func (n *Node[T]) FindByID(id string) *Node[T] {
	if n.id == id {
		return n
	}

	for _, child := range n.children {
		found := child.FindByID(id)
		if found != nil {
			return found
		}
	}

	return nil
}

// Depth returns the depth of this node in the tree (0 for root)
func (n *Node[T]) Depth() int {
	depth := 0
	current := n.parent
	for current != nil {
		depth++
		current = current.parent
	}
	return depth
}

// IsLastChild returns true if this node is the last child of its parent
func (n *Node[T]) IsLastChild() bool {
	if n.parent == nil {
		return true
	}

	siblings := n.parent.children
	return len(siblings) > 0 && siblings[len(siblings)-1] == n
}

// AncestorIsLastChild returns a boolean slice indicating whether each ancestor was a last child
// Used for building tree structure prefixes
func (n *Node[T]) AncestorIsLastChild() []bool {
	var result []bool
	current := n.parent

	for current != nil && current.parent != nil {
		result = append([]bool{current.IsLastChild()}, result...)
		current = current.parent
	}

	return result
}
