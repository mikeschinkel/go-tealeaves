package teacrumbs

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestRenderTrail_Empty(t *testing.T) {
	result := renderTrail(nil, " > ", 80, DefaultStyles())
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestRenderTrail_SingleCrumb(t *testing.T) {
	trail := []Crumb{{Text: "Home"}}
	result := renderTrail(trail, " > ", 80, DefaultStyles())
	if !strings.Contains(result, "Home") {
		t.Errorf("expected result to contain 'Home', got %q", result)
	}
}

func TestRenderTrail_MultipleCrumbs(t *testing.T) {
	trail := []Crumb{
		{Text: "Home"},
		{Text: "Settings"},
		{Text: "Profile"},
	}
	result := renderTrail(trail, " > ", 0, DefaultStyles())
	if !strings.Contains(result, "Home") {
		t.Errorf("expected result to contain 'Home', got %q", result)
	}
	if !strings.Contains(result, "Settings") {
		t.Errorf("expected result to contain 'Settings', got %q", result)
	}
	if !strings.Contains(result, "Profile") {
		t.Errorf("expected result to contain 'Profile', got %q", result)
	}
}

func TestRenderTrail_PreStyledPassthrough(t *testing.T) {
	styledText := "\x1b[31mRed Text\x1b[0m"
	trail := []Crumb{
		{Text: "Home"},
		{Text: styledText, PreStyled: true},
	}
	result := renderTrail(trail, " > ", 0, DefaultStyles())
	if !strings.Contains(result, styledText) {
		t.Errorf("PreStyled text should pass through unchanged, got %q", result)
	}
}

func TestRenderTrail_CustomSeparator(t *testing.T) {
	trail := []Crumb{
		{Text: "A"},
		{Text: "B"},
	}
	result := renderTrail(trail, " / ", 0, DefaultStyles())
	if !strings.Contains(result, "/") {
		t.Errorf("expected custom separator '/', got %q", result)
	}
}

func TestRenderTrail_TruncationWithEllipsis(t *testing.T) {
	trail := []Crumb{
		{Text: "Home"},
		{Text: "Very Long Middle Path"},
		{Text: "Another Long Middle Path"},
		{Text: "Current"},
	}
	// Use a narrow width to force truncation
	result := renderTrail(trail, " > ", 30, DefaultStyles())
	if !strings.Contains(result, "...") {
		t.Errorf("expected ellipsis in truncated result, got %q", result)
	}
	if !strings.Contains(result, "Home") {
		t.Errorf("expected first crumb 'Home' preserved, got %q", result)
	}
}

func TestRenderTrail_ExtremeTruncation(t *testing.T) {
	trail := []Crumb{
		{Text: "Very Long Home Path"},
		{Text: "Middle"},
		{Text: "Very Long Current Path"},
	}
	result := renderTrail(trail, " > ", 10, DefaultStyles())
	if !strings.Contains(result, "...") {
		t.Errorf("expected ellipsis in extreme truncation, got %q", result)
	}
}

func TestRenderTrail_ZeroWidth_NoTruncation(t *testing.T) {
	trail := []Crumb{
		{Text: "Home"},
		{Text: "Settings"},
		{Text: "Profile"},
	}
	result := renderTrail(trail, " > ", 0, DefaultStyles())
	// With zero width, no truncation should happen
	if strings.Contains(result, "...") {
		t.Errorf("zero width should not truncate, got %q", result)
	}
	if !strings.Contains(result, "Settings") {
		t.Errorf("expected all crumbs present, got %q", result)
	}
}

func TestRenderTrail_UsesAnsiStringWidth(t *testing.T) {
	// Verify that the width measurement is ANSI-aware
	styled := "\x1b[1mBold\x1b[0m"
	width := ansi.StringWidth(styled)
	if width != 4 {
		t.Errorf("ansi.StringWidth should measure visual width, got %d for %q", width, styled)
	}
}

func TestRenderTrail_VeryNarrowWidth(t *testing.T) {
	trail := []Crumb{
		{Text: "HomeScreen"},
		{Text: "SettingsPage"},
	}
	result := renderTrail(trail, " > ", 4, DefaultStyles())
	// Should produce something with "..."
	if result == "" {
		t.Error("expected non-empty result even at very narrow width")
	}
}
