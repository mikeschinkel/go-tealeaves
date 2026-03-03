package teamodal

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
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
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 after Up, got %d", m.Cursor())
	}
	if cmd == nil {
		t.Error("expected non-nil cmd")
	}

	// At top, Up should stay at 0
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 (clamped), got %d", m.Cursor())
	}
}

func TestListModel_KeyDown(t *testing.T) {
	m := newTestListModel(testItems())
	// Start at 1, move down
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.Cursor() != 2 {
		t.Errorf("expected cursor=2, got %d", m.Cursor())
	}

	// Move to last
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.Cursor() != 3 {
		t.Errorf("expected cursor=3, got %d", m.Cursor())
	}

	// At bottom, stays at last
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.Cursor() != 3 {
		t.Errorf("expected cursor=3 (clamped), got %d", m.Cursor())
	}
}

func TestListModel_KeySpace(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})

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
	_, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	// cmd is tea.Batch, execute it
	if cmd == nil {
		t.Fatal("expected non-nil cmd from Enter")
	}
}

func TestListModel_KeyNew(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})

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
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})

	if !m.isEditing {
		t.Error("expected isEditing=true after 'e'")
	}
	if m.editBuffer != "Beta" {
		t.Errorf("expected editBuffer='Beta', got %q", m.editBuffer)
	}
}

func TestListModel_KeyDelete(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyPressMsg{Code: 'd', Text: "d"})

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
	m, _ = m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	if !m.showHelp {
		t.Error("expected showHelp=true after '?'")
	}

	// Toggle off
	m, _ = m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	if m.showHelp {
		t.Error("expected showHelp=false after second '?'")
	}
}

func TestListModel_HelpVisibleEscClosesHelp(t *testing.T) {
	m := newTestListModel(testItems())
	// Open help
	m, _ = m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	if !m.showHelp {
		t.Fatal("expected help visible")
	}

	// Esc closes help, not the modal
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if m.showHelp {
		t.Error("expected help closed after Esc")
	}
	if !m.IsOpen() {
		t.Error("expected modal to remain open (Esc only closed help)")
	}
}

func TestListModel_KeyCancel(t *testing.T) {
	m := newTestListModel(testItems())
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})

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
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	// Type new text (first keystroke overwrites)
	m, _ = m.Update(tea.KeyPressMsg{Code: 'N', Text: "N"})
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	m, _ = m.Update(tea.KeyPressMsg{Code: 'w', Text: "w"})
	// Complete edit
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

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
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	// Type something
	m, _ = m.Update(tea.KeyPressMsg{Code: 'X', Text: "X"})
	// Cancel
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})

	if m.isEditing {
		t.Error("expected isEditing=false after Esc")
	}
	if m.editBuffer != "" {
		t.Errorf("expected editBuffer cleared, got %q", m.editBuffer)
	}
}

func TestListModel_EditFirstKeystrokeOverwrites(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	// editOverwrite should be true
	if !m.editOverwrite {
		t.Error("expected editOverwrite=true initially")
	}
	// First keystroke
	m, _ = m.Update(tea.KeyPressMsg{Code: 'X', Text: "X"})
	if m.editBuffer != "X" {
		t.Errorf("expected editBuffer='X' after overwrite, got %q", m.editBuffer)
	}
	if m.editOverwrite {
		t.Error("expected editOverwrite=false after first keystroke")
	}
}

func TestListModel_EditSubsequentInsertion(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	// First keystroke overwrites
	m, _ = m.Update(tea.KeyPressMsg{Code: 'A', Text: "A"})
	// Subsequent inserts
	m, _ = m.Update(tea.KeyPressMsg{Code: 'B', Text: "B"})
	if m.editBuffer != "AB" {
		t.Errorf("expected editBuffer='AB', got %q", m.editBuffer)
	}
}

func TestListModel_EditCursorMovement(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	// Disable overwrite by moving cursor
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.editOverwrite {
		t.Error("expected editOverwrite=false after cursor move")
	}

	// editBuffer should be "Beta", cursor at 1
	// Backspace at position 1
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if m.editBuffer != "eta" {
		t.Errorf("expected 'eta' after backspace, got %q", m.editBuffer)
	}
}

// --- Migration-sensitive tests (v1→v2 regression guards) ---

// LST-EDIT-SPACE: Guards tea.KeySpace branch in updateEditing (list_model.go:372-384)
// Space in edit mode must insert a space character into the edit buffer.
func TestListModel_EditSpaceInsertion(t *testing.T) {
	m := newTestListModel(testItems())
	// Enter edit mode (cursor on "Beta")
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	if !m.isEditing {
		t.Fatal("expected isEditing=true")
	}

	// First keystroke overwrites: type "A"
	m, _ = m.Update(tea.KeyPressMsg{Code: 'A', Text: "A"})
	if m.editBuffer != "A" {
		t.Fatalf("expected editBuffer='A', got %q", m.editBuffer)
	}

	// Space should insert at cursor position (after "A")
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})
	if m.editBuffer != "A " {
		t.Errorf("expected editBuffer='A ', got %q", m.editBuffer)
	}
	if m.editCursor != 2 {
		t.Errorf("expected editCursor=2 after space, got %d", m.editCursor)
	}

	// Type "B" after space
	m, _ = m.Update(tea.KeyPressMsg{Code: 'B', Text: "B"})
	if m.editBuffer != "A B" {
		t.Errorf("expected editBuffer='A B', got %q", m.editBuffer)
	}
}

// LST-EDIT-DELETE: Guards keyMsg.Type == tea.KeyDelete branch in updateEditing (list_model.go:342-350)
// Delete key in edit mode must remove the character at cursor position.
func TestListModel_EditDeleteKey(t *testing.T) {
	m := newTestListModel(testItems())
	// Enter edit mode (cursor on "Beta")
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	if !m.isEditing {
		t.Fatal("expected isEditing=true")
	}

	// Move cursor right to disable overwrite, cursor is at position 0 in "Beta"
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.editOverwrite {
		t.Fatal("expected editOverwrite=false after cursor move")
	}
	// Now cursor=1 within "Beta"

	// Move cursor to position 0 (left)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.editCursor != 0 {
		t.Fatalf("expected editCursor=0, got %d", m.editCursor)
	}

	// Delete at position 0 should remove 'B' → "eta"
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDelete})
	if m.editBuffer != "eta" {
		t.Errorf("expected editBuffer='eta' after delete, got %q", m.editBuffer)
	}
	if m.editCursor != 0 {
		t.Errorf("expected editCursor=0 (unchanged after delete), got %d", m.editCursor)
	}

	// Delete again at position 0 should remove 'e' → "ta"
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDelete})
	if m.editBuffer != "ta" {
		t.Errorf("expected editBuffer='ta' after second delete, got %q", m.editBuffer)
	}
}

func TestListModel_SetItems(t *testing.T) {
	m := newTestListModel(testItems())
	// Move cursor to position 3
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})

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
	if view.Content == "" {
		t.Error("expected non-empty view when open")
	}
	if !strings.Contains(view.Content, "Test List") {
		t.Error("expected view to contain title")
	}
	if !strings.Contains(view.Content, "Alpha") {
		t.Error("expected view to contain item label 'Alpha'")
	}
}

func TestListModel_View_ActiveItem(t *testing.T) {
	m := newTestListModel(testItems())
	view := m.View()
	if !strings.Contains(view.Content, "ACTIVE") {
		t.Error("expected view to show ACTIVE indicator for active item")
	}
}

func TestListModel_Scrolling(t *testing.T) {
	// Create more items than maxVisible (3)
	manyItems := []testItem{
		{id: "1", label: "Item-1"},
		{id: "2", label: "Item-2"},
		{id: "3", label: "Item-3"},
		{id: "4", label: "Item-4"},
		{id: "5", label: "Item-5"},
		{id: "6", label: "Item-6"},
	}
	m := NewListModel(manyItems, &ListModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Scroll Test",
		MaxVisible:   3,
	})
	m = m.Open()

	// Initially offset should be 0, cursor at 0
	if m.Offset() != 0 {
		t.Errorf("expected offset=0 initially, got %d", m.Offset())
	}

	// Move cursor down past maxVisible
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	// cursor=3, should trigger scroll (offset should adjust)
	if m.Cursor() != 3 {
		t.Errorf("expected cursor=3, got %d", m.Cursor())
	}
	if m.Offset() == 0 {
		t.Error("expected offset > 0 after scrolling past viewport")
	}

	// Continue to bottom
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if m.Cursor() != 5 {
		t.Errorf("expected cursor=5, got %d", m.Cursor())
	}

	// Now scroll back up
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if m.Cursor() != 0 {
		t.Errorf("expected cursor=0 after scrolling up, got %d", m.Cursor())
	}
	if m.Offset() != 0 {
		t.Errorf("expected offset=0 at top, got %d", m.Offset())
	}
}

func TestListModel_View_SelectedItem(t *testing.T) {
	m := newTestListModel(testItems())
	// Cursor is on Beta (active item, index 1)
	view := m.View()

	// Should have the cursor indicator character
	if !strings.Contains(view.Content, "\u25b6") {
		t.Error("expected view to contain cursor indicator (▶)")
	}

	// Move cursor to Gamma (index 2)
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	view = m.View()
	if !strings.Contains(view.Content, "Gamma") {
		t.Error("expected view to contain 'Gamma'")
	}
}

func TestListModel_View_EditMode(t *testing.T) {
	m := newTestListModel(testItems())
	// Enter edit mode
	m, _ = m.Update(tea.KeyPressMsg{Code: 'e', Text: "e"})
	if !m.isEditing {
		t.Fatal("expected isEditing=true")
	}

	view := m.View()
	if view.Content == "" {
		t.Fatal("expected non-empty view in edit mode")
	}

	// The footer should change to show edit-mode keys
	if !strings.Contains(view.Content, "accept") {
		t.Error("expected edit-mode footer with 'accept' hint")
	}
}

func TestListModel_View_StatusMessage(t *testing.T) {
	m := newTestListModel(testItems())
	m = m.SetStatus("Item saved")
	view := m.View()

	if !strings.Contains(view.Content, "Item saved") {
		t.Error("expected view to contain status message 'Item saved'")
	}

	// Clear status
	m = m.ClearStatus()
	view = m.View()
	if strings.Contains(view.Content, "Item saved") {
		t.Error("expected status message cleared from view")
	}
}

func TestListModel_View_HelpVisor(t *testing.T) {
	m := newTestListModel(testItems())
	m, _ = m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	view := m.View()
	if !strings.Contains(view.Content, "Keyboard Shortcuts") {
		t.Error("expected help visor to contain 'Keyboard Shortcuts'")
	}
}
