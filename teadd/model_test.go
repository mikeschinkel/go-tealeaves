package teadd

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func testOptions() []Option {
	return []Option{
		{Text: "Alpha", Value: "a"},
		{Text: "Beta", Value: "b"},
		{Text: "Gamma", Value: "c"},
	}
}

func newTestDropdown(options []Option) DropdownModel {
	m := NewModel(options, 5, 10, &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m, _ = m.Open()
	return m
}

func extractMsg(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

// --- Layer 1 Tests ---

func TestNewModel_Defaults(t *testing.T) {
	m := NewModel(testOptions(), 5, 10, nil)
	if m.IsOpen {
		t.Error("expected IsOpen=false by default")
	}
	if m.Selected != 0 {
		t.Errorf("expected Selected=0, got %d", m.Selected)
	}
	if m.FieldRow != 5 {
		t.Errorf("expected FieldRow=5, got %d", m.FieldRow)
	}
	if m.FieldCol != 10 {
		t.Errorf("expected FieldCol=10, got %d", m.FieldCol)
	}
}

func TestNewModel_WithOptions(t *testing.T) {
	opts := testOptions()
	m := NewModel(opts, 0, 0, nil)
	if len(m.Options) != 3 {
		t.Errorf("expected 3 options, got %d", len(m.Options))
	}
}

func TestNewModel_EmptyOptions(t *testing.T) {
	m := NewModel([]Option{}, 0, 0, nil)
	if len(m.Options) != 0 {
		t.Errorf("expected 0 options, got %d", len(m.Options))
	}
}

func TestDropdownModel_Open(t *testing.T) {
	m := NewModel(testOptions(), 5, 10, &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m, _ = m.Open()
	if !m.IsOpen {
		t.Error("expected IsOpen=true after Open()")
	}
}

func TestDropdownModel_Close(t *testing.T) {
	m := newTestDropdown(testOptions())
	m, _ = m.Close()
	if m.IsOpen {
		t.Error("expected IsOpen=false after Close()")
	}
}

func TestDropdownModel_UpdateWhenClosed(t *testing.T) {
	m := NewModel(testOptions(), 5, 10, nil)
	// Not opened — closed by default
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	updated := result.(DropdownModel)
	if cmd != nil {
		t.Error("expected nil cmd when dropdown is closed")
	}
	if updated.Selected != 0 {
		t.Errorf("expected Selected unchanged (0), got %d", updated.Selected)
	}
}

func TestDropdownModel_KeyUp(t *testing.T) {
	m := newTestDropdown(testOptions())
	// Move down first
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DropdownModel)
	if m.Selected != 1 {
		t.Fatalf("expected Selected=1 after Down, got %d", m.Selected)
	}

	// Now move up
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = result.(DropdownModel)
	if m.Selected != 0 {
		t.Errorf("expected Selected=0 after Up, got %d", m.Selected)
	}
	if cmd == nil {
		t.Error("expected non-nil cmd after Up")
	}

	// Up at 0 should stay at 0
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = result.(DropdownModel)
	if m.Selected != 0 {
		t.Errorf("expected Selected=0 (clamped), got %d", m.Selected)
	}
}

func TestDropdownModel_KeyDown(t *testing.T) {
	m := newTestDropdown(testOptions())

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DropdownModel)
	if m.Selected != 1 {
		t.Errorf("expected Selected=1 after Down, got %d", m.Selected)
	}
	if cmd == nil {
		t.Error("expected non-nil cmd after Down")
	}

	// Move to last
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DropdownModel)
	if m.Selected != 2 {
		t.Errorf("expected Selected=2, got %d", m.Selected)
	}

	// Down at last should stay at last
	result, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DropdownModel)
	if m.Selected != 2 {
		t.Errorf("expected Selected=2 (clamped), got %d", m.Selected)
	}
}

func TestDropdownModel_KeySelect(t *testing.T) {
	m := newTestDropdown(testOptions())
	// Move to Beta
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DropdownModel)

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = result.(DropdownModel)

	if m.IsOpen {
		t.Error("expected dropdown closed after selection")
	}

	msg := extractMsg(cmd)
	selected, ok := msg.(OptionSelectedMsg)
	if !ok {
		t.Fatalf("expected OptionSelectedMsg, got %T", msg)
	}
	if selected.Index != 1 {
		t.Errorf("expected Index=1, got %d", selected.Index)
	}
	if selected.Text != "Beta" {
		t.Errorf("expected Text='Beta', got %q", selected.Text)
	}
	if selected.Value != "b" {
		t.Errorf("expected Value='b', got %v", selected.Value)
	}
}

func TestDropdownModel_KeyCancel(t *testing.T) {
	m := newTestDropdown(testOptions())

	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	m = result.(DropdownModel)

	if m.IsOpen {
		t.Error("expected dropdown closed after cancel")
	}

	msg := extractMsg(cmd)
	_, ok := msg.(DropdownCancelledMsg)
	if !ok {
		t.Fatalf("expected DropdownCancelledMsg, got %T", msg)
	}
}

func TestDropdownModel_WindowSizeMsg(t *testing.T) {
	m := newTestDropdown(testOptions())

	result, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = result.(DropdownModel)

	if m.ScreenWidth != 120 {
		t.Errorf("expected ScreenWidth=120, got %d", m.ScreenWidth)
	}
	if m.ScreenHeight != 40 {
		t.Errorf("expected ScreenHeight=40, got %d", m.ScreenHeight)
	}
	if cmd == nil {
		t.Error("expected non-nil cmd for handled WindowSizeMsg")
	}
}

func TestDropdownModel_ScrollOffset(t *testing.T) {
	// Create many options to force scrolling
	var manyOptions []Option
	for i := 0; i < 20; i++ {
		manyOptions = append(manyOptions, Option{Text: strings.Repeat("x", 5), Value: i})
	}
	m := NewModel(manyOptions, 5, 10, &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m, _ = m.Open()

	// Move down repeatedly to trigger scrolling
	for i := 0; i < 15; i++ {
		result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		m = result.(DropdownModel)
	}
	if m.Selected != 15 {
		t.Errorf("expected Selected=15 after 15 Down presses, got %d", m.Selected)
	}
}

func TestDropdownModel_WithPosition(t *testing.T) {
	m := NewModel(testOptions(), 0, 0, nil)
	m = m.WithPosition(10, 20)
	if m.FieldRow != 10 {
		t.Errorf("expected FieldRow=10, got %d", m.FieldRow)
	}
	if m.FieldCol != 20 {
		t.Errorf("expected FieldCol=20, got %d", m.FieldCol)
	}
}

func TestDropdownModel_WithOptions(t *testing.T) {
	m := newTestDropdown(testOptions())
	// Select second option
	result, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = result.(DropdownModel)

	// Replace with new options
	newOpts := []Option{
		{Text: "One", Value: 1},
		{Text: "Two", Value: 2},
	}
	m = m.WithOptions(newOpts)

	if len(m.Options) != 2 {
		t.Errorf("expected 2 options, got %d", len(m.Options))
	}
	if m.Selected != 1 {
		t.Errorf("expected Selected clamped to 1, got %d", m.Selected)
	}
	if m.ScrollOffset != 0 {
		t.Errorf("expected ScrollOffset reset to 0, got %d", m.ScrollOffset)
	}
}

func TestDropdownModel_WithScreenSize(t *testing.T) {
	m := NewModel(testOptions(), 0, 0, nil)
	m = m.WithScreenSize(100, 50)
	if m.ScreenWidth != 100 {
		t.Errorf("expected ScreenWidth=100, got %d", m.ScreenWidth)
	}
	if m.ScreenHeight != 50 {
		t.Errorf("expected ScreenHeight=50, got %d", m.ScreenHeight)
	}
}

// --- Layer 2 Tests ---

func TestDropdownModel_View_Closed(t *testing.T) {
	m := NewModel(testOptions(), 0, 0, nil)
	view := m.View()
	if view.Content != "" {
		t.Errorf("expected empty view when closed, got %q", view.Content)
	}
}

func TestDropdownModel_View_Open(t *testing.T) {
	m := newTestDropdown(testOptions())
	view := m.View()
	if view.Content == "" {
		t.Error("expected non-empty view when open")
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
}
