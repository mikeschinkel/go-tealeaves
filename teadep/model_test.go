package teadep

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikeschinkel/go-tealeaves/teadd"
)

func extractMsg(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

// firstChildSelector always selects the first child at each level.
func firstChildSelector(_ *Tree, children []*Tree) (*Tree, error) {
	if len(children) == 0 {
		return nil, nil
	}
	return children[0], nil
}

func newInitializedPathViewer() PathViewerModel {
	tree, _ := buildTestDepTree()
	m := NewPathViewer(tree, PathViewerArgs{
		SelectorFunc: firstChildSelector,
		Prompt:       "Dependencies",
	})
	m.Width = 80
	m.Height = 24
	m, _ = m.Initialize()
	return m
}

// --- Layer 1: PathViewer Tests ---

func TestNewPathViewer(t *testing.T) {
	tree, _ := buildTestDepTree()
	m := NewPathViewer(tree, PathViewerArgs{
		SelectorFunc: firstChildSelector,
	})

	if m.Root == nil {
		t.Error("expected Root to be set")
	}
	if m.SelectorFunc == nil {
		t.Error("expected SelectorFunc to be set")
	}
	// Path should not be built yet
	if len(m.Path) != 0 {
		t.Errorf("expected empty path before Initialize, got %d", len(m.Path))
	}
}

func TestPathViewer_Initialize(t *testing.T) {
	m := newInitializedPathViewer()

	// Path should be built: A → B → D (first child at each level)
	if len(m.Path) != 3 {
		t.Fatalf("expected path length 3, got %d", len(m.Path))
	}
	if m.Path[0].Node.DisplayName() != "A" {
		t.Errorf("expected path[0]='A', got %q", m.Path[0].Node.DisplayName())
	}
	if m.Path[1].Node.DisplayName() != "B" {
		t.Errorf("expected path[1]='B', got %q", m.Path[1].Node.DisplayName())
	}
	if m.Path[2].Node.DisplayName() != "D" {
		t.Errorf("expected path[2]='D', got %q", m.Path[2].Node.DisplayName())
	}
	// SelectedLevel should be at leaf (last)
	if m.SelectedLevel != 2 {
		t.Errorf("expected SelectedLevel=2, got %d", m.SelectedLevel)
	}
}

func TestPathViewer_Initialize_NilSelector(t *testing.T) {
	tree, _ := buildTestDepTree()
	m := NewPathViewer(tree, PathViewerArgs{
		SelectorFunc: nil,
	})
	_, err := m.Initialize()
	if err == nil {
		t.Error("expected error for nil selector")
	}
}

func TestPathViewer_KeyUp(t *testing.T) {
	m := newInitializedPathViewer()
	// SelectedLevel starts at 2 (leaf D)

	result, cmd := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = result.(PathViewerModel)

	if m.SelectedLevel != 1 {
		t.Errorf("expected SelectedLevel=1 after KeyUp, got %d", m.SelectedLevel)
	}

	msg := extractMsg(cmd)
	focus, ok := msg.(FocusNodeMsg)
	if !ok {
		t.Fatalf("expected FocusNodeMsg, got %T", msg)
	}
	if focus.Level != 1 {
		t.Errorf("expected FocusNodeMsg.Level=1, got %d", focus.Level)
	}

	// Move up again
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = result.(PathViewerModel)
	if m.SelectedLevel != 0 {
		t.Errorf("expected SelectedLevel=0, got %d", m.SelectedLevel)
	}

	// At top, should stay at 0
	result, cmd = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = result.(PathViewerModel)
	if m.SelectedLevel != 0 {
		t.Errorf("expected SelectedLevel=0 (clamped), got %d", m.SelectedLevel)
	}
	if cmd != nil {
		t.Error("expected nil cmd at top boundary")
	}
}

func TestPathViewer_KeyDown(t *testing.T) {
	m := newInitializedPathViewer()
	// Start at level 0
	m.SelectedLevel = 0

	result, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = result.(PathViewerModel)

	if m.SelectedLevel != 1 {
		t.Errorf("expected SelectedLevel=1 after KeyDown, got %d", m.SelectedLevel)
	}

	msg := extractMsg(cmd)
	focus, ok := msg.(FocusNodeMsg)
	if !ok {
		t.Fatalf("expected FocusNodeMsg, got %T", msg)
	}
	if focus.Level != 1 {
		t.Errorf("expected FocusNodeMsg.Level=1, got %d", focus.Level)
	}

	// Move to end
	result, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = result.(PathViewerModel)
	if m.SelectedLevel != 2 {
		t.Errorf("expected SelectedLevel=2, got %d", m.SelectedLevel)
	}

	// At bottom, should stay
	result, cmd = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = result.(PathViewerModel)
	if m.SelectedLevel != 2 {
		t.Errorf("expected SelectedLevel=2 (clamped), got %d", m.SelectedLevel)
	}
	if cmd != nil {
		t.Error("expected nil cmd at bottom boundary")
	}
}

func TestPathViewer_OpenDropdown(t *testing.T) {
	m := newInitializedPathViewer()
	// Navigate to level 1 (B), which has alternatives (B, C)
	m.SelectedLevel = 1

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = result.(PathViewerModel)

	if !m.DropdownOpen {
		t.Error("expected DropdownOpen=true after Space on node with alternatives")
	}
}

func TestPathViewer_OpenDropdown_NoAlternatives(t *testing.T) {
	m := newInitializedPathViewer()
	// Navigate to level 0 (root A), which has no alternatives
	m.SelectedLevel = 0

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = result.(PathViewerModel)

	if m.DropdownOpen {
		t.Error("expected DropdownOpen=false when no alternatives exist")
	}
}

func TestPathViewer_EnterOnLeaf(t *testing.T) {
	m := newInitializedPathViewer()
	// SelectedLevel=2 (leaf D)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	msg := extractMsg(cmd)
	sel, ok := msg.(SelectNodeMsg)
	if !ok {
		t.Fatalf("expected SelectNodeMsg, got %T", msg)
	}
	if sel.Tree.Node.DisplayName() != "D" {
		t.Errorf("expected selected node 'D', got %q", sel.Tree.Node.DisplayName())
	}
}

func TestPathViewer_EnterOnNonLeaf(t *testing.T) {
	m := newInitializedPathViewer()
	// Navigate to level 0 (root A, not a leaf)
	m.SelectedLevel = 0

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd != nil {
		t.Error("expected nil cmd for Enter on non-leaf node")
	}
}

func TestPathViewer_DropdownSelection(t *testing.T) {
	m := newInitializedPathViewer()
	// Select level 1 (B) and open dropdown
	m.SelectedLevel = 1

	// Open dropdown first
	result, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = result.(PathViewerModel)
	if !m.DropdownOpen {
		t.Fatal("expected dropdown to be open")
	}

	// Send OptionSelectedMsg for second alternative (C, index=1)
	result, cmd := m.Update(teadd.OptionSelectedMsg{
		Index: 1,
		Text:  "C",
	})
	m = result.(PathViewerModel)

	if m.DropdownOpen {
		t.Error("expected dropdown closed after selection")
	}

	// Path should be rebuilt: A → C → E
	if len(m.Path) != 3 {
		t.Fatalf("expected path length 3 after rebuild, got %d", len(m.Path))
	}
	if m.Path[1].Node.DisplayName() != "C" {
		t.Errorf("expected path[1]='C' after selection, got %q", m.Path[1].Node.DisplayName())
	}
	if m.Path[2].Node.DisplayName() != "E" {
		t.Errorf("expected path[2]='E' after rebuild, got %q", m.Path[2].Node.DisplayName())
	}

	msg := extractMsg(cmd)
	change, ok := msg.(ChangeNodeMsg)
	if !ok {
		t.Fatalf("expected ChangeNodeMsg, got %T", msg)
	}
	if change.Level != 1 {
		t.Errorf("expected ChangeNodeMsg.Level=1, got %d", change.Level)
	}
}

func TestPathViewer_DropdownCancellation(t *testing.T) {
	m := newInitializedPathViewer()
	m.SelectedLevel = 1

	// Open dropdown
	result, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = result.(PathViewerModel)
	if !m.DropdownOpen {
		t.Fatal("expected dropdown open")
	}

	// Cancel dropdown
	result, _ = m.Update(teadd.DropdownCancelledMsg{})
	m = result.(PathViewerModel)

	if m.DropdownOpen {
		t.Error("expected dropdown closed after cancellation")
	}
	// Path should be unchanged
	if m.Path[1].Node.DisplayName() != "B" {
		t.Errorf("expected path[1] still 'B', got %q", m.Path[1].Node.DisplayName())
	}
}

func TestPathViewer_WindowSizeMsg(t *testing.T) {
	m := newInitializedPathViewer()

	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = result.(PathViewerModel)

	if m.Width != 120 {
		t.Errorf("expected Width=120, got %d", m.Width)
	}
	if m.Height != 40 {
		t.Errorf("expected Height=40, got %d", m.Height)
	}
}

// --- Layer 2: View Tests ---

func TestPathViewer_View(t *testing.T) {
	m := newInitializedPathViewer()
	view := m.View()

	if view == "" {
		t.Fatal("expected non-empty view")
	}
	// Path items should be rendered
	if !strings.Contains(view, "A") {
		t.Error("expected view to contain 'A'")
	}
	if !strings.Contains(view, "B") {
		t.Error("expected view to contain 'B'")
	}
	if !strings.Contains(view, "D") {
		t.Error("expected view to contain 'D'")
	}
	// Prompt should be visible
	if !strings.Contains(view, "Dependencies") {
		t.Error("expected view to contain prompt 'Dependencies'")
	}
}

func TestPathViewer_View_DropdownOpen(t *testing.T) {
	m := newInitializedPathViewer()
	m.SelectedLevel = 1

	// Open dropdown
	result, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = result.(PathViewerModel)

	view := m.View()
	// When dropdown is open, it should overlay the view
	// The dropdown items (alternatives B, C) should be visible
	if !strings.Contains(view, "B") {
		t.Error("expected dropdown to show 'B'")
	}
	if !strings.Contains(view, "C") {
		t.Error("expected dropdown to show 'C'")
	}
}

func TestPathViewer_View_BorderGeometry(t *testing.T) {
	m := newInitializedPathViewer()
	m.Width = 40
	m.Height = 10
	view := m.View()

	lines := strings.Split(view, "\n")
	// Should not exceed Height
	if len(lines) > m.Height {
		t.Errorf("expected at most %d lines, got %d", m.Height, len(lines))
	}
}
