package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

const (
	colName   = "name"
	colRole   = "role"
	colStatus = "status"
)

type model struct {
	grid      teagrid.GridModel
	statusBar teastatus.StatusBarModel
	theme     teautils.Theme
	isDark    bool
	width     int
	height    int
}

func newModel() model {
	m := model{isDark: true}
	m = m.applyTheme()
	return m
}

func (m model) applyTheme() model {
	var sys teautils.SystemPalette
	if m.isDark {
		sys = teautils.DarkSystemPalette(nil)
	} else {
		sys = teautils.LightSystemPalette(nil)
	}
	m.theme = teautils.NewTheme(sys)

	// Build grid with theme
	m.grid = teagrid.NewGridModel([]teagrid.Column{
		teagrid.NewColumn(colName, "Name", 15),
		teagrid.NewColumn(colRole, "Role", 14),
		teagrid.NewColumn(colStatus, "Status", 10),
	}).WithRows(sampleRows()).
		WithTheme(m.theme).
		WithSelectableRows(true)

	if m.width > 0 {
		m.grid = m.grid.WithTargetWidth(m.width)
	}

	// Build status bar with theme
	themeLabel := "Dark"
	if !m.isDark {
		themeLabel = "Light"
	}
	m.statusBar = teastatus.NewStatusBarModel().
		WithTheme(m.theme).
		SetMenuItems([]teastatus.MenuItem{
			teastatus.NewMenuItemFromBinding(
				key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "Toggle theme")),
				"Theme",
			),
			teastatus.NewMenuItemFromBinding(
				key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit")),
				"Quit",
			),
		}).
		SetIndicators([]teastatus.StatusIndicator{
			{
				Text:  themeLabel,
				Style: m.theme.System.Accent.Foreground().Bold(true),
			},
		})

	if m.width > 0 {
		m.statusBar = m.statusBar.SetSize(m.width)
	}

	return m
}

func sampleRows() []teagrid.Row {
	return []teagrid.Row{
		teagrid.NewRow(teagrid.RowData{colName: "Alice", colRole: "Engineer", colStatus: "Active"}),
		teagrid.NewRow(teagrid.RowData{colName: "Bob", colRole: "Designer", colStatus: "Away"}),
		teagrid.NewRow(teagrid.RowData{colName: "Charlie", colRole: "Manager", colStatus: "Active"}),
		teagrid.NewRow(teagrid.RowData{colName: "Diana", colRole: "Engineer", colStatus: "Busy"}),
		teagrid.NewRow(teagrid.RowData{colName: "Eve", colRole: "Analyst", colStatus: "Active"}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.grid = m.grid.WithTargetWidth(m.width)
		m.statusBar = m.statusBar.SetSize(m.width)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		case "t":
			m.isDark = !m.isDark
			m = m.applyTheme()
			return m, nil
		}
	}

	m.grid, cmd = m.grid.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}

	var body strings.Builder

	title := m.theme.System.Accent.Foreground().Bold(true).
		Render("Theming Example")

	hint := m.theme.System.TextMuted.Render("Press [t] to toggle dark/light theme")

	body.WriteString(fmt.Sprintf("\n  %s  %s\n\n", title, hint))
	body.WriteString(m.grid.View().Content)

	// Fill remaining height to push status bar to bottom
	lines := strings.Split(body.String(), "\n")
	for len(lines) < m.height-1 {
		lines = append(lines, "")
	}
	lines = append(lines, m.statusBar.View().Content)

	v := tea.NewView(strings.Join(lines, "\n"))
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
