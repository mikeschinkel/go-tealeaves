package teatxtsnip

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

func newTestModel() Model {
	m := NewTextSnipModel(nil)
	m.Model.SetWidth(80)
	m.Model.SetHeight(10)
	m.Model.SetValue("Hello World\nSecond line\nThird line")
	m.Model.Focus()
	// Move cursor to start
	m.moveCursorTo(0, 0)
	return m
}

// --- Layer 1: Model Tests ---

func TestNewTextSnipModel(t *testing.T) {
	m := NewTextSnipModel(nil)
	if m.HasSelection() {
		t.Error("expected no selection initially")
	}
	if m.IsSingleLine() {
		t.Error("expected IsSingleLine=false for NewTextSnipModel(nil)")
	}
}

func TestNewTextSnipModel_SingleLine(t *testing.T) {
	m := NewTextSnipModel(&TextSnipModelArgs{SingleLine: true})
	if !m.IsSingleLine() {
		t.Error("expected IsSingleLine=true")
	}
}

func TestNewTextSnipModel_FromTextarea(t *testing.T) {
	ta := textarea.New()
	ta.SetValue("existing content")
	m := NewTextSnipModel(&TextSnipModelArgs{Textarea: &ta})

	if m.Value() != "existing content" {
		t.Errorf("expected value='existing content', got %q", m.Value())
	}
	if m.HasSelection() {
		t.Error("expected no selection on new model")
	}
}

func TestModel_SelectLeft(t *testing.T) {
	m := newTestModel()
	// Move cursor to position (0, 5) - after "Hello"
	m.moveCursorTo(0, 5)

	m = m.extendSelectionLeft(1)

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	// Selection should extend left by 1
	if sel.End.Col != 4 {
		t.Errorf("expected End.Col=4 after extend left, got %d", sel.End.Col)
	}
}

func TestModel_SelectRight(t *testing.T) {
	m := newTestModel()
	// Cursor at (0, 0)
	m = m.extendSelectionRight(1)

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	if sel.End.Col != 1 {
		t.Errorf("expected End.Col=1 after extend right, got %d", sel.End.Col)
	}
}

func TestModel_SelectUp(t *testing.T) {
	m := newTestModel()
	// Move to second line
	m.moveCursorTo(1, 3)

	m = m.extendSelectionUp()

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	if sel.End.Row != 0 {
		t.Errorf("expected End.Row=0 after extend up, got %d", sel.End.Row)
	}
}

func TestModel_SelectDown(t *testing.T) {
	m := newTestModel()
	// Cursor at (0, 3)
	m.moveCursorTo(0, 3)

	m = m.extendSelectionDown()

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	if sel.End.Row != 1 {
		t.Errorf("expected End.Row=1 after extend down, got %d", sel.End.Row)
	}
}

func TestModel_SelectWordLeft(t *testing.T) {
	m := newTestModel()
	// Move cursor to end of "World" (0, 11)
	m.moveCursorTo(0, 11)

	m = m.extendSelectionWordLeft()

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	// Should have moved to start of "World" (col 6)
	if sel.End.Col != 6 {
		t.Errorf("expected End.Col=6 after word left, got %d", sel.End.Col)
	}
}

func TestModel_SelectWordRight(t *testing.T) {
	m := newTestModel()
	// Cursor at (0, 0)
	m = m.extendSelectionWordRight()

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	// Should have moved past "Hello " to col 6
	if sel.End.Col != 6 {
		t.Errorf("expected End.Col=6 after word right, got %d", sel.End.Col)
	}
}

func TestModel_SelectToLineStart(t *testing.T) {
	m := newTestModel()
	m.moveCursorTo(0, 5)

	m = m.extendSelectionToLineStart()

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	if sel.End.Col != 0 {
		t.Errorf("expected End.Col=0, got %d", sel.End.Col)
	}
}

func TestModel_SelectToLineEnd(t *testing.T) {
	m := newTestModel()

	m = m.extendSelectionToLineEnd()

	sel := m.Selection()
	if !sel.Active {
		t.Error("expected selection active")
	}
	// "Hello World" has 11 chars
	if sel.End.Col != 11 {
		t.Errorf("expected End.Col=11, got %d", sel.End.Col)
	}
}

func TestModel_SelectToStart(t *testing.T) {
	m := newTestModel()
	m.moveCursorTo(1, 5)

	m = m.extendSelectionToStart()

	sel := m.Selection()
	if sel.End.Row != 0 || sel.End.Col != 0 {
		t.Errorf("expected End at {0,0}, got {%d,%d}", sel.End.Row, sel.End.Col)
	}
}

func TestModel_SelectToEnd(t *testing.T) {
	m := newTestModel()

	m = m.extendSelectionToEnd()

	sel := m.Selection()
	// Last line is "Third line" (10 chars), row 2
	if sel.End.Row != 2 {
		t.Errorf("expected End.Row=2, got %d", sel.End.Row)
	}
	if sel.End.Col != 10 {
		t.Errorf("expected End.Col=10, got %d", sel.End.Col)
	}
}

func TestModel_SelectAll(t *testing.T) {
	m := newTestModel()
	m = m.selectAll()

	if !m.HasSelection() {
		t.Error("expected HasSelection=true after SelectAll")
	}
	sel := m.Selection()
	if sel.Start.Row != 0 || sel.Start.Col != 0 {
		t.Error("expected Start at {0,0}")
	}
	if sel.End.Row != 2 {
		t.Errorf("expected End.Row=2, got %d", sel.End.Row)
	}
}

func TestModel_ClearSelection(t *testing.T) {
	m := newTestModel()
	m = m.selectAll()
	m = m.ClearSelection()

	if m.HasSelection() {
		t.Error("expected HasSelection=false after ClearSelection")
	}
}

func TestModel_CursorMovementClearsSelection(t *testing.T) {
	m := newTestModel()
	// Create a selection
	m = m.extendSelectionRight(3)
	if !m.HasSelection() {
		t.Fatal("expected selection before cursor move")
	}

	// Arrow key without shift should clear selection
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.HasSelection() {
		t.Error("expected selection cleared after cursor movement")
	}
}

func TestModel_TypingReplacesSelection(t *testing.T) {
	m := newTestModel()
	// Select "Hello"
	m = m.extendSelectionRight(5)
	if !m.HasSelection() {
		t.Fatal("expected selection before typing")
	}

	// Type a character - should delete selection and insert
	m, _ = m.Update(tea.KeyPressMsg{Code: 'X', Text: "X"})

	if m.HasSelection() {
		t.Error("expected selection cleared after typing")
	}
	// "Hello" should be replaced with "X"
	value := m.Value()
	if len(value) < 1 || value[0] != 'X' {
		t.Errorf("expected value to start with 'X', got %q", value)
	}
}

func TestModel_SingleLine_BlocksEnter(t *testing.T) {
	m := NewTextSnipModel(&TextSnipModelArgs{SingleLine: true})
	m.Model.SetValue("test")
	m.Model.Focus()

	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	// Value should not contain newline
	if m.lineCount() > 1 {
		t.Error("expected single-line model to block Enter")
	}
}

func TestModel_SingleLine_BlocksVerticalSelection(t *testing.T) {
	m := NewTextSnipModel(&TextSnipModelArgs{SingleLine: true})
	m.Model.SetValue("test content")
	m.Model.Focus()

	// SelectUp and SelectDown should be disabled
	km := m.SelectionKeyMap()
	if km.SelectUp.Enabled() {
		t.Error("expected SelectUp disabled in single-line mode")
	}
	if km.SelectDown.Enabled() {
		t.Error("expected SelectDown disabled in single-line mode")
	}
}

// --- Migration-sensitive tests (v1→v2 regression guards) ---

// TSL-SPACE: Guards msg.Type == tea.KeySpace in isTypingKey() (model.go:554)
// Space must be recognized as a typing key so it replaces active selection.
func TestModel_SpaceIsTypingKey(t *testing.T) {
	m := newTestModel()
	// Select "Hello"
	m = m.extendSelectionRight(5)
	if !m.HasSelection() {
		t.Fatal("expected selection before space")
	}

	// Space should delete selection then insert space
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeySpace, Text: " "})

	if m.HasSelection() {
		t.Error("expected selection cleared after space")
	}
	value := m.Value()
	if len(value) < 1 || value[0] != ' ' {
		t.Errorf("expected value to start with space, got %q", value)
	}
	// "Hello" (5 chars) replaced with " " (1 char), rest is " World\n..."
	if !strings.Contains(value, " World") {
		t.Errorf("expected remaining text after selection replacement, got %q", value)
	}
}

// TSL-ENTER: Guards msg.Type == tea.KeyEnter in isTypingKey() (model.go:559)
// Enter must be recognized as a typing key in multi-line mode so it replaces
// active selection with a newline.
func TestModel_EnterIsTypingKey_MultiLine(t *testing.T) {
	m := newTestModel()
	// Select "Hello"
	m = m.extendSelectionRight(5)
	if !m.HasSelection() {
		t.Fatal("expected selection before enter")
	}

	// Enter should delete selection then insert newline
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if m.HasSelection() {
		t.Error("expected selection cleared after enter")
	}
	// "Hello" should be replaced with a newline, so line count increases
	if m.lineCount() < 3 {
		t.Errorf("expected at least 3 lines after enter (was 3, 'Hello' replaced with newline), got %d", m.lineCount())
	}
}

// TSL-SHIFT: Guards strings.Contains(msg.String(), "shift") in isCursorMovement() (model.go:572)
// Shift+arrow must NOT clear selection — it extends it. This tests the inverse
// of TestModel_CursorMovementClearsSelection.
func TestModel_ShiftArrowExtendsSelection(t *testing.T) {
	m := newTestModel()
	// Create a selection via shift+right
	m = m.extendSelectionRight(3)
	if !m.HasSelection() {
		t.Fatal("expected selection active")
	}

	startSel := m.Selection()

	// Send another shift+right through Update() — should extend, not clear
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight, Mod: tea.ModShift})

	if !m.HasSelection() {
		t.Error("expected selection still active after shift+arrow")
	}
	sel := m.Selection()
	// Selection should have extended (End.Col increased by 1)
	if sel.End.Col <= startSel.End.Col {
		t.Errorf("expected selection extended (End.Col > %d), got %d",
			startSel.End.Col, sel.End.Col)
	}
}

// TSL-STYLE: Guards FocusedStyle.CursorLine and BlurredStyle.Base access in view.go
// When selection is active, view.go renders through renderFocused/renderBlurred which
// access these style fields. This test ensures both code paths execute and produce output.
func TestModel_View_FocusedStylePath(t *testing.T) {
	m := newTestModel()
	m = m.extendSelectionRight(5) // Select "Hello"

	// Focused path: m.Model.Focused() is true (set in newTestModel)
	if !m.Model.Focused() {
		t.Fatal("expected model to be focused")
	}
	focusedView := m.View()
	if focusedView.Content == "" {
		t.Fatal("expected non-empty focused view with selection")
	}
	if !strings.Contains(focusedView.Content, "World") {
		t.Error("expected focused view to contain non-selected text 'World'")
	}

	// Blurred path: after Blur()
	m.Model.Blur()
	if m.Model.Focused() {
		t.Fatal("expected model to be blurred after Blur()")
	}
	blurredView := m.View()
	if blurredView.Content == "" {
		t.Fatal("expected non-empty blurred view with selection")
	}
	if !strings.Contains(blurredView.Content, "World") {
		t.Error("expected blurred view to contain non-selected text 'World'")
	}
}

// --- Layer 2: View Tests ---

func TestModel_View_NoSelection(t *testing.T) {
	m := newTestModel()
	view := m.View()

	if view.Content == "" {
		t.Error("expected non-empty view")
	}
}

func TestModel_View_WithSelection(t *testing.T) {
	m := newTestModel()
	m = m.extendSelectionRight(5) // Select "Hello"

	view := m.View()
	if view.Content == "" {
		t.Error("expected non-empty view with selection")
	}
	// View should contain the original text (rendered differently with selection style)
}

func TestModel_View_MultiLineSelection(t *testing.T) {
	m := newTestModel()
	m = m.selectAll()

	view := m.View()
	if view.Content == "" {
		t.Error("expected non-empty view with multi-line selection")
	}
}
