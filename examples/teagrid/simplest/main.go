package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	columnKeyName    = "name"
	columnKeyElement = "element"
)

type model struct {
	simpleTable teagrid.Model
}

func newModel() model {
	return model{
		simpleTable: teagrid.New([]teagrid.Column{
			teagrid.NewColumn(columnKeyName, "Name", 13),
			teagrid.NewColumn(columnKeyElement, "Element", 10),
		}).WithRows([]teagrid.Row{
			teagrid.NewRow(teagrid.RowData{
				columnKeyName:    "Pikachu",
				columnKeyElement: "Electric",
			}),
			teagrid.NewRow(teagrid.RowData{
				columnKeyName:    "Charmander",
				columnKeyElement: "Fire",
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
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var body strings.Builder

	body.WriteString("A very simple default table (non-interactive)\nPress q or ctrl+c to quit\n\n")
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
