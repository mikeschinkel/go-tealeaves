package teatree

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestModel() TreeModel[string] {
	tree, _ := buildTestTree()
	return NewTreeModel(tree, 10)
}

// --- Layer 1: Model Tests ---

func TestNewTreeModel(t *testing.T) {
	m := newTestModel()
	if m.Tree() == nil {
		t.Fatal("expected tree to be set")
	}
	if m.height != 10 {
		t.Errorf("expected height=10, got %d", m.height)
	}
}

func TestModel_Init(t *testing.T) {
	m := newTestModel()
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init() to return nil")
	}
}

func TestModel_Update_UnknownMsg(t *testing.T) {
	m := newTestModel()
	// A custom message type should be delegated to viewport without error
	type customMsg struct{}
	m, cmd := m.Update(customMsg{})
	// Should not panic, and model should remain functional
	if m.Tree() == nil {
		t.Error("expected tree to still be set after unknown msg")
	}
	// cmd may or may not be nil (viewport may return something), but no panic
	_ = cmd
}

func TestModel_KeyUp(t *testing.T) {
	m := newTestModel()
	// Expand root1 so we have more nodes to navigate
	m.Tree().ExpandFocused()
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown}) // Move to child1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})   // Move back to root1

	if m.FocusedNode().ID() != "root1" {
		t.Errorf("expected focus on 'root1' after KeyUp, got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyDown(t *testing.T) {
	m := newTestModel()
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	if m.FocusedNode().ID() != "root2" {
		t.Errorf("expected focus on 'root2' after KeyDown, got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyRight_Expand(t *testing.T) {
	m := newTestModel()
	// Focus on root1 (has children, collapsed)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})

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
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	// Already expanded, KeyRight again should move to first child
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})

	if m.FocusedNode().ID() != "child1" {
		t.Errorf("expected focus on 'child1', got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyLeft_Collapse(t *testing.T) {
	m := newTestModel()
	// Expand root1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if !m.FocusedNode().IsExpanded() {
		t.Fatal("expected root1 expanded")
	}

	// Collapse with KeyLeft
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.FocusedNode().IsExpanded() {
		t.Error("expected root1 collapsed after KeyLeft")
	}
}

func TestModel_KeyLeft_MoveToParent(t *testing.T) {
	m := newTestModel()
	// Expand root1, move to child1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight}) // expand root1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight}) // move to child1

	if m.FocusedNode().ID() != "child1" {
		t.Fatalf("expected focus on child1, got %q", m.FocusedNode().ID())
	}

	// child1 is collapsed (no grandchildren expanded), KeyLeft should go to parent
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.FocusedNode().ID() != "root1" {
		t.Errorf("expected focus on 'root1' after KeyLeft from child1, got %q", m.FocusedNode().ID())
	}
}

func TestModel_KeyToggle(t *testing.T) {
	m := newTestModel()

	// Toggle (Enter/Space) should expand
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !m.FocusedNode().IsExpanded() {
		t.Error("expected node expanded after Enter toggle")
	}

	// Toggle again should collapse
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
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

// --- Migration-sensitive tests (v1→v2 regression guards) ---

// TRE-VIEWPORT: Guards m.viewport.Width = msg.Width and m.viewport.Height = msg.Height
// direct field assignments in the WindowSizeMsg handler (model.go:98-99).
// After a WindowSizeMsg, viewport dimensions must match and view output must
// respect the new height.
func TestModel_WindowSizeMsg_ViewportSync(t *testing.T) {
	tree, _ := buildTestTree()
	tree.ExpandAll() // 6 visible nodes
	m := NewTreeModel(tree, 10)

	// Send WindowSizeMsg with specific dimensions
	m, _ = m.Update(tea.WindowSizeMsg{Width: 60, Height: 4})

	if m.width != 60 {
		t.Errorf("expected width=60, got %d", m.width)
	}
	if m.height != 4 {
		t.Errorf("expected height=4, got %d", m.height)
	}

	// Viewport should be synced — view should show at most 4 lines
	view := m.View()
	lines := strings.Split(view.Content, "\n")
	if len(lines) > 4 {
		t.Errorf("expected at most 4 visible lines after WindowSizeMsg(height=4), got %d", len(lines))
	}
	// Must still contain content
	if view.Content == "" {
		t.Error("expected non-empty view after WindowSizeMsg")
	}
}

// TRE-SCROLL-RESIZE: Guards m.viewport.Height field read in View() (model.go:122)
// and ensureFocusedVisible(). After scrolling in a small viewport, then resizing,
// the view must respect the new height (viewport.Height drives the slice in View()).
func TestModel_ScrollAfterResize(t *testing.T) {
	tree, _ := buildTestTree()
	tree.ExpandAll() // 6 visible nodes: Root1, Child1, GC1, GC2, Child2, Root2
	m := NewTreeModel(tree, 3)
	m = m.SetSize(80, 3) // Only 3 lines visible

	// Navigate down 4 times to scroll the viewport
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	focusedName := m.FocusedNode().Name()
	view := m.View()

	// Focused node must be visible in the viewport after scroll
	if !strings.Contains(view.Content,focusedName) {
		t.Errorf("expected focused node %q visible after scrolling, view=%q",
			focusedName, view.Content)
	}

	// View should respect the 3-line height limit
	lines := strings.Split(view.Content, "\n")
	if len(lines) > 3 {
		t.Errorf("expected at most 3 lines, got %d", len(lines))
	}

	// Resize to 2 lines — viewport.Height must be updated so View() slices correctly
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 2})
	view = m.View()
	lines = strings.Split(view.Content, "\n")
	if len(lines) > 2 {
		t.Errorf("expected at most 2 lines after resize to height=2, got %d", len(lines))
	}
	if view.Content == "" {
		t.Error("expected non-empty view after resize")
	}

	// Navigate after resize to verify ensureFocusedVisible still works
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	view = m.View()
	newFocused := m.FocusedNode().Name()
	if !strings.Contains(view.Content,newFocused) {
		t.Errorf("expected focused node %q visible after navigate post-resize, view=%q",
			newFocused, view.Content)
	}
}

// --- Layer 2: View Tests ---

func TestModel_View_BasicTree(t *testing.T) {
	m := newTestModel()
	view := m.View()

	if view.Content == "" {
		t.Fatal("expected non-empty view")
	}
	if !strings.Contains(view.Content,"Root1") {
		t.Error("expected view to contain 'Root1'")
	}
	if !strings.Contains(view.Content,"Root2") {
		t.Error("expected view to contain 'Root2'")
	}
}

func TestModel_View_ExpandedTree(t *testing.T) {
	m := newTestModel()
	// Expand root1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	view := m.View()

	if !strings.Contains(view.Content,"Child1") {
		t.Error("expected view to contain 'Child1' after expansion")
	}
	if !strings.Contains(view.Content,"Child2") {
		t.Error("expected view to contain 'Child2' after expansion")
	}
	// Tree prefix characters should be present
	if !strings.Contains(view.Content,"├") && !strings.Contains(view.Content,"└") {
		t.Error("expected tree prefix characters in expanded view")
	}
}

func TestModel_View_FocusedNode(t *testing.T) {
	m := newTestModel()
	view := m.View()

	// The focused node (Root1) should be visually distinct
	// We can't easily check style, but the node name should appear
	if !strings.Contains(view.Content,"Root1") {
		t.Error("expected focused node 'Root1' in view")
	}
}

func TestModel_View_Scrolling(t *testing.T) {
	// Create a model with very small viewport
	tree, _ := buildTestTree()
	tree.ExpandAll()
	m := NewTreeModel(tree, 3) // Only 3 lines visible

	view := m.View()
	lines := strings.Split(view.Content, "\n")

	// View should be limited to viewport height
	if len(lines) > 3 {
		t.Errorf("expected at most 3 lines in view, got %d", len(lines))
	}

	// Navigate down several times to trigger scrolling
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})

	view = m.View()
	// Focused node should still be visible
	focusedName := m.FocusedNode().Name()
	if !strings.Contains(view.Content,focusedName) {
		t.Errorf("expected focused node '%s' to be visible after scrolling", focusedName)
	}
}
