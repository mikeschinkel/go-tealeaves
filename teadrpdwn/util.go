package teadrpdwn

import (
	"github.com/charmbracelet/x/ansi"
)

// maxLength returns visual display width of longest item text (ANSI-aware)
func maxLength(items []Option) (max int) {
	var item Option
	var width int

	for _, item = range items {
		width = ansi.StringWidth(item.Text)
		if width > max {
			max = width
		}
	}

	return max
}

// truncateWithEllipsis truncates string to maxWidth, adding ellipsis
func truncateWithEllipsis(text string, maxWidth int) (result string) {
	runes := []rune(text)
	if len(runes) <= maxWidth {
		return text
	}

	if maxWidth <= 1 {
		return "…"
	}

	return string(runes[:maxWidth-1]) + "…"
}
