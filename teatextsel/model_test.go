package teatextsel

import (
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel() Model {
	m := New()
	m.Model.SetWidth(80)
	m.Model.SetHeight(10)
	m.Model.SetValue("Hello World\nSecond line\nThird line")
	m.Model.Focus()
	// Move cursor to start
	m.moveCursorTo(0, 0)
	return m
}

// --- Layer 1: Model Tests ---

func TestNew(t *testing.T) {
	m := New()
	if m.HasSelection() {
		t.Error("expected no selection initially")
	}
	if m.IsSingleLine() {
		t.Error("expected IsSingleLine=false for New()")
	}
}

func TestNewSingleLine(t *testing.T) {
	m := NewSingleLine()
	if !m.IsSingleLine() {
		t.Error("expected IsSingleLine=true")
	}
}

func TestNewFromTextarea(t *testing.T) {
	ta := textarea.New()
	ta.SetValue("existing content")
	m := NewFromTextarea(ta)

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
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
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
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'X'}})

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
	m := NewSingleLine()
	m.Model.SetValue("test")
	m.Model.Focus()

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Value should not contain newline
	if m.lineCount() > 1 {
		t.Error("expected single-line model to block Enter")
	}
}

func TestModel_SingleLine_BlocksVerticalSelection(t *testing.T) {
	m := NewSingleLine()
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

// --- Layer 2: View Tests ---

func TestModel_View_NoSelection(t *testing.T) {
	m := newTestModel()
	view := m.View()

	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestModel_View_WithSelection(t *testing.T) {
	m := newTestModel()
	m = m.extendSelectionRight(5) // Select "Hello"

	view := m.View()
	if view == "" {
		t.Error("expected non-empty view with selection")
	}
	// View should contain the original text (rendered differently with selection style)
}

func TestModel_View_MultiLineSelection(t *testing.T) {
	m := newTestModel()
	m = m.selectAll()

	view := m.View()
	if view == "" {
		t.Error("expected non-empty view with multi-line selection")
	}
}
