// Source: site/src/content/docs/cookbook/grid-with-filtering.mdx:150,163,178
package main_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

// Column key constants mirroring those in main.go (inaccessible from main_test).
const (
	testColName  = "name"
	testColLang  = "language"
	testColStars = "stars"
)

// TestCompile_FilteredColumnDefinitions verifies that Column.WithFiltered(bool)
// is callable and that columns can be constructed with filtering flags.
// Source line 150.
func TestCompile_FilteredColumnDefinitions(t *testing.T) {
	columns := []teagrid.Column{
		teagrid.NewColumn(testColName, "Project", 20).WithFiltered(true),
		teagrid.NewColumn(testColLang, "Language", 15).WithFiltered(true),
		teagrid.NewColumn(testColStars, "Stars", 10), // not filterable
	}
	_ = columns
}

// TestCompile_GridWithFilterEnabled verifies that GridModel.WithFiltered(bool)
// is callable on a newly created grid model.
// Source line 163.
func TestCompile_GridWithFilterEnabled(t *testing.T) {
	columns := []teagrid.Column{
		teagrid.NewColumn(testColName, "Project", 20),
	}
	m := teagrid.NewGridModel(columns).WithFiltered(true)
	_ = m
}

// TestCompile_FilterFocusCheck verifies that IsFilterInputFocused() exists on
// GridModel and returns a bool usable in a conditional.
// Source line 178.
func TestCompile_FilterFocusCheck(t *testing.T) {
	columns := []teagrid.Column{
		teagrid.NewColumn(testColName, "Project", 20),
	}
	m := teagrid.NewGridModel(columns).WithFiltered(true)

	var cmds []tea.Cmd
	switch "ctrl+c" {
	case "ctrl+c", "q":
		if !m.IsFilterInputFocused() {
			cmds = append(cmds, tea.Quit)
		}
	}
	_ = cmds
}
