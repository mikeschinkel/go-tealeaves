package teatextsel

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// SelectionStyle is the default style for selected text (inverted colors)
var SelectionStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("39")). // Light blue background
	Foreground(lipgloss.Color("232")) // Dark text

// View renders the textarea with selection highlighting
func (m Model) View() string {
	if !m.HasSelection() {
		// No selection - use default textarea view
		return m.Model.View()
	}

	// With selection, we need to render with highlighting
	return m.renderWithSelection()
}

// renderWithSelection renders the textarea content with selection highlighting
func (m Model) renderWithSelection() string {
	// Get the textarea's raw view
	// We need to work with the textarea's internal rendering, but add selection highlighting
	// Since textarea.Model doesn't expose its internal rendering details,
	// we'll render the content ourselves when selection is active

	lines := m.lines()
	if len(lines) == 0 {
		return m.Model.View()
	}

	start, end := m.selection.Normalized()
	var result strings.Builder

	// Get viewport dimensions from textarea
	width := m.Model.Width()
	height := m.Model.Height()

	// Calculate visible line range (simplified - textarea handles scrolling)
	// For now, render all lines and let textarea handle viewport

	for row, line := range lines {
		if row > 0 {
			result.WriteString("\n")
		}

		runes := []rune(line)
		lineLen := len(runes)

		// Check if this line has any selection
		lineStart := Position{Row: row, Col: 0}
		lineEnd := Position{Row: row, Col: lineLen}

		if !m.selection.Contains(lineStart) && !m.selection.Contains(lineEnd) &&
			(row < start.Row || row > end.Row) {
			// Line is not in selection at all
			result.WriteString(line)
			continue
		}

		// Line has some selection - render character by character
		for col := 0; col <= lineLen; col++ {
			pos := Position{Row: row, Col: col}

			// Determine if this character is selected
			isSelected := m.isPositionSelected(pos, start, end)

			if col < lineLen {
				char := string(runes[col])
				if isSelected {
					result.WriteString(SelectionStyle.Render(char))
				} else {
					result.WriteString(char)
				}
			}
		}
	}

	// Pad to width/height if needed
	rendered := result.String()
	renderedLines := strings.Split(rendered, "\n")

	// Pad lines to width
	var padded strings.Builder
	for i, line := range renderedLines {
		if i > 0 {
			padded.WriteString("\n")
		}
		padded.WriteString(line)
		// Don't pad width - let textarea handle it
	}

	// Pad to height
	for len(renderedLines) < height {
		padded.WriteString("\n")
		renderedLines = append(renderedLines, "")
	}

	// Apply textarea's focus/blur styling
	content := padded.String()

	// Use textarea's base style
	if m.Model.Focused() {
		return m.renderFocused(content, width, height)
	}
	return m.renderBlurred(content, width, height)
}

// isPositionSelected returns true if the given position is within the selection
func (m Model) isPositionSelected(pos Position, start, end Position) bool {
	// Before start
	if pos.Row < start.Row {
		return false
	}
	if pos.Row == start.Row && pos.Col < start.Col {
		return false
	}

	// After end
	if pos.Row > end.Row {
		return false
	}
	if pos.Row == end.Row && pos.Col >= end.Col {
		return false
	}

	return true
}

// renderFocused applies focused styling to the content
func (m Model) renderFocused(content string, width, height int) string {
	style := m.Model.FocusedStyle

	// Apply prompt and line styling
	lines := strings.Split(content, "\n")
	var result strings.Builder

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}

		// Apply cursor line style if this is the current line
		if i == m.Model.Line() {
			result.WriteString(style.CursorLine.Render(line))
		} else {
			result.WriteString(line)
		}
	}

	// Apply base style
	return style.Base.Render(result.String())
}

// renderBlurred applies blurred styling to the content
func (m Model) renderBlurred(content string, width, height int) string {
	style := m.Model.BlurredStyle
	return style.Base.Render(content)
}
