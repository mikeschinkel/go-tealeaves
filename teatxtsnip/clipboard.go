package teatxtsnip

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"golang.design/x/clipboard"
)

var clipboardInitialized bool

// initClipboard initializes the system clipboard (must be called once)
func initClipboard() bool {
	if clipboardInitialized {
		return true
	}

	err := clipboard.Init()
	if err == nil {
		clipboardInitialized = true
	}
	return clipboardInitialized
}

// Copy copies the selected text to the clipboard
func (m Model) Copy() Model {
	if !m.HasSelection() {
		return m
	}

	text := m.SelectedText()
	m = m.writeToClipboard(text)
	return m
}

// Cut copies the selected text to clipboard and deletes it
func (m Model) Cut() Model {
	if !m.HasSelection() {
		m.log("Cut: no selection, returning early")
		return m
	}

	text := m.SelectedText()
	m.log("Cut: selected text",
		"text", text,
		"text_len", len(text),
	)
	m = m.writeToClipboard(text)
	m.log("Cut: wrote to clipboard, now deleting selection")
	valueBefore := m.Value()
	m = m.deleteSelection()
	valueAfter := m.Value()
	m.log("Cut: deleteSelection completed",
		"value_before_len", len(valueBefore),
		"value_after_len", len(valueAfter),
		"value_changed", valueBefore != valueAfter,
	)
	return m
}

// Paste inserts text from clipboard, replacing any selection
func (m Model) Paste() (Model, tea.Cmd) {
	text := m.readFromClipboard()
	if text == "" {
		return m, nil
	}

	// Delete selection first if present
	if m.HasSelection() {
		m = m.deleteSelection()
	}

	// Insert the text by simulating typing
	// This handles newlines correctly through the textarea
	m.Model.InsertString(text)

	return m, nil
}

// SelectedText returns the currently selected text
func (m Model) SelectedText() string {
	if !m.HasSelection() {
		return ""
	}

	start, end := m.selection.Normalized()
	lines := m.lines()

	if len(lines) == 0 {
		return ""
	}

	// Single line selection
	if start.Row == end.Row {
		runes := []rune(lines[start.Row])
		endCol := end.Col
		if endCol > len(runes) {
			endCol = len(runes)
		}
		startCol := start.Col
		if startCol > len(runes) {
			startCol = len(runes)
		}
		if startCol > endCol {
			startCol, endCol = endCol, startCol
		}
		return string(runes[startCol:endCol])
	}

	// Multi-line selection
	var result strings.Builder

	// First line (from start.Col to end of line)
	if start.Row < len(lines) {
		runes := []rune(lines[start.Row])
		startCol := start.Col
		if startCol > len(runes) {
			startCol = len(runes)
		}
		result.WriteString(string(runes[startCol:]))
	}

	// Middle lines (complete lines)
	for row := start.Row + 1; row < end.Row && row < len(lines); row++ {
		result.WriteString("\n")
		result.WriteString(lines[row])
	}

	// Last line (from start of line to end.Col)
	if end.Row < len(lines) {
		result.WriteString("\n")
		runes := []rune(lines[end.Row])
		endCol := end.Col
		if endCol > len(runes) {
			endCol = len(runes)
		}
		result.WriteString(string(runes[:endCol]))
	}

	return result.String()
}

// deleteSelection deletes the selected text and clears the selection
func (m Model) deleteSelection() Model {
	if !m.HasSelection() {
		m.log("deleteSelection: no selection, returning early")
		return m
	}

	start, end := m.selection.Normalized()
	lines := m.lines()

	m.log("deleteSelection: starting",
		"start", start,
		"end", end,
		"num_lines", len(lines),
	)

	if len(lines) == 0 {
		m.selection = m.selection.Clear()
		return m
	}

	// Build new content without the selected text
	var result strings.Builder

	// Text before selection
	for row := 0; row < start.Row && row < len(lines); row++ {
		result.WriteString(lines[row])
		result.WriteString("\n")
	}

	// First line of selection (keep text before start.Col)
	if start.Row < len(lines) {
		runes := []rune(lines[start.Row])
		startCol := start.Col
		if startCol > len(runes) {
			startCol = len(runes)
		}
		result.WriteString(string(runes[:startCol]))
	}

	// Last line of selection (keep text after end.Col)
	if end.Row < len(lines) {
		runes := []rune(lines[end.Row])
		endCol := end.Col
		if endCol > len(runes) {
			endCol = len(runes)
		}
		result.WriteString(string(runes[endCol:]))
	}

	// Text after selection
	for row := end.Row + 1; row < len(lines); row++ {
		result.WriteString("\n")
		result.WriteString(lines[row])
	}

	newValue := result.String()
	m.log("deleteSelection: built new value",
		"new_value_len", len(newValue),
		"new_value_preview", truncateForLog(newValue, 50),
	)

	// Set new value and move cursor to start of selection
	m.Model.SetValue(newValue)
	m.log("deleteSelection: SetValue called, verifying",
		"actual_value_len", len(m.Model.Value()),
		"values_match", m.Model.Value() == newValue,
	)
	m.moveCursorTo(start.Row, start.Col)
	m.selection = m.selection.Clear()

	return m
}

// truncateForLog truncates a string for logging purposes
func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// writeToClipboard writes text to system clipboard with fallback to internal
func (m Model) writeToClipboard(text string) Model {
	if initClipboard() {
		clipboard.Write(clipboard.FmtText, []byte(text))
	}
	// Always update internal clipboard as fallback
	m.internalClip = text
	return m
}

// readFromClipboard reads text from system clipboard with fallback to internal
func (m Model) readFromClipboard() string {
	if initClipboard() {
		data := clipboard.Read(clipboard.FmtText)
		if len(data) > 0 {
			return string(data)
		}
	}
	// Fallback to internal clipboard
	return m.internalClip
}
