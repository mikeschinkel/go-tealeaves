package teatextsel

import (
	"testing"
)

// --- Layer 1: Selection Tests ---

func TestNewSelection(t *testing.T) {
	s := NewSelection()
	if s.Active {
		t.Error("expected Active=false")
	}
	if s.Start.Row != 0 || s.Start.Col != 0 {
		t.Error("expected Start at zero position")
	}
	if s.End.Row != 0 || s.End.Col != 0 {
		t.Error("expected End at zero position")
	}
}

func TestSelection_Begin(t *testing.T) {
	s := NewSelection()
	s = s.Begin(Position{Row: 2, Col: 5})

	if !s.Active {
		t.Error("expected Active=true after Begin")
	}
	if s.Start.Row != 2 || s.Start.Col != 5 {
		t.Errorf("expected Start={2,5}, got {%d,%d}", s.Start.Row, s.Start.Col)
	}
	if s.End.Row != 2 || s.End.Col != 5 {
		t.Errorf("expected End={2,5}, got {%d,%d}", s.End.Row, s.End.Col)
	}
}

func TestSelection_Extend(t *testing.T) {
	s := NewSelection()
	s = s.Begin(Position{Row: 1, Col: 3})
	s = s.Extend(Position{Row: 1, Col: 8})

	if !s.Active {
		t.Error("expected Active=true after Extend")
	}
	if s.End.Row != 1 || s.End.Col != 8 {
		t.Errorf("expected End={1,8}, got {%d,%d}", s.End.Row, s.End.Col)
	}
	// Start should be unchanged
	if s.Start.Row != 1 || s.Start.Col != 3 {
		t.Errorf("expected Start unchanged at {1,3}, got {%d,%d}", s.Start.Row, s.Start.Col)
	}
}

func TestSelection_Clear(t *testing.T) {
	s := NewSelection()
	s = s.Begin(Position{Row: 1, Col: 3})
	s = s.Extend(Position{Row: 2, Col: 5})
	s = s.Clear()

	if s.Active {
		t.Error("expected Active=false after Clear")
	}
}

func TestSelection_Normalized(t *testing.T) {
	// Forward selection
	s := NewSelection()
	s = s.Begin(Position{Row: 1, Col: 3})
	s = s.Extend(Position{Row: 2, Col: 5})
	start, end := s.Normalized()

	if start.Row != 1 || start.Col != 3 {
		t.Errorf("expected normalized start={1,3}, got {%d,%d}", start.Row, start.Col)
	}
	if end.Row != 2 || end.Col != 5 {
		t.Errorf("expected normalized end={2,5}, got {%d,%d}", end.Row, end.Col)
	}

	// Backward selection (Start after End)
	s2 := NewSelection()
	s2 = s2.Begin(Position{Row: 3, Col: 8})
	s2 = s2.Extend(Position{Row: 1, Col: 2})
	start2, end2 := s2.Normalized()

	if start2.Row != 1 || start2.Col != 2 {
		t.Errorf("expected normalized start={1,2}, got {%d,%d}", start2.Row, start2.Col)
	}
	if end2.Row != 3 || end2.Col != 8 {
		t.Errorf("expected normalized end={3,8}, got {%d,%d}", end2.Row, end2.Col)
	}
}

func TestSelection_Contains(t *testing.T) {
	s := NewSelection()
	s = s.Begin(Position{Row: 1, Col: 3})
	s = s.Extend(Position{Row: 3, Col: 5})

	// Inside range
	if !s.Contains(Position{Row: 2, Col: 0}) {
		t.Error("expected Contains=true for position inside range")
	}
	// At start boundary
	if !s.Contains(Position{Row: 1, Col: 3}) {
		t.Error("expected Contains=true at start boundary")
	}
	// At end boundary
	if !s.Contains(Position{Row: 3, Col: 5}) {
		t.Error("expected Contains=true at end boundary")
	}
	// Before start
	if s.Contains(Position{Row: 0, Col: 0}) {
		t.Error("expected Contains=false before start")
	}
	// After end
	if s.Contains(Position{Row: 4, Col: 0}) {
		t.Error("expected Contains=false after end")
	}

	// Inactive selection contains nothing
	inactive := NewSelection()
	if inactive.Contains(Position{Row: 0, Col: 0}) {
		t.Error("expected inactive selection Contains=false")
	}
}

func TestSelection_IsEmpty(t *testing.T) {
	// Empty (not active)
	s := NewSelection()
	if !s.IsEmpty() {
		t.Error("expected IsEmpty=true for new selection")
	}

	// Active but Start==End
	s = s.Begin(Position{Row: 1, Col: 3})
	if !s.IsEmpty() {
		t.Error("expected IsEmpty=true when Start==End")
	}

	// Active with different End
	s = s.Extend(Position{Row: 1, Col: 5})
	if s.IsEmpty() {
		t.Error("expected IsEmpty=false with different Start and End")
	}
}

func TestSelectAll(t *testing.T) {
	lines := []string{"Hello", "World", "Foo"}
	s := SelectAll(lines)

	if !s.Active {
		t.Error("expected Active=true")
	}
	if s.Start.Row != 0 || s.Start.Col != 0 {
		t.Error("expected Start at {0,0}")
	}
	if s.End.Row != 2 || s.End.Col != 3 {
		t.Errorf("expected End at {2,3}, got {%d,%d}", s.End.Row, s.End.Col)
	}

	// Empty content
	empty := SelectAll([]string{})
	if empty.Active {
		t.Error("expected inactive selection for empty content")
	}
}

func TestSelection_BeforeAfterEqual(t *testing.T) {
	a := Position{Row: 1, Col: 3}
	b := Position{Row: 2, Col: 5}
	c := Position{Row: 1, Col: 3}

	if !a.Before(b) {
		t.Error("expected a Before b")
	}
	if a.After(b) {
		t.Error("expected a not After b")
	}
	if !b.After(a) {
		t.Error("expected b After a")
	}
	if !a.Equal(c) {
		t.Error("expected a Equal c")
	}
	if a.Before(c) {
		t.Error("expected a not Before c (same position)")
	}

	// Same row, different col
	d := Position{Row: 1, Col: 5}
	if !a.Before(d) {
		t.Error("expected a Before d (same row, smaller col)")
	}
}
