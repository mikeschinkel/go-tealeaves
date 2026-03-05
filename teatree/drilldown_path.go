package teatree

// DrillDownSelectorFunc determines which child to follow at each level
// when building or rebuilding a drill-down path.
// It receives the current node and its children, and returns the selected child.
type DrillDownSelectorFunc[T any] func(current *Node[T], children []*Node[T]) (*Node[T], error)

// buildDrillDownPath constructs a path from root to leaf using the selector.
// The selector determines which child to follow at each level.
// Returns a slice of pointers into the node tree.
func buildDrillDownPath[T any](root *Node[T], selector DrillDownSelectorFunc[T]) (path []*Node[T], err error) {
	var current *Node[T]
	var bestChild *Node[T]
	var visited map[*Node[T]]bool

	if root == nil {
		err = NewErr(ErrDrillDown, ErrInvalidNode, "reason", "nil root")
		goto end
	}

	visited = make(map[*Node[T]]bool)
	path = []*Node[T]{root}
	visited[root] = true
	current = root

	// Walk down to leaf, using selector to choose best child at each level
	for current.HasChildren() {
		bestChild, err = selector(current, current.Children())
		if err != nil {
			goto end
		}
		if bestChild == nil {
			break
		}
		// Cycle detection
		if visited[bestChild] {
			break
		}
		visited[bestChild] = true
		path = append(path, bestChild)
		current = bestChild
	}

end:
	return path, err
}

// rebuildDrillDownPath rebuilds a path from the given level with a new node.
// Strategy: keep existingPath[0:level], then rebuild from newNode to leaf using selector.
func rebuildDrillDownPath[T any](existingPath []*Node[T], level int, newNode *Node[T], selector DrillDownSelectorFunc[T]) (newPath []*Node[T], err error) {
	var current *Node[T]
	var bestChild *Node[T]
	var visited map[*Node[T]]bool

	if level < 0 || level >= len(existingPath) {
		err = NewErr(ErrDrillDown, ErrInvalidLevel, "level", level, "max", len(existingPath)-1)
		goto end
	}

	// Keep path up to (but not including) selected level
	newPath = make([]*Node[T], level)
	copy(newPath, existingPath[:level])

	// Add new node at this level
	newPath = append(newPath, newNode)

	// Rebuild best path from new node to leaf using selector
	visited = make(map[*Node[T]]bool)
	for _, n := range newPath {
		visited[n] = true
	}

	current = newNode
	for current.HasChildren() {
		bestChild, err = selector(current, current.Children())
		if err != nil {
			goto end
		}
		if bestChild == nil {
			break
		}
		if visited[bestChild] {
			break
		}
		visited[bestChild] = true
		newPath = append(newPath, bestChild)
		current = bestChild
	}

end:
	return newPath, err
}

// hasAlternatives returns true if the node has siblings (parent has more than one child)
func hasAlternatives[T any](node *Node[T]) bool {
	if node.Parent() == nil {
		return false
	}
	return len(node.Parent().Children()) > 1
}

// alternatives returns all sibling nodes at this level (parent's children)
func alternatives[T any](node *Node[T]) []*Node[T] {
	if node.Parent() == nil {
		return nil
	}
	return node.Parent().Children()
}
