package main

import (
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	colName = "name"
	colLang = "language"
	colStars = "stars"
)

type project struct {
	name     string
	language string
	stars    int
}

var projects = []project{
	{"bubbletea", "Go", 28000},
	{"lipgloss", "Go", 8200},
	{"charm", "Go", 3100},
	{"glow", "Go", 16000},
	{"vhs", "Go", 15000},
	{"react", "JavaScript", 230000},
	{"vue", "JavaScript", 208000},
	{"svelte", "JavaScript", 80000},
	{"django", "Python", 81000},
	{"flask", "Python", 68000},
	{"fastapi", "Python", 78000},
	{"rails", "Ruby", 56000},
}

type model struct {
	table teagrid.GridModel
}

func newModel() model {
	columns := []teagrid.Column{
		teagrid.NewColumn(colName, "Project", 20).WithFiltered(true),
		teagrid.NewColumn(colLang, "Language", 15).WithFiltered(true),
		teagrid.NewColumn(colStars, "Stars", 10),
	}

	rows := make([]teagrid.Row, len(projects))
	for i, p := range projects {
		rows[i] = teagrid.NewRow(teagrid.RowData{
			colName:  p.name,
			colLang:  p.language,
			colStars: p.stars,
		})
	}

	return model{
		table: teagrid.
			NewGridModel(columns).
			WithFiltered(true).
			WithFocused(true).
			WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Left)).
			WithPageSize(15).
			WithRows(rows),
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Let the grid handle the message (navigation, filtering, sorting)
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table = m.table.WithTargetWidth(msg.Width)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Only quit if the filter input is NOT focused.
			// When filtering, "q" is a valid character to type.
			if !m.table.IsFilterInputFocused() {
				cmds = append(cmds, tea.Quit)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var b strings.Builder

	b.WriteString("Filterable Project Grid\n")
	b.WriteString("Press / to filter, Escape to clear, q to quit\n\n")
	b.WriteString(m.table.View().Content)

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
