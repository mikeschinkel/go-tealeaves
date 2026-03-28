package tealayout

import "testing"

func TestPane_SetDimension(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}

	col1 := NewColumn(Flex(1), NewElement(w1))
	col2 := NewColumn(Flex(1), NewElement(w2))

	root := NewRow(Percent100, col1, col2)
	root.SetSize(80, 10)

	// Initially 40/40
	root.Render()
	if w1.width != 40 || w2.width != 40 {
		t.Fatalf("initial: w1=%d w2=%d, want 40/40", w1.width, w2.width)
	}

	// Change col1 to Fixed(20) at runtime
	col1.SetDimension(Fixed(20))
	root.Render()
	if w1.width != 20 {
		t.Errorf("after SetDimension: w1=%d, want 20", w1.width)
	}
	if w2.width != 60 {
		t.Errorf("after SetDimension: w2=%d, want 60", w2.width)
	}
}

func TestPane_SetMinSize_Runtime(t *testing.T) {
	col1 := NewColumn(Flex(1))
	col2 := NewColumn(Flex(1))
	root := NewRow(Percent100, col1, col2)
	root.SetSize(80, 10)
	root.Resolve()

	col1.SetMinSize(50)
	sizes, _ := root.Resolve()
	if sizes[0] < 50 {
		t.Errorf("after SetMinSize(50): sizes[0]=%d, want >= 50", sizes[0])
	}
}

func TestPane_SetMaxSize_Runtime(t *testing.T) {
	col1 := NewColumn(Flex(1))
	col2 := NewColumn(Flex(1))
	root := NewRow(Percent100, col1, col2)
	root.SetSize(80, 10)

	col1.SetMaxSize(20)
	sizes, _ := root.Resolve()
	if sizes[0] > 20 {
		t.Errorf("after SetMaxSize(20): sizes[0]=%d, want <= 20", sizes[0])
	}
}
