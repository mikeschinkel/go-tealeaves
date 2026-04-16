package teagrid

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
)

func newTestGrid(width, height int) GridModel {
	cols := []Column{
		NewColumn("name", "Name", 20),
		NewColumn("role", "Role", 15),
		NewFlexColumn("email", "Email", 1),
	}
	rows := []Row{
		NewRow(RowData{"name": "Alice Johnson", "role": "Engineer", "email": "alice@example.com"}),
		NewRow(RowData{"name": "Bob Smith", "role": "Designer", "email": "bob@example.com"}),
		NewRow(RowData{"name": "Carol White", "role": "Manager", "email": "carol@example.com"}),
		NewRow(RowData{"name": "Dave Brown", "role": "Engineer", "email": "dave@example.com"}),
		NewRow(RowData{"name": "Eve Davis", "role": "Analyst", "email": "eve@example.com"}),
	}
	return NewGridModel(cols).WithSize(width, height).WithRows(rows)
}

func TestGridModel_Golden_80x24(t *testing.T) {
	m := newTestGrid(80, 24)
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}

func TestGridModel_Golden_120x40(t *testing.T) {
	m := newTestGrid(120, 40)
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}
