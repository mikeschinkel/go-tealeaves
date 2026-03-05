package teatree

// DrillDownSelectMsg is sent when the user confirms selection with Enter on a leaf node
type DrillDownSelectMsg[T any] struct {
	Node *Node[T]
}

// DrillDownChangeMsg is sent when the user picks a different node from the dropdown
type DrillDownChangeMsg[T any] struct {
	Level int
	Node  *Node[T]
}

// DrillDownFocusMsg is sent when the selection changes via navigation (up/down arrows)
type DrillDownFocusMsg[T any] struct {
	Level int
	Node  *Node[T]
}
