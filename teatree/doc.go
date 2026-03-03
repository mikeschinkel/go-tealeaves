// Package teatree provides a generic Bubble Tea v2 component for rendering
// and navigating expandable/collapsible tree structures.
//
// Trees are built from [Node] values parameterized by a data type T. A
// pluggable [NodeProvider] interface controls how each node is rendered.
//
// Usage:
//
//	root := teatree.NewNode("1", "Root", myData)
//	tree := teatree.NewTree([]*teatree.Node[MyType]{root}, nil)
//	model := teatree.NewTreeModel(tree, 20)
//
// # Stability
//
// This package is provisional as of v0.3.0. The public API may change in
// minor releases until promoted to stable.
package teatree
