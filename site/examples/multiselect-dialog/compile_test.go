package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

type fileItem struct {
	id   string
	name string
}

func (f fileItem) ID() string    { return f.id }
func (f fileItem) Label() string { return f.name }

// TestCompile_MultiSelectQuickExample verifies the quick example from multiselect-dialog.mdx.
func TestCompile_MultiSelectQuickExample(t *testing.T) {
	modal := teamodal.NewMultiSelectModel([]fileItem{
		{id: "main.go", name: "main.go"},
		{id: "util.go", name: "util.go"},
		{id: "config.go", name: "config.go"},
	}, &teamodal.MultiSelectModelArgs{
		Title:   "Select files to delete",
		Message: "The following files will be permanently removed.",
		Footer:  "This action cannot be undone.",
		Buttons: []teamodal.MultiSelectButton{
			{Label: "Delete", Hotkey: 'd', ID: "delete"},
		},
	})

	modal, _ = modal.Open()
	_ = modal.IsOpen()
	_ = modal.Selected()
	_ = modal.Cursor()
}
