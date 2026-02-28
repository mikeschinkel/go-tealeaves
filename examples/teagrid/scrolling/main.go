package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	columnKeyID = "id"

	numCols = 100
	numRows = 10
	idWidth = 5

	colWidth = 3
	maxWidth = 30
)

type model struct {
	scrollableTable teagrid.Model
}

func colKey(colNum int) string {
	return fmt.Sprintf("%d", colNum)
}

func genRow(id int) teagrid.Row {
	data := teagrid.RowData{
		columnKeyID: fmt.Sprintf("ID %d", id),
	}

	for i := 0; i < numCols; i++ {
		data[colKey(i)] = colWidth
	}

	return teagrid.NewRow(data)
}

func newModel() model {
	rows := make([]teagrid.Row, 0, numRows)
	for i := 0; i < numRows; i++ {
		rows = append(rows, genRow(i))
	}

	cols := []teagrid.Column{
		teagrid.NewColumn(columnKeyID, "ID", idWidth),
	}

	for i := 0; i < numCols; i++ {
		cols = append(cols, teagrid.NewColumn(colKey(i), colKey(i+1), colWidth))
	}

	t := teagrid.New(cols).
		WithRows(rows).
		WithMaxTotalWidth(maxWidth).
		WithHorizontalFreezeColumnCount(1).
		WithStaticFooter("A footer").
		Focused(true)

	return model{
		scrollableTable: t,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.scrollableTable, cmd = m.scrollableTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var body strings.Builder

	body.WriteString("A scrollable table\n")
	body.WriteString("Press shift+left or shift+right to scroll\n")
	body.WriteString("Press q or ctrl+c to quit\n\n")
	body.WriteString(m.scrollableTable.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
