package tealayout

import "testing"

func TestComponent_NewRow(t *testing.T) {
	comp := NewRow(Percent50)
	if comp.direction != Horizontal {
		t.Errorf("direction = %v, want Horizontal", comp.direction)
	}
	if comp.dim != Percent50 {
		t.Errorf("dim = %+v, want Percent50", comp.dim)
	}
	if comp.maxSize != -1 {
		t.Errorf("maxSize = %d, want -1", comp.maxSize)
	}
}

func TestComponent_NewColumn(t *testing.T) {
	comp := NewColumn(Percent100)
	if comp.direction != Vertical {
		t.Errorf("direction = %v, want Vertical", comp.direction)
	}
	if comp.dim != Percent100 {
		t.Errorf("dim = %+v, want Percent100", comp.dim)
	}
}

func TestComponent_WithMinSize(t *testing.T) {
	comp := NewRow(Flex(1)).WithMinSize(20)
	if comp.minSize != 20 {
		t.Errorf("minSize = %d, want 20", comp.minSize)
	}
}

func TestComponent_WithMaxSize(t *testing.T) {
	comp := NewRow(Flex(1)).WithMaxSize(60)
	if comp.maxSize != 60 {
		t.Errorf("maxSize = %d, want 60", comp.maxSize)
	}
}

func TestComponent_WithOptional(t *testing.T) {
	comp := NewRow(Flex(1)).WithOptional(true)
	if !comp.optional {
		t.Error("optional = false, want true")
	}
}

func TestComponent_WithGap(t *testing.T) {
	comp := NewRow(Flex(1)).WithGap(2)
	if comp.gap != 2 {
		t.Errorf("gap = %d, want 2", comp.gap)
	}
}

func TestComponent_Chaining(t *testing.T) {
	comp := NewColumn(Percent25).WithMinSize(20).WithMaxSize(60).WithOptional(true).WithGap(1)
	if comp.minSize != 20 {
		t.Errorf("minSize = %d, want 20", comp.minSize)
	}
	if comp.maxSize != 60 {
		t.Errorf("maxSize = %d, want 60", comp.maxSize)
	}
	if !comp.optional {
		t.Error("optional = false, want true")
	}
	if comp.gap != 1 {
		t.Errorf("gap = %d, want 1", comp.gap)
	}
}

func TestComponent_ToConstraint_WithModifiers(t *testing.T) {
	comp := NewColumn(Percent25).WithMinSize(10).WithMaxSize(60).WithOptional(true)
	cs := comp.toConstraint()

	if cs.kind != constraintFlex {
		t.Errorf("kind = %v, want constraintFlex", cs.kind)
	}
	if cs.flexWeight != 25 {
		t.Errorf("flexWeight = %f, want 25", cs.flexWeight)
	}
	if cs.minSize != 10 {
		t.Errorf("minSize = %d, want 10", cs.minSize)
	}
	if cs.maxSize != 60 {
		t.Errorf("maxSize = %d, want 60", cs.maxSize)
	}
	if !cs.optional {
		t.Error("optional = false, want true")
	}
}

func TestComponent_NestedResolution(t *testing.T) {
	// Build a nested tree and verify resolution matches expected sizes.
	// Outer column (100): row(fixed 3) + row(flex 1) + row(fixed 1)
	w1 := &mockWidget{char: 'H'}
	w2 := &mockWidget{char: 'C'}
	w3 := &mockWidget{char: 'F'}

	root := NewColumn(Percent100,
		NewRow(Fixed(3), w1),
		NewRow(Flex(1), w2),
		NewRow(Fixed(1), w3),
	)
	root.SetSize(80, 24)

	_, err := root.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w1.width != 80 || w1.height != 3 {
		t.Errorf("w1 = %dx%d, want 80x3", w1.width, w1.height)
	}
	if w2.width != 80 || w2.height != 20 {
		t.Errorf("w2 = %dx%d, want 80x20", w2.width, w2.height)
	}
	if w3.width != 80 || w3.height != 1 {
		t.Errorf("w3 = %dx%d, want 80x1", w3.width, w3.height)
	}
}

func TestComponent_PercentDistribution(t *testing.T) {
	// 50/25/25 split should distribute like flex weights 50/25/25
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}
	w3 := &mockWidget{char: 'C'}

	root := NewRow(Percent100,
		NewColumn(Percent50, w1),
		NewColumn(Percent25, w2),
		NewColumn(Percent25, w3),
	)
	root.SetSize(100, 10)

	_, err := root.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w1.width != 50 {
		t.Errorf("w1.width = %d, want 50", w1.width)
	}
	if w2.width != 25 {
		t.Errorf("w2.width = %d, want 25", w2.width)
	}
	if w3.width != 25 {
		t.Errorf("w3.width = %d, want 25", w3.width)
	}
}

func TestComponent_OptionalRemoval(t *testing.T) {
	// Third column is optional with min 30 — in 50 total can't fit.
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}
	w3 := &mockWidget{char: 'C'}

	root := NewRow(Percent100,
		NewColumn(Flex(1), w1),
		NewColumn(Flex(1), w2),
		NewColumn(Flex(1), w3).WithMinSize(30).WithOptional(true),
	)
	root.SetSize(50, 10)

	_, err := root.Render()
	if err != nil {
		t.Fatal(err)
	}

	// w3 should have been removed (size 0), w1+w2 split the space
	if w1.width+w2.width != 50 {
		t.Errorf("w1+w2 = %d, want 50", w1.width+w2.width)
	}
	if w3.width != 0 && w3.height != 0 {
		t.Errorf("w3 should not have been rendered, got %dx%d", w3.width, w3.height)
	}
}

func TestComponent_NilElementsSkipped(t *testing.T) {
	w := &mockWidget{char: 'X'}
	root := NewRow(Percent100, nil, NewColumn(Flex(1), w), nil)
	root.SetSize(80, 10)

	_, err := root.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w.width != 80 {
		t.Errorf("w.width = %d, want 80", w.width)
	}
}
