package teacrumbs

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNewBreadcrumbsModel_Defaults(t *testing.T) {
	m := NewBreadcrumbsModel()
	if m.Len() != 0 {
		t.Errorf("expected empty trail, got %d", m.Len())
	}
	if m.separator != " > " {
		t.Errorf("expected default separator ' > ', got %q", m.separator)
	}
}

func TestPush(t *testing.T) {
	m := NewBreadcrumbsModel()
	m2 := m.Push(Crumb{Text: "Home"})
	if m2.Len() != 1 {
		t.Errorf("expected 1 crumb, got %d", m2.Len())
	}
	trail := m2.Trail()
	if trail[0].Text != "Home" {
		t.Errorf("expected 'Home', got %q", trail[0].Text)
	}
}

func TestPush_CopySemantics(t *testing.T) {
	m := NewBreadcrumbsModel().Push(Crumb{Text: "Home"})
	m2 := m.Push(Crumb{Text: "Settings"})

	if m.Len() != 1 {
		t.Errorf("original should still have 1 crumb, got %d", m.Len())
	}
	if m2.Len() != 2 {
		t.Errorf("new model should have 2 crumbs, got %d", m2.Len())
	}
}

func TestPop(t *testing.T) {
	m := NewBreadcrumbsModel().
		Push(Crumb{Text: "Home"}).
		Push(Crumb{Text: "Settings"})
	m2 := m.Pop()

	if m2.Len() != 1 {
		t.Errorf("expected 1 crumb after pop, got %d", m2.Len())
	}
	if m2.Trail()[0].Text != "Home" {
		t.Errorf("expected 'Home', got %q", m2.Trail()[0].Text)
	}
}

func TestPop_Empty(t *testing.T) {
	m := NewBreadcrumbsModel()
	m2 := m.Pop()
	if m2.Len() != 0 {
		t.Errorf("pop on empty should be no-op, got %d", m2.Len())
	}
}

func TestPop_CopySemantics(t *testing.T) {
	m := NewBreadcrumbsModel().
		Push(Crumb{Text: "Home"}).
		Push(Crumb{Text: "Settings"})
	m2 := m.Pop()

	if m.Len() != 2 {
		t.Errorf("original should still have 2 crumbs, got %d", m.Len())
	}
	if m2.Len() != 1 {
		t.Errorf("new model should have 1 crumb, got %d", m2.Len())
	}
}

func TestSetTrail(t *testing.T) {
	m := NewBreadcrumbsModel()
	trail := []Crumb{{Text: "A"}, {Text: "B"}, {Text: "C"}}
	m2 := m.SetTrail(trail)

	if m2.Len() != 3 {
		t.Errorf("expected 3 crumbs, got %d", m2.Len())
	}

	// Verify the original slice isn't shared
	trail[0].Text = "Modified"
	if m2.Trail()[0].Text == "Modified" {
		t.Error("SetTrail should copy the slice, not share it")
	}
}

func TestSetCrumb(t *testing.T) {
	m := NewBreadcrumbsModel().
		Push(Crumb{Text: "Home"}).
		Push(Crumb{Text: "Settings"})
	m2 := m.SetCrumb(1, Crumb{Text: "Preferences"})

	if m2.Trail()[1].Text != "Preferences" {
		t.Errorf("expected 'Preferences', got %q", m2.Trail()[1].Text)
	}
	// Original should be unchanged
	if m.Trail()[1].Text != "Settings" {
		t.Errorf("original should still be 'Settings', got %q", m.Trail()[1].Text)
	}
}

func TestSetCrumb_OutOfRange(t *testing.T) {
	m := NewBreadcrumbsModel().Push(Crumb{Text: "Home"})

	m2 := m.SetCrumb(-1, Crumb{Text: "Bad"})
	if m2.Len() != 1 {
		t.Error("negative index should be no-op")
	}

	m3 := m.SetCrumb(5, Crumb{Text: "Bad"})
	if m3.Len() != 1 {
		t.Error("out of range index should be no-op")
	}
}

func TestTrail_ReturnsCopy(t *testing.T) {
	m := NewBreadcrumbsModel().
		Push(Crumb{Text: "Home"}).
		Push(Crumb{Text: "Settings"})

	trail := m.Trail()
	trail[0].Text = "Modified"

	if m.Trail()[0].Text == "Modified" {
		t.Error("Trail() should return a copy, not a reference")
	}
}

func TestInit_ReturnsNil(t *testing.T) {
	m := NewBreadcrumbsModel()
	if cmd := m.Init(); cmd != nil {
		t.Error("Init() should return nil")
	}
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := NewBreadcrumbsModel()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m2 := updated.(BreadcrumbsModel)
	if m2.width != 120 {
		t.Errorf("expected width 120, got %d", m2.width)
	}
}

func TestUpdate_PushCrumbMsg(t *testing.T) {
	m := NewBreadcrumbsModel()
	updated, _ := m.Update(PushCrumbMsg{Crumb: Crumb{Text: "Home"}})
	m2 := updated.(BreadcrumbsModel)
	if m2.Len() != 1 {
		t.Errorf("expected 1 crumb, got %d", m2.Len())
	}
}

func TestUpdate_PopCrumbMsg(t *testing.T) {
	m := NewBreadcrumbsModel().Push(Crumb{Text: "Home"})
	updated, _ := m.Update(PopCrumbMsg{})
	m2 := updated.(BreadcrumbsModel)
	if m2.Len() != 0 {
		t.Errorf("expected 0 crumbs, got %d", m2.Len())
	}
}

func TestUpdate_SetTrailMsg(t *testing.T) {
	m := NewBreadcrumbsModel()
	trail := []Crumb{{Text: "A"}, {Text: "B"}}
	updated, _ := m.Update(SetTrailMsg{Trail: trail})
	m2 := updated.(BreadcrumbsModel)
	if m2.Len() != 2 {
		t.Errorf("expected 2 crumbs, got %d", m2.Len())
	}
}

func TestUpdate_SetCrumbMsg(t *testing.T) {
	m := NewBreadcrumbsModel().Push(Crumb{Text: "Home"})
	updated, _ := m.Update(SetCrumbMsg{Index: 0, Crumb: Crumb{Text: "Dashboard"}})
	m2 := updated.(BreadcrumbsModel)
	if m2.Trail()[0].Text != "Dashboard" {
		t.Errorf("expected 'Dashboard', got %q", m2.Trail()[0].Text)
	}
}

func TestUpdate_UnknownMsg_NoOp(t *testing.T) {
	m := NewBreadcrumbsModel().Push(Crumb{Text: "Home"})
	updated, cmd := m.Update("unknown message")
	m2 := updated.(BreadcrumbsModel)
	if m2.Len() != 1 {
		t.Error("unknown message should be no-op")
	}
	if cmd != nil {
		t.Error("unknown message should return nil cmd")
	}
}

func TestView_ContainsCrumbText(t *testing.T) {
	m := NewBreadcrumbsModel().
		Push(Crumb{Text: "Home"}).
		Push(Crumb{Text: "Settings"}).
		SetSize(80)

	view := m.View()
	if !strings.Contains(view.Content, "Home") {
		t.Errorf("view should contain 'Home', got %q", view.Content)
	}
	if !strings.Contains(view.Content, "Settings") {
		t.Errorf("view should contain 'Settings', got %q", view.Content)
	}
}

func TestView_EmptyTrail(t *testing.T) {
	m := NewBreadcrumbsModel().SetSize(80)
	view := m.View()
	if view.Content != "" {
		t.Errorf("empty trail should render empty string, got %q", view.Content)
	}
}

func TestWithSeparator(t *testing.T) {
	m := NewBreadcrumbsModel().WithSeparator(" / ")
	if m.separator != " / " {
		t.Errorf("expected separator ' / ', got %q", m.separator)
	}
}

func TestWithStyles(t *testing.T) {
	custom := DefaultStyles()
	m := NewBreadcrumbsModel().WithStyles(custom)
	if m.Styles.ParentStyle.GetForeground() != custom.ParentStyle.GetForeground() {
		t.Error("WithStyles should apply custom styles")
	}
}
