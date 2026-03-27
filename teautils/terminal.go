package teautils

import (
	"os"
	"strings"
)

// IsJediTerm returns true when running inside JetBrains' JediTerm terminal
// emulator (GoLand, IntelliJ, etc.). JediTerm mishandles CBT escape
// sequences that ultraviolet emits under normal TERM values (CBT means
// Cursor Backward Tab (CSI Z) and it that moves cursor back to previous
// tab stop.
//
// Checks TERMINAL_EMULATOR (set in JediTerm's built-in terminal) and falls
// back to __CFBundleIdentifier and XPC_SERVICE_NAME (set by macOS for
// Run/Debug configurations where TERMINAL_EMULATOR is not propagated).
func IsJediTerm() (isJedi bool) {
	isJedi = os.Getenv("TERMINAL_EMULATOR") == "JetBrains-JediTerm"
	if isJedi {
		goto end
	}
	isJedi = strings.HasPrefix(os.Getenv("__CFBundleIdentifier"), "com.jetbrains.")
	if isJedi {
		goto end
	}
	isJedi = strings.Contains(os.Getenv("XPC_SERVICE_NAME"), "com.jetbrains.")
	if isJedi {
		goto end
	}
end:
	return
}
