package teamodal

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestChoiceModel(options []ChoiceOption, defaultIndex int) ChoiceModel {
	m := NewChoiceModel(&ChoiceModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Message:      "Test message",
		Options:      options,
		DefaultIndex: defaultIndex,
	})
	m, _ = m.Open()
	return m
}

func threeOptions() []ChoiceOption {
	return []ChoiceOption{
		{Label: "Reorganize & Exit", Hotkey: 'r', ID: "reorganize"},
		{Label: "Save & Exit", Hotkey: 's', ID: "save"},
		{Label: "Cancel", ID: "cancel"},
	}
}

func extractMsg(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

func TestChoiceModel_NavigationForward(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	// Tab moves focus forward
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = result.(ChoiceModel)
	if m.FocusButton() != 1 {
		t.Errorf("expected focusButton=1 after Tab, got %d", m.FocusButton())
	}
	if cmd == nil {
		t.Error("expected non-nil cmd after Tab (message consumed)")
	}

	// Tab again
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = result.(ChoiceModel)
	if m.FocusButton() != 2 {
		t.Errorf("expected focusButton=2 after second Tab, got %d", m.FocusButton())
	}

	// Tab wraps to 0
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = result.(ChoiceModel)
	if m.FocusButton() != 0 {
		t.Errorf("expected focusButton=0 after wrap, got %d", m.FocusButton())
	}
}

func TestChoiceModel_NavigationBackward(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	// Shift+Tab from 0 wraps to last
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m = result.(ChoiceModel)
	if m.FocusButton() != 2 {
		t.Errorf("expected focusButton=2 after Shift+Tab from 0, got %d", m.FocusButton())
	}
	if cmd == nil {
		t.Error("expected non-nil cmd after Shift+Tab (message consumed)")
	}

	// Shift+Tab again
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	m = result.(ChoiceModel)
	if m.FocusButton() != 1 {
		t.Errorf("expected focusButton=1 after second Shift+Tab, got %d", m.FocusButton())
	}
}

func TestChoiceModel_NavigationWithRightLeft(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	// Right arrow moves forward
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	m = result.(ChoiceModel)
	if m.FocusButton() != 1 {
		t.Errorf("expected focusButton=1 after Right, got %d", m.FocusButton())
	}

	// Left arrow moves backward
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	m = result.(ChoiceModel)
	if m.FocusButton() != 0 {
		t.Errorf("expected focusButton=0 after Left, got %d", m.FocusButton())
	}
}

func TestChoiceModel_SelectWithEnter(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	// Move to second option then select
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = result.(ChoiceModel)

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(ChoiceModel)

	if m.IsOpen() {
		t.Error("expected modal to be closed after Enter")
	}

	msg := extractMsg(cmd)
	selectedMsg, ok := msg.(ChoiceSelectedMsg)
	if !ok {
		t.Fatalf("expected ChoiceSelectedMsg, got %T", msg)
	}
	if selectedMsg.OptionID != "save" {
		t.Errorf("expected OptionID='save', got %q", selectedMsg.OptionID)
	}
	if selectedMsg.Index != 1 {
		t.Errorf("expected Index=1, got %d", selectedMsg.Index)
	}
}

func TestChoiceModel_CancelWithEsc(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = result.(ChoiceModel)

	if m.IsOpen() {
		t.Error("expected modal to be closed after Esc")
	}

	msg := extractMsg(cmd)
	_, ok := msg.(ChoiceCancelledMsg)
	if !ok {
		t.Fatalf("expected ChoiceCancelledMsg, got %T", msg)
	}
}

func TestChoiceModel_HotkeyLowercase(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 2) // Focus on Cancel

	// Press 'r' (lowercase) should select "Reorganize & Exit"
	result, cmd := m.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})
	m = result.(ChoiceModel)

	if m.IsOpen() {
		t.Error("expected modal to be closed after hotkey press")
	}

	msg := extractMsg(cmd)
	selectedMsg, ok := msg.(ChoiceSelectedMsg)
	if !ok {
		t.Fatalf("expected ChoiceSelectedMsg, got %T", msg)
	}
	if selectedMsg.OptionID != "reorganize" {
		t.Errorf("expected OptionID='reorganize', got %q", selectedMsg.OptionID)
	}
	if selectedMsg.Index != 0 {
		t.Errorf("expected Index=0, got %d", selectedMsg.Index)
	}
}

func TestChoiceModel_HotkeyUppercase(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	// Press 'S' (uppercase) should select "Save & Exit" (case-insensitive)
	result, cmd := m.Update(tea.KeyPressMsg{Code: 'S', Text: "S"})
	m = result.(ChoiceModel)

	if m.IsOpen() {
		t.Error("expected modal to be closed after uppercase hotkey")
	}

	msg := extractMsg(cmd)
	selectedMsg, ok := msg.(ChoiceSelectedMsg)
	if !ok {
		t.Fatalf("expected ChoiceSelectedMsg, got %T", msg)
	}
	if selectedMsg.OptionID != "save" {
		t.Errorf("expected OptionID='save', got %q", selectedMsg.OptionID)
	}
}

func TestChoiceModel_HotkeyNoMatch(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	// Press 'z' which has no matching hotkey
	result, cmd := m.Update(tea.KeyPressMsg{Code: 'z', Text: "z"})
	m = result.(ChoiceModel)

	if !m.IsOpen() {
		t.Error("expected modal to remain open when no hotkey matches")
	}
	if cmd != nil {
		t.Error("expected nil cmd when no hotkey matches")
	}
}

func TestChoiceModel_DefaultIndex(t *testing.T) {
	t.Run("DefaultIndex0", func(t *testing.T) {
		m := newTestChoiceModel(threeOptions(), 0)
		if m.FocusButton() != 0 {
			t.Errorf("expected focusButton=0, got %d", m.FocusButton())
		}
	})

	t.Run("DefaultIndex2", func(t *testing.T) {
		m := newTestChoiceModel(threeOptions(), 2)
		if m.FocusButton() != 2 {
			t.Errorf("expected focusButton=2, got %d", m.FocusButton())
		}
	})

	t.Run("DefaultIndexOutOfBounds", func(t *testing.T) {
		m := newTestChoiceModel(threeOptions(), 99)
		if m.FocusButton() != 0 {
			t.Errorf("expected focusButton=0 for out-of-bounds index, got %d", m.FocusButton())
		}
	})

	t.Run("DefaultIndexNegative", func(t *testing.T) {
		m := newTestChoiceModel(threeOptions(), -1)
		if m.FocusButton() != 0 {
			t.Errorf("expected focusButton=0 for negative index, got %d", m.FocusButton())
		}
	})
}

func TestChoiceModel_ClosedModalIgnoresInput(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)
	m, _ = m.Close()

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if cmd != nil {
		t.Error("expected nil cmd when modal is closed")
	}
}

func TestChoiceModel_OverlayModal(t *testing.T) {
	m := newTestChoiceModel(threeOptions(), 0)

	// Build a background large enough for the modal to overlay onto.
	// The modal is centered based on screenHeight (24), so the background
	// needs at least that many lines for the overlay to be visible.
	bgLines := make([]string, 24)
	for i := range bgLines {
		bgLines[i] = strings.Repeat(" ", 80)
	}
	background := strings.Join(bgLines, "\n")

	view := m.OverlayModal(background)

	// When open, overlay should be different from background
	if view == background {
		t.Error("expected overlay to modify background when modal is open")
	}
	if view == "" {
		t.Error("expected non-empty overlay view")
	}

	// When closed, should return background unchanged
	m, _ = m.Close()
	view = m.OverlayModal(background)
	if view != background {
		t.Error("expected background unchanged when modal is closed")
	}
}
