package examples_test

import (
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	dt "github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-tealeaves/teatree"
)

// TestCompile_TreeViewQuickExample verifies the quick example from tree-view.mdx.
func TestCompile_TreeViewQuickExample(t *testing.T) {
	files := []*teatree.File{
		teatree.NewFile(dt.RelFilepath("cmd/main.go"), nil),
		teatree.NewFile(dt.RelFilepath("internal/handler/auth.go"), nil),
		teatree.NewFile(dt.RelFilepath("pkg/config/config.go"), nil),
		teatree.NewFile(dt.RelFilepath("go.mod"), nil),
	}

	nodes := teatree.BuildFileTree(files, teatree.BuildFileTreeArgs{
		RootPath: dt.PathSegment("myproject"),
	})

	provider := teatree.NewCompactNodeProvider[teatree.File](teatree.TriangleExpanderControls)
	tree := teatree.NewTree(nodes, &teatree.TreeArgs[teatree.File]{
		NodeProvider:     provider,
		ExpanderControls: &teatree.TriangleExpanderControls,
	})

	treeModel := teatree.NewTreeModel(tree, 20)
	_ = treeModel
}

// TestCompile_ExpanderControls verifies expander control presets from tree-view.mdx.
func TestCompile_ExpanderControls(t *testing.T) {
	_ = teatree.TriangleExpanderControls
	_ = teatree.PlusExpanderControls
	_ = teatree.NoExpanderControls
}

// TestCompile_DefaultTreeKeyMap verifies DefaultTreeKeyMap from tree-view.mdx.
func TestCompile_DefaultTreeKeyMap(t *testing.T) {
	km := teatree.DefaultTreeKeyMap()
	_ = km.Up
	_ = km.Down
	_ = km.ExpandOrEnter
	_ = km.CollapseOrUp
	_ = km.Toggle
}

// TestCompile_CustomNodeProvider verifies custom NodeProvider from tree-view.mdx.
func TestCompile_CustomNodeProvider(t *testing.T) {
	type Package struct {
		Name     string
		Version  string
		Outdated bool
	}

	type PackageProvider struct {
		teatree.CompactNodeProvider[Package]
	}

	// Verify the CompactNodeProvider can be embedded
	pp := &PackageProvider{}
	_ = pp

	// Build a basic tree with custom type
	nodes := []*teatree.Node[Package]{
		teatree.NewNode("root", "root", Package{Name: "myapp", Version: "1.0.0"}),
	}
	tree := teatree.NewTree(nodes, &teatree.TreeArgs[Package]{})
	treeModel := teatree.NewTreeModel(tree, 20)
	_ = treeModel
}

// TestCompile_NewFile verifies NewFile from tree-view.mdx.
func TestCompile_NewFile(t *testing.T) {
	f := teatree.NewFile(dt.RelFilepath("cmd/main.go"), nil)
	_ = f
}

// TestCompile_TreeWithStyleOption verifies custom NodeProvider Style method from tree-view.mdx.
func TestCompile_TreeWithStyleOption(t *testing.T) {
	type Package struct {
		Name     string
		Version  string
		Outdated bool
	}

	// Build using the style method signature from MDX
	nodes := []*teatree.Node[Package]{
		teatree.NewNode("root", "root", Package{Name: "myapp", Version: "1.0.0", Outdated: true}),
	}

	type PackageProvider struct {
		teatree.CompactNodeProvider[Package]
	}

	var p PackageProvider
	// Verify Style signature matches: Style(node *Node[T], tree *Tree[T]) lipgloss.Style
	tree := teatree.NewTree(nodes, &teatree.TreeArgs[Package]{
		NodeProvider: &p,
	})

	// Access node data to verify it compiles
	for _, n := range nodes {
		style := p.Style(n, tree)
		if n.Data() != nil && n.Data().Outdated {
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
		}
		_ = style
	}
}
