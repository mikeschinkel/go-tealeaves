package teadepview

// Node defines the interface for dependency nodes
// This allows different node types (repos, modules, etc.)
type Node interface {
	DisplayName() string
	SetDisplayName(string)
	Dependencies() []Node
	SetDependencies([]Node)
}
