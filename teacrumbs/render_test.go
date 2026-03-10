package teacrumbs

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// helper builds a BreadcrumbsModel from crumbs, separator, and width for tests.
func testModel(crumbs []Crumb, sep string, width int) BreadcrumbsModel {
	return NewBreadcrumbsModel().
		SetCrumbs(crumbs).
		WithSeparator(sep).
		SetSize(width)
}

func TestRender_Empty(t *testing.T) {
	result := testModel(nil, " > ", 80).render()
	if result.content != "" {
		t.Errorf("expected empty string, got %q", result.content)
	}
}

func TestRender_SingleCrumb(t *testing.T) {
	crumbs := []Crumb{{Text: "Home"}}
	result := testModel(crumbs, " > ", 80).render()
	if !strings.Contains(result.content, "Home") {
		t.Errorf("expected result to contain 'Home', got %q", result.content)
	}
}

func TestRender_MultipleCrumbs(t *testing.T) {
	crumbs := []Crumb{
		{Text: "Home"},
		{Text: "Settings"},
		{Text: "Profile"},
	}
	result := testModel(crumbs, " > ", 0).render()
	if !strings.Contains(result.content, "Home") {
		t.Errorf("expected result to contain 'Home', got %q", result.content)
	}
	if !strings.Contains(result.content, "Settings") {
		t.Errorf("expected result to contain 'Settings', got %q", result.content)
	}
	if !strings.Contains(result.content, "Profile") {
		t.Errorf("expected result to contain 'Profile', got %q", result.content)
	}
}

func TestRender_CustomSeparator(t *testing.T) {
	crumbs := []Crumb{
		{Text: "A"},
		{Text: "B"},
	}
	result := testModel(crumbs, " / ", 0).render()
	if !strings.Contains(result.content, "/") {
		t.Errorf("expected custom separator '/', got %q", result.content)
	}
}

func TestRender_TruncationWithEllipsis(t *testing.T) {
	crumbs := []Crumb{
		{Text: "Home"},
		{Text: "Very Long Middle Path"},
		{Text: "Another Long Middle Path"},
		{Text: "Current"},
	}
	result := testModel(crumbs, " > ", 30).render()
	if !strings.Contains(result.content, "...") {
		t.Errorf("expected ellipsis in truncated result, got %q", result.content)
	}
	if !strings.Contains(result.content, "Home") {
		t.Errorf("expected first crumb 'Home' preserved, got %q", result.content)
	}
}

func TestRender_ExtremeTruncation(t *testing.T) {
	crumbs := []Crumb{
		{Text: "Very Long Home Path"},
		{Text: "Middle"},
		{Text: "Very Long Current Path"},
	}
	result := testModel(crumbs, " > ", 10).render()
	if !strings.Contains(result.content, "...") {
		t.Errorf("expected ellipsis in extreme truncation, got %q", result.content)
	}
}

func TestRender_ZeroWidth_NoTruncation(t *testing.T) {
	crumbs := []Crumb{
		{Text: "Home"},
		{Text: "Settings"},
		{Text: "Profile"},
	}
	result := testModel(crumbs, " > ", 0).render()
	if strings.Contains(result.content, "...") {
		t.Errorf("zero width should not truncate, got %q", result.content)
	}
	if !strings.Contains(result.content, "Settings") {
		t.Errorf("expected all crumbs present, got %q", result.content)
	}
}

func TestRender_UsesAnsiStringWidth(t *testing.T) {
	styled := "\x1b[1mBold\x1b[0m"
	width := ansi.StringWidth(styled)
	if width != 4 {
		t.Errorf("ansi.StringWidth should measure visual width, got %d for %q", width, styled)
	}
}

func TestRender_VeryNarrowWidth(t *testing.T) {
	crumbs := []Crumb{
		{Text: "HomeScreen"},
		{Text: "SettingsPage"},
	}
	result := testModel(crumbs, " > ", 4).render()
	if result.content == "" {
		t.Error("expected non-empty result even at very narrow width")
	}
}

// --- Short text tests ---

func TestRender_ShortText_DefaultBehavior(t *testing.T) {
	crumbs := []Crumb{
		{Text: "go-dt", Short: "dt"},
		{Text: "go-utils", Short: "utils"},
	}
	// Wide enough for all Text — non-current uses Short, last uses Text
	result := testModel(crumbs, " > ", 80).render()
	if !strings.Contains(result.content, "dt") {
		t.Errorf("non-current should use Short text, got %q", result.content)
	}
	if !strings.Contains(result.content, "go-utils") {
		t.Errorf("current (last) should use Text, got %q", result.content)
	}
}

func TestRender_ShortFallsBackToText(t *testing.T) {
	crumbs := []Crumb{
		{Text: "Home"}, // No Short set
		{Text: "go-dt", Short: "dt"},
	}
	result := testModel(crumbs, " > ", 80).render()
	if !strings.Contains(result.content, "Home") {
		t.Errorf("empty Short should fall back to Text, got %q", result.content)
	}
	if !strings.Contains(result.content, "go-dt") {
		t.Errorf("current should use Text, got %q", result.content)
	}
}

func TestRender_CompactFallback(t *testing.T) {
	crumbs := []Crumb{
		{Text: "github.com/mikeschinkel/go-dt", Short: "dt"},
		{Text: "github.com/mikeschinkel/go-utils", Short: "utils"},
	}
	// Too narrow for default text but wide enough for Short
	result := testModel(crumbs, " > ", 20).render()
	// Should fall back to compact mode (all Short)
	if strings.Contains(result.content, "github.com") {
		t.Errorf("compact fallback should use Short text, got %q", result.content)
	}
}

func TestRender_PerCrumbStyleOverride(t *testing.T) {
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	crumbs := []Crumb{
		{Text: "Home"},
		{Text: "Special", Style: &red},
	}
	result := testModel(crumbs, " > ", 0).render()
	if !strings.Contains(result.content, "Special") {
		t.Errorf("per-crumb style should still render text, got %q", result.content)
	}
}

func TestRender_MixedStyles(t *testing.T) {
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	crumbs := []Crumb{
		{Text: "Normal"},
		{Text: "Styled", Style: &green},
		{Text: "Current"},
	}
	result := testModel(crumbs, " > ", 0).render()
	if !strings.Contains(result.content, "Normal") {
		t.Errorf("expected Normal in output, got %q", result.content)
	}
	if !strings.Contains(result.content, "Styled") {
		t.Errorf("expected Styled in output, got %q", result.content)
	}
	if !strings.Contains(result.content, "Current") {
		t.Errorf("expected Current in output, got %q", result.content)
	}
}

// --- Renderer tests ---

func TestRendererFunc_Adapter(t *testing.T) {
	called := false
	fn := RendererFunc(func(index int, model BreadcrumbsModel) string {
		called = true
		return "[" + model.crumbs[index].Text + "]"
	})
	m := testModel([]Crumb{{Text: "test"}}, " > ", 80)
	result := fn.Render(0, m)
	if !called {
		t.Error("RendererFunc should call the underlying function")
	}
	if result != "[test]" {
		t.Errorf("expected '[test]', got %q", result)
	}
}

func TestRenderer_OverridesStandardStyling(t *testing.T) {
	crumbs := []Crumb{
		{Text: "Home"},
		{
			Text: "Custom",
			Renderer: RendererFunc(func(index int, model BreadcrumbsModel) string {
				return "<<" + model.crumbs[index].Text + ">>"
			}),
		},
	}
	result := testModel(crumbs, " > ", 0).render()
	if !strings.Contains(result.content, "<<Custom>>") {
		t.Errorf("renderer should override styling, got %q", result.content)
	}
}

func TestRenderer_ContextFields(t *testing.T) {
	var capturedIndex int
	var capturedModel BreadcrumbsModel
	crumbs := []Crumb{
		{Text: "First"},
		{
			Text:  "Last",
			Short: "L",
			Data:  "payload",
			Renderer: RendererFunc(func(index int, model BreadcrumbsModel) string {
				capturedIndex = index
				capturedModel = model
				return model.crumbs[index].Text
			}),
		},
	}
	testModel(crumbs, " > ", 80).render()

	if capturedIndex != 1 {
		t.Errorf("expected index=1, got %d", capturedIndex)
	}
	if capturedModel.Len() != 2 {
		t.Errorf("expected model with 2 crumbs, got %d", capturedModel.Len())
	}
	if capturedModel.Width() != 80 {
		t.Errorf("expected Width=80, got %d", capturedModel.Width())
	}
	crumb := capturedModel.Crumbs()[1]
	if crumb.Data != "payload" {
		t.Errorf("expected Data='payload', got %v", crumb.Data)
	}
}

func TestRenderer_IsCurrent_FirstCrumbNotCurrent(t *testing.T) {
	var capturedIndex int
	crumbs := []Crumb{
		{
			Text: "First",
			Renderer: RendererFunc(func(index int, model BreadcrumbsModel) string {
				capturedIndex = index
				return model.crumbs[index].Text
			}),
		},
		{Text: "Last"},
	}
	testModel(crumbs, " > ", 0).render()
	if capturedIndex != 0 {
		t.Errorf("expected index=0, got %d", capturedIndex)
	}
}

func TestRenderer_NilFallsThrough(t *testing.T) {
	crumbs := []Crumb{
		{Text: "Normal", Renderer: nil},
	}
	result := testModel(crumbs, " > ", 0).render()
	if !strings.Contains(result.content, "Normal") {
		t.Errorf("nil renderer should fall through to standard styling, got %q", result.content)
	}
}

func TestRenderer_TruncationWithWideOutput(t *testing.T) {
	crumbs := []Crumb{
		{Text: "A"},
		{Text: "B"},
		{
			Text: "C",
			Renderer: RendererFunc(func(index int, model BreadcrumbsModel) string {
				return "VeryVeryVeryLongRenderedOutput"
			}),
		},
	}
	result := testModel(crumbs, " > ", 20).render()
	if result.content == "" {
		t.Error("expected non-empty result with wide renderer output")
	}
}
