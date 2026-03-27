package teamodal

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// msItem implements MultiSelectItem for testing
type msItem struct {
	id    string
	label string
}

func (t msItem) ID() string    { return t.id }
func (t msItem) Label() string { return t.label }

func threeItems() []msItem {
	return []msItem{
		{id: "a", label: "Alpha"},
		{id: "b", label: "Beta"},
		{id: "c", label: "Gamma"},
	}
}

func twoButtons() []MultiSelectButton {
	return []MultiSelectButton{
		{Label: "Update Selected", Hotkey: 'u', ID: "update"},
		{Label: "Skip", Hotkey: 's', ID: "skip"},
	}
}

func newTestMultiSelectModel(items []msItem, allChecked bool) MultiSelectModel[msItem] {
	m := NewMultiSelectModel(items, &MultiSelectModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Test Title",
		Buttons:      twoButtons(),
		AllChecked:   allChecked,
	})
	m, _ = m.Open()
	return m
}

func extractMultiSelectMsg(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

func TestMultiSelectModel_ToggleWithSpace(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// All items start checked
	selected := m.Selected()
	if len(selected) != 3 {
		t.Errorf("expected 3 selected, got %d", len(selected))
	}

	// Space toggles the first item off
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	if cmd == nil {
		t.Error("expected non-nil cmd after space (message consumed)")
	}

	selected = m.Selected()
	if len(selected) != 2 {
		t.Errorf("expected 2 selected after toggle, got %d", len(selected))
	}

	// Space toggles it back on
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	selected = m.Selected()
	if len(selected) != 3 {
		t.Errorf("expected 3 selected after re-toggle, got %d", len(selected))
	}
}

func TestMultiSelectModel_CursorNavigation(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Starts at 0
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0, got %d", m.Cursor())
	}

	// Down
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.Cursor() != 1 {
		t.Errorf("expected cursor=1 after down, got %d", m.Cursor())
	}
	if cmd == nil {
		t.Error("expected non-nil cmd after navigation")
	}

	// Down again
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.Cursor() != 2 {
		t.Errorf("expected cursor=2 after second down, got %d", m.Cursor())
	}

	// Down at boundary — clamped
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.Cursor() != 2 {
		t.Errorf("expected cursor=2 at boundary, got %d", m.Cursor())
	}

	// Up
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.Cursor() != 1 {
		t.Errorf("expected cursor=1 after up, got %d", m.Cursor())
	}

	// Up to 0
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0, got %d", m.Cursor())
	}

	// Up at boundary — clamped
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 at boundary, got %d", m.Cursor())
	}
}

func TestMultiSelectModel_FocusCycling(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Start: focus on list
	if m.focus != focusList {
		t.Errorf("expected focus=focusList, got %d", m.focus)
	}

	// Tab: list → button0
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusButton || m.buttonIdx != 0 {
		t.Errorf("expected focus=focusButton idx=0, got focus=%d idx=%d", m.focus, m.buttonIdx)
	}
	if cmd == nil {
		t.Error("expected non-nil cmd after tab")
	}

	// Tab: button0 → button1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusButton || m.buttonIdx != 1 {
		t.Errorf("expected focus=focusButton idx=1, got focus=%d idx=%d", m.focus, m.buttonIdx)
	}

	// Tab: button1 → button2 (auto-Cancel)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusButton || m.buttonIdx != 2 {
		t.Errorf("expected focus=focusButton idx=2, got focus=%d idx=%d", m.focus, m.buttonIdx)
	}

	// Tab: button2 → list (wrap)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.focus != focusList {
		t.Errorf("expected focus=focusList after wrap, got %d", m.focus)
	}

	// Shift+Tab: list → last button (idx=2, auto-Cancel)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	if m.focus != focusButton || m.buttonIdx != 2 {
		t.Errorf("expected focus=focusButton idx=2 after shift+tab, got focus=%d idx=%d", m.focus, m.buttonIdx)
	}

	// Shift+Tab: button2 → button1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	if m.focus != focusButton || m.buttonIdx != 1 {
		t.Errorf("expected focus=focusButton idx=1 after shift+tab, got focus=%d idx=%d", m.focus, m.buttonIdx)
	}

	// Shift+Tab: button1 → button0
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	if m.focus != focusButton || m.buttonIdx != 0 {
		t.Errorf("expected focus=focusButton idx=0 after shift+tab, got focus=%d idx=%d", m.focus, m.buttonIdx)
	}

	// Shift+Tab: button0 → list
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	if m.focus != focusList {
		t.Errorf("expected focus=focusList after reverse wrap, got %d", m.focus)
	}
}

func TestMultiSelectModel_ButtonActivation(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Tab to button 0
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})

	// Enter activates button
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if m.IsOpen() {
		t.Error("expected modal closed after button activation")
	}

	msg := extractMultiSelectMsg(cmd)
	acceptedMsg, ok := msg.(MultiSelectButtonPressedMsg[msItem])
	if !ok {
		t.Fatalf("expected MultiSelectButtonPressedMsg, got %T", msg)
	}
	if acceptedMsg.ButtonID != "update" {
		t.Errorf("expected ButtonID='update', got %q", acceptedMsg.ButtonID)
	}
	if len(acceptedMsg.Selected) != 3 {
		t.Errorf("expected 3 selected items, got %d", len(acceptedMsg.Selected))
	}
}

func TestMultiSelectModel_ButtonHotkey(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Press 's' hotkey (focus still on list)
	m, cmd := m.Update(tea.KeyPressMsg{Code: 's', Text: "s"})
	if m.IsOpen() {
		t.Error("expected modal closed after hotkey")
	}

	msg := extractMultiSelectMsg(cmd)
	acceptedMsg, ok := msg.(MultiSelectButtonPressedMsg[msItem])
	if !ok {
		t.Fatalf("expected MultiSelectButtonPressedMsg, got %T", msg)
	}
	if acceptedMsg.ButtonID != "skip" {
		t.Errorf("expected ButtonID='skip', got %q", acceptedMsg.ButtonID)
	}
}

func TestMultiSelectModel_ButtonHotkeyUppercase(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Press 'U' (uppercase) — case-insensitive match to 'u' hotkey
	m, cmd := m.Update(tea.KeyPressMsg{Code: 'U', Text: "U"})
	if m.IsOpen() {
		t.Error("expected modal closed after uppercase hotkey")
	}

	msg := extractMultiSelectMsg(cmd)
	acceptedMsg, ok := msg.(MultiSelectButtonPressedMsg[msItem])
	if !ok {
		t.Fatalf("expected MultiSelectButtonPressedMsg, got %T", msg)
	}
	if acceptedMsg.ButtonID != "update" {
		t.Errorf("expected ButtonID='update', got %q", acceptedMsg.ButtonID)
	}
}

func TestMultiSelectModel_EscCancellation(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if m.IsOpen() {
		t.Error("expected modal closed after esc")
	}

	msg := extractMultiSelectMsg(cmd)
	_, ok := msg.(MultiSelectCancelledMsg)
	if !ok {
		t.Fatalf("expected MultiSelectCancelledMsg, got %T", msg)
	}
}

func TestMultiSelectModel_AllCheckedDefault(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)
	selected := m.Selected()
	if len(selected) != 3 {
		t.Errorf("expected 3 selected with AllChecked=true, got %d", len(selected))
	}
}

func TestMultiSelectModel_NoneCheckedDefault(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), false)
	selected := m.Selected()
	if len(selected) != 0 {
		t.Errorf("expected 0 selected with AllChecked=false, got %d", len(selected))
	}
}

func TestMultiSelectModel_SelectedReturnsOnlyChecked(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Uncheck first item
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})

	selected := m.Selected()
	if len(selected) != 2 {
		t.Errorf("expected 2 selected, got %d", len(selected))
	}
	// Should be Beta and Gamma
	for _, s := range selected {
		if s.ID() == "a" {
			t.Error("item 'a' should not be in selected list after uncheck")
		}
	}
}

func TestMultiSelectModel_ScrollingWhenMoreThanMaxVisible(t *testing.T) {
	items := make([]msItem, 12)
	for i := range items {
		items[i] = msItem{id: string(rune('a' + i)), label: strings.Repeat("x", i+1)}
	}

	m := NewMultiSelectModel(items, &MultiSelectModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Buttons:      twoButtons(),
		AllChecked:   true,
		MaxVisible:   5,
	})
	m, _ = m.Open()

	// Navigate down past viewport
	for i := 0; i < 6; i++ {
		m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}

	if m.Cursor() != 6 {
		t.Errorf("expected cursor=6 after 6 downs, got %d", m.Cursor())
	}

	// Offset should have adjusted to keep cursor visible
	if m.offset <= 0 {
		t.Errorf("expected offset > 0 after scrolling, got %d", m.offset)
	}
}

func TestMultiSelectModel_ClosedModalIgnoresInput(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)
	m, _ = m.Close()

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if cmd != nil {
		t.Error("expected nil cmd when modal is closed")
	}
}

func TestMultiSelectModel_ViewContainsContent(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)
	view := m.View()

	if !strings.Contains(view.Content, "Test Title") {
		t.Error("expected view to contain title")
	}
	if !strings.Contains(view.Content, "Alpha") {
		t.Error("expected view to contain 'Alpha'")
	}
	if !strings.Contains(view.Content, "Beta") {
		t.Error("expected view to contain 'Beta'")
	}
	if !strings.Contains(view.Content, "Gamma") {
		t.Error("expected view to contain 'Gamma'")
	}
	if !strings.Contains(view.Content, "Update Selected") {
		t.Error("expected view to contain button label 'Update Selected'")
	}
	if !strings.Contains(view.Content, "Skip") {
		t.Error("expected view to contain button label 'Skip'")
	}
}

func TestMultiSelectModel_EnterOnListTogglesCheckbox(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Enter on list toggles first item
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected non-nil cmd (consumed)")
	}

	selected := m.Selected()
	if len(selected) != 2 {
		t.Errorf("expected 2 selected after enter toggle, got %d", len(selected))
	}
}

func TestMultiSelectModel_OverlayModal(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	bgLines := make([]string, 24)
	for i := range bgLines {
		bgLines[i] = strings.Repeat(" ", 80)
	}
	background := strings.Join(bgLines, "\n")

	view := m.OverlayModal(background)
	if view == background {
		t.Error("expected overlay to modify background when modal is open")
	}

	// Closed modal returns background unchanged
	m, _ = m.Close()
	view = m.OverlayModal(background)
	if view != background {
		t.Error("expected background unchanged when modal is closed")
	}
}

func TestMultiSelectModel_ButtonNavigationLeftRight(t *testing.T) {
	m := newTestMultiSelectModel(threeItems(), true)

	// Tab to buttons
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.buttonIdx != 0 {
		t.Errorf("expected buttonIdx=0, got %d", m.buttonIdx)
	}

	// Right moves to next button
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.buttonIdx != 1 {
		t.Errorf("expected buttonIdx=1 after right, got %d", m.buttonIdx)
	}

	// Left moves back
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.buttonIdx != 0 {
		t.Errorf("expected buttonIdx=0 after left, got %d", m.buttonIdx)
	}
}
