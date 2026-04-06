package tealayout

import "testing"

func TestPane_DefaultVisible(t *testing.T) {
	p := NewRow(Flex(1))
	if !p.Visible() {
		t.Error("new pane should be visible by default")
	}
}

func TestPane_HideShow(t *testing.T) {
	p := NewRow(Flex(1))
	p.Hide()
	if p.Visible() {
		t.Error("Hide() should make pane invisible")
	}
	p.Show()
	if !p.Visible() {
		t.Error("Show() should make pane visible")
	}
}

func TestPane_HiddenPaneGetsZeroSize(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}

	col1 := NewColumn(Flex(1), NewElement(w1))
	col2 := NewColumn(Flex(1), NewElement(w2))
	col2.Hide()

	root := NewRow(Percent100, col1, col2)
	root.SetSize(80, 10)

	sizes, err := root.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	// col2 hidden → gets 0, col1 gets all 80
	if sizes[0] != 80 {
		t.Errorf("sizes[0] = %d, want 80", sizes[0])
	}
	if sizes[1] != 0 {
		t.Errorf("sizes[1] = %d, want 0", sizes[1])
	}
}

func TestPane_ShowRestoresSize(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}

	col1 := NewColumn(Flex(1), NewElement(w1))
	col2 := NewColumn(Flex(1), NewElement(w2))

	root := NewRow(Percent100, col1, col2)
	root.SetSize(80, 10)

	// Initially both visible → 40/40
	sizes, _ := root.Resolve()
	if sizes[0] != 40 || sizes[1] != 40 {
		t.Errorf("initial: %v, want [40 40]", sizes)
	}

	// Hide col2
	col2.Hide()
	sizes, _ = root.Resolve()
	if sizes[0] != 80 || sizes[1] != 0 {
		t.Errorf("after hide: %v, want [80 0]", sizes)
	}

	// Show col2
	col2.Show()
	sizes, _ = root.Resolve()
	if sizes[0] != 40 || sizes[1] != 40 {
		t.Errorf("after show: %v, want [40 40]", sizes)
	}
}

func TestPane_HidePropagatesDirtyUp(t *testing.T) {
	child := NewColumn(Flex(1))
	root := NewRow(Percent100, child)
	root.SetSize(80, 10)
	if _, err := root.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	if !root.resolved {
		t.Fatal("root should be resolved")
	}

	child.Hide()
	if root.resolved {
		t.Error("hiding child should mark root as dirty (not resolved)")
	}
}

func TestPane_HiddenPaneNotRendered(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}

	col1 := NewColumn(Flex(1), NewElement(w1))
	col2 := NewColumn(Flex(1), NewElement(w2))
	col2.Hide()

	root := NewRow(Percent100, col1, col2)
	root.SetSize(80, 10)

	_, err := root.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w1.width != 80 {
		t.Errorf("w1 width = %d, want 80", w1.width)
	}
	// w2 should not have been rendered (SetSize not called with positive values)
	if w2.width > 0 {
		t.Errorf("hidden widget w2 should not have been sized, got width=%d", w2.width)
	}
}

func TestPane_WithName(t *testing.T) {
	p := NewRow(Flex(1)).WithName("sidebar")
	if p.Name() != "sidebar" {
		t.Errorf("Name() = %q, want %q", p.Name(), "sidebar")
	}
}

func TestPane_ParentPointers(t *testing.T) {
	child := NewColumn(Flex(1))
	root := NewRow(Percent100, child)

	if child.parent != root {
		t.Error("child.parent should point to root")
	}
}

func TestPane_HiddenGapReclaimed(t *testing.T) {
	col1 := NewColumn(Flex(1))
	col2 := NewColumn(Flex(1))
	col3 := NewColumn(Flex(1))
	col3.Hide()

	root := NewRow(Percent100, col1, col2, col3).WithGap(2)
	root.SetSize(80, 10)

	sizes, _ := root.Resolve()
	// 3 children, 1 hidden: gap counted only between 2 visible → 80-2=78, 78/2=39
	total := sizes[0] + sizes[1]
	if total != 78 {
		t.Errorf("sum of visible = %d, want 78 (80 - 2 gap)", total)
	}
	if sizes[2] != 0 {
		t.Errorf("hidden pane size = %d, want 0", sizes[2])
	}
}
