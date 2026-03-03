package teagrid

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// limitStr truncates a string to maxLen printable characters,
// handling ANSI escape sequences and newlines correctly.
func limitStr(str string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}

	// Replace newlines with ellipsis
	idx := strings.Index(str, "\n")
	if idx > -1 {
		str = str[:idx] + "…"
	}

	if ansi.StringWidth(str) > maxLen {
		return ansi.Truncate(str, maxLen, "…")
	}

	return str
}
