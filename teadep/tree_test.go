package teadep

import (
	"testing"
)

// buildTestDepTree creates a dependency tree for testing:
//
//	A
//	├── B
//	│   └── D (leaf)
//	└── C
//	    └── E (leaf)
func buildTestDepTree() (*Tree, *BaseNode) {
	nodeD := NewBaseNode("D", nil)
	nodeE := NewBaseNode("E", nil)
	nodeB := NewBaseNode("B", &BaseNodeArgs{Dependencies: []Node{nodeD}})
	nodeC := NewBaseNode("C", &BaseNodeArgs{Dependencies: []Node{nodeE}})
	nodeA := NewBaseNode("A", &BaseNodeArgs{Dependencies: []Node{nodeB, nodeC}})

	tree := NewTree(nodeA)
	return tree, nodeA
}

// --- Layer 1: Tree Tests ---

func TestNewTree(t *testing.T) {
	tree, nodeA := buildTestDepTree()

	if tree == nil {
		t.Fatal("expected non-nil tree")
	}
	if tree.Node != nodeA {
		t.Error("expected root node to match")
	}
	if len(tree.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(tree.Children))
	}
	if tree.Children[0].Node.DisplayName() != "B" {
		t.Errorf("expected first child 'B', got %q", tree.Children[0].Node.DisplayName())
	}
	if tree.Children[1].Node.DisplayName() != "C" {
		t.Errorf("expected second child 'C', got %q", tree.Children[1].Node.DisplayName())
	}
	// Check grandchildren
	if len(tree.Children[0].Children) != 1 {
		t.Fatalf("expected B to have 1 child, got %d", len(tree.Children[0].Children))
	}
	if tree.Children[0].Children[0].Node.DisplayName() != "D" {
		t.Errorf("expected B's child 'D', got %q", tree.Children[0].Children[0].Node.DisplayName())
	}
}

func TestTree_Alternatives(t *testing.T) {
	tree, _ := buildTestDepTree()

	// B and C are siblings (children of A)
	childB := tree.Children[0]
	alts := childB.Alternatives()
	if len(alts) != 2 {
		t.Fatalf("expected 2 alternatives, got %d", len(alts))
	}
	if alts[0].Node.DisplayName() != "B" {
		t.Errorf("expected first alt 'B', got %q", alts[0].Node.DisplayName())
	}
	if alts[1].Node.DisplayName() != "C" {
		t.Errorf("expected second alt 'C', got %q", alts[1].Node.DisplayName())
	}

	// Root has no alternatives
	rootAlts := tree.Alternatives()
	if rootAlts != nil {
		t.Error("expected root to have nil alternatives")
	}
}

func TestTree_IsLeaf(t *testing.T) {
	tree, _ := buildTestDepTree()

	// Root is not a leaf
	if tree.IsLeaf() {
		t.Error("expected root to not be a leaf")
	}

	// D is a leaf (no children)
	leafD := tree.Children[0].Children[0]
	if !leafD.IsLeaf() {
		t.Error("expected D to be a leaf")
	}
}

func TestTree_HasAlternatives(t *testing.T) {
	tree, _ := buildTestDepTree()

	// B has alternatives (sibling C)
	childB := tree.Children[0]
	if !childB.HasAlternatives() {
		t.Error("expected B to have alternatives")
	}

	// Root has no alternatives (no parent)
	if tree.HasAlternatives() {
		t.Error("expected root to have no alternatives")
	}

	// D has no alternatives (only child of B)
	leafD := tree.Children[0].Children[0]
	if leafD.HasAlternatives() {
		t.Error("expected D to have no alternatives (only child)")
	}
}

func TestNewTree_CircularDeps(t *testing.T) {
	// Create circular: A → B → A
	nodeA := NewBaseNode("A", nil)
	nodeB := NewBaseNode("B", nil)
	nodeA.SetDependencies([]Node{nodeB})
	nodeB.SetDependencies([]Node{nodeA})

	// Should not infinite loop
	tree := NewTree(nodeA)

	if tree == nil {
		t.Fatal("expected non-nil tree even with circular deps")
	}
	if tree.Node.DisplayName() != "A" {
		t.Errorf("expected root 'A', got %q", tree.Node.DisplayName())
	}
	// B should be a child, but B's child A should be pruned (circular)
	if len(tree.Children) != 1 {
		t.Fatalf("expected 1 child (B), got %d", len(tree.Children))
	}
	childB := tree.Children[0]
	if childB.Node.DisplayName() != "B" {
		t.Errorf("expected child 'B', got %q", childB.Node.DisplayName())
	}
	// B's dependency on A should be pruned
	if len(childB.Children) != 0 {
		t.Errorf("expected B to have 0 children (circular pruned), got %d", len(childB.Children))
	}
}
