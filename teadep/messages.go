package teadep

// SelectNodeMsg sent when user confirms selection with Enter on a leaf node
type SelectNodeMsg struct {
	Tree *Tree
}

// ChangeNodeMsg sent when user picks different dependency from dropdown
type ChangeNodeMsg struct {
	Level int
	Tree  *Tree
}

// FocusNodeMsg sent when selection changes via navigation (up/down arrows)
type FocusNodeMsg struct {
	Level int
	Tree  *Tree
}
