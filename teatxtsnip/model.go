package teatxtsnip

import (
	"log/slog"
	"strings"
	"unicode"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
)

// Model wraps textarea.Model with text selection and clipboard support
type Model struct {
	textarea.Model // Embedded textarea

	// Selection state
	selection    Selection
	selectionKey SelectionKeyMap

	// Internal clipboard fallback (used when system clipboard unavailable)
	internalClip string

	// Single-line mode (prevents newlines, used for single-line inputs)
	singleLine bool

	// Optional logger for debugging
	Logger *slog.Logger
}

// New creates a new Model wrapping a textarea.Model
func New() Model {
	ta := textarea.New()
	ta.Prompt = "" // Remove the default "> " prompt
	return Model{
		Model:        ta,
		selection:    NewSelection(),
		selectionKey: DefaultSelectionKeyMap(),
	}
}

// NewSingleLine creates a single-line input with selection and clipboard support.
// This is an alternative to textinput.Model that supports Ctrl+C/X/V clipboard operations.
// The textarea is configured to prevent newlines and display as a single line.
func NewSingleLine() Model {
	ta := textarea.New()
	ta.SetHeight(1)
	ta.ShowLineNumbers = false
	ta.Prompt = "" // Remove the default "> " prompt
	// Disable soft wrap to behave like a single-line input
	ta.SetWidth(80) // Default width, can be overridden with SetWidth()

	// Create a keymap that doesn't include multi-line selection operations
	km := DefaultSelectionKeyMap()
	// Disable multi-line selection keys for single-line input
	km.SelectUp.SetEnabled(false)
	km.SelectDown.SetEnabled(false)
	km.SelectToStart.SetEnabled(false)
	km.SelectToEnd.SetEnabled(false)

	return Model{
		Model:        ta,
		selection:    NewSelection(),
		selectionKey: km,
		singleLine:   true,
	}
}

// IsSingleLine returns true if this is a single-line input
func (m Model) IsSingleLine() bool {
	return m.singleLine
}

// NewFromTextarea creates a Model from an existing textarea.Model
func NewFromTextarea(ta textarea.Model) Model {
	return Model{
		Model:        ta,
		selection:    NewSelection(),
		selectionKey: DefaultSelectionKeyMap(),
	}
}

// SelectionKeyMap returns the current selection key map
func (m Model) SelectionKeyMap() SelectionKeyMap {
	return m.selectionKey
}

// SetSelectionKeyMap sets the selection key map
func (m Model) SetSelectionKeyMap(km SelectionKeyMap) Model {
	m.selectionKey = km
	return m
}

// Selection returns the current selection state
func (m Model) Selection() Selection {
	return m.selection
}

// HasSelection returns true if there is an active non-empty selection
func (m Model) HasSelection() bool {
	return m.selection.Active && !m.selection.IsEmpty()
}

// ClearSelection clears the current selection
func (m Model) ClearSelection() Model {
	m.selection = m.selection.Clear()
	return m
}

// SetSelection sets the selection state (for saving/restoring between contexts)
func (m Model) SetSelection(sel Selection) Model {
	m.selection = sel
	return m
}

// SetLogger sets an optional logger for debugging
func (m Model) SetLogger(logger *slog.Logger) Model {
	m.Logger = logger
	return m
}

// log logs a debug message if logger is set
func (m Model) log(msg string, args ...any) {
	if m.Logger != nil {
		m.Logger.Debug(msg, args...)
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	keyMsg, isKey := msg.(tea.KeyPressMsg)
	if !isKey {
		// Not a key message - pass to textarea
		m.Model, cmd = m.Model.Update(msg)
		return m, cmd
	}

	// Handle selection keys
	switch {
	case key.Matches(keyMsg, m.selectionKey.SelectAll):
		m = m.selectAll()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.Copy):
		m = m.Copy()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.Cut):
		m.log("Cut key matched",
			"key", keyMsg.String(),
			"hasSelection", m.HasSelection(),
			"selection", m.selection,
		)
		m = m.Cut()
		m.log("Cut completed",
			"hasSelection", m.HasSelection(),
			"value_len", len(m.Value()),
		)
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.Paste):
		m, cmd = m.Paste()
		return m, cmd

	case key.Matches(keyMsg, m.selectionKey.ClearSelection):
		if m.HasSelection() {
			m = m.ClearSelection()
			return m, nil
		}
		// Let escape propagate if no selection

	case key.Matches(keyMsg, m.selectionKey.SelectLeft):
		m = m.extendSelectionLeft(1)
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectRight):
		m = m.extendSelectionRight(1)
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectUp):
		m = m.extendSelectionUp()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectDown):
		m = m.extendSelectionDown()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectWordLeft):
		m = m.extendSelectionWordLeft()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectWordRight):
		m = m.extendSelectionWordRight()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectToLineStart):
		m = m.extendSelectionToLineStart()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectToLineEnd):
		m = m.extendSelectionToLineEnd()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectToStart):
		m = m.extendSelectionToStart()
		return m, nil

	case key.Matches(keyMsg, m.selectionKey.SelectToEnd):
		m = m.extendSelectionToEnd()
		return m, nil
	}

	// In single-line mode, block Enter key to prevent newlines
	if m.singleLine && keyMsg.Code == tea.KeyEnter {
		return m, nil
	}

	// Check if this is a typing key that should replace selection
	if m.HasSelection() && isTypingKey(keyMsg) {
		m = m.deleteSelection()
		// Fall through to let textarea handle the typed character
	}

	// Clear selection on cursor movement (non-shift arrow keys)
	if isCursorMovement(keyMsg) && m.HasSelection() {
		m = m.ClearSelection()
	}

	// Pass to textarea
	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// cursorPosition returns the current cursor position
func (m Model) cursorPosition() Position {
	return Position{
		Row: m.Model.Line(),
		Col: m.Model.LineInfo().CharOffset,
	}
}

// lines returns the textarea content as a slice of lines
func (m Model) lines() []string {
	return strings.Split(m.Model.Value(), "\n")
}

// lineRunes returns the runes for a specific line
func (m Model) lineRunes(row int) []rune {
	lines := m.lines()
	if row < 0 || row >= len(lines) {
		return nil
	}
	return []rune(lines[row])
}

// lineLength returns the length (in runes) of the given line
func (m Model) lineLength(row int) int {
	return len(m.lineRunes(row))
}

// lineCount returns the number of lines
func (m Model) lineCount() int {
	return len(m.lines())
}

// selectAll selects all text
func (m Model) selectAll() Model {
	m.selection = SelectAll(m.lines())
	return m
}

// extendSelectionLeft extends selection left by n characters
func (m Model) extendSelectionLeft(n int) Model {
	pos := m.cursorPosition()

	// Start selection if not active
	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	// Calculate new position
	newCol := pos.Col - n
	newRow := pos.Row

	if newCol < 0 {
		// Wrap to previous line
		if newRow > 0 {
			newRow--
			newCol = m.lineLength(newRow)
		} else {
			newCol = 0
		}
	}

	m.selection = m.selection.Extend(Position{Row: newRow, Col: newCol})

	// Move cursor to match selection end
	m.moveCursorTo(newRow, newCol)

	return m
}

// extendSelectionRight extends selection right by n characters
func (m Model) extendSelectionRight(n int) Model {
	pos := m.cursorPosition()

	// Start selection if not active
	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	lineLen := m.lineLength(pos.Row)
	newCol := pos.Col + n
	newRow := pos.Row

	if newCol > lineLen {
		// Wrap to next line
		if newRow < m.lineCount()-1 {
			newRow++
			newCol = 0
		} else {
			newCol = lineLen
		}
	}

	m.selection = m.selection.Extend(Position{Row: newRow, Col: newCol})

	// Move cursor to match selection end
	m.moveCursorTo(newRow, newCol)

	return m
}

// extendSelectionUp extends selection up one line
func (m Model) extendSelectionUp() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	newRow := pos.Row - 1
	if newRow < 0 {
		newRow = 0
	}

	newCol := pos.Col
	lineLen := m.lineLength(newRow)
	if newCol > lineLen {
		newCol = lineLen
	}

	m.selection = m.selection.Extend(Position{Row: newRow, Col: newCol})
	m.moveCursorTo(newRow, newCol)

	return m
}

// extendSelectionDown extends selection down one line
func (m Model) extendSelectionDown() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	newRow := pos.Row + 1
	lastRow := m.lineCount() - 1
	if newRow > lastRow {
		newRow = lastRow
	}

	newCol := pos.Col
	lineLen := m.lineLength(newRow)
	if newCol > lineLen {
		newCol = lineLen
	}

	m.selection = m.selection.Extend(Position{Row: newRow, Col: newCol})
	m.moveCursorTo(newRow, newCol)

	return m
}

// extendSelectionWordLeft extends selection to the previous word boundary
func (m Model) extendSelectionWordLeft() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	// Find word boundary
	runes := m.lineRunes(pos.Row)
	col := pos.Col
	row := pos.Row

	if col == 0 && row > 0 {
		// At line start, go to end of previous line
		row--
		col = m.lineLength(row)
		runes = m.lineRunes(row)
	}

	if len(runes) > 0 && col > 0 {
		// Skip whitespace
		for col > 0 && unicode.IsSpace(runes[col-1]) {
			col--
		}
		// Skip word characters
		for col > 0 && !unicode.IsSpace(runes[col-1]) {
			col--
		}
	}

	m.selection = m.selection.Extend(Position{Row: row, Col: col})
	m.moveCursorTo(row, col)

	return m
}

// extendSelectionWordRight extends selection to the next word boundary
func (m Model) extendSelectionWordRight() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	runes := m.lineRunes(pos.Row)
	col := pos.Col
	row := pos.Row
	lineLen := len(runes)

	if col >= lineLen && row < m.lineCount()-1 {
		// At line end, go to start of next line
		row++
		col = 0
		runes = m.lineRunes(row)
		lineLen = len(runes)
	}

	if lineLen > 0 && col < lineLen {
		// Skip word characters
		for col < lineLen && !unicode.IsSpace(runes[col]) {
			col++
		}
		// Skip whitespace
		for col < lineLen && unicode.IsSpace(runes[col]) {
			col++
		}
	}

	m.selection = m.selection.Extend(Position{Row: row, Col: col})
	m.moveCursorTo(row, col)

	return m
}

// extendSelectionToLineStart extends selection to start of current line
func (m Model) extendSelectionToLineStart() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	m.selection = m.selection.Extend(Position{Row: pos.Row, Col: 0})
	m.moveCursorTo(pos.Row, 0)

	return m
}

// extendSelectionToLineEnd extends selection to end of current line
func (m Model) extendSelectionToLineEnd() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	lineLen := m.lineLength(pos.Row)
	m.selection = m.selection.Extend(Position{Row: pos.Row, Col: lineLen})
	m.moveCursorTo(pos.Row, lineLen)

	return m
}

// extendSelectionToStart extends selection to start of document
func (m Model) extendSelectionToStart() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	m.selection = m.selection.Extend(Position{Row: 0, Col: 0})
	m.moveCursorTo(0, 0)

	return m
}

// extendSelectionToEnd extends selection to end of document
func (m Model) extendSelectionToEnd() Model {
	pos := m.cursorPosition()

	if !m.selection.Active {
		m.selection = m.selection.Begin(pos)
	}

	lastRow := m.lineCount() - 1
	lastCol := m.lineLength(lastRow)
	m.selection = m.selection.Extend(Position{Row: lastRow, Col: lastCol})
	m.moveCursorTo(lastRow, lastCol)

	return m
}

// moveCursorTo moves the textarea cursor to the specified position
func (m *Model) moveCursorTo(row, col int) {
	// Move to correct row
	currentRow := m.Model.Line()
	for currentRow < row && currentRow < m.lineCount()-1 {
		m.Model.CursorDown()
		currentRow = m.Model.Line()
	}
	for currentRow > row && currentRow > 0 {
		m.Model.CursorUp()
		currentRow = m.Model.Line()
	}

	// Move to correct column
	m.Model.SetCursorColumn(col)
}

// isTypingKey returns true if the key would insert text
func isTypingKey(msg tea.KeyPressMsg) bool {
	// Text input (runes, space) are typing keys
	if msg.Text != "" {
		return true
	}

	// Enter is a typing key (inserts newline)
	if msg.Code == tea.KeyEnter {
		return true
	}

	return false
}

// isCursorMovement returns true if the key moves the cursor without shift
func isCursorMovement(msg tea.KeyPressMsg) bool {
	switch msg.Code {
	case tea.KeyLeft, tea.KeyRight, tea.KeyUp, tea.KeyDown,
		tea.KeyHome, tea.KeyEnd, tea.KeyPgUp, tea.KeyPgDown:
		// Only clear selection if not extending (shift not held)
		return !msg.Mod.Contains(tea.ModShift)
	}
	return false
}
