package teautils

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestRenderAlignedLine_Left(t *testing.T) {
	style := lipgloss.NewStyle()
	result := RenderAlignedLine("Hello", style, 20, lipgloss.Left)
	if !strings.Contains(result, "Hello") {
		t.Errorf("expected result to contain 'Hello', got %q", result)
	}
	// Left-aligned text should start at the beginning (after possible padding)
	trimmed := strings.TrimRight(result, " ")
	if !strings.HasPrefix(trimmed, "Hello") {
		t.Errorf("expected left-aligned text to start with 'Hello', got %q", trimmed)
	}
}

func TestRenderAlignedLine_Center(t *testing.T) {
	style := lipgloss.NewStyle()
	result := RenderAlignedLine("Hi", style, 20, lipgloss.Center)
	if !strings.Contains(result, "Hi") {
		t.Errorf("expected result to contain 'Hi', got %q", result)
	}
	// Centered text should have leading spaces
	idx := strings.Index(result, "Hi")
	if idx == 0 {
		t.Error("expected centered text to have leading spaces")
	}
}

func TestRenderAlignedLine_Right(t *testing.T) {
	style := lipgloss.NewStyle()
	result := RenderAlignedLine("Hi", style, 20, lipgloss.Right)
	if !strings.Contains(result, "Hi") {
		t.Errorf("expected result to contain 'Hi', got %q", result)
	}
	// Right-aligned text should have leading spaces
	idx := strings.Index(result, "Hi")
	if idx < 10 {
		t.Errorf("expected right-aligned text to be near end, 'Hi' found at index %d", idx)
	}
}

func TestRenderCenteredLine(t *testing.T) {
	style := lipgloss.NewStyle()
	centered := RenderCenteredLine("Test", style, 20)
	aligned := RenderAlignedLine("Test", style, 20, lipgloss.Center)
	if centered != aligned {
		t.Errorf("RenderCenteredLine should match RenderAlignedLine with Center alignment\ncentered: %q\naligned:  %q", centered, aligned)
	}
}

func TestApplyBoxBorder(t *testing.T) {
	borderStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	result := ApplyBoxBorder(borderStyle, "Content")
	if !strings.Contains(result, "Content") {
		t.Errorf("expected result to contain 'Content', got %q", result)
	}
	// Should contain border characters (rounded border uses ╭, ╮, ╰, ╯)
	if !strings.ContainsAny(result, "╭╮╰╯│─") {
		t.Error("expected result to contain border characters")
	}
}

func TestApplyBoxBorder_MultiLine(t *testing.T) {
	borderStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	content := "Line 1\nLine 2\nLine 3"
	result := ApplyBoxBorder(borderStyle, content)
	if !strings.Contains(result, "Line 1") {
		t.Errorf("expected result to contain 'Line 1', got %q", result)
	}
	if !strings.Contains(result, "Line 3") {
		t.Errorf("expected result to contain 'Line 3', got %q", result)
	}
	// Multi-line content should produce more output lines than single-line
	lines := strings.Split(result, "\n")
	// At minimum: top border + padding + 3 content lines + padding + bottom border = 7
	if len(lines) < 5 {
		t.Errorf("expected at least 5 lines for bordered multi-line content, got %d", len(lines))
	}
}
