package tealayout

import "testing"

func TestPane_NewRow(t *testing.T) {
	p := NewRow(Percent50)
	if p.direction != Horizontal {
		t.Errorf("direction = %v, want Horizontal", p.direction)
	}
	if p.dim != Percent50 {
		t.Errorf("dim = %+v, want Percent50", p.dim)
	}
	if p.maxSize != -1 {
		t.Errorf("maxSize = %d, want -1", p.maxSize)
	}
}

func TestPane_NewColumn(t *testing.T) {
	p := NewColumn(Percent100)
	if p.direction != Vertical {
		t.Errorf("direction = %v, want Vertical", p.direction)
	}
	if p.dim != Percent100 {
		t.Errorf("dim = %+v, want Percent100", p.dim)
	}
}

func TestPane_WithMinSize(t *testing.T) {
	p := NewRow(Flex(1)).WithMinSize(20)
	if p.minSize != 20 {
		t.Errorf("minSize = %d, want 20", p.minSize)
	}
}

func TestPane_WithMaxSize(t *testing.T) {
	p := NewRow(Flex(1)).WithMaxSize(60)
	if p.maxSize != 60 {
		t.Errorf("maxSize = %d, want 60", p.maxSize)
	}
}

func TestPane_WithOptional(t *testing.T) {
	p := NewRow(Flex(1)).WithOptional(true)
	if !p.optional {
		t.Error("optional = false, want true")
	}
}

func TestPane_WithGap(t *testing.T) {
	p := NewRow(Flex(1)).WithGap(2)
	if p.gap != 2 {
		t.Errorf("gap = %d, want 2", p.gap)
	}
}

func TestPane_Chaining(t *testing.T) {
	p := NewColumn(Percent25).WithMinSize(20).WithMaxSize(60).WithOptional(true).WithGap(1)
	if p.minSize != 20 {
		t.Errorf("minSize = %d, want 20", p.minSize)
	}
	if p.maxSize != 60 {
		t.Errorf("maxSize = %d, want 60", p.maxSize)
	}
	if !p.optional {
		t.Error("optional = false, want true")
	}
	if p.gap != 1 {
		t.Errorf("gap = %d, want 1", p.gap)
	}
}

func TestPane_ToConstraint_WithModifiers(t *testing.T) {
	p := NewColumn(Percent25).WithMinSize(10).WithMaxSize(60).WithOptional(true)
	cs := p.toConstraint()

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

func TestPane_NestedResolution(t *testing.T) {
	w1 := &mockWidget{char: 'H'}
	w2 := &mockWidget{char: 'C'}
	w3 := &mockWidget{char: 'F'}

	root := NewColumn(Percent100,
		NewRow(Fixed(3), NewElement(w1)),
		NewRow(Flex(1), NewElement(w2)),
		NewRow(Fixed(1), NewElement(w3)),
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

func TestPane_PercentDistribution(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}
	w3 := &mockWidget{char: 'C'}

	root := NewRow(Percent100,
		NewColumn(Percent50, NewElement(w1)),
		NewColumn(Percent25, NewElement(w2)),
		NewColumn(Percent25, NewElement(w3)),
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

func TestPane_OptionalRemoval(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}
	w3 := &mockWidget{char: 'C'}

	root := NewRow(Percent100,
		NewColumn(Flex(1), NewElement(w1)),
		NewColumn(Flex(1), NewElement(w2)),
		NewColumn(Flex(1), NewElement(w3)).WithMinSize(30).WithOptional(true),
	)
	root.SetSize(50, 10)

	_, err := root.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w1.width+w2.width != 50 {
		t.Errorf("w1+w2 = %d, want 50", w1.width+w2.width)
	}
	if w3.width != 0 && w3.height != 0 {
		t.Errorf("w3 should not have been rendered, got %dx%d", w3.width, w3.height)
	}
}

func TestPane_NilElementsSkipped(t *testing.T) {
	w := &mockWidget{char: 'X'}
	root := NewRow(Percent100, nil, NewColumn(Flex(1), NewElement(w)), nil)
	root.SetSize(80, 10)

	_, err := root.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w.width != 80 {
		t.Errorf("w.width = %d, want 80", w.width)
	}
}
