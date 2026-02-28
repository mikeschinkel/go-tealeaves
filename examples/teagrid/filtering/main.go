package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	columnKeyTitle       = "title"
	columnKeyAuthor      = "author"
	columnKeyDescription = "description"
)

type model struct {
	table teagrid.Model
}

func newModel() model {
	columns := []teagrid.Column{
		teagrid.NewColumn(columnKeyTitle, "Title", 13).WithFiltered(true),
		teagrid.NewColumn(columnKeyAuthor, "Author", 13).WithFiltered(true),
		teagrid.NewColumn(columnKeyDescription, "Description", 50),
	}
	return model{
		table: teagrid.
			New(columns).
			Filtered(true).
			Focused(true).
			WithPageSize(10).
			SelectableRows(true).
			WithRows([]teagrid.Row{
				teagrid.NewRow(teagrid.RowData{
					columnKeyTitle:       "Computer Systems : A Programmer's Perspective",
					columnKeyAuthor:      "Randal E. Bryant, David R. O'Hallaron / Prentice Hall",
					columnKeyDescription: "This book explains the important and enduring concepts underlying all computer...",
				}),
				teagrid.NewRow(teagrid.RowData{
					columnKeyTitle:       "Effective Java : 3rd Edition",
					columnKeyAuthor:      "Joshua Bloch",
					columnKeyDescription: "The Definitive Guide to Java Platform Best Practices-Updated for Java 9 Java ...",
				}),
				teagrid.NewRow(teagrid.RowData{
					columnKeyTitle:       "Structure and Interpretation of Computer Programs - 2nd Edition (MIT)",
					columnKeyAuthor:      "Harold Abelson, Gerald Jay Sussman",
					columnKeyDescription: "Structure and Interpretation of Computer Programs has had a dramatic impact on...",
				}),
				teagrid.NewRow(teagrid.RowData{
					columnKeyTitle:       "Game Programming Patterns",
					columnKeyAuthor:      "Robert Nystrom / Genever Benning",
					columnKeyDescription: "The biggest challenge facing many game programmers is completing their game. M...",
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

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.table.GetIsFilterInputFocused() {
				cmds = append(cmds, tea.Quit)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var body strings.Builder

	body.WriteString("A filtered simple default table\n")
	body.WriteString("Currently filter by Title and Author, press / + letters to start filtering, and escape to clear filter.\n")
	body.WriteString("Press q or ctrl+c to quit\n\n")
	body.WriteString(m.table.View())

	return body.String()
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
