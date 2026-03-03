package teatree

import (
	"testing"
)

// --- Layer 1: Node Tests ---

func TestNewNode(t *testing.T) {
	n := NewNode[string]("id1", "Root", "payload")
	if n.ID() != "id1" {
		t.Errorf("expected ID='id1', got %q", n.ID())
	}
	if n.Name() != "Root" {
		t.Errorf("expected Name='Root', got %q", n.Name())
	}
	if *n.Data() != "payload" {
		t.Errorf("expected Data='payload', got %q", *n.Data())
	}
	if !n.IsRoot() {
		t.Error("expected IsRoot=true for new node")
	}
	if n.IsExpanded() {
		t.Error("expected IsExpanded=false initially")
	}
	if !n.IsVisible() {
		t.Error("expected IsVisible=true initially")
	}
	if n.HasChildren() {
		t.Error("expected HasChildren=false initially")
	}
}

func TestNode_AddChild(t *testing.T) {
	parent := NewNode[string]("p", "Parent", "")
	child := NewNode[string]("c", "Child", "")

	parent.AddChild(child)

	if len(parent.Children()) != 1 {
		t.Fatalf("expected 1 child, got %d", len(parent.Children()))
	}
	if parent.Children()[0] != child {
		t.Error("expected child to be in parent's children")
	}
	if child.Parent() != parent {
		t.Error("expected child's parent to be set")
	}
	if child.IsRoot() {
		t.Error("expected child IsRoot=false after AddChild")
	}
}

func TestNode_SetChildren(t *testing.T) {
	parent := NewNode[string]("p", "Parent", "")
	c1 := NewNode[string]("c1", "Child1", "")
	c2 := NewNode[string]("c2", "Child2", "")

	parent.SetChildren([]*Node[string]{c1, c2})

	if len(parent.Children()) != 2 {
		t.Fatalf("expected 2 children, got %d", len(parent.Children()))
	}
	if c1.Parent() != parent {
		t.Error("expected c1's parent to be set")
	}
	if c2.Parent() != parent {
		t.Error("expected c2's parent to be set")
	}
}

func TestNode_RemoveChild(t *testing.T) {
	parent := NewNode[string]("p", "Parent", "")
	c1 := NewNode[string]("c1", "Child1", "")
	c2 := NewNode[string]("c2", "Child2", "")
	parent.AddChild(c1)
	parent.AddChild(c2)

	removed := parent.RemoveChild("c1")
	if !removed {
		t.Error("expected RemoveChild to return true for existing child")
	}
	if len(parent.Children()) != 1 {
		t.Errorf("expected 1 child after removal, got %d", len(parent.Children()))
	}
	if parent.Children()[0] != c2 {
		t.Error("expected remaining child to be c2")
	}
	if c1.Parent() != nil {
		t.Error("expected removed child's parent to be nil")
	}

	// Unknown ID returns false
	removed = parent.RemoveChild("nonexistent")
	if removed {
		t.Error("expected RemoveChild to return false for unknown ID")
	}
}

func TestNode_InsertChildSorted(t *testing.T) {
	parent := NewNode[string]("p", "Parent", "")
	cB := NewNode[string]("b", "Bravo", "")
	cD := NewNode[string]("d", "Delta", "")
	parent.AddChild(cB)
	parent.AddChild(cD)

	// Insert "Charlie" between Bravo and Delta
	cC := NewNode[string]("c", "Charlie", "")
	parent.InsertChildSorted(cC, nil) // nil = alphabetical by name

	if len(parent.Children()) != 3 {
		t.Fatalf("expected 3 children, got %d", len(parent.Children()))
	}
	if parent.Children()[0].Name() != "Bravo" {
		t.Errorf("expected first child 'Bravo', got %q", parent.Children()[0].Name())
	}
	if parent.Children()[1].Name() != "Charlie" {
		t.Errorf("expected second child 'Charlie', got %q", parent.Children()[1].Name())
	}
	if parent.Children()[2].Name() != "Delta" {
		t.Errorf("expected third child 'Delta', got %q", parent.Children()[2].Name())
	}
	if cC.Parent() != parent {
		t.Error("expected inserted child's parent to be set")
	}
}

func TestNode_FindByID(t *testing.T) {
	root := NewNode[string]("r", "Root", "")
	child := NewNode[string]("c", "Child", "")
	grandchild := NewNode[string]("gc", "Grandchild", "")
	root.AddChild(child)
	child.AddChild(grandchild)

	found := root.FindByID("gc")
	if found == nil {
		t.Fatal("expected to find grandchild by ID")
	}
	if found.ID() != "gc" {
		t.Errorf("expected found node ID='gc', got %q", found.ID())
	}

	// Unknown ID returns nil
	notFound := root.FindByID("nonexistent")
	if notFound != nil {
		t.Error("expected nil for unknown ID")
	}
}

func TestNode_Depth(t *testing.T) {
	root := NewNode[string]("r", "Root", "")
	child := NewNode[string]("c", "Child", "")
	grandchild := NewNode[string]("gc", "Grandchild", "")
	root.AddChild(child)
	child.AddChild(grandchild)

	if root.Depth() != 0 {
		t.Errorf("expected root depth=0, got %d", root.Depth())
	}
	if child.Depth() != 1 {
		t.Errorf("expected child depth=1, got %d", child.Depth())
	}
	if grandchild.Depth() != 2 {
		t.Errorf("expected grandchild depth=2, got %d", grandchild.Depth())
	}
}

func TestNode_IsLastChild(t *testing.T) {
	root := NewNode[string]("r", "Root", "")
	c1 := NewNode[string]("c1", "First", "")
	c2 := NewNode[string]("c2", "Last", "")
	root.AddChild(c1)
	root.AddChild(c2)

	if c1.IsLastChild() {
		t.Error("expected first child IsLastChild=false")
	}
	if !c2.IsLastChild() {
		t.Error("expected last child IsLastChild=true")
	}
	// Root with no parent returns true
	if !root.IsLastChild() {
		t.Error("expected root IsLastChild=true (no parent)")
	}
}

func TestNode_AncestorIsLastChild(t *testing.T) {
	root := NewNode[string]("r", "Root", "")
	c1 := NewNode[string]("c1", "First", "")
	c2 := NewNode[string]("c2", "Second", "")
	gc := NewNode[string]("gc", "Grandchild", "")
	root.AddChild(c1)
	root.AddChild(c2)
	c1.AddChild(gc)

	result := gc.AncestorIsLastChild()

	// gc's parent (c1) is not last child → false
	if len(result) != 1 {
		t.Fatalf("expected 1 ancestor entry, got %d", len(result))
	}
	if result[0] {
		t.Error("expected ancestor (c1) is not last child")
	}
}

func TestNode_ExpandCollapse(t *testing.T) {
	n := NewNode[string]("n", "Node", "")
	if n.IsExpanded() {
		t.Error("expected not expanded initially")
	}

	n.Expand()
	if !n.IsExpanded() {
		t.Error("expected expanded after Expand()")
	}

	n.Collapse()
	if n.IsExpanded() {
		t.Error("expected not expanded after Collapse()")
	}

	n.Toggle()
	if !n.IsExpanded() {
		t.Error("expected expanded after Toggle() from collapsed")
	}

	n.Toggle()
	if n.IsExpanded() {
		t.Error("expected not expanded after Toggle() from expanded")
	}
}

func TestNode_HasGrandChildren(t *testing.T) {
	root := NewNode[string]("r", "Root", "")
	child := NewNode[string]("c", "Child", "")
	grandchild := NewNode[string]("gc", "Grandchild", "")
	root.AddChild(child)

	if root.HasGrandChildren() {
		t.Error("expected HasGrandChildren=false when child has no children")
	}

	child.AddChild(grandchild)
	// Reset cached value
	root.hasGrandChildren = nil
	if !root.HasGrandChildren() {
		t.Error("expected HasGrandChildren=true after adding grandchild")
	}
}

func TestNode_SetName(t *testing.T) {
	n := NewNode[string]("id1", "Original", "data")
	if n.Name() != "Original" {
		t.Errorf("expected Name='Original', got %q", n.Name())
	}

	n.SetName("Updated")
	if n.Name() != "Updated" {
		t.Errorf("expected Name='Updated' after SetName, got %q", n.Name())
	}
}

func TestNode_SetExpanded(t *testing.T) {
	n := NewNode[string]("id1", "Node", "data")
	if n.IsExpanded() {
		t.Error("expected not expanded initially")
	}

	n.SetExpanded(true)
	if !n.IsExpanded() {
		t.Error("expected expanded after SetExpanded(true)")
	}

	n.SetExpanded(false)
	if n.IsExpanded() {
		t.Error("expected not expanded after SetExpanded(false)")
	}
}

func TestNode_SetVisible(t *testing.T) {
	n := NewNode[string]("id1", "Node", "data")
	if !n.IsVisible() {
		t.Error("expected visible initially")
	}

	n.SetVisible(false)
	if n.IsVisible() {
		t.Error("expected not visible after SetVisible(false)")
	}

	n.SetVisible(true)
	if !n.IsVisible() {
		t.Error("expected visible after SetVisible(true)")
	}
}

func TestNode_Text_FallsBackToName(t *testing.T) {
	n := NewNode[string]("n", "NodeName", "")
	if n.Text() != "NodeName" {
		t.Errorf("expected Text() to fall back to name 'NodeName', got %q", n.Text())
	}

	n.SetText("CustomText")
	if n.Text() != "CustomText" {
		t.Errorf("expected Text()='CustomText', got %q", n.Text())
	}
}
