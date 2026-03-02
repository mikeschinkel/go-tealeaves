package teamodal

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestOKModal() ModalModel {
	m := NewOKModal("Test alert message", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m, _ = m.Open()
	return m
}

func newTestYesNoModal() ModalModel {
	m := NewYesNoModal("Are you sure?", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		DefaultYes:   true,
	})
	m, _ = m.Open()
	return m
}

// --- Layer 1: OK Modal ---

func TestNewOKModal(t *testing.T) {
	m := NewOKModal("Alert!", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	if m.Type() != ModalTypeOK {
		t.Errorf("expected ModalTypeOK, got %d", m.Type())
	}
	if m.Message() != "Alert!" {
		t.Errorf("expected message='Alert!', got %q", m.Message())
	}
	if m.IsOpen() {
		t.Error("expected modal not open initially")
	}
}

func TestNewYesNoModal(t *testing.T) {
	m := NewYesNoModal("Continue?", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		DefaultYes:   true,
	})
	if m.Type() != ModalTypeYesNo {
		t.Errorf("expected ModalTypeYesNo, got %d", m.Type())
	}
	if m.FocusButton() != 0 {
		t.Errorf("expected focusButton=0 (Yes) with DefaultYes, got %d", m.FocusButton())
	}
}

func TestNewYesNoModal_DefaultNo(t *testing.T) {
	m := NewYesNoModal("Continue?", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		DefaultYes:   false,
	})
	if m.FocusButton() != 1 {
		t.Errorf("expected focusButton=1 (No) without DefaultYes, got %d", m.FocusButton())
	}
}

func TestModalModel_Open(t *testing.T) {
	m := NewOKModal("Test", &ModelArgs{ScreenWidth: 80, ScreenHeight: 24})
	m, _ = m.Open()
	if !m.IsOpen() {
		t.Error("expected IsOpen=true after Open()")
	}
}

func TestModalModel_Close(t *testing.T) {
	m := newTestOKModal()
	m, _ = m.Close()
	if m.IsOpen() {
		t.Error("expected IsOpen=false after Close()")
	}
}

func TestModalModel_SetSize(t *testing.T) {
	m := NewOKModal("Test", nil)
	m = m.SetSize(120, 40)
	if m.ScreenWidth() != 120 {
		t.Errorf("expected ScreenWidth=120, got %d", m.ScreenWidth())
	}
	if m.ScreenHeight() != 40 {
		t.Errorf("expected ScreenHeight=40, got %d", m.ScreenHeight())
	}
}

func TestOKModal_EnterClosesAndSendsClosedMsg(t *testing.T) {
	m := newTestOKModal()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(ModalModel)

	if m.IsOpen() {
		t.Error("expected modal closed after Enter")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(ClosedMsg)
	if !ok {
		t.Fatalf("expected ClosedMsg, got %T", msg)
	}
}

func TestOKModal_EscClosesAndSendsClosedMsg(t *testing.T) {
	m := newTestOKModal()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = result.(ModalModel)

	if m.IsOpen() {
		t.Error("expected modal closed after Esc")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(ClosedMsg)
	if !ok {
		t.Fatalf("expected ClosedMsg, got %T", msg)
	}
}

// --- Layer 1: YesNo Modal ---

func TestYesNoModal_TabTogglesFocus(t *testing.T) {
	m := newTestYesNoModal()
	if m.FocusButton() != 0 {
		t.Fatalf("expected focus=0 initially, got %d", m.FocusButton())
	}

	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = result.(ModalModel)
	if m.FocusButton() != 1 {
		t.Errorf("expected focus=1 after Tab, got %d", m.FocusButton())
	}

	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = result.(ModalModel)
	if m.FocusButton() != 0 {
		t.Errorf("expected focus=0 after second Tab, got %d", m.FocusButton())
	}
}

func TestYesNoModal_EnterOnYes(t *testing.T) {
	m := newTestYesNoModal() // focus=0 (Yes)
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(ModalModel)

	if m.IsOpen() {
		t.Error("expected modal closed")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(AnsweredYesMsg)
	if !ok {
		t.Fatalf("expected AnsweredYesMsg, got %T", msg)
	}
}

func TestYesNoModal_EnterOnNo(t *testing.T) {
	m := newTestYesNoModal()
	// Move to No
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m = result.(ModalModel)

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(ModalModel)

	if m.IsOpen() {
		t.Error("expected modal closed")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(AnsweredNoMsg)
	if !ok {
		t.Fatalf("expected AnsweredNoMsg, got %T", msg)
	}
}

func TestYesNoModal_EscSendsAnsweredNo(t *testing.T) {
	m := newTestYesNoModal()
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = result.(ModalModel)

	if m.IsOpen() {
		t.Error("expected modal closed")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(AnsweredNoMsg)
	if !ok {
		t.Fatalf("expected AnsweredNoMsg, got %T", msg)
	}
}

func TestYesNoModal_ArrowKeysFocus(t *testing.T) {
	m := newTestYesNoModal()

	// Right moves to No
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	m = result.(ModalModel)
	if m.FocusButton() != 1 {
		t.Errorf("expected focus=1 after Right, got %d", m.FocusButton())
	}

	// Left moves to Yes
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	m = result.(ModalModel)
	if m.FocusButton() != 0 {
		t.Errorf("expected focus=0 after Left, got %d", m.FocusButton())
	}
}

func TestYesNoModal_MouseClickYes(t *testing.T) {
	m := newTestYesNoModal()

	// Approximate button position — Yes button is near the center
	result, cmd := m.Update(tea.MouseClickMsg{
		Button: tea.MouseLeft,
		X:      m.lastCol + m.width/2 - 5,
		Y:      m.lastRow + m.height - 3,
	})
	m = result.(ModalModel)

	if m.IsOpen() {
		t.Error("expected modal closed after mouse click")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(AnsweredYesMsg)
	if !ok {
		t.Fatalf("expected AnsweredYesMsg from click, got %T", msg)
	}
}

func TestYesNoModal_MouseClickNo(t *testing.T) {
	m := newTestYesNoModal()

	// Click on the No button (right side of button row)
	result, cmd := m.Update(tea.MouseClickMsg{
		Button: tea.MouseLeft,
		X:      m.lastCol + m.width/2 + 5,
		Y:      m.lastRow + m.height - 3,
	})
	m = result.(ModalModel)

	if m.IsOpen() {
		t.Error("expected modal closed after mouse click")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(AnsweredNoMsg)
	if !ok {
		t.Fatalf("expected AnsweredNoMsg from click, got %T", msg)
	}
}

// --- Migration-sensitive tests (v1→v2 regression guards) ---

// MOD-MOUSE-MOTION: Guards switch mouseMsg.Type + tea.MouseMotion (model.go:279)
// Mouse motion over a YesNo modal button should update focus without closing.
func TestYesNoModal_MouseMotionHover(t *testing.T) {
	m := newTestYesNoModal()
	if m.FocusButton() != 0 {
		t.Fatalf("expected initial focus=0 (Yes), got %d", m.FocusButton())
	}

	// Send mouse motion over the No button area (right side of button row)
	result, _ := m.Update(tea.MouseMotionMsg{
		X: m.lastCol + m.width/2 + 5,
		Y: m.lastRow + m.height - 3,
	})
	m = result.(ModalModel)

	// Modal should remain open
	if !m.IsOpen() {
		t.Error("expected modal to remain open after mouse motion")
	}
	// Focus should shift to the No button (1)
	if m.FocusButton() != 1 {
		t.Errorf("expected focus=1 (No) after hover, got %d", m.FocusButton())
	}
}

// MOD-MOUSE-TYPE: Guards mouseMsg.Type == tea.MouseLeft + bounds check (model.go:269-277)
// A mouse click that misses the button row should NOT close the modal.
func TestOKModal_MouseClickMiss(t *testing.T) {
	m := newTestOKModal()

	// Click well outside the button row (Y=0 is above the modal)
	result, cmd := m.Update(tea.MouseClickMsg{
		Button: tea.MouseLeft,
		X:      0,
		Y:      0,
	})
	m = result.(ModalModel)

	if !m.IsOpen() {
		t.Error("expected modal to remain open when click misses buttons")
	}
	if cmd != nil {
		msg := extractMsg(cmd)
		if msg != nil {
			t.Errorf("expected no message from missed click, got %T", msg)
		}
	}
}

func TestModalModel_ClosedIgnoresInput(t *testing.T) {
	m := newTestOKModal()
	m, _ = m.Close()

	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected nil cmd when modal is closed")
	}
}

// --- Layer 2 ---

func TestModalModel_View_Closed(t *testing.T) {
	m := NewOKModal("Test", nil)
	view := m.View()
	if view.Content != "" {
		t.Errorf("expected empty view when closed, got %q", view.Content)
	}
}

func TestModalModel_View_OKOpen(t *testing.T) {
	m := newTestOKModal()
	view := m.View()
	if !strings.Contains(view.Content, "Test alert message") {
		t.Error("expected view to contain message text")
	}
	if !strings.Contains(view.Content, "OK") {
		t.Error("expected view to contain OK button label")
	}
}

func TestModalModel_View_YesNoOpen(t *testing.T) {
	m := newTestYesNoModal()
	view := m.View()
	if !strings.Contains(view.Content, "Are you sure?") {
		t.Error("expected view to contain message text")
	}
	if !strings.Contains(view.Content, "Yes") {
		t.Error("expected view to contain Yes label")
	}
	if !strings.Contains(view.Content, "No") {
		t.Error("expected view to contain No label")
	}
}

func TestModalModel_OverlayModal(t *testing.T) {
	m := newTestOKModal()

	bgLines := make([]string, 24)
	for i := range bgLines {
		bgLines[i] = strings.Repeat(" ", 80)
	}
	background := strings.Join(bgLines, "\n")

	view := m.OverlayModal(background)
	if view == background {
		t.Error("expected overlay to modify background when modal is open")
	}

	// Closed should return background unchanged
	m, _ = m.Close()
	view = m.OverlayModal(background)
	if view != background {
		t.Error("expected background unchanged when modal is closed")
	}
}

func TestModalModel_CustomLabels(t *testing.T) {
	m := NewYesNoModal("Delete?", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		YesLabel:     "Delete",
		NoLabel:      "Keep",
	})
	m, _ = m.Open()

	view := m.View()
	if !strings.Contains(view.Content, "Delete") {
		t.Error("expected custom Yes label 'Delete' in view")
	}
	if !strings.Contains(view.Content, "Keep") {
		t.Error("expected custom No label 'Keep' in view")
	}
}

func TestModalModel_CustomLabels_Withers(t *testing.T) {
	m := NewOKModal("Test", &ModelArgs{ScreenWidth: 80, ScreenHeight: 24})
	m = m.WithOKLabel("Got it")
	m, _ = m.Open()

	view := m.View()
	if !strings.Contains(view.Content, "Got it") {
		t.Error("expected custom OK label 'Got it' in view")
	}
}
