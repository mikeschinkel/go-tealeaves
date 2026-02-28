package teamodal

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newTestProgressModal(bgEnabled bool) ProgressModal {
	m := NewProgressModal(&ProgressModalArgs{
		ScreenWidth:       80,
		ScreenHeight:      24,
		Title:             "Commit Message",
		BackgroundEnabled: bgEnabled,
	})
	m = m.Open()
	return m
}

// --- Layer 1 ---

func TestNewProgressModal(t *testing.T) {
	m := NewProgressModal(&ProgressModalArgs{Title: "Test"})
	if m.IsOpen() {
		t.Error("expected not open initially")
	}
}

func TestProgressModal_Open(t *testing.T) {
	m := NewProgressModal(&ProgressModalArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Test",
	})
	m = m.Open()
	if !m.IsOpen() {
		t.Error("expected IsOpen=true after Open()")
	}
}

func TestProgressModal_Close(t *testing.T) {
	m := newTestProgressModal(false)
	m = m.Close()
	if m.IsOpen() {
		t.Error("expected IsOpen=false after Close()")
	}
}

func TestProgressModal_EscCancels(t *testing.T) {
	m := newTestProgressModal(false)
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if m.IsOpen() {
		t.Error("expected modal closed after Esc")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(ProgressCancelledMsg)
	if !ok {
		t.Fatalf("expected ProgressCancelledMsg, got %T", msg)
	}
}

func TestProgressModal_BackgroundKey(t *testing.T) {
	m := newTestProgressModal(true)
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

	if m.IsOpen() {
		t.Error("expected modal closed after 'b'")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(ProgressBackgroundMsg)
	if !ok {
		t.Fatalf("expected ProgressBackgroundMsg, got %T", msg)
	}
}

func TestProgressModal_BackgroundDisabled(t *testing.T) {
	m := newTestProgressModal(false)
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

	// Should remain open when background is disabled
	if !m.IsOpen() {
		t.Error("expected modal to remain open when background disabled")
	}
	if cmd != nil {
		t.Error("expected nil cmd when background disabled")
	}
}

func TestProgressModal_ClosedIgnoresInput(t *testing.T) {
	m := newTestProgressModal(false)
	m = m.Close()
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected nil cmd when modal is closed")
	}
}

// --- Layer 2 ---

func TestProgressModal_View_Open(t *testing.T) {
	m := newTestProgressModal(false)
	view := m.View()
	if !strings.Contains(view, "Commit Message") {
		t.Error("expected view to contain title")
	}
	if !strings.Contains(view, "cancel") || !strings.Contains(view, "esc") {
		t.Error("expected view to contain cancel hint")
	}
}

func TestProgressModal_View_BackgroundHint(t *testing.T) {
	t.Run("Enabled", func(t *testing.T) {
		m := newTestProgressModal(true)
		view := m.View()
		if !strings.Contains(view, "Background") {
			t.Error("expected 'Background' hint when enabled")
		}
	})

	t.Run("Disabled", func(t *testing.T) {
		m := newTestProgressModal(false)
		view := m.View()
		if strings.Contains(view, "Background") {
			t.Error("expected no 'Background' hint when disabled")
		}
	})
}

func TestProgressModal_OverlayModal(t *testing.T) {
	m := newTestProgressModal(false)

	bgLines := make([]string, 24)
	for i := range bgLines {
		bgLines[i] = strings.Repeat(".", 80)
	}
	background := strings.Join(bgLines, "\n")

	view := m.OverlayModal(background)
	if view == background {
		t.Error("expected overlay to modify background when open")
	}

	m = m.Close()
	view = m.OverlayModal(background)
	if view != background {
		t.Error("expected background unchanged when closed")
	}
}
