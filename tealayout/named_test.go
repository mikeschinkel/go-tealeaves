package tealayout

import "testing"

func TestLayout_PaneByName(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree")
	code := NewColumn(Flex(1)).WithName("code")
	root := NewRow(Percent100, tree, code)
	layout := NewLayout(root)

	if p := layout.Pane("tree"); p != tree {
		t.Error("Pane(\"tree\") should return tree pane")
	}
	if p := layout.Pane("code"); p != code {
		t.Error("Pane(\"code\") should return code pane")
	}
}

func TestLayout_PaneNotFound(t *testing.T) {
	root := NewRow(Percent100, NewColumn(Flex(1)).WithName("a"))
	layout := NewLayout(root)

	if p := layout.Pane("nonexistent"); p != nil {
		t.Errorf("Pane(\"nonexistent\") = %v, want nil", p)
	}
}

func TestLayout_PaneNested(t *testing.T) {
	inner := NewColumn(Flex(1)).WithName("inner")
	middle := NewRow(Flex(1), inner).WithName("middle")
	root := NewColumn(Percent100, middle)
	layout := NewLayout(root)

	if p := layout.Pane("inner"); p != inner {
		t.Error("nested pane lookup failed")
	}
	if p := layout.Pane("middle"); p != middle {
		t.Error("middle pane lookup failed")
	}
}
