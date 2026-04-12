// Source: site/src/content/docs/components/diff-viewer.mdx:115,154,308,322,339,360,382,395,406
package examples_test

import (
	"strings"
	"testing"

	"github.com/mikeschinkel/go-diffutils"
	"github.com/mikeschinkel/go-tealeaves/teacolor"
	"github.com/mikeschinkel/go-tealeaves/teadiff"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TestCompile_SplitDiffModelConstruction verifies SplitDiffModel construction from diff-viewer.mdx.
func TestCompile_SplitDiffModelConstruction(t *testing.T) {
	model := teadiff.NewSplitDiffModel(&teadiff.SplitDiffModelArgs{
		Width:  80,
		Height: 24,
	})
	_ = model.CursorIndex()
	_ = model.LineCount()
}

// TestCompile_RowAnnotations verifies row annotations from diff-viewer.mdx.
func TestCompile_RowAnnotations(t *testing.T) {
	model := teadiff.NewSplitDiffModel(&teadiff.SplitDiffModelArgs{
		Width:  80,
		Height: 24,
	})

	annotations := map[int]teadiff.RowAnnotation{
		42: {Char: '●', Color: teacolor.Coral},
		17: {Char: '!', Color: teacolor.Gold},
	}
	model = model.SetAnnotations(annotations)
	_ = model
}

// TestCompile_CommitGroupColors verifies CommitGroupColors usage from diff-viewer.mdx.
func TestCompile_CommitGroupColors(t *testing.T) {
	model := teadiff.NewSplitDiffModel(&teadiff.SplitDiffModelArgs{
		Width:  80,
		Height: 24,
	})

	lineGroups := map[int]int{
		0: 0,
		1: 1,
	}

	annotations := map[int]teadiff.RowAnnotation{}
	for lineIdx, groupIdx := range lineGroups {
		annotations[lineIdx] = teadiff.RowAnnotation{
			Char:  '●',
			Color: teadiff.CommitGroupColors[groupIdx%len(teadiff.CommitGroupColors)],
		}
	}
	model = model.SetAnnotations(annotations)
	_ = model
}

// TestCompile_TUIRenderer verifies TUIRenderer from diff-viewer.mdx.
func TestCompile_TUIRenderer(t *testing.T) {
	renderer := teadiff.NewTUIRenderer(nil)

	diffs := []teadiff.FileDiff{
		{
			Path:   "cmd/main.go",
			Status: teadiff.FileModified,
			Blocks: []teadiff.CondensedBlock{
				{Type: "added", LineCount: 3, ChangedLines: []string{"+func foo() {}"}},
			},
		},
	}

	termWidth := 80
	lines := teadiff.RenderFileDiffs(diffs, renderer, termWidth)
	output := strings.Join(lines, "\n")
	_ = output
}

// TestCompile_ThemedTUIRenderer verifies NewThemedTUIRenderer from diff-viewer.mdx.
func TestCompile_ThemedTUIRenderer(t *testing.T) {
	theme := teautils.DefaultTheme()
	renderer := teadiff.NewThemedTUIRenderer(theme)
	_ = renderer
}

// TestCompile_SplitDiffWithHighlight verifies SetContent and HighlightFunc from diff-viewer.mdx.
func TestCompile_SplitDiffWithHighlight(t *testing.T) {
	model := teadiff.NewSplitDiffModel(&teadiff.SplitDiffModelArgs{
		Width:  80,
		Height: 24,
		HighlightFunc: func(text, lang string) string {
			return text // identity; real usage would call teahilite.Highlight
		},
	})

	content, err := diffutils.DiffLines(
		[]string{"func main() {}"},
		[]string{"func main() { fmt.Println() }"},
		"main.go",
	)
	if err != nil {
		t.Fatal(err)
	}
	model = model.SetContent(content)
	_ = model
}

// TestCompile_PaneLineTypes verifies pane line types from diff-viewer.mdx.
func TestCompile_PaneLineTypes(t *testing.T) {
	textLine := teadiff.NewTextLine(1, "func main() {}")
	_ = textLine.LineNo()

	blockMarker := teadiff.NewBlockMarker(1, 5)
	_ = blockMarker.IsBlockStart()
	_ = blockMarker.LineNo()

	placeholder := teadiff.NewPlaceholderLine(1, 2)
	_ = placeholder.IsWithinHunk()
	_ = placeholder.LineNo()
}
