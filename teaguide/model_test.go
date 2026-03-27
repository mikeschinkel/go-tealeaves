package teaguide

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func testData() GuideData {
	return GuideData{
		Title: "What's Next?",
		Sections: []GuideSection{
			{
				Priority: PriorityRecommended,
				Heading:  "Recommended",
				Items: []GuideItem{
					{
						ActionKey:  "t",
						KeyDisplay: "[T]",
						Label:      "Run Tests",
						Prose:      "Tests have not been run yet.",
					},
				},
			},
			{
				Priority: PriorityAvailable,
				Heading:  "Also Available",
				Items: []GuideItem{
					{
						ActionKey:  "r",
						KeyDisplay: "[R]",
						Label:      "Refresh",
					},
					{
						ActionKey:  "q",
						KeyDisplay: "[Q]",
						Label:      "Quit",
					},
				},
			},
			{
				Priority: PriorityBlocked,
				Heading:  "Not Yet Available",
				Items: []GuideItem{
					{
						Label:       "Deploy",
						BlockReason: "run tests first",
					},
				},
			},
		},
	}
}

// keyPress creates a KeyPressMsg for a printable rune.
func keyPress(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: r, Text: string(r)}
}

// specialKeyPress creates a KeyPressMsg for a special key (no text).
func specialKeyPress(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code}
}

func TestNewGuideModel(t *testing.T) {
	m := NewGuideModel()
	if m.IsOpen() {
		t.Error("new guide should not be open")
	}
}

func TestOpenClose(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)

	m, _ = m.Open(testData())
	if !m.IsOpen() {
		t.Error("guide should be open after Open()")
	}

	m = m.Close()
	if m.IsOpen() {
		t.Error("guide should be closed after Close()")
	}
}

func TestDismissWithEsc(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	result, cmd := m.Update(specialKeyPress(tea.KeyEscape))
	updated := result.(GuideModel)

	if updated.IsOpen() {
		t.Error("guide should close on Esc")
	}
	if cmd == nil {
		t.Fatal("expected GuideDismissedMsg command")
	}
	msg := cmd()
	if _, ok := msg.(GuideDismissedMsg); !ok {
		t.Errorf("expected GuideDismissedMsg, got %T", msg)
	}
}

func TestNKeyDoesNotDismiss(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	result, cmd := m.Update(keyPress('n'))
	updated := result.(GuideModel)

	if !updated.IsOpen() {
		t.Error("guide should remain open on 'n' (only Esc dismisses)")
	}
	if cmd != nil {
		t.Error("no command should be emitted for 'n' key")
	}
}

func TestActionKeyDispatch(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	// Press "t" which is an action key
	result, cmd := m.Update(keyPress('t'))
	updated := result.(GuideModel)

	if updated.IsOpen() {
		t.Error("guide should close on action key dispatch")
	}
	if cmd == nil {
		t.Fatal("expected ActionSelectedMsg command")
	}
	msg := cmd()
	actionMsg, ok := msg.(ActionSelectedMsg)
	if !ok {
		t.Fatalf("expected ActionSelectedMsg, got %T", msg)
	}
	if actionMsg.ActionKey != "t" {
		t.Errorf("expected ActionKey 't', got %q", actionMsg.ActionKey)
	}
}

func TestBlockedKeyNotDispatched(t *testing.T) {
	// "Deploy" has no ActionKey, so pressing "d" should not dispatch
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	result, cmd := m.Update(keyPress('d'))
	updated := result.(GuideModel)

	if !updated.IsOpen() {
		t.Error("guide should remain open for unrecognized key")
	}
	if cmd != nil {
		t.Error("no command should be emitted for unrecognized key")
	}
}

func TestBlockedSectionToggle(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	if m.blockedExpanded {
		t.Error("blocked section should be collapsed by default")
	}

	// Toggle with space
	result, _ := m.Update(keyPress(' '))
	updated := result.(GuideModel)

	if !updated.blockedExpanded {
		t.Error("blocked section should be expanded after toggle")
	}

	// Toggle again
	result, _ = updated.Update(keyPress(' '))
	updated = result.(GuideModel)

	if updated.blockedExpanded {
		t.Error("blocked section should be collapsed after second toggle")
	}
}

func TestScrollUp(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	// Scroll up from 0 should stay at 0
	result, _ := m.Update(specialKeyPress(tea.KeyUp))
	updated := result.(GuideModel)

	if updated.scrollOffset != 0 {
		t.Errorf("scroll offset should stay at 0, got %d", updated.scrollOffset)
	}
}

func TestScrollDown(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	result, _ := m.Update(specialKeyPress(tea.KeyDown))
	updated := result.(GuideModel)

	if updated.scrollOffset != 1 {
		t.Errorf("scroll offset should be 1, got %d", updated.scrollOffset)
	}
}

func TestClosedModelIgnoresKeys(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)

	// Should not panic or produce commands when closed
	result, cmd := m.Update(keyPress('t'))
	updated := result.(GuideModel)

	if updated.IsOpen() {
		t.Error("closed guide should not open on key press")
	}
	if cmd != nil {
		t.Error("closed guide should not produce commands")
	}
}

func TestSetSize(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(120, 40)

	if m.screenWidth != 120 {
		t.Errorf("expected screenWidth 120, got %d", m.screenWidth)
	}
	if m.screenHeight != 40 {
		t.Errorf("expected screenHeight 40, got %d", m.screenHeight)
	}
}

func TestViewWhenClosed(t *testing.T) {
	m := NewGuideModel()
	view := m.View()
	if view.Content != "" {
		t.Errorf("closed guide should render empty, got %q", view.Content)
	}
}

func TestViewWhenOpen(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	view := m.View()
	if view.Content == "" {
		t.Error("open guide should render content")
	}
}

func TestActionMapBuiltCorrectly(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	// Should have "t", "r", "q" from Recommended + Available
	expected := []string{"t", "r", "q"}
	for _, k := range expected {
		if !m.actionMap[k] {
			t.Errorf("expected action key %q in actionMap", k)
		}
	}

	// Should NOT have blocked items' keys
	if m.actionMap[""] {
		t.Error("empty string should not be in actionMap")
	}
}

func TestOverlayModalWhenClosed(t *testing.T) {
	m := NewGuideModel()
	bg := "hello world"
	result := m.OverlayModal(bg)
	if result != bg {
		t.Error("closed guide OverlayModal should return background unchanged")
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(testData())

	result, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	updated := result.(GuideModel)

	if updated.screenWidth != 120 {
		t.Errorf("expected screenWidth 120 after resize, got %d", updated.screenWidth)
	}
	if updated.screenHeight != 40 {
		t.Errorf("expected screenHeight 40 after resize, got %d", updated.screenHeight)
	}
}
