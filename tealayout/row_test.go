package tealayout

import "testing"

func TestRow_Resolve(t *testing.T) {
	row := NewRow(Percent100,
		NewColumn(Fixed(20)),
		NewColumn(Flex(1)),
		NewColumn(Flex(1)),
	)
	row.SetSize(80, 24)

	sizes, err := row.Resolve()
	if err != nil {
		t.Fatal(err)
	}
	if sizes[0] != 20 {
		t.Errorf("sizes[0] = %d, want 20", sizes[0])
	}
	if sizes[1] != 30 {
		t.Errorf("sizes[1] = %d, want 30", sizes[1])
	}
	if sizes[2] != 30 {
		t.Errorf("sizes[2] = %d, want 30", sizes[2])
	}
}

func TestRow_ResolveWithGap(t *testing.T) {
	row := NewRow(Percent100,
		NewColumn(Flex(1)),
		NewColumn(Flex(1)),
		NewColumn(Flex(1)),
	).WithGap(2)
	row.SetSize(80, 24)

	sizes, err := row.Resolve()
	if err != nil {
		t.Fatal(err)
	}
	total := sum(sizes)
	if total != 76 {
		t.Errorf("sum = %d, want 76 (80 - 4 gap)", total)
	}
}

func TestRow_ChildRect(t *testing.T) {
	row := NewRow(Percent100,
		NewColumn(Fixed(20)),
		NewColumn(Fixed(30)),
		NewColumn(Fixed(10)),
	)
	row.SetSize(80, 24)
	_, err := row.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	r0 := row.ChildRect(0)
	if r0.X != 0 || r0.Width != 20 || r0.Height != 24 {
		t.Errorf("ChildRect(0) = %+v, want X=0 W=20 H=24", r0)
	}

	r1 := row.ChildRect(1)
	if r1.X != 20 || r1.Width != 30 {
		t.Errorf("ChildRect(1) = %+v, want X=20 W=30", r1)
	}

	r2 := row.ChildRect(2)
	if r2.X != 50 || r2.Width != 10 {
		t.Errorf("ChildRect(2) = %+v, want X=50 W=10", r2)
	}
}

func TestRow_ChildRectWithGap(t *testing.T) {
	row := NewRow(Percent100,
		NewColumn(Fixed(20)),
		NewColumn(Fixed(30)),
	).WithGap(2)
	row.SetSize(80, 24)
	_, err := row.Resolve()
	if err != nil {
		t.Fatal(err)
	}

	r0 := row.ChildRect(0)
	if r0.X != 0 || r0.Width != 20 {
		t.Errorf("ChildRect(0) = %+v, want X=0 W=20", r0)
	}

	r1 := row.ChildRect(1)
	if r1.X != 22 || r1.Width != 30 {
		t.Errorf("ChildRect(1) = %+v, want X=22 W=30", r1)
	}
}

func TestRow_CachedResolve(t *testing.T) {
	row := NewRow(Percent100, NewColumn(Flex(1)))
	row.SetSize(80, 24)

	s1, _ := row.Resolve()
	s2, _ := row.Resolve()
	if s1[0] != s2[0] {
		t.Error("cached resolve returned different result")
	}
}

func TestRow_SetSizeInvalidatesCache(t *testing.T) {
	p := NewRow(Percent100, NewColumn(Flex(1)))
	p.SetSize(80, 24)
	if _, err := p.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	p.SetSize(100, 24)
	if p.resolved {
		t.Error("SetSize should invalidate cache")
	}
}

func TestRow_Direction(t *testing.T) {
	row := NewRow(Percent100)
	if row.Direction() != Horizontal {
		t.Errorf("Direction() = %v, want Horizontal", row.Direction())
	}
}
