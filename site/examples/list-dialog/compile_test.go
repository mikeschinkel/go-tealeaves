package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

type myItem struct {
	id     string
	name   string
	active bool
}

func (m myItem) ID() string     { return m.id }
func (m myItem) Label() string  { return m.name }
func (m myItem) IsActive() bool { return m.active }

// TestCompile_ListDialogQuickExample verifies the quick example from list-dialog.mdx.
func TestCompile_ListDialogQuickExample(t *testing.T) {
	items := []myItem{
		{id: "p1", name: "Profile 1", active: true},
		{id: "p2", name: "Profile 2", active: false},
	}

	modal := teamodal.NewListModel[myItem](items, &teamodal.ListModelArgs{
		Title:      "Select a Profile",
		MaxVisible: 8,
		LabelWidth: 20,
	})

	modal = modal.Open()
	_ = modal.IsOpen()
}
