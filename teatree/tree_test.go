package teatree

import (
	"testing"
)

// buildTestTree creates a tree for testing:
//
//	Root1
//	  ├── Child1
//	  │   ├── Grandchild1
//	  │   └── Grandchild2
//	  └── Child2
//	Root2
func buildTestTree() (*Tree[string], []*Node[string]) {
	root1 := NewNode[string]("root1", "Root1", "r1")
	child1 := NewNode[string]("child1", "Child1", "c1")
	child2 := NewNode[string]("child2", "Child2", "c2")
	gc1 := NewNode[string]("gc1", "Grandchild1", "gc1")
	gc2 := NewNode[string]("gc2", "Grandchild2", "gc2")
	root2 := NewNode[string]("root2", "Root2", "r2")

	root1.AddChild(child1)
	root1.AddChild(child2)
	child1.AddChild(gc1)
	child1.AddChild(gc2)

	nodes := []*Node[string]{root1, root2}
	tree := NewTree(nodes, nil)
	return tree, nodes
}

// --- Layer 1: Tree Tests ---

func TestNewTree(t *testing.T) {
	tree, nodes := buildTestTree()

	if len(tree.Nodes()) != 2 {
		t.Fatalf("expected 2 root nodes, got %d", len(tree.Nodes()))
	}
	if tree.Nodes()[0] != nodes[0] {
		t.Error("expected first node to match")
	}

	// Focus should auto-set to first visible node
	focused := tree.FocusedNode()
	if focused == nil {
		t.Fatal("expected focused node to be set")
	}
	if focused.ID() != "root1" {
		t.Errorf("expected focused on 'root1', got %q", focused.ID())
	}
}

func TestTree_MoveUp(t *testing.T) {
	tree, _ := buildTestTree()

	// Focus starts on root1; MoveUp should return false (already at top)
	moved := tree.MoveUp()
	if moved {
		t.Error("expected MoveUp to return false at top")
	}

	// Move down first, then back up
	tree.MoveDown()
	moved = tree.MoveUp()
	if !moved {
		t.Error("expected MoveUp to return true")
	}
	if tree.FocusedNode().ID() != "root1" {
		t.Errorf("expected focus on 'root1' after MoveUp, got %q", tree.FocusedNode().ID())
	}
}

func TestTree_MoveDown(t *testing.T) {
	tree, _ := buildTestTree()

	// Initially collapsed, only root nodes are visible
	moved := tree.MoveDown()
	if !moved {
		t.Error("expected MoveDown to return true")
	}
	if tree.FocusedNode().ID() != "root2" {
		t.Errorf("expected focus on 'root2', got %q", tree.FocusedNode().ID())
	}

	// At bottom, should return false
	moved = tree.MoveDown()
	if moved {
		t.Error("expected MoveDown to return false at bottom")
	}
}

func TestTree_ExpandFocused(t *testing.T) {
	tree, _ := buildTestTree()

	// Focus on root1, expand it
	expanded := tree.ExpandFocused()
	if !expanded {
		t.Error("expected ExpandFocused to return true")
	}
	if !tree.FocusedNode().IsExpanded() {
		t.Error("expected focused node to be expanded")
	}

	// Children should now be visible
	visible := tree.VisibleNodes()
	found := false
	for _, n := range visible {
		if n.ID() == "child1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected child1 to be visible after expanding root1")
	}
}

func TestTree_CollapseFocused(t *testing.T) {
	tree, _ := buildTestTree()

	// Expand first, then collapse
	tree.ExpandFocused()
	collapsed := tree.CollapseFocused()
	if !collapsed {
		t.Error("expected CollapseFocused to return true")
	}
	if tree.FocusedNode().IsExpanded() {
		t.Error("expected focused node to be collapsed")
	}
}

func TestTree_ToggleFocused(t *testing.T) {
	tree, _ := buildTestTree()

	toggled := tree.ToggleFocused()
	if !toggled {
		t.Error("expected ToggleFocused to return true")
	}
	if !tree.FocusedNode().IsExpanded() {
		t.Error("expected focused node expanded after toggle")
	}

	toggled = tree.ToggleFocused()
	if !toggled {
		t.Error("expected ToggleFocused to return true on second toggle")
	}
	if tree.FocusedNode().IsExpanded() {
		t.Error("expected focused node collapsed after second toggle")
	}
}

func TestTree_ExpandAll(t *testing.T) {
	tree, _ := buildTestTree()

	tree.ExpandAll()

	// All nodes with children should be expanded
	visible := tree.VisibleNodes()
	// root1, child1, gc1, gc2, child2, root2 = 6 visible nodes
	if len(visible) != 6 {
		t.Errorf("expected 6 visible nodes after ExpandAll, got %d", len(visible))
	}
}

func TestTree_CollapseAll(t *testing.T) {
	tree, _ := buildTestTree()

	tree.ExpandAll()
	tree.CollapseAll()

	// Only root nodes should be visible
	visible := tree.VisibleNodes()
	if len(visible) != 2 {
		t.Errorf("expected 2 visible nodes after CollapseAll, got %d", len(visible))
	}
}

func TestTree_VisibleNodes(t *testing.T) {
	tree, _ := buildTestTree()

	// Initially only roots are visible (nothing expanded)
	visible := tree.VisibleNodes()
	if len(visible) != 2 {
		t.Errorf("expected 2 visible nodes initially, got %d", len(visible))
	}

	// Expand root1 → root1, child1, child2, root2
	tree.ExpandFocused()
	visible = tree.VisibleNodes()
	if len(visible) != 4 {
		t.Errorf("expected 4 visible nodes after expanding root1, got %d", len(visible))
	}
}

func TestTree_SetFocusedNode(t *testing.T) {
	tree, _ := buildTestTree()

	found := tree.SetFocusedNode("root2")
	if !found {
		t.Error("expected SetFocusedNode to return true for existing ID")
	}
	if tree.FocusedNode().ID() != "root2" {
		t.Errorf("expected focus on 'root2', got %q", tree.FocusedNode().ID())
	}

	// Unknown ID returns false
	found = tree.SetFocusedNode("nonexistent")
	if found {
		t.Error("expected SetFocusedNode to return false for unknown ID")
	}
}

func TestTree_FindByID(t *testing.T) {
	tree, _ := buildTestTree()

	node := tree.FindByID("gc1")
	if node == nil {
		t.Fatal("expected to find gc1")
	}
	if node.Name() != "Grandchild1" {
		t.Errorf("expected name='Grandchild1', got %q", node.Name())
	}

	// Unknown ID
	notFound := tree.FindByID("nonexistent")
	if notFound != nil {
		t.Error("expected nil for unknown ID")
	}
}
