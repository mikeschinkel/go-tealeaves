package teamodal

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

func TestOverlayModal_BasicOverlay(t *testing.T) {
	bg := makeBackground(20, 10)
	fg := "XXXX\nXXXX"
	result := OverlayModal(bg, fg, 3, 5)

	lines := strings.Split(result, "\n")
	if len(lines) != 10 {
		t.Fatalf("expected 10 lines, got %d", len(lines))
	}
	// Lines before overlay should be unchanged
	if lines[0] != strings.Repeat(".", 20) {
		t.Errorf("line 0 should be unchanged background")
	}
	// Overlay lines should contain foreground
	if !strings.Contains(lines[3], "XXXX") {
		t.Errorf("line 3 should contain foreground")
	}
	if !strings.Contains(lines[4], "XXXX") {
		t.Errorf("line 4 should contain foreground")
	}
}

func TestOverlayModal_ANSIContent(t *testing.T) {
	bg := makeBackground(30, 5)
	fg := "\x1b[31mRED\x1b[0m"
	result := OverlayModal(bg, fg, 1, 5)

	lines := strings.Split(result, "\n")
	if !strings.Contains(lines[1], "RED") {
		t.Error("expected overlaid line to contain 'RED'")
	}
}

func TestOverlayModal_BoundaryConditions(t *testing.T) {
	bg := makeBackground(20, 5)

	t.Run("AtOrigin", func(t *testing.T) {
		result := OverlayModal(bg, "AB", 0, 0)
		lines := strings.Split(result, "\n")
		if !strings.HasPrefix(lines[0], "AB") {
			t.Errorf("expected line 0 to start with 'AB'")
		}
	})

	t.Run("RowOutOfBounds", func(t *testing.T) {
		result := OverlayModal(bg, "AB", 10, 0)
		if result != bg {
			t.Error("expected unchanged background when row out of bounds")
		}
	})

	t.Run("ColAtEdge", func(t *testing.T) {
		result := OverlayModal(bg, "AB", 0, 18)
		lines := strings.Split(result, "\n")
		if !strings.Contains(lines[0], "AB") {
			t.Error("expected foreground at edge to be visible")
		}
	})
}
