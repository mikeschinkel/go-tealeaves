package teatree

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel() Model[string] {
	tree, _ := buildTestTree()
	return NewModel(tree, 10)
}

// --- Layer 1: Model Tests ---

func TestNewModel(t *testing.T) {
	m := newTestModel()
	if m.Tree() == nil {
		t.Fatal("expected tree to be set")
	}
	if m.height != 10 {
		t.Errorf("expected height=10, got %d", m.height)
	}
}

func TestModel_KeyUp(t *testing.T) {
	m := newTestModel()
	// Expand root1 so we have more nodes to navigate
	m.Tree().ExpandFocused()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to child1
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})   // Move back to root1

	if m.FocusedNode().ID() != "root1" {
		t.Errorf("expected focus on 'root1' after KeyUp, got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyDown(t *testing.T) {
	m := newTestModel()
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	if m.FocusedNode().ID() != "root2" {
		t.Errorf("expected focus on 'root2' after KeyDown, got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyRight_Expand(t *testing.T) {
	m := newTestModel()
	// Focus on root1 (has children, collapsed)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})

	if !m.FocusedNode().IsExpanded() {
		t.Error("expected focused node to be expanded after KeyRight")
	}
	// Children should be visible
	visible := m.Tree().VisibleNodes()
	found := false
	for _, n := range visible {
		if n.ID() == "child1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected child1 visible after expanding root1")
	}
}

func TestModel_KeyRight_EnterChild(t *testing.T) {
	m := newTestModel()
	// Expand root1 first
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	// Already expanded, KeyRight again should move to first child
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})

	if m.FocusedNode().ID() != "child1" {
		t.Errorf("expected focus on 'child1', got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyLeft_Collapse(t *testing.T) {
	m := newTestModel()
	// Expand root1
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	if !m.FocusedNode().IsExpanded() {
		t.Fatal("expected root1 expanded")
	}

	// Collapse with KeyLeft
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if m.FocusedNode().IsExpanded() {
		t.Error("expected root1 collapsed after KeyLeft")
	}
}

func TestModel_KeyLeft_MoveToParent(t *testing.T) {
	m := newTestModel()
	// Expand root1, move to child1
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight}) // expand root1
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight}) // move to child1

	if m.FocusedNode().ID() != "child1" {
		t.Fatalf("expected focus on child1, got %q", m.FocusedNode().ID())
	}

	// child1 is collapsed (no grandchildren expanded), KeyLeft should go to parent
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if m.FocusedNode().ID() != "root1" {
		t.Errorf("expected focus on 'root1' after KeyLeft from child1, got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyToggle(t *testing.T) {
	m := newTestModel()

	// Toggle (Enter/Space) should expand
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !m.FocusedNode().IsExpanded() {
		t.Error("expected node expanded after Enter toggle")
	}

	// Toggle again should collapse
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeySpace})
	if m.FocusedNode().IsExpanded() {
		t.Error("expected node collapsed after Space toggle")
	}
}

func TestModel_WindowSizeMsg(t *testing.T) {
	m := newTestModel()
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if m.width != 120 {
		t.Errorf("expected width=120, got %d", m.width)
	}
	if m.height != 40 {
		t.Errorf("expected height=40, got %d", m.height)
	}
}

func TestModel_SetSize(t *testing.T) {
	m := newTestModel()
	m = m.SetSize(100, 50)

	if m.width != 100 {
		t.Errorf("expected width=100, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("expected height=50, got %d", m.height)
	}
}

func TestModel_SetFocusedNode(t *testing.T) {
	m := newTestModel()
	m = m.SetFocusedNode("root2")

	if m.FocusedNode().ID() != "root2" {
		t.Errorf("expected focus on 'root2', got %q", m.FocusedNode().ID())
	}
}

// --- Layer 2: View Tests ---

func TestModel_View_BasicTree(t *testing.T) {
	m := newTestModel()
	view := m.View()

	if view == "" {
		t.Fatal("expected non-empty view")
	}
	if !strings.Contains(view, "Root1") {
		t.Error("expected view to contain 'Root1'")
	}
	if !strings.Contains(view, "Root2") {
		t.Error("expected view to contain 'Root2'")
	}
}

func TestModel_View_ExpandedTree(t *testing.T) {
	m := newTestModel()
	// Expand root1
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	view := m.View()

	if !strings.Contains(view, "Child1") {
		t.Error("expected view to contain 'Child1' after expansion")
	}
	if !strings.Contains(view, "Child2") {
		t.Error("expected view to contain 'Child2' after expansion")
	}
	// Tree prefix characters should be present
	if !strings.Contains(view, "├") && !strings.Contains(view, "└") {
		t.Error("expected tree prefix characters in expanded view")
	}
}

func TestModel_View_FocusedNode(t *testing.T) {
	m := newTestModel()
	view := m.View()

	// The focused node (Root1) should be visually distinct
	// We can't easily check style, but the node name should appear
	if !strings.Contains(view, "Root1") {
		t.Error("expected focused node 'Root1' in view")
	}
}

func TestModel_View_Scrolling(t *testing.T) {
	// Create a model with very small viewport
	tree, _ := buildTestTree()
	tree.ExpandAll()
	m := NewModel(tree, 3) // Only 3 lines visible

	view := m.View()
	lines := strings.Split(view, "\n")

	// View should be limited to viewport height
	if len(lines) > 3 {
		t.Errorf("expected at most 3 lines in view, got %d", len(lines))
	}

	// Navigate down several times to trigger scrolling
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	view = m.View()
	// Focused node should still be visible
	focusedName := m.FocusedNode().Name()
	if !strings.Contains(view, focusedName) {
		t.Errorf("expected focused node '%s' to be visible after scrolling", focusedName)
	}
}
