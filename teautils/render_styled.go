package teautils

import (
	"charm.land/lipgloss/v2"
)

// RenderAlignedLine renders text with a style and wraps it in alignment.
// This replaces the common 3-line pattern:
//
//	line = style.Render(text)
//	line = lipgloss.NewStyle().Width(width).AlignHorizontal(align).Render(line)
func RenderAlignedLine(text string, style lipgloss.Style, width int, align lipgloss.Position) string {
	styledText := style.Render(text)
	return lipgloss.NewStyle().
		Width(width).
		AlignHorizontal(align).
		Render(styledText)
}

// RenderCenteredLine renders text with a style and centers it within width.
// Convenience wrapper around RenderAlignedLine with lipgloss.Center.
func RenderCenteredLine(text string, style lipgloss.Style, width int) string {
	return RenderAlignedLine(text, style, width, lipgloss.Center)
}

// ApplyBoxBorder applies a standard border with padding to content.
// This replaces the common pattern: borderStyle.Padding(1, 2).Render(content)
func ApplyBoxBorder(borderStyle lipgloss.Style, content string) string {
	return borderStyle.Padding(1, 2).Render(content)
}
