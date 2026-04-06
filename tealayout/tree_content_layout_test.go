package tealayout

import (
	"strings"
	"testing"
)

// mockTreeWidget implements ContentProvider, SetSizer, and SizeHinter.
type mockTreeWidget struct {
	width  int
	height int
	hintW  int
}

func (m *mockTreeWidget) SetSize(w, h int) { m.width = w; m.height = h }
func (m *mockTreeWidget) Content() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	line := strings.Repeat("T", m.width)
	lines := make([]string, m.height)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}
func (m *mockTreeWidget) SizeHint(availW, availH int) SizeHint {
	return SizeHint{
		Desired: Size{Width: m.hintW, Height: availH},
		Max:     Size{Width: -1, Height: -1},
	}
}

// mockContentWidget implements ContentProvider and SetSizer.
type mockContentWidget struct {
	width  int
	height int
}

func (m *mockContentWidget) SetSize(w, h int) { m.width = w; m.height = h }
func (m *mockContentWidget) Content() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	line := strings.Repeat("C", m.width)
	lines := make([]string, m.height)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func TestTreeContentLayout_BasicRender(t *testing.T) {
	tree := &mockTreeWidget{hintW: 30}
	content := &mockContentWidget{}

	tcl := NewTreeContentLayout(tree, content)
	tcl.SetSize(100, 20)

	output, err := tcl.Render()
	if err != nil {
		t.Fatal(err)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}

	// Tree should get its desired width
	if tree.width != 30 {
		t.Errorf("tree.width = %d, want 30", tree.width)
	}
	// Content should get remaining
	if content.width != 70 {
		t.Errorf("content.width = %d, want 70", content.width)
	}
	// Both get full height
	if tree.height != 20 {
		t.Errorf("tree.height = %d, want 20", tree.height)
	}
	if content.height != 20 {
		t.Errorf("content.height = %d, want 20", content.height)
	}
}

func TestTreeContentLayout_TypedAccess(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}

	tcl := NewTreeContentLayout(tree, content)

	if tcl.Tree() != tree {
		t.Error("Tree() should return the tree widget")
	}
	if tcl.Content() != content {
		t.Error("Content() should return the content widget")
	}
}

func TestTreeContentLayout_Focus(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}

	tcl := NewTreeContentLayout(tree, content)

	// Initial focus is tree (first in depth-first order)
	if !tcl.TreeFocused() {
		t.Error("initial focus should be tree")
	}
	if tcl.ContentFocused() {
		t.Error("content should not be focused initially")
	}

	tcl.FocusContent()
	if !tcl.ContentFocused() {
		t.Error("after FocusContent, content should be focused")
	}

	tcl.FocusTree()
	if !tcl.TreeFocused() {
		t.Error("after FocusTree, tree should be focused")
	}

	tcl.ToggleFocus()
	if !tcl.ContentFocused() {
		t.Error("after ToggleFocus from tree, content should be focused")
	}

	tcl.ToggleFocus()
	if !tcl.TreeFocused() {
		t.Error("after ToggleFocus from content, tree should be focused")
	}
}

func TestTreeContentLayout_Visibility(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}

	tcl := NewTreeContentLayout(tree, content)
	tcl.SetSize(100, 20)

	tcl.HideTree()
	pane := tcl.PaneLayout().Pane("tree")
	if pane.Visible() {
		t.Error("tree pane should be hidden")
	}

	tcl.ShowTree()
	if !pane.Visible() {
		t.Error("tree pane should be visible")
	}

	tcl.HideContent()
	cpane := tcl.PaneLayout().Pane("content")
	if cpane.Visible() {
		t.Error("content pane should be hidden")
	}

	tcl.ShowContent()
	if !cpane.Visible() {
		t.Error("content pane should be visible")
	}
}

func TestTreeContentLayout_ToggleTree(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	tcl.ToggleTree()
	if tcl.PaneLayout().Pane("tree").Visible() {
		t.Error("tree should be hidden after ToggleTree")
	}

	tcl.ToggleTree()
	if !tcl.PaneLayout().Pane("tree").Visible() {
		t.Error("tree should be visible after second ToggleTree")
	}
}

func TestTreeContentLayout_ToggleTree_LastVisible(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	tcl.HideContent()
	tcl.ToggleTree() // should be no-op — tree is last visible
	if !tcl.PaneLayout().Pane("tree").Visible() {
		t.Error("toggle should be no-op when tree is last visible pane")
	}
}

func TestTreeContentLayout_ToggleContent(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	tcl.ToggleContent()
	if tcl.PaneLayout().Pane("content").Visible() {
		t.Error("content should be hidden after ToggleContent")
	}

	tcl.ToggleContent()
	if !tcl.PaneLayout().Pane("content").Visible() {
		t.Error("content should be visible after second ToggleContent")
	}
}

func TestTreeContentLayout_ToggleContent_LastVisible(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	tcl.HideTree()
	tcl.ToggleContent() // should be no-op
	if !tcl.PaneLayout().Pane("content").Visible() {
		t.Error("toggle should be no-op when content is last visible pane")
	}
}

func TestTreeContentLayout_SoloTree(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	tcl.FocusContent()
	tcl.SoloTree()
	if !tcl.PaneLayout().Pane("tree").Visible() {
		t.Error("tree should be visible")
	}
	if tcl.PaneLayout().Pane("content").Visible() {
		t.Error("content should be hidden")
	}
	if !tcl.TreeFocused() {
		t.Error("tree should be focused")
	}
}

func TestTreeContentLayout_SoloContent(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	tcl.SoloContent()
	if tcl.PaneLayout().Pane("tree").Visible() {
		t.Error("tree should be hidden")
	}
	if !tcl.PaneLayout().Pane("content").Visible() {
		t.Error("content should be visible")
	}
	if !tcl.ContentFocused() {
		t.Error("content should be focused")
	}
}

func TestTreeContentLayout_ShowBoth(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	tcl.SoloTree()
	tcl.ShowBoth()
	if !tcl.PaneLayout().Pane("tree").Visible() {
		t.Error("tree should be visible")
	}
	if !tcl.PaneLayout().Pane("content").Visible() {
		t.Error("content should be visible")
	}
}

func TestTreeContentLayout_ToggleTree_FocusFollows(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}
	tcl := NewTreeContentLayout(tree, content)

	// Focus tree, then toggle it off — focus should move to content.
	tcl.FocusTree()
	tcl.ToggleTree()
	if !tcl.ContentFocused() {
		t.Error("focus should move to content when tree is hidden")
	}
}

func TestTreeContentLayout_PaneLayout(t *testing.T) {
	tree := &mockTreeWidget{hintW: 25}
	content := &mockContentWidget{}

	tcl := NewTreeContentLayout(tree, content)
	if tcl.PaneLayout() == nil {
		t.Error("PaneLayout() should not return nil")
	}
}

// contentViewer is an interface type to test TreeContentLayout with interface C.
type contentViewer interface {
	ContentProvider
	SetSizer
}

func TestTreeContentLayout_InterfaceType(t *testing.T) {
	tree := &mockTreeWidget{hintW: 20}
	var content contentViewer = &mockContentWidget{}

	tcl := NewTreeContentLayout[*mockTreeWidget, contentViewer](tree, content)
	tcl.SetSize(80, 20)

	_, err := tcl.Render()
	if err != nil {
		t.Fatal(err)
	}

	if tree.width != 20 {
		t.Errorf("tree.width = %d, want 20", tree.width)
	}
}
