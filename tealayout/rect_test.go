package tealayout

import "testing"

func TestRect_Fields(t *testing.T) {
	r := Rect{X: 5, Y: 3, Width: 40, Height: 20}
	if r.X != 5 || r.Y != 3 || r.Width != 40 || r.Height != 20 {
		t.Errorf("Rect fields = %+v, want {5, 3, 40, 20}", r)
	}
}
