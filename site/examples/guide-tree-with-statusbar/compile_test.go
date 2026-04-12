// Source: site/src/content/docs/cookbook/tree-with-statusbar.mdx:222,242
package main_test

import (
	"testing"

	dt "github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
	"github.com/mikeschinkel/go-tealeaves/teatree"
)

// Row-count constants mirroring those in main.go (inaccessible from main_test).
const (
	headerLines    = 3
	statusBarLines = 1
)

// buildTestTree constructs a minimal file tree for use in compile tests.
func buildTestTree() *teatree.Tree[teatree.File] {
	files := []*teatree.File{
		teatree.NewFile(dt.RelFilepath("cmd/main.go"), nil),
	}
	nodes := teatree.BuildFileTree(files, teatree.BuildFileTreeArgs{
		RootPath: dt.PathSegment("testproject"),
	})
	provider := teatree.NewCompactNodeProvider[teatree.File](teatree.TriangleExpanderControls)
	return teatree.NewTree(nodes, &teatree.TreeArgs[teatree.File]{
		NodeProvider:     provider,
		ExpanderControls: &teatree.TriangleExpanderControls,
	})
}

// TestCompile_TreeStatusBarSizing verifies that SetSize is callable on
// TreeModel and StatusBarModel with appropriate arguments.
// Source line 222.
func TestCompile_TreeStatusBarSizing(t *testing.T) {
	tree := buildTestTree()
	treeModel := teatree.NewTreeModel(tree, 20)
	statusBar := teastatus.NewStatusBarModel()

	msgWidth := 80
	msgHeight := 24

	treeHeight := msgHeight - headerLines - statusBarLines
	treeModel = treeModel.SetSize(msgWidth, treeHeight)
	statusBar = statusBar.SetSize(msgWidth)

	_ = treeModel
	_ = statusBar
}

// TestCompile_StatusBarFocusedNode verifies that FocusedNode() returns a node
// pointer and that SetIndicators accepts a slice of NewStatusIndicator values.
// Source line 242.
func TestCompile_StatusBarFocusedNode(t *testing.T) {
	tree := buildTestTree()
	treeModel := teatree.NewTreeModel(tree, 20)
	statusBar := teastatus.NewStatusBarModel()

	if focused := treeModel.FocusedNode(); focused != nil {
		statusBar = statusBar.SetIndicators([]teastatus.StatusIndicator{
			teastatus.NewStatusIndicator(focused.Name()),
		})
	}

	_ = statusBar
}
