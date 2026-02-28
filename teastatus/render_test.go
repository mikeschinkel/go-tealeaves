package teastatus

import (
	"strings"
	"testing"
)

func TestRenderMenuLine(t *testing.T) {
	items := []MenuItem{
		{Key: "?", Label: "Help"},
		{Key: "q", Label: "Quit"},
	}
	styles := DefaultStyles()
	result := RenderMenuLine(items, styles)

	if !strings.Contains(result, "?") {
		t.Error("expected result to contain key '?'")
	}
	if !strings.Contains(result, "Help") {
		t.Error("expected result to contain label 'Help'")
	}
	if !strings.Contains(result, "q") {
		t.Error("expected result to contain key 'q'")
	}
	if !strings.Contains(result, "Quit") {
		t.Error("expected result to contain label 'Quit'")
	}
}

func TestRenderMenuLine_Empty(t *testing.T) {
	styles := DefaultStyles()
	result := RenderMenuLine([]MenuItem{}, styles)
	if result != "" {
		t.Errorf("expected empty output for empty items, got %q", result)
	}
}
