package examples_test

import (
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
)

// testWidget is a minimal widget for testing layout.
type testWidget struct {
	width, height int
}

func (w *testWidget) View() string              { return "" }
func (w *testWidget) SetSize(width, height int) { w.width = width; w.height = height }
func (w *testWidget) Style() lipgloss.Style     { return lipgloss.Style{} }

// TestCompile_LayoutQuickExample verifies the quick example from layout-engine.mdx.
func TestCompile_LayoutQuickExample(t *testing.T) {
	header := tealayout.NewElement(&testWidget{})
	sidebar := tealayout.NewElement(&testWidget{})
	main := tealayout.NewElement(&testWidget{})

	root := tealayout.NewColumn(tealayout.Percent100,
		tealayout.NewRow(tealayout.Fixed(1), header),
		tealayout.NewRow(tealayout.Flex(1),
			tealayout.NewColumn(tealayout.Flex(0.25), sidebar),
			tealayout.NewColumn(tealayout.Flex(0.75), main),
		),
	)

	layout := tealayout.NewLayout(root)
	layout.SetSize(80, 24)
	_, err := layout.Render()
	if err != nil {
		t.Fatal(err)
	}
}

// TestCompile_PercentConstants verifies percent dimension constants from layout-engine.mdx.
func TestCompile_PercentConstants(t *testing.T) {
	_ = tealayout.Percent100
	_ = tealayout.Percent75
	_ = tealayout.Percent50
	_ = tealayout.Percent33
	_ = tealayout.Percent25
	_ = tealayout.Percent20
}

// TestCompile_PaneTreeModifiers verifies pane modifier methods from layout-engine.mdx.
func TestCompile_PaneTreeModifiers(t *testing.T) {
	el := tealayout.NewElement(&testWidget{})

	sidebar := tealayout.NewColumn(tealayout.Percent(25), el).
		WithOptional(true).
		WithMinSize(20)

	main := tealayout.NewColumn(tealayout.Flex(1), tealayout.NewElement(&testWidget{}))

	root := tealayout.NewRow(tealayout.Percent100, sidebar, main)
	layout := tealayout.NewLayout(root)
	layout.SetSize(80, 24)
	_, err := layout.Render()
	if err != nil {
		t.Fatal(err)
	}
}

// TestCompile_ElementAccess verifies Element type-safe access from layout-engine.mdx.
func TestCompile_ElementAccess(t *testing.T) {
	myWidget := &testWidget{}
	el := tealayout.NewElement(myWidget)
	root := tealayout.NewRow(tealayout.Flex(1), el)
	layout := tealayout.NewLayout(root)
	layout.SetSize(80, 24)
	_, _ = layout.Render()

	_ = el.Widget()
}

// TestCompile_StringElement verifies StringElement from layout-engine.mdx.
func TestCompile_StringElement(t *testing.T) {
	el := tealayout.StringElement("hello")
	root := tealayout.NewRow(tealayout.Flex(1), el)
	layout := tealayout.NewLayout(root)
	layout.SetSize(80, 24)
	_, err := layout.Render()
	if err != nil {
		t.Fatal(err)
	}
}

// TestCompile_PaneLayout verifies PaneLayout from layout-engine.mdx.
func TestCompile_PaneLayout(t *testing.T) {
	sidebarEl := tealayout.NewElement(&testWidget{})
	mainEl := tealayout.NewElement(&testWidget{})

	root := tealayout.NewRow(tealayout.Flex(1),
		tealayout.NewColumn(tealayout.Flex(0.3), sidebarEl).WithName("sidebar").WithFocusable(),
		tealayout.NewColumn(tealayout.Flex(0.7), mainEl).WithName("main").WithFocusable(),
	)
	pl := tealayout.NewPaneLayout(root)
	pl.SetSize(80, 24)

	_, _ = pl.Render()
	pl.FocusNext()
	pl.FocusPrev()
	_ = pl.Focused("main")
}

// TestCompile_MultiPaneLayout verifies MultiPaneLayout from layout-engine.mdx.
func TestCompile_MultiPaneLayout(t *testing.T) {
	sidebar := &testWidget{}
	editor := &testWidget{}
	refs := &testWidget{}

	panes := []tealayout.PaneDef{
		{Name: "sidebar", Element: tealayout.NewElement(sidebar), Dim: tealayout.Flex(0.3), MinSize: 20},
		{Name: "editor", Element: tealayout.NewElement(editor), Dim: tealayout.Flex(0.7)},
		{Name: "refs", Element: tealayout.NewElement(refs), Dim: tealayout.Flex(0.3), Optional: true},
	}

	mpl := tealayout.NewMultiPaneLayout(panes)
	mpl.SetSize(80, 24)
	_, err := mpl.Render()
	if err != nil {
		t.Fatal(err)
	}
}

// TestCompile_Alignment verifies alignment constants from layout-engine.mdx.
func TestCompile_Alignment(t *testing.T) {
	el := tealayout.NewElement(&testWidget{})
	sidebar := tealayout.NewColumn(tealayout.Flex(1), el)
	sidebar.WithAlignment(tealayout.MiddleCenter)

	el2 := tealayout.NewElement(&testWidget{})
	header := tealayout.NewRow(tealayout.Fixed(1), el2)
	header.WithAlignment(tealayout.TopLeft)
}
