package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teahilite"
)

const goCode = `package main

func main() {}`

// TestCompile_SyntaxHighlightingQuickExample verifies the quick example from syntax-highlighting.mdx.
func TestCompile_SyntaxHighlightingQuickExample(t *testing.T) {
	// Package-level convenience (uses monokai + terminal256 defaults)
	highlighted, err := teahilite.Highlight(goCode, "go")
	if err != nil {
		t.Fatal(err)
	}
	_ = highlighted

	// Create a configured Highlighter instance
	h := teahilite.NewHighlighter(teahilite.HighlighterArgs{
		StyleName: "dracula",
	})
	highlighted, err = h.Highlight(goCode, "go")
	if err != nil {
		t.Fatal(err)
	}
	_ = highlighted

	lines, err := h.HighlightLines(goCode, "go")
	if err != nil {
		t.Fatal(err)
	}
	_ = lines
}

// TestCompile_HighlighterArgs verifies HighlighterArgs fields from syntax-highlighting.mdx.
func TestCompile_HighlighterArgs(t *testing.T) {
	h := teahilite.NewHighlighter(teahilite.HighlighterArgs{
		StyleName:     "monokai",
		FormatterName: "terminal256",
	})

	result, err := h.Highlight(goCode, "go")
	if err != nil {
		t.Fatal(err)
	}
	_ = result

	lines, err := h.HighlightLines(goCode, "go")
	if err != nil {
		t.Fatal(err)
	}
	_ = lines
}

// TestCompile_PackageLevelFunctions verifies package-level functions from syntax-highlighting.mdx.
func TestCompile_PackageLevelFunctions(t *testing.T) {
	result, err := teahilite.Highlight(goCode, "python")
	if err != nil {
		t.Fatal(err)
	}
	_ = result

	lines, err := teahilite.HighlightLines(goCode, "python")
	if err != nil {
		t.Fatal(err)
	}
	_ = lines
}

// TestCompile_DetectLanguage verifies DetectLanguage from syntax-highlighting.mdx.
func TestCompile_DetectLanguage(t *testing.T) {
	lang := teahilite.DetectLanguage("main.go")
	_ = lang

	lang = teahilite.DetectLanguage("schema.sql")
	_ = lang

	lang = teahilite.DetectLanguage("app.ts")
	_ = lang

	lang = teahilite.DetectLanguage("unknown.xyz")
	_ = lang
}

// TestCompile_DetectLanguageWithHighlight verifies combining DetectLanguage with Highlight.
func TestCompile_DetectLanguageWithHighlight(t *testing.T) {
	filePath := "main.go"
	h := teahilite.NewHighlighter(teahilite.HighlighterArgs{})
	lang := teahilite.DetectLanguage(filePath)
	result, err := h.Highlight(goCode, lang)
	if err != nil {
		t.Fatal(err)
	}
	_ = result
}

// TestCompile_DefaultConstants verifies default name constants from syntax-highlighting.mdx.
func TestCompile_DefaultConstants(t *testing.T) {
	_ = teahilite.DefaultStyleName
	_ = teahilite.DefaultFormatterName
}
