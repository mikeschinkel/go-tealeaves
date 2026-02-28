package teadd

import (
	"strings"
	"testing"
)

func makeBackground(width, height int) string {
	lines := make([]string, height)
	for i := range lines {
		lines[i] = strings.Repeat(".", width)
	}
	return strings.Join(lines, "\n")
}

func TestOverlayDropdown_BasicOverlay(t *testing.T) {
	bg := makeBackground(20, 10)
	fg := "XXXX\nXXXX"
	result := OverlayDropdown(bg, fg, 3, 5)

	lines := strings.Split(result, "\n")
	if len(lines) != 10 {
		t.Fatalf("expected 10 lines, got %d", len(lines))
	}
	// Lines 0-2 should be unchanged background
	if lines[0] != strings.Repeat(".", 20) {
		t.Errorf("line 0 should be unchanged background, got %q", lines[0])
	}
	// Lines 3-4 should contain foreground
	if !strings.Contains(lines[3], "XXXX") {
		t.Errorf("line 3 should contain foreground, got %q", lines[3])
	}
	if !strings.Contains(lines[4], "XXXX") {
		t.Errorf("line 4 should contain foreground, got %q", lines[4])
	}
}

func TestOverlayDropdown_AtOrigin(t *testing.T) {
	bg := makeBackground(20, 5)
	fg := "AB"
	result := OverlayDropdown(bg, fg, 0, 0)

	lines := strings.Split(result, "\n")
	if !strings.HasPrefix(lines[0], "AB") {
		t.Errorf("expected line 0 to start with 'AB', got %q", lines[0])
	}
}

func TestOverlayDropdown_OutOfBounds(t *testing.T) {
	bg := makeBackground(10, 3)
	fg := "OVERLAY"

	// Row beyond background — should return background unchanged
	result := OverlayDropdown(bg, fg, 10, 0)
	if result != bg {
		t.Errorf("expected unchanged background when overlay row is out of bounds")
	}
}

func TestOverlayDropdown_EmptyForeground(t *testing.T) {
	bg := makeBackground(10, 3)
	result := OverlayDropdown(bg, "", 0, 0)
	// Empty foreground splits to [""], so line 0 gets overlaid with empty
	// Just verify it doesn't crash and background structure is preserved
	lines := strings.Split(result, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestOverlayDropdown_ANSIContent(t *testing.T) {
	bg := makeBackground(20, 5)
	fg := "\x1b[31mRED\x1b[0m"
	result := OverlayDropdown(bg, fg, 1, 5)

	lines := strings.Split(result, "\n")
	// The overlaid line should contain our ANSI content
	if !strings.Contains(lines[1], "RED") {
		t.Errorf("expected line 1 to contain 'RED', got %q", lines[1])
	}
}
