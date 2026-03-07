package teamodal

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestConfirmModel_WithTheme(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewOKModal("test", &ConfirmModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	}).WithTheme(theme)

	if m.titleStyle.GetForeground() == nil {
		t.Error("themed titleStyle has no foreground")
	}
	if m.messageStyle.GetForeground() == nil {
		t.Error("themed messageStyle has no foreground")
	}
	if m.focusedButtonStyle.GetBackground() == nil {
		t.Error("themed focusedButtonStyle has no background")
	}
}

func TestConfirmModel_WithTheme_PreservesContent(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewOKModal("hello", &ConfirmModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Test",
	}).WithTheme(theme)

	if m.message != "hello" {
		t.Errorf("message changed: %q", m.message)
	}
	if m.title != "Test" {
		t.Errorf("title changed: %q", m.title)
	}
}

func TestChoiceModel_WithTheme(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewChoiceModel(&ChoiceModelArgs{
		Title:   "Pick one",
		Message: "Choose wisely",
		Options: []ChoiceOption{
			{Label: "A"},
			{Label: "B"},
		},
	}).WithTheme(theme)

	if m.titleStyle.GetForeground() == nil {
		t.Error("themed titleStyle has no foreground")
	}
	if m.messageStyle.GetForeground() == nil {
		t.Error("themed messageStyle has no foreground")
	}
	if m.focusedButtonStyle.GetBackground() == nil {
		t.Error("themed focusedButtonStyle has no background")
	}
	if m.cancelKeyStyle.GetForeground() == nil {
		t.Error("themed cancelKeyStyle has no foreground")
	}
	if m.cancelTextStyle.GetForeground() == nil {
		t.Error("themed cancelTextStyle has no foreground")
	}
}

func TestChoiceModel_WithTheme_PreservesContent(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m := NewChoiceModel(&ChoiceModelArgs{
		Title:   "Title",
		Message: "Msg",
		Options: []ChoiceOption{
			{Label: "X"},
			{Label: "Y"},
		},
	}).WithTheme(theme)

	if m.title != "Title" {
		t.Errorf("title changed: %q", m.title)
	}
	if m.message != "Msg" {
		t.Errorf("message changed: %q", m.message)
	}
	if len(m.options) != 2 {
		t.Errorf("options count changed: %d", len(m.options))
	}
}

func TestListModel_WithTheme(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	items := []testThemeItem{
		{id: "1", label: "Alpha"},
		{id: "2", label: "Beta"},
	}
	m := NewListModel(items, &ListModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Test",
	}).WithTheme(theme)

	if m.titleStyle.GetForeground() == nil {
		t.Error("themed titleStyle has no foreground")
	}
	if m.itemStyle.GetForeground() == nil {
		t.Error("themed itemStyle has no foreground")
	}
	if m.selectedItemStyle.GetBackground() == nil {
		t.Error("themed selectedItemStyle has no background")
	}
	if m.activeItemStyle.GetForeground() == nil {
		t.Error("themed activeItemStyle has no foreground")
	}
	if m.editItemStyle.GetBackground() == nil {
		t.Error("themed editItemStyle has no background")
	}
}

func TestListModel_WithTheme_PreservesContent(t *testing.T) {
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	items := []testThemeItem{
		{id: "1", label: "Alpha"},
	}
	m := NewListModel(items, &ListModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "My List",
	}).WithTheme(theme)

	if m.title != "My List" {
		t.Errorf("title changed: %q", m.title)
	}
}

// testThemeItem implements ListItem for theme tests
type testThemeItem struct {
	id     string
	label  string
	active bool
}

func (ti testThemeItem) ID() string     { return ti.id }
func (ti testThemeItem) Label() string  { return ti.label }
func (ti testThemeItem) IsActive() bool { return ti.active }
