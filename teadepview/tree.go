package teadepview

// Tree is a navigation wrapper around a Node
// Provides parent/child relationships for tree traversal
type Tree struct {
	Node     Node    // The actual dependency
	Parent   *Tree   // Parent in THIS traversal (nil for root)
	Children []*Tree // Tree wrappers for each dependency (built upfront or on-demand)
}

// NewTree creates a Tree structure from a Node, building all Children upfront
// Detects and breaks circular dependencies by checking if a node is in the ancestor chain
func NewTree(node Node) *Tree {
	return buildTree(node, nil)
}

// buildTree creates a Tree structure from a Node, building all Children upfront
// Detects and breaks circular dependencies by checking if a node is in the ancestor chain
func buildTree(node Node, parent *Tree) *Tree {
	if node == nil {
		return nil
	}

	// Check for circular dependencies - if this node is already in ancestor chain, stop
	if parent.isAncestor(node) {
		return nil
	}

	tree := &Tree{
		Node:   node,
		Parent: parent,
	}

	// Build children upfront
	dependencies := node.Dependencies()
	tree.Children = make([]*Tree, 0, len(dependencies))
	for _, dep := range dependencies {
		child := buildTree(dep, tree)
		if child != nil {
			tree.Children = append(tree.Children, child)
		}
	}

	return tree
}

// isAncestor checks if a node appears in the ancestor chain
// This prevents circular dependencies from causing infinite recursion
func (t *Tree) isAncestor(node Node) (inChain bool) {
	for cur := t; cur != nil; cur = cur.Parent {
		if cur.Node != node {
			continue
		}
		inChain = true
		goto end
	}
end:
	return inChain
}

// Alternatives returns all options at this level (siblings)
// Returns Parent's children, which includes this node
func (t *Tree) Alternatives() []*Tree {
	if t.Parent == nil {
		return nil // Root has no alternatives
	}
	return t.Parent.Children
}

// IsLeaf returns true if node has no dependencies
func (t *Tree) IsLeaf() bool {
	return t.Node == nil || len(t.Children) == 0
}

// HasAlternatives returns true if there are other options at this level
func (t *Tree) HasAlternatives() bool {
	if t.Parent == nil {
		return false // Root has no alternatives
	}
	return len(t.Parent.Children) > 1
}

// buildPath constructs a path from root to leaf using the model's selector
// The selector determines which child to follow at each level
// Returns a slice of pointers into the tree structure
func (t *Tree) buildPath(m PathViewerModel) (path []*Tree, err error) {
	var current *Tree
	var bestChild *Tree

	if t == nil || t.Node == nil {
		err = NewErr(ErrDependency, ErrInvalidNode, "reason", "nil t or node")
		goto end
	}

	path = []*Tree{t}
	current = t

	// Walk down to leaf, using selector to choose best child at each level
	for len(current.Children) > 0 {
		bestChild, err = m.SelectorFunc(current, current.Children)
		if err != nil {
			goto end
		}

		path = append(path, bestChild)
		current = bestChild
	}

end:
	return path, err
}

// rebuildPath rebuilds path from selected level with new tree node using model state
// Strategy: Keep model.Path[0:model.SelectedLevel], then rebuild from this node to leaf using selector
// This is the core of your suggestion - preserve 0..<n>-1, rebuild from <n>
func (t *Tree) rebuildPath(m PathViewerModel) (newPath []*Tree, err error) {
	var current *Tree
	var bestChild *Tree

	level := m.SelectedLevel

	if level < 0 || level >= len(m.Path) {
		err = NewErr(ErrDependency, ErrInvalidLevel, "level", level, "max", len(m.Path)-1)
		goto end
	}

	// Keep path up to (but not including) selected level: m.Path[0:level]
	// This preserves indices 0 through level-1
	newPath = make([]*Tree, level)
	copy(newPath, m.Path[:level])

	// Add new node at this level
	newPath = append(newPath, t)

	// Rebuild best path from new node to leaf using selector
	current = t
	for len(current.Children) > 0 {
		bestChild, err = m.SelectorFunc(current, current.Children)
		if err != nil {
			goto end
		}

		newPath = append(newPath, bestChild)
		current = bestChild
	}

end:
	return newPath, err
}
