package teamodal

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// testItem implements ListItem for testing
type testItem struct {
	id     string
	label  string
	active bool
}

func (ti testItem) ID() string      { return ti.id }
func (ti testItem) Label() string   { return ti.label }
func (ti testItem) IsActive() bool  { return ti.active }

func testItems() []testItem {
	return []testItem{
		{id: "1", label: "Alpha", active: false},
		{id: "2", label: "Beta", active: true},
		{id: "3", label: "Gamma", active: false},
		{id: "4", label: "Delta", active: false},
	}
}

func newTestListModel(items []testItem) ListModel[testItem] {
	m := NewListModel(items, &ListModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Test List",
		MaxVisible:   8,
	})
	m = m.Open()
	return m
}

// --- Layer 1 ---

func TestNewListModel(t *testing.T) {
	items := testItems()
	m := NewListModel(items, &ListModelArgs{Title: "Test"})
	if len(m.Items()) != 4 {
		t.Errorf("expected 4 items, got %d", len(m.Items()))
	}
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0, got %d", m.Cursor())
	}
	if m.IsOpen() {
		t.Error("expected not open initially")
	}
}

func TestListModel_Open_FocusesActiveItem(t *testing.T) {
	items := testItems() // Item at index 1 is active
	m := NewListModel(items, &ListModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m = m.Open()
	if m.Cursor() != 1 {
		t.Errorf("expected cursor=1 (active item), got %d", m.Cursor())
	}
}

func TestListModel_Open_NoActiveItem(t *testing.T) {
	items := []testItem{
		{id: "1", label: "One"},
		{id: "2", label: "Two"},
	}
	m := NewListModel(items, &ListModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m = m.Open()
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 when no active item, got %d", m.Cursor())
	}
}

func TestListModel_KeyUp(t *testing.T) {
	m := newTestListModel(testItems())
	// Cursor starts at 1 (active item). Move up.
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 after Up, got %d", m.Cursor())
	}
	if cmd == nil {
		t.Error("expected non-nil cmd")
	}

	// At top, Up should stay at 0
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 (clamped), got %d", m.Cursor())
	}
}

func TestListModel_KeyDown(t *testing.T) {
	m := newTestListModel(testItems())
	// Start at 1, move down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.Cursor() != 2 {
		t.Errorf("expected cursor=2, got %d", m.Cursor())
	}

	// Move to last
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.Cursor() != 3 {
		t.Errorf("expected cursor=3, got %d", m.Cursor())
	}

	// At bottom, stays at last
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.Cursor() != 3 {
		t.Errorf("expected cursor=3 (clamped), got %d", m.Cursor())
	}
}

func TestListModel_KeySpace(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeySpace})

	// Modal should stay open
	if !m.IsOpen() {
		t.Error("expected modal to stay open after Space")
	}

	msg := extractMsg(cmd)
	selected, ok := msg.(ItemSelectedMsg[testItem])
	if !ok {
		t.Fatalf("expected ItemSelectedMsg, got %T", msg)
	}
	if selected.Item.ID() != "2" {
		t.Errorf("expected item ID='2' (cursor item), got %q", selected.Item.ID())
	}
}

func TestListModel_KeyEnter(t *testing.T) {
	m := newTestListModel(testItems())
	// Cursor at active item (1), Enter should close
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// cmd is tea.Batch, execute it
	if cmd == nil {
		t.Fatal("expected non-nil cmd from Enter")
	}
}

func TestListModel_KeyNew(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	msg := extractMsg(cmd)
	_, ok := msg.(NewItemRequestedMsg)
	if !ok {
		t.Fatalf("expected NewItemRequestedMsg, got %T", msg)
	}
	// Modal should stay open
	if !m.IsOpen() {
		t.Error("expected modal to stay open after 'a'")
	}
}

func TestListModel_KeyEdit(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	if !m.isEditing {
		t.Error("expected isEditing=true after 'e'")
	}
	if m.editBuffer != "Beta" {
		t.Errorf("expected editBuffer='Beta', got %q", m.editBuffer)
	}
}

func TestListModel_KeyDelete(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	msg := extractMsg(cmd)
	del, ok := msg.(DeleteItemRequestedMsg[testItem])
	if !ok {
		t.Fatalf("expected DeleteItemRequestedMsg, got %T", msg)
	}
	if del.Item.ID() != "2" {
		t.Errorf("expected deleted item ID='2', got %q", del.Item.ID())
	}
}

func TestListModel_KeyHelp(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !m.showHelp {
		t.Error("expected showHelp=true after '?'")
	}

	// Toggle off
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if m.showHelp {
		t.Error("expected showHelp=false after second '?'")
	}
}

func TestListModel_HelpVisibleEscClosesHelp(t *testing.T) {
	m := newTestListModel(testItems())
	// Open help
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	if !m.showHelp {
		t.Fatal("expected help visible")
	}

	// Esc closes help, not the modal
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.showHelp {
		t.Error("expected help closed after Esc")
	}
	if !m.IsOpen() {
		t.Error("expected modal to remain open (Esc only closed help)")
	}
}

func TestListModel_KeyCancel(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if m.IsOpen() {
		t.Error("expected modal closed after Esc")
	}
	msg := extractMsg(cmd)
	_, ok := msg.(ListCancelledMsg)
	if !ok {
		t.Fatalf("expected ListCancelledMsg, got %T", msg)
	}
}

func TestListModel_EditEnterCompletes(t *testing.T) {
	m := newTestListModel(testItems())
	// Enter edit mode
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	// Type new text (first keystroke overwrites)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}})
	// Complete edit
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.isEditing {
		t.Error("expected isEditing=false after Enter")
	}
	msg := extractMsg(cmd)
	edit, ok := msg.(EditCompletedMsg[testItem])
	if !ok {
		t.Fatalf("expected EditCompletedMsg, got %T", msg)
	}
	if edit.NewLabel != "New" {
		t.Errorf("expected NewLabel='New', got %q", edit.NewLabel)
	}
}

func TestListModel_EditEscCancels(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	// Type something
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})
	// Cancel
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if m.isEditing {
		t.Error("expected isEditing=false after Esc")
	}
	if m.editBuffer != "" {
		t.Errorf("expected editBuffer cleared, got %q", m.editBuffer)
	}
}

func TestListModel_EditFirstKeystrokeOverwrites(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	// editOverwrite should be true
	if !m.editOverwrite {
		t.Error("expected editOverwrite=true initially")
	}
	// First keystroke
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})
	if m.editBuffer != "X" {
		t.Errorf("expected editBuffer='X' after overwrite, got %q", m.editBuffer)
	}
	if m.editOverwrite {
		t.Error("expected editOverwrite=false after first keystroke")
	}
}

func TestListModel_EditSubsequentInsertion(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	// First keystroke overwrites
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}})
	// Subsequent inserts
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'B'}})
	if m.editBuffer != "AB" {
		t.Errorf("expected editBuffer='AB', got %q", m.editBuffer)
	}
}

func TestListModel_EditCursorMovement(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	// Disable overwrite by moving cursor
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	if m.editOverwrite {
		t.Error("expected editOverwrite=false after cursor move")
	}

	// editBuffer should be "Beta", cursor at 1
	// Backspace at position 1
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	if m.editBuffer != "eta" {
		t.Errorf("expected 'eta' after backspace, got %q", m.editBuffer)
	}
}

func TestListModel_SetItems(t *testing.T) {
	m := newTestListModel(testItems())
	// Move cursor to position 3
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Replace with fewer items
	newItems := []testItem{
		{id: "a", label: "One"},
		{id: "b", label: "Two"},
	}
	m = m.SetItems(newItems)

	if len(m.Items()) != 2 {
		t.Errorf("expected 2 items, got %d", len(m.Items()))
	}
	// Cursor should be clamped
	if m.Cursor() > 1 {
		t.Errorf("expected cursor clamped to <=1, got %d", m.Cursor())
	}
}

func TestListModel_SetCursor(t *testing.T) {
	m := newTestListModel(testItems())
	m = m.SetCursor(2)
	if m.Cursor() != 2 {
		t.Errorf("expected cursor=2, got %d", m.Cursor())
	}

	// Out of bounds clamped
	m = m.SetCursor(99)
	if m.Cursor() != 3 {
		t.Errorf("expected cursor=3 (clamped to last), got %d", m.Cursor())
	}

	m = m.SetCursor(-1)
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 (clamped to first), got %d", m.Cursor())
	}
}

func TestListModel_SetCursorToLast(t *testing.T) {
	m := newTestListModel(testItems())
	m = m.SetCursorToLast()
	if m.Cursor() != 3 {
		t.Errorf("expected cursor=3, got %d", m.Cursor())
	}
}

// --- Layer 2 ---

func TestListModel_View_Open(t *testing.T) {
	m := newTestListModel(testItems())
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view when open")
	}
	if !strings.Contains(view, "Test List") {
		t.Error("expected view to contain title")
	}
	if !strings.Contains(view, "Alpha") {
		t.Error("expected view to contain item label 'Alpha'")
	}
}

func TestListModel_View_ActiveItem(t *testing.T) {
	m := newTestListModel(testItems())
	view := m.View()
	if !strings.Contains(view, "ACTIVE") {
		t.Error("expected view to show ACTIVE indicator for active item")
	}
}

func TestListModel_View_HelpVisor(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	view := m.View()
	if !strings.Contains(view, "Keyboard Shortcuts") {
		t.Error("expected help visor to contain 'Keyboard Shortcuts'")
	}
}
