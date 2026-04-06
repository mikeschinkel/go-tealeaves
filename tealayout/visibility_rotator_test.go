package tealayout

import (
	"testing"
)

func newTestMPL3() *MultiPaneLayout {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})
	c := NewElement(&mockWidget{char: 'C'})
	return NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
		{Name: "c", Element: c},
	})
}

func TestVisibilityRotator_NextWraps(t *testing.T) {
	mpl := newTestMPL3()
	combos := [][]string{{"a", "b", "c"}, {"a", "b"}, {"a"}}
	r := NewVisibilityRotator(mpl, combos)

	if r.Index() != 0 {
		t.Fatalf("initial index = %d, want 0", r.Index())
	}

	r.Next() // → 1
	if r.Index() != 1 {
		t.Errorf("after Next, index = %d, want 1", r.Index())
	}
	vis := mpl.VisiblePaneNames()
	if len(vis) != 2 || vis[0] != "a" || vis[1] != "b" {
		t.Errorf("visible = %v, want [a b]", vis)
	}

	r.Next() // → 2
	r.Next() // → 0 (wrap)
	if r.Index() != 0 {
		t.Errorf("after wrap, index = %d, want 0", r.Index())
	}
	vis = mpl.VisiblePaneNames()
	if len(vis) != 3 {
		t.Errorf("visible = %v, want [a b c]", vis)
	}
}

func TestVisibilityRotator_PrevWraps(t *testing.T) {
	mpl := newTestMPL3()
	combos := [][]string{{"a", "b", "c"}, {"a"}, {"b"}}
	r := NewVisibilityRotator(mpl, combos)

	r.Prev() // → 2 (wrap backwards)
	if r.Index() != 2 {
		t.Errorf("after Prev from 0, index = %d, want 2", r.Index())
	}
	vis := mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "b" {
		t.Errorf("visible = %v, want [b]", vis)
	}

	r.Prev() // → 1
	if r.Index() != 1 {
		t.Errorf("index = %d, want 1", r.Index())
	}
}

func TestVisibilityRotator_Apply(t *testing.T) {
	mpl := newTestMPL3()
	combos := [][]string{{"a"}, {"b"}}
	r := NewVisibilityRotator(mpl, combos)

	// Apply initial combo (index 0) without Next/Prev
	r.Apply()
	vis := mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "a" {
		t.Errorf("after Apply, visible = %v, want [a]", vis)
	}
}

func TestVisibilityRotator_SetIndex(t *testing.T) {
	mpl := newTestMPL3()
	combos := [][]string{{"a", "b", "c"}, {"b"}, {"c"}}
	r := NewVisibilityRotator(mpl, combos)

	r.SetIndex(2)
	if r.Index() != 2 {
		t.Errorf("index = %d, want 2", r.Index())
	}
	vis := mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "c" {
		t.Errorf("visible = %v, want [c]", vis)
	}
}

func TestVisibilityRotator_SetIndex_Panic(t *testing.T) {
	mpl := newTestMPL3()
	r := NewVisibilityRotator(mpl, [][]string{{"a"}})

	defer func() {
		if recover() == nil {
			t.Error("SetIndex out of range should panic")
		}
	}()
	r.SetIndex(5)
}

func TestVisibilityRotator_EmptyCombos_Panic(t *testing.T) {
	mpl := newTestMPL3()
	defer func() {
		if recover() == nil {
			t.Error("empty combos should panic")
		}
	}()
	NewVisibilityRotator(mpl, [][]string{})
}

func TestVisibilityRotator_Current(t *testing.T) {
	mpl := newTestMPL3()
	combos := [][]string{{"a", "b"}, {"c"}}
	r := NewVisibilityRotator(mpl, combos)

	cur := r.Current()
	if len(cur) != 2 || cur[0] != "a" || cur[1] != "b" {
		t.Errorf("Current = %v, want [a b]", cur)
	}

	// Mutating returned slice should not affect rotator
	cur[0] = "x"
	cur2 := r.Current()
	if cur2[0] != "a" {
		t.Error("Current should return a copy")
	}
}

func TestVisibilityRotator_Len(t *testing.T) {
	mpl := newTestMPL3()
	r := NewVisibilityRotator(mpl, [][]string{{"a"}, {"b"}, {"c"}})
	if r.Len() != 3 {
		t.Errorf("Len = %d, want 3", r.Len())
	}
}

func TestVisibilityRotator_FocusFollows(t *testing.T) {
	mpl := newTestMPL3()
	combos := [][]string{{"a", "b", "c"}, {"b", "c"}}
	r := NewVisibilityRotator(mpl, combos)

	// Focus a, then switch to combo that hides a
	if !mpl.Focused("a") {
		t.Fatal("initial focus should be a")
	}

	r.Next() // hides a → focus should move
	if mpl.Focused("a") {
		t.Error("focus should not remain on hidden pane a")
	}
	fp := mpl.FocusedPane()
	if fp == nil || !fp.visible {
		t.Error("focused pane should be visible")
	}
}

func TestVisibilityRotator_UnknownPaneInCombo(t *testing.T) {
	mpl := newTestMPL3()
	// "z" doesn't exist — should be silently ignored
	combos := [][]string{{"a", "z"}}
	r := NewVisibilityRotator(mpl, combos)

	r.Apply()
	vis := mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "a" {
		t.Errorf("visible = %v, want [a]", vis)
	}
}
