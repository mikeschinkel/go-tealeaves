package tealayout

import "testing"

func TestSizeHint_Fields(t *testing.T) {
	hint := SizeHint{
		Min:     Size{10, 5},
		Desired: Size{40, 20},
		Max:     Size{-1, -1},
	}
	if hint.Min.Width != 10 {
		t.Errorf("Min.Width = %d, want 10", hint.Min.Width)
	}
	if hint.Desired.Height != 20 {
		t.Errorf("Desired.Height = %d, want 20", hint.Desired.Height)
	}
	if hint.Max.Width != -1 {
		t.Errorf("Max.Width = %d, want -1", hint.Max.Width)
	}
}
