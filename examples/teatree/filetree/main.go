package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-tealeaves/teatree"
)

type model struct {
	tree     teatree.TreeModel[teatree.File]
	width    int
	height   int
	quitting bool
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func main() {
	// Build a sample file tree representing a Go project
	files := []*teatree.File{
		teatree.NewFile(dt.RelFilepath("cmd/main.go"), nil),
		teatree.NewFile(dt.RelFilepath("cmd/version.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/handler/auth.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/handler/user.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/model/user.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/model/session.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/store/database.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/store/cache.go"), nil),
		teatree.NewFile(dt.RelFilepath("pkg/config/config.go"), nil),
		teatree.NewFile(dt.RelFilepath("pkg/config/defaults.go"), nil),
		teatree.NewFile(dt.RelFilepath("pkg/logger/logger.go"), nil),
		teatree.NewFile(dt.RelFilepath("go.mod"), nil),
		teatree.NewFile(dt.RelFilepath("go.sum"), nil),
		teatree.NewFile(dt.RelFilepath("README.md"), nil),
		teatree.NewFile(dt.RelFilepath("Makefile"), nil),
	}

	// Build hierarchical tree from flat file list
	nodes := teatree.BuildFileTree(files, teatree.BuildFileTreeArgs{
		RootPath: dt.PathSegment("myproject"),
	})

	// Create tree with triangle expander controls
	provider := teatree.NewCompactNodeProvider[teatree.File](teatree.TriangleExpanderControls)
	tree := teatree.NewTree(nodes, &teatree.TreeArgs[teatree.File]{
		NodeProvider:     provider,
		ExpanderControls: &teatree.TriangleExpanderControls,
	})

	// Expand root nodes for initial view
	for _, node := range tree.Nodes() {
		node.Expand()
	}

	treeModel := teatree.NewTreeModel(tree, 20)

	m := model{
		tree: treeModel,
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
		treeHeight := msg.Height - 6
		if treeHeight < 5 {
			treeHeight = 5
		}
		m.tree = m.tree.SetSize(msg.Width, treeHeight)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "e":
			// Expand all
			m.tree.Tree().ExpandAll()
			return m, nil

		case "c":
			// Collapse all
			m.tree.Tree().CollapseAll()
			return m, nil
		}
	}

	m.tree, cmd = m.tree.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("teatree File Tree Example"))
	b.WriteString("\n\n")

	b.WriteString(m.tree.View().Content)
	b.WriteString("\n\n")

	// Show focused node info
	if focused := m.tree.FocusedNode(); focused != nil {
		b.WriteString(fmt.Sprintf("Focused: %s", focused.Name()))
	}
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑↓: navigate | →: expand | ←: collapse | Space: toggle | e: expand all | c: collapse all | q: quit"))

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}
