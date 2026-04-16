// Source: site/src/content/docs/cookbook/tree-with-statusbar.mdx:25#c2568167
package main

import (
	"fmt"
	"os"
	"strings"

	key "charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	dt "github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
	"github.com/mikeschinkel/go-tealeaves/teatree"
)

// headerLines is the number of rows used by the title and spacing above the tree.
const headerLines = 3

// statusBarLines is the number of rows used by the status bar at the bottom.
const statusBarLines = 1

type model struct {
	tree      teatree.TreeModel[teatree.File]
	statusBar teastatus.StatusBarModel
	width     int
	height    int
}

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

func main() {
	// Build a sample file tree
	files := []*teatree.File{
		teatree.NewFile(dt.RelFilepath("cmd/main.go"), nil),
		teatree.NewFile(dt.RelFilepath("cmd/version.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/handler/auth.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/handler/user.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/model/user.go"), nil),
		teatree.NewFile(dt.RelFilepath("pkg/config/config.go"), nil),
		teatree.NewFile(dt.RelFilepath("go.mod"), nil),
		teatree.NewFile(dt.RelFilepath("README.md"), nil),
	}

	nodes := teatree.BuildFileTree(files, teatree.BuildFileTreeArgs{
		RootPath: dt.PathSegment("myproject"),
	})

	provider := teatree.NewCompactNodeProvider[teatree.File](teatree.TriangleExpanderControls)
	tree := teatree.NewTree(nodes, &teatree.TreeArgs[teatree.File]{
		NodeProvider:     provider,
		ExpanderControls: &teatree.TriangleExpanderControls,
	})

	// Expand root nodes so the tree is visible on launch
	for _, node := range tree.Nodes() {
		node.Expand()
	}

	treeModel := teatree.NewTreeModel(tree, 20)

	// Build the status bar with key menus
	helpBinding := key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "Help"))
	expandBinding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "Expand all"))
	collapseBinding := key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "Collapse all"))
	quitBinding := key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit"))

	sb := teastatus.NewStatusBarModel().
		SetMenuItems([]teastatus.MenuItem{
			teastatus.NewMenuItem(helpBinding, &teastatus.MenuItemOpts{Label: "Help"}),
			teastatus.NewMenuItem(expandBinding, &teastatus.MenuItemOpts{Label: "Expand all"}),
			teastatus.NewMenuItem(collapseBinding, &teastatus.MenuItemOpts{Label: "Collapse all"}),
			teastatus.NewMenuItem(quitBinding, &teastatus.MenuItemOpts{Label: "Quit"}),
		}).
		SetIndicators([]teastatus.StatusIndicator{
			teastatus.NewStatusIndicator("Ready"),
		})

	m := model{
		tree:      treeModel,
		statusBar: sb,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return m.tree.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Tree gets all height except header and status bar
		treeHeight := msg.Height - headerLines - statusBarLines
		if treeHeight < 5 {
			treeHeight = 5
		}
		m.tree = m.tree.SetSize(msg.Width, treeHeight)

		// Status bar gets the full width
		m.statusBar = m.statusBar.SetSize(msg.Width)

		// Update the file count indicator
		focused := m.tree.FocusedNode()
		if focused != nil {
			m.statusBar = m.statusBar.SetIndicators([]teastatus.StatusIndicator{
				teastatus.NewStatusIndicator(focused.Name()),
			})
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "e":
			m.tree.Tree().ExpandAll()
			return m, nil
		case "c":
			m.tree.Tree().CollapseAll()
			return m, nil
		}
	}

	m.tree, cmd = m.tree.Update(msg)

	// Update status bar indicator with focused node name
	if focused := m.tree.FocusedNode(); focused != nil {
		m.statusBar = m.statusBar.SetIndicators([]teastatus.StatusIndicator{
			teastatus.NewStatusIndicator(focused.Name()),
		})
	}

	return m, cmd
}

func (m model) View() tea.View {
	var b strings.Builder

	// Header
	b.WriteString(titleStyle.Render("File Navigator"))
	b.WriteString("\n\n")

	// Tree view
	b.WriteString(m.tree.View().Content)
	b.WriteString("\n")

	// Pad to push the status bar to the bottom
	lines := strings.Split(b.String(), "\n")
	targetLines := m.height - statusBarLines
	for len(lines) < targetLines {
		lines = append(lines, "")
	}
	// Trim if content is too tall
	if len(lines) > targetLines {
		lines = lines[:targetLines]
	}

	// Append the status bar as the final row
	lines = append(lines, m.statusBar.View().Content)

	v := tea.NewView(strings.Join(lines, "\n"))
	v.AltScreen = true
	return v
}
