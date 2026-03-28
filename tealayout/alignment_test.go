package tealayout

import (
	"strings"
	"testing"
)

func TestAlignment_MergeSingleAxis(t *testing.T) {
	// Setting Top preserves existing horizontal
	a := mergeAlignment(Left, Top)
	if a != TopLeft {
		t.Errorf("merge(Left, Top) = %d, want TopLeft(%d)", a, TopLeft)
	}
}

func TestAlignment_MergeBothAxes(t *testing.T) {
	a := mergeAlignment(TopLeft, BottomRight)
	if a != BottomRight {
		t.Errorf("merge(TopLeft, BottomRight) = %d, want BottomRight(%d)", a, BottomRight)
	}
}

func TestAlignment_MergePreservesOtherAxis(t *testing.T) {
	// Start with BottomLeft, set Center (horizontal only)
	a := mergeAlignment(BottomLeft, Center)
	if a != BottomCenter {
		t.Errorf("merge(BottomLeft, Center) = %d, want BottomCenter(%d)", a, BottomCenter)
	}
}

func TestAlignment_WithAlignment_Composable(t *testing.T) {
	p := NewRow(Flex(1)).WithAlignment(Top).WithAlignment(Right)
	if p.alignment != TopRight {
		t.Errorf("alignment = %d, want TopRight(%d)", p.alignment, TopRight)
	}
}

func TestAlignContent_Right(t *testing.T) {
	content := "hi"
	result := alignContent(content, 10, 1, Right)
	if !strings.HasPrefix(result, "        ") {
		t.Errorf("Right alignment should pad left: got %q", result)
	}
	if !strings.HasSuffix(result, "hi") {
		t.Errorf("Right alignment content missing: got %q", result)
	}
}

func TestAlignContent_Center(t *testing.T) {
	content := "hi"
	result := alignContent(content, 10, 1, Center)
	// "hi" is 2 chars, 10-2=8 pad, center = 4 left pad
	if !strings.HasPrefix(result, "    ") {
		t.Errorf("Center alignment should pad: got %q", result)
	}
}

func TestAlignContent_Bottom(t *testing.T) {
	content := "hi"
	result := alignContent(content, 10, 5, Bottom)
	lines := strings.Split(result, "\n")
	// 1 content line + 4 blank lines above = 5 total
	if len(lines) != 5 {
		t.Errorf("Bottom: got %d lines, want 5", len(lines))
	}
	if strings.TrimSpace(lines[0]) != "" {
		t.Error("Bottom: first line should be blank")
	}
	if !strings.Contains(lines[4], "hi") {
		t.Errorf("Bottom: last line should contain 'hi', got %q", lines[4])
	}
}

func TestAlignContent_Middle(t *testing.T) {
	content := "hi"
	result := alignContent(content, 10, 5, Middle)
	lines := strings.Split(result, "\n")
	// 1 content line, vPad=4, topPad=2
	if len(lines) < 3 {
		t.Fatalf("Middle: got %d lines, want >= 3", len(lines))
	}
	if strings.TrimSpace(lines[0]) != "" {
		t.Error("Middle: first line should be blank")
	}
	if !strings.Contains(lines[2], "hi") {
		t.Errorf("Middle: line[2] should contain 'hi', got %q", lines[2])
	}
}

func TestAlignContent_MiddleCenter(t *testing.T) {
	content := "hi"
	result := alignContent(content, 10, 5, MiddleCenter)
	lines := strings.Split(result, "\n")
	// Vertical: topPad=2, so content at line index 2
	if len(lines) < 3 {
		t.Fatalf("MiddleCenter: got %d lines, want >= 3", len(lines))
	}
	// Horizontal: "hi" is 2 wide, 10-2=8, center=4
	if !strings.HasPrefix(lines[2], "    ") {
		t.Errorf("MiddleCenter: content line should be centered, got %q", lines[2])
	}
}
