package main

import "github.com/mikeschinkel/go-tealeaves/teadepview"

// nodeMeta contains pre-computed information about a tree node
type nodeMeta struct {
	ReferenceCount int  // How many nodes depend on this
	MaxDepth       int  // Longest path from this node to any leaf
	ModuleKind     int  // Module kind: exe, lib, test, etc.
	ModuleKindSet  bool // True if ModuleKind has been set
}

// analyzeTree performs a complete analysis of the tree and returns metadata for each node
func analyzeTree(t *teadepview.Tree, determineKind func(teadepview.Node) (int, bool)) map[*teadepview.Tree]*nodeMeta {
	meta := make(map[*teadepview.Tree]*nodeMeta)

	// Pass 1: Calculate max depths (bottom-up)
	calculateDepths(t, meta)

	// Pass 2: Count references (top-down, tracking visited to handle shared dependencies)
	visited := make(map[*teadepview.Tree]bool)
	countReferences(t, meta, visited)

	// Pass 3: Determine module kinds (if callback provided)
	if determineKind != nil {
		determineKinds(t, meta, determineKind)
	}

	return meta
}

// calculateDepths computes MaxDepth for each node (bottom-up traversal)
func calculateDepths(t *teadepview.Tree, meta map[*teadepview.Tree]*nodeMeta) (maxDepth int) {
	var childDepth int
	var child *teadepview.Tree

	if t == nil {
		goto end
	}

	// Ensure meta entry exists
	if meta[t] == nil {
		meta[t] = &nodeMeta{}
	}

	// Base case: leaf node has depth 0
	if len(t.Children) == 0 {
		meta[t].MaxDepth = 0
		maxDepth = 0
		goto end
	}

	// Recursive case: max of children's depths + 1
	maxDepth = 0
	for _, child = range t.Children {
		childDepth = calculateDepths(child, meta)
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	}
	maxDepth++
	meta[t].MaxDepth = maxDepth

end:
	return maxDepth
}

// countReferences counts how many times each node is referenced in the tree
func countReferences(t *teadepview.Tree, meta map[*teadepview.Tree]*nodeMeta, visited map[*teadepview.Tree]bool) {
	var child *teadepview.Tree

	if t == nil {
		goto end
	}

	// Skip if we've already visited this node in this traversal path
	if visited[t] {
		goto end
	}

	// Mark as visited for this path
	visited[t] = true

	// Ensure meta entry exists
	if meta[t] == nil {
		meta[t] = &nodeMeta{}
	}

	// Increment reference count
	meta[t].ReferenceCount++

	// Recurse to children
	for _, child = range t.Children {
		countReferences(child, meta, visited)
	}

	// Unmark for backtracking (allows counting in different paths)
	visited[t] = false

end:
	return
}

// determineKinds sets ModuleKind for each node using the provided callback
func determineKinds(t *teadepview.Tree, meta map[*teadepview.Tree]*nodeMeta, determiner func(teadepview.Node) (int, bool)) {
	var child *teadepview.Tree
	var kind int
	var ok bool

	if t == nil {
		goto end
	}

	// Ensure metadata entry exists
	if meta[t] == nil {
		meta[t] = &nodeMeta{}
	}

	// Use callback to determine kind
	kind, ok = determiner(t.Node)
	if ok {
		meta[t].ModuleKind = kind
		meta[t].ModuleKindSet = true
	}

	// Recurse to children
	for _, child = range t.Children {
		determineKinds(child, meta, determiner)
	}

end:
	return
}
