package teadepview

// BaseNode provides common functionality for node implementations
// Apps should embed this in their custom node types
type BaseNode struct {
	displayName  string
	dependencies []Node
}

type BaseNodeArgs struct {
	Dependencies []Node
}

func NewBaseNode(displayName string, args *BaseNodeArgs) *BaseNode {
	if args == nil {
		args = &BaseNodeArgs{}
	}
	return &BaseNode{
		displayName:  displayName,
		dependencies: args.Dependencies}
}

// DisplayName implements Node interface
func (n *BaseNode) DisplayName() string {
	return n.displayName
}

// SetDisplayName accepts a new name to set for display
func (n *BaseNode) SetDisplayName(name string) {
	n.displayName = name
}

// Dependencies implements Node interface
func (n *BaseNode) Dependencies() []Node {
	return n.dependencies
}

// SetDependencies accepts a list of nodes to set as the dependency property
func (n *BaseNode) SetDependencies(nodes []Node) {
	n.dependencies = nodes
}
