package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	columnKeyTitle       = "title"
	columnKeyAuthor      = "author"
	columnKeyDescription = "description"
)

type model struct {
	table teagrid.GridModel
}

type book struct {
	title       string
	author      string
	description string
}

var books = []book{
	{"Computer Systems : A Programmer's Perspective", "Randal E. Bryant", "This book explains the important and enduring concepts underlying all computer systems."},
	{"Effective Java : 3rd Edition", "Joshua Bloch", "The Definitive Guide to Java Platform Best Practices, Updated for Java 9."},
	{"Structure and Interpretation of Computer Programs", "Harold Abelson", "Structure and Interpretation of Computer Programs has had a dramatic impact on computer science curricula."},
	{"Game Programming Patterns", "Robert Nystrom", "The biggest challenge facing many game programmers is completing their game."},
	{"The Go Programming Language", "Alan Donovan", "A thorough guide to the Go language, following in the tradition of K&R C."},
	{"Clean Code", "Robert C. Martin", "A handbook of agile software craftsmanship with examples in Java."},
}

func newModel() model {
	columns := []teagrid.Column{
		teagrid.NewColumn(columnKeyTitle, "Title", 50).WithFiltered(true),
		teagrid.NewColumn(columnKeyAuthor, "Author", 20).WithFiltered(true),
		teagrid.NewFlexColumn(columnKeyDescription, "Description", 1),
	}

	rows := make([]teagrid.Row, len(books))
	for i, b := range books {
		rows[i] = teagrid.NewRow(teagrid.RowData{
			columnKeyTitle:       b.title,
			columnKeyAuthor:      b.author,
			columnKeyDescription: b.description,
		})
	}

	return model{
		table: teagrid.
			NewGridModel(columns).
			Filtered(true).
			Focused(true).
			WithPageSize(10).
			WithRows(rows),
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
	case tea.WindowSizeMsg:
		m.table = m.table.SetSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.table.GetIsFilterInputFocused() {
				cmds = append(cmds, tea.Quit)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var body strings.Builder

	body.WriteString("Filter by Title and Author: press / to start filtering, escape to clear.\n")
	body.WriteString("Press q or ctrl+c to quit\n\n")
	body.WriteString(m.table.View().Content)

	v := tea.NewView(body.String())
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
