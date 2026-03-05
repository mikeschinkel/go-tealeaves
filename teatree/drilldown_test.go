package teatree

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teadrpdwn"
)

func extractMsg(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

// firstChildSelector always selects the first child at each level.
func firstChildSelector(_ *Node[string], children []*Node[string]) (*Node[string], error) {
	if len(children) == 0 {
		return nil, nil
	}
	return children[0], nil
}

// buildTestTree creates a tree for testing:
//
//	A
//	├── B
//	│   └── D (leaf)
//	└── C
//	    └── E (leaf)
func buildDrillDownTestTree() *Node[string] {
	nodeA := NewNode("a", "A", "root node")
	nodeB := NewNode("b", "B", "child B")
	nodeC := NewNode("c", "C", "child C")
	nodeD := NewNode("d", "D", "leaf D")
	nodeE := NewNode("e", "E", "leaf E")

	nodeA.AddChild(nodeB)
	nodeA.AddChild(nodeC)
	nodeB.AddChild(nodeD)
	nodeC.AddChild(nodeE)

	return nodeA
}

func newInitializedDrillDown() DrillDownModel[string] {
	root := buildDrillDownTestTree()
	m := NewDrillDownModel(root, DrillDownArgs[string]{
		SelectorFunc: firstChildSelector,
		Prompt:       "Drill Down",
	})
	m.Width = 80
	m.Height = 24
	m, _ = m.Initialize()
	return m
}

// --- Path Building Tests ---

func TestBuildDrillDownPath(t *testing.T) {
	root := buildDrillDownTestTree()
	path, err := buildDrillDownPath(root, firstChildSelector)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Path should be: A → B → D (first child at each level)
	if len(path) != 3 {
		t.Fatalf("expected path length 3, got %d", len(path))
	}
	if path[0].Name() != "A" {
		t.Errorf("expected path[0]='A', got %q", path[0].Name())
	}
	if path[1].Name() != "B" {
		t.Errorf("expected path[1]='B', got %q", path[1].Name())
	}
	if path[2].Name() != "D" {
		t.Errorf("expected path[2]='D', got %q", path[2].Name())
	}
}

func TestBuildDrillDownPath_NilRoot(t *testing.T) {
	_, err := buildDrillDownPath[string](nil, firstChildSelector)
	if err == nil {
		t.Error("expected error for nil root")
	}
}

func TestRebuildDrillDownPath(t *testing.T) {
	root := buildDrillDownTestTree()
	path, _ := buildDrillDownPath(root, firstChildSelector)

	// Path is A → B → D. Rebuild from level 1, choosing C instead of B
	nodeC := root.Children()[1] // C
	newPath, err := rebuildDrillDownPath(path, 1, nodeC, firstChildSelector)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// New path should be A → C → E
	if len(newPath) != 3 {
		t.Fatalf("expected path length 3, got %d", len(newPath))
	}
	if newPath[0].Name() != "A" {
		t.Errorf("expected path[0]='A', got %q", newPath[0].Name())
	}
	if newPath[1].Name() != "C" {
		t.Errorf("expected path[1]='C', got %q", newPath[1].Name())
	}
	if newPath[2].Name() != "E" {
		t.Errorf("expected path[2]='E', got %q", newPath[2].Name())
	}
}

func TestHasAlternatives(t *testing.T) {
	root := buildDrillDownTestTree()

	// Root has no alternatives (no parent)
	if hasAlternatives(root) {
		t.Error("expected root to have no alternatives")
	}

	// B has alternatives (sibling C)
	childB := root.Children()[0]
	if !hasAlternatives(childB) {
		t.Error("expected B to have alternatives")
	}

	// D has no alternatives (only child of B)
	leafD := childB.Children()[0]
	if hasAlternatives(leafD) {
		t.Error("expected D to have no alternatives (only child)")
	}
}

func TestAlternatives(t *testing.T) {
	root := buildDrillDownTestTree()

	// Root has nil alternatives
	if alternatives(root) != nil {
		t.Error("expected root to have nil alternatives")
	}

	// B and C are siblings
	childB := root.Children()[0]
	alts := alternatives(childB)
	if len(alts) != 2 {
		t.Fatalf("expected 2 alternatives, got %d", len(alts))
	}
	if alts[0].Name() != "B" {
		t.Errorf("expected first alt 'B', got %q", alts[0].Name())
	}
	if alts[1].Name() != "C" {
		t.Errorf("expected second alt 'C', got %q", alts[1].Name())
	}
}

// --- DrillDownModel Tests ---

func TestNewDrillDownModel(t *testing.T) {
	root := buildDrillDownTestTree()
	m := NewDrillDownModel(root, DrillDownArgs[string]{
		SelectorFunc: firstChildSelector,
	})

	if m.Root() == nil {
		t.Error("expected Root to be set")
	}
	if m.SelectorFunc == nil {
		t.Error("expected SelectorFunc to be set")
	}
	if len(m.Path) != 0 {
		t.Errorf("expected empty path before Initialize, got %d", len(m.Path))
	}
}

func TestDrillDownModel_Initialize(t *testing.T) {
	m := newInitializedDrillDown()

	// Path should be: A → B → D
	if len(m.Path) != 3 {
		t.Fatalf("expected path length 3, got %d", len(m.Path))
	}
	if m.Path[0].Name() != "A" {
		t.Errorf("expected path[0]='A', got %q", m.Path[0].Name())
	}
	if m.Path[1].Name() != "B" {
		t.Errorf("expected path[1]='B', got %q", m.Path[1].Name())
	}
	if m.Path[2].Name() != "D" {
		t.Errorf("expected path[2]='D', got %q", m.Path[2].Name())
	}
	if m.SelectedLevel != 2 {
		t.Errorf("expected SelectedLevel=2, got %d", m.SelectedLevel)
	}
}

func TestDrillDownModel_Initialize_NilSelector(t *testing.T) {
	root := buildDrillDownTestTree()
	m := NewDrillDownModel(root, DrillDownArgs[string]{
		SelectorFunc: nil,
	})
	_, err := m.Initialize()
	if err == nil {
		t.Error("expected error for nil selector")
	}
}

func TestDrillDownModel_KeyUp(t *testing.T) {
	m := newInitializedDrillDown()
	// SelectedLevel starts at 2 (leaf D)

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = result.(DrillDownModel[string])

	if m.SelectedLevel != 1 {
		t.Errorf("expected SelectedLevel=1 after KeyUp, got %d", m.SelectedLevel)
	}

	msg := extractMsg(cmd)
	focus, ok := msg.(DrillDownFocusMsg[string])
	if !ok {
		t.Fatalf("expected DrillDownFocusMsg, got %T", msg)
	}
	if focus.Level != 1 {
		t.Errorf("expected FocusMsg.Level=1, got %d", focus.Level)
	}

	// Move up again
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = result.(DrillDownModel[string])
	if m.SelectedLevel != 0 {
		t.Errorf("expected SelectedLevel=0, got %d", m.SelectedLevel)
	}

	// At top, should stay at 0
	result, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = result.(DrillDownModel[string])
	if m.SelectedLevel != 0 {
		t.Errorf("expected SelectedLevel=0 (clamped), got %d", m.SelectedLevel)
	}
	if cmd != nil {
		t.Error("expected nil cmd at top boundary")
	}
}

func TestDrillDownModel_KeyDown(t *testing.T) {
	m := newInitializedDrillDown()
	m.SelectedLevel = 0

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DrillDownModel[string])

	if m.SelectedLevel != 1 {
		t.Errorf("expected SelectedLevel=1 after KeyDown, got %d", m.SelectedLevel)
	}

	msg := extractMsg(cmd)
	focus, ok := msg.(DrillDownFocusMsg[string])
	if !ok {
		t.Fatalf("expected DrillDownFocusMsg, got %T", msg)
	}
	if focus.Level != 1 {
		t.Errorf("expected FocusMsg.Level=1, got %d", focus.Level)
	}

	// Move to end
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DrillDownModel[string])
	if m.SelectedLevel != 2 {
		t.Errorf("expected SelectedLevel=2, got %d", m.SelectedLevel)
	}

	// At bottom, should stay
	result, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DrillDownModel[string])
	if m.SelectedLevel != 2 {
		t.Errorf("expected SelectedLevel=2 (clamped), got %d", m.SelectedLevel)
	}
	if cmd != nil {
		t.Error("expected nil cmd at bottom boundary")
	}
}

func TestDrillDownModel_OpenDropdown(t *testing.T) {
	m := newInitializedDrillDown()
	m.SelectedLevel = 1 // B, which has alternatives (B, C)

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	m = result.(DrillDownModel[string])

	if !m.DropdownOpen {
		t.Error("expected DropdownOpen=true after Space on node with alternatives")
	}
}

func TestDrillDownModel_OpenDropdown_NoAlternatives(t *testing.T) {
	m := newInitializedDrillDown()
	m.SelectedLevel = 0 // Root A, which has no alternatives

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	m = result.(DrillDownModel[string])

	if m.DropdownOpen {
		t.Error("expected DropdownOpen=false when no alternatives exist")
	}
}

func TestDrillDownModel_EnterOnLeaf(t *testing.T) {
	m := newInitializedDrillDown()
	// SelectedLevel=2 (leaf D)

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	msg := extractMsg(cmd)
	sel, ok := msg.(DrillDownSelectMsg[string])
	if !ok {
		t.Fatalf("expected DrillDownSelectMsg, got %T", msg)
	}
	if sel.Node.Name() != "D" {
		t.Errorf("expected selected node 'D', got %q", sel.Node.Name())
	}
}

func TestDrillDownModel_EnterOnNonLeaf(t *testing.T) {
	m := newInitializedDrillDown()
	m.SelectedLevel = 0 // Root A, not a leaf

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if cmd != nil {
		t.Error("expected nil cmd for Enter on non-leaf node")
	}
}

func TestDrillDownModel_DropdownSelection(t *testing.T) {
	m := newInitializedDrillDown()
	m.SelectedLevel = 1 // B

	// Open dropdown
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	m = result.(DrillDownModel[string])
	if !m.DropdownOpen {
		t.Fatal("expected dropdown to be open")
	}

	// Select C (index=1)
	result, cmd := m.Update(teadrpdwn.OptionSelectedMsg{
		Index: 1,
		Text:  "C",
	})
	m = result.(DrillDownModel[string])

	if m.DropdownOpen {
		t.Error("expected dropdown closed after selection")
	}

	// Path should be rebuilt: A → C → E
	if len(m.Path) != 3 {
		t.Fatalf("expected path length 3 after rebuild, got %d", len(m.Path))
	}
	if m.Path[1].Name() != "C" {
		t.Errorf("expected path[1]='C' after selection, got %q", m.Path[1].Name())
	}
	if m.Path[2].Name() != "E" {
		t.Errorf("expected path[2]='E' after rebuild, got %q", m.Path[2].Name())
	}

	msg := extractMsg(cmd)
	change, ok := msg.(DrillDownChangeMsg[string])
	if !ok {
		t.Fatalf("expected DrillDownChangeMsg, got %T", msg)
	}
	if change.Level != 1 {
		t.Errorf("expected ChangeMsg.Level=1, got %d", change.Level)
	}
}

func TestDrillDownModel_DropdownCancellation(t *testing.T) {
	m := newInitializedDrillDown()
	m.SelectedLevel = 1

	// Open dropdown
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	m = result.(DrillDownModel[string])
	if !m.DropdownOpen {
		t.Fatal("expected dropdown open")
	}

	// Cancel
	result, _ = m.Update(teadrpdwn.DropdownCancelledMsg{})
	m = result.(DrillDownModel[string])

	if m.DropdownOpen {
		t.Error("expected dropdown closed after cancellation")
	}
	if m.Path[1].Name() != "B" {
		t.Errorf("expected path[1] still 'B', got %q", m.Path[1].Name())
	}
}

func TestDrillDownModel_WindowSizeMsg(t *testing.T) {
	m := newInitializedDrillDown()

	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = result.(DrillDownModel[string])

	if m.Width != 120 {
		t.Errorf("expected Width=120, got %d", m.Width)
	}
	if m.Height != 40 {
		t.Errorf("expected Height=40, got %d", m.Height)
	}
}

// --- View Tests ---

func TestDrillDownModel_View(t *testing.T) {
	m := newInitializedDrillDown()
	view := m.View()

	if view.Content == "" {
		t.Fatal("expected non-empty view")
	}
	if !strings.Contains(view.Content, "A") {
		t.Error("expected view to contain 'A'")
	}
	if !strings.Contains(view.Content, "B") {
		t.Error("expected view to contain 'B'")
	}
	if !strings.Contains(view.Content, "D") {
		t.Error("expected view to contain 'D'")
	}
	if !strings.Contains(view.Content, "Drill Down") {
		t.Error("expected view to contain prompt 'Drill Down'")
	}
}

func TestDrillDownModel_View_DropdownOpen(t *testing.T) {
	m := newInitializedDrillDown()
	m.SelectedLevel = 1

	// Open dropdown
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	m = result.(DrillDownModel[string])

	view := m.View()
	if !strings.Contains(view.Content, "B") {
		t.Error("expected dropdown to show 'B'")
	}
	if !strings.Contains(view.Content, "C") {
		t.Error("expected dropdown to show 'C'")
	}
}

func TestDrillDownModel_View_BorderGeometry(t *testing.T) {
	m := newInitializedDrillDown()
	m.Width = 40
	m.Height = 10
	view := m.View()

	lines := strings.Split(view.Content, "\n")
	if len(lines) > m.Height {
		t.Errorf("expected at most %d lines, got %d", m.Height, len(lines))
	}
}
