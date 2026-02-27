package teadep

// SelectorFunc is a function that chooses the best child to follow
// Given a parent tree node and its children, returns which child to select for the path
// This allows the caller to define what "best" means (longest path, deepest, most in-flux, etc.)
// The function can capture any necessary context (metadata, configuration, etc.) in a closure
type SelectorFunc func(parent *Tree, children []*Tree) (best *Tree, err error)
