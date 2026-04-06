package tealayout

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teacrumbs"
)

// mockStackView is a minimal StackView for testing.
type mockStackView struct {
	name      string
	width     int
	height    int
	entered   int
	exited    int
	initCalls int
}

func (m *mockStackView) Init() tea.Cmd {
	m.initCalls++
	return nil
}

func (m *mockStackView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *mockStackView) View() tea.View {
	return tea.NewView(m.name + " view")
}

func (m *mockStackView) OnEnter() tea.Cmd {
	m.entered++
	return nil
}

func (m *mockStackView) OnExit() tea.Cmd {
	m.exited++
	return nil
}

func (m *mockStackView) Breadcrumb() teacrumbs.Crumb {
	return teacrumbs.NewCrumb(m.name, nil)
}

func (m *mockStackView) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func newMockView(name string) *mockStackView {
	return &mockStackView{name: name}
}

func TestStackLayoutModel_NewAndCurrent(t *testing.T) {
	view := newMockView("root")
	m := NewStackLayoutModel(view, teacrumbs.DefaultStyles())

	if m.Current() != view {
		t.Error("Current() should return initial view")
	}
	if m.Depth() != 1 {
		t.Errorf("Depth() = %d, want 1", m.Depth())
	}
	if m.CanPop() {
		t.Error("CanPop() should be false with single view")
	}
}

func TestStackLayoutModel_PushPop(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	child := newMockView("child")
	m.Push(child, "")

	if m.Current() != child {
		t.Error("Current() should be child after push")
	}
	if m.Depth() != 2 {
		t.Errorf("Depth() = %d, want 2", m.Depth())
	}
	if !m.CanPop() {
		t.Error("CanPop() should be true with 2 views")
	}

	// Verify lifecycle: root.OnExit called, child.OnEnter called
	if root.exited != 1 {
		t.Errorf("root.exited = %d, want 1", root.exited)
	}
	if child.entered != 1 {
		t.Errorf("child.entered = %d, want 1", child.entered)
	}

	_, err := m.Pop()
	if err != nil {
		t.Fatal(err)
	}

	if m.Current() != root {
		t.Error("Current() should be root after pop")
	}
	if m.Depth() != 1 {
		t.Errorf("Depth() = %d, want 1", m.Depth())
	}

	// child.OnExit and root.OnEnter called
	if child.exited != 1 {
		t.Errorf("child.exited = %d, want 1", child.exited)
	}
	if root.entered != 1 {
		t.Errorf("root.entered = %d, want 1", root.entered)
	}
}

func TestStackLayoutModel_PopEmpty(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	_, err := m.Pop()
	if err == nil {
		t.Fatal("expected error popping single view")
	}
	if !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected ErrStackUnderflow, got %v", err)
	}
}

func TestStackLayoutModel_PopTo(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	v1 := newMockView("v1")
	v2 := newMockView("v2")
	v3 := newMockView("v3")
	m.Push(v1, "")
	m.Push(v2, "")
	m.Push(v3, "")

	if m.Depth() != 4 {
		t.Fatalf("Depth() = %d, want 4", m.Depth())
	}

	_, err := m.PopTo(2)
	if err != nil {
		t.Fatal(err)
	}
	if m.Depth() != 2 {
		t.Errorf("Depth() = %d, want 2", m.Depth())
	}
	if m.Current() != v1 {
		t.Error("Current() should be v1 after PopTo(2)")
	}
}

func TestStackLayoutModel_PopToError(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	_, err := m.PopTo(0)
	if err == nil {
		t.Fatal("expected error for PopTo(0)")
	}
	if !errors.Is(err, ErrStackUnderflow) {
		t.Errorf("expected ErrStackUnderflow, got %v", err)
	}

	_, err = m.PopTo(5)
	if err == nil {
		t.Fatal("expected error for PopTo(5) with depth 1")
	}
}

func TestStackLayoutModel_Cache(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	child := newMockView("child")
	m.Push(child, "child-key")

	cached, ok := m.GetCached("child-key")
	if !ok {
		t.Fatal("expected cached view")
	}
	if cached != child {
		t.Error("cached view should be the same instance")
	}

	m.DeleteCached("child-key")
	_, ok = m.GetCached("child-key")
	if ok {
		t.Error("expected cache miss after delete")
	}
}

func TestStackLayoutModel_SetSize(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())
	m.width = 80
	m.height = 24
	m.Push(root, "") // re-push to trigger SetSize

	// After push, view gets height minus breadcrumbHeight
	if root.width != 80 {
		t.Errorf("root.width = %d, want 80", root.width)
	}
	if root.height != 23 {
		t.Errorf("root.height = %d, want 23 (24 - 1 breadcrumb)", root.height)
	}
}

func TestStackLayoutModel_UpdateDelegation(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	// WindowSizeMsg should set dimensions
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = updated.(StackLayoutModel)
	if m.width != 100 || m.height != 30 {
		t.Errorf("size = %dx%d, want 100x30", m.width, m.height)
	}
}

func TestStackLayoutModel_UpdateCurrent(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	replacement := newMockView("replacement")
	m.UpdateCurrent(replacement)

	if m.Current() != replacement {
		t.Error("Current() should be replacement after UpdateCurrent")
	}
}

func TestStackLayoutModel_UpdateCurrentWithCache(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	child := newMockView("child")
	m.Push(child, "child-key")

	replacement := newMockView("replacement")
	m.UpdateCurrent(replacement)

	// Cache should be updated too
	cached, ok := m.GetCached("child-key")
	if !ok {
		t.Fatal("expected cached view")
	}
	if cached != replacement {
		t.Error("cache should be updated after UpdateCurrent")
	}
}

func TestStackLayoutModel_View(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())
	m.width = 80
	m.height = 24
	root.SetSize(80, 23)

	v := m.View()
	if v.Content == "" {
		t.Error("View should return non-empty content")
	}
}

func TestStackLayoutModel_PushViewMsg(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	child := newMockView("child")
	updated, _ := m.Update(PushViewMsg{View: child, CacheKey: "test"})
	m = updated.(StackLayoutModel)

	if m.Depth() != 2 {
		t.Errorf("Depth() = %d, want 2", m.Depth())
	}
	if m.Current() != child {
		t.Error("Current() should be child after PushViewMsg")
	}
}

func TestStackLayoutModel_PopViewMsg(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	child := newMockView("child")
	m.Push(child, "")

	updated, _ := m.Update(PopViewMsg{})
	m = updated.(StackLayoutModel)

	if m.Depth() != 1 {
		t.Errorf("Depth() = %d, want 1", m.Depth())
	}
}

func TestStackLayoutModel_Init(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	m.Init()
	if root.initCalls != 1 {
		t.Errorf("root.initCalls = %d, want 1", root.initCalls)
	}
	if root.entered != 1 {
		t.Errorf("root.entered = %d, want 1 (OnEnter called in Init)", root.entered)
	}
}

func TestStackLayoutModel_Breadcrumbs(t *testing.T) {
	root := newMockView("root")
	m := NewStackLayoutModel(root, teacrumbs.DefaultStyles())

	child := newMockView("child")
	m.Push(child, "")

	// Breadcrumbs should have 2 entries
	crumbs := m.crumbs.Crumbs()
	if len(crumbs) != 2 {
		t.Errorf("crumbs count = %d, want 2", len(crumbs))
	}
	if crumbs[0].Text != "root" {
		t.Errorf("crumbs[0].Text = %q, want %q", crumbs[0].Text, "root")
	}
	if crumbs[1].Text != "child" {
		t.Errorf("crumbs[1].Text = %q, want %q", crumbs[1].Text, "child")
	}

	if _, err := m.Pop(); err != nil {
		t.Fatalf("Pop: %v", err)
	}
	crumbs = m.crumbs.Crumbs()
	if len(crumbs) != 1 {
		t.Errorf("crumbs count = %d, want 1 after pop", len(crumbs))
	}
}
