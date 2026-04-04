package examples_test

import (
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teadiffr"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TestCompile_DiffRendererQuickExample verifies the quick example from diff-renderer.mdx.
func TestCompile_DiffRendererQuickExample(t *testing.T) {
	diffs := []teadiffr.FileDiff{
		{
			Path:   "cmd/main.go",
			Status: teadiffr.FileModified,
			Blocks: []teadiffr.CondensedBlock{
				{
					Type:         "added",
					LineCount:    2,
					ChangedLines: []string{"func newHandler() http.Handler {", "    return mux"},
				},
			},
		},
		{
			Path:   "internal/config.go",
			Status: teadiffr.FileNew,
			Blocks: []teadiffr.CondensedBlock{
				{Type: "added", LineCount: 10, ChangedLines: []string{"package internal"}},
			},
		},
	}

	renderer := teadiffr.NewTUIRenderer(nil)
	termWidth := 80
	lines := teadiffr.RenderFileDiffs(diffs, renderer, termWidth)
	output := strings.Join(lines, "\n")
	_ = output
}

// TestCompile_TUIRendererWithArgs verifies TUIRenderer with color args from diff-renderer.mdx.
func TestCompile_TUIRendererWithArgs(t *testing.T) {
	_ = teadiffr.NewTUIRenderer(&teadiffr.TUIRendererArgs{
		AddedColor:   lipgloss.Color("46"),
		DeletedColor: lipgloss.Color("196"),
	})
}

// TestCompile_ThemedRenderer verifies NewThemedTUIRenderer from diff-renderer.mdx.
func TestCompile_ThemedRenderer(t *testing.T) {
	theme := teautils.DefaultTheme()
	renderer := teadiffr.NewThemedTUIRenderer(theme)
	_ = renderer
}
