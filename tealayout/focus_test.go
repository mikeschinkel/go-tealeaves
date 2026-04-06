package tealayout

import (
	"errors"
	"testing"
)

func TestFocusManager_InitialFocus(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	if fm.FocusedPane() != tree {
		t.Error("initial focus should be first pane (tree)")
	}
	if !tree.Focused() {
		t.Error("tree should be focused")
	}
}

func TestFocusManager_Next(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	diff := NewColumn(Flex(1)).WithName("diff").WithFocusable()
	root := NewRow(Percent100, tree, code, diff)
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	fm.FocusNext()
	if fm.FocusedPane() != code {
		t.Error("after FocusNext, should be code")
	}
	if tree.Focused() {
		t.Error("tree should be blurred")
	}
	if !code.Focused() {
		t.Error("code should be focused")
	}

	fm.FocusNext()
	if fm.FocusedPane() != diff {
		t.Error("after second FocusNext, should be diff")
	}

	// Wrap around
	fm.FocusNext()
	if fm.FocusedPane() != tree {
		t.Error("should wrap to tree")
	}
}

func TestFocusManager_Prev(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	fm.FocusPrev()
	if fm.FocusedPane() != code {
		t.Error("FocusPrev from first should wrap to last (code)")
	}
}

func TestFocusManager_FocusByName(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	err := fm.FocusPane("code")
	if err != nil {
		t.Fatal(err)
	}
	if fm.FocusedPane() != code {
		t.Error("should be focused on code")
	}
}

func TestFocusManager_FocusByName_NotFound(t *testing.T) {
	root := NewRow(Percent100, NewColumn(Flex(1)).WithName("a").WithFocusable())
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	err := fm.FocusPane("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrPaneNotFound) {
		t.Errorf("expected ErrPaneNotFound, got %v", err)
	}
}

func TestFocusManager_SkipsHiddenPanes(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	diff := NewColumn(Flex(1)).WithName("diff").WithFocusable()
	root := NewRow(Percent100, tree, code, diff)
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	// Hide code, then FocusNext should skip to diff
	code.Hide()
	fm.FocusNext()
	if fm.FocusedPane() != diff {
		t.Errorf("should skip hidden code, got %q", fm.FocusedPane().Name())
	}
}

func TestFocusManager_EnsureFocusedVisible(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	// Focus is on tree, hide it
	tree.Hide()
	fm.EnsureFocusedVisible()
	if fm.FocusedPane() != code {
		t.Error("should auto-advance to code when tree is hidden")
	}
}

func TestFocusManager_Focused_Convenience(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	if !fm.Focused("tree") {
		t.Error("Focused(\"tree\") should be true")
	}
	if fm.Focused("code") {
		t.Error("Focused(\"code\") should be false")
	}
}

func TestFocusManager_NamedNotFocusable(t *testing.T) {
	header := NewRow(Fixed(1)).WithName("header") // named but NOT focusable
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewColumn(Percent100, header, NewRow(Flex(1), tree, code))
	layout := NewLayout(root)
	fm := NewFocusManager(layout)

	// Initial focus should be tree, not header
	if fm.FocusedPane() != tree {
		t.Errorf("initial focus should be tree, got %q", fm.FocusedPane().Name())
	}

	// Tab should go to code, not header
	fm.FocusNext()
	if fm.FocusedPane() != code {
		t.Errorf("FocusNext should go to code, got %q", fm.FocusedPane().Name())
	}

	// Tab wraps back to tree
	fm.FocusNext()
	if fm.FocusedPane() != tree {
		t.Errorf("FocusNext should wrap to tree, got %q", fm.FocusedPane().Name())
	}
}
