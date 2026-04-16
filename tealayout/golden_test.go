package tealayout

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
)

// goldenWidget implements SetSizer + ContentProvider for golden layout testing.
type goldenWidget struct {
	name          string
	width, height int
}

func (w *goldenWidget) SetSize(width, height int) { w.width = width; w.height = height }
func (w *goldenWidget) Content() string {
	if w.width == 0 || w.height == 0 {
		return ""
	}
	var lines []string
	header := fmt.Sprintf("[%s %dx%d]", w.name, w.width, w.height)
	lines = append(lines, header)
	for i := 1; i < w.height && i < 5; i++ {
		lines = append(lines, fmt.Sprintf("  %s line %d", w.name, i))
	}
	return strings.Join(lines, "\n")
}
func (w *goldenWidget) Focus() {}
func (w *goldenWidget) Blur()  {}

func TestTreeContentLayout_Golden_80x24(t *testing.T) {
	tree := &goldenWidget{name: "tree"}
	content := &goldenWidget{name: "content"}
	layout := NewTreeContentLayout(tree, content)
	layout.SetSize(80, 24)
	rendered, err := layout.Render()
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	output := ansi.Strip(rendered)
	golden.RequireEqual(t, []byte(output))
}

func TestTreeContentLayout_Golden_120x40(t *testing.T) {
	tree := &goldenWidget{name: "tree"}
	content := &goldenWidget{name: "content"}
	layout := NewTreeContentLayout(tree, content)
	layout.SetSize(120, 40)
	rendered, err := layout.Render()
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	output := ansi.Strip(rendered)
	golden.RequireEqual(t, []byte(output))
}
