package tealayout

import (
	"errors"
	"testing"
)

func TestPaneLayout_SetSizeAndRender(t *testing.T) {
	w := &mockWidget{char: 'X'}
	root := NewRow(Percent100,
		NewColumn(Flex(1), NewElement(w)).WithName("main"),
	)
	pl := NewPaneLayout(root)
	pl.SetSize(80, 24)

	output, err := pl.Render()
	if err != nil {
		t.Fatal(err)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}
	if w.width != 80 || w.height != 24 {
		t.Errorf("widget size = %dx%d, want 80x24", w.width, w.height)
	}
}

func TestPaneLayout_FocusNavigation(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	pl := NewPaneLayout(root)

	if !pl.Focused("tree") {
		t.Error("initial focus should be tree")
	}

	pl.FocusNext()
	if !pl.Focused("code") {
		t.Error("after FocusNext, should focus code")
	}

	pl.FocusPrev()
	if !pl.Focused("tree") {
		t.Error("after FocusPrev, should focus tree")
	}

	err := pl.FocusPane("code")
	if err != nil {
		t.Fatal(err)
	}
	if pl.FocusedPane() != code {
		t.Error("FocusPane should focus code")
	}
}

func TestPaneLayout_FocusPaneNotFound(t *testing.T) {
	root := NewRow(Percent100, NewColumn(Flex(1)).WithName("a").WithFocusable())
	pl := NewPaneLayout(root)

	err := pl.FocusPane("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrPaneNotFound) {
		t.Errorf("expected ErrPaneNotFound, got %v", err)
	}
}

func TestPaneLayout_ShowHide(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	pl := NewPaneLayout(root)

	pl.HidePane("tree")
	if tree.Visible() {
		t.Error("tree should be hidden")
	}

	pl.ShowPane("tree")
	if !tree.Visible() {
		t.Error("tree should be visible")
	}

	pl.SetPaneVisible("code", false)
	if code.Visible() {
		t.Error("code should be hidden")
	}
}

func TestPaneLayout_EnsureFocusedVisible(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree").WithFocusable()
	code := NewColumn(Flex(1)).WithName("code").WithFocusable()
	root := NewRow(Percent100, tree, code)
	pl := NewPaneLayout(root)

	pl.HidePane("tree")
	pl.EnsureFocusedVisible()
	if !pl.Focused("code") {
		t.Error("should auto-advance to code when tree is hidden")
	}
}

func TestPaneLayout_Pane(t *testing.T) {
	tree := NewColumn(Flex(1)).WithName("tree")
	root := NewRow(Percent100, tree)
	pl := NewPaneLayout(root)

	if pl.Pane("tree") != tree {
		t.Error("Pane should return named pane")
	}
	if pl.Pane("nonexistent") != nil {
		t.Error("Pane should return nil for unknown name")
	}
}

func TestPaneLayout_Layout(t *testing.T) {
	root := NewRow(Percent100)
	pl := NewPaneLayout(root)
	if pl.Layout() == nil {
		t.Error("Layout() should not return nil")
	}
}

// --- ResizeFocused tests ---

func TestPaneLayout_ResizeFocused_GrowClamped(t *testing.T) {
	// Three flex panes: a=0.5, b=0.3, c=0.2 with minFlexWeight=0.1 each.
	a := NewColumn(Flex(0.5)).WithName("a").WithFocusable().WithMinFlexWeight(0.1)
	b := NewColumn(Flex(0.3)).WithName("b").WithFocusable().WithMinFlexWeight(0.1)
	c := NewColumn(Flex(0.2)).WithName("c").WithFocusable().WithMinFlexWeight(0.1)
	root := NewRow(Percent100, a, b, c)
	pl := NewPaneLayout(root)
	pl.SetSize(100, 24)

	// Focus "a" and grow by a huge delta — should clamp.
	// Max for "a" = 0.5 + (0.3-0.1) + (0.2-0.1) = 0.8
	pl.ResizeFocused(10.0)
	gotA := a.dim.value
	wantA := 0.8
	if abs(gotA-wantA) > 0.001 {
		t.Errorf("after huge grow, a.dim.value = %f, want %f", gotA, wantA)
	}
	// Siblings should have shrunk to their floors.
	if abs(b.dim.value-0.1) > 0.001 {
		t.Errorf("b should be at floor 0.1, got %f", b.dim.value)
	}
	if abs(c.dim.value-0.1) > 0.001 {
		t.Errorf("c should be at floor 0.1, got %f", c.dim.value)
	}
}

func TestPaneLayout_ResizeFocused_WeightConserved(t *testing.T) {
	// Verify total weight is preserved after resize.
	a := NewColumn(Flex(0.4)).WithName("a").WithFocusable().WithMinFlexWeight(0.05)
	b := NewColumn(Flex(0.35)).WithName("b").WithFocusable().WithMinFlexWeight(0.05)
	c := NewColumn(Flex(0.25)).WithName("c").WithFocusable().WithMinFlexWeight(0.05)
	root := NewRow(Percent100, a, b, c)
	pl := NewPaneLayout(root)
	pl.SetSize(100, 24)

	totalBefore := a.dim.value + b.dim.value + c.dim.value

	pl.ResizeFocused(-0.1)

	totalAfter := a.dim.value + b.dim.value + c.dim.value
	if abs(totalAfter-totalBefore) > 0.001 {
		t.Errorf("total weight changed: before=%f, after=%f", totalBefore, totalAfter)
	}
	// "a" should have shrunk
	if a.dim.value >= 0.4 {
		t.Errorf("focused pane should have shrunk, got %f", a.dim.value)
	}
	// siblings should have grown
	if b.dim.value <= 0.35 {
		t.Errorf("sibling b should have grown, got %f", b.dim.value)
	}
}

func TestPaneLayout_ResizeFocused_ShrinkClamped(t *testing.T) {
	a := NewColumn(Flex(0.5)).WithName("a").WithFocusable().WithMinFlexWeight(0.15)
	b := NewColumn(Flex(0.5)).WithName("b").WithFocusable().WithMinFlexWeight(0.1)
	root := NewRow(Percent100, a, b)
	pl := NewPaneLayout(root)
	pl.SetSize(100, 24)

	// Shrink "a" massively — should clamp to its own minFlexWeight.
	pl.ResizeFocused(-10.0)
	gotA := a.dim.value
	wantA := 0.15
	if abs(gotA-wantA) > 0.001 {
		t.Errorf("after huge shrink, a.dim.value = %f, want %f", gotA, wantA)
	}
	// "b" should have absorbed the freed weight.
	wantB := 0.85 // 0.5 + (0.5 - 0.15)
	if abs(b.dim.value-wantB) > 0.001 {
		t.Errorf("b should have grown to %f, got %f", wantB, b.dim.value)
	}
}

func TestPaneLayout_ResizeFocused_FixedPaneNoop(t *testing.T) {
	a := NewColumn(Fixed(20)).WithName("a").WithFocusable()
	b := NewColumn(Flex(1)).WithName("b").WithFocusable()
	root := NewRow(Percent100, a, b)
	pl := NewPaneLayout(root)
	pl.SetSize(100, 24)

	// Focus "a" (fixed) and try to resize — should be no-op.
	pl.ResizeFocused(0.1)
	if a.dim.kind != dimensionFixed || a.dim.value != 20 {
		t.Errorf("fixed pane should not change, got %+v", a.dim)
	}
}

func TestPaneLayout_ResizeFocused_SingleVisibleNoop(t *testing.T) {
	a := NewColumn(Flex(0.5)).WithName("a").WithFocusable()
	b := NewColumn(Flex(0.5)).WithName("b").WithFocusable()
	root := NewRow(Percent100, a, b)
	pl := NewPaneLayout(root)
	pl.SetSize(100, 24)

	// Hide b so only a is visible.
	b.SetVisible(false)
	pl.ResizeFocused(0.1)
	if a.dim.value != 0.5 {
		t.Errorf("single visible pane should not change, got %f", a.dim.value)
	}
}

func TestPaneLayout_ResizeFocused_NoFocusNoop(t *testing.T) {
	root := NewRow(Percent100)
	pl := NewPaneLayout(root)

	// No focused pane — should not panic.
	pl.ResizeFocused(0.1)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// mockHintWidget implements SetSizer, ContentProvider, and SizeHinter for resize tests.
type mockHintWidget struct {
	width, height int
	minWidth      int
}

func (m *mockHintWidget) SetSize(w, h int) { m.width = w; m.height = h }
func (m *mockHintWidget) Content() string    { return "" }
func (m *mockHintWidget) SizeHint(availW, availH int) SizeHint {
	return SizeHint{
		Min:     Size{Width: m.minWidth},
		Desired: Size{Width: m.minWidth},
		Max:     Size{Width: -1, Height: -1},
	}
}

func TestPaneLayout_ResizeFocused_CellClampRollback(t *testing.T) {
	// Three panes with minSizeFit. Pane A has minWidth=30, B and C have
	// minWidth=10. At 90 cols with equal weights, each gets exactly 30
	// cells. A is already at its cell minimum, so shrinking it should
	// produce no visual change — weights must be rolled back.
	wA := &mockHintWidget{minWidth: 30}
	wB := &mockHintWidget{minWidth: 10}
	wC := &mockHintWidget{minWidth: 10}

	a := NewColumn(Flex(1), NewElement(wA)).WithName("a").WithFocusable().WithMinSizeFit().WithMinFlexWeight(0.05)
	b := NewColumn(Flex(1), NewElement(wB)).WithName("b").WithFocusable().WithMinSizeFit().WithMinFlexWeight(0.05)
	c := NewColumn(Flex(1), NewElement(wC)).WithName("c").WithFocusable().WithMinSizeFit().WithMinFlexWeight(0.05)
	root := NewRow(Percent100, a, b, c)
	pl := NewPaneLayout(root)
	pl.SetSize(90, 24)

	// Force initial resolve so sizes are populated.
	if _, err := pl.Render(); err != nil {
		t.Fatal(err)
	}

	// Record pre-shrink state.
	origA := a.dim.value
	origB := b.dim.value
	origC := c.dim.value

	// Try to shrink "a" — it's at 30 cells (its minSizeFit), so the
	// resolver will clamp it right back. No visual change → rollback.
	pl.ResizeFocused(-0.3)

	if a.dim.value != origA {
		t.Errorf("a weight changed from %f to %f; expected rollback", origA, a.dim.value)
	}
	if b.dim.value != origB {
		t.Errorf("b weight changed from %f to %f; expected rollback", origB, b.dim.value)
	}
	if c.dim.value != origC {
		t.Errorf("c weight changed from %f to %f; expected rollback", origC, c.dim.value)
	}

	// Now verify that growing DOES work: B and C can shrink from 30 to 10,
	// so A has room to grow.
	pl.ResizeFocused(0.3)
	if a.dim.value <= origA {
		t.Errorf("grow should have increased a's weight, got %f (was %f)", a.dim.value, origA)
	}
}

func TestPaneLayout_ResizeFocused_BeforeRender(t *testing.T) {
	// ResizeFocused before any SetSize/Render — parent.sizes is nil.
	// Should not panic, and weights should be applied (no rollback).
	a := NewColumn(Flex(0.5)).WithName("a").WithFocusable().WithMinFlexWeight(0.1)
	b := NewColumn(Flex(0.5)).WithName("b").WithFocusable().WithMinFlexWeight(0.1)
	root := NewRow(Percent100, a, b)
	pl := NewPaneLayout(root)

	// No SetSize or Render — sizes slice is nil.
	pl.ResizeFocused(0.1)

	wantA := 0.6
	if abs(a.dim.value-wantA) > 0.001 {
		t.Errorf("a.dim.value = %f, want %f", a.dim.value, wantA)
	}
	wantB := 0.4
	if abs(b.dim.value-wantB) > 0.001 {
		t.Errorf("b.dim.value = %f, want %f", b.dim.value, wantB)
	}
}

func TestPaneLayout_MarkDirty(t *testing.T) {
	w := &mockWidget{char: 'X'}
	root := NewRow(Percent100, NewColumn(Flex(1), NewElement(w)).WithName("main"))
	pl := NewPaneLayout(root)
	pl.SetSize(80, 24)
	if _, err := pl.Render(); err != nil {
		t.Fatalf("initial Render: %v", err)
	}

	pl.MarkDirty()
	// Should re-render without error
	_, err := pl.Render()
	if err != nil {
		t.Fatal(err)
	}
}
