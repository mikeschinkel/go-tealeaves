package teadiffview

import (
	"testing"
)

func TestRenderFileDiffs_Empty(t *testing.T) {
	renderer := NewTUIRenderer(nil)
	lines := RenderFileDiffs(nil, renderer, 80)
	if len(lines) != 0 {
		t.Errorf("expected 0 lines for nil input, got %d", len(lines))
	}
}

func TestRenderFileDiffs_SingleFile(t *testing.T) {
	renderer := NewTUIRenderer(nil)
	files := []FileDiff{
		{
			Path:   "main.go",
			Status: FileModified,
			Blocks: []CondensedBlock{
				{
					Type:          "added",
					LineCount:     2,
					ContextBefore: []string{"package main"},
					ChangedLines:  []string{"import \"fmt\"", "import \"os\""},
					ContextAfter:  []string{"func main() {"},
				},
			},
		},
	}

	lines := RenderFileDiffs(files, renderer, 80)

	if len(lines) == 0 {
		t.Fatal("expected non-empty output")
	}
	if len(lines) != 6 {
		t.Errorf("expected 6 lines, got %d", len(lines))
	}
}

func TestRenderFileDiffs_MultipleFiles(t *testing.T) {
	renderer := NewTUIRenderer(nil)
	files := []FileDiff{
		{Path: "a.go", Status: FileNew, Blocks: nil},
		{Path: "b.go", Status: FileDeleted, Blocks: nil},
	}

	lines := RenderFileDiffs(files, renderer, 80)

	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestRenderFileDiffs_Truncation(t *testing.T) {
	renderer := NewTUIRenderer(nil)
	files := []FileDiff{
		{
			Path:   "truncated.go",
			Status: FileModified,
			Blocks: []CondensedBlock{
				{
					Type:         "deleted",
					LineCount:    10,
					ChangedLines: []string{"old line"},
					IsTruncated:  true,
				},
			},
		},
	}

	lines := RenderFileDiffs(files, renderer, 80)

	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestNewTUIRenderer_Defaults(t *testing.T) {
	renderer := NewTUIRenderer(nil)
	if renderer.FileHeaderColor == nil {
		t.Error("expected non-nil FileHeaderColor")
	}
	if renderer.AddedColor == nil {
		t.Error("expected non-nil AddedColor")
	}
}

func TestNewTUIRenderer_CustomArgs(t *testing.T) {
	custom := lipglossColor("99")
	renderer := NewTUIRenderer(&TUIRendererArgs{
		AddedColor: custom,
	})
	if renderer.AddedColor != custom {
		t.Error("expected custom AddedColor")
	}
	if renderer.FileHeaderColor == nil {
		t.Error("expected default FileHeaderColor")
	}
}

// lipglossColor is a test helper to create a color value.
func lipglossColor(s string) ansiColor {
	return ansiColor(s)
}

// ansiColor implements color.Color for testing.
type ansiColor string

func (c ansiColor) RGBA() (r, g, b, a uint32) {
	return 0, 0, 0, 0xffff
}
