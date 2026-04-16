// Source: site/src/content/docs/components/layout-engine.mdx:26#ab81f9f0,127#20867427,163#16248c0d,180#2c3f1373,196#f5342dd2,218#bc5e7584,229#b75a3215,283#1821bbdc,295#b96d6152,325#587b6502,343#93189eb2,405#130c8e87,422#59347abe,473#c002d05d,525#a8517c21,544#ab8ac507,571#f83d0b01,595#b16ed7ad,624#1c22e32e
package examples_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teacrumbs"
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

// TestCompile_FocusManager verifies FocusManager creation and FocusPane from layout-engine.mdx.
func TestCompile_FocusManager(t *testing.T) {
	sidebarEl := tealayout.NewElement(&testWidget{})
	mainEl := tealayout.NewElement(&testWidget{})

	root := tealayout.NewRow(tealayout.Flex(1),
		tealayout.NewColumn(tealayout.Flex(0.3), sidebarEl).WithName("sidebar").WithFocusable(),
		tealayout.NewColumn(tealayout.Flex(0.7), mainEl).WithName("main").WithFocusable(),
	)
	layout := tealayout.NewLayout(root)
	layout.SetSize(80, 24)

	fm := tealayout.NewFocusManager(layout)
	err := fm.FocusPane("main")
	if err != nil {
		t.Fatal(err)
	}
	_ = fm.Focused("sidebar")
	_ = fm.FocusedPane()
}

// TestCompile_RectType verifies that tealayout.Rect and tealayout.Size types exist with correct fields.
func TestCompile_RectType(t *testing.T) {
	r := tealayout.Rect{X: 1, Y: 2, Width: 80, Height: 24}
	_ = r.X
	_ = r.Y
	_ = r.Width
	_ = r.Height

	s := tealayout.Size{Width: 80, Height: 24}
	_ = s.Width
	_ = s.Height
}

// TestCompile_SizeHintType verifies that tealayout.SizeHint exists with Min/Desired/Max fields.
func TestCompile_SizeHintType(t *testing.T) {
	hint := tealayout.SizeHint{
		Min:     tealayout.Size{Width: 10, Height: 5},
		Desired: tealayout.Size{Width: 40, Height: 20},
		Max:     tealayout.Size{Width: -1, Height: -1},
	}
	_ = hint.Min
	_ = hint.Desired
	_ = hint.Max
}

// TestCompile_PaneLayoutExtended verifies FocusPane, ShowPane, HidePane, SetPaneVisible, ResizeFocused on PaneLayout.
func TestCompile_PaneLayoutExtended(t *testing.T) {
	sidebarEl := tealayout.NewElement(&testWidget{})
	mainEl := tealayout.NewElement(&testWidget{})

	root := tealayout.NewRow(tealayout.Flex(1),
		tealayout.NewColumn(tealayout.Flex(0.3), sidebarEl).WithName("sidebar").WithFocusable(),
		tealayout.NewColumn(tealayout.Flex(0.7), mainEl).WithName("main").WithFocusable(),
	)
	pl := tealayout.NewPaneLayout(root)
	pl.SetSize(80, 24)
	_, _ = pl.Render()

	err := pl.FocusPane("sidebar")
	if err != nil {
		t.Fatal(err)
	}
	pl.ShowPane("sidebar")
	pl.HidePane("sidebar")
	pl.SetPaneVisible("sidebar", true)
	pl.ResizeFocused(0.1)
}

// TestCompile_MultiPaneLayoutWithOptions verifies NewMultiPaneLayout with WithHeader, WithFooter, WithContentGap.
func TestCompile_MultiPaneLayoutWithOptions(t *testing.T) {
	header := &testWidget{}
	footer := &testWidget{}
	editor := &testWidget{}
	refs := &testWidget{}

	panes := []tealayout.PaneDef{
		{Name: "editor", Element: tealayout.NewElement(editor), Dim: tealayout.Flex(0.7)},
		{Name: "refs", Element: tealayout.NewElement(refs), Dim: tealayout.Flex(0.3), Optional: true},
	}

	mpl := tealayout.NewMultiPaneLayout(panes,
		tealayout.WithHeader(tealayout.NewElement(header)),
		tealayout.WithFooter(tealayout.NewElement(footer)),
		tealayout.WithContentGap(1),
	)
	mpl.SetSize(80, 24)
	_, err := mpl.Render()
	if err != nil {
		t.Fatal(err)
	}
}

// stackTestView is a minimal StackView implementation for compile testing.
type stackTestView struct {
	width, height int
}

func (v *stackTestView) Init() tea.Cmd                          { return nil }
func (v *stackTestView) Update(tea.Msg) (tea.Model, tea.Cmd)   { return v, nil }
func (v *stackTestView) View() tea.View                         { return tea.View{} }
func (v *stackTestView) OnEnter() tea.Cmd                       { return nil }
func (v *stackTestView) OnExit() tea.Cmd                        { return nil }
func (v *stackTestView) Breadcrumb() teacrumbs.Crumb            { return teacrumbs.NewCrumb("test", nil) }
func (v *stackTestView) SetSize(width, height int)              { v.width = width; v.height = height }

// TestCompile_StackLayoutModel verifies NewStackLayoutModel, Push, Pop, GetCached, DeleteCached.
func TestCompile_StackLayoutModel(t *testing.T) {
	initial := &stackTestView{}
	styles := teacrumbs.DefaultStyles()

	m := tealayout.NewStackLayoutModel(initial, styles)
	m.SetSize(80, 24)

	second := &stackTestView{}
	_ = m.Push(second, "second")

	_, ok := m.GetCached("second")
	_ = ok

	m.DeleteCached("second")

	_, err := m.Pop()
	if err != nil {
		t.Fatal(err)
	}
}

// TestCompile_VisibilityRotator verifies NewVisibilityRotator, Next, Prev, SetIndex, Index, Current, Len, Apply.
func TestCompile_VisibilityRotator(t *testing.T) {
	sidebar := &testWidget{}
	editor := &testWidget{}
	refs := &testWidget{}

	panes := []tealayout.PaneDef{
		{Name: "sidebar", Element: tealayout.NewElement(sidebar), Dim: tealayout.Flex(0.3)},
		{Name: "editor", Element: tealayout.NewElement(editor), Dim: tealayout.Flex(0.7)},
		{Name: "refs", Element: tealayout.NewElement(refs), Dim: tealayout.Flex(0.3), Optional: true},
	}
	mpl := tealayout.NewMultiPaneLayout(panes)
	mpl.SetSize(80, 24)

	combos := [][]string{
		{"sidebar", "editor"},
		{"editor"},
		{"sidebar", "editor", "refs"},
	}
	vr := tealayout.NewVisibilityRotator(mpl, combos)

	vr.Next()
	vr.Prev()
	vr.SetIndex(2)
	_ = vr.Index()
	_ = vr.Current()
	_ = vr.Len()
	vr.Apply()
}

// TestCompile_TreeContentLayout verifies NewTreeContentLayout with typed widget accessors.
func TestCompile_TreeContentLayout(t *testing.T) {
	tree := &testWidget{}
	content := &testWidget{}

	tcl := tealayout.NewTreeContentLayout(tree, content)
	tcl.SetSize(80, 24)
	_, err := tcl.Render()
	if err != nil {
		t.Fatal(err)
	}

	_ = tcl.Tree()
	_ = tcl.Content()
	tcl.FocusTree()
	tcl.FocusContent()
	tcl.ToggleFocus()
	_ = tcl.TreeFocused()
	_ = tcl.ContentFocused()
}
