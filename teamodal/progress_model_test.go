package teamodal

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestProgressModel(bgEnabled bool) ProgressModel {
	m := NewProgressModel(&ProgressModelArgs{
		ScreenWidth:       80,
		ScreenHeight:      24,
		Title:             "Commit Message",
		BackgroundEnabled: bgEnabled,
	})
	m = m.Open()
	return m
}

// --- Layer 1 ---

func TestNewProgressModel(t *testing.T) {
	m := NewProgressModel(&ProgressModelArgs{Title: "Test"})
	if m.IsOpen() {
		t.Error("expected not open initially")
	}
}

func TestProgressModel_Open(t *testing.T) {
	m := NewProgressModel(&ProgressModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Test",
	})
	m = m.Open()
	if !m.IsOpen() {
		t.Error("expected IsOpen=true after Open()")
	}
}

func TestProgressModel_Close(t *testing.T) {
	m := newTestProgressModel(false)
	m = m.Close()
	if m.IsOpen() {
		t.Error("expected IsOpen=false after Close()")
	}
}

func TestProgressModel_EscCancels(t *testing.T) {
	m := newTestProgressModel(false)
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})

	if m.IsOpen() {
		t.Error("expected modal closed after Esc")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(ProgressCancelledMsg)
	if !ok {
		t.Fatalf("expected ProgressCancelledMsg, got %T", msg)
	}
}

func TestProgressModel_BackgroundKey(t *testing.T) {
	m := newTestProgressModel(true)
	m, cmd := m.Update(tea.KeyPressMsg{Code: 'b', Text: "b"})

	if m.IsOpen() {
		t.Error("expected modal closed after 'b'")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(ProgressBackgroundMsg)
	if !ok {
		t.Fatalf("expected ProgressBackgroundMsg, got %T", msg)
	}
}

func TestProgressModel_BackgroundDisabled(t *testing.T) {
	m := newTestProgressModel(false)
	m, cmd := m.Update(tea.KeyPressMsg{Code: 'b', Text: "b"})

	// Should remain open when background is disabled
	if !m.IsOpen() {
		t.Error("expected modal to remain open when background disabled")
	}
	if cmd != nil {
		t.Error("expected nil cmd when background disabled")
	}
}

func TestProgressModel_ClosedIgnoresInput(t *testing.T) {
	m := newTestProgressModel(false)
	m = m.Close()
	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected nil cmd when modal is closed")
	}
}

// --- Layer 2 ---

func TestProgressModel_View_Open(t *testing.T) {
	m := newTestProgressModel(false)
	view := m.View()
	if !strings.Contains(view.Content, "Commit Message") {
		t.Error("expected view to contain title")
	}
	if !strings.Contains(view.Content, "cancel") || !strings.Contains(view.Content, "esc") {
		t.Error("expected view to contain cancel hint")
	}
}

func TestProgressModel_View_BackgroundHint(t *testing.T) {
	t.Run("Enabled", func(t *testing.T) {
		m := newTestProgressModel(true)
		view := m.View()
		if !strings.Contains(view.Content, "Background") {
			t.Error("expected 'Background' hint when enabled")
		}
	})

	t.Run("Disabled", func(t *testing.T) {
		m := newTestProgressModel(false)
		view := m.View()
		if strings.Contains(view.Content, "Background") {
			t.Error("expected no 'Background' hint when disabled")
		}
	})
}

func TestProgressModel_OverlayModal(t *testing.T) {
	m := newTestProgressModel(false)

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
