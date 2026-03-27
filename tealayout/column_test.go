package tealayout

import "testing"

func TestColumn_Resolve(t *testing.T) {
	col := NewColumn(Percent100,
		NewRow(Fixed(3)),
		NewRow(Flex(1)),
		NewRow(Fixed(1)),
	)
	col.SetSize(80, 24)

	sizes, err := col.Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 3 {
		t.Errorf("sizes[0] = %d, want 3", sizes[0])
	}
	if sizes[1] != 20 {
		t.Errorf("sizes[1] = %d, want 20", sizes[1])
	}
	if sizes[2] != 1 {
		t.Errorf("sizes[2] = %d, want 1", sizes[2])
	}
}

func TestColumn_ResolveWithGap(t *testing.T) {
	col := NewColumn(Percent100,
		NewRow(Flex(1)),
		NewRow(Flex(1)),
	).WithGap(1)
	col.SetSize(80, 24)

	sizes, err := col.Resolve()
	if err != nil {
		t.Fatal(err)
	}
	// 24 - 1(gap) = 23. 23/2 = 12+11 or 11+12
	if sum(sizes) != 23 {
		t.Errorf("sum = %d, want 23", sum(sizes))
	}
}

func TestColumn_ChildRect(t *testing.T) {
	col := NewColumn(Percent100,
		NewRow(Fixed(3)),
		NewRow(Fixed(20)),
		NewRow(Fixed(1)),
	)
	col.SetSize(80, 24)
	_, err := col.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	r0 := col.ChildRect(0)
	if r0.Y != 0 || r0.Height != 3 || r0.Width != 80 {
		t.Errorf("ChildRect(0) = %+v, want Y=0 H=3 W=80", r0)
	}

	r1 := col.ChildRect(1)
	if r1.Y != 3 || r1.Height != 20 {
		t.Errorf("ChildRect(1) = %+v, want Y=3 H=20", r1)
	}

	r2 := col.ChildRect(2)
	if r2.Y != 23 || r2.Height != 1 {
		t.Errorf("ChildRect(2) = %+v, want Y=23 H=1", r2)
	}
}

func TestColumn_Direction(t *testing.T) {
	col := NewColumn(Percent100)
	if col.Direction() != Vertical {
		t.Errorf("Direction() = %v, want Vertical", col.Direction())
	}
}
