package teahilite

import (
	"strings"
	"testing"
)

func TestHighlighter_Highlight_ContainsANSI(t *testing.T) {
	h := NewHighlighter(HighlighterArgs{})
	code := `package main

func main() {
	println("hello")
}
`
	result, err := h.Highlight(code, "go")
	if err != nil {
		t.Fatalf("Highlight returned error: %v", err)
	}
	if !strings.Contains(result, "\x1b[") {
		t.Error("expected ANSI escape sequences in highlighted output")
	}
}

func TestHighlighter_HighlightLines_SplitsLines(t *testing.T) {
	h := NewHighlighter(HighlighterArgs{})
	code := "line1\nline2\nline3"
	lines, err := h.HighlightLines(code, "text")
	if err != nil {
		t.Fatalf("HighlightLines returned error: %v", err)
	}
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d", len(lines))
	}
}

func TestHighlighter_UnknownLanguageFallback(t *testing.T) {
	h := NewHighlighter(HighlighterArgs{})
	result, err := h.Highlight("hello world", "nonexistent_lang_xyz")
	if err != nil {
		t.Fatalf("Highlight with unknown language returned error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty output for unknown language")
	}
}

func TestHighlighter_CustomStyle(t *testing.T) {
	h := NewHighlighter(HighlighterArgs{StyleName: "github"})
	result, err := h.Highlight(`x = 1`, "python")
	if err != nil {
		t.Fatalf("Highlight with github style returned error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty output with github style")
	}
}

func TestHighlighter_UnknownStyleFallback(t *testing.T) {
	h := NewHighlighter(HighlighterArgs{StyleName: "nonexistent_style_xyz"})
	result, err := h.Highlight(`x = 1`, "python")
	if err != nil {
		t.Fatalf("Highlight with unknown style returned error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty output with fallback style")
	}
}

func TestConvenience_Highlight(t *testing.T) {
	result, err := Highlight(`fmt.Println("hi")`, "go")
	if err != nil {
		t.Fatalf("convenience Highlight returned error: %v", err)
	}
	if !strings.Contains(result, "\x1b[") {
		t.Error("expected ANSI escape sequences from convenience Highlight")
	}
}

func TestConvenience_HighlightLines(t *testing.T) {
	lines, err := HighlightLines("a\nb", "text")
	if err != nil {
		t.Fatalf("convenience HighlightLines returned error: %v", err)
	}
	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines, got %d", len(lines))
	}
}
