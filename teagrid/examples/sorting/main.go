package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	columnKeyName = "name"
	columnKeyType = "type"
	columnKeyWins = "wins"
)

type model struct {
	simpleTable   teagrid.GridModel
	columnSortKey string
}

func newModel() model {
	return model{
		simpleTable: teagrid.NewGridModel([]teagrid.Column{
			teagrid.NewColumn(columnKeyName, "Name", 13),
			teagrid.NewColumn(columnKeyType, "Type", 13),
			teagrid.NewColumn(columnKeyWins, "Win %", 8).
				WithFormatString("%.1f%%"),
		}).Focused(true).WithRows([]teagrid.Row{
			teagrid.NewRow(teagrid.RowData{
				columnKeyName: "ピカピカ",
				columnKeyType: "Pikachu",
				columnKeyWins: 78.3,
			}),
			teagrid.NewRow(teagrid.RowData{
				columnKeyName: "Zapmouse",
				columnKeyType: "Pikachu",
				columnKeyWins: 3.3,
			}),
			teagrid.NewRow(teagrid.RowData{
				columnKeyName: "Burninator",
				columnKeyType: "Charmander",
				columnKeyWins: 32.1,
			}),
			teagrid.NewRow(teagrid.RowData{
				columnKeyName: "Alphonse",
				columnKeyType: "Pikachu",
				columnKeyWins: 13.8,
			}),
			teagrid.NewRow(teagrid.RowData{
				columnKeyName: "Trogdor",
				columnKeyType: "Charmander",
				columnKeyWins: 99.9,
			}),
			teagrid.NewRow(teagrid.RowData{
				columnKeyName: "Dihydrogen Monoxide",
				columnKeyType: "Squirtle",
				columnKeyWins: 31.348,
			}),
		}),
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

	m.simpleTable, cmd = m.simpleTable.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "n":
			m.columnSortKey = columnKeyName
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey)

		case "t":
			m.columnSortKey = columnKeyType
			m.simpleTable = m.simpleTable.SortByAsc(m.columnSortKey).ThenSortByDesc(columnKeyWins)

		case "w":
			m.columnSortKey = columnKeyWins
			m.simpleTable = m.simpleTable.SortByDesc(m.columnSortKey)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var body strings.Builder

	body.WriteString("A sorted simple default table\n")
	body.WriteString("Sort by (n)ame, (t)ype->wins combo, or (w)ins\n")
	body.WriteString("Currently sorting by: " + m.columnSortKey + "\n")
	body.WriteString("Press q or ctrl+c to quit\n\n")
	body.WriteString(m.simpleTable.View().Content)

	return tea.NewView(body.String())
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
