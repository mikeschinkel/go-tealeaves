package tealayout

import (
	"testing"
)

func TestMultiPaneLayout_Construction(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a, Dim: Percent(60)},
		{Name: "b", Element: b, Dim: Percent(40)},
	})

	names := mpl.PaneNames()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("PaneNames = %v, want [a b]", names)
	}

	if mpl.Pane("a") == nil {
		t.Error("Pane(a) should not be nil")
	}
	if mpl.Pane("b") == nil {
		t.Error("Pane(b) should not be nil")
	}
}

func TestMultiPaneLayout_WithHeaderFooter(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	h := NewElement(&mockWidget{char: 'H'})
	f := NewElement(&mockWidget{char: 'F'})

	mpl := NewMultiPaneLayout(
		[]PaneDef{{Name: "a", Element: a}},
		WithHeader(h),
		WithFooter(f),
	)

	if mpl.Pane("header") == nil {
		t.Error("header pane should exist")
	}
	if mpl.Pane("footer") == nil {
		t.Error("footer pane should exist")
	}
}

func TestMultiPaneLayout_Render(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a, Dim: Percent(50)},
		{Name: "b", Element: b, Dim: Percent(50)},
	})
	mpl.SetSize(80, 24)

	output, err := mpl.Render()
	if err != nil {
		t.Fatal(err)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestMultiPaneLayout_ResizeFocused(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a, Dim: Percent(50), MinFlexWeight: 0.1},
		{Name: "b", Element: b, Dim: Percent(50), MinFlexWeight: 0.1},
	})
	mpl.SetSize(100, 24)

	// Focus should be on "a" initially.
	if !mpl.Focused("a") {
		t.Fatal("initial focus should be on a")
	}

	mpl.ResizeFocused(0.1)
	pa := mpl.Pane("a")
	if pa.dim.value < 0.5 {
		t.Errorf("after grow, a.dim.value = %f, should be > 0.5", pa.dim.value)
	}
}

func TestMultiPaneLayout_Visibility(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
	})

	mpl.HidePane("b")
	vis := mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "a" {
		t.Errorf("VisiblePaneNames = %v, want [a]", vis)
	}

	mpl.ShowPane("b")
	vis = mpl.VisiblePaneNames()
	if len(vis) != 2 {
		t.Errorf("VisiblePaneNames = %v, want [a b]", vis)
	}
}

func TestMultiPaneLayout_FocusCycling(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})
	c := NewElement(&mockWidget{char: 'C'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
		{Name: "c", Element: c},
	})

	if !mpl.Focused("a") {
		t.Error("initial focus should be a")
	}

	mpl.FocusNext()
	if !mpl.Focused("b") {
		t.Error("after FocusNext, should be b")
	}

	mpl.FocusNext()
	if !mpl.Focused("c") {
		t.Error("after second FocusNext, should be c")
	}

	mpl.FocusPrev()
	if !mpl.Focused("b") {
		t.Error("after FocusPrev, should be b")
	}
}

func TestMultiPaneLayout_EnsureFocusedVisible(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
	})

	mpl.HidePane("a")
	mpl.EnsureFocusedVisible()
	if !mpl.Focused("b") {
		t.Error("should advance to b when a is hidden")
	}
}

func TestMultiPaneLayout_PaneLayout(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
	})
	if mpl.PaneLayout() == nil {
		t.Error("PaneLayout() should not be nil")
	}
}

func TestMultiPaneLayout_DefaultDim(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a}, // Dim defaults to zero → should become Flex(1)
	})
	pa := mpl.Pane("a")
	if pa.dim.value != 1.0 {
		t.Errorf("default dim value = %f, want 1.0", pa.dim.value)
	}
}

func TestMultiPaneLayout_VisibleFlexPercent(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})
	c := NewElement(&mockWidget{char: 'C'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a, Dim: Percent(25)},
		{Name: "b", Element: b, Dim: Percent(50)},
		{Name: "c", Element: c, Dim: Percent(25)},
	})

	// All visible: a=25%, b=50%, c=25%.
	if got := mpl.VisibleFlexPercent("a"); abs(got-25) > 0.1 {
		t.Errorf("a = %.1f%%, want 25%%", got)
	}
	if got := mpl.VisibleFlexPercent("b"); abs(got-50) > 0.1 {
		t.Errorf("b = %.1f%%, want 50%%", got)
	}

	// Hide c: a=33.3%, b=66.7%.
	mpl.HidePane("c")
	if got := mpl.VisibleFlexPercent("a"); abs(got-33.33) > 0.1 {
		t.Errorf("after hiding c, a = %.1f%%, want 33.3%%", got)
	}
	if got := mpl.VisibleFlexPercent("b"); abs(got-66.67) > 0.1 {
		t.Errorf("after hiding c, b = %.1f%%, want 66.7%%", got)
	}

	// Hidden pane returns 0.
	if got := mpl.VisibleFlexPercent("c"); got != 0 {
		t.Errorf("hidden c = %.1f%%, want 0%%", got)
	}

	// Unknown pane returns 0.
	if got := mpl.VisibleFlexPercent("nope"); got != 0 {
		t.Errorf("unknown pane = %.1f%%, want 0%%", got)
	}
}

func TestMultiPaneLayout_VisibleFlexPercents(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})
	c := NewElement(&mockWidget{char: 'C'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a, Dim: Percent(25)},
		{Name: "b", Element: b, Dim: Percent(50)},
		{Name: "c", Element: c, Dim: Percent(25)},
	})

	pcts := mpl.VisibleFlexPercents()
	if len(pcts) != 3 {
		t.Fatalf("len = %d, want 3", len(pcts))
	}
	if abs(pcts["a"]-25) > 0.1 || abs(pcts["b"]-50) > 0.1 || abs(pcts["c"]-25) > 0.1 {
		t.Errorf("percentages = %v, want a=25 b=50 c=25", pcts)
	}

	// Hide b: a=50%, c=50%.
	mpl.HidePane("b")
	pcts = mpl.VisibleFlexPercents()
	if len(pcts) != 2 {
		t.Fatalf("after hide, len = %d, want 2", len(pcts))
	}
	if _, ok := pcts["b"]; ok {
		t.Error("hidden pane b should not appear in map")
	}
	if abs(pcts["a"]-50) > 0.1 || abs(pcts["c"]-50) > 0.1 {
		t.Errorf("after hide, percentages = %v, want a=50 c=50", pcts)
	}
}

func TestMultiPaneLayout_TogglePane(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})
	c := NewElement(&mockWidget{char: 'C'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
		{Name: "c", Element: c},
	})

	// Toggle b off.
	mpl.TogglePane("b")
	vis := mpl.VisiblePaneNames()
	if len(vis) != 2 || vis[0] != "a" || vis[1] != "c" {
		t.Errorf("after toggle b off, visible = %v, want [a c]", vis)
	}

	// Toggle b back on.
	mpl.TogglePane("b")
	vis = mpl.VisiblePaneNames()
	if len(vis) != 3 {
		t.Errorf("after toggle b on, visible = %v, want [a b c]", vis)
	}

	// Toggle down to one, then try to toggle the last one — should be no-op.
	mpl.TogglePane("a")
	mpl.TogglePane("b")
	vis = mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "c" {
		t.Fatalf("visible = %v, want [c]", vis)
	}
	mpl.TogglePane("c") // should be no-op
	vis = mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "c" {
		t.Errorf("toggle last pane should be no-op, visible = %v", vis)
	}

	// Toggle unknown name — no-op.
	mpl.TogglePane("nope")
}

func TestMultiPaneLayout_SoloPane(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})
	c := NewElement(&mockWidget{char: 'C'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
		{Name: "c", Element: c},
	})

	mpl.SoloPane("b")
	vis := mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "b" {
		t.Errorf("after SoloPane(b), visible = %v, want [b]", vis)
	}
	if !mpl.Focused("b") {
		t.Error("SoloPane should focus the solo pane")
	}

	// Unknown name — no-op.
	mpl.SoloPane("nope")
	vis = mpl.VisiblePaneNames()
	if len(vis) != 1 || vis[0] != "b" {
		t.Errorf("SoloPane(nope) should be no-op, visible = %v", vis)
	}
}

func TestMultiPaneLayout_ShowAllPanes(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})
	c := NewElement(&mockWidget{char: 'C'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
		{Name: "c", Element: c},
	})

	mpl.HidePane("a")
	mpl.HidePane("c")
	mpl.ShowAllPanes()
	vis := mpl.VisiblePaneNames()
	if len(vis) != 3 {
		t.Errorf("after ShowAllPanes, visible = %v, want [a b c]", vis)
	}
}

func TestMultiPaneLayout_TogglePane_FocusFollows(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	b := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "b", Element: b},
	})

	// Focus a, then toggle it off — focus should move to b.
	if !mpl.Focused("a") {
		t.Fatal("initial focus should be a")
	}
	mpl.TogglePane("a")
	if !mpl.Focused("b") {
		t.Error("after toggling a off, focus should move to b")
	}
}

func TestMultiPaneLayout_VisibleFlexPercents_AllHidden(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
	})
	mpl.HidePane("a")

	if pcts := mpl.VisibleFlexPercents(); pcts != nil {
		t.Errorf("all hidden should return nil, got %v", pcts)
	}
	if got := mpl.VisibleFlexPercent("a"); got != 0 {
		t.Errorf("all hidden, a = %f, want 0", got)
	}
}

func TestMultiPaneLayout_NestedChildren_Construction(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	ct := NewElement(&mockWidget{char: 'T'})
	cb := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "group", Dim: Percent(50), Children: []PaneDef{
			{Name: "child-top", Element: ct, Dim: Percent(62)},
			{Name: "child-bot", Element: cb, Dim: Percent(38)},
		}},
	})

	// Group pane exists.
	if mpl.Pane("group") == nil {
		t.Fatal("group pane should exist")
	}
	// Child panes exist and are focusable.
	topPane := mpl.Pane("child-top")
	if topPane == nil {
		t.Fatal("child-top pane should exist")
	}
	if !topPane.IsFocusable() {
		t.Error("child-top should be focusable")
	}
	botPane := mpl.Pane("child-bot")
	if botPane == nil {
		t.Fatal("child-bot pane should exist")
	}
	if !botPane.IsFocusable() {
		t.Error("child-bot should be focusable")
	}

	// Group itself is NOT focusable.
	groupPane := mpl.Pane("group")
	if groupPane.IsFocusable() {
		t.Error("group pane should NOT be focusable")
	}

	// PaneNames includes group, not children.
	names := mpl.PaneNames()
	if len(names) != 2 || names[0] != "a" || names[1] != "group" {
		t.Errorf("PaneNames = %v, want [a group]", names)
	}
}

func TestMultiPaneLayout_NestedChildren_FocusCycling(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	ct := NewElement(&mockWidget{char: 'T'})
	cb := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "group", Children: []PaneDef{
			{Name: "child-top", Element: ct},
			{Name: "child-bot", Element: cb},
		}},
	})

	// Focus order: a → child-top → child-bot (depth-first walk).
	if !mpl.Focused("a") {
		t.Fatal("initial focus should be a")
	}
	mpl.FocusNext()
	if !mpl.Focused("child-top") {
		t.Error("after FocusNext, should be child-top")
	}
	mpl.FocusNext()
	if !mpl.Focused("child-bot") {
		t.Error("after second FocusNext, should be child-bot")
	}
	mpl.FocusPrev()
	if !mpl.Focused("child-top") {
		t.Error("after FocusPrev, should be child-top")
	}
}

func TestMultiPaneLayout_NestedChildren_GroupVisibility(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	ct := NewElement(&mockWidget{char: 'T'})
	cb := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "group", Children: []PaneDef{
			{Name: "child-top", Element: ct},
			{Name: "child-bot", Element: cb},
		}},
	})
	mpl.SetSize(80, 24)

	// Hiding the group hides it (children inside are structurally hidden).
	mpl.HidePane("group")
	groupPane := mpl.Pane("group")
	if groupPane.Visible() {
		t.Error("group should be hidden")
	}

	// Focus should leave children since group is hidden.
	mpl.EnsureFocusedVisible()
	if mpl.Focused("child-top") || mpl.Focused("child-bot") {
		t.Error("focus should not be on a child of a hidden group")
	}

	mpl.ShowPane("group")
	if !groupPane.Visible() {
		t.Error("group should be visible after ShowPane")
	}
}

func TestMultiPaneLayout_NestedChildren_Render(t *testing.T) {
	ct := NewElement(&mockWidget{char: 'T'})
	cb := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "group", Children: []PaneDef{
			{Name: "child-top", Element: ct, Dim: Percent(60)},
			{Name: "child-bot", Element: cb, Dim: Percent(40)},
		}},
	})
	mpl.SetSize(40, 20)

	output, err := mpl.Render()
	if err != nil {
		t.Fatal(err)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestMultiPaneLayout_NestedChildren_ResizeWithinGroup(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	ct := NewElement(&mockWidget{char: 'T'})
	cb := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a},
		{Name: "group", Children: []PaneDef{
			{Name: "child-top", Element: ct, Dim: Percent(50), MinFlexWeight: 0.1},
			{Name: "child-bot", Element: cb, Dim: Percent(50), MinFlexWeight: 0.1},
		}},
	})
	mpl.SetSize(80, 40)

	// Focus child-top and resize it.
	if err := mpl.FocusPane("child-top"); err != nil {
		t.Fatalf("FocusPane: %v", err)
	}
	topPane := mpl.Pane("child-top")
	oldWeight := topPane.dim.value

	mpl.ResizeFocused(0.1)
	if topPane.dim.value <= oldWeight {
		t.Errorf("child-top weight should increase: was %f, now %f", oldWeight, topPane.dim.value)
	}
}

func TestMultiPaneLayout_NestedChildren_MaxSize(t *testing.T) {
	a := NewElement(&mockWidget{char: 'A'})
	ct := NewElement(&mockWidget{char: 'T'})
	cb := NewElement(&mockWidget{char: 'B'})

	mpl := NewMultiPaneLayout([]PaneDef{
		{Name: "a", Element: a, MaxSize: 30},
		{Name: "group", MaxSize: 50, Children: []PaneDef{
			{Name: "child-top", Element: ct},
			{Name: "child-bot", Element: cb},
		}},
	})

	// Verify MaxSize was wired.
	aPane := mpl.Pane("a")
	if aPane.maxSize != 30 {
		t.Errorf("a.maxSize = %d, want 30", aPane.maxSize)
	}
	groupPane := mpl.Pane("group")
	if groupPane.maxSize != 50 {
		t.Errorf("group.maxSize = %d, want 50", groupPane.maxSize)
	}
}

func TestMultiPaneLayout_NestedChildren_PanicOnDeepNesting(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nested Children within Children")
		}
	}()

	ct := NewElement(&mockWidget{char: 'T'})
	NewMultiPaneLayout([]PaneDef{
		{Name: "group", Children: []PaneDef{
			{Name: "child", Children: []PaneDef{
				{Name: "deep", Element: ct},
			}},
		}},
	})
}
