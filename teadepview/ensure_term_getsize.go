package teadepview

import (
	"os"
	"strconv"
	"time"

	"github.com/mikeschinkel/go-dt"
	"golang.org/x/term"
)

// EnsureTermGetSize tries hard to get a real terminal size.
// It treats (0,0,nil) as "not ready yet" and retries briefly.
func EnsureTermGetSize(fd uintptr) (w int, h int, ok bool) {
	var deadline time.Time

	// Fast path: try immediately a few times with tiny sleeps.
	for {
		w, h, ok = getSizeOnce(fd)
		if ok {
			return
		}

		if deadline.IsZero() {
			deadline = time.Now().Add(250 * time.Millisecond)
		}

		if time.Now().After(deadline) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Fallbacks that often work better in IDEs.
	// 1) /dev/tty (controlling terminal) if available
	tty, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
	if err == nil {
		defer dt.CloseOrLog(tty)
		w, h, ok = getSizeOnce(tty.Fd())
		if ok {
			return
		}
	}

	// 2) Environment variables (some terminals/IDEs set these)
	w, h, ok = getSizeFromEnv()
	return
}

func getSizeOnce(fd uintptr) (w int, h int, ok bool) {
	w, h, err := term.GetSize(int(fd))
	if err != nil {
		return 0, 0, false
	}
	if w <= 0 || h <= 0 {
		return 0, 0, false
	}
	return w, h, true
}

func getSizeFromEnv() (w int, h int, ok bool) {
	cols := os.Getenv("COLUMNS")
	lines := os.Getenv("LINES")

	if cols == "" || lines == "" {
		return 0, 0, false
	}

	w64, err := strconv.Atoi(cols)
	if err != nil {
		return 0, 0, false
	}

	h64, err := strconv.Atoi(lines)
	if err != nil {
		return 0, 0, false
	}

	if w64 <= 0 || h64 <= 0 {
		return 0, 0, false
	}

	return w64, h64, true
}
