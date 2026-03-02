package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teadd"
	"github.com/mikeschinkel/go-tealeaves/teadep"
)

type model struct {
	pathViewer teadep.PathViewerModel
	width      int
	height     int
}

func main() {
	// Ensure that term.GetSize() is initialized before continuing.
	// This is needed in GoLand terminal for debugging, but is not harmful if not needed.
	teadd.EnsureTermGetSize(os.Stdout.Fd())

	// Build dependency tree from generated code
	root := ExampleTree()

	// Analyze tree to compute metadata
	metadata := analyzeTree(root, kindDeterminer)

	// Create selector that captures metadata in closure
	selector := func(parent *teadep.Tree, children []*teadep.Tree) (*teadep.Tree, error) {
		return selectBestChild(parent, children, metadata)
	}

	// Create path viewer with selector (metadata captured in closure)
	pathViewer := teadep.NewPathViewer(root, teadep.PathViewerArgs{
		SelectorFunc: selector,
		Prompt:       "Select a Commit Target:",
	})

	// Initialize the path viewer (validates and builds initial path)
	pathViewer, err := pathViewer.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing path viewer: %v\n", err)
		os.Exit(1)
	}

	m := model{
		pathViewer: pathViewer,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var pathViewer tea.Model

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Let pathViewer handle the size update via its Update method
		pathViewer, cmd = m.pathViewer.Update(msg)
		m.pathViewer = pathViewer.(teadep.PathViewerModel)
		return m, cmd

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case teadep.SelectNodeMsg:
		fmt.Printf("Selected: %s\n", msg.Tree.Node.DisplayName())
		return m, tea.Quit

	case teadep.ChangeNodeMsg:
		// Alternative was selected, path updated automatically
		return m, nil
	}

	// Delegate to path viewer
	pathViewer, cmd = m.pathViewer.Update(msg)
	m.pathViewer = pathViewer.(teadep.PathViewerModel)
	return m, cmd
}

func (m model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}

	content := m.pathViewer.View().Content
	content += "\n\n"
	content += "↑↓: navigate | Space/→: alternatives | Enter: select leaf | q: quit"

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// selectBestChild chooses the best child using metadata
// Priority: Executables > Libraries > Tests
// Then by reference count (higher better)
// Then by max depth (longer better)
func selectBestChild(parent *teadep.Tree, children []*teadep.Tree, meta map[*teadep.Tree]*nodeMeta) (*teadep.Tree, error) {
	var best *teadep.Tree
	var bestMeta *nodeMeta
	var childMeta *nodeMeta
	var child *teadep.Tree

	if len(children) == 0 {
		return nil, fmt.Errorf("no children")
	}

	best = children[0]
	bestMeta = meta[best]

	for _, child = range children[1:] {
		childMeta = meta[child]

		// Priority 1: Module kind (Exe > Lib > Test > Unspecified)
		if childMeta.ModuleKindSet && bestMeta.ModuleKindSet {
			childKind := childMeta.ModuleKind
			bestKind := bestMeta.ModuleKind

			// Convert testKind to lowest priority (-1)
			if childKind == kindTest {
				childKind = -1
			}
			if bestKind == kindTest {
				bestKind = -1
			}

			if childKind > bestKind {
				best = child
				bestMeta = childMeta
				continue
			}
			if childKind < bestKind {
				continue
			}
		}

		// Priority 2: Reference count (higher better)
		if childMeta.ReferenceCount > bestMeta.ReferenceCount {
			best = child
			bestMeta = childMeta
			continue
		}
		if childMeta.ReferenceCount < bestMeta.ReferenceCount {
			continue
		}

		// Priority 3: Max depth (longer better)
		if childMeta.MaxDepth > bestMeta.MaxDepth {
			best = child
			bestMeta = childMeta
		}
	}

	return best, nil
}

// kindDeterminer determines module kind based on the node type
func kindDeterminer(node teadep.Node) (kind int, ok bool) {
	var en *exampleNode
	var isExample bool

	en, isExample = node.(*exampleNode)
	if !isExample {
		goto end
	}
	if !en.KindSet {
		goto end
	}

	kind = en.Kind
	ok = true
end:
	return kind, ok
}
