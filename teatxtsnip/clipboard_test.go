package teatxtsnip

import (
	"testing"
)

// --- Layer 1: Clipboard Tests ---

func TestSelectedText_SingleLine(t *testing.T) {
	m := newTestModel()
	// Select "Hello" (first 5 chars on line 0)
	m.selection = m.selection.Begin(Position{Row: 0, Col: 0})
	m.selection = m.selection.Extend(Position{Row: 0, Col: 5})

	text := m.SelectedText()
	if text != "Hello" {
		t.Errorf("expected 'Hello', got %q", text)
	}
}

func TestSelectedText_MultiLine(t *testing.T) {
	m := newTestModel()
	// Select from "World" to "Second"
	m.selection = m.selection.Begin(Position{Row: 0, Col: 6})
	m.selection = m.selection.Extend(Position{Row: 1, Col: 6})

	text := m.SelectedText()
	expected := "World\nSecond"
	if text != expected {
		t.Errorf("expected %q, got %q", expected, text)
	}
}

func TestSelectedText_NoSelection(t *testing.T) {
	m := newTestModel()
	text := m.SelectedText()
	if text != "" {
		t.Errorf("expected empty string with no selection, got %q", text)
	}
}

func TestCopy(t *testing.T) {
	m := newTestModel()
	// Select "Hello"
	m.selection = m.selection.Begin(Position{Row: 0, Col: 0})
	m.selection = m.selection.Extend(Position{Row: 0, Col: 5})

	m = m.Copy()

	// Check internal clipboard (mock)
	if m.internalClip != "Hello" {
		t.Errorf("expected internalClip='Hello', got %q", m.internalClip)
	}
	// Selection should still be active after copy
	if !m.HasSelection() {
		t.Error("expected selection to remain active after Copy")
	}
}

func TestCut(t *testing.T) {
	m := newTestModel()
	// Select "Hello"
	m.selection = m.selection.Begin(Position{Row: 0, Col: 0})
	m.selection = m.selection.Extend(Position{Row: 0, Col: 5})

	m = m.Cut()

	// Check internal clipboard
	if m.internalClip != "Hello" {
		t.Errorf("expected internalClip='Hello', got %q", m.internalClip)
	}
	// Selection should be cleared after cut
	if m.HasSelection() {
		t.Error("expected selection cleared after Cut")
	}
	// "Hello" should be removed from content
	value := m.Value()
	if len(value) > 0 && value[0] == 'H' {
		t.Error("expected 'Hello' to be removed from content")
	}
}

func TestPaste_NoSelection(t *testing.T) {
	m := newTestModel()
	m.moveCursorTo(0, 0)
	// Write to clipboard via writeToClipboard (writes to both system and internal)
	m = m.writeToClipboard("PASTED")

	m, _ = m.Paste()

	value := m.Value()
	if len(value) < 6 || value[:6] != "PASTED" {
		t.Errorf("expected value to start with 'PASTED', got %q", value[:min(20, len(value))])
	}
}

func TestPaste_WithSelection(t *testing.T) {
	m := newTestModel()
	// Select "Hello"
	m.selection = m.selection.Begin(Position{Row: 0, Col: 0})
	m.selection = m.selection.Extend(Position{Row: 0, Col: 5})
	// Write to clipboard
	m = m.writeToClipboard("Goodbye")

	m, _ = m.Paste()

	value := m.Value()
	if len(value) < 7 || value[:7] != "Goodbye" {
		t.Errorf("expected value to start with 'Goodbye', got %q", value[:min(20, len(value))])
	}
	// Selection should be cleared
	if m.HasSelection() {
		t.Error("expected selection cleared after paste with selection")
	}
}
