// Source: site/src/content/docs/components/drilldown-view.mdx:22#9a2374ab,157#9d797105
package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teatree"
)

type myData struct{ name string }

// TestCompile_DrilldownQuickExample verifies the quick example from drilldown-view.mdx.
func TestCompile_DrilldownQuickExample(t *testing.T) {
	root := teatree.NewNode("myapp", "My App", myData{"myapp"})
	libA := teatree.NewNode("libA", "Library A", myData{"libA"})
	libB := teatree.NewNode("libB", "Library B", myData{"libB"})
	libC := teatree.NewNode("libC", "Library C", myData{"libC"})

	libA.SetChildren([]*teatree.Node[myData]{libC})
	root.SetChildren([]*teatree.Node[myData]{libA, libB})

	viewer := teatree.NewDrillDownModel(root, teatree.DrillDownArgs[myData]{
		SelectorFunc: func(parent *teatree.Node[myData], children []*teatree.Node[myData]) (*teatree.Node[myData], error) {
			return children[0], nil
		},
		Prompt: "Dependency path:",
	})

	viewer, err := viewer.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	_ = viewer
}

// TestCompile_SelectorFunc verifies the SelectorFunc example from drilldown-view.mdx.
func TestCompile_SelectorFunc(t *testing.T) {
	var countDepth func(n *teatree.Node[myData]) int
	countDepth = func(n *teatree.Node[myData]) int {
		children := n.Children()
		if len(children) == 0 {
			return 0
		}
		max := 0
		for _, c := range children {
			if d := countDepth(c); d > max {
				max = d
			}
		}
		return max + 1
	}

	selector := func(parent *teatree.Node[myData], children []*teatree.Node[myData]) (*teatree.Node[myData], error) {
		best := children[0]
		bestDepth := countDepth(best)
		for _, child := range children[1:] {
			if d := countDepth(child); d > bestDepth {
				best = child
				bestDepth = d
			}
		}
		return best, nil
	}

	root := teatree.NewNode("root", "Root", myData{"root"})
	child := teatree.NewNode("child", "Child", myData{"child"})
	root.SetChildren([]*teatree.Node[myData]{child})

	viewer := teatree.NewDrillDownModel(root, teatree.DrillDownArgs[myData]{
		SelectorFunc: selector,
	})
	_, err := viewer.Initialize()
	if err != nil {
		t.Fatal(err)
	}
}
