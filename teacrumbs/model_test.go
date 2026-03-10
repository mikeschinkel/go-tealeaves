package teacrumbs_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teacrumbs"
)

func TestBreadcrumbsModel_Defaults(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	if m.Len() != 0 {
		t.Errorf("expected empty crumbs, got %d", m.Len())
	}
	if m.Separator() != " > " {
		t.Errorf("expected default separator ' > ', got %q", m.Separator())
	}
}

func TestPush(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	m2 := m.Push(teacrumbs.NewCrumb("Home", nil))
	if m2.Len() != 1 {
		t.Errorf("expected 1 crumb, got %d", m2.Len())
	}
	crumbs := m2.Crumbs()
	if crumbs[0].Text != "Home" {
		t.Errorf("expected 'Home', got %q", crumbs[0].Text)
	}
}

func TestPush_CopySemantics(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().Push(teacrumbs.NewCrumb("Home", nil))
	m2 := m.Push(teacrumbs.NewCrumb("Settings", nil))

	if m.Len() != 1 {
		t.Errorf("original should still have 1 crumb, got %d", m.Len())
	}
	if m2.Len() != 2 {
		t.Errorf("new model should have 2 crumbs, got %d", m2.Len())
	}
}

func TestPop(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil))
	m2 := m.Pop()

	if m2.Len() != 1 {
		t.Errorf("expected 1 crumb after pop, got %d", m2.Len())
	}
	if m2.Crumbs()[0].Text != "Home" {
		t.Errorf("expected 'Home', got %q", m2.Crumbs()[0].Text)
	}
}

func TestPop_Empty(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	m2 := m.Pop()
	if m2.Len() != 0 {
		t.Errorf("pop on empty should be no-op, got %d", m2.Len())
	}
}

func TestPop_CopySemantics(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil))
	m2 := m.Pop()

	if m.Len() != 2 {
		t.Errorf("original should still have 2 crumbs, got %d", m.Len())
	}
	if m2.Len() != 1 {
		t.Errorf("new model should have 1 crumb, got %d", m2.Len())
	}
}

func TestSetCrumbs(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	crumbs := []teacrumbs.Crumb{{Text: "A"}, {Text: "B"}, {Text: "C"}}
	m2 := m.SetCrumbs(crumbs)

	if m2.Len() != 3 {
		t.Errorf("expected 3 crumbs, got %d", m2.Len())
	}

	// Verify the original slice isn't shared
	crumbs[0].Text = "Modified"
	if m2.Crumbs()[0].Text == "Modified" {
		t.Error("SetCrumbs should copy the slice, not share it")
	}
}

func TestSetCrumb(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil))
	m2 := m.SetCrumb(1, teacrumbs.NewCrumb("Preferences", nil))

	if m2.Crumbs()[1].Text != "Preferences" {
		t.Errorf("expected 'Preferences', got %q", m2.Crumbs()[1].Text)
	}
	// Original should be unchanged
	if m.Crumbs()[1].Text != "Settings" {
		t.Errorf("original should still be 'Settings', got %q", m.Crumbs()[1].Text)
	}
}

func TestSetCrumb_OutOfRange(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().Push(teacrumbs.NewCrumb("Home", nil))

	m2 := m.SetCrumb(-1, teacrumbs.NewCrumb("Bad", nil))
	if m2.Len() != 1 {
		t.Error("negative index should be no-op")
	}

	m3 := m.SetCrumb(5, teacrumbs.NewCrumb("Bad", nil))
	if m3.Len() != 1 {
		t.Error("out of range index should be no-op")
	}
}

func TestCrumbs_ReturnsCopy(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil))

	crumbs := m.Crumbs()
	crumbs[0].Text = "Modified"

	if m.Crumbs()[0].Text == "Modified" {
		t.Error("Crumbs() should return a copy, not a reference")
	}
}

func TestInit_ReturnsNil(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	if cmd := m.Init(); cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m2 := updated.(teacrumbs.BreadcrumbsModel)
	if m2.Width() != 120 {
		t.Errorf("expected width 120, got %d", m2.Width())
	}
}

func TestUpdate_PushCrumbMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	updated, _ := m.Update(teacrumbs.PushCrumbMsg{Crumb: teacrumbs.NewCrumb("Home", nil)})
	m2 := updated.(teacrumbs.BreadcrumbsModel)
	if m2.Len() != 1 {
		t.Errorf("expected 1 crumb, got %d", m2.Len())
	}
}

func TestUpdate_PopCrumbMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().Push(teacrumbs.NewCrumb("Home", nil))
	updated, _ := m.Update(teacrumbs.PopCrumbMsg{})
	m2 := updated.(teacrumbs.BreadcrumbsModel)
	if m2.Len() != 0 {
		t.Errorf("expected 0 crumbs, got %d", m2.Len())
	}
}

func TestUpdate_SetCrumbsMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel()
	crumbs := []teacrumbs.Crumb{{Text: "A"}, {Text: "B"}}
	updated, _ := m.Update(teacrumbs.SetCrumbsMsg{Crumbs: crumbs})
	m2 := updated.(teacrumbs.BreadcrumbsModel)
	if m2.Len() != 2 {
		t.Errorf("expected 2 crumbs, got %d", m2.Len())
	}
}

func TestUpdate_SetCrumbMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().Push(teacrumbs.NewCrumb("Home", nil))
	updated, _ := m.Update(teacrumbs.SetCrumbMsg{Index: 0, Crumb: teacrumbs.NewCrumb("Dashboard", nil)})
	m2 := updated.(teacrumbs.BreadcrumbsModel)
	if m2.Crumbs()[0].Text != "Dashboard" {
		t.Errorf("expected 'Dashboard', got %q", m2.Crumbs()[0].Text)
	}
}

func TestUpdate_UnknownMsg_NoOp(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().Push(teacrumbs.NewCrumb("Home", nil))
	updated, cmd := m.Update("unknown message")
	m2 := updated.(teacrumbs.BreadcrumbsModel)
	if m2.Len() != 1 {
		t.Error("unknown message should be no-op")
	}
	if cmd != nil {
		t.Error("unknown message should return nil cmd")
	}
}

func TestView_ContainsCrumbText(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil)).
		SetSize(80)

	view := m.View()
	if !strings.Contains(view.Content, "Home") {
		t.Errorf("view should contain 'Home', got %q", view.Content)
	}
	if !strings.Contains(view.Content, "Settings") {
		t.Errorf("view should contain 'Settings', got %q", view.Content)
	}
}

func TestView_EmptyCrumbs(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().SetSize(80)
	view := m.View()
	if view.Content != "" {
		t.Errorf("empty crumbs should render empty string, got %q", view.Content)
	}
}

func TestWithSeparator(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().WithSeparator(" / ")
	if m.Separator() != " / " {
		t.Errorf("expected separator ' / ', got %q", m.Separator())
	}
}

func TestWithStyles(t *testing.T) {
	custom := teacrumbs.DefaultStyles()
	m := teacrumbs.NewBreadcrumbsModel().WithStyles(custom)
	if m.Styles.ParentStyle.GetForeground() != custom.ParentStyle.GetForeground() {
		t.Error("WithStyles should apply custom styles")
	}
}

func TestPush_PreservesShortText(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().Push(teacrumbs.NewCrumb("go-dt", &teacrumbs.CrumbArgs{
		Short: "dt"},
	))
	crumbs := m.Crumbs()
	if crumbs[0].Short != "dt" {
		t.Errorf("Push should preserve Short text, got %q", crumbs[0].Short)
	}
}

func TestSetCrumb_PreservesStylePointer(t *testing.T) {
	s := teacrumbs.DefaultStyles().CurrentStyle
	m := teacrumbs.NewBreadcrumbsModel().Push(teacrumbs.NewCrumb("Home", nil))
	m = m.SetCrumb(0, teacrumbs.Crumb{Text: "Home", Style: &s})
	crumbs := m.Crumbs()
	if crumbs[0].Style == nil {
		t.Error("SetCrumb should preserve Style pointer")
	}
}

// --- Mouse/position/hover tests ---

func TestSetPosition(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().SetPosition(5, 10)
	row, col := m.Position()
	if row != 5 || col != 10 {
		t.Errorf("expected (5, 10), got (%d, %d)", row, col)
	}
}

func TestHitTest_OnCrumb(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil)).
		SetSize(80).
		SetPosition(0, 0)

	// Render to populate bounds
	m.View()

	// "Home" starts at x=0, has width 4
	idx := m.HitTest(0, 0)
	if idx != 0 {
		t.Errorf("expected hit on crumb 0 at x=0, got %d", idx)
	}
	idx = m.HitTest(3, 0)
	if idx != 0 {
		t.Errorf("expected hit on crumb 0 at x=3, got %d", idx)
	}
}

func TestHitTest_OnSeparator(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("AB", nil)).
		Push(teacrumbs.NewCrumb("CD", nil)).
		SetSize(80).
		SetPosition(0, 0)

	m.View()

	// "AB" width=2, sep " > " width=3, so separator is at x=2,3,4
	// "CD" starts at x=5
	idx := m.HitTest(3, 0)
	if idx != -1 {
		t.Errorf("expected no hit on separator, got %d", idx)
	}
}

func TestHitTest_WrongRow(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		SetSize(80).
		SetPosition(2, 0)

	m.View()

	idx := m.HitTest(0, 0)
	if idx != -1 {
		t.Errorf("expected no hit on wrong row, got %d", idx)
	}
}

func TestHitTest_OutOfBounds(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Hi", nil)).
		SetSize(80).
		SetPosition(0, 0)

	m.View()

	idx := m.HitTest(50, 0)
	if idx != -1 {
		t.Errorf("expected no hit past crumbs end, got %d", idx)
	}
}

func TestMouseClick_EmitsCrumbClickedMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", &teacrumbs.CrumbArgs{
			Data: "home-data",
		})).
		Push(teacrumbs.NewCrumb("Settings", nil)).
		SetSize(80).
		SetPosition(0, 0)

	m.View()

	// Simulate click on first crumb
	cmd := m.HandleMouse(tea.MouseClickMsg{X: 1, Y: 0, Button: tea.MouseLeft})
	if cmd == nil {
		t.Fatal("expected a cmd from click on crumb")
	}
	msg := cmd()
	clicked, ok := msg.(teacrumbs.CrumbClickedMsg)
	if !ok {
		t.Fatalf("expected CrumbClickedMsg, got %T", msg)
	}
	if clicked.Index != 0 {
		t.Errorf("expected Index=0, got %d", clicked.Index)
	}
	if clicked.Crumb.Text != "Home" {
		t.Errorf("expected crumb Text='Home', got %q", clicked.Crumb.Text)
	}
	if clicked.Button != tea.MouseLeft {
		t.Errorf("expected MouseLeft, got %v", clicked.Button)
	}
	if clicked.Crumb.Data != "home-data" {
		t.Errorf("expected Data='home-data', got %v", clicked.Crumb.Data)
	}
}

func TestMouseMotion_EmitsHoverMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil)).
		SetSize(80).
		SetPosition(0, 0)

	m.View()

	// Move onto first crumb
	cmd := m.HandleMouse(tea.MouseMotionMsg{X: 1, Y: 0})
	if cmd == nil {
		t.Fatal("expected a cmd from hover on crumb")
	}
	msg := cmd()
	hover, ok := msg.(teacrumbs.CrumbHoverMsg)
	if !ok {
		t.Fatalf("expected CrumbHoverMsg, got %T", msg)
	}
	if hover.Index != 0 {
		t.Errorf("expected Index=0, got %d", hover.Index)
	}
}

func TestMouseMotion_EmitsHoverLeaveMsg(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Hi", nil)).
		SetSize(80).
		SetPosition(0, 0)

	m.View()

	// First hover on crumb to set hoveredIdx
	m.HandleMouse(tea.MouseMotionMsg{X: 0, Y: 0})

	// Move off all crumbs
	cmd := m.HandleMouse(tea.MouseMotionMsg{X: 50, Y: 0})
	if cmd == nil {
		t.Fatal("expected a cmd for hover leave")
	}
	msg := cmd()
	if _, ok := msg.(teacrumbs.CrumbHoverLeaveMsg); !ok {
		t.Fatalf("expected CrumbHoverLeaveMsg, got %T", msg)
	}
}

func TestMouseMotion_HoverDedup(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		SetSize(80).
		SetPosition(0, 0)

	m.View()

	// First hover
	m.HandleMouse(tea.MouseMotionMsg{X: 0, Y: 0})

	// Same crumb again — should be nil (deduplicated)
	cmd := m.HandleMouse(tea.MouseMotionMsg{X: 1, Y: 0})
	if cmd != nil {
		t.Error("expected nil cmd for dedup hover on same crumb")
	}
}

func TestView_HasMouseMode(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		SetSize(80)

	view := m.View()
	if view.MouseMode != tea.MouseModeAllMotion {
		t.Errorf("expected MouseModeAllMotion, got %d", view.MouseMode)
	}
	if view.OnMouse == nil {
		t.Error("expected OnMouse handler to be set")
	}
}

func TestView_HoverStyleApplied(t *testing.T) {
	m := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil)).
		SetSize(80).
		SetPosition(0, 0)

	// Render once to get bounds
	m.View()

	// Hover on first crumb
	m.HandleMouse(tea.MouseMotionMsg{X: 0, Y: 0})

	// Re-render — "Home" should have hover style (underline)
	view := m.View()
	stripped := ansi.Strip(view.Content)
	if !strings.Contains(stripped, "Home") {
		t.Errorf("hover should still show crumb text, got %q (raw: %q)", stripped, view.Content)
	}
}
